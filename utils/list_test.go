package utils

import (
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/collections/implementation"
	"github.com/stretchr/testify/assert"
)

type iList = collections.IGenericList
type list = implementation.ListTypeName

func TestFormatList(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
		want   iList
	}{
		{"Empty List", `"%v"`, []interface{}{}, list{}},
		{"Single element", `"%v"`, []interface{}{1}, list{`"1"`}},
		{"quote", `"%v"`, []interface{}{1, 2}, list{`"1"`, `"2"`}},
		{"greating", "Hello %v", []interface{}{1, 2}, list{"Hello 1", "Hello 2"}},
		{"greating list", "Hello %v", []interface{}{[]int{1, 2}}, list{"Hello 1", "Hello 2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatList(tt.format, tt.args...))
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
			assert.Equal(t, tt.want, MergeLists(tt.args...))
		})
	}
}
