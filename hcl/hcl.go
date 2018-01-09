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
	"github.com/fatih/color"
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
func Marshal(value interface{}) ([]byte, error) {
	value = toBase(value)
	return marshalHCL(value, "", ""), nil
}

// MarshalIndent serialize values to hcl format with indentation
func MarshalIndent(value interface{}, prefix, indent string) ([]byte, error) {
	fmt.Println(1, value)
	value = toBase(value)
	fmt.Println(2, value)
	return marshalHCL(value, prefix, indent), nil
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
		return result

	case reflect.Map:
		result := make(map[string]interface{}, val.Len())
		for _, key := range val.MapKeys() {
			result[fmt.Sprintf("%v", key)] = toBase(val.MapIndex(key).Interface())
		}
		return result

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

func marshalHCL(value interface{}, prefix, indent string) []byte {
	if value == nil {
		return []byte("null")
	}

	ifIndent := func(vTrue, vFalse interface{}) interface{} {
		if indent != "" {
			return vTrue
		}
		return vFalse
	}

	switch value := value.(type) {
	case []interface{}:
		result := make([]string, len(value))
		var (
			totalLength int
			newLine     bool
		)
		for i := range value {
			result[i] = string(marshalHCL(value[i], "", indent))
			totalLength += len(result[i])
			newLine = newLine || strings.Contains(result[i], "\n")
		}

		if totalLength > 60 && indent != "" || newLine {
			output := ""
			for i := range result {
				output += "\n"
				for _, s := range strings.Split(result[i], "\n") {
					output += fmt.Sprintf("%s%s", prefix, s)
				}
				output += ","
			}
			fmt.Println("array 1", color.GreenString(fmt.Sprintf("[%s\n]", output)))
			return []byte(fmt.Sprintf("[%s\n]", output))
		}
		fmt.Println("array 2", color.GreenString(fmt.Sprintf("[%s]", strings.Join(result, ifIndent(", ", ",").(string)))))
		return []byte(fmt.Sprintf("[%s]", strings.Join(result, ifIndent(", ", ",").(string))))

	case map[string]interface{}:
		keys := make([]string, 0, len(value))
		rendered := make(map[string]string, len(value))
		keyLen := 0

		for key, val := range value {
			keys = append(keys, key)
			if len(key) > keyLen && indent != "" {
				keyLen = len(key)
			}
			rendered[key] = string(marshalHCL(val, indent, indent))
		}
		sort.Strings(keys)

		items := make([]string, 0, len(value))

		if indent == "" {
			for _, key := range keys {
				switch value[key].(type) {
				case map[string]interface{}:
					items = append(items, fmt.Sprintf("%s{%s}", id(key), rendered[key]))
				default:
					items = append(items, fmt.Sprintf("%s=%s", id(key), rendered[key]))
				}
			}
			fmt.Println("map 1", color.GreenString(strings.Join(items, " ")))
			return []byte(strings.Join(items, " "))
		}

		for _, multiline := range []bool{false, true} {
			for _, key := range keys {
				lines := strings.Split(rendered[key], "\n")

				// We process the multilines elements after the single line one
				if len(lines) > 1 && !multiline || len(lines) == 1 && multiline {
					continue
				}

				_, isMap := value[key].(map[string]interface{})
				if !multiline {
					if isMap {
						lines[0] = fmt.Sprintf("{ %s }", lines[0])
					}
					fmt.Println(color.MagentaString(lines[0]))
					items = append(items, fmt.Sprintf("%*s = %s", -keyLen, id(key), lines[0]))
					continue
				} else if len(items) > 0 {
					items = append(items, "")
				}

				if isMap {
					items = append(items, fmt.Sprintf("%s {", id(key)))
					for _, line := range lines {
						items = append(items, fmt.Sprintf("%s%s", prefix+indent, line))
					}
					items = append(items, "}")
				} else {
					items = append(items, fmt.Sprintf("%*s = %s", -keyLen, id(key), lines[0]))
					for _, line := range lines[1:] {
						items = append(items, fmt.Sprintf("%s%s", prefix+indent, line))
					}
				}
			}
		}

		fmt.Println("map 2", color.GreenString(strings.Join(items, ifIndent("\n", " ").(string))))
		return []byte(strings.Join(items, ifIndent("\n", " ").(string)))

	case string:
		if indent != "" && strings.Contains(value, "\\n") {
			prefix += indent

			// We unquote the value, other escape chars are \a, \b, \f and \v are not managed
			value = value[1 : len(value)-1]
			value = strings.Replace(value, `\\`, "\\", -1)
			value = strings.Replace(value, `\"`, "\"", -1)
			value = strings.Replace(value, `\r`, "\r", -1)
			value = strings.Replace(value, `\t`, "\t", -1)

			// We indent each line
			lines := strings.Join(strings.Split(value, "\\n"), "\n"+prefix)
			value = fmt.Sprintf("<<-EOF\n%[1]s%[2]s\n%[1]sEOF", prefix, lines)
		}
		fmt.Println("string", color.GreenString(value))
		return []byte(value)
	}
	return []byte(fmt.Sprintf("Not evaluated %T", value))
}

var identifierRegex = regexp.MustCompile(`^[A-za-z][\w-]*$`)

func id(key string) string {
	if identifierRegex.MatchString(key) {
		return key
	}
	// The identifier contains characters that may be considered invalid, we have to quote it
	return fmt.Sprintf("%q", key)
}
