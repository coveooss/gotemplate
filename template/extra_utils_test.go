package template

import (
	"testing"

	"github.com/coveooss/gotemplate/v3/json"
	"github.com/stretchr/testify/assert"
)

func Test_rawPrint(t *testing.T) {
	type il = []interface{}
	tests := []struct {
		name string
		args il
		want interface{}
	}{
		{"no argument", il{}, ""},
		{"1 argument", il{`"hello"`}, `"hello"`},
		{"2 arguments", il{`"1"`, `"2"`}, json.List{`!Q!"1"!Q!`, `!Q!"2"!Q!`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, rawPrint(tt.args...))
		})
	}
}

func Test_rawList(t *testing.T) {
	type il = []interface{}
	tests := []struct {
		name string
		args il
		want interface{}
	}{
		{"no argument", il{}, json.List{}},
		{"1 argument", il{`"hello"`}, json.List{`!Q!"hello"!Q!`}},
		{"2 arguments", il{`"1"`, `"2"`}, json.List{`!Q!"1"!Q!`, `!Q!"2"!Q!`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, rawList(tt.args...))
		})
	}
}
