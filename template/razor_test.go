package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/coveo/gotemplate/v3/collections"
	"github.com/coveo/gotemplate/v3/json"
	logging "github.com/op/go-logging"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestTemplate_applyRazor(t *testing.T) {
	t.Parallel()
	dmp := diffmatchpatch.New()
	SetLogLevel(logging.WARNING)
	template := MustNewTemplate("../docs/doc_test", nil, "", nil)
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

	collections.ListHelper = json.GenericListHelper
	collections.DictionaryHelper = json.DictionaryHelper
	template.options[AcceptNoValue] = true

	load := func(path string) []byte { return must(ioutil.ReadFile(path)).([]byte) }

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
				got, _ := template.applyRazor(content)
				if !reflect.DeepEqual(got, result) {
					diffs := dmp.DiffMain(string(result), string(got), true)
					t.Errorf("Differences on Razor result for %s\n%s", tt.razor, dmp.DiffPrettyText(diffs))
				}
			}

			var got string
			var err error
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						err = fmt.Errorf("Template.ProcessContent() panic=%v\n%s", rec, string(debug.Stack()))
					}
				}()
				got, err = template.ProcessContent(string(content), tt.path)
			}()

			if err != nil {
				t.Errorf("Template.ProcessContent(), err=%v", err)
			} else if tt.render != "" {
				result := string(load(tt.render))
				if !reflect.DeepEqual(got, result) {
					diffs := dmp.DiffMain(string(result), string(got), true)
					t.Errorf("Differences on Rendered for %s\n%s", tt.render, dmp.DiffPrettyText(diffs))
				}
			}
		})
	}
}

func TestBase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{"Empty", "", ""},
		{
			"Simple global variable",
			"@Hello",
			"{{ $.Hello }}",
		},
		{
			"Email",
			"Hello john.doe@company.com",
			"Hello john.doe@company.com",
		},
		{
			"Literal",
			"Hello john.doe@@company",
			"Hello john.doe@company",
		},
		{
			"No razor",
			"{{ gotemplate }}",
			"{{ gotemplate }}",
		},
		{
			"Mix",
			"@test {{ gotemplate }}",
			"{{ $.test }} {{ gotemplate }}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			if got, _ := template.applyRazor([]byte(tt.razor)); string(got) != tt.want {
				t.Errorf("applyRazor() = got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestInvocation(t *testing.T) {
	tests := []struct {
		name       string
		debugLevel int
		razor      string
		want       string
	}{
		{
			"Func call", 2,
			"@func(1,2,3)",
			"{{ func 1 2 3 }}",
		},
		{
			"Method call", 2,
			"@object.func(1,2,3)",
			"{{ $.object.func 1 2 3 }}",
		},
		{
			"Method call on result", 2,
			"@object.func(1,2).func2(3)",
			"{{ ($.object.func 1 2).func2 3 }}",
		},
		{
			"Double invocation", 2,
			"@func1().func2()",
			"{{ func1.func2 }}",
		},
		{
			"Double invocation with params", 6,
			"@func1(1).func2(2)",
			"{{ (func1 1).func2 2 }}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging.SetLevel(logging.Level(tt.debugLevel), loggerInternal)
			defer func() { logging.SetLevel(logging.Level(2), loggerInternal) }()
			template := MustNewTemplate(".", nil, "", nil)
			if got, _ := template.applyRazor([]byte(tt.razor)); string(got) != tt.want {
				t.Errorf("applyRazor() = got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAssign(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{
			"Local assign",
			"@{a} := 2",
			"{{- $a := 2 }}",
		},
		{
			"Local assign 2",
			"@{a := 2}",
			"{{- $a := 2 }}",
		},
		{
			"Local assign 3",
			"@$a := 2",
			`{{- $a := 2 }}`,
		},
		{
			"Global assign",
			`@a := "test"`,
			`{{- assertWarning (isNil $.a) "$.a has already been declared, use = to overwrite existing value" }}{{- set $ "a" "test" }}`,
		},
		{
			"Deprecated local assign with no other razor",
			`$a := "test"`,
			`$a := "test"`,
		},
		{
			"Deprecated local assign",
			`@test; $a := $.test`,
			`{{ $.test }} {{- $a := $.test }}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			if got, _ := template.applyRazor([]byte(tt.razor)); string(got) != tt.want {
				t.Errorf("applyRazor() = got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAutoWrap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{
			"Base",
			"Before @autoWrap(to(10)) after",
			`{{ join "" (formatList "Before %v after" (to 10)) }}`,
		},
		{
			"With newline",
			"Before @<aWrap(to(10)) after",
			`{{- $.NEWLINE }}{{ join "\n" (formatList "Before %v after" (to 10)) }}`,
		},
		{
			"With space eater",
			"Before @--awrap(to(10)) after",
			`{{- join "" (formatList "Before %v after" (to 10)) -}}`,
		},
		{
			"With error",
			"Before @--awrap(to(10) after",
			"Before {{- awrap to(10 -}} after",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			if got, _ := template.applyRazor([]byte(tt.razor)); string(got) != tt.want {
				t.Errorf("applyRazor() = got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSpaceEater(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{
			"Base",
			"@value",
			`{{ $.value }}`,
		},
		{
			"Before",
			"@-value",
			`{{- $.value }}`,
		},
		{
			"After",
			"@_-value",
			`{{ $.value -}}`,
		},
		{
			"Both",
			"@--value",
			`{{- $.value -}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			if got, _ := template.applyRazor([]byte(tt.razor)); string(got) != tt.want {
				t.Errorf("applyRazor() = got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		code string
		want string
		err  error
	}{
		{"Empty", `@data("")`, "", nil},
		{"Integer", `@data("1")`, "1", nil},
		{"Hcl", `@data("a = 1 b = 2")`, "a=1 b=2", nil},
		{"Hcl type", `@typeOf(data("a = 1 b = 2"))`, "hcl.hclDict", nil},
		{"Hcl kind", `@kindOf(data("a = 1 b = 2"))`, "map", nil},
		{"Invalid", "@typeOf(data(`\"a\": 1, \"b\": 2`))", `"<RUN_ERROR>"`, fmt.Errorf("")},
		{"Json", "@typeOf(data(`{\"a\": 1, \"b\": 2}`))", "json.jsonDict", nil},
		{"Yaml", "@typeOf(data(`a: 1\nb: 2`))", "yaml.yamlDict", nil},
		{"Flexible Hcl", "@typeOf(data(`a = 1 b = hello`))", "yaml.yamlDict", nil}, // TODO: Change that to hcl
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", DefaultOptions().Set(StrictErrorCheck))
			got, err := template.ProcessContent(tt.code, "")
			assert.Equal(t, tt.want, got)
			if tt.err == nil {
				assert.NoError(t, err)
			} else if tt.err.Error() == "" {
				assert.Error(t, err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
