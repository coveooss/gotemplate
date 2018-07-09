package hcl

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/collections/implementation"
	"github.com/coveo/gotemplate/errors"
	"github.com/hashicorp/hcl"
)

// Expose hcl public objects
var (
	Decode       = hcl.Decode
	DecodeObject = hcl.DecodeObject
	Parse        = hcl.Parse
	ParseBytes   = hcl.ParseBytes
	ParseString  = hcl.ParseString
)

func (l hclList) String() string {
	result, err := MarshalInternal(l.AsArray())
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(result))
}

func (d hclDict) String() string {
	result, err := Marshal(d.AsMap())
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(result))
}

var _ = func() int {
	collections.TypeConverters["hcl"] = Unmarshal
	return 0
}()

// Unmarshal adds support to single array and struct representation
func Unmarshal(bs []byte, out interface{}) (err error) {
	defer func() { err = errors.Trap(err, recover()) }()
	bs = bytes.TrimSpace(bs)

	if err = hcl.Unmarshal(bs, out); err != nil {
		bs = append([]byte("_="), bs...)
		var temp hclDict
		if errInternalHcl := hcl.Unmarshal(bs, &temp); errInternalHcl != nil {
			return err
		}
		err = nil
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(temp["_"]))
	}

	transform(out)
	return
}

// Load loads hcl file into variable
func Load(filename string) (result interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		err = Unmarshal(content, &result)
	}
	return
}

// Marshal serialize values to hcl format
func Marshal(value interface{}) ([]byte, error) { return MarshalIndent(value, "", "") }

// MarshalIndent serialize values to hcl format with indentation
func MarshalIndent(value interface{}, prefix, indent string) ([]byte, error) {
	result, err := marshalHCL(collections.ToNativeRepresentation(value), true, true, prefix, indent)
	return []byte(result), err
}

// MarshalInternal serialize values to hcl format for result used in outer hcl struct
func MarshalInternal(value interface{}) ([]byte, error) {
	result, err := marshalHCL(collections.ToNativeRepresentation(value), false, false, "", "")
	return []byte(result), err
}

// MarshalTFVars serialize values to hcl format (without hcl map format)
func MarshalTFVars(value interface{}) ([]byte, error) { return MarshalTFVarsIndent(value, "", "") }

// MarshalTFVarsIndent serialize values to hcl format with indentation (without hcl map format)
func MarshalTFVarsIndent(value interface{}, prefix, indent string) ([]byte, error) {
	result, err := marshalHCL(collections.ToNativeRepresentation(value), false, true, prefix, indent)
	return []byte(result), err
}

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
}

type helperBase = implementation.BaseHelper
type helperList = implementation.ListHelper
type helperDict = implementation.DictHelper

//go:generate genny -pkg=hcl -in=../collections/implementation/generic.go -out=generated_impl.go gen "ListTypeName=List DictTypeName=Dictionary base=hcl"
//go:generate genny -pkg=hcl -in=../collections/implementation/generic_test.go -out=generated_test.go gen "base=hcl"
