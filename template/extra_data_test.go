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

func Test_Python(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		data       interface{}
		want       string
		wantPretty string
	}{
		{"true", hcl.Dictionary{"a": true},
			`{"a":True}`,
			"{\n  \"a\": True\n}"},
		{"false", yaml.Dictionary{"a": false},
			`{"a":False}`,
			"{\n  \"a\": False\n}"},
		{"null", json.Dictionary{"a": nil},
			`{"a":None}`,
			"{\n  \"a\": None\n}"},
		{"list", json.List{"Hello", 1, 2.3, true, false, nil},
			`["Hello",1,2.3,True,False,None]`,
			"[\n  \"Hello\",\n  1,\n  2.3,\n  True,\n  False,\n  None\n]"},
		{"combo", json.List{"Hello", hcl.Dictionary{"a": true, "b": false, "c": nil}},
			`["Hello",{"a":True,"b":False,"c":None}]`,
			"[\n  \"Hello\",\n  {\n    \"a\": True,\n    \"b\": False,\n    \"c\": None\n  }\n]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normal, err := toPython(tt.data)
			assert.Equal(t, tt.want, normal)
			assert.NoError(t, err)
			pretty, err := toPrettyPython(tt.data)
			assert.Equal(t, tt.wantPretty, pretty)
			assert.NoError(t, err)
		})
	}
}
