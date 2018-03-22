package hcl

import (
	"github.com/coveo/gotemplate/types"
)

// List is a specialized type for HCL Generic List representation
type List = list

type list types.GenericList
type iList = types.IGenericList
type pList = types.GenericList

func (l list) String() string                                  { result, _ := MarshalInternal(pList(l)); return string(result) }
func (l list) Set(index int, value interface{}) (iList, error) { return pList(l).Set(index, value) }
func (l list) Get(index int) interface{}                       { return pList(l).Get(index) }
func (l list) Len() int                                        { return pList(l).Len() }
func (l list) AsList() []interface{}                           { return pList(l).AsList() }
func (l list) Clone() iList                                    { return list(pList(l).Clone().AsList()) }
