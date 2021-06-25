package implementation

import (
	"strings"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/multilogger/errors"
)

// ListTypeName implementation of IGenericList for baseList
type ListTypeName = baseList
type baseIList = collections.IGenericList
type baseList []interface{}

var baseLower = strings.ToLower("base") // This is required because genny capitalize the type name in strings

func (l baseList) AsArray() []interface{}               { return []interface{}(l) }
func (l baseList) Cap() int                             { return cap(l) }
func (l baseList) Capacity() int                        { return cap(l) }
func (l baseList) Clone() baseIList                     { return baseListHelper.Clone(l) }
func (l baseList) Count() int                           { return len(l) }
func (l baseList) Create(args ...int) baseIList         { return baseListHelper.CreateList(args...) }
func (l baseList) CreateDict(args ...int) baseIDict     { return baseListHelper.CreateDictionary(args...) }
func (l baseList) Find(element interface{}) baseIList   { return baseListHelper.Find(l, element, false) }
func (l baseList) First() interface{}                   { return baseListHelper.GetIndexes(l, 0) }
func (l baseList) Get(indexes ...int) interface{}       { return baseListHelper.GetIndexes(l, indexes...) }
func (l baseList) GetKinds() baseIList                  { return baseListHelper.GetTypes(l, true) }
func (l baseList) GetTypes() baseIList                  { return baseListHelper.GetTypes(l, false) }
func (l baseList) Has(values ...interface{}) bool       { return l.Contains(values...) }
func (l baseList) HasStrict(values ...interface{}) bool { return l.ContainsStrict(values...) }
func (l baseList) Join(sep interface{}) str             { return l.StringArray().Join(sep) }
func (l baseList) Last() interface{}                    { return baseListHelper.GetIndexes(l, len(l)-1) }
func (l baseList) Len() int                             { return len(l) }
func (l baseList) New(args ...interface{}) baseIList    { return baseListHelper.NewList(args...) }
func (l baseList) Reverse() baseIList                   { return baseListHelper.Reverse(l) }
func (l baseList) RemoveEmpty() baseIList               { return baseListHelper.RemoveEmpty(l) }
func (l baseList) RemoveNil() baseIList                 { return baseListHelper.RemoveNil(l) }
func (l baseList) StringArray() strArray                { return baseListHelper.GetStringArray(l) }
func (l baseList) Strings() []string                    { return baseListHelper.GetStrings(l) }
func (l baseList) Type() str                            { return baseListHelper.Type(l) }
func (l baseList) TypeName() str                        { return str(baseLower) }
func (l baseList) Unique() baseIList                    { return baseListHelper.Unique(l) }

func (l baseList) GetHelpers() (collections.IDictionaryHelper, collections.IListHelper) {
	return baseDictHelper, baseListHelper
}

func (l baseList) Append(values ...interface{}) baseIList {
	return baseListHelper.Add(l, false, values...)
}

func (l baseList) Contains(values ...interface{}) bool {
	return baseListHelper.Contains(l, false, values...)
}

func (l baseList) ContainsStrict(values ...interface{}) bool {
	return baseListHelper.Contains(l, true, values...)
}

func (l baseList) FindStrict(element interface{}) baseIList {
	return baseListHelper.Find(l, element, true)
}

func (l baseList) Intersect(values ...interface{}) baseIList {
	return baseListHelper.Intersect(l, values...)
}

func (l baseList) Pop(indexes ...int) (interface{}, baseIList) {
	if len(indexes) == 0 {
		indexes = []int{len(l) - 1}
	}
	return l.Get(indexes...), l.Remove(indexes...)
}

func (l baseList) Prepend(values ...interface{}) baseIList {
	return baseListHelper.Add(l, true, values...)
}

func (l baseList) Remove(indexes ...int) baseIList {
	return baseListHelper.Remove(l, indexes...)
}

func (l baseList) Set(i int, v interface{}) (baseIList, error) {
	return baseListHelper.SetIndex(l, i, v)
}

func (l baseList) Union(values ...interface{}) baseIList {
	return baseListHelper.Add(l, false, values...).Unique()
}

func (l baseList) Without(values ...interface{}) baseIList {
	return baseListHelper.Without(l, values...)
}

// DictTypeName implementation of IDictionary for baseDict
type DictTypeName = baseDict
type baseIDict = collections.IDictionary
type baseDict map[string]interface{}

func (d baseDict) Add(key, v interface{}) baseIDict    { return baseDictHelper.Add(d, key, v) }
func (d baseDict) AsMap() map[string]interface{}       { return (map[string]interface{})(d) }
func (d baseDict) Clone(keys ...interface{}) baseIDict { return baseDictHelper.Clone(d, keys) }
func (d baseDict) Count() int                          { return len(d) }
func (d baseDict) Create(args ...int) baseIDict        { return baseListHelper.CreateDictionary(args...) }
func (d baseDict) CreateList(args ...int) baseIList    { return baseHelper.CreateList(args...) }
func (d baseDict) Flush(keys ...interface{}) baseIDict { return baseDictHelper.Flush(d, keys) }
func (d baseDict) Get(keys ...interface{}) interface{} { return baseDictHelper.Get(d, keys) }
func (d baseDict) GetKeys() baseIList                  { return baseDictHelper.GetKeys(d) }
func (d baseDict) GetKinds() baseIDict                 { return baseDictHelper.GetTypes(d, true) }
func (d baseDict) GetTypes() baseIDict                 { return baseDictHelper.GetTypes(d, false) }
func (d baseDict) GetValues() baseIList                { return baseDictHelper.GetValues(d) }
func (d baseDict) Has(keys ...interface{}) bool        { return baseDictHelper.Has(d, keys) }
func (d baseDict) KeysAsString() strArray              { return baseDictHelper.KeysAsString(d) }
func (d baseDict) Len() int                            { return len(d) }
func (d baseDict) Native() interface{}                 { return must(collections.MarshalGo(d)) }
func (d baseDict) Pop(keys ...interface{}) interface{} { return baseDictHelper.Pop(d, keys) }
func (d baseDict) Set(key, v interface{}) baseIDict    { return baseDictHelper.Set(d, key, v) }
func (d baseDict) Transpose() baseIDict                { return baseDictHelper.Transpose(d) }
func (d baseDict) Type() str                           { return baseDictHelper.Type(d) }
func (d baseDict) TypeName() str                       { return str(baseLower) }

func (d baseDict) GetHelpers() (collections.IDictionaryHelper, collections.IListHelper) {
	return baseDictHelper, baseListHelper
}

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

// Generic helpers to simplify physical implementation
func baseListConvert(list baseIList) baseIList { return baseList(list.AsArray()) }
func baseDictConvert(dict baseIDict) baseIDict { return baseDict(dict.AsMap()) }
func needConversion(object interface{}, strict bool) bool {
	return needConversionImpl(object, strict, "base")
}

var baseHelper = helperBase{ConvertList: baseListConvert, ConvertDict: baseDictConvert, NeedConversion: needConversion}
var baseListHelper = helperList{BaseHelper: baseHelper}
var baseDictHelper = helperDict{BaseHelper: baseHelper}

// DictionaryHelper gives public access to the basic dictionary functions
var DictionaryHelper collections.IDictionaryHelper = baseDictHelper

// GenericListHelper gives public access to the basic list functions
var GenericListHelper collections.IListHelper = baseListHelper

type (
	str      = collections.String
	strArray = collections.StringArray
)

var must = errors.Must
