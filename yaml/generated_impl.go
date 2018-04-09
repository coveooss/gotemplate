// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package yaml

import "github.com/coveo/gotemplate/types"

// List implementation of IGenericList for yamlList
type List = yamlList
type yamlIList = types.IGenericList
type yamlList []interface{}

func (l *yamlList) Append(values ...interface{}) yamlIList { return yamlListHelper.Append(l, values...) }
func (l *yamlList) AsArray() *[]interface{}                { return (*[]interface{})(l) }
func (l yamlList) Cap() int                                { return cap(l) }
func (l yamlList) Capacity() int                           { return cap(l) }
func (l yamlList) Clone() yamlIList                        { return yamlListHelper.Clone(&l) }
func (l yamlList) Count() int                              { return len(l) }
func (l yamlList) Get(index int) interface{}               { return yamlListHelper.GetIndex(&l, index) }
func (l yamlList) Len() int                                { return len(l) }
func (l yamlList) Reverse() yamlIList                      { return yamlListHelper.Reverse(&l) }

func (l *yamlList) Set(i int, v interface{}) (yamlIList, error) {
	return yamlListHelper.SetIndex(l, i, v)
}

// Dictionary implementation of IDictionary for yamlDict
type Dictionary = yamlDict
type yamlIDict = types.IDictionary
type yamlDict map[string]interface{}

func (d yamlDict) AsMap() map[string]interface{}       { return (map[string]interface{})(d) }
func (d yamlDict) Count() int                          { return len(d) }
func (d yamlDict) Len() int                            { return len(d) }
func (d yamlDict) Clone(keys ...interface{}) yamlIDict { return yamlDictHelper.Clone(d, keys) }
func (d yamlDict) CreateList(args ...int) yamlIList    { return yamlHelper.CreateList(args...) }
func (d yamlDict) Flush(keys ...interface{}) yamlIDict { return yamlDictHelper.Flush(d, keys) }
func (d yamlDict) Get(key interface{}) interface{}     { return yamlDictHelper.Get(d, key) }
func (d yamlDict) Has(key interface{}) bool            { return yamlDictHelper.Has(d, key) }
func (d yamlDict) Keys() yamlIList                     { return yamlDictHelper.Keys(d) }
func (d yamlDict) KeysAsString() []string              { return yamlDictHelper.KeysAsString(d) }

func (d yamlDict) Default(key, defVal interface{}) interface{} {
	return yamlDictHelper.Default(d, key, defVal)
}

func (d yamlDict) Delete(key interface{}, otherKeys ...interface{}) (yamlIDict, error) {
	return yamlDictHelper.Delete(d, append([]interface{}{key}, otherKeys...))
}

func (d yamlDict) Merge(dict yamlIDict, otherDicts ...yamlIDict) yamlIDict {
	return yamlDictHelper.Merge(d, append([]yamlIDict{dict}, otherDicts...))
}

func (d yamlDict) Omit(key interface{}, otherKeys ...interface{}) yamlIDict {
	return yamlDictHelper.Omit(d, append([]interface{}{key}, otherKeys...))
}

func (d yamlDict) Set(key interface{}, v interface{}) yamlIDict {
	return yamlDictHelper.Set(d, key, v)
}

// Generic helpers to simplify physical implementation
func yamlListConvert(list yamlIList) yamlIList {
	array := yamlList(*list.AsArray())
	return &array
}
func yamlDictConvert(dict yamlIDict) yamlIDict { return yamlDict(dict.AsMap()) }

var yamlHelper = helperBase{ConvertList: yamlListConvert, ConvertDict: yamlDictConvert}
var yamlListHelper = helperList{BaseHelper: yamlHelper}
var yamlDictHelper = helperDict{BaseHelper: yamlHelper}