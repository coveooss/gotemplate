package main

import (
	"os"
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
			if got := formatList(tt.args.format, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_execFunc(t *testing.T) {
	cwd, _ := os.Getwd()
	type args struct {
		command string
		args    []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Current folder", args{"pwd", nil}, cwd, false},
		{"Echo", args{"echo", []interface{}{"Hello", "World!"}}, "Hello World!", false},
		{"Echo One arg", args{command: "echo Hello World!"}, "Hello World!", false},
		{"Echo separated arg", args{"echo Hello", []interface{}{"World!"}}, "Hello World!", false},
		{"Non working", args{command: "non existent"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := execFunc(tt.args.command, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("execFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("execFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeLists(t *testing.T) {
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
			if got := mergeLists(tt.args.lists...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeLists() = %v, want %v", got, tt.want)
			}
		})
	}
}
