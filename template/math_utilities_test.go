package template

import (
	"reflect"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/collections/implementation"
	"github.com/coveooss/gotemplate/v3/json"
)

type a = []interface{}
type l = implementation.ListTypeName
type j = json.List

func Test_convertArgs(t *testing.T) {
	// t.Parallel()
	collections.DictionaryHelper = implementation.DictionaryHelper
	collections.ListHelper = implementation.GenericListHelper
	tests := []struct {
		name string
		arg1 interface{}
		args a
		want iList
	}{
		{"Nil", nil, nil, l{}},
		{"Single int", 5, nil, l{5}},
		{"Two int", 2, a{3}, l{2, 3}},
		{"First nil", nil, a{3}, l{3}},
		{"nil+values", nil, a{3, 4, 5}, l{3, 4, 5}},
		{"nil+array", nil, a{a{3, 4, 5}}, l{3, 4, 5}},
		{"array+nil", a{3, 4, 5}, nil, l{3, 4, 5}},
		{"json+empty", j{3, 4, 5}, a{}, j{3, 4, 5}},
		{"nil+json exp", nil, j{3, 4, 5}, l{3, 4, 5}},
		{"nil+json list", nil, a{j{3, 4, 5}}, j{3, 4, 5}},
		{"value+json exp", 2, j{3, 4}, l{2, 3, 4}},
		{"value+json list", 2, a{j{3, 4}}, l{2, j{3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertArgs(tt.arg1, tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertArgs() = %[1]v (%[1]T), want %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_toListOfFloats(t *testing.T) {
	// t.Parallel()
	collections.DictionaryHelper = implementation.DictionaryHelper
	collections.ListHelper = implementation.GenericListHelper
	tests := []struct {
		name       string
		values     iList
		wantResult iList
		wantErr    bool
	}{
		{"Nil", nil, l{}, false},
		{"Empty", l{}, l{}, false},
		{"Array of int", l{1, 2, 3}, l{float64(1), float64(2), float64(3)}, false},
		{"Array of string", l{"1", "2", "3"}, l{float64(1), float64(2), float64(3)}, false},
		{"Invalid value", l{"1", "bad"}, nil, true},
		{"Json list", l{j{1.2, 2.3}}, j{1.2, 2.3}, false},
		{"Two json list", l{j{1.2, 2.3}, j{3.4}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := toListOfFloats(tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("toArrayOfFloats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("toArrayOfFloats() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_asFloats(t *testing.T) {
	t.Parallel()
	type a = []interface{}
	type l = implementation.ListTypeName
	type j = json.List
	tests := []struct {
		name       string
		values     iList
		wantResult []float64
		wantErr    bool
	}{
		{"Nil", nil, []float64{}, false},
		{"Empty", l{}, []float64{}, false},
		{"Array of int", l{1, 2, 3}, []float64{1, 2, 3}, false},
		{"Array of string", l{"1", "2", "3"}, []float64{1, 2, 3}, false},
		{"Invalid value", l{"1", "bad"}, nil, true},
		{"Json list", l{j{1.2, 2.3}}, []float64{1.2, 2.3}, false},
		{"Two json list", l{j{1.2, 2.3}, j{3.4}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := asFloats(tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("toArrayOfFloats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("toArrayOfFloats() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
