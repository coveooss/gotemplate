package implementation

import (
	"fmt"
	"reflect"

	"github.com/coveooss/gotemplate/v3/collections"
)

type helperBase = BaseHelper
type helperList = ListHelper
type helperDict = DictHelper

// BaseHelper implements basic functionalities required for both IGenericList & IDictionary
type BaseHelper struct {
	ConvertList    func(list baseIList) baseIList
	ConvertDict    func(dict baseIDict) baseIDict
	NeedConversion func(object interface{}, strict bool) bool
}

// NewList creates a new IGenericList from supplied arguments.
func (bh BaseHelper) NewList(items ...interface{}) baseIList {
	if len(items) == 1 && items[0] != nil {
		v := reflect.ValueOf(items[0])
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			// There is only one items and it is an array or a slice
			items = make([]interface{}, v.Len())
			for i := 0; i < v.Len(); i++ {
				items[i] = v.Index(i).Interface()
			}
		}
	}
	newList := bh.CreateList(0, len(items))
	for i := range items {
		newList = newList.Append(items[i])
	}
	return newList
}

// NewStringList creates a new IGenericList from supplied arguments.
func (bh BaseHelper) NewStringList(items ...string) baseIList {
	newList := bh.CreateList(0, len(items))
	for i := range items {
		newList = newList.Append(items[i])
	}
	return newList
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
	return bh.tryAsDictionary(object, false)
}

// TryAsDictionaryStrict tries to convert any object to IDictionary object.
func (bh BaseHelper) TryAsDictionaryStrict(object interface{}) (baseIDict, error) {
	return bh.tryAsDictionary(object, true)
}

func (bh BaseHelper) tryAsDictionary(object interface{}, strict bool) (baseIDict, error) {
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
				return nil, fmt.Errorf("object cannot be converted to dictionary: %T", object)
			}
		}
	}

	if bh.NeedConversion(result, strict) {
		newDict := bh.CreateDictionary()
		for key, val := range result.AsMap() {
			// We loop on the key/values to ensure that all values are converted to the
			// desired type.
			newDict.Set(key, val)
		}
		result = newDict
	}

	return result, nil
}

// TryAsList tries to convert any object to IGenericList object.
func (bh BaseHelper) TryAsList(object interface{}) (baseIList, error) {
	return bh.tryAsList(object, false)
}

// TryAsListStrict tries to convert any object to IGenericList object.
func (bh BaseHelper) TryAsListStrict(object interface{}) (baseIList, error) {
	return bh.tryAsList(object, true)
}

func (bh BaseHelper) tryAsList(object interface{}, strict bool) (baseIList, error) {
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
				return nil, fmt.Errorf("object cannot be converted to generic list: %T", object)
			}
		}
	}
	if bh.NeedConversion(result, false) {
		newList := bh.CreateList(result.Len())
		for i, val := range result.AsArray() {
			newList.Set(i, val)
		}
		result = newList
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

// NeedConversion determine if the object need deep conversion.
//    strict indicates that the type must be converted to the desired type
//    even if the object implements the Dictionary or List interface.
func NeedConversion(object interface{}, strict bool, typeName string) bool {
	if object == nil {
		return false
	}
	objectType := reflect.TypeOf(object)
	switch objectType.Kind() {
	case reflect.Map:
		if dict, ok := object.(baseIDict); !ok || strict && dict.TypeName().Str() != typeName {
			return true
		}

		value := reflect.ValueOf(object)
		keys := value.MapKeys()
		for i := range keys {
			if NeedConversion(value.MapIndex(keys[i]).Interface(), strict, typeName) {
				return true
			}
		}
	case reflect.Slice, reflect.Array:
		if list, ok := object.(baseIList); !ok || strict && list.TypeName().Str() != typeName {
			return true
		}
		value := reflect.ValueOf(object)
		for i := 0; i < value.Len(); i++ {
			if NeedConversion(value.Index(i).Interface(), strict, typeName) {
				return true
			}
		}
	}
	return false
}

var (
	needConversionImpl = NeedConversion
	iif                = collections.IIf
)
