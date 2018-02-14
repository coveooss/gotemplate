package utils

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/imdario/mergo"
)

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
		if values[0] == nil {
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

// MergeMaps merges several maps into one privileging the leftmost
func MergeMaps(destination map[string]interface{}, sources ...map[string]interface{}) (map[string]interface{}, error) {
	for i := range sources {
		if err := mergo.Merge(&destination, sources[i]); err != nil {
			return destination, err
		}
	}
	return destination, nil
}
