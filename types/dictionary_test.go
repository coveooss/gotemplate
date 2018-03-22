package types

import (
	"reflect"
	"strings"
	"testing"
)

var simpleDict = Dictionary{
	"a": 1,
	"b": "2",
	"c": []int{1, 2, 3},
}

var complexDict = Dictionary{
	"int":    123,
	"float":  1.23,
	"string": "Foo",
	"map": map[int]int{
		1: 1 * 1,
		2: 2 * 2,
		3: 3 * 3,
	},
	"listFloat":  []float64{1.1, 2.2, 3.3, 4.4},
	"listString": strings.Fields("Hello World!, life is beautiful!"),
	"dict":       simpleDict,
	"list":       GenericList{1, 1.2, "3", "quatre"},
}

type args = GenericList

func TestDictionary_Clone(t *testing.T) {
	tests := []struct {
		name string
		d    Dictionary
		keys args
		want IDictionary
	}{
		{"Simple Dictionary", simpleDict, nil, simpleDict},
		{"Simple Dictionary a", simpleDict, args{"a"}, Dictionary{"a": 1}},
		{"Complex Dictionary", complexDict, nil, complexDict},
		{"Complex Dictionary int", complexDict, args{"int"}, Dictionary{"int": 123}},
		{"Complex Dictionary int, dict", complexDict, args{"int", "dict"}, Dictionary{"int": 123, "dict": simpleDict.Clone()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Clone(tt.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dictionary.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_Omit(t *testing.T) {
	tests := []struct {
		name string
		d    Dictionary
		keys args
		want IDictionary
	}{
		{"Simple Dictionary", simpleDict, nil, simpleDict},
		{"Simple Dictionary a", simpleDict, args{"a"}, simpleDict.Clone("b", "c")},
		{"Complex Dictionary", complexDict, nil, complexDict},
		{"Complex Dictionary int", complexDict, args{"int"}, complexDict.Clone("float", "string", "map", "listFloat", "listString", "dict", "list")},
		{"Complex Dictionary int, dict", complexDict, args{"int", "dict"}, complexDict.Clone("float", "string", "map", "listFloat", "listString", "list")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Omit(tt.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dictionary.Omit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_Get(t *testing.T) {
	tests := []struct {
		name string
		d    Dictionary
		key  interface{}
		want interface{}
	}{
		{"Simple a", simpleDict, "a", 1},
		{"Simple b", simpleDict, "b", "2"},
		{"Simple c", simpleDict, "c", simpleDict["c"]},
		{"Comlect dict", complexDict, "dict", complexDict["dict"]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Get(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dictionary.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDictionary_Delete(t *testing.T) {
	tests := []struct {
		name string
		d    Dictionary
		key  interface{}
		want IDictionary
	}{
		{"Simple a", simpleDict, "a", simpleDict.Omit("a")},
		{"Simple b", simpleDict, "b", simpleDict.Omit("b")},
		{"Simple c", simpleDict, "c", simpleDict.Omit("c")},
		{"Comlect dict", complexDict, "dict", complexDict.Omit("dict")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := tt.d.Clone()
			source.Delete(tt.key)
			if !reflect.DeepEqual(source, tt.want) {
				t.Errorf("Dictionary.Delete() = %v, want %v", source, tt.want)
			}
			if source.Len() != tt.want.Len() {
				t.Errorf("Dictionary.Len() = %v, want %v", source.Len(), tt.want.Len())
			}
		})
	}
}

func TestAsDictionary(t *testing.T) {
	tests := []struct {
		name       string
		object     interface{}
		wantResult IDictionary
		wantErr    bool
	}{
		{"Simple a", map[string]interface{}{"a": 1}, simpleDict.Clone("a"), false},
		{"Nil", nil, Dictionary{}, false},
		{"Not working string", "From string", nil, true},
		{"Not working int", 1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := AsDictionary(tt.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AsDictionary() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestDictionary_AsMap(t *testing.T) {
	tests := []struct {
		name string
		d    Dictionary
		want map[string]interface{}
	}{
		{"Simple", simpleDict, simpleDict},
		{"Nil", nil, Dictionary{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.AsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dictionary.AsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	a := Dictionary{
		"a": 123,
		"b": 1.23,
		"c": "Foo",
	}
	b := Dictionary{
		"a": 321,
		"b": 3.21,
		"d": "Bar",
	}
	type args struct {
		destination Dictionary
		sources     []IDictionary
	}
	tests := []struct {
		name    string
		args    args
		want    Dictionary
		wantErr bool
	}{
		{"Merge A <= B", args{a, []IDictionary{b}}, Dictionary{
			"a": 123,
			"b": 1.23,
			"c": "Foo",
			"d": "Bar",
		}, false},
		{"Merge B <= A", args{b, []IDictionary{a}}, Dictionary{
			"a": 321,
			"b": 3.21,
			"c": "Foo",
			"d": "Bar",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.args.destination.Clone()
			_, err := d.Merge(tt.args.sources...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeMaps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("MergeMaps() = %v, want %v", d, tt.want)
			}
		})
	}
}
