package collections

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dictionary = map[string]interface{}

func TestMarshalGo(t *testing.T) {
	t.Parallel()

	type SubStruct struct {
		U int64
		I interface{}
	}
	type a struct {
		private int
		I       int
		F       float64
		S       string
		A       []interface{}
		M       dictionary
		SS      SubStruct
	}
	tests := []struct {
		name string
		args interface{}
		want interface{}
	}{
		{"Struct conversion", a{
			private: 0,
			I:       123,
			F:       1.23,
			S:       "123",
			A:       []interface{}{1, "2"},
			M: dictionary{
				"a": "a",
				"b": 2,
			},
			SS: SubStruct{64, "Foo"},
		}, dictionary{
			"I": 123,
			"F": float64(1.23),
			"S": "123",
			"A": []interface{}{1, "2"},
			"M": dictionary{
				"a": "a",
				"b": 2,
			},
			"SS": dictionary{
				"U": int64(64),
				"I": "Foo",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalGo(tt.args)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func Test_quote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"Simple value", "Foo", "Foo"},
		{"Simple value", "Foo Bar", `"Foo Bar"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, quote(tt.arg))
		})
	}
}

// This object has only private attributes.
type customMarshallStruct struct {
	privateInteger     int
	privateString      string
	privateBoolPointer *bool
}

// MarshalGo implements a custom marshaler that allows the object to convert its private attributes.
func (s customMarshallStruct) MarshalGo(v interface{}) (interface{}, error) {
	return map[string]interface{}{
		"int":    s.privateInteger,
		"string": s.privateString,
		"bool":   *s.privateBoolPointer,
	}, nil
}

func Test_CustomMarshal(t *testing.T) {
	b := true
	e := customMarshallStruct{10, "Hello", &b}
	converted, err := MarshalGo(e)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"int": 10, "string": "Hello", "bool": true}, converted)
}

func ExampleMarshalGo() {
	type Struct struct {
		Integer       int     `hcl:"int,omitempty"`
		StringPointer *string `hcl:"string_pointer,omitempty"`
		BoolPointer   *bool   `hcl:"bool_pointer,omitempty"`
	}

	var main Struct

	fmt.Println(MarshalGo(main))
	main.Integer = 10
	fmt.Println(MarshalGo(main))
	word := "World"
	main.StringPointer = &word
	bool := false
	main.BoolPointer = &bool
	fmt.Println(MarshalGo(main))

	// Output:
	// map[] <nil>
	// map[int:10] <nil>
	// map[bool_pointer:false int:10 string_pointer:World] <nil>
}

func ExampleMarshalGo_withSubStruct() {
	type WithoutStructTag struct {
		String1 string
	}

	type WithStructTag struct {
		String2 string `hcl:"string,omitempty"`
	}

	type MainStruct struct {
		WithoutStructTag
		WithStructTag `hcl:",squash"`
		Omit          WithStructTag `hcl:",omitempty"`
		Integer       int           `hcl:"int,omitempty"`
	}

	var main MainStruct

	main.Integer = 10
	fmt.Println(MarshalGo(main))
	main.String1 = "Hello"
	main.String2 = "World"
	fmt.Println(MarshalGo(main))

	// Output:
	// map[WithoutStructTag:map[String1:] int:10] <nil>
	// map[WithoutStructTag:map[String1:Hello] int:10 string:World] <nil>
}

func ExampleMarshalGo_withList() {
	type Element struct {
		Name  string
		Value int
	}

	type Struct struct {
		Elements []Element `hcl:"elements,omitempty"`
	}

	var main Struct
	main.Elements = []Element{{"value1", 1}, {"value2", 2}, {"value3", 3}}
	fmt.Println(MarshalGo(main))

	// Output:
	// map[elements:[map[Name:value1 Value:1] map[Name:value2 Value:2] map[Name:value3 Value:3]]] <nil>
}

func ExampleMarshalGo_withKey() {
	type Element struct {
		Name  string `hcl:",key"`
		Value int    `hcl:"value"`
	}

	type Struct struct {
		Elements []Element `hcl:"elements,omitempty"`
	}

	var main Struct
	main.Elements = []Element{{"v1", 1}, {"v2", 2}, {"v3", 3}}
	fmt.Println(MarshalGo(main))

	// Output:
	// map[elements:[map[v1:map[value:1]] map[v2:map[value:2]] map[v3:map[value:3]]]] <nil>
}

func ExampleMarshalGo_withError() {
	type Element struct {
		Name  string `hcl:",key"`
		Value int    `hcl:"value,key"`
	}

	type Struct struct {
		Elements []Element `hcl:"elements,omitempty"`
	}

	var main Struct
	main.Elements = []Element{{"v1", 1}, {"v2", 2}, {"v3", 3}}
	fmt.Println(MarshalGo(main))

	// Output:
	// <nil> Multiple keys defined on struct 'Element' ('Name' and 'Value')
}
