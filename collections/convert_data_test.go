package collections

import (
	"fmt"
	"reflect"
	"testing"
)

type dictionary = map[string]interface{}

func TestToNativeRepresentation(t *testing.T) {
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
			if got := ToNativeRepresentation(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNativeRepresentation()\ngot : %v\nwant: %v", got, tt.want)
				for k, v := range tt.want.(dictionary) {
					if reflect.DeepEqual(v, got.(dictionary)[k]) {
						continue
					}
					t.Errorf("key %v: %T vs %T", k, v, got.(dictionary)[k])
				}

			}
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
			if got := quote(tt.arg); got != tt.want {
				t.Errorf("quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleToNativeRepresentation() {
	type Y struct {
		S string `hcl:"string,omitempty"`
	}
	type X struct {
		Y  `hcl:",squash"`
		A  int     `hcl:"a,omitempty"`
		PS *string `hcl:"string_pointer,omitempty"`
		PB *bool   `hcl:"bool_pointer,omitempty"`
	}

	var a X
	fmt.Println(ToNativeRepresentation(a))
	a.A = 10
	fmt.Println(ToNativeRepresentation(a))
	a.S = "Hello"
	fmt.Println(ToNativeRepresentation(a))
	x := "World"
	a.PS = &x
	b := false
	a.PB = &b
	fmt.Println(ToNativeRepresentation(a))

	// Output:
	// map[]
	// map[a:10]
	// map[a:10 string:Hello]
	// map[a:10 bool_pointer:false string:Hello string_pointer:World]
}
