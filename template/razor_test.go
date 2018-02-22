package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/errors"
)

func TestTemplate_applyRazor(t *testing.T) {
	context := make(map[string]interface{})
	template := NewTemplate(context, "", Math|Sprig, true, true)

	files, err := filepath.Glob("../doc_test/*.md")
	if err != nil {
		t.Fatalf("Unable to read test files (documentation in ../doc)")
		t.Fail()
	}

	type test struct {
		name   string
		path   string
		razor  string
		render string
	}

	ifExist := func(path string) string {
		if _, err := os.Stat(path); err != nil {
			return ""
		}
		return path
	}

	load := func(path string) []byte { return errors.Must(ioutil.ReadFile(path)).([]byte) }

	tests := make([]test, 0, len(files))
	for _, file := range files {
		path := strings.TrimSuffix(file, ".md")
		tests = append(tests, test{
			name:   filepath.Base(path),
			path:   file,
			razor:  ifExist(path + ".razor"),
			render: ifExist(path + ".rendered"),
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template.RazorSyntax = tt.razor != ""

			content := load(tt.path)
			if tt.razor != "" {
				result := load(tt.razor)
				got := template.applyRazor(content)
				if !reflect.DeepEqual(got, result) {
					t.Errorf("Template.applyRazor()\n\nExpected:\n%s\n\nGot:\n%s\n", string(result), string(got))
				}
			}

			got, err := template.ProcessContent(string(content), tt.name)
			if err != nil {
				t.Errorf("Template.ProcessContent(), err=%v", err)
			}

			if tt.render != "" {
				result := string(load(tt.render))
				if !reflect.DeepEqual(got, result) {
					t.Errorf("Template.ProcessContent()\n\nExpected:\n%s\n\nGot:\n%s\n", result, got)
				}
			}
		})
	}
}
