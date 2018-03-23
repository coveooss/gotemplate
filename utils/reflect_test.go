package utils

import (
	"reflect"
	"testing"
)

type empty struct{}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want bool
	}{
		{"False", false, true},
		{"True", true, false},
		{"Nil", nil, true},
		{"String", "Hello", false},
		{"Empty string", "", true},
		{"Zero", 0, true},
		{"Uint", uint(10), false},
		{"Floating point zero", 0.0, true},
		{"Floating point", 10.0, false},
		{"Floating point negative", -10.0, false},
		{"Interface to nil", interface{}(nil), true},
		{"Interface to zero", interface{}(0), true},
		{"Interface to empty string", interface{}(""), true},
		{"Interface to string", interface{}("Foo"), false},
		{"Integer", 10, false},
		{"Empty struct", empty{}, false},
		{"Empty list", []string{}, true},
		{"List", []string{"Hello"}, false},
		{"Empty map", map[string]int{}, true},
		{"Map", map[string]int{"Hello": 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyValue(reflect.ValueOf(tt.arg)); got != tt.want {
				t.Errorf("IsEmptyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExported(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"Private", "privateValue", false},
		{"Public", "PublicValue", true},
		{"Zero number", "0", false},
		{"Positive number", "1", false},
		{"Underscore", "_test", false},
		{"Id with space", "ID 1", true},
		{"AllCap", "ALL_CAP_ID", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExported(tt.id); got != tt.want {
				t.Errorf("IsExported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIfUndef(t *testing.T) {
	const def = "default"
	type empty struct{}
	tests := []struct {
		name string
		arg  interface{}
		want interface{}
	}{
		{"False", false, false},
		{"True", true, true},
		{"Nil", nil, def},
		{"String", "Hello", "Hello"},
		{"Empty string", "", ""},
		{"Zero", 0, 0},
		{"Integer", 10, 10},
		{"Empty struct", empty{}, empty{}},
		{"Empty list", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfUndef(def, tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IfUndef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIfUndefNoValue(t *testing.T) {
	def := "default"
	t.Run("No Value", func(t *testing.T) {
		if got := IfUndef(def); !reflect.DeepEqual(got, def) {
			t.Errorf("IfUndef() = %v, want %v", got, def)
		}
	})
}

func TestIfManyValues(t *testing.T) {
	const def = "default"
	list := []interface{}{1, 2, 3}
	t.Run("Many values", func(t *testing.T) {
		if got := IfUndef(def, list...); !reflect.DeepEqual(got, list) {
			t.Errorf("IfUndef() = %v, want %v", got, list)
		}
	})
}

func TestIIf(t *testing.T) {
	tests := []struct {
		name      string
		testValue interface{}
		want      interface{}
	}{
		{"False", false, 2},
		{"True", true, 1},
		{"Nil", nil, 2},
		{"String", "Hello", 1},
		{"Empty string", "", 2},
		{"Zero", 0, 2},
		{"Integer", 10, 1},
		{"Empty struct", empty{}, 1},
		{"Empty list", []string{}, 2},
		{"List", []string{"Hello"}, 1},
		{"Empty map", map[string]int{}, 2},
		{"Map", map[string]int{"Hello": 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IIf(tt.testValue, 1, 2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IIf() = %v, want %v", got, tt.want)
			}
		})
	}
}
