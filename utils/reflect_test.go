package utils

import (
	"reflect"
	"strings"
	"testing"
)

type empty struct{}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want bool
	}{
		{"False", false, true},
		{"True", true, false},
		{"Nil", nil, true},
		{"String", "Hello", false},
		{"Empty string", "", true},
		{"Zero", 0, true},
		{"Uint", uint(10), false},
		{"Floating point zero", 0.0, true},
		{"Floating point", 10.0, false},
		{"Floating point negative", -10.0, false},
		{"Interface to nil", interface{}(nil), true},
		{"Interface to zero", interface{}(0), true},
		{"Interface to empty string", interface{}(""), true},
		{"Interface to string", interface{}("Foo"), false},
		{"Integer", 10, false},
		{"Empty struct", empty{}, false},
		{"Empty list", []string{}, true},
		{"List", []string{"Hello"}, false},
		{"Empty map", map[string]int{}, true},
		{"Map", map[string]int{"Hello": 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyValue(reflect.ValueOf(tt.arg)); got != tt.want {
				t.Errorf("IsEmptyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExported(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"Private", "privateValue", false},
		{"Public", "PublicValue", true},
		{"Zero number", "0", false},
		{"Positive number", "1", false},
		{"Underscore", "_test", false},
		{"Id with space", "ID 1", true},
		{"AllCap", "ALL_CAP_ID", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExported(tt.id); got != tt.want {
				t.Errorf("IsExported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIfUndef(t *testing.T) {
	const def = "default"
	type empty struct{}
	tests := []struct {
		name string
		arg  interface{}
		want interface{}
	}{
		{"False", false, false},
		{"True", true, true},
		{"Nil", nil, def},
		{"String", "Hello", "Hello"},
		{"Empty string", "", ""},
		{"Zero", 0, 0},
		{"Integer", 10, 10},
		{"Empty struct", empty{}, empty{}},
		{"Empty list", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfUndef(def, tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IfUndef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIfUndefNoValue(t *testing.T) {
	def := "default"
	t.Run("No Value", func(t *testing.T) {
		if got := IfUndef(def); !reflect.DeepEqual(got, def) {
			t.Errorf("IfUndef() = %v, want %v", got, def)
		}
	})
}

func TestIfManyValues(t *testing.T) {
	const def = "default"
	list := []interface{}{1, 2, 3}
	t.Run("Many values", func(t *testing.T) {
		if got := IfUndef(def, list...); !reflect.DeepEqual(got, list) {
			t.Errorf("IfUndef() = %v, want %v", got, list)
		}
	})
}

func TestIIf(t *testing.T) {
	tests := []struct {
		name      string
		testValue interface{}
		want      interface{}
	}{
		{"False", false, 2},
		{"True", true, 1},
		{"Nil", nil, 2},
		{"String", "Hello", 1},
		{"Empty string", "", 2},
		{"Zero", 0, 2},
		{"Integer", 10, 1},
		{"Empty struct", empty{}, 1},
		{"Empty list", []string{}, 2},
		{"List", []string{"Hello"}, 1},
		{"Empty map", map[string]int{}, 2},
		{"Map", map[string]int{"Hello": 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IIf(tt.testValue, 1, 2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IIf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	a := map[string]interface{}{
		"a": 123,
		"b": 1.23,
		"c": "Foo",
	}
	b := map[string]interface{}{
		"a": 321,
		"b": 3.21,
		"d": "Bar",
	}
	type args struct {
		destination map[string]interface{}
		sources     []map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"Merge A <= B", args{a, []map[string]interface{}{b}}, map[string]interface{}{
			"a": 123,
			"b": 1.23,
			"c": "Foo",
			"d": "Bar",
		}, false},
		{"Merge B <= A", args{b, []map[string]interface{}{a}}, map[string]interface{}{
			"a": 321,
			"b": 3.21,
			"c": "Foo",
			"d": "Bar",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeMaps(tt.args.destination, tt.args.sources...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeMaps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNativeRepresentation(t *testing.T) {
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
		M       map[string]interface{}
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
			M: map[string]interface{}{
				"a": "a",
				"b": 2,
			},
			SS: SubStruct{64, "Foo"},
		}, map[string]interface{}{
			"I": "123",
			"F": "1.23",
			"S": `"123"`,
			"A": []interface{}{"1", `"2"`},
			"M": map[string]interface{}{
				"a": `"a"`,
				"b": "2",
			},
			"SS": map[string]interface{}{
				"U": "64",
				"I": `"Foo"`,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNativeRepresentation(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNativeRepresentation()\ngot : %v\nwant: %v", got, tt.want)
			}
		})
	}
}

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
		{"YAML", args{"a: 10", &out1}, map[string]interface{}{"a": 10}, false},
		{"HCL", args{`a = 10 b = "Foo"`, &out1}, map[string]interface{}{"a": 10, "b": "Foo"}, false},
		{"JSON", args{`{ "a": 10, "b": "Foo" }`, &out1}, map[string]interface{}{"a": 10, "b": "Foo"}, false},
		{"Flexible", args{`a = 10 b = Foo`, &out1}, map[string]interface{}{"a": 10, "b": "Foo"}, false},
		{"No change", args{"NoChange", &out1}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertData(tt.args.data, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("ConvertData() error = %v, wantErr %v", err, tt.wantErr)
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
		M       map[string]interface{}
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
			M: map[string]interface{}{
				"a": "a",
				"b": 2,
			},
			SS: SubStruct{64, "Foo"},
		}, strings.TrimSpace(UnIndent(`
		declare -a A
		A=(1 "2")
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
			if got := ToBash(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNativeRepresentation()\ngot : %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func Test_quote(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"Simple value", "Foo", "Foo"},
		{"Simple value", "Foo Bar", `"Foo Bar"`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quote(tt.arg); got != tt.want {
				t.Errorf("quote() = %v, want %v", got, tt.want)
			}
		})
	}
}
