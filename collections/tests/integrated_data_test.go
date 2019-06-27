package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	impl "github.com/coveooss/gotemplate/v3/collections/implementation"
	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/yaml"
	"github.com/stretchr/testify/assert"
)

type dictionary = map[string]interface{}

var hclHelper = hcl.DictionaryHelper
var yamlHelper = yaml.DictionaryHelper
var jsonHelper = json.DictionaryHelper
var genHelper = impl.DictionaryHelper

func TestConvertData(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    interface{}
		wantErr error
	}{
		{"Simple value", "a = 10", map[string]interface{}{"a": 10}, nil},
		{"YAML", "a: 10", dictionary{"a": 10}, nil},
		{"HCL", `a = 10 b = "Foo"`, dictionary{"a": 10, "b": "Foo"}, nil},
		{"JSON", `{ "a": 10, "b": "Foo" }`, dictionary{"a": 10, "b": "Foo"}, nil},
		{"Flexible", `a = 10 b = Foo`, dictionary{"a": 10, "b": "Foo"}, nil},
		{"No change", "NoChange", "NoChange", nil},
		{"Invalid", "a = 'value", nil, fmt.Errorf("Trying !json: invalid character 'a' looking for beginning of value\nTrying hcl: At 1:5: illegal char")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out interface{}
			err := collections.ConvertData(tt.data, &out)
			assert.EqualValues(t, tt.want, out)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
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
			assert.Equal(t, tt.want, collections.ToBash(tt.args))
		})
	}
}
