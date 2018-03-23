package xml

import (
	"github.com/coveo/gotemplate/types"
)

// Dictionary is a specialized type for HCL Dictionary representation
type Dictionary = dict

type dict types.Dictionary
type iDict = types.IDictionary
type pDict = types.Dictionary

func (d dict) Clone(keys ...interface{}) iDict       { return dict(pDict(d).Clone().AsMap()) }
func (d dict) Omit(keys ...interface{}) iDict        { return dict(pDict(d).Omit().AsMap()) }
func (d dict) Get(key interface{}) interface{}       { return pDict(d).Get(key) }
func (d dict) Has(key interface{}) bool              { return pDict(d).Has(key) }
func (d dict) Delete(key interface{}) (iDict, error) { _, err := pDict(d).Delete(key); return d, err }
func (d dict) Len() int                              { return pDict(d).Len() }
func (d dict) Keys() iList                           { return list(pDict(d).Keys().(pList)) }
func (d dict) KeysAsString() []string                { return pDict(d).KeysAsString() }
func (d dict) AsMap() map[string]interface{}         { return pDict(d).AsMap() }
func (d dict) String() string                        { result, _ := Marshal(pDict(d)); return string(result) }
func (d dict) Merge(dicts ...iDict) (iDict, error)   { return pDict(d).Merge(dicts...) }
func (d dict) Set(key interface{}, value interface{}) iDict {
	return dict(pDict(d).Set(key, value).AsMap())
}
