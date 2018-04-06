package types

// IDictionary represents objects that act as map[string]interface.
type IDictionary interface {
	// Cast(interface{}) IDictionary                              // Returns the object casted as same dictionary type.
	// CreateList(...int) IGenericList                            // Allocates a new list of the same type implementation as this dictionary. Optional arguments are size and capacity.
	// NewDictionary(interface{}) IDictionary                     // Creates a new dictionary of the same type from the supplied arguments.
	// TryCast(interface{}) (IDictionary, error)                  // Returns the object casted as same dictionary type if possible.
	// TryMerge(IDictionary, ...IDictionary) (IDictionary, error) // Merges the other dictionaries into the current dictionary.
	// TryNewDictionary(interface{}) (IDictionary, error)         // Creates a new dictionary of the same type from the supplied arguments.

	AsMap() map[string]interface{}                                    // Returns the object casted as map[string]interface{}.
	Clone(keys ...interface{}) IDictionary                            // Returns a distinct copy of the object with only supplied keys. If no keys are supplied, all keys from d are copied.
	Count() int                                                       // Simply an alias for Len.
	CreateList(...int) IGenericList                                   // Instantiates a list of the same type as current dictionary with optional size and capacity.
	Default(key, defVal interface{}) interface{}                      // Returns defVal if dictionary doesn't contain key, otherwise, simply returns entry corresponding to key.
	Delete(key interface{}, keys ...interface{}) (IDictionary, error) // Removes the entry value associated with key. The entry must exist.
	Flush(keys ...interface{}) IDictionary                            // Removes all specified keys from the dictionary. If no key is specified, all keys are removed.
	Get(key interface{}) interface{}                                  // Returns the value associated with key.
	Has(key interface{}) bool                                         // Returns true if the dictionary object contains the key.
	Keys() IGenericList                                               // Returns the keys in the dictionary in alphabetical order.
	KeysAsString() []string                                           // Returns the keys in the dictionary in alphabetical order.
	Len() int                                                         // Returns the number of keys in the dictionary.
	Merge(IDictionary, ...IDictionary) IDictionary                    // Merges the other dictionaries into the current dictionary.
	Omit(key interface{}, keys ...interface{}) IDictionary            // Returns a distinct copy of the object including all keys except specified ones.
	Set(key, value interface{}) IDictionary                           // Sets key to value in the dictionary.
	String() string                                                   // Returns the string representation of the dictionary.
}

// AsDictionary returns the object casted as IDictionary.
var AsDictionary func(interface{}) IDictionary

// CreateDictionary instantiates a new dictionary with optional size.
var CreateDictionary func(size ...int) IDictionary

// TryAsDictionary returns the object casted as IDictionary if possible.
var TryAsDictionary func(interface{}) (IDictionary, error)
