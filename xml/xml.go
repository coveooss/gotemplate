package xml

import (
	"encoding/xml"
	"math"
	"reflect"

	"github.com/coveooss/gotemplate/v3/collections/implementation"
)

// Expose xml public objects
var (
	CopyToken       = xml.CopyToken
	Escape          = xml.Escape
	EscapeText      = xml.EscapeText
	Marshal         = xml.Marshal
	MarshalIndent   = xml.MarshalIndent
	NewDecoder      = xml.NewDecoder
	NewEncoder      = xml.NewEncoder
	NativeUnmarshal = xml.Unmarshal
)

func (l xmlList) String() string { result, _ := Marshal(l.AsArray()); return string(result) }
func (d xmlDict) String() string { result, _ := Marshal(d.AsMap()); return string(result) }

func (l xmlList) PrettyPrint() string {
	result, _ := MarshalIndent(l.AsArray(), "", "  ")
	return string(result)
}
func (d xmlDict) PrettyPrint() string {
	result, _ := MarshalIndent(d.AsMap(), "", "  ")
	return string(result)
}

// TODO activate when Marshall will be enabled on maps
// var _ = func() int {
// 	collections.TypeConverters["xml"] = Unmarshal
// 	return 0
// }()

// Unmarshal calls the native Unmarshal but transform the results
// to returns Dictionary and GenerecList instead of go native collections.
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
		result = result.(xmlIDict).Native()
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
}

func transformElement(source interface{}) interface{} {
	if value, err := xmlHelper.TryAsDictionary(source); err == nil {
		for _, key := range value.KeysAsString() {
			value.Set(key, transformElement(value.Get(key)))
		}
		source = value
	} else if value, err := xmlHelper.TryAsList(source); err == nil {
		for i, sub := range value.AsArray() {
			value.Set(i, transformElement(sub))
		}
		source = value
	} else if value, ok := source.(float64); ok {
		// xml.Unmarshal returns all int values as float64
		if math.Floor(value) == value {
			source = int(value)
		}
	}
	return source
}

type (
	helperBase = implementation.BaseHelper
	helperList = implementation.ListHelper
	helperDict = implementation.DictHelper
)

var needConversionImpl = implementation.NeedConversion

//go:generate genny -pkg=xml -in=../collections/implementation/generic.go -out=generated_impl.go gen "ListTypeName=List DictTypeName=Dictionary base=xml"
//go:generate genny -pkg=xml -in=../collections/implementation/generic_test.go -out=generated_test.go gen "base=xml"
