package template

import (
	"reflect"
	"testing"
)

func Test_fromData(t *testing.T) {
	template := NewTemplate("", nil, "", nil)

	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple hcl", "a = 1", map[string]interface{}{"a": 1}, false},
		{"Simple yaml", "b: 2", map[string]interface{}{"b": 2}, false},
		{"Simple json", `"c": 3`, map[string]interface{}{"c": 3}, false},
		{"Simple string", "string", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.fromData(tt.test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.fromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Template.fromData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromYAML(t *testing.T) {
	template := NewTemplate("", nil, "", nil)

	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple yaml", "b: 2", map[string]interface{}{"b": 2}, false},
		{"Simple quoted string", `"string"`, "string", false},
		{"Simple string", "string", "string", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.fromYAML(tt.test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.fromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Template.fromData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromHCL(t *testing.T) {
	template := NewTemplate("", nil, "", nil)

	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple hcl", "a = 1", map[string]interface{}{"a": 1}, false},
		{"Simple string", `"string"`, "string", false},
		{"Simple string", "string", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.fromHCL(tt.test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.fromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Template.fromData() = %v, want %v", got, tt.want)
			}
		})
	}
}
