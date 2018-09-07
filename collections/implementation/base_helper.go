package implementation

import (
	"fmt"
	"reflect"

	"github.com/coveo/gotemplate/errors"
)

type helperBase = BaseHelper
type helperList = ListHelper
type helperDict = DictHelper

var must = errors.Must

// BaseHelper implements basic functionalities required for both IGenericList & IDictionary
type BaseHelper struct {
	ConvertList func(baseIList) baseIList
	ConvertDict func(baseIDict) baseIDict
}

// AsList converts object to IGenericList object. It panics if conversion is impossible.
func (bh BaseHelper) AsList(object interface{}) baseIList {
	return must(bh.TryAsList(object)).(baseIList)
}

// AsDictionary converts object to IDictionary object. It panics if conversion is impossible.
func (bh BaseHelper) AsDictionary(object interface{}) baseIDict {
	return must(bh.TryAsDictionary(object)).(baseIDict)
}

// Convert tries to convert the supplied object into IDictionary or IGenericList.
// Returns the supplied object if not conversion occurred.
func (bh BaseHelper) Convert(object interface{}) interface{} {
	object, _ = bh.TryConvert(object)
	return object
}

// CreateList creates a new IGenericList with optional size/capacity arguments.
func (bh BaseHelper) CreateList(args ...int) baseIList {
	var size, capacity int
	switch len(args) {
	case 0:
	case 1:
		size = args[0]
	case 2:
		size, capacity = args[0], args[1]
	default:
		panic(fmt.Errorf("CreateList only accept 2 arguments, size and capacity"))
	}
	if capacity < size {
		capacity = size
	}
	return bh.ConvertList(make(baseList, size, capacity))
}

// CreateDictionary creates a new IDictionary with optional capacity arguments.
func (bh BaseHelper) CreateDictionary(args ...int) baseIDict {
	var capacity int
	switch len(args) {
	case 0:
	case 1:
		capacity = args[0]
	default:
		panic(fmt.Errorf("CreateList only accept 1 argument for size"))
	}
	return bh.ConvertDict(make(baseDict, capacity))
}

// TryAsDictionary tries to convert any object to IDictionary object.
func (bh BaseHelper) TryAsDictionary(object interface{}) (baseIDict, error) {
	if object != nil && reflect.TypeOf(object).Kind() == reflect.Ptr {
		object = reflect.ValueOf(object).Elem().Interface()
	}

	var result baseIDict
	if dict, ok := object.(baseIDict); ok {
		// The object is already a IDictionary
		result = dict
	} else if object == nil {
		result = bh.CreateDictionary()
	} else {
		target := reflect.TypeOf(baseDict{})
		objectType := reflect.TypeOf(object)
		if objectType.ConvertibleTo(target) {
			result = bh.ConvertDict(reflect.ValueOf(object).Convert(target).Interface().(baseIDict))
		} else {
			switch objectType.Kind() {
			case reflect.Map:
				result = bh.CreateDictionary()
				value := reflect.ValueOf(object)
				keys := value.MapKeys()
				for i := range keys {
					result.Set(fmt.Sprint(keys[i]), value.MapIndex(keys[i]).Interface())
				}
			default:
				return nil, fmt.Errorf("Object cannot be converted to dictionary: %T", object)
			}
		}
	}

	for key, val := range result.AsMap() {
		// We loop on the key/values to ensure that all values are converted to the
		// desired type.
		result.Set(key, val)
	}

	return result, nil
}

// TryAsList tries to convert any object to IGenericList object.
func (bh BaseHelper) TryAsList(object interface{}) (baseIList, error) {
	if object != nil && reflect.TypeOf(object).Kind() == reflect.Ptr {
		object = reflect.ValueOf(object).Elem().Interface()
	}

	var result baseIList
	if list, ok := object.(baseIList); ok {
		// The object is already a IGenericList
		result = list
	} else if object == nil {
		result = bh.CreateList()
	} else {
		target := reflect.TypeOf(baseList{})
		objectType := reflect.TypeOf(object)
		if objectType.ConvertibleTo(target) {
			result = bh.ConvertList(reflect.ValueOf(object).Convert(target).Interface().(baseIList))
		} else {
			switch objectType.Kind() {
			case reflect.Slice, reflect.Array:
				value := reflect.ValueOf(object)
				result = bh.CreateList(value.Len())
				for i := 0; i < result.Len(); i++ {
					result.Set(i, value.Index(i).Interface())
				}
			default:
				return nil, fmt.Errorf("Object cannot be converted to generic list: %T", object)
			}
		}
	}
	for i, val := range result.AsArray() {
		result.Set(i, val)
	}

	return result, nil
}

// TryConvert tries to convert any object to IGenericList or IDictionary object.
// Returns true if a conversion occurred.
func (bh BaseHelper) TryConvert(object interface{}) (interface{}, bool) {
	if object != nil {
		if o, err := bh.TryAsDictionary(object); err == nil {
			return o, true
		} else if o, err := bh.TryAsList(object); err == nil {
			return o, true
		}
	}
	return object, false
}
