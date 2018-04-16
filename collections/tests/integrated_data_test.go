package tests

import (
	"reflect"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/collections"
	impl "github.com/coveo/gotemplate/collections/implementation"
	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/json"
	"github.com/coveo/gotemplate/yaml"
)

type dictionary = map[string]interface{}

var hclHelper = hcl.DictionaryHelper
var yamlHelper = yaml.DictionaryHelper
var jsonHelper = json.DictionaryHelper
var genHelper = impl.DictionaryHelper

func TestConvertData(t *testing.T) {
	var out1 interface{}
	type args struct {
		data string
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"Simple value", args{"a = 10", &out1}, map[string]interface{}{"a": 10}, false},
		{"YAML", args{"a: 10", &out1}, dictionary{"a": 10}, false},
		{"HCL", args{`a = 10 b = "Foo"`, &out1}, dictionary{"a": 10, "b": "Foo"}, false},
		{"JSON", args{`{ "a": 10, "b": "Foo" }`, &out1}, dictionary{"a": 10, "b": "Foo"}, false},
		{"Flexible", args{`a = 10 b = Foo`, &out1}, dictionary{"a": 10, "b": "Foo"}, false},
		{"No change", args{"NoChange", &out1}, nil, false},
		{"Invalid", args{"a = 'value", &out1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := collections.ConvertData(tt.args.data, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("ConvertData() error = %v, wantErr %v\n%v", err, tt.wantErr, reflect.ValueOf(tt.args.out).Elem())
			}
		})
	}
}

func TestToBash(t *testing.T) {
	type SubStruct struct {
		U int64
		I interface{}
	}
	type a struct {
		private int
		I       int
		F       float64
		S       string
		A       []interface{}
		M       dictionary
		SS      SubStruct
	}
	tests := []struct {
		name string
		args interface{}
		want interface{}
	}{
		{"Struct conversion", a{
			private: 0,
			I:       123,
			F:       1.23,
			S:       "123",
			A:       []interface{}{1, "2"},
			M: dictionary{
				"a": "a",
				"b": 2,
			},
			SS: SubStruct{64, "Foo"},
		}, strings.TrimSpace(collections.UnIndent(`
		declare -a A
		A=(1 2)
		F=1.23
		I=123
		declare -A M
		M=([a]=a [b]=2)
		S=123
		declare -A SS
		SS=([I]=Foo [U]=64)
		`))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collections.ToBash(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNativeRepresentation()\ngot : %q\nwant: %q", got, tt.want)
			}
		})
	}
}
