package template

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/bmatcuk/doublestar"
	"github.com/coveooss/multilogger"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestTemplate_applyRazor(t *testing.T) {
	t.Parallel()
	dmp := diffmatchpatch.New()
	TemplateLog = multilogger.New("test")
	template := MustNewTemplate("../docs_tests", nil, "", nil)
	files, err := doublestar.Glob(filepath.Join(template.folder, "**/*.md"))
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

	template.options[AcceptNoValue] = true

	load := func(path string) []byte { return must(os.ReadFile(path)).([]byte) }

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

			processContent := func(renderingDisabled bool) string {
				var got string
				var err error
				func() {
					defer func() {
						if rec := recover(); rec != nil {
							err = fmt.Errorf("Template.ProcessContent() panic=%v\n%s", rec, string(debug.Stack()))
						}
					}()
					template.options[RenderingDisabled] = renderingDisabled
					got, err = template.ProcessContent(string(load(tt.path)), tt.path)
				}()
				if err != nil {
					t.Errorf("Template.ProcessContent(), err=%v", err)
				}
				return got
			}

			if tt.razor != "" {
				result := string(load(tt.razor))
				got := processContent(true)
				if !reflect.DeepEqual(got, result) {
					diffs := dmp.DiffMain(string(result), string(got), true)
					t.Errorf("Differences on Razor result for %s\n%s", tt.razor, dmp.DiffPrettyText(diffs))
				}
			}

			if tt.render != "" {
				result := string(load(tt.render))
				got := processContent(false)
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
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
		})
	}
}

func TestInvocation(t *testing.T) {
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{
			"Func call",
			"@func(1,2,3)",
			"{{ func 1 2 3 }}",
		},
		{
			"Method call",
			"@object.func(1,2,3)",
			"{{ $.object.func 1 2 3 }}",
		},
		{
			"Method call on result",
			"@object.func(1,2).func2(3)",
			"{{ ($.object.func 1 2).func2 3 }}",
		},
		{
			"Double invocation",
			"@func1().func2()",
			"{{ func1.func2 }}",
		},
		{
			"Double invocation with params",
			"@func1(1).func2(2)",
			"{{ (func1 1).func2 2 }}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
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
			"Local replacement 1",
			"@$a = 2",
			`{{- $a = 2 }}`,
		},
		{
			"Local replacement 2",
			"@{a = 2}",
			"{{- $a = 2 }}",
		},
		{
			"Local replacement 3",
			"@$a = 2",
			`{{- $a = 2 }}`,
		},
		{
			"Local replacement 4",
			"@{a.b.c} = 2",
			`{{- set $a.b "c" 2 }}`,
		},
		{
			"Global assign 1",
			`@a := "test"`,
			`{{- set $ "a" "test" }}`,
		},
		{
			"Global assign 2",
			`@.a := "test"`,
			`{{- set . "a" "test" }}`,
		},
		{
			"Global assign 3",
			`@$.a := "test"`,
			`{{- set $ "a" "test" }}`,
		},
		{
			"Replacement of global value",
			`@a = "test"`,
			`{{- assertWarning (not (isNil $.a)) "$.a does not exist, use := to declare new variable" }}{{- set $ "a" "test" }}`,
		},
		{
			"Global assign with non standard identifier characters",
			`@12t%!e#st- := "test"`,
			`{{- set $ "12t%!e#st-" "test" }}`,
		},
		{
			"Global assign with sub objects",
			`@a.b.c.d.e := "test"`,
			`{{- set $.a.b.c.d "e" "test" }}`,
		},
		{
			"Assignment operator 1",
			`@{a} += 10`,
			`{{- $a = add $a 10 }}`,
		},
		{
			"Assignment operator 2",
			`@{a *= 10}`,
			`{{- $a = mul $a 10 }}`,
		},
		{
			"Assignment operator 3",
			`@a <<= 10`,
			`{{- assertWarning (not (isNil $.a)) "$.a does not exist, use := to declare new variable" }}{{- set $ "a" (lshift $.a 10) }}`,
		},
		{
			"Assignment operator 3",
			`@a.b.c <<= 10`,
			`{{- assertWarning (not (isNil $.a.b.c)) "$.a.b.c does not exist, use := to declare new variable" }}{{- set $.a.b "c" (lshift $.a.b.c 10) }}`,
		},
		{
			"Assignment operator 4",
			`@{a} »= 10`,
			`{{- $a = rshift $a 10 }}`,
		},
		{
			"Assignment operator 5",
			`@{a} ÷= 2`,
			`{{- $a = div $a 2 }}`,
		},
		{
			"Assignment operator 6",
			`@.a.b *= 4*2`,
			`{{- assertWarning (not (isNil .a.b)) ".a.b does not exist, use := to declare new variable" }}{{- set .a "b" (mul .a.b (mul 4 2)) }}`,
		},
		{
			"Assignment operator 7",
			`@$.a.b *= 4`,
			`{{- assertWarning (not (isNil $.a.b)) "$.a.b does not exist, use := to declare new variable" }}{{- set $.a "b" (mul $.a.b 4) }}`,
		},
		{
			"Assignment operator 8",
			`@$a *= 4`,
			`{{- $a = mul $a 4 }}`,
		},
		{
			"Assignment operator local sub",
			`@{a.b.c} ÷= 2`,
			`{{- set $a.b "c" (div $a.b.c 2) }}`,
		},
		{
			"Assignment operator with expression",
			`@{a} /= 2 * 3`,
			`{{- $a = div $a (mul 2 3) }}`,
		},
		{
			"Global assignment operator with expression",
			`@a %= 2 / 3`,
			`{{- assertWarning (not (isNil $.a)) "$.a does not exist, use := to declare new variable" }}{{- set $ "a" (mod $.a (div 2 3)) }}`,
		},
		{
			"Assignment operator with index",
			`@{a} += $text[3:]`,
			`{{- $a = add $a (slice $text 3 -1) }}`,
		},
		{
			"Assignment with @",
			`@a := "How do you @handle this"`,
			`{{- set $ "a" "How do you @handle this" }}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
		})
	}
}

func TestAssignWithValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		razor  string
		want   string
		result string
	}{
		{
			"Assignment with with",
			`
			@d := dict("v0", 0)
			@-with (d)
				@.v1 := 1
				@.v2 := 2
			@-end
			@--d
			`,
			`
			{{- set $ "d" (dict "v0" 0) }}
			{{- with $.d }}
				{{- set . "v1" 1 }}
				{{- set . "v2" 2 }}
			{{- end }}
			{{- $.d -}}
			`,
			`{"v0":0,"v1":1,"v2":2}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			got, changed := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
			assert.True(t, changed)
			r, err := template.ProcessContent(string(got), ".")
			assert.NoError(t, err)
			assert.Equal(t, r, string(tt.result))
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
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
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
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
		})
	}
}

func TestMultilineStringProtect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		razor string
		want  string
	}{
		{
			"Empty",
			"This is an empty string ``",
			"This is an empty string ``",
		},
		{
			"Expression within quote",
			"`@(1+2)`",
			"`{{ add 1 2 }}`",
		},
		{
			"String withing expression",
			"@func(`@(1+2)`)",
			"{{ func `@(1+2)` }}",
		},
		{
			"Expression within multiline quote",
			"`\n@(1+2)\n`",
			"`\n@(1+2)\n`",
		},
		{
			"Expression within empty quotes",
			"``\n@(1+2)\n``",
			"``\n{{ add 1 2 }}\n``",
		},
		{
			"Expression within markdown (md)",
			"```razor\n@(1+2)\n```",
			"```razor\n{{ add 1 2 }}\n```",
		},
		{
			"Expression with escaped @ in multiline string",
			"`\n@@Not changed\n`",
			"`\n@@Not changed\n`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			got, _ := template.applyRazor([]byte(tt.razor))
			assert.Equal(t, tt.want, string(got), tt.razor)
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

func TestReservedKeywords(t *testing.T) {
	// Ensure that protected go keyword are processed correctly
	t.Parallel()

	template := MustNewTemplate(".", nil, "", nil)
	for i, keyword := range reservedKeywords {
		if keyword == "..." {
			continue
		}
		t.Run(keyword, func(t *testing.T) {
			code := fmt.Sprintf("@var := %s + %d", keyword, i)
			got, _ := template.applyRazor([]byte(code))
			switch keyword {
			case "$":
				assert.Equal(t, `{{- set $ "var" (add $ 0) }}`, string(got), code)
			default:
				assert.Equal(t, fmt.Sprintf(`{{- set $ "var" (add $.%s %d) }}`, keyword, i), string(got), code)
			}
		})
	}
}
