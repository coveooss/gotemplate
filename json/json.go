package json

import (
	"encoding/json"
	"math"
	"reflect"

	"github.com/coveo/gotemplate/types"
	"github.com/coveo/gotemplate/utils"
)

// Expose json public objects
var (
	Compact         = json.Compact
	HTMLEscape      = json.HTMLEscape
	Indent          = json.Indent
	Marshal         = json.Marshal
	MarshalIndent   = json.MarshalIndent
	NewDecoder      = json.NewDecoder
	NewEncoder      = json.NewEncoder
	NativeUnmarshal = json.Unmarshal
)

var _ = func() int {
	utils.TypeConverters["!json"] = Unmarshal
	return 0
}()

// Unmarshal calls the native Unmarshal but transform the results
// to returns Dictionary and GenerecList instead of go native types.
func Unmarshal(data []byte, out interface{}) (err error) {
	if err = NativeUnmarshal(data, out); err != nil {
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
	} else if value, ok := source.(float64); ok {
		// json.Unmarshal returns all int values as float64
		if math.Floor(value) == value {
			source = int(value)
		}
	}
	return source
}
