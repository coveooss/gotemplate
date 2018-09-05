package implementation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/errors"
	"github.com/imdario/mergo"
)

func (d baseDict) String() string {
	// Unlike go maps, we render dictionary keys in order
	keys := d.KeysAsString()
	values := make([]string, d.Len())
	for i := range values {
		values[i] = fmt.Sprintf("%s:%v", keys[i], d.Get(keys[i]))
	}
	return fmt.Sprintf("dict[%s]", strings.Join(values, " "))
}

// DictHelper implements basic functionalities required for IDictionary.
type DictHelper struct {
	BaseHelper
}

// AsDictionary returns the object casted as IDictionary.
func (dh DictHelper) AsDictionary(object interface{}) baseIDict {
	return errors.Must(dh.TryAsDictionary(object)).(baseIDict)
}

// Clone returns a distinct copy of the object with only supplied keys. If no keys are supplied, all keys from d are copied.
func (dh DictHelper) Clone(dict baseIDict, keys []interface{}) baseIDict {
	if len(keys) == 0 {
		keys = dict.GetKeys().AsArray()
	}
	newDict := dh.CreateDictionary(dict.Len())
	for i := range keys {
		value := dict.Get(keys[i])
		if value != nil {
			if v, err := dh.TryAsDictionary(value); err == nil {
				value = dh.Clone(v, nil)
			} else if v, err := dh.TryAsList(value); err == nil {
				value = dh.ConvertList(v.Clone())
			}
		}
		newDict.Set(keys[i], value)
	}
	return newDict
}

// Default returns defVal if dictionary doesn't contain key, otherwise, simply returns entry corresponding to key.
func (dh DictHelper) Default(dict baseIDict, key, defVal interface{}) interface{} {
	if !dict.Has(key) {
		return defVal
	}
	return dict.Get(key)
}

// Delete removes the entry value associated with key. The entry must exist.
func (dh DictHelper) Delete(dict baseIDict, keys []interface{}) (baseIDict, error) {
	return dh.delete(dict, keys, true)
}

// Flush removes all specified keys from the dictionary. If no key is specified, all keys are removed.
func (dh DictHelper) Flush(dict baseIDict, keys []interface{}) baseIDict {
	if len(keys) == 0 {
		keys = dict.GetKeys().AsArray()
	}
	dh.delete(dict, keys, false)
	return dict
}

// Get returns the value associated with key.
func (dh DictHelper) Get(dict baseIDict, key interface{}) interface{} {
	return dict.AsMap()[fmt.Sprint(key)]
}

// Has returns true if the dictionary object contains the key.
func (dh DictHelper) Has(dict baseIDict, key interface{}) bool {
	_, ok := dict.AsMap()[fmt.Sprint(key)]
	return ok
}

// GetKeys returns the keys in the dictionary in alphabetical order.
func (dh DictHelper) GetKeys(dict baseIDict) baseIList {
	keys := dict.KeysAsString()
	result := dh.CreateList(dict.Len())

	for i := range keys {
		result.Set(i, keys[i])
	}
	return result
}

// KeysAsString returns the keys in the dictionary in alphabetical order.
func (dh DictHelper) KeysAsString(dict baseIDict) []string {
	keys := make([]string, 0, dict.Len())
	for key := range dict.AsMap() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Merge merges the other dictionaries into the current dictionary.
func (dh DictHelper) Merge(target baseIDict, sources []baseIDict) baseIDict {
	m := target.AsMap()
	for i := range sources {
		if sources[i] == nil {
			continue
		}
		if err := mergo.Merge(&m, dh.ConvertDict(sources[i]).AsMap()); err != nil {
			panic(err)
		}
	}
	return target
}

// Omit returns a distinct copy of the object including all keys except specified ones.
func (dh DictHelper) Omit(dict baseIDict, keys []interface{}) baseIDict {
	omitKeys := make(map[string]bool, len(keys))
	for i := range keys {
		omitKeys[fmt.Sprint(keys[i])] = true
	}
	keep := make([]interface{}, 0, dict.Len())
	for key := range dict.AsMap() {
		if !omitKeys[key] {
			keep = append(keep, key)
		}
	}
	return dh.Clone(dict, keep)
}

// Set sets key to value in the dictionary.
func (dh DictHelper) Set(dict baseIDict, key interface{}, value interface{}) baseIDict {
	dict.AsMap()[fmt.Sprint(key)] = dh.Convert(value)
	return dict
}

// GetValues returns the values in the dictionary in key alphabetical order.
func (dh DictHelper) GetValues(dict baseIDict) baseIList {
	result := dh.CreateList(dict.Len())
	for i, key := range dict.KeysAsString() {
		result.Set(i, dict.Get(key))
	}
	return result
}

func (dh DictHelper) delete(dict baseIDict, keys []interface{}, mustExist bool) (baseIDict, error) {
	for i := range keys {
		if mustExist && !dict.Has(keys[i]) {
			return dict, fmt.Errorf("key %v not found", keys[i])
		}
		delete(dict.AsMap(), fmt.Sprint(keys[i]))
	}
	return dict, nil
}

// Register the default implementation of dictionary helper
var _ = func() int {
	collections.DictionaryHelper = baseDictHelper
	return 0
}()
