package collections

import "reflect"

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
