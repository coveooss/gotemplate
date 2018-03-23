package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/errors"
	logging "github.com/op/go-logging"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestTemplate_applyRazor(t *testing.T) {
	dmp := diffmatchpatch.New()
	SetLogLevel(logging.WARNING)
	template := NewTemplate("../docs/doc_test", nil, "", nil)
	files, err := filepath.Glob(filepath.Join(template.folder, "*.md"))
	if err != nil {
		t.Fatalf("Unable to read test files (documentation in %s)", template.folder)
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
			template.options[Razor] = tt.razor != ""

			content := load(tt.path)
			if tt.razor != "" {
				result := load(tt.razor)
				got := template.applyRazor(content)
				if !reflect.DeepEqual(got, result) {
					diffs := dmp.DiffMain(string(result), string(got), true)
					t.Errorf("Differences on Razor result for %s\n%s", tt.razor, dmp.DiffPrettyText(diffs))
				}
			}

			got, err := template.ProcessContent(string(content), tt.path)
			if err != nil {
				t.Errorf("Template.ProcessContent(), err=%v", err)
			}

			if tt.render != "" {
				result := string(load(tt.render))
				if !reflect.DeepEqual(got, result) {
					diffs := dmp.DiffMain(string(result), string(got), true)
					t.Errorf("Differences on Rendered for %s\n%s", tt.render, dmp.DiffPrettyText(diffs))
				}
			}
		})
	}
}
