package implementation

import (
	"reflect"
	"strings"
	"testing"
)

func al(l baseList) *baseList { return &l }

var strFixture = baseList(*baseListHelper.NewStringList(strings.Split("Hello World, I'm Foo Bar!", " ")...).AsArray())

func Test_list_Append(t *testing.T) {
	tests := []struct {
		name   string
		l      baseIList
		values []interface{}
		want   baseIList
	}{
		{"Empty", al(nil), []interface{}{1, 2, 3}, al(baseList{1, 2, 3})},
		{"List of int", al(baseList{1, 2, 3}), []interface{}{4}, al(baseList{1, 2, 3, 4})},
		{"List of string", &strFixture, []interface{}{"That's all folks!"}, al(baseList{"Hello", "World,", "I'm", "Foo", "Bar!", "That's all folks!"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			if got := l.Append(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.Append():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if !reflect.DeepEqual(l, tt.want) {
				t.Errorf("baseList.Append():\nsrc  %[1]v (%[1]T)\nwant %[2]v (%[2]T)", l, tt.want)
			}
		})
	}
}

func Test_list_AsList(t *testing.T) {
	tests := []struct {
		name string
		l    baseList
		want []interface{}
	}{
		{"Nil", nil, []interface{}{}},
		{"Empty List", baseList{}, []interface{}{}},
		{"List of int", baseList{1, 2, 3}, []interface{}{1, 2, 3}},
		{"List of string", strFixture, []interface{}{"Hello", "World,", "I'm", "Foo", "Bar!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			got := l.AsArray()
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("baseList.AsList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}

			// We add an element to the baseList and we check that bot baseList are modified
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
		l    baseIList
		want int
	}{
		{"Empty List with 100 spaces", baseListHelper.CreateList(0, 100), 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Capacity(); got != tt.want {
				t.Errorf("baseList.Capacity() = %v, want %v", got, tt.want)
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
		l    baseList
		want baseIList
	}{
		{"Empty List", baseList{}, al(baseList{})},
		{"List of int", baseList{1, 2, 3}, al(baseList{1, 2, 3})},
		{"List of string", strFixture, al(baseList{"Hello", "World,", "I'm", "Foo", "Bar!"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Get(t *testing.T) {
	tests := []struct {
		name  string
		l     baseList
		index int
		want  interface{}
	}{
		{"Empty List", baseList{}, 0, nil},
		{"Negative index", baseList{}, -1, nil},
		{"List of int", baseList{1, 2, 3}, 0, 1},
		{"List of string", strFixture, 1, "World,"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Get(tt.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_list_Len(t *testing.T) {
	tests := []struct {
		name string
		l    baseList
		want int
	}{
		{"Empty List", baseList{}, 0},
		{"List of int", baseList{1, 2, 3}, 3},
		{"List of string", strFixture, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Len(); got != tt.want {
				t.Errorf("baseList.Len() = %v, want %v", got, tt.want)
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
		want baseIList
	}{
		{"Empty", args{0, 0}, &baseList{}},
		{"With nil elements", args{10, 0}, al(make(baseList, 10))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := baseListHelper.CreateList(tt.args.size, tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.CreateList():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_list_Reverse(t *testing.T) {
	tests := []struct {
		name string
		l    baseList
		want baseIList
	}{
		{"Empty List", baseList{}, al(baseList{})},
		{"List of int", baseList{1, 2, 3}, al(baseList{3, 2, 1})},
		{"List of string", strFixture, al(baseList{"Bar!", "Foo", "I'm", "World,", "Hello"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			if got := l.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.Reverse():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
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
		l       baseIList
		args    args
		want    baseIList
		wantErr bool
	}{
		{"Empty", al(nil), args{2, 1}, al(baseList{nil, nil, 1}), false},
		{"List of int", al(baseList{1, 2, 3}), args{0, 10}, al(baseList{10, 2, 3}), false},
		{"List of string", al(strFixture), args{2, "You're"}, al(baseList{"Hello", "World,", "You're", "Foo", "Bar!"}), false},
		{"Negative", al(nil), args{-1, "negative value"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.l.Clone()
			got, err := l.Set(tt.args.i, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseList.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseList.Set():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
			if err == nil && !reflect.DeepEqual(l, tt.want) {
				t.Errorf("baseList.Set():\nsrc  %[1]v (%[1]T)\nwant %[2]v (%[2]T)", l, tt.want)
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

var dictFixture = baseDict(baseDictHelper.AsDictionary(mapFixture).AsMap())

func dumpKeys(t *testing.T, d1, d2 baseIDict) {
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
		d    baseDict
		want map[string]interface{}
	}{
		{"Nil", nil, nil},
		{"Empty", baseDict{}, map[string]interface{}{}},
		{"Map", dictFixture, map[string]interface{}(dictFixture)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.AsMap():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Clone(t *testing.T) {
	tests := []struct {
		name string
		d    baseDict
		keys []interface{}
		want baseIDict
	}{
		{"Nil", nil, nil, baseDict{}},
		{"Empty", baseDict{}, nil, baseDict{}},
		{"Map", dictFixture, nil, dictFixture},
		{"Map with Fields", dictFixture, []interface{}{"int", "list"}, baseDict(dictFixture).Omit("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Clone(tt.keys...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.Clone():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
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

func Test_baseDict_CreateList(t *testing.T) {
	tests := []struct {
		name         string
		d            baseDict
		args         []int
		want         baseIList
		wantLen      int
		wantCapacity int
	}{
		{"Nil", nil, nil, al(baseList{}), 0, 0},
		{"Empty", baseDict{}, nil, al(baseList{}), 0, 0},
		{"Map", dictFixture, nil, al(baseList{}), 0, 0},
		{"Map with size", dictFixture, []int{3}, al(baseList{nil, nil, nil}), 3, 3},
		{"Map with capacity", dictFixture, []int{0, 10}, al(baseList{}), 0, 10},
		{"Map with size&capacity", dictFixture, []int{3, 10}, al(baseList{nil, nil, nil}), 3, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.CreateList(tt.args...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.CreateList() = %v, want %v", got, tt.want)
			}
			if got.Len() != tt.wantLen || got.Cap() != tt.wantCapacity {
				t.Errorf("baseDict.CreateList() size: %d, %d vs %d, %d", got.Len(), got.Cap(), tt.wantLen, tt.wantCapacity)
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
		d    baseDict
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
				t.Errorf("baseDict.Default() = %v, want %v", got, tt.want)
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
		d       baseDict
		args    args
		want    baseIDict
		wantErr bool
	}{
		{"Empty", nil, args{}, baseDict{}, true},
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
				t.Errorf("baseDict.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.Delete():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}

func Test_dict_Flush(t *testing.T) {
	tests := []struct {
		name string
		d    baseDict
		keys []interface{}
		want baseIDict
	}{
		{"Empty", nil, nil, baseDict{}},
		{"Map", dictFixture, nil, baseDict{}},
		{"Non existant key", dictFixture, []interface{}{"Test"}, dictFixture},
		{"Map with keys", dictFixture, []interface{}{"int", "list"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
		{"Map with keys + non existant", dictFixture, []interface{}{"int", "list", "Test"}, dictFixture.Clone("float", "string", "listInt", "map", "mapInt")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Flush(tt.keys...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.Flush():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
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
		d    baseDict
		want baseIList
	}{
		{"Empty", nil, al(baseList{})},
		{"Map", dictFixture, al(baseList{"float", "int", "list", "listInt", "map", "mapInt", "string"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.Keys():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_KeysAsString(t *testing.T) {
	tests := []struct {
		name string
		d    baseDict
		want []string
	}{
		{"Empty", nil, []string{}},
		{"Map", dictFixture, []string{"float", "int", "list", "listInt", "map", "mapInt", "string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.KeysAsString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.KeysAsString():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_dict_Merge(t *testing.T) {
	adding1 := baseDict{
		"int":        1000,
		"Add1Int":    1,
		"Add1String": "string",
	}
	adding2 := baseDict{
		"Add2Int":    1,
		"Add2String": "string",
		"map": map[string]interface{}{
			"sub1":   2,
			"newVal": "NewValue",
		},
	}
	type args struct {
		baseDict baseIDict
		dicts    []baseIDict
	}
	tests := []struct {
		name string
		d    baseDict
		args args
		want baseIDict
	}{
		{"Empty", nil, args{nil, []baseIDict{}}, baseDict{}},
		{"Add map to empty", nil, args{dictFixture, []baseIDict{}}, dictFixture},
		{"Add map to same map", dictFixture, args{dictFixture, []baseIDict{}}, dictFixture},
		{"Add empty to map", dictFixture, args{nil, []baseIDict{}}, dictFixture},
		{"Add new1 to map", dictFixture, args{adding1, []baseIDict{}}, dictFixture.Clone().Merge(adding1)},
		{"Add new2 to map", dictFixture, args{adding2, []baseIDict{}}, dictFixture.Clone().Merge(adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []baseIDict{adding2}}, dictFixture.Clone().Merge(adding1, adding2)},
		{"Add new1 & new2 to map", dictFixture, args{adding1, []baseIDict{adding2}}, dictFixture.Clone().Merge(adding1).Merge(adding2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.d.Clone()
			got := d.Merge(tt.args.baseDict, tt.args.dicts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseDict.Merge():\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
				dumpKeys(t, got, tt.want)
			}
		})
	}
}
