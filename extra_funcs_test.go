package main

import (
	"reflect"
	"testing"
)

func Test_listFormat(t *testing.T) {
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
			if got := listFormat(tt.args.format, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
