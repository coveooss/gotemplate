package hcl

import (
	"fmt"
	"testing"

	"github.com/coveo/gotemplate/v3/collections"
	"github.com/stretchr/testify/assert"
)

func Test_list_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    hclList
		want string
	}{
		{"Nil", nil, "[]"},
		{"Empty list", hclList{}, "[]"},
		{"List of int", hclList{1, 2, 3}, "[1,2,3]"},
		{"List of string", strFixture, `["Hello","World,","I'm","Foo","Bar!"]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.l.String())
		})
	}
}

func Test_dict_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    hclDict
		want string
	}{
		{"nil", nil, ""},
		{"Empty dict", hclDict{}, ""},
		{"Map", dictFixture, `float=1.23 int=123 list=[1,"two"] listInt=[1,2,3] map{sub1=1 sub2="two"} mapInt{"1"=1 "2"="two"} string="Foo bar"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.d.String())
		})
	}
}

func TestMarshalHCLVars(t *testing.T) {
	t.Parallel()

	type test struct {
		Name   string `hcl:",omitempty"`
		Value  int    `hcl:",omitempty"`
		Public bool   `hcl:"public"`
	}
	const (
		noIndent = ""
		indent   = "  "
	)
	var testNilPtr *test

	type args struct {
		value  interface{}
		indent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Integer", args{2, noIndent}, "2"},
		{"Boolean", args{true, noIndent}, "true"},
		{"String", args{"Hello world", noIndent}, `"Hello world"`},
		{"String with newline", args{"Hello\nworld\n", noIndent}, `"Hello\nworld\n"`},
		{"String with newline (pretty)", args{"Hello\n\"world\"\n", indent}, "<<-EOF\nHello\n\"world\"\nEOF"},
		{"Null value", args{nil, noIndent}, "null"},
		{"Null struct", args{testNilPtr, noIndent}, "null"},
		{"List of integer", args{[]int{0, 1, 2, 3}, noIndent}, "[0,1,2,3]"},
		{"One level map", args{hclDict{"a": hclDict{"b": 10}}, noIndent}, "a {b=10}"},
		{"One level map (pretty)", args{hclDict{"a": hclDict{"b": 10}}, indent}, "a {\n  b = 10\n}"},
		{"Two level map 1", args{hclDict{"a": hclDict{"b": hclDict{"c": 10, "d": 20}}}, noIndent}, "a b {c=10 d=20}"},
		{"Two level map 1 (pretty)", args{hclDict{"a": hclDict{"b": hclDict{"c": 10, "d": 20}}}, indent}, "a b {\n  c = 10\n  d = 20\n}"},
		{"Two level map 2", args{hclDict{"a": hclDict{"b": hclDict{"c": 10, "d": 20}}, "e": 30}, noIndent}, "a b {c=10 d=20} e=30"},
		{"Two level map 2 (pretty)", args{hclDict{"a": hclDict{"b": hclDict{"c": 10, "d": 20}}, "e": 30}, indent}, "e = 30\n\na b {\n  c = 10\n  d = 20\n}"},
		{"Map", args{hclDict{"a": 0, "bb": 1}, noIndent}, "a=0 bb=1"},
		{"Map (pretty)", args{hclDict{"a": 0, "bb": 1}, indent}, "a  = 0\nbb = 1"},
		{"Structure (pretty)", args{test{"name", 1, true}, indent}, "Name   = \"name\"\nValue  = 1\npublic = true"},
		{"Structure Ptr (pretty)", args{&test{"name", 1, true}, indent}, "Name   = \"name\"\nValue  = 1\npublic = true"},
		{"Array of 1 structure (pretty)", args{[]test{{"name", 1, false}}, indent}, "Name   = \"name\"\nValue  = 1\npublic = false"},
		{"Array of 2 structures (pretty)", args{[]test{{"val1", 1, false}, {"val2", 1, true}}, indent}, "[\n  {\n    Name   = \"val1\"\n    Value  = 1\n    public = false\n  },\n  {\n    Name   = \"val2\"\n    Value  = 1\n    public = true\n  },\n]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := collections.MarshalGo(tt.args.value)
			assert.NoError(t, err)
			got, err := marshalHCL(value, true, true, "", tt.args.indent)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		hcl  string
		want interface{}
	}{
		{"Empty", "", hclDict{}},
		{"Empty list", "[]", hclList{}},
		{"List of int", "[1,2,3]", hclList{1, 2, 3}},
		{"Array of map", "a { b { c { d = 1 e = 2 }}}", hclDict{"a": hclDict{"b": hclDict{"c": hclDict{"d": 1, "e": 2}}}}},
		{"Map", fmt.Sprint(dictFixture), dictFixture},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out interface{}
			err := Unmarshal([]byte(tt.hcl), &out)
			assert.Equal(t, tt.want, out)
			assert.NoError(t, err)
		})
	}
}

func TestUnmarshalStrict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		hcl     string
		want    interface{}
		wantErr error
	}{
		{"Empty", "", map[string]interface{}{}, nil},
		{"Empty list", "[]", map[string]interface{}(nil), fmt.Errorf("reflect.Set: value of type []interface {} is not assignable to type map[string]interface {}")},
		{"List of int", "[1,2,3]", map[string]interface{}(nil), fmt.Errorf("reflect.Set: value of type []interface {} is not assignable to type map[string]interface {}")},
		{"Array of map", "a { b { c { d = 1 e = 2 }}}", map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": map[string]interface{}{"d": 1, "e": 2}}}}, nil},
		{"Map", fmt.Sprint(dictFixture), dictFixture.Native(), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out map[string]interface{}
			err := Unmarshal([]byte(tt.hcl), &out)
			assert.Equal(t, tt.want, out)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
