package implementation

import (
	"fmt"
	"reflect"

	"github.com/coveo/gotemplate/collections"
)

func (l baseList) String() string { return fmt.Sprint(l.AsArray()) }

// ListHelper implements basic functionalities required for IGenericList.
type ListHelper struct {
	BaseHelper
}

// Append adds elements at the end of the supplied list.
func (lh ListHelper) Append(list baseIList, objects ...interface{}) baseIList {
	array := list.AsArray()
	for i := range objects {
		objects[i] = lh.Convert(objects[i])
	}
	return lh.AsList(append(array, objects...))
}

// Clone returns a copy of the supplied list.
func (lh ListHelper) Clone(list baseIList) baseIList {
	return lh.NewList(list.AsArray()...)
}

// GetIndex returns the element at position index in the list. If index is out of bound, nil is returned.
func (lh ListHelper) GetIndex(list baseIList, index int) interface{} {
	if index < 0 || index >= list.Len() {
		return nil
	}
	return (list.AsArray())[index]
}

// GetStrings returns a string array representation of the array.
func (lh ListHelper) GetStrings(list baseIList) []string {
	result := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		result[i] = fmt.Sprint(list.Get(i))
	}
	return result
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

// Reverse returns a copy of the current list in reverse order.
func (lh ListHelper) Reverse(list baseIList) baseIList {
	source := list.AsArray()
	target := lh.CreateList(list.Len())
	for i := range source {
		target.Set(target.Len()-i-1, source[i])
	}
	return lh.ConvertList(target)
}

// SetIndex sets the value at position index into the list. If list is not large enough, it is enlarged to fit the index.
func (lh ListHelper) SetIndex(list baseIList, index int, value interface{}) (baseIList, error) {
	if index < 0 {
		return nil, fmt.Errorf("index must be positive number")
	}
	for list.Len() <= index {
		list = lh.Append(list, nil)
	}
	list.AsArray()[index] = lh.Convert(value)
	return list, nil
}

// Register the implementation of list functions
var _ = func() int {
	collections.ListHelper = baseListHelper
	return 0
}()
