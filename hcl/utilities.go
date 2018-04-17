package hcl

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/utils"
)

// flatten converts array of map to single map if there is only one element in the array.
// By default, hc.Unmarshal returns array of map even if there is only a single map in the definition.
func flatten(source interface{}) interface{} {
	switch value := source.(type) {
	case []map[string]interface{}:
		switch len(value) {
		case 1:
			source = flatten(value[0])
		default:
			result := make([]map[string]interface{}, len(value))
			for i := range value {
				result[i] = flatten(value[i]).(map[string]interface{})
			}
			source = result
		}
	case map[string]interface{}:
		for key := range value {
			value[key] = flatten(value[key])
		}
		source = value
	case []interface{}:
		for i, sub := range value {
			value[i] = flatten(sub)
		}
		source = value
	}
	return source
}

func transform(out interface{}) {
	result := transformElement(flatten(reflect.ValueOf(out).Elem().Interface()))
	if _, isMap := out.(*map[string]interface{}); isMap {
		// If the result is expected to be map[string]interface{}, we convert it back from internal dict type.
		result = result.(hclIDict).Native()
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
}

func transformElement(source interface{}) interface{} {
	if value, err := hclHelper.TryAsDictionary(source); err == nil {
		for _, key := range value.KeysAsString() {
			value.Set(key, transformElement(value.Get(key)))
		}
		source = value
	} else if value, err := hclHelper.TryAsList(source); err == nil {
		for i, sub := range value.AsArray() {
			value.Set(i, transformElement(sub))
		}
		source = value
	}
	return source
}

func marshalHCL(value interface{}, fullHcl, head bool, prefix, indent string) (result string, err error) {
	if value == nil {
		result = "null"
		return
	}

	ifIndent := func(vTrue, vFalse interface{}) interface{} { return utils.IIf(indent, vTrue, vFalse) }
	const specialFormat = "#HCL_ARRAY_MAP#!"

	switch value := value.(type) {
	case int, float64, bool:
		result = fmt.Sprint(value)
	case string:
		value = fmt.Sprintf("%q", value)
		if indent != "" && strings.Contains(value, "\\n") {
			// We unquote the value
			unIndented := value[1 : len(value)-1]
			// Then replace escaped characters, other escape chars are \a, \b, \f and \v are not managed
			unIndented = strings.Replace(unIndented, `\n`, "\n", -1)
			unIndented = strings.Replace(unIndented, `\\`, "\\", -1)
			unIndented = strings.Replace(unIndented, `\"`, "\"", -1)
			unIndented = strings.Replace(unIndented, `\r`, "\r", -1)
			unIndented = strings.Replace(unIndented, `\t`, "\t", -1)
			unIndented = collections.UnIndent(unIndented)
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
			result = fmt.Sprintf("[\n%s,\n]", collections.Indent(strings.Join(results, ",\n"), prefix+indent))
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
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		rendered := make(map[string]string, len(keys))
		keyLen := 0

		for key, val := range value {
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
				if _, err := hclHelper.TryAsDictionary(value[key]); err == nil {
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

		result = fmt.Sprintf("{\n%s\n}", collections.Indent(strings.Join(items, "\n"), prefix+indent))

	default:
		debug.PrintStack()
		err = fmt.Errorf("marshalHCL Unknown type %[1]T %[1]v", value)
	}
	return
}

func isArrayOfMap(array []interface{}) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		if item, isMap := item.([]map[string]interface{}); !isMap || len(item) != 1 {
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
