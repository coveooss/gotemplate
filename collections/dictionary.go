package collections

import (
	"fmt"
)

// IDictionary represents objects that act as map[string]interface.
type IDictionary interface {
	Add(key, value interface{}) IDictionary                  // Add value to key (if key exist, convert the key value into list and append the new value).
	AsMap() map[string]interface{}                           // Returns the object casted as map[string]interface{}.
	Clone(...interface{}) IDictionary                        // Returns a distinct copy of the object with only supplied keys. If no keys are supplied, all keys from d are copied.
	Count() int                                              // Simply an alias for Len.
	Create(...int) IDictionary                               // Instantiates a new dictionary of the same type with optional size.
	CreateList(...int) IGenericList                          // Instantiates a list of the same type as current dictionary with optional size and capacity.
	Default(key, defVal interface{}) interface{}             // Returns defVal if dictionary doesn't contain key, otherwise, simply returns entry corresponding to key.
	Delete(interface{}, ...interface{}) (IDictionary, error) // Removes the entry value associated with key. The entry must exist.
	Flush(...interface{}) IDictionary                        // Removes all specified keys from the dictionary. If no key is specified, all keys are removed.
	Get(...interface{}) interface{}                          // Returns the values associated with key.
	GetHelpers() (IDictionaryHelper, IListHelper)            // Returns the helpers implementation associated with the current type.
	GetKeys() IGenericList                                   // Returns the keys in the dictionary in alphabetical order.
	GetValues() IGenericList                                 // Returns the values in the dictionary in alphabetical order of keys.
	Has(...interface{}) bool                                 // Returns true if the dictionary object contains all the key.
	KeysAsString() StringArray                               // Returns the keys in the dictionary in alphabetical order.
	Len() int                                                // Returns the number of keys in the dictionary.
	Merge(IDictionary, ...IDictionary) IDictionary           // Merges the other dictionaries into the current dictionary.
	Native() interface{}                                     // Returns the object casted as native go type (applied recursively).
	Omit(interface{}, ...interface{}) IDictionary            // Returns a distinct copy of the object including all keys except specified ones.
	Pop(...interface{}) interface{}                          // Returns and remove the objects with the specified keys.
	PrettyPrint() string                                     // Returns the pretty string representation of the dictionary.
	Set(key, value interface{}) IDictionary                  // Sets key to value in the dictionary.
	String() string                                          // Returns the string representation of the dictionary.
	Transpose() IDictionary                                  // Transpose keys/values and return the resulting dictionary
	TypeName() String                                        // Returns the actual type name
}

// IDictionaryHelper represents objects that implement IDictionary compatible objects
type IDictionaryHelper interface {
	AsDictionary(interface{}) IDictionary                    // Returns the object casted as IDictionary.
	Convert(object interface{}) interface{}                  // Tries to convert the supplied object into IDictionary or IGenericList.
	CreateDictionary(args ...int) IDictionary                // Creates a new IDictionary with optional capacity arguments.
	TryAsDictionary(object interface{}) (IDictionary, error) // Tries to convert any object to IDictionary objects
	TryConvert(object interface{}) (interface{}, bool)       // Tries to convert any object to IGenericList or IDictionary object.
}

// DictionaryHelper configures the default dictionary manager.
var DictionaryHelper IDictionaryHelper

func assertDictionaryHelper() {
	if DictionaryHelper == nil {
		panic(fmt.Errorf("DictionaryHelper not configured"))
	}
}

// AsDictionary returns the object casted as IDictionary.
func AsDictionary(object interface{}) IDictionary {
	return must(TryAsDictionary(object)).(IDictionary)
}

// CreateDictionary instantiates a new dictionary with optional size.
func CreateDictionary(size ...int) IDictionary {
	assertDictionaryHelper()
	return DictionaryHelper.CreateDictionary(size...)
}

// TryAsDictionary returns the object casted as IDictionary if possible.
func TryAsDictionary(object interface{}) (IDictionary, error) {
	if result, ok := object.(IDictionary); ok {
		return result, nil
	}
	assertDictionaryHelper()
	return DictionaryHelper.TryAsDictionary(object)
}
