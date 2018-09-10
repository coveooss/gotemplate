// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package hcl

import "github.com/coveo/gotemplate/collections"

// List implementation of IGenericList for hclList
type List = hclList
type hclIList = collections.IGenericList
type hclList []interface{}

func (l hclList) AsArray() []interface{} { return []interface{}(l) }
func (l hclList) Cap() int               { return cap(l) }
func (l hclList) Capacity() int          { return cap(l) }
func (l hclList) Clone() hclIList        { return hclListHelper.Clone(l) }
func (l hclList) Contains(values ...interface{}) bool {
	return hclListHelper.Contains(l, values...)
}
func (l hclList) Count() int                  { return len(l) }
func (l hclList) Create(args ...int) hclIList { return hclListHelper.CreateList(args...) }
func (l hclList) CreateDict(args ...int) hclIDict {
	return hclListHelper.CreateDictionary(args...)
}
func (l hclList) Get(index int) interface{}        { return hclListHelper.GetIndex(l, index) }
func (l hclList) Len() int                         { return len(l) }
func (l hclList) New(args ...interface{}) hclIList { return hclListHelper.NewList(args...) }
func (l hclList) Reverse() hclIList                { return hclListHelper.Reverse(l) }
func (l hclList) Strings() []string                { return hclListHelper.GetStrings(l) }
func (l hclList) Unique() hclIList                 { return hclListHelper.Unique(l) }

func (l hclList) Append(values ...interface{}) hclIList {
	return hclListHelper.Add(l, false, values...)
}

func (l hclList) Prepend(values ...interface{}) hclIList {
	return hclListHelper.Add(l, true, values...)
}

func (l hclList) Set(i int, v interface{}) (hclIList, error) {
	return hclListHelper.SetIndex(l, i, v)
}

func (l hclList) Without(values ...interface{}) hclIList {
	return hclListHelper.Without(l, values...)
}

// Dictionary implementation of IDictionary for hclDict
type Dictionary = hclDict
type hclIDict = collections.IDictionary
type hclDict map[string]interface{}

func (d hclDict) Add(key, v interface{}) hclIDict    { return hclDictHelper.Add(d, key, v) }
func (d hclDict) AsMap() map[string]interface{}      { return (map[string]interface{})(d) }
func (d hclDict) Native() interface{}                { return collections.ToNativeRepresentation(d) }
func (d hclDict) Count() int                         { return len(d) }
func (d hclDict) Len() int                           { return len(d) }
func (d hclDict) Clone(keys ...interface{}) hclIDict { return hclDictHelper.Clone(d, keys) }
func (d hclDict) Create(args ...int) hclIDict        { return hclListHelper.CreateDictionary(args...) }
func (d hclDict) CreateList(args ...int) hclIList    { return hclHelper.CreateList(args...) }
func (d hclDict) Flush(keys ...interface{}) hclIDict { return hclDictHelper.Flush(d, keys) }
func (d hclDict) Get(key interface{}) interface{}    { return hclDictHelper.Get(d, key) }
func (d hclDict) Has(key interface{}) bool           { return hclDictHelper.Has(d, key) }
func (d hclDict) GetKeys() hclIList                  { return hclDictHelper.GetKeys(d) }
func (d hclDict) KeysAsString() []string             { return hclDictHelper.KeysAsString(d) }
func (d hclDict) GetValues() hclIList                { return hclDictHelper.GetValues(d) }
func (d hclDict) Set(key, v interface{}) hclIDict    { return hclDictHelper.Set(d, key, v) }
func (d hclDict) Transpose() hclIDict                { return hclDictHelper.Transpose(d) }

func (d hclDict) Default(key, defVal interface{}) interface{} {
	return hclDictHelper.Default(d, key, defVal)
}

func (d hclDict) Delete(key interface{}, otherKeys ...interface{}) (hclIDict, error) {
	return hclDictHelper.Delete(d, append([]interface{}{key}, otherKeys...))
}

func (d hclDict) Merge(dict hclIDict, otherDicts ...hclIDict) hclIDict {
	return hclDictHelper.Merge(d, append([]hclIDict{dict}, otherDicts...))
}

func (d hclDict) Omit(key interface{}, otherKeys ...interface{}) hclIDict {
	return hclDictHelper.Omit(d, append([]interface{}{key}, otherKeys...))
}

// Generic helpers to simplify physical implementation
func hclListConvert(list hclIList) hclIList { return hclList(list.AsArray()) }
func hclDictConvert(dict hclIDict) hclIDict { return hclDict(dict.AsMap()) }

var hclHelper = helperBase{ConvertList: hclListConvert, ConvertDict: hclDictConvert}
var hclListHelper = helperList{BaseHelper: hclHelper}
var hclDictHelper = helperDict{BaseHelper: hclHelper}

// DictionaryHelper gives public access to the basic dictionary functions
var DictionaryHelper collections.IDictionaryHelper = hclDictHelper

// GenericListHelper gives public access to the basic list functions
var GenericListHelper collections.IListHelper = hclListHelper
