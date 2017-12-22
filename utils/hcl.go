package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/hcl"
)

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
}

// FlattenHCL - Convert array of map to single map if there is only one element in the array
// By default, the hcl.Unmarshal returns array of map even if there is only a single map in the definition
func FlattenHCL(source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		switch value := value.(type) {
		case []map[string]interface{}:
			switch len(value) {
			case 1:
				source[key] = FlattenHCL(value[0])
			default:
				for i, subMap := range value {
					value[i] = FlattenHCL(subMap)
				}
			}
		}
	}
	return source
}

// LoadHCL loads hcl file into variable
func LoadHCL(filename string) (result map[string]interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		if err = hcl.Unmarshal(content, &result); err == nil {
			result = FlattenHCL(result)
		}
	}
	return
}

// ToHCL serialize values to hcl format
func ToHCL(value interface{}) []byte {
	return marshalHCL(value, false, 0)
}

// ToPrettyHCL serialize values to hcl format with indentation
func ToPrettyHCL(value interface{}) []byte {
	return marshalHCL(value, true, 0)
}

func marshalHCL(value interface{}, pretty bool, indent int) []byte {
	if value == nil {
		return []byte("null")
	}

	var buffer bytes.Buffer
	typ, val := reflect.TypeOf(value), reflect.ValueOf(value)

	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return []byte("null")
		}
		val = val.Elem()
		typ = val.Type()
	}

	switch typ.Kind() {
	case reflect.String:
		buffer.WriteString(fmt.Sprintf("%q", strings.Replace(val.String(), `\`, `\\`, -1)))

	case reflect.Int:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Bool:
		buffer.WriteString(fmt.Sprintf("%v", value))

	case reflect.Slice:
		fallthrough
	case reflect.Array:
		switch val.Len() {
		case 0:
			buffer.WriteString("[]")
		case 1:
			buffer.WriteByte('[')
			buffer.Write(marshalHCL(reflect.ValueOf(value).Index(0).Interface(), true, indent+1))
			buffer.WriteByte(']')
		default:
			buffer.WriteString("[")
			indent++
			if pretty {
				buffer.WriteString("\n")
			}
			for i := 0; i < val.Len(); i++ {
				if pretty {
					buffer.WriteString(strings.Repeat(" ", indent*2))
				}
				buffer.Write(marshalHCL(val.Index(i).Interface(), pretty, indent+1))
				if i < val.Len()-1 || pretty {
					buffer.WriteString(",")
				}
				if pretty {
					buffer.WriteString("\n")
				}
			}
			if pretty && indent > 0 {
				buffer.WriteString(strings.Repeat(" ", (indent-1)*2))
			}
			buffer.WriteString("]")
		}

	case reflect.Map:
		if indent == 0 {
			value = FlattenHCL(value.(map[string]interface{}))
		}
		switch value := value.(type) {
		case map[string]interface{}:
			keys := make([]string, 0, len(value))

			for _, key := range val.MapKeys() {
				keys = append(keys, key.String())
			}
			sort.Strings(keys)

			if indent > 0 {
				buffer.WriteString("{")
				if pretty {
					buffer.WriteString("\n")
				}
			}

			for i, key := range keys {
				if pretty {
					buffer.WriteString(strings.Repeat(" ", indent*2))
				}
				if identifierRegex.MatchString(key) {
					buffer.WriteString(key)
				} else {
					// The identifier contains characters that may be considered invalid, we have to quote it
					buffer.WriteString(fmt.Sprintf("%q", key))
				}
				if pretty {
					buffer.WriteString(" = ")
				} else {
					buffer.WriteString("=")
				}

				buffer.Write(marshalHCL(value[key], pretty, indent+1))
				if pretty {
					buffer.WriteString("\n")
				} else if i < len(keys)-1 {
					buffer.WriteString(" ")
				}
			}

			if indent > 0 {
				if pretty {
					buffer.WriteString(strings.Repeat(" ", (indent-1)*2))
				}
				buffer.WriteString("}")
			}
		}
	case reflect.Struct:
		buffer.Write(marshalHCL(structs.Map(value), pretty, indent))
	default:
		fmt.Fprintf(os.Stderr, "Unknown type %T %v : %v\n", value, typ.Kind(), value)
		buffer.WriteString(fmt.Sprintf("%v", value))
	}

	return buffer.Bytes()
}

var identifierRegex = regexp.MustCompile(`^[A-za-z][\w-]*$`)
