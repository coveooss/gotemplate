// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package xml

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var strFixture = xmlList(xmlListHelper.NewStringList(strings.Split("Hello World, I'm Foo Bar!", " ")...).AsArray())

func Test_list_Append(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      xmlIList
		values []interface{}
		want   xmlIList
	}{
		{"Empty", xmlList{}, []interface{}{1, 2, 3}, xmlList{1, 2, 3}},
		{"List of int", xmlList{1, 2, 3}, []interface{}{4, 5}, xmlList{1, 2, 3, 4, 5}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, xmlList{"Hello", "World,", "I'm", "Foo", "Bar!", "That's all folks!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Append(tt.values...))
		})
	}
}

func Test_list_Prepend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		l      xmlIList
		values []interface{}
		want   xmlIList
	}{
		{"Empty", xmlList{}, []interface{}{1, 2, 3}, xmlList{1, 2, 3}},
		{"List of int", xmlList{1, 2, 3}, []interface{}{4, 5}, xmlList{4, 5, 1, 2, 3}},
		{"List of string", strFixture, []interface{}{"That's all folks!"}, xmlList{"That's all folks!", "Hello", "World,", "I'm", "Foo", "Bar!"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Prepend(tt.values...))
		})
	}
}

func Test_list_AsArray(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want []interface{}
	}{
		{"Empty List", xmlList{}, []interface{}{}},
		{"List of int", xmlList{1, 2, 3}, []interface{}{1, 2, 3}},
		{"List of string", strFixture, []interface{}{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.AsArray())
		})
	}
}

func Test_XmlList_Strings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want []string
	}{
		{"Empty List", xmlList{}, []string{}},
		{"List of int", xmlList{1, 2, 3}, []string{"1", "2", "3"}},
		{"List of string", strFixture, []string{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Strings())
		})
	}
}

func Test_list_Capacity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlIList
		want int
	}{
		{"Empty List with 100 spaces", xmlListHelper.CreateList(0, 100), 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Capacity())
			assert.Equal(t, tt.want, tt.l.Cap())
		})
	}
}

func Test_list_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want xmlIList
	}{
		{"Empty List", xmlList{}, xmlList{}},
		{"List of int", xmlList{1, 2, 3}, xmlList{1, 2, 3}},
		{"List of string", strFixture, xmlList{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Clone())
		})
	}
}

func Test_list_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       xmlList
		indexes []int
		want    interface{}
	}{
		{"Empty List", xmlList{}, []int{0}, nil},
		{"Negative index", xmlList{}, []int{-1}, nil},
		{"List of int", xmlList{1, 2, 3}, []int{0}, 1},
		{"List of string", strFixture, []int{1}, "World,"},
		{"Get last", strFixture, []int{-1}, "Bar!"},
		{"Get before last", strFixture, []int{-2}, "Foo"},
		{"A way to before last", strFixture, []int{-12}, nil},
		{"Get nothing", strFixture, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Get(tt.indexes...))
		})
	}
}

func Test_list_GetTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		kind bool
		l    xmlList
		want interface{}
	}{
		{"Empty", false, nil, xmlList{}},
		{"Fixture", false, strFixture, xmlList{"string", "string", "string", "string", "string"}},
		{"Mixed Types", false, xmlList{1, 1.2, true, "Hello", xmlList{}, xmlDict{}}, xmlList{"int", "float64", "bool", "string", xmlLower + "List", xmlLower + "Dict"}},
		{"Mixed Kinds", true, xmlList{1, 1.2, true, "Hello", xmlList{}, xmlDict{}}, xmlList{"int", "float64", "bool", "string", "slice", "map"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFunc := tt.l.GetTypes
			if tt.kind {
				testFunc = tt.l.GetKinds
			}
			assert.Equal(t, tt.want, testFunc())
		})
	}
}

func Test_list_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want int
	}{
		{"Empty List", xmlList{}, 0},
		{"List of int", xmlList{1, 2, 3}, 3},
		{"List of string", strFixture, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Count())
			assert.Equal(t, tt.want, tt.l.Len())
		})
	}
}

func Test_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []int
		want    xmlIList
		wantErr string
	}{
		{"Empty", nil, xmlList{}, ""},
		{"With nil elements", []int{10}, make(xmlList, 10), ""},
		{"With capacity", []int{0, 10}, make(xmlList, 0, 10), ""},
		{"Too many args", []int{0, 10, 1}, nil, "CreateList only accept 2 arguments, size and capacity"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.EqualError(t, err.(error), tt.wantErr)
				}
			}()

			got := xmlListHelper.CreateList(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want.Cap(), got.Capacity())
		})
	}
}

func Test_list_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		args []int
		want xmlIList
	}{
		{"Empty", nil, nil, xmlList{}},
		{"Existing List", xmlList{1, 2}, nil, xmlList{}},
		{"With Empty spaces", xmlList{1, 2}, []int{5}, xmlList{nil, nil, nil, nil, nil}},
		{"With Capacity", xmlList{1, 2}, []int{0, 5}, xmlListHelper.CreateList(0, 5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.l.Create(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want.Capacity(), got.Cap())
		})
	}
}

func Test_list_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		args []interface{}
		want xmlIList
	}{
		{"Empty", nil, nil, xmlList{}},
		{"Existing List", xmlList{1, 2}, nil, xmlList{}},
		{"With elements", xmlList{1, 2}, []interface{}{3, 4, 5}, xmlList{3, 4, 5}},
		{"With strings", xmlList{1, 2}, []interface{}{"Hello", "World"}, xmlList{"Hello", "World"}},
		{"With nothing", xmlList{1, 2}, []interface{}{}, xmlList{}},
		{"With nil", xmlList{1, 2}, nil, xmlList{}},
		{"Adding array", xmlList{1, 2}, []interface{}{xmlList{3, 4}}, xmlList{3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.New(tt.args...))
		})
	}
}

func Test_list_CreateDict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       xmlList
		args    []int
		want    xmlIDict
		wantErr string
	}{
		{"Empty", nil, nil, xmlDict{}, ""},
		{"With capacity", nil, []int{10}, xmlDict{}, ""},
		{"With too many parameters", nil, []int{10, 1}, nil, "CreateList only accept 1 argument for size"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.EqualError(t, err.(error), tt.wantErr)
				}
			}()
			assert.Equal(t, tt.want, tt.l.CreateDict(tt.args...))
		})
	}
}

func Test_list_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		l                xmlList
		args             []interface{}
		want, wantStrict bool
	}{
		{"Empty List", nil, []interface{}{}, false, false},
		{"Search nothing", xmlList{1}, nil, true, true},
		{"Search nothing 2", xmlList{1}, []interface{}{}, true, true},
		{"Not there", xmlList{1}, []interface{}{2}, false, false},
		{"Included", xmlList{1, 2}, []interface{}{2}, true, true},
		{"Partially there", xmlList{1, 2}, []interface{}{2, 3}, false, false},
		{"Different types", xmlList{1, 2, "3"}, []interface{}{"2", 3}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Contains(tt.args...))
			assert.Equal(t, tt.wantStrict, tt.l.ContainsStrict(tt.args...))
			assert.Equal(t, tt.want, tt.l.Has(tt.args...))
		})
	}
}

func Test_list_Find(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		l                xmlList
		element          interface{}
		want, wantStrict xmlList
	}{
		{"Empty List", nil, 2, xmlList{}, xmlList{}},
		{"Not found", xmlList{0, 1, 2, 3}, 4, xmlList{}, xmlList{}},
		{"Fist", xmlList{0, 1, 2, 3}, 0, xmlList{0}, xmlList{0}},
		{"Last", xmlList{0, 1, 2, 3}, 3, xmlList{3}, xmlList{3}},
		{"Many", xmlList{0, 1, 2, 3, 0, 1, 2, 3}, 3, xmlList{3, 7}, xmlList{3, 7}},
		{"Different type", xmlList{0, 1, 2, 3, "2", 2.0}, 2.0, xmlList{2, 4, 5}, xmlList{5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Find(tt.element))
			assert.Equal(t, tt.wantStrict, tt.l.FindStrict(tt.element))
		})
	}
}

func Test_list_First_Last(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		l         xmlList
		wantFirst interface{}
		wantLast  interface{}
	}{
		{"Nil", nil, nil, nil},
		{"Empty", xmlList{}, nil, nil},
		{"One element", xmlList{1}, 1, 1},
		{"Many element ", xmlList{1, "two", 3.1415, "four"}, 1, "four"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFirst, tt.l.First())
			assert.Equal(t, tt.wantLast, tt.l.Last())
		})
	}
}

func Test_list_Pop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		l        xmlList
		args     []int
		want     interface{}
		wantList xmlList
	}{
		{"Nil", nil, nil, nil, xmlList{}},
		{"Empty", xmlList{}, nil, nil, xmlList{}},
		{"Non existent", xmlList{}, []int{1}, nil, xmlList{}},
		{"Empty with args", xmlList{}, []int{1, 3}, xmlList{nil, nil}, xmlList{}},
		{"List with bad index", xmlList{0, 1, 2, 3, 4, 5}, []int{1, 3, 8}, xmlList{1, 3, nil}, xmlList{0, 2, 4, 5}},
		{"Pop last element", xmlList{0, 1, 2, 3, 4, 5}, nil, 5, xmlList{0, 1, 2, 3, 4}},
		{"Pop before last", xmlList{0, 1, 2, 3, 4, 5}, []int{-2}, 4, xmlList{0, 1, 2, 3, 5}},
		{"Pop first element", xmlList{0, 1, 2, 3, 4, 5}, []int{0}, 0, xmlList{1, 2, 3, 4, 5}},
		{"Pop all", xmlList{0, 1, 2, 3}, []int{0, 1, 2, 3}, xmlList{0, 1, 2, 3}, xmlList{}},
		{"Pop same element many time", xmlList{0, 1, 2, 3}, []int{1, 1, 2, 2}, xmlList{1, 1, 2, 2}, xmlList{0, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, listAfter := tt.l.Pop(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantList, listAfter)
		})
	}
}

func Test_list_Intersect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		args []interface{}
		want xmlList
	}{
		{"Empty List", nil, []interface{}{}, xmlList{}},
		{"Intersect nothing", xmlList{1}, nil, xmlList{}},
		{"Intersect nothing 2", xmlList{1}, []interface{}{}, xmlList{}},
		{"Not there", xmlList{1}, []interface{}{2}, xmlList{}},
		{"Included", xmlList{1, 2}, []interface{}{2}, xmlList{2}},
		{"Partially there", xmlList{1, 2}, []interface{}{2, 3}, xmlList{2}},
		{"With duplicates", xmlList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{3, 4, 5, 6, 7, 8, 7, 6, 5, 5, 4, 3}, xmlList{3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Intersect(tt.args...))
		})
	}
}

func Test_list_Union(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		args []interface{}
		want xmlList
	}{
		{"Empty List", nil, []interface{}{}, xmlList{}},
		{"Intersect nothing", xmlList{1}, nil, xmlList{1}},
		{"Intersect nothing 2", xmlList{1}, []interface{}{}, xmlList{1}},
		{"Not there", xmlList{1}, []interface{}{2}, xmlList{1, 2}},
		{"Included", xmlList{1, 2}, []interface{}{2}, xmlList{1, 2}},
		{"Partially there", xmlList{1, 2}, []interface{}{2, 3}, xmlList{1, 2, 3}},
		{"With duplicates", xmlList{1, 2, 3, 4, 5, 4, 3, 2, 1}, []interface{}{8, 7, 6, 5, 6, 7, 8}, xmlList{1, 2, 3, 4, 5, 8, 7, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Union(tt.args...))
		})
	}
}

func Test_list_Without(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		args []interface{}
		want xmlList
	}{
		{"Empty List", nil, []interface{}{}, xmlList{}},
		{"Remove nothing", xmlList{1}, nil, xmlList{1}},
		{"Remove nothing 2", xmlList{1}, []interface{}{}, xmlList{1}},
		{"Not there", xmlList{1}, []interface{}{2}, xmlList{1}},
		{"Included", xmlList{1, 2}, []interface{}{2}, xmlList{1}},
		{"Partially there", xmlList{1, 2}, []interface{}{2, 3}, xmlList{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Without(tt.args...))
		})
	}
}

func Test_list_Unique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want xmlList
	}{
		{"Empty List", nil, xmlList{}},
		{"Remove nothing", xmlList{1}, xmlList{1}},
		{"Duplicates following", xmlList{1, 1, 2, 3}, xmlList{1, 2, 3}},
		{"Duplicates not following", xmlList{1, 2, 3, 1, 2, 3, 4}, xmlList{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Unique())
		})
	}
}

func Test_list_RemoveEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want xmlIList
	}{
		{"Empty List", xmlList{}, xmlList{}},
		{"List of int", xmlList{1, 2, 3}, xmlList{3, 2, 1}},
		{"List of string", strFixture, xmlList{"Bar!", "Foo", "I'm", "World,", "Hello"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Clone().Reverse())
		})
	}
}

func Test_list_Reverse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    xmlList
		want xmlIList
	}{
		{"Empty List", xmlList{}, xmlList{}},
		{"List of int", xmlList{1, 2, 3}, xmlList{3, 2, 1}},
		{"List of string", strFixture, xmlList{"Bar!", "Foo", "I'm", "World,", "Hello"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.Clone().Reverse())
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
		l       xmlIList
		args    args
		want    xmlIList
		wantErr string
	}{
		{"Empty", xmlList{}, args{2, 1}, xmlList{nil, nil, 1}, ""},
		{"List of int", xmlList{1, 2, 3}, args{0, 10}, xmlList{10, 2, 3}, ""},
		{"List of string", strFixture, args{2, "You're"}, xmlList{"Hello", "World,", "You're", "Foo", "Bar!"}, ""},
		{"Negative", xmlList{}, args{-1, "negative value"}, nil, "index must be positive number"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Clone().Set(tt.args.i, tt.args.v)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
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

var dictFixture = xmlDict(xmlDictHelper.AsDictionary(mapFixture).AsMap())

func Test_dict_AsMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		want map[string]interface{}
	}{
		{"Nil", nil, nil},
		{"Empty", xmlDict{}, map[string]interface{}{}},
		{"Map", dictFixture, map[string]interface{}(dictFixture)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.AsMap())
		})
	}
}

func Test_dict_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		keys []interface{}
		want xmlIDict
	}{
		{"Nil", nil, nil, xmlDict{}},
		{"Empty", xmlDict{}, nil, xmlDict{}},
		{"Map", dictFixture, nil, dictFixture},
		{"Map with Fields", dictFixture, []interface{}{"int", "list"}, xmlDict(dictFixture).Omit("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Clone(tt.keys...)
			assert.Equal(t, tt.want, got)

			// Ensure that the copy is distinct from the original
			got.Set("NewFields", "Test")
			assert.NotEqual(t, tt.want, got)
			assert.True(t, got.Has("NewFields"))
			assert.Equal(t, "Test", got.Get("NewFields"))
			assert.Equal(t, tt.want.Count()+1, got.Len())
		})
	}
}

func Test_XmlDict_CreateList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		d            xmlDict
		args         []int
		want         xmlIList
		wantLen      int
		wantCapacity int
	}{
		{"Nil", nil, nil, xmlList{}, 0, 0},
		{"Empty", xmlDict{}, nil, xmlList{}, 0, 0},
		{"Map", dictFixture, nil, xmlList{}, 0, 0},
		{"Map with size", dictFixture, []int{3}, xmlList{nil, nil, nil}, 3, 3},
		{"Map with capacity", dictFixture, []int{0, 10}, xmlList{}, 0, 10},
		{"Map with size&capacity", dictFixture, []int{3, 10}, xmlList{nil, nil, nil}, 3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.CreateList(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantLen, got.Len())
			assert.Equal(t, tt.wantCapacity, got.Capacity())
		})
	}
}

func Test_dict_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		d       xmlDict
		args    []int
		want    xmlIDict
		wantErr string
	}{
		{"Empty", nil, nil, xmlDict{}, ""},
		{"With capacity", nil, []int{10}, xmlDict{}, ""},
		{"With too much parameter", nil, []int{10, 1}, nil, "CreateList only accept 1 argument for size"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.EqualError(t, err.(error), tt.wantErr)
				}
			}()
			assert.Equal(t, tt.want, tt.d.Create(tt.args...))
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
		d    xmlDict
		args args
		want interface{}
	}{
		{"Empty", nil, args{"Foo", "Bar"}, "Bar"},
		{"Map int", dictFixture, args{"int", 1}, 123},
		{"Map float", dictFixture, args{"float", 1}, 1.23},
		{"Map Non existent", dictFixture, args{"Foo", "Bar"}, "Bar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Default(tt.args.key, tt.args.defVal))
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
		d       xmlDict
		args    args
		want    xmlIDict
		wantErr string
	}{
		{"Empty", nil, args{}, xmlDict{}, "key <nil> not found"},
		{"Map", dictFixture, args{}, dictFixture, "key <nil> not found"},
		{"Non existent key", dictFixture, args{"Test", nil}, dictFixture, "key Test not found"},
		{"Map with keys", dictFixture, args{"int", []interface{}{"list"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), ""},
		{"Map with keys + non existent", dictFixture, args{"int", []interface{}{"list", "Test"}}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt"), "key Test not found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Clone().Delete(tt.args.key, tt.args.keys...)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func Test_dict_Flush(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		keys []interface{}
		want xmlIDict
	}{
		{"Empty", nil, nil, xmlDict{}},
		{"Map", dictFixture, nil, xmlDict{}},
		{"Non existent key", dictFixture, []interface{}{"Test"}, dictFixture},
		{"Map with keys", dictFixture, []interface{}{"int", "list"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
		{"Map with keys + non existent", dictFixture, []interface{}{"int", "list", "Test"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Flush(tt.keys...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, d, got)
		})
	}
}

func Test_dict_Keys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		want xmlIList
	}{
		{"Empty", nil, xmlList{}},
		{"Map", dictFixture, xmlList{str("float"), str("int"), str("list"), str("listInt"), str("map"), str("mapInt"), str("string")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.GetKeys())
		})
	}
}

func Test_dict_KeysAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		want strArray
	}{
		{"Empty", nil, strArray{}},
		{"Map", dictFixture, strArray{"float", "int", "list", "listInt", "map", "mapInt", "string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.KeysAsString())
		})
	}
}

func Test_dict_Merge(t *testing.T) {
	t.Parallel()

	adding1 := xmlDict{
		"int":        1000,
		"Add1Int":    1,
		"Add1String": "string",
	}
	adding2 := xmlDict{
		"Add2Int":    1,
		"Add2String": "string",
		"map": map[string]interface{}{
			"sub1":   2,
			"newVal": "NewValue",
		},
	}
	type args struct {
		xmlDict xmlIDict
		dicts   []xmlIDict
	}
	tests := []struct {
		name string
		d    xmlDict
		args args
		want xmlIDict
	}{
		{"Empty", nil, args{nil, []xmlIDict{}}, xmlDict{}},
		{"Add map to empty", nil, args{dictFixture, []xmlIDict{}}, dictFixture},
		{"Add map to same map", dictFixture, args{dictFixture, []xmlIDict{}}, dictFixture},
		{"Add empty to map", dictFixture, args{nil, []xmlIDict{}}, dictFixture},
		{"Add new1 to map", dictFixture, args{adding1, []xmlIDict{}}, dictFixture.Clone().Merge(adding1)},
		{"Add new2 to map", dictFixture, args{adding2, []xmlIDict{}}, dictFixture.Clone().Merge(adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []xmlIDict{adding2}}, dictFixture.Clone().Merge(adding1, adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []xmlIDict{adding2}}, dictFixture.Clone().Merge(adding1).Merge(adding2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := tt.d.Clone()
			got := d.Merge(tt.args.xmlDict, tt.args.dicts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_dict_Values(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		want xmlIList
	}{
		{"Empty", nil, xmlList{}},
		{"Map", dictFixture, xmlList{1.23, 123, xmlList{1, "two"}, xmlList{1, 2, 3}, xmlDict{"sub1": 1, "sub2": "two"}, xmlDict{"1": 1, "2": "two"}, "Foo bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.GetValues())
		})
	}
}

func Test_dict_Pop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		d          xmlDict
		args       []interface{}
		want       interface{}
		wantObject xmlIDict
	}{
		{"Nil", dictFixture, nil, nil, dictFixture},
		{"Pop one element", dictFixture, []interface{}{"float"}, 1.23, dictFixture.Omit("float")},
		{"Pop missing element", dictFixture, []interface{}{"undefined"}, nil, dictFixture},
		{"Pop element twice", dictFixture, []interface{}{"int", "int", "string"}, xmlList{123, 123, "Foo bar"}, dictFixture.Omit("int", "string")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Pop(tt.args...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantObject, d)
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
		d    xmlDict
		args args
		want xmlIDict
	}{
		{"Empty", nil, args{"A", 1}, xmlDict{"A": 1}},
		{"With element", xmlDict{"A": 1}, args{"A", 2}, xmlDict{"A": xmlList{1, 2}}},
		{"With element, another value", xmlDict{"A": 1}, args{"B", 2}, xmlDict{"A": 1, "B": 2}},
		{"With list element", xmlDict{"A": xmlList{1, 2}}, args{"A", 3}, xmlDict{"A": xmlList{1, 2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Add(tt.args.key, tt.args.v))
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
		d    xmlDict
		args args
		want xmlIDict
	}{
		{"Empty", nil, args{"A", 1}, xmlDict{"A": 1}},
		{"With element", xmlDict{"A": 1}, args{"A", 2}, xmlDict{"A": 2}},
		{"With element, another value", xmlDict{"A": 1}, args{"B", 2}, xmlDict{"A": 1, "B": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Set(tt.args.key, tt.args.v))
		})
	}
}

func Test_dict_Transpose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    xmlDict
		want xmlIDict
	}{
		{"Empty", nil, xmlDict{}},
		{"Base", xmlDict{"A": 1}, xmlDict{"1": str("A")}},
		{"Multiple", xmlDict{"A": 1, "B": 2, "C": 1}, xmlDict{"1": xmlList{str("A"), str("C")}, "2": str("B")}},
		{"List", xmlDict{"A": []int{1, 2, 3}, "B": 2, "C": 3}, xmlDict{"1": str("A"), "2": xmlList{str("A"), str("B")}, "3": xmlList{str("A"), str("C")}}},
		{"Complex", xmlDict{"A": xmlDict{"1": 1, "2": 2}, "B": 2, "C": 3}, xmlDict{"2": str("B"), "3": str("C"), fmt.Sprint(xmlDict{"1": 1, "2": 2}): str("A")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.Transpose())
		})
	}
}

func Test_dict_GetTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		kind bool
		d    xmlDict
		want interface{}
	}{
		{"Empty", false, nil, xmlDict{}},
		{"Fixture Types", false, dictFixture, xmlDict{
			"float":   "float64",
			"int":     "int",
			"list":    xmlLower + "List",
			"listInt": xmlLower + "List",
			"map":     xmlLower + "Dict",
			"mapInt":  xmlLower + "Dict",
			"string":  "string",
		}},
		{"Fixture Kinds", true, dictFixture, xmlDict{
			"float":   "float64",
			"int":     "int",
			"list":    "slice",
			"listInt": "slice",
			"map":     "map",
			"mapInt":  "map",
			"string":  "string",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFunc := tt.d.GetTypes
			if tt.kind {
				testFunc = tt.d.GetKinds
			}
			assert.Equal(t, tt.want, testFunc())
		})
	}
}

func Test_Xml_Type(t *testing.T) {
	t.Run("list", func(t *testing.T) { assert.Equal(t, str(xmlLower+"List"), xmlList{}.Type()) })
	t.Run("dict", func(t *testing.T) { assert.Equal(t, str(xmlLower+"Dict"), xmlDict{}.Type()) })
}

func Test_Xml_TypeName(t *testing.T) {
	t.Run("list", func(t *testing.T) { assert.Equal(t, str(xmlLower), xmlList{}.TypeName()) })
	t.Run("dict", func(t *testing.T) { assert.Equal(t, str(xmlLower), xmlDict{}.TypeName()) })
}

func Test_Xml_GetHelper(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		gotD, gotL := xmlList{}.GetHelpers()
		assert.Equal(t, xmlDictHelper.CreateDictionary().TypeName(), gotD.CreateDictionary().TypeName())
		assert.Equal(t, xmlListHelper.CreateList().TypeName(), gotL.CreateList().TypeName())
	})
	t.Run("dict", func(t *testing.T) {
		gotD, gotL := xmlDict{}.GetHelpers()
		assert.Equal(t, xmlDictHelper.CreateDictionary().TypeName(), gotD.CreateDictionary().TypeName())
		assert.Equal(t, xmlListHelper.CreateList().TypeName(), gotL.CreateList().TypeName())
	})
}
