package hcl

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/coveo/gotemplate/utils"
	"github.com/hashicorp/hcl"
)

// Expose hcl public objects
var (
	Decode       = hcl.Decode
	DecodeObject = hcl.DecodeObject
	Parse        = hcl.Parse
	ParseBytes   = hcl.ParseBytes
	ParseString  = hcl.ParseString
	Unmarshal    = hcl.Unmarshal
)

// Load loads hcl file into variable
func Load(filename string) (result map[string]interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		if err = Unmarshal(content, &result); err == nil {
			result = Flatten(result)
		}
	}
	return
}

// Marshal serialize values to hcl format
func Marshal(value interface{}) ([]byte, error) { return MarshalIndent(value, "", "") }

// MarshalIndent serialize values to hcl format with indentation
func MarshalIndent(value interface{}, prefix, indent string) ([]byte, error) {
	value = toBase(value)
	result, err := marshalHCL(value, true, prefix, indent)
	return []byte(result), err
}

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
}

// Flatten - Convert array of map to single map if there is only one element in the array
// By default, Unmarshal returns array of map even if there is only a single map in the definition
func Flatten(source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		switch value := value.(type) {
		case []map[string]interface{}:
			switch len(value) {
			case 1:
				source[key] = Flatten(value[0])
			default:
				for i, subMap := range value {
					value[i] = Flatten(subMap)
				}
			}
		}
	}
	return source
}

func toBase(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	typ, val := reflect.TypeOf(value), reflect.ValueOf(value)
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = val.Type()
	}
	switch typ.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return fmt.Sprintf("%v", value)

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, val.Len())
		for i := range result {
			result[i] = toBase(reflect.ValueOf(value).Index(i).Interface())
		}
		if len(result) == 1 && reflect.TypeOf(result[0]).Kind() == reflect.Map {
			// If the result is an array of one map, we just return the inner element
			return result[0]
		}
		return result

	case reflect.Map:
		result := make(map[string]interface{}, val.Len())
		for _, key := range val.MapKeys() {
			result[fmt.Sprintf("%v", key)] = toBase(val.MapIndex(key).Interface())
		}
		return Flatten(result)

	case reflect.Struct:
		result := make(map[string]interface{}, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			sf := typ.Field(i)
			if sf.Anonymous {
				t := sf.Type
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				// If embedded, StructField.PkgPath is not a reliable
				// indicator of whether the field is exported.
				// See https://golang.org/issue/21122
				if !utils.IsExported(t.Name()) && t.Kind() != reflect.Struct {
					// Ignore embedded fields of unexported non-struct types.
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
					continue
				}
			} else if sf.PkgPath != "" {
				// Ignore unexported non-embedded fields.
				continue
			}
			tag := sf.Tag.Get("hcl")
			if tag == "" {
				// If there is no hcl specific tag, we rely on json tag if there is
				tag = sf.Tag.Get("json")
			}
			if tag == "-" {
				continue
			}
			if tag == "" {
				tag = sf.Name
			}

			split := strings.Split(tag, ",")
			name := split[0]
			options := make(map[string]bool, len(split[1:]))
			for i := range split[1:] {
				options[split[i+1]] = true
			}

			if options["omitempty"] && utils.IsEmptyValue(val.Field(i)) {
				continue
			}

			result[name] = toBase(val.Field(i).Interface())
		}
		return result
	default:
		fmt.Fprintf(os.Stderr, "Unknown type %T %v : %v\n", value, typ.Kind(), value)
		return fmt.Sprintf("%v", value)
	}
}

func marshalHCL(value interface{}, head bool, prefix, indent string) (result string, err error) {
	//defer func() { fmt.Println(color.GreenString(result)) }()
	if value == nil {
		result = "null"
		return
	}

	ifIndent := func(vTrue, vFalse interface{}) interface{} { return utils.IIf(indent, vTrue, vFalse) }
	const specialFormat = "#HCL_ARRAY_MAP#!"

	switch value := value.(type) {
	case string:
		if indent != "" && strings.Contains(value, "\\n") {
			// We unquote the value
			unIndented := value[1 : len(value)-1]
			// Then replace escaped characters, other escape chars are \a, \b, \f and \v are not managed
			unIndented = strings.Replace(unIndented, `\n`, "\n", -1)
			unIndented = strings.Replace(unIndented, `\\`, "\\", -1)
			unIndented = strings.Replace(unIndented, `\"`, "\"", -1)
			unIndented = strings.Replace(unIndented, `\r`, "\r", -1)
			unIndented = strings.Replace(unIndented, `\t`, "\t", -1)
			unIndented = utils.UnIndent(unIndented)
			if strings.HasSuffix(unIndented, "\n") {
				value = fmt.Sprintf("<<-EOF\n%sEOF", unIndented)
			}
		}
		result = value

	case []interface{}:
		results := make([]string, len(value))
		if isArrayOfMap(value) {
			for i, element := range value {
				element := element.(map[string]interface{})
				for key := range element {
					results[i], err = marshalHCL(element[key], false, "", indent)
					results[i] = fmt.Sprintf(`%s "%s" %s`, specialFormat, key, results[i])
				}
			}
			result = strings.Join(results, ifIndent("\n\n", " ").(string))
			break
		}
		var totalLength int
		var newLine bool
		for i := range value {
			results[i], err = marshalHCL(value[i], false, "", indent)
			totalLength += len(results[i])
			newLine = newLine || strings.Contains(results[i], "\n")
		}
		if totalLength > 60 && indent != "" || newLine {
			result = fmt.Sprintf("[\n%s,\n]", utils.Indent(strings.Join(results, ",\n"), prefix+indent))
		} else {
			result = fmt.Sprintf("[%s]", strings.Join(results, ifIndent(", ", ",").(string)))
		}

	case map[string]interface{}:
		keys := make([]string, 0, len(value))
		rendered := make(map[string]string, len(value))
		keyLen := 0

		for key, val := range value {
			keys = append(keys, key)
			if len(key) > keyLen && indent != "" {
				keyLen = len(key)
			}
			rendered[key], err = marshalHCL(val, false, "", indent)
		}
		sort.Strings(keys)

		items := make([]string, 0, len(value)+2)
		for _, multiline := range []bool{false, true} {
			for _, key := range keys {
				rendered := rendered[key]
				lines := strings.Split(rendered, "\n")

				// We process the multilines elements after the single line one
				if len(lines) > 1 && !multiline || len(lines) == 1 && multiline {
					continue
				}

				if multiline && len(items) > 0 {
					// Add a blank line between multilines elements
					items = append(items, "")
					keyLen = 0
				}

				equal := ifIndent(" = ", "=").(string)
				if _, isMap := value[key].(map[string]interface{}); isMap {
					if multiline {
						equal = " "
					} else if indent == "" {
						equal = ""
					}
				}

				if strings.Contains(rendered, specialFormat) {
					items = append(items, strings.Replace(rendered, specialFormat, id(key), -1))

				} else {
					items = append(items, fmt.Sprintf("%*s%s%s", -keyLen, id(key), equal, rendered))
				}
			}
		}

		if head {
			result = strings.Join(items, ifIndent("\n", " ").(string))
			break
		}

		if indent == "" || len(items) == 0 {
			result = fmt.Sprintf("{%s}", strings.Join(items, " "))
			break
		}

		result = fmt.Sprintf("{\n%s\n}", utils.Indent(strings.Join(items, "\n"), prefix+indent))

	default:
		err = fmt.Errorf("Unknown type %[1]T %[1]v", value)
	}
	return
}

func isArrayOfMap(array []interface{}) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		mapItem, isMap := item.(map[string]interface{})
		if !isMap || len(mapItem) != 1 {
			return false
		}
	}
	return true
}

var identifierRegex = regexp.MustCompile(`^[A-za-z][\w-]*$`)

func id(key string) string {
	if identifierRegex.MatchString(key) {
		return key
	}
	// The identifier contains characters that may be considered invalid, we have to quote it
	return fmt.Sprintf("%q", key)
}
