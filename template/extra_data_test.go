package template

import (
	"reflect"
	"testing"

	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/yaml"
)

func Test_Data(t *testing.T) {
	t.Parallel()
	template := MustNewTemplate("", nil, "", nil)
	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple hcl", "a = 1", hcl.Dictionary{"a": 1}, false},
		{"Simple yaml", "b: 2", yaml.Dictionary{"b": 2}, false},
		{"Simple json", `{"c": 3}`, json.Dictionary{"c": 3}, false},
		{"Simple string", "string", "string", false},
		{"Error", "a = '", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.dataConverter(tt.test)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.fromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Template.fromData()\ngot : %[1]v (%[1]T)\nwant: %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}

func Test_YAML(t *testing.T) {
	t.Parallel()
	template := MustNewTemplate("", nil, "", nil)
	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple yaml", "b: 2", yaml.Dictionary{"b": 2}, false},
		{"Simple quoted string", `"string"`, "string", false},
		{"Simple string", "string", "string", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.yamlConverter(tt.test)
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

func Test_HCL(t *testing.T) {
	t.Parallel()
	template := MustNewTemplate("", nil, "", nil)
	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr bool
	}{
		{"Simple hcl", "a = 1", hcl.Dictionary{"a": 1}, false},
		{"Simple string", `"string"`, "string", false},
		{"Simple string", "string", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.hclConverter(tt.test)
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
