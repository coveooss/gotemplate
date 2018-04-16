package utils

import (
	"reflect"

	"github.com/coveo/gotemplate/collections"
)

// String is simply an alias of collections.String
type String = collections.String

var isEmpty = collections.IsEmptyValue

// IfUndef returns default value if nothing is supplied as values
func IfUndef(def interface{}, values ...interface{}) interface{} {
	switch len(values) {
	case 0:
		return def
	case 1:
		if values[0] == nil || reflect.TypeOf(values[0]).Kind() == reflect.Ptr && isEmpty(reflect.ValueOf(values[0])) {
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
	if isEmpty(reflect.ValueOf(testValue)) {
		return valueFalse
	}
	return valueTrue
}

// Default returns the value if it is not empty or default value.
func Default(value, defaultValue interface{}) interface{} {
	return IIf(value, value, defaultValue)
}

// MergeDictionary merges multiple dictionaries into a single one prioritizing the first ones.
func MergeDictionary(args ...map[string]interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return make(map[string]interface{}), nil
	}
	dicts := make([]collections.IDictionary, len(args))
	for i := range dicts {
		var err error
		dicts[i], err = collections.TryAsDictionary(args[i])
		if err != nil {
			return nil, err
		}
	}

	result := collections.CreateDictionary()
	return result.Merge(dicts[0], dicts[1:]...).Native().(map[string]interface{}), nil
}
