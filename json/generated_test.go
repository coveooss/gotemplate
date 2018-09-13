// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package json

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/errors"
)

var strFixture = jsonList(jsonListHelper.NewStringList(strings.Split("Hello World, I'm Foo Bar!", " ")...).AsArray())

func Test_list_Append(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      jsonIList
		values []interface{}
		want   jsonIList
	}{
		{"Empty", jsonList{}, []interface{}{1, 2, 3}, jsonList{1, 2, 3}},
		{"List of int", jsonList{1, 2, 3}, []interface{}{4, 5}, jsonList{1, 2, 3, 4, 5}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, jsonList{"Hello", "World,", "I'm", "Foo", "Bar!", "That's all folks!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Append(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Append():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Prepend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      jsonIList
		values []interface{}
		want   jsonIList
	}{
		{"Empty", jsonList{}, []interface{}{1, 2, 3}, jsonList{1, 2, 3}},
		{"List of int", jsonList{1, 2, 3}, []interface{}{4, 5}, jsonList{4, 5, 1, 2, 3}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, jsonList{"That's all folks!", "Hello", "World,", "I'm", "Foo", "Bar!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Prepend(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Prepend():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_AsArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want []interface{}
	}{
		{"Empty List", jsonList{}, []interface{}{}},
		{"List of int", jsonList{1, 2, 3}, []interface{}{1, 2, 3}},
		{"List of string", strFixture, []interface{}{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.AsArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.AsList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_JsonList_Strings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want []string
	}{
		{"Empty List", jsonList{}, []string{}},
		{"List of int", jsonList{1, 2, 3}, []string{"1", "2", "3"}},
		{"List of string", strFixture, []string{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Strings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Strings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_list_Capacity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonIList
		want int
	}{
		{"Empty List with 100 spaces", jsonListHelper.CreateList(0, 100), 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Capacity(); got != tt.want {
				t.Errorf("JsonList.Capacity() = %v, want %v", got, tt.want)
			}
			if tt.l.Capacity() != tt.l.Cap() {
				t.Errorf("Cap and Capacity return different values")
			}
		})
	}
}

func Test_list_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want jsonIList
	}{
		{"Empty List", jsonList{}, jsonList{}},
		{"List of int", jsonList{1, 2, 3}, jsonList{1, 2, 3}},
		{"List of string", strFixture, jsonList{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		l     jsonList
		index int
		want  interface{}
	}{
		{"Empty List", jsonList{}, 0, nil},
		{"Negative index", jsonList{}, -1, nil},
		{"List of int", jsonList{1, 2, 3}, 0, 1},
		{"List of string", strFixture, 1, "World,"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Get(tt.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_list_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want int
	}{
		{"Empty List", jsonList{}, 0},
		{"List of int", jsonList{1, 2, 3}, 3},
		{"List of string", strFixture, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Len(); got != tt.want {
				t.Errorf("JsonList.Len() = %v, want %v", got, tt.want)
			}
			if tt.l.Len() != tt.l.Count() {
				t.Errorf("Len and Count return different values")
			}
		})
	}
}

func Test_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []int
		want    jsonIList
		wantErr bool
	}{
		{"Empty", nil, jsonList{}, false},
		{"With nil elements", []int{10}, make(jsonList, 10), false},
		{"With capacity", []int{0, 10}, make(jsonList, 0, 10), false},
		{"Too much args", []int{0, 10, 1}, nil, true},
	}
	for _, tt := range tests {
		var got jsonIList
		var err error
		func() {
			defer func() { err = errors.Trap(err, recover()) }()
			got = jsonListHelper.CreateList(tt.args...)
		}()
		if (err != nil) != tt.wantErr {
			t.Errorf("CreateList() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if err != nil {
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CreateList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
		}
		if got.Capacity() != tt.want.Cap() {
			t.Errorf("CreateList() capacity:\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got.Cap(), tt.want.Capacity())
		}
	}
}

func Test_list_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []int
		want jsonIList
	}{
		{"Empty", nil, nil, jsonList{}},
		{"Existing List", jsonList{1, 2}, nil, jsonList{}},
		{"With Empty spaces", jsonList{1, 2}, []int{5}, jsonList{nil, nil, nil, nil, nil}},
		{"With Capacity", jsonList{1, 2}, []int{0, 5}, jsonListHelper.CreateList(0, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.Create(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Create():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if got.Capacity() != tt.want.Capacity() {
				t.Errorf("JsonList.Create() capacity:\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got.Capacity(), tt.want.Capacity())
			}
		})
	}
}

func Test_list_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []interface{}
		want jsonIList
	}{
		{"Empty", nil, nil, jsonList{}},
		{"Existing List", jsonList{1, 2}, nil, jsonList{}},
		{"With elements", jsonList{1, 2}, []interface{}{3, 4, 5}, jsonList{3, 4, 5}},
		{"With strings", jsonList{1, 2}, []interface{}{"Hello", "World"}, jsonList{"Hello", "World"}},
		{"With nothing", jsonList{1, 2}, []interface{}{}, jsonList{}},
		{"With nil", jsonList{1, 2}, nil, jsonList{}},
		{"Adding array", jsonList{1, 2}, []interface{}{jsonList{3, 4}}, jsonList{3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.New(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Create():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_CreateDict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       jsonList
		args    []int
		want    jsonIDict
		wantErr bool
	}{
		{"Empty", nil, nil, jsonDict{}, false},
		{"With capacity", nil, []int{10}, jsonDict{}, false},
		{"With too much parameter", nil, []int{10, 1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got jsonIDict
			var err error
			func() {
				defer func() { err = errors.Trap(err, recover()) }()
				got = tt.l.CreateDict(tt.args...)
			}()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.CreateDict():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonList.CreateDict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_list_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []interface{}
		want bool
	}{
		{"Empty List", nil, []interface{}{}, false},
		{"Search nothing", jsonList{1}, nil, true},
		{"Search nothing 2", jsonList{1}, []interface{}{}, true},
		{"Not there", jsonList{1}, []interface{}{2}, false},
		{"Included", jsonList{1, 2}, []interface{}{2}, true},
		{"Partially there", jsonList{1, 2}, []interface{}{2, 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Contains(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Contains():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Intersect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []interface{}
		want jsonList
	}{
		{"Empty List", nil, []interface{}{}, jsonList{}},
		{"Intersect nothing", jsonList{1}, nil, jsonList{}},
		{"Intersect nothing 2", jsonList{1}, []interface{}{}, jsonList{}},
		{"Not there", jsonList{1}, []interface{}{2}, jsonList{}},
		{"Included", jsonList{1, 2}, []interface{}{2}, jsonList{2}},
		{"Partially there", jsonList{1, 2}, []interface{}{2, 3}, jsonList{2}},
		{"With duplicates", jsonList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{3, 4, 5, 6, 7, 8, 7, 6, 5, 5, 4, 3}, jsonList{3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Intersect(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Intersect():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Union(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []interface{}
		want jsonList
	}{
		{"Empty List", nil, []interface{}{}, jsonList{}},
		{"Intersect nothing", jsonList{1}, nil, jsonList{1}},
		{"Intersect nothing 2", jsonList{1}, []interface{}{}, jsonList{1}},
		{"Not there", jsonList{1}, []interface{}{2}, jsonList{1, 2}},
		{"Included", jsonList{1, 2}, []interface{}{2}, jsonList{1, 2}},
		{"Partially there", jsonList{1, 2}, []interface{}{2, 3}, jsonList{1, 2, 3}},
		{"With duplicates", jsonList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{8, 7, 6, 5, 6, 7, 8}, jsonList{1, 2, 3, 4, 5, 8, 7, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Union(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Union():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Without(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		args []interface{}
		want jsonList
	}{
		{"Empty List", nil, []interface{}{}, jsonList{}},
		{"Remove nothing", jsonList{1}, nil, jsonList{1}},
		{"Remove nothing 2", jsonList{1}, []interface{}{}, jsonList{1}},
		{"Not there", jsonList{1}, []interface{}{2}, jsonList{1}},
		{"Included", jsonList{1, 2}, []interface{}{2}, jsonList{1}},
		{"Partially there", jsonList{1, 2}, []interface{}{2, 3}, jsonList{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Without(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Without():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Unique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want jsonList
	}{
		{"Empty List", nil, jsonList{}},
		{"Remove nothing", jsonList{1}, jsonList{1}},
		{"Duplicates following", jsonList{1, 1, 2, 3}, jsonList{1, 2, 3}},
		{"Duplicates not following", jsonList{1, 2, 3, 1, 2, 3, 4}, jsonList{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Unique(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Unique():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}
func Test_list_Reverse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    jsonList
		want jsonIList
	}{
		{"Empty List", jsonList{}, jsonList{}},
		{"List of int", jsonList{1, 2, 3}, jsonList{3, 2, 1}},
		{"List of string", strFixture, jsonList{"Bar!", "Foo", "I'm", "World,", "Hello"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			if got := l.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Reverse():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		i int
		v interface{}
	}
	tests := []struct {
		name    string
		l       jsonIList
		args    args
		want    jsonIList
		wantErr bool
	}{
		{"Empty", jsonList{}, args{2, 1}, jsonList{nil, nil, 1}, false},
		{"List of int", jsonList{1, 2, 3}, args{0, 10}, jsonList{10, 2, 3}, false},
		{"List of string", strFixture, args{2, "You're"}, jsonList{"Hello", "World,", "You're", "Foo", "Bar!"}, false},
		{"Negative", jsonList{}, args{-1, "negative value"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Clone().Set(tt.args.i, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonList.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonList.Set():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

var mapFixture = map[string]interface{}{
	"int":     123,
	"float":   1.23,
	"string":  "Foo bar",
	"list":    []interface{}{1, "two"},
	"listInt": []int{1, 2, 3},
	"map": map[string]interface{}{
		"sub1": 1,
		"sub2": "two",
	},
	"mapInt": map[int]interface{}{
		1: 1,
		2: "two",
	},
}

var dictFixture = jsonDict(jsonDictHelper.AsDictionary(mapFixture).AsMap())

func dumpKeys(t *testing.T, d1, d2 jsonIDict) {
	t.Parallel()

	for key := range d1.AsMap() {
		v1, v2 := d1.Get(key), d2.Get(key)
		if reflect.DeepEqual(v1, v2) {
			continue
		}
		t.Logf("'%[1]v' = %[2]v (%[2]T) vs %[3]v (%[3]T)", key, v1, v2)
	}
}

func Test_dict_AsMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		want map[string]interface{}
	}{
		{"Nil", nil, nil},
		{"Empty", jsonDict{}, map[string]interface{}{}},
		{"Map", dictFixture, map[string]interface{}(dictFixture)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.AsMap():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		keys []interface{}
		want jsonIDict
	}{
		{"Nil", nil, nil, jsonDict{}},
		{"Empty", jsonDict{}, nil, jsonDict{}},
		{"Map", dictFixture, nil, dictFixture},
		{"Map with Fields", dictFixture, []interface{}{"int", "list"}, jsonDict(dictFixture).Omit("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Clone(tt.keys...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}

			// Ensure that the copy is distinct from the original
			got.Set("NewFields", "Test")
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("Should be different:\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if !got.Has("NewFields") || !reflect.DeepEqual(got.Get("NewFields"), "Test") {
				t.Errorf("Element has not been added")
			}
			if got.Len() != tt.want.Count()+1 {
				t.Errorf("Len and Count don't return the same value")
			}
		})
	}
}

func Test_JsonDict_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		d            jsonDict
		args         []int
		want         jsonIList
		wantLen      int
		wantCapacity int
	}{
		{"Nil", nil, nil, jsonList{}, 0, 0},
		{"Empty", jsonDict{}, nil, jsonList{}, 0, 0},
		{"Map", dictFixture, nil, jsonList{}, 0, 0},
		{"Map with size", dictFixture, []int{3}, jsonList{nil, nil, nil}, 3, 3},
		{"Map with capacity", dictFixture, []int{0, 10}, jsonList{}, 0, 10},
		{"Map with size&capacity", dictFixture, []int{3, 10}, jsonList{nil, nil, nil}, 3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.CreateList(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.CreateList() = %v, want %v", got, tt.want)
			}
			if got.Len() != tt.wantLen || got.Cap() != tt.wantCapacity {
				t.Errorf("JsonDict.CreateList() size: %d, %d vs %d, %d", got.Len(), got.Cap(), tt.wantLen, tt.wantCapacity)
			}
		})
	}
}

func Test_dict_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       jsonDict
		args    []int
		want    jsonIDict
		wantErr bool
	}{
		{"Empty", nil, nil, jsonDict{}, false},
		{"With capacity", nil, []int{10}, jsonDict{}, false},
		{"With too much parameter", nil, []int{10, 1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got jsonIDict
			var err error
			func() {
				defer func() { err = errors.Trap(err, recover()) }()
				got = tt.d.Create(tt.args...)
			}()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Create():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonList.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_dict_Default(t *testing.T) {
	t.Parallel()

	type args struct {
		key    interface{}
		defVal interface{}
	}
	tests := []struct {
		name string
		d    jsonDict
		args args
		want interface{}
	}{
		{"Empty", nil, args{"Foo", "Bar"}, "Bar"},
		{"Map int", dictFixture, args{"int", 1}, 123},
		{"Map float", dictFixture, args{"float", 1}, 1.23},
		{"Map Non existant", dictFixture, args{"Foo", "Bar"}, "Bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Default(tt.args.key, tt.args.defVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_Delete(t *testing.T) {
	t.Parallel()

	type args struct {
		key  interface{}
		keys []interface{}
	}
	tests := []struct {
		name    string
		d       jsonDict
		args    args
		want    jsonIDict
		wantErr bool
	}{
		{"Empty", nil, args{}, jsonDict{}, true},
		{"Map", dictFixture, args{}, dictFixture, true},
		{"Non existant key", dictFixture, args{"Test", nil}, dictFixture, true},
		{"Map with keys", dictFixture, args{"int", []interface{}{"list"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), false},
		{"Map with keys + non existant", dictFixture, args{"int", []interface{}{"list", "Test"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got, err := d.Delete(tt.args.key, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonDict.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Delete():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}

func Test_dict_Flush(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		keys []interface{}
		want jsonIDict
	}{
		{"Empty", nil, nil, jsonDict{}},
		{"Map", dictFixture, nil, jsonDict{}},
		{"Non existant key", dictFixture, []interface{}{"Test"}, dictFixture},
		{"Map with keys", dictFixture, []interface{}{"int", "list"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
		{"Map with keys + non existant", dictFixture, []interface{}{"int", "list", "Test"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Flush(tt.keys...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Flush():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
			if !reflect.DeepEqual(d, got) {
				t.Errorf("Should be equal after: %v, want %v", d, got)
				dumpKeys(t, d, got)
			}
		})
	}
}

func Test_dict_Keys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		want jsonIList
	}{
		{"Empty", nil, jsonList{}},
		{"Map", dictFixture, jsonList{"float", "int", "list", "listInt", "map", "mapInt", "string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.GetKeys():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_KeysAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		want []string
	}{
		{"Empty", nil, []string{}},
		{"Map", dictFixture, []string{"float", "int", "list", "listInt", "map", "mapInt", "string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.KeysAsString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.KeysAsString():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Merge(t *testing.T) {
	t.Parallel()

	adding1 := jsonDict{
		"int":        1000,
		"Add1Int":    1,
		"Add1String": "string",
	}
	adding2 := jsonDict{
		"Add2Int":    1,
		"Add2String": "string",
		"map": map[string]interface{}{
			"sub1":   2,
			"newVal": "NewValue",
		},
	}
	type args struct {
		jsonDict jsonIDict
		dicts    []jsonIDict
	}
	tests := []struct {
		name string
		d    jsonDict
		args args
		want jsonIDict
	}{
		{"Empty", nil, args{nil, []jsonIDict{}}, jsonDict{}},
		{"Add map to empty", nil, args{dictFixture, []jsonIDict{}}, dictFixture},
		{"Add map to same map", dictFixture, args{dictFixture, []jsonIDict{}}, dictFixture},
		{"Add empty to map", dictFixture, args{nil, []jsonIDict{}}, dictFixture},
		{"Add new1 to map", dictFixture, args{adding1, []jsonIDict{}}, dictFixture.Clone().Merge(adding1)},
		{"Add new2 to map", dictFixture, args{adding2, []jsonIDict{}}, dictFixture.Clone().Merge(adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []jsonIDict{adding2}}, dictFixture.Clone().Merge(adding1, adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []jsonIDict{adding2}}, dictFixture.Clone().Merge(adding1).Merge(adding2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Merge(tt.args.jsonDict, tt.args.dicts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Merge():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}

func Test_dict_Values(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		want jsonIList
	}{
		{"Empty", nil, jsonList{}},
		{"Map", dictFixture, jsonList{1.23, 123, jsonList{1, "two"}, jsonList{1, 2, 3}, jsonDict{"sub1": 1, "sub2": "two"}, jsonDict{"1": 1, "2": "two"}, "Foo bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.GetValues():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Add(t *testing.T) {
	t.Parallel()

	type args struct {
		key interface{}
		v   interface{}
	}
	tests := []struct {
		name string
		d    jsonDict
		args args
		want jsonIDict
	}{
		{"Empty", nil, args{"A", 1}, jsonDict{"A": 1}},
		{"With element", jsonDict{"A": 1}, args{"A", 2}, jsonDict{"A": jsonList{1, 2}}},
		{"With element, another value", jsonDict{"A": 1}, args{"B", 2}, jsonDict{"A": 1, "B": 2}},
		{"With list element", jsonDict{"A": jsonList{1, 2}}, args{"A", 3}, jsonDict{"A": jsonList{1, 2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Add(tt.args.key, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		key interface{}
		v   interface{}
	}
	tests := []struct {
		name string
		d    jsonDict
		args args
		want jsonIDict
	}{
		{"Empty", nil, args{"A", 1}, jsonDict{"A": 1}},
		{"With element", jsonDict{"A": 1}, args{"A", 2}, jsonDict{"A": 2}},
		{"With element, another value", jsonDict{"A": 1}, args{"B", 2}, jsonDict{"A": 1, "B": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Set(tt.args.key, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_Transpose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    jsonDict
		want jsonIDict
	}{
		{"Empty", nil, jsonDict{}},
		{"Base", jsonDict{"A": 1}, jsonDict{"1": "A"}},
		{"Multiple", jsonDict{"A": 1, "B": 2, "C": 1}, jsonDict{"1": jsonList{"A", "C"}, "2": "B"}},
		{"List", jsonDict{"A": []int{1, 2, 3}, "B": 2, "C": 3}, jsonDict{"1": "A", "2": jsonList{"A", "B"}, "3": jsonList{"A", "C"}}},
		{"Complex", jsonDict{"A": jsonDict{"1": 1, "2": 2}, "B": 2, "C": 3}, jsonDict{"2": "B", "3": "C", fmt.Sprint(jsonDict{"1": 1, "2": 2}): "A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Transpose(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDict.Transpose() = %v, want %v", got, tt.want)
			}
		})
	}
}
