package collections

import (
	"fmt"

	"github.com/coveo/gotemplate/errors"
)

// IDictionary represents objects that act as map[string]interface.
type IDictionary interface {
	AsMap() map[string]interface{}                                    // Returns the object casted as map[string]interface{}.
	Native() interface{}                                              // Returns the object casted as native go type (applied recursively).
	Clone(keys ...interface{}) IDictionary                            // Returns a distinct copy of the object with only supplied keys. If no keys are supplied, all keys from d are copied.
	Count() int                                                       // Simply an alias for Len.
	CreateList(...int) IGenericList                                   // Instantiates a list of the same type as current dictionary with optional size and capacity.
	Default(key, defVal interface{}) interface{}                      // Returns defVal if dictionary doesn't contain key, otherwise, simply returns entry corresponding to key.
	Delete(key interface{}, keys ...interface{}) (IDictionary, error) // Removes the entry value associated with key. The entry must exist.
	Flush(keys ...interface{}) IDictionary                            // Removes all specified keys from the dictionary. If no key is specified, all keys are removed.
	Get(key interface{}) interface{}                                  // Returns the value associated with key.
	Has(key interface{}) bool                                         // Returns true if the dictionary object contains the key.
	GetKeys() IGenericList                                            // Returns the keys in the dictionary in alphabetical order.
	KeysAsString() []string                                           // Returns the keys in the dictionary in alphabetical order.
	Len() int                                                         // Returns the number of keys in the dictionary.
	Merge(IDictionary, ...IDictionary) IDictionary                    // Merges the other dictionaries into the current dictionary.
	Omit(key interface{}, keys ...interface{}) IDictionary            // Returns a distinct copy of the object including all keys except specified ones.
	Set(key, value interface{}) IDictionary                           // Sets key to value in the dictionary.
	String() string                                                   // Returns the string representation of the dictionary.
	GetValues() IGenericList                                          // Returns the values in the dictionary in alphabetical order of keys.
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
	return errors.Must(TryAsDictionary(object)).(IDictionary)
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
