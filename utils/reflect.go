package utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/coveo/gotemplate/types"
)

type dictionary = types.Dictionary

// String is simply an alias of types.String
type String = types.String

// IsEmptyValue determines if a value is a zero value
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}

// IsExported reports whether the identifier is exported.
func IsExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}

// IfUndef returns default value if nothing is supplied as values
func IfUndef(def interface{}, values ...interface{}) interface{} {
	switch len(values) {
	case 0:
		return def
	case 1:
		if values[0] == nil || reflect.TypeOf(values[0]).Kind() == reflect.Ptr && IsEmptyValue(reflect.ValueOf(values[0])) {
			return def
		}
		return values[0]
	default:
		return values
	}
}

// IIf acts as a generic ternary operator. It returns valueTrue if testValue is not empty,
// otherwise, it returns valueFalse
func IIf(testValue, valueTrue, valueFalse interface{}) interface{} {
	if IsEmptyValue(reflect.ValueOf(testValue)) {
		return valueFalse
	}
	return valueTrue
}

// Default returns the value if it is not empty or default value.
func Default(value, defaultValue interface{}) interface{} {
	return IIf(value, value, defaultValue)
}

// ToNativeRepresentation converts any object to native (literals, maps, slices)
func ToNativeRepresentation(value interface{}) interface{} {
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
			result[i] = ToNativeRepresentation(reflect.ValueOf(value).Index(i).Interface())
		}
		if len(result) == 1 && reflect.TypeOf(result[0]).Kind() == reflect.Map {
			// If the result is an array of one map, we just return the inner element
			return result[0]
		}
		return result

	case reflect.Map:
		result := make(dictionary, val.Len())
		for _, key := range val.MapKeys() {
			result[fmt.Sprintf("%v", key)] = ToNativeRepresentation(val.MapIndex(key).Interface())
		}
		return result

	case reflect.Struct:
		result := make(dictionary, typ.NumField())
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
				if !IsExported(t.Name()) && t.Kind() != reflect.Struct {
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

			split := strings.Split(tag, ",")
			name := split[0]
			if name == "" {
				name = sf.Name
			}
			options := make(map[string]bool, len(split[1:]))
			for i := range split[1:] {
				options[split[i+1]] = true
			}

			if options["omitempty"] && IsEmptyValue(val.Field(i)) {
				continue
			}

			if options["inline"] {
				for key, value := range ToNativeRepresentation(val.Field(i).Interface()).(map[string]interface{}) {
					result[key] = value
				}
			} else {
				result[name] = ToNativeRepresentation(val.Field(i).Interface())
			}
		}
		return result
	default:
		fmt.Fprintf(os.Stderr, "Unknown type %T %v : %v\n", value, typ.Kind(), value)
		return fmt.Sprintf("%v", value)
	}
}
