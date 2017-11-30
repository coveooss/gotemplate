package utils

import (
	"reflect"
	"testing"
)

func TestFormatList(t *testing.T) {
	type args struct {
		format string
		v      interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"quote", args{"\"%v\"", []int{1, 2}}, []string{"\"1\"", "\"2\""}},
		{"greating", args{"Hello %v", []int{1, 2}}, []string{"Hello 1", "Hello 2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatList(tt.args.format, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeLists(t *testing.T) {
	type args struct {
		lists [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{"Empty list", args{nil}, nil},
		{"Simple list", args{[][]interface{}{{1, 2, 3}}}, []interface{}{1, 2, 3}},
		{"Two lists", args{[][]interface{}{{1, 2, 3}, {4, 5, 6}}}, []interface{}{1, 2, 3, 4, 5, 6}},
		{"Three lists mixed", args{[][]interface{}{{"One", 2, "3"}, {4, 5, 6}, {"7", "8", "9"}}}, []interface{}{"One", 2, "3", 4, 5, 6, "7", "8", "9"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeLists(tt.args.lists...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeLists() = %v, want %v", got, tt.want)
			}
		})
	}
}
