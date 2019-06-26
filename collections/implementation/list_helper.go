package implementation

import (
	"fmt"
	"reflect"

	"github.com/coveooss/gotemplate/v3/collections"
)

func (l baseList) String() string      { return fmt.Sprint(l.AsArray()) }
func (l baseList) PrettyPrint() string { return l.String() }

// ListHelper implements basic functionalities required for IGenericList.
type ListHelper struct {
	BaseHelper
}

// Add adds elements at the end of the supplied list.
func (lh ListHelper) Add(list baseIList, prepend bool, objects ...interface{}) baseIList {
	array := list.AsArray()
	for i := range objects {
		objects[i] = lh.Convert(objects[i])
	}
	if prepend {
		array, objects = objects, array
	}
	return lh.AsList(append(array, objects...))
}

// Clone returns a copy of the supplied list.
func (lh ListHelper) Clone(list baseIList) baseIList {
	return lh.NewList(list.AsArray()...)
}

// GetIndexes returns the element at position index in the list. If index is out of bound, nil is returned.
func (lh ListHelper) GetIndexes(list baseIList, indexes ...int) interface{} {
	switch len(indexes) {
	case 0:
		return nil
	case 1:
		index := indexes[0]
		if index < 0 {
			// If index is negative, we try to get from the end
			index += list.Len()
		}
		if index < 0 || index >= list.Len() {
			return nil
		}
		return (list.AsArray())[index]
	}
	result := list.Create(len(indexes))
	for i := range indexes {
		result.Set(i, lh.GetIndexes(list, indexes[i]))
	}
	return result
}

// GetStrings returns a string array representation of the list.
func (lh ListHelper) GetStrings(list baseIList) []string {
	return collections.ToStrings(list.AsArray())
}

// GetStringArray returns a StringArray representation of the list.
func (lh ListHelper) GetStringArray(list baseIList) strArray {
	result := make(strArray, list.Len())
	for i := 0; i < list.Len(); i++ {
		result[i] = str(fmt.Sprint(list.Get(i)))
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
		list = lh.Add(list, false, nil)
	}
	list.AsArray()[index] = lh.Convert(value)
	return list, nil
}

// Register the implementation of list functions
var _ = func() int {
	collections.ListHelper = baseListHelper
	return 0
}()

// Unique returns a copy of the list removing all duplicate elements.
func (lh ListHelper) Unique(list baseIList) baseIList {
	source := list.AsArray()
	target := lh.CreateList(0, list.Len())
	for i := range source {
		if !target.Contains(source[i]) {
			target = target.Append(source[i])
		}
	}
	return target
}

// Intersect returns a new list that is the result of the intersection of the list and the parameters.
func (lh ListHelper) Intersect(list baseIList, values ...interface{}) baseIList {
	source := list.Unique().AsArray()
	include := collections.AsList(values)
	target := lh.CreateList(0, include.Len())
	for i := range source {
		if include.Contains(source[i]) {
			target = target.Append(source[i])
		}
	}
	return target
}

// Remove returns a new list without the element specified.
func (lh ListHelper) Remove(list baseIList, indexes ...int) baseIList {
	for i, index := range indexes {
		if index < 0 {
			indexes[i] += list.Len()
		}
	}
	discard := collections.AsList(indexes)
	target := list.Create(0, list.Len())
	for i := range list.AsArray() {
		if !discard.Contains(i) {
			target = target.Append(list.Get(i))
		}
	}
	return target
}

// Without returns a copy of the list removing specified elements.
func (lh ListHelper) Without(list baseIList, values ...interface{}) baseIList {
	source := list.AsArray()
	exclude := collections.AsList(values)
	target := lh.CreateList(0, list.Len())
	for i := range source {
		if !exclude.Contains(source[i]) {
			target = target.Append(source[i])
		}
	}
	return target
}

// Contains indicates if the list contains all specified elements
func (lh ListHelper) Contains(list baseIList, values ...interface{}) bool {
	source := list.AsArray()
	for _, value := range values {
		match := false
		for _, item := range source {
			if fmt.Sprint(value) == fmt.Sprint(item) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return len(source) > 0
}
