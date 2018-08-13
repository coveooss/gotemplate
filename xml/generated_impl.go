// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package xml

import "github.com/coveo/gotemplate/collections"

// List implementation of IGenericList for xmlList
type List = xmlList
type xmlIList = collections.IGenericList
type xmlList []interface{}

func (l xmlList) Append(values ...interface{}) xmlIList { return xmlListHelper.Append(l, values...) }
func (l xmlList) AsArray() []interface{}                { return []interface{}(l) }
func (l xmlList) Cap() int                              { return cap(l) }
func (l xmlList) Capacity() int                         { return cap(l) }
func (l xmlList) Clone() xmlIList                       { return xmlListHelper.Clone(l) }
func (l xmlList) Count() int                            { return len(l) }
func (l xmlList) Create(args ...int) xmlIList           { return xmlListHelper.CreateList(args...) }
func (l xmlList) Get(index int) interface{}             { return xmlListHelper.GetIndex(l, index) }
func (l xmlList) New(args ...interface{}) xmlIList      { return xmlListHelper.NewList(args...) }
func (l xmlList) Len() int                              { return len(l) }
func (l xmlList) Reverse() xmlIList                     { return xmlListHelper.Reverse(l) }
func (l xmlList) Strings() []string                     { return xmlListHelper.GetStrings(l) }

func (l xmlList) Set(i int, v interface{}) (xmlIList, error) {
	return xmlListHelper.SetIndex(l, i, v)
}

// Dictionary implementation of IDictionary for xmlDict
type Dictionary = xmlDict
type xmlIDict = collections.IDictionary
type xmlDict map[string]interface{}

func (d xmlDict) AsMap() map[string]interface{}      { return (map[string]interface{})(d) }
func (d xmlDict) Native() interface{}                { return collections.ToNativeRepresentation(d) }
func (d xmlDict) Count() int                         { return len(d) }
func (d xmlDict) Len() int                           { return len(d) }
func (d xmlDict) Clone(keys ...interface{}) xmlIDict { return xmlDictHelper.Clone(d, keys) }
func (d xmlDict) CreateList(args ...int) xmlIList    { return xmlHelper.CreateList(args...) }
func (d xmlDict) Flush(keys ...interface{}) xmlIDict { return xmlDictHelper.Flush(d, keys) }
func (d xmlDict) Get(key interface{}) interface{}    { return xmlDictHelper.Get(d, key) }
func (d xmlDict) Has(key interface{}) bool           { return xmlDictHelper.Has(d, key) }
func (d xmlDict) Keys() xmlIList                     { return xmlDictHelper.Keys(d) }
func (d xmlDict) KeysAsString() []string             { return xmlDictHelper.KeysAsString(d) }
func (d xmlDict) Values() xmlIList                   { return xmlDictHelper.Values(d) }

func (d xmlDict) Default(key, defVal interface{}) interface{} {
	return xmlDictHelper.Default(d, key, defVal)
}

func (d xmlDict) Delete(key interface{}, otherKeys ...interface{}) (xmlIDict, error) {
	return xmlDictHelper.Delete(d, append([]interface{}{key}, otherKeys...))
}

func (d xmlDict) Merge(dict xmlIDict, otherDicts ...xmlIDict) xmlIDict {
	return xmlDictHelper.Merge(d, append([]xmlIDict{dict}, otherDicts...))
}

func (d xmlDict) Omit(key interface{}, otherKeys ...interface{}) xmlIDict {
	return xmlDictHelper.Omit(d, append([]interface{}{key}, otherKeys...))
}

func (d xmlDict) Set(key interface{}, v interface{}) xmlIDict {
	return xmlDictHelper.Set(d, key, v)
}

// Generic helpers to simplify physical implementation
func xmlListConvert(list xmlIList) xmlIList { return xmlList(list.AsArray()) }
func xmlDictConvert(dict xmlIDict) xmlIDict { return xmlDict(dict.AsMap()) }

var xmlHelper = helperBase{ConvertList: xmlListConvert, ConvertDict: xmlDictConvert}
var xmlListHelper = helperList{BaseHelper: xmlHelper}
var xmlDictHelper = helperDict{BaseHelper: xmlHelper}

// DictionaryHelper gives public access to the basic dictionary functions
var DictionaryHelper collections.IDictionaryHelper = xmlDictHelper

// GenericListHelper gives public access to the basic list functions
var GenericListHelper collections.IListHelper = xmlListHelper
