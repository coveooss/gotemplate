package utils

import (
	"reflect"
	"testing"
)

func TestMarshalHCLVars(t *testing.T) {
	type test struct {
		Name  string
		Value int
	}
	var testNilPtr *test

	type args struct {
		value  interface{}
		pretty bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Integer", args{2, false}, "2"},
		{"Boolean", args{true, false}, "true"},
		{"String", args{"Hello world", false}, `"Hello world"`},
		{"List of integer", args{[]int{0, 1, 2, 3}, false}, `[0,1,2,3]`},
		{"Map", args{map[string]interface{}{"a": 0, "b": 1}, false}, `a=0 b=1`},
		{"Map (pretty)", args{map[string]interface{}{"a": 0, "b": 1}, true}, "a = 0\nb = 1\n"},
		{"Struct (pretty)", args{test{"name", 1}, true}, "Name = \"name\"\nValue = 1\n"},
		{"Struct Ptr (pretty)", args{&test{"name", 1}, true}, "Name = \"name\"\nValue = 1\n"},
		{"Array of 1 struct (pretty)", args{[]test{{"name", 1}}, true}, "[{\n  Name = \"name\"\n  Value = 1\n}]"},
		{"Array of 2 structs (pretty)", args{[]test{{"val1", 1}, {"val2", 1}}, true}, "[\n  {\n    Name = \"val1\"\n    Value = 1\n  },\n  {\n    Name = \"val2\"\n    Value = 1\n  },\n]"},
		{"Null value", args{nil, true}, "null"},
		{"Null struct", args{testNilPtr, true}, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(marshalHCL(tt.args.value, tt.args.pretty, 0)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalHCLVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
