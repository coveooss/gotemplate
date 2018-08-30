package template

import (
	"reflect"
	"testing"
)

func TestProtectString(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		wantArray []string
	}{
		{"", "", nil},
		{`"This is a string"`, `"♠0"`, []string{`"This is a string"`}},
		{`A test with a "single string"`, `A test with a "♠0"`, []string{`"single string"`}},
		{"A test with `backtick string`", `A test with "♠0"`, []string{"`backtick string`"}},
		{`Non closed "string`, `Non closed "string`, nil},
		{`This contains two "string1" and "string2"`, `This contains two "♠0" and "♠1"`, []string{`"string1"`, `"string2"`}},
		{"A mix of `backtick` and \"regular\" string", `A mix of "♠0" and "♠1" string`, []string{"`backtick`", `"regular"`}},
		{"A confused one of `backtick with \"` and \"regular with \\\" quoted and ` inside\" string", "A confused one of \"♠0\" and \"♠1\" string", []string{"`backtick with \"`", "\"regular with \\\" quoted and ` inside\""}},
		{`A string with "false \\\\\\" inside"`, `A string with "♠0" inside"`, []string{`"false \\\\\\"`}},
		{`A string with "true \\\\\\\" inside"`, `A string with "♠0"`, []string{`"true \\\\\\\" inside"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, array := ProtectString(tt.name)
			if got != tt.want {
				t.Errorf("ProtectString() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(array, tt.wantArray) {
				t.Errorf("ProtectString() array = %v, want %v", array, tt.wantArray)
			}

			restored := RestoreProtectedString(got, array)
			if tt.name != restored {
				t.Errorf("RestoreProtectedString() got %v, want %v", restored, tt.name)
			}
		})
	}
}
