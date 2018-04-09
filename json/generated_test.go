// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package json

import (
	"reflect"
	"strings"
	"testing"
)

func al(l jsonList) *jsonList { return &l }

var strFixture = jsonList(*jsonListHelper.NewStringList(strings.Split("Hello World, I'm Foo Bar!", " ")...).AsArray())

func Test_list_Append(t *testing.T) {
	tests := []struct {
		name   string
		l      jsonIList
		values []interface{}
		want   jsonIList
	}{
		{"Empty", al(nil), []interface{}{1, 2, 3}, al(jsonList{1, 2, 3})},
		{"List of int", al(jsonList{1, 2, 3}), []interface{}{4}, al(jsonList{1, 2, 3, 4})},
		{"List of string", &strFixture, []interface{}{"That's all folks!"}, al(jsonList{"Hello", "World,", "I'm", "Foo", "Bar!", "That's all folks!"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			if got := l.Append(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonList.Append():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if !reflect.DeepEqual(l, tt.want) {
				t.Errorf("jsonList.Append():\nsrc %[1]v (%[1]T)\nwant %[2]v (%[2]T)", l, tt.want)
			}
		})
	}
}

func Test_list_AsList(t *testing.T) {
	tests := []struct {
		name string
		l    jsonList
		want []interface{}
	}{
		{"Nil", nil, []interface{}{}},
		{"Empty List", jsonList{}, []interface{}{}},
		{"List of int", jsonList{1, 2, 3}, []interface{}{1, 2, 3}},
		{"List of string", strFixture, []interface{}{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			got := l.AsArray()
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("jsonList.AsList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}

			// We add an element to the jsonList and we check that bot jsonList are modified
			l.Append("Modified", 1, 2, 3)
			pList := l.AsArray()
			if !reflect.DeepEqual(pList, got) {
				t.Errorf("After modification::\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, pList)
			}
		})
	}
}

func Test_list_Capacity(t *testing.T) {
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
				t.Errorf("jsonList.Capacity() = %v, want %v", got, tt.want)
			}
			if tt.l.Capacity() != tt.l.Cap() {
				t.Errorf("Cap and Capacity return different values")
			}
		})
	}
}

func Test_list_Clone(t *testing.T) {
	tests := []struct {
		name string
		l    jsonList
		want jsonIList
	}{
		{"Empty List", jsonList{}, al(jsonList{})},
		{"List of int", jsonList{1, 2, 3}, al(jsonList{1, 2, 3})},
		{"List of string", strFixture, al(jsonList{"Hello", "World,", "I'm", "Foo", "Bar!"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonList.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Get(t *testing.T) {
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
				t.Errorf("jsonList.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_list_Len(t *testing.T) {
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
				t.Errorf("jsonList.Len() = %v, want %v", got, tt.want)
			}
			if tt.l.Len() != tt.l.Count() {
				t.Errorf("Len and Count return different values")
			}
		})
	}
}

func Test_NewList(t *testing.T) {
	type args struct {
		size     int
		capacity int
	}
	tests := []struct {
		name string
		args args
		want jsonIList
	}{
		{"Empty", args{0, 0}, &jsonList{}},
		{"With nil elements", args{10, 0}, al(make(jsonList, 10))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonListHelper.CreateList(tt.args.size, tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonList.CreateList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Reverse(t *testing.T) {
	tests := []struct {
		name string
		l    jsonList
		want jsonIList
	}{
		{"Empty List", jsonList{}, al(jsonList{})},
		{"List of int", jsonList{1, 2, 3}, al(jsonList{3, 2, 1})},
		{"List of string", strFixture, al(jsonList{"Bar!", "Foo", "I'm", "World,", "Hello"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			if got := l.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonList.Reverse():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Set(t *testing.T) {
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
		{"Empty", al(nil), args{2, 1}, al(jsonList{nil, nil, 1}), false},
		{"List of int", al(jsonList{1, 2, 3}), args{0, 10}, al(jsonList{10, 2, 3}), false},
		{"List of string", al(strFixture), args{2, "You're"}, al(jsonList{"Hello", "World,", "You're", "Foo", "Bar!"}), false},
		{"Negative", al(nil), args{-1, "negative value"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			got, err := l.Set(tt.args.i, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonList.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonList.Set():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if err == nil && !reflect.DeepEqual(l, tt.want) {
				t.Errorf("jsonList.Set():\nsrc %[1]v (%[1]T)\nwant %[2]v (%[2]T)", l, tt.want)
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
	for key := range d1.AsMap() {
		v1, v2 := d1.Get(key), d2.Get(key)
		if reflect.DeepEqual(v1, v2) {
			continue
		}
		t.Logf("'%[1]v' = %[2]v (%[2]T) vs %[3]v (%[3]T)", key, v1, v2)
	}
}

func Test_dict_AsMap(t *testing.T) {
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
				t.Errorf("jsonDict.AsMap():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Clone(t *testing.T) {
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
				t.Errorf("jsonDict.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
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
	tests := []struct {
		name         string
		d            jsonDict
		args         []int
		want         jsonIList
		wantLen      int
		wantCapacity int
	}{
		{"Nil", nil, nil, al(jsonList{}), 0, 0},
		{"Empty", jsonDict{}, nil, al(jsonList{}), 0, 0},
		{"Map", dictFixture, nil, al(jsonList{}), 0, 0},
		{"Map with size", dictFixture, []int{3}, al(jsonList{nil, nil, nil}), 3, 3},
		{"Map with capacity", dictFixture, []int{0, 10}, al(jsonList{}), 0, 10},
		{"Map with size&capacity", dictFixture, []int{3, 10}, al(jsonList{nil, nil, nil}), 3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.CreateList(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonDict.CreateList() = %v, want %v", got, tt.want)
			}
			if got.Len() != tt.wantLen || got.Cap() != tt.wantCapacity {
				t.Errorf("jsonDict.CreateList() size: %d, %d vs %d, %d", got.Len(), got.Cap(), tt.wantLen, tt.wantCapacity)
			}
		})
	}
}

func Test_dict_Default(t *testing.T) {
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
				t.Errorf("jsonDict.Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_Delete(t *testing.T) {
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
				t.Errorf("jsonDict.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonDict.Delete():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}

func Test_dict_Flush(t *testing.T) {
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
				t.Errorf("jsonDict.Flush():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
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
	tests := []struct {
		name string
		d    jsonDict
		want jsonIList
	}{
		{"Empty", nil, al(jsonList{})},
		{"Map", dictFixture, al(jsonList{"float", "int", "list", "listInt", "map", "mapInt", "string"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonDict.Keys():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_KeysAsString(t *testing.T) {
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
				t.Errorf("jsonDict.KeysAsString():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Merge(t *testing.T) {
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
				t.Errorf("jsonDict.Merge():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}