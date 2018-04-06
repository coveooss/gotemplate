package yaml

import (
	"testing"

	"github.com/coveo/gotemplate/types"
)

func Test_list_String(t *testing.T) {
	tests := []struct {
		name string
		l    yamlList
		want string
	}{
		{"Nil", nil, "[]\n"},
		{"Empty List", yamlList{}, "[]\n"},
		{"List of int", yamlList{1, 2, 3}, types.UnIndent(`
			- 1
			- 2
			- 3
			`)[1:]},
		{"List of string", strFixture, types.UnIndent(`
			- Hello
			- World,
			- I'm
			- Foo
			- Bar!
			`)[1:]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("yamlList.String():\ngot:\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}

func Test_dict_String(t *testing.T) {
	tests := []struct {
		name string
		d    yamlDict
		want string
	}{
		{"nil", nil, "{}\n"},
		{"Empty List", yamlDict{}, "{}\n"},
		{"Map", dictFixture, types.UnIndent(`
			float: 1.23
			int: 123
			list:
			- 1
			- two
			listInt:
			- 1
			- 2
			- 3
			map:
			  sub1: 1
			  sub2: two
			mapInt:
			  "1": 1
			  "2": two
			string: Foo bar
			`)[1:]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("yamlDict.String():\ngot:\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
}
