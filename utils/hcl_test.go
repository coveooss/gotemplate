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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(marshalHCL(tt.args.value, tt.args.pretty, 0)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalHCLVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
