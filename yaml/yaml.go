package yaml

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/collections/implementation"
	"gopkg.in/yaml.v2"
)

// Expose yaml public objects
var (
	Marshal               = yaml.Marshal
	NewDecoder            = yaml.NewDecoder
	NewEncoder            = yaml.NewEncoder
	NativeUnmarshal       = yaml.Unmarshal
	NativeUnmarshalStrict = yaml.UnmarshalStrict
)

func (l yamlList) String() string      { result, _ := Marshal(l.AsArray()); return string(result) }
func (d yamlDict) String() string      { result, _ := Marshal(d.AsMap()); return string(result) }
func (l yamlList) PrettyPrint() string { return l.String() }
func (d yamlDict) PrettyPrint() string { return d.String() }

var _ = func() int {
	collections.TypeConverters["yaml"] = Unmarshal
	return 0
}()

// Unmarshal calls the native Unmarshal but transform the results
// to returns Dictionary and GenericList instead of go native collections.
func Unmarshal(data []byte, out interface{}) (err error) {
	// Yaml does not support tab, so we replace tabs by spaces if there are
	data = bytes.Replace(data, []byte("\t"), []byte("    "), -1)
	if err = NativeUnmarshal(data, out); err != nil {
		return
	}
	transform(out)
	return
}

// UnmarshalStrict calls the native UnmarshalStrict but transform the results
// to returns Dictionary and GenericList instead of go native collections.
func UnmarshalStrict(data []byte, out interface{}) (err error) {
	// Yaml does not support tab, so we replace tabs by spaces if there are
	data = bytes.Replace(data, []byte("\t"), []byte("    "), -1)
	if err = NativeUnmarshalStrict(data, out); err != nil {
		return
	}
	transform(out)
	return
}

func transform(out interface{}) {
	result := transformElement(reflect.ValueOf(out).Elem().Interface())
	if _, isMap := out.(*map[string]interface{}); isMap {
		// If the result is expected to be map[string]interface{}, we convert it back from internal dict type.
		result = result.(yamlIDict).Native()
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
}

func transformElement(source interface{}) interface{} {
	switch value := source.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(value))
		for key, val := range value {
			result[fmt.Sprint(key)] = val
		}
		source = result
	}

	if value, err := yamlHelper.TryAsDictionary(source); err == nil {
		for _, key := range value.KeysAsString() {
			value.Set(key, transformElement(value.Get(key)))
		}
		source = value
	} else if value, err := yamlHelper.TryAsList(source); err == nil {
		for i, sub := range value.AsArray() {
			value.Set(i, transformElement(sub))
		}
		source = value
	}
	return source
}

type (
	helperBase = implementation.BaseHelper
	helperList = implementation.ListHelper
	helperDict = implementation.DictHelper
)

var needConversionImpl = implementation.NeedConversion

//go:generate genny -pkg=yaml -in=../collections/implementation/generic.go -out=generated_impl.go gen "ListTypeName=List DictTypeName=Dictionary base=yaml"
//go:generate genny -pkg=yaml -in=../collections/implementation/generic_test.go -out=generated_test.go gen "base=yaml"
