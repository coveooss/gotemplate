package implementation

import (
	"github.com/coveo/gotemplate/types"
)

// ListTypeName implementation of IGenericList for baseList
type ListTypeName = baseList
type baseIList = types.IGenericList
type baseList []interface{}

func (l *baseList) Append(values ...interface{}) baseIList { return baseListHelper.Append(l, values...) }
func (l *baseList) AsArray() *[]interface{}                { return (*[]interface{})(l) }
func (l baseList) Cap() int                                { return cap(l) }
func (l baseList) Capacity() int                           { return cap(l) }
func (l baseList) Clone() baseIList                        { return baseListHelper.Clone(&l) }
func (l baseList) Count() int                              { return len(l) }
func (l baseList) Get(index int) interface{}               { return baseListHelper.GetIndex(&l, index) }
func (l baseList) Len() int                                { return len(l) }
func (l baseList) Reverse() baseIList                      { return baseListHelper.Reverse(&l) }

func (l *baseList) Set(i int, v interface{}) (baseIList, error) {
	return baseListHelper.SetIndex(l, i, v)
}

// DictTypeName implementation of IDictionary for baseDict
type DictTypeName = baseDict
type baseIDict = types.IDictionary
type baseDict map[string]interface{}

func (d baseDict) AsMap() map[string]interface{}       { return (map[string]interface{})(d) }
func (d baseDict) Count() int                          { return len(d) }
func (d baseDict) Len() int                            { return len(d) }
func (d baseDict) Clone(keys ...interface{}) baseIDict { return baseDictHelper.Clone(d, keys) }
func (d baseDict) CreateList(args ...int) baseIList    { return baseHelper.CreateList(args...) }
func (d baseDict) Flush(keys ...interface{}) baseIDict { return baseDictHelper.Flush(d, keys) }
func (d baseDict) Get(key interface{}) interface{}     { return baseDictHelper.Get(d, key) }
func (d baseDict) Has(key interface{}) bool            { return baseDictHelper.Has(d, key) }
func (d baseDict) Keys() baseIList                     { return baseDictHelper.Keys(d) }
func (d baseDict) KeysAsString() []string              { return baseDictHelper.KeysAsString(d) }

func (d baseDict) Default(key, defVal interface{}) interface{} {
	return baseDictHelper.Default(d, key, defVal)
}

func (d baseDict) Delete(key interface{}, otherKeys ...interface{}) (baseIDict, error) {
	return baseDictHelper.Delete(d, append([]interface{}{key}, otherKeys...))
}

func (d baseDict) Merge(dict baseIDict, otherDicts ...baseIDict) baseIDict {
	return baseDictHelper.Merge(d, append([]baseIDict{dict}, otherDicts...))
}

func (d baseDict) Omit(key interface{}, otherKeys ...interface{}) baseIDict {
	return baseDictHelper.Omit(d, append([]interface{}{key}, otherKeys...))
}

func (d baseDict) Set(key interface{}, v interface{}) baseIDict {
	return baseDictHelper.Set(d, key, v)
}

// Generic helpers to simplify physical implementation
func baseListConvert(list baseIList) baseIList {
	array := baseList(*list.AsArray())
	return &array
}
func baseDictConvert(dict baseIDict) baseIDict { return baseDict(dict.AsMap()) }

var baseHelper = helperBase{ConvertList: baseListConvert, ConvertDict: baseDictConvert}
var baseListHelper = helperList{BaseHelper: baseHelper}
var baseDictHelper = helperDict{BaseHelper: baseHelper}
