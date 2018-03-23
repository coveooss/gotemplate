package hcl

import (
	"bytes"
	"io/ioutil"
	"reflect"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
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

var _ = func() int {
	utils.TypeConverters["hcl"] = Unmarshal
	return 0
}()

// Unmarshal adds support to single array and struct representation
func Unmarshal(bs []byte, out interface{}) (err error) {
	defer func() { err = errors.Trap(err, recover()) }()
	bs = bytes.TrimSpace(bs)

	if err = hcl.Unmarshal(bs, out); err != nil {
		bs = append([]byte("_="), bs...)
		var temp dict
		if err2 := hcl.Unmarshal(bs, &temp); err2 != nil {
			return err
		}
		err = nil
		reflect.ValueOf(out).Elem().Set(reflect.ValueOf(temp["_"]))
	}
	result := flatten(reflect.ValueOf(out).Elem().Interface())

	if _, isMap := out.(*map[string]interface{}); isMap {
		// If the result is expected to be map[string]interface{}, we convert it back from internal dict type.
		result = result.(dict).AsMap()
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
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
	result, err := marshalHCL(utils.ToNativeRepresentation(value), true, true, prefix, indent)
	return []byte(result), err
}

// MarshalInternal serialize values to hcl format for result used in outer hcl struct
func MarshalInternal(value interface{}) ([]byte, error) {
	result, err := marshalHCL(utils.ToNativeRepresentation(value), false, false, "", "")
	return []byte(result), err
}

// MarshalTFVars serialize values to hcl format (without hcl map format)
func MarshalTFVars(value interface{}) ([]byte, error) { return MarshalTFVarsIndent(value, "", "") }

// MarshalTFVarsIndent serialize values to hcl format with indentation (without hcl map format)
func MarshalTFVarsIndent(value interface{}, prefix, indent string) ([]byte, error) {
	result, err := marshalHCL(utils.ToNativeRepresentation(value), false, true, prefix, indent)
	return []byte(result), err
}

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
}
