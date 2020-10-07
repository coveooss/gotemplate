package template

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTargetFile(t *testing.T) {
	t.Parallel()

	type args struct {
		fileName   string
		sourcePath string
		targetPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple move", args{"/source/file", "/source", "/target"}, "/target/file"},
		{"Relative", args{"source/file", "/source", "/target"}, "/target/source/file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTargetFile(tt.args.fileName, tt.args.sourcePath, tt.args.targetPath); got != tt.want {
				t.Errorf("targetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateFilesOverwrite(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "No templating",
			content: "Test",
			result:  "Test",
		},
		{
			name:    "Basic Razor",
			content: "@(3+4)",
			result:  "7",
		},
		{
			name:    "Basic Gotemplate",
			content: "{{ add 3 4 }}",
			result:  "7",
		},
		{
			name:    "Razor with double delimiter",
			content: "@(3+4)\n@@testValue",
			result:  "7\n@testValue",
		},
		{
			name:    "Razor only double delimiter",
			content: "@@testValue",
			result:  "@testValue",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tempfile, err := ioutil.TempFile("", "")
			assert.Nil(t, err)
			defer os.Remove(tempfile.Name())

			_, err = tempfile.WriteString(tt.content)
			assert.Nil(t, err)

			template, _ := NewTemplate(path.Dir(tempfile.Name()), nil, "", nil)
			template.SetOption(Overwrite, true)
			template.ProcessTemplates("", "", tempfile.Name())

			result, err := ioutil.ReadFile(tempfile.Name())
			assert.Nil(t, err)
			assert.Equal(t, tt.result, string(result))
		})
	}
}

func TestTemplateAddFunctions(t *testing.T) {
	t.Parallel()

	getValue := func() interface{} {
		return "This Is My Value"
	}

	options := DefaultOptions()
	options[StrictErrorCheck] = true
	template, _ := NewTemplate(".", map[string]interface{}{}, "", options)
	template.AddFunctions(map[string]interface{}{"getValue": getValue}, "Inline", nil)
	_, ok := template.functions["getValue"]
	assert.True(t, ok, "getValue is not defined")
	result, err := template.ProcessContent("@getValue()", ".")
	assert.Nil(t, err)
	assert.Equal(t, "This Is My Value", result)
}
