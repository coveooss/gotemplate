package utils

import (
	"reflect"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/collections/implementation"
)

type iList = collections.IGenericList
type list = implementation.ListTypeName

func TestFormatList(t *testing.T) {
	type args struct {
		format string
		v      interface{}
	}
	tests := []struct {
		name string
		args args
		want iList
	}{
		{"quote", args{`"%v"`, []int{1, 2}}, list{`"1"`, `"2"`}},
		{"greating", args{"Hello %v", []int{1, 2}}, list{"Hello 1", "Hello 2"}},
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
	tests := []struct {
		name string
		args []iList
		want iList
	}{
		{"Empty list", nil, nil},
		{"Simple list", []iList{list{1, 2, 3}}, list{1, 2, 3}},
		{"Two lists", []iList{list{1, 2, 3}, list{4, 5, 6}}, list{1, 2, 3, 4, 5, 6}},
		{"Three lists mixed", []iList{list{"One", 2, "3"}, list{4, 5, 6}, list{"7", "8", "9"}}, list{"One", 2, "3", 4, 5, 6, "7", "8", "9"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeLists(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeLists() = %v, want %v", got, tt.want)
			}
		})
	}
}
