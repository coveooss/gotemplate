package collections

import (
	"fmt"

	"github.com/coveo/gotemplate/errors"
)

// IGenericList represents objects that act as []interface{}.
type IGenericList interface {
	AsArray() []interface{}                                 // Returns the current list as standard array of interface{}.
	Append(...interface{}) IGenericList                     // Add elements to the current list. If list is not large enough, it is enlarged to fit the required size.
	Cap() int                                               // Returns the capacity of the list.
	Capacity() int                                          // Simply an alias for Cap.
	Clone() IGenericList                                    // Returns a distinct copy of the object.
	Count() int                                             // Simply an alias for Len.
	Contains(...interface{}) bool                           // Indicates if the list contains all specified elements
	Create(...int) IGenericList                             // Allocates a new list of the same type implementation as this list. Optional arguments are size and capacity.
	Get(index int) interface{}                              // Returns the element at position index in the list. If index is out of bound, nil is returned.
	Len() int                                               // Returns the number of elements in the list.
	New(...interface{}) IGenericList                        // Creates a new generic list from the supplied arguments.
	Prepend(...interface{}) IGenericList                    // Add elements to the beginning of the current list. If list is not large enough, it is enlarged to fit the required size.
	Reverse() IGenericList                                  // Returns a copy of the current list in reverse order.
	Set(index int, value interface{}) (IGenericList, error) // Sets the value at position index into the list. If list is not large enough, it is enlarged to fit the index.
	String() string                                         // Returns the string representation of the list.
	Strings() []string                                      // Returns the current list as list of strings.
	Unique() IGenericList                                   // Returns a copy of the list removing all duplicate elements.
	Without(...interface{}) IGenericList                    // Returns a copy of the list removing specified elements.
}

// IListHelper represents objects that implement IGenericList compatible objects
type IListHelper interface {
	AsList(interface{}) IGenericList                    // Converts object to IGenericList object. It panics if conversion is impossible.
	Convert(object interface{}) interface{}             // Tries to convert the supplied object into IDictionary or IGenericList.
	CreateList(...int) IGenericList                     // Creates a new IGenericList with optional size/capacity arguments.
	NewList(...interface{}) IGenericList                // Creates a new IGenericList from supplied arguments.
	NewStringList(...string) IGenericList               // Creates a new IGenericList from supplied arguments.
	TryAsList(object interface{}) (IGenericList, error) // Tries to convert any object to IGenericList object.
	TryConvert(object interface{}) (interface{}, bool)  // Tries to convert any object to IGenericList or IDictionary object.
}

// ListHelper configures the default list manager.
var ListHelper IListHelper

func assertListHelper() {
	if ListHelper == nil {
		panic(fmt.Errorf("ListHelper not configured"))
	}
}

// AsList returns the object casted as IGenericList.
func AsList(object interface{}) IGenericList {
	return errors.Must(TryAsList(object)).(IGenericList)
}

// CreateList instantiates a new generic list with optional size and capacity.
func CreateList(args ...int) IGenericList {
	assertListHelper()
	return ListHelper.CreateList(args...)
}

// NewList instantiates a new generic list from supplied arguments.
func NewList(objects ...interface{}) IGenericList {
	assertListHelper()
	return ListHelper.NewList(objects...)
}

// NewStringList creates a new IGenericList from supplied arguments.
func NewStringList(objects ...string) IGenericList {
	assertListHelper()
	return ListHelper.NewStringList(objects...)
}

// TryAsList returns the object casted as IGenericList if possible.
func TryAsList(object interface{}) (IGenericList, error) {
	if result, ok := object.(IGenericList); ok {
		return result, nil
	}
	assertListHelper()
	return ListHelper.TryAsList(object)
}
