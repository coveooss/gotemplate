package yaml

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/coveo/gotemplate/types"
	"github.com/coveo/gotemplate/utils"
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

var _ = func() int {
	utils.TypeConverters["yaml"] = Unmarshal
	return 0
}()

// Unmarshal calls the native Unmarshal but transform the results
// to returns Dictionary and GenerecList instead of go native types.
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
// to returns Dictionary and GenerecList instead of go native types.
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
		result = result.(dict).AsMap()
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

	if value, err := types.AsDictionary(source); err == nil {
		for _, key := range value.Keys().AsList() {
			value.Set(key, transformElement(value.Get(key)))
		}
		source = dict(value.AsMap())
	} else if value, err := types.AsGenericList(source); err == nil {
		for i, sub := range value.AsList() {
			value.Set(i, transformElement(sub))
		}
		source = list(value.AsList())
	}
	return source
}
