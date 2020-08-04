package template

import (
	"testing"

	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/yaml"
	"github.com/stretchr/testify/assert"
)

func Test_Data(t *testing.T) {
	t.Parallel()
	template := MustNewTemplate("", nil, "", nil)
	tests := []struct {
		name    string
		test    string
		want    interface{}
		wantErr string
	}{
		{"Simple hcl", "a = 1", hcl.Dictionary{"a": 1}, ""},
		{"Simple yaml", "b: 2", yaml.Dictionary{"b": 2}, ""},
		{"Simple json", `{"c": 3}`, json.Dictionary{"c": 3}, ""},
		{"Simple string", "string", "string", ""},
		{"Error", "a = '", nil, "\n   1 a = '\n\nTrying !json: invalid character 'a' looking for beginning of value\nTrying hcl: At 1:5: illegal char"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.dataConverter(tt.test)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func Test_YAML(t *testing.T) {
	t.Parallel()
	template := MustNewTemplate("", nil, "", nil)
	tests := []struct {
		name string
		test string
		want interface{}
	}{
		{"Simple yaml", "b: 2", yaml.Dictionary{"b": 2}},
		{"Simple quoted string", `"string"`, "string"},
		{"Simple string", "string", "string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.yamlConverter(tt.test)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
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
		wantErr string
	}{
		{"Simple hcl", "a = 1", hcl.Dictionary{"a": 1}, ""},
		{"Simple string", `"string"`, "string", ""},
		{"Simple string", "string", nil, "\n   1 string\n\nAt 2:1: key 'string' expected start of object ('{') or assignment ('=')"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := template.hclConverter(tt.test)
			assert.Equal(t, tt.want, got)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
