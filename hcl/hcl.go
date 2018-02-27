package hcl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/coveo/gotemplate/errors"
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
)

var _ = func() int {
	utils.HCLConvert = Unmarshal
	return 0
}()

// flatten converts array of map to single map if there is only one element in the array
// By default, hc.Unmarshal returns array of map even if there is only a single map in the definition
func flatten(source interface{}) interface{} {
	switch value := source.(type) {
	case []interface{}:
		for i, sub := range value {
			value[i] = flatten(sub)
		}
	case map[string]interface{}:
		for key := range value {
			value[key] = flatten(value[key])
		}
	case []map[string]interface{}:
		switch len(value) {
		case 1:
			source = flatten(value[0])
		default:
			for i := range value {
				value[i] = flatten(value[i]).(map[string]interface{})
			}
		}
	}
	return source
}

// Unmarshal adds support to single array and struct representation
func Unmarshal(bs []byte, out interface{}) (err error) {
	defer func() { err = errors.Trap(err, recover()) }()
	bs = bytes.TrimSpace(bs)

	if err = hcl.Unmarshal(bs, out); err != nil {
		bs = append([]byte("_="), bs...)
		var temp map[string]interface{}
		if err = hcl.Unmarshal(bs, &temp); err != nil {
			return err
		}
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(temp["_"]))
	}
	result := flatten(reflect.ValueOf(out).Elem().Interface())
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
	return
}

// Load loads hcl file into variable
func Load(filename string) (result interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		err = Unmarshal(content, &result)
	}
	return
}

// Marshal serialize values to hcl format
func Marshal(value interface{}) ([]byte, error) { return MarshalIndent(value, "", "") }

// MarshalIndent serialize values to hcl format with indentation
func MarshalIndent(value interface{}, prefix, indent string) ([]byte, error) {
	result, err := marshalHCL(utils.ToNativeRepresentation(value), true, true, prefix, indent)
	return []byte(result), err
}

// MarshalTFVars serialize values to hcl format (without hcl map format)
func MarshalTFVars(value interface{}) ([]byte, error) { return MarshalTFVarsIndent(value, "", "") }

// MarshalTFVarsIndent serialize values to hcl format with indentation (without hcl map format)
func MarshalTFVarsIndent(value interface{}, prefix, indent string) ([]byte, error) {
	result, err := marshalHCL(utils.ToNativeRepresentation(value), false, true, prefix, indent)
	return []byte(result), err
}

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
}

func marshalHCL(value interface{}, fullHcl, head bool, prefix, indent string) (result string, err error) {
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
		if fullHcl && isArrayOfMap(value) {
			for i, element := range value {
				element := element.(map[string]interface{})
				for key := range element {
					if results[i], err = marshalHCL(element[key], fullHcl, false, "", indent); err != nil {
						return
					}
					if head {
						results[i] = fmt.Sprintf(`%s%s%s`, id(key), ifIndent(" = ", ""), results[i])
					} else {
						results[i] = fmt.Sprintf(`%s %s %s`, specialFormat, id(key), results[i])
					}
				}
			}
			result = strings.Join(results, ifIndent("\n\n", " ").(string))
			break
		}
		var totalLength int
		var newLine bool
		for i := range value {
			if results[i], err = marshalHCL(value[i], fullHcl, false, "", indent); err != nil {
				return
			}
			totalLength += len(results[i])
			newLine = newLine || strings.Contains(results[i], "\n")
		}
		if totalLength > 60 && indent != "" || newLine {
			result = fmt.Sprintf("[\n%s,\n]", utils.Indent(strings.Join(results, ",\n"), prefix+indent))
		} else {
			result = fmt.Sprintf("[%s]", strings.Join(results, ifIndent(", ", ",").(string)))
		}

	case map[string]interface{}:
		if key := singleMap(value); fullHcl && key != "" {
			var element string
			if element, err = marshalHCL(value[key], fullHcl, false, "", indent); err != nil {
				return
			}
			result = fmt.Sprintf(`%s %s`, id(key), element)
			break
		}

		keys := make([]string, 0, len(value))
		rendered := make(map[string]string, len(value))
		keyLen := 0

		for key, val := range value {
			keys = append(keys, key)
			if rendered[key], err = marshalHCL(val, fullHcl, false, "", indent); err != nil {
				return
			}
			if strings.Contains(rendered[key], "\n") {
				continue
			}
			if len(key) > keyLen && indent != "" {
				keyLen = len(key)
			}
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
					if indent == "" && strings.HasPrefix(rendered, `"`) && equal == "" {
						keyLen = len(id(key)) + 1
					}
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

func singleMap(m map[string]interface{}) string {
	if len(m) != 1 {
		return ""
	}
	for k := range m {
		if _, isMap := m[k].(map[string]interface{}); isMap {
			return k
		}
	}
	return ""
}

var identifierRegex = regexp.MustCompile(`^[A-za-z][\w-]*$`)

func id(key string) string {
	if identifierRegex.MatchString(key) {
		return key
	}
	// The identifier contains characters that may be considered invalid, we have to quote it
	return fmt.Sprintf("%q", key)
}
