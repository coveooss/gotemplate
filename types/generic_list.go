package types

import (
	"fmt"
	"reflect"
)

// IGenericList represents objects that act as []interface{}.
type IGenericList interface {
	Clone() IGenericList
	Set(index int, value interface{}) (IGenericList, error)
	Get(index int) interface{}
	Len() int
	AsList() []interface{}
	String() string
}

// GenericList implements base IGenericList
type GenericList []interface{}

// String returns the string representation of the list.
func (l GenericList) String() string { return fmt.Sprint(l.AsList()) }

// Clone returns a distinct copy of the object.
func (l GenericList) Clone() IGenericList { return NewGenericList(l...) }

// Len returns the length of the list.
func (l GenericList) Len() int { return len(l) }

// AsList returns the current list as standard array of interface{}
func (l GenericList) AsList() []interface{} { return l }

// Set sets the value at position index int the list.
// If list is not large enough, it is enlarged to fit the index.
func (l GenericList) Set(index int, value interface{}) (IGenericList, error) {
	if index < 0 {
		return nil, fmt.Errorf("index must be positive number")
	}
	if index > len(l) {
		newList := make(GenericList, index+1)
		for i := range l {
			newList[i] = l[i]
		}
		l = newList
	}
	l[index] = value
	return l, nil
}

// Get returns the element at position index in the list.
// If index is out of bound, nil is returned
func (l GenericList) Get(index int) interface{} {
	if index < 0 || index >= len(l) {
		return nil
	}
	return l[index]
}

// NewGenericList instantiates a new GenericList from supplied arguments
func NewGenericList(items ...interface{}) IGenericList {
	newList := make(GenericList, len(items))
	for i := range items {
		newList[i] = items[i]
	}
	return newList
}

// NewGenericListFromStrings instantiates a new GenericList from supplied arguments
func NewGenericListFromStrings(items ...string) IGenericList {
	newList := make(GenericList, len(items))
	for i := range items {
		newList[i] = items[i]
	}
	return newList
}

// AsGenericList returns the object casted as IGenericList if possible
func AsGenericList(object interface{}) (result IGenericList, err error) {
	if object == nil {
		return nil, nil
	}

	target := reflect.TypeOf(GenericList{})
	t := reflect.TypeOf(object)
	if !t.ConvertibleTo(target) {
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			v := reflect.ValueOf(object)
			list := make(GenericList, v.Len())
			for i := range list {
				list[i] = v.Index(i).Interface()
			}
			return list, nil
		default:
			return nil, fmt.Errorf("Object cannot be converted to generic list: %T", object)
		}
	}

	return reflect.ValueOf(object).Convert(target).Interface().(IGenericList), nil
}
