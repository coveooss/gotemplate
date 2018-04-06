package types

import "github.com/coveo/gotemplate/errors"

// IGenericList represents objects that act as []interface{}.
type IGenericList interface {
	// Cast(interface{}) IGenericList             // Returns the object casted as same list type.
	// CreateList(...int) IGenericList            // Allocates a new list of the same type implementation as this list. Optional arguments are size and capacity.
	// NewList(...interface{}) IGenericList       // Creates a new generic list from the supplied arguments.
	// TryCast(interface{}) (IGenericList, error) // Returns the object casted as same list type if possible.

	AsArray() *[]interface{}                                // Returns the current list as standard array of interface{}.
	Append(...interface{}) IGenericList                     // Add elements to to current list. If list is not large enough, it is enlarged to fit the required size.
	Cap() int                                               // Returns the capacity of the list.
	Capacity() int                                          // Simply an alias for Cap.
	Clone() IGenericList                                    // Returns a distinct copy of the object.
	Count() int                                             // Simply an alias for Len.
	Get(index int) interface{}                              // Returns the element at position index in the list. If index is out of bound, nil is returned.
	Len() int                                               // Returns the number of elements in the list.
	Reverse() IGenericList                                  // Returns a copy of the current list in reverse order.
	Set(index int, value interface{}) (IGenericList, error) // Sets the value at position index into the list. If list is not large enough, it is enlarged to fit the index.
	String() string                                         // Returns the string representation of the list.
}

// AsList returns the object casted as IGenericList.
func AsList(object interface{}) IGenericList {
	return errors.Must(TryAsList(object)).(IGenericList)
}

// CreateList instantiates a new generic list with optional size and capacity.
var CreateList func(...int) IGenericList

// NewList instantiates a new generic list from supplied arguments.
var NewList func(...interface{}) IGenericList

// NewStringList creates a new IGenericList from supplied arguments.
var NewStringList func(...string) IGenericList

// TryAsList returns the object casted as IGenericList if possible.
var TryAsList func(interface{}) (IGenericList, error)
