package types

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/imdario/mergo"
)

// IDictionary represents objects that act as map[string]interface.
type IDictionary interface {
	Clone(keys ...interface{}) IDictionary
	Omit(keys ...interface{}) IDictionary
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}
	Has(key interface{}) bool
	Delete(key interface{})
	KeysAsString() []string
	Keys() IGenericList
	Len() int
	AsMap() map[string]interface{}
	Merge(dicts ...IDictionary) (IDictionary, error)
	String() string
}

// Dictionary implements base IDictionary.
type Dictionary map[string]interface{}

// String returns the string representation of the dictionary.
func (d Dictionary) String() string { return fmt.Sprint(d.AsMap()) }

// Set sets key to value in the dictionary.
func (d Dictionary) Set(key interface{}, value interface{}) { d[fmt.Sprint(key)] = value }

// Get returns the value associated with key.
func (d Dictionary) Get(key interface{}) interface{} { return d[fmt.Sprint(key)] }

// Has returns true if the dictionary object contains the key.
func (d Dictionary) Has(key interface{}) bool { _, ok := d[fmt.Sprint(key)]; return ok }

// Delete removes the entry value associated with key.
func (d Dictionary) Delete(key interface{}) { delete(d, fmt.Sprint(key)) }

// Len returns the number of keys in the dictionary
func (d Dictionary) Len() int { return len(d) }

// Keys returns the keys in the dictionary in alphabetical order.
func (d Dictionary) Keys() IGenericList { return NewGenericListFromStrings(d.KeysAsString()...) }

// AsMap returns the object casted as map[string]interface{}.
func (d Dictionary) AsMap() map[string]interface{} {
	if d == nil {
		return make(map[string]interface{})
	}
	return d
}

// Clone returns a distinct copy of the object with only supplied keys.
// If no keys are supplied, all keys from d are copied.
func (d Dictionary) Clone(keys ...interface{}) IDictionary {
	if len(keys) == 0 {
		keys = d.Keys().AsList()
	}
	newDict := make(Dictionary, len(d))
	for i := range keys {
		value := d.Get(keys[i])
		switch value := value.(type) {
		case IDictionary:
			value = value.Clone()
		case IGenericList:
			value = value.Clone()
		}
		newDict.Set(keys[i], value)
	}
	return newDict
}

// Omit returns a distinct copy of the object including all keys except specified ones.
func (d Dictionary) Omit(keys ...interface{}) IDictionary {
	omitKeys := make(map[interface{}]bool)
	for i := range keys {
		omitKeys[keys[i]] = true
	}
	keep := make([]interface{}, 0, len(d))
	for key := range d {
		if !omitKeys[key] {
			keep = append(keep, key)
		}
	}
	return d.Clone(keep...)
}

// KeysAsString returns the keys in the dictionary in alphabetical order.
func (d Dictionary) KeysAsString() []string {
	keys := make([]string, 0, len(d))
	for key := range d {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Merge merges the other dictionaries into the current dictionary.
func (d Dictionary) Merge(dicts ...IDictionary) (IDictionary, error) {
	m := d.AsMap()
	for i := range dicts {
		if dicts[i] == nil {
			continue
		}
		if err := mergo.Merge(&m, dicts[i].AsMap()); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// AsDictionary returns the object casted as IDictionary if possible.
func AsDictionary(object interface{}) (IDictionary, error) {
	if object == nil {
		return make(Dictionary), nil
	}

	if reflect.TypeOf(object).Kind() == reflect.Ptr {
		object = reflect.ValueOf(object).Elem().Interface()
		if object == nil {
			return make(Dictionary), nil
		}
	}

	target := reflect.TypeOf(Dictionary{})
	if !reflect.TypeOf(object).ConvertibleTo(target) {
		return nil, fmt.Errorf("Object cannot be converted to map: %T", object)
	}

	return reflect.ValueOf(object).Convert(target).Interface().(IDictionary), nil
}

// NewDictionary creates a dictionary from object if is is a map
func NewDictionary(object interface{}) (dict IDictionary, err error) {
	dict = make(Dictionary)
	if object == nil {
		return
	}

	if reflect.TypeOf(object).Kind() != reflect.Map {
		err = fmt.Errorf("Object cannot be converted to map: %T", object)
		return
	}

	value := reflect.ValueOf(object)
	keys := value.MapKeys()
	for i := range keys {
		dict.Set(keys[i].String(), value.MapIndex(keys[i]).Interface())
	}
	return
}

// MustNewDictionary creates a dictionary from object if is is a map
func MustNewDictionary(object interface{}) IDictionary {
	result, err := NewDictionary(object)
	if err != nil {
		panic(err)
	}
	return result
}
