package implementation

import (
	"fmt"
	"reflect"

	"github.com/coveooss/gotemplate/v3/collections"
)

func (d baseDict) String() string {
	// Unlike go maps, we render dictionary keys in order
	keys := d.KeysAsString()
	for i, k := range keys {
		keys[i] = str(fmt.Sprintf("%s:%v", k, d.Get(k)))
	}
	return fmt.Sprintf("dict[%s]", keys.Join(" "))
}

func (d baseDict) PrettyPrint() string { return d.String() }

// DictHelper implements basic functionalities required for IDictionary.
type DictHelper struct {
	BaseHelper
}

// AsDictionary returns the object casted as IDictionary.
func (dh DictHelper) AsDictionary(object interface{}) baseIDict {
	return must(dh.TryAsDictionary(object)).(baseIDict)
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
func (dh DictHelper) Get(dict baseIDict, keys []interface{}) interface{} {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		return dict.AsMap()[fmt.Sprint(keys[0])]
	}
	result := dict.CreateList(len(keys))
	for i := range result.AsArray() {
		result.Set(i, dict.Get(keys[i]))
	}
	return result
}

// Has returns true if the dictionary object contains all the keys.
func (dh DictHelper) Has(dict baseIDict, keys []interface{}) bool {
	for _, key := range keys {
		if _, ok := dict.AsMap()[fmt.Sprint(key)]; !ok {
			return false
		}
	}
	return dict.Len() > 0
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

// GetTypes returns a dictionary with key and type (or kind) for each element.
func (dh DictHelper) GetTypes(dict baseIDict, kind bool) baseIDict {
	result := dh.CreateDictionary(dict.Len())

	for key, value := range dict.AsMap() {
		result.Set(key, iif(kind, reflect.TypeOf(value).Kind().String(), reflect.TypeOf(value).Name()))
	}
	return result
}

// KeysAsString returns the keys in the dictionary in alphabetical order.
func (dh DictHelper) KeysAsString(dict baseIDict) collections.StringArray {
	keys := make(collections.StringArray, 0, dict.Len())
	for key := range dict.AsMap() {
		keys = append(keys, str(key))
	}
	return keys.Sorted()
}

// Merge merges the other dictionaries into the current dictionary.
func (dh DictHelper) Merge(target baseIDict, sources []baseIDict) baseIDict {
	for i := range sources {
		if sources[i] == nil {
			continue
		}
		target = dh.deepMerge(target, dh.ConvertDict(sources[i]))
	}
	return target
}

func (dh DictHelper) deepMerge(target baseIDict, source baseIDict) baseIDict {
	targetMap := target.AsMap()
	sourceMap := source.AsMap()
	for key := range sourceMap {
		sourceValue, sourceHasKey := sourceMap[key]
		targetValue, targetHasKey := targetMap[key]

		if sourceHasKey && !targetHasKey {
			targetMap[key] = sourceValue
			continue
		}

		targetValueDict, targetValueIsDict := targetValue.(baseIDict)
		sourceValueDict, sourceValueIsDict := sourceValue.(baseIDict)

		if sourceValueIsDict && targetValueIsDict {
			targetMap[key] = dh.deepMerge(targetValueDict, sourceValueDict)
		}
	}

	return target
}

// Omit returns a distinct copy of the object including all keys except specified ones.
func (dh DictHelper) Omit(dict baseIDict, keys []interface{}) baseIDict {
	if len(keys) == 0 {
		return dict.Clone()
	}
	return dict.Clone().Flush(keys...)
}

// Pop returns and remove the objects with the specified keys.
func (dh DictHelper) Pop(dict baseIDict, keys []interface{}) interface{} {
	if len(keys) == 0 {
		return nil
	}
	result := dh.Get(dict, keys)
	dh.delete(dict, keys, false)
	return result
}

// Set sets key to value in the dictionary.
func (dh DictHelper) Set(dict baseIDict, key interface{}, value interface{}) baseIDict {
	if dict.AsMap() == nil {
		dict = dh.CreateDictionary()
	}
	dict.AsMap()[fmt.Sprint(key)] = dh.Convert(value)
	return dict
}

// Add adds value to an existing key instead of replacing the value as done by set.
func (dh DictHelper) Add(dict baseIDict, key interface{}, value interface{}) baseIDict {
	if dict.AsMap() == nil {
		dict = dh.CreateDictionary()
	}
	m := dict.AsMap()
	k := fmt.Sprint(key)

	if current, ok := m[k]; ok {
		if list, err := collections.TryAsList(current); err == nil {
			m[k] = list.Append(value)
		} else {
			// Convert the current value into a list
			m[k] = dict.CreateList().Append(current, value)
		}
	} else {
		m[k] = dh.Convert(value)
	}
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

// Transpose returns a dictionary with values as keys and keys as values. The resulting dictionary
// is a dictionary where each key could contains single value or list of values if there are multiple matches.
func (dh DictHelper) Transpose(dict baseIDict) baseIDict {
	result := dh.CreateDictionary()
	for _, key := range dict.GetKeys().AsArray() {
		value := dict.Get(key)
		if list, err := collections.TryAsList(value); err == nil {
			// If the element is a list, we scan each element
			for _, value := range list.AsArray() {
				result.Add(value, key)
			}
		} else {
			result.Add(value, key)
		}
	}
	return result
}

// Type returns the actual type of object.
func (dh DictHelper) Type(dict baseIDict) str {
	return str(reflect.TypeOf(dict).Name())
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
func init() { collections.SetDictionaryHelper(baseDictHelper) }
