package template

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/yaml"
	"github.com/stretchr/testify/assert"
)

func init() {
	collections.SetListHelper(json.GenericListHelper)
	collections.SetDictionaryHelper(json.DictionaryHelper)
}

func TestRuntime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "Global context variable must be available",
			content: `@define("func")@Math.Pi@end @-include("func")`,
			result:  "3.141592653589793",
		},
		{
			name:    "With Args",
			content: `@define("func")@Math.Pi@end @-include("func", 1, 2, 3)`,
			result:  "3.141592653589793",
		},
		{
			name: "Get context",
			content: `
			@-define("func")
				@base = 2
				@-println("base =", base)
				@-println("ARGS =", ARGS)
				@-println("_.base =", _.base)
			@-end

			@-include("func", 1, 2, 3)
			@--println("base =", base)
			`,
			result: "base = 2\nARGS = [1,2,3]\n_.base = 1\nbase = 1\n",
		},
		{
			name: "Override parent value",
			content: `
			@-define("func")
				@-println("base =", base)
				@-println("_.base =", _.base)
			@-end
			@-include("func", data("base=over"))`,
			result: "base = over\n_.base = 1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			template.Add("base", 1)
			result, err := template.ProcessContent(tt.content, tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestMultilineError(t *testing.T) {
	// Ensure that multiline errors are not truncated after the first line
	t.Parallel()

	template := MustNewTemplate(".", nil, "", nil)
	template.SetOption(StrictErrorCheck, true)
	_, err := template.ProcessContent(`@run("ls -DONT EXIST")`, "bad param")
	assert.Error(t, err)
	if err != nil {
		assert.GreaterOrEqual(t, len(toStringClass(err.Error()).Lines()), 2)
	}
}

func TestInclude(t *testing.T) {
	t.Parallel()

	// This test confirm that validation mode is not changing when we invoke a sub template.
	// In the second call, value does not exist, so if strict validation is applied, the call would fail.
	template := MustNewTemplate(".", nil, "", nil)
	template.SetOption(StrictErrorCheck, true)
	x, err := template.ProcessContent(`
		@--define("sub")
			Hello @if (value) @value;
		@-end

		@-include("sub", dict("value", 123))
		@<--include("sub")
	`, "include")
	assert.Equal(t, "Hello 123\nHello ", x)
	assert.NoError(t, err)
}

func TestIncludeContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "Include without arguments resolves a value from the global context",
			content: `{{- define "sub" }}{{ .base }}{{ end }}{{ include "sub" }}`,
			result:  "1",
		},
		{
			name:    "Include arguments take precedence over the include context",
			content: `{{- define "sub" }}{{ .NL }}{{ end }}{{ include "sub" (dict "NL" "argument") }}`,
			result:  "argument",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			template.SetOption(StrictErrorCheck, true)
			template.Add("base", 1)
			result, err := template.ProcessContent(tt.content, tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestIncludeResolvesGlobalValueAfterNestedInclude(t *testing.T) {
	t.Parallel()

	template := MustNewTemplate(".", nil, "", nil)
	template.Add("kubernetes", map[string]interface{}{"namespace": "monitoring"})
	result, err := template.ProcessContent(`{{- define "labels" -}}
managed-by: gotemplate
{{- end -}}
{{- define "alert-labels" -}}
{{ include "labels" }}
namespace: {{ .kubernetes.namespace }}
{{- end -}}
{{ include "alert-labels" }}`, "nested include")
	assert.NoError(t, err)
	assert.Equal(t, "managed-by: gotemplate\nnamespace: monitoring", result)
}

func TestIncludeContextResolution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "Caller context passed as argument",
			content: `{{- define "sub" }}{{ .base }}{{ end }}{{ include "sub" . }}`,
			result:  "1",
		},
		{
			name:    "Missing key falls back to the global context",
			content: `{{- define "sub" }}{{ .base }}{{ end }}{{ include "sub" }}`,
			result:  "1",
		},
		{
			name: "Parent chain through nested includes",
			content: `
				{{- define "inner" }}{{ .base }} {{ ._._.base }}{{ end }}
				{{- define "outer" }}{{ include "inner" (dict "base" 3) }}{{ end }}
				{{- include "outer" (dict "base" 2) }}`,
			result: "3 1",
		},
		{
			name: "Caller context forwarded through nested includes",
			content: `
				{{- define "inner" }}{{ .base }}{{ end }}
				{{- define "outer" }}{{ include "inner" . }}{{ end }}
				{{- include "outer" . }}`,
			result: "1",
		},
		{
			name: "Missing key in a nested include falls back to the global context",
			content: `
				{{- define "inner" }}{{ .base }}{{ end }}
				{{- define "outer" }}{{ include "inner" }}{{ end }}
				{{- include "outer" }}`,
			result: "1",
		},
		{
			name: "Includes after an include of the caller context",
			content: `
				{{- define "first" }}{{ .base }}{{ end }}
				{{- define "second" }}{{ .base }}{{ end }}
				{{- include "first" . }} {{ include "second" }} {{ include "first" . }}`,
			result: "1 1 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			template.Add("base", 1)
			result, err := template.ProcessContent(tt.content, tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestIncludeCallerContextDoesNotReferenceItself(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "Repeated include of a defined template with the caller context",
			content: `{{- define "sub" }}{{ .base }}{{ end }}` + strings.Repeat(`{{ include "sub" . }}`, 8),
			result:  strings.Repeat("1", 8),
		},
		{
			name: "Repeated include of the parent context",
			content: `
				{{- define "inner" }}{{ .base }}{{ end }}
				{{- define "outer" }}{{ include "inner" ._ }}{{ include "inner" ._ }}{{ include "inner" ._ }}{{ include "inner" ._ }}{{ end }}
				{{- include "outer" . }}{{ include "outer" . }}`,
			result: "11111111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			template.Add("base", 1)
			keysBefore := fmt.Sprint(template.Context().GetKeys())
			result, err := template.ProcessContent(tt.content, tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
			assert.Equal(t, keysBefore, fmt.Sprint(template.Context().GetKeys()))
		})
	}
}

func TestMergeIntoCopyMatchesMergeWithoutModifyingTarget(t *testing.T) {
	t.Parallel()

	newTarget := func() collections.IDictionary {
		return collections.AsDictionary(map[string]interface{}{
			"kept":           "target",
			"scalar":         "target",
			"nested":         map[string]interface{}{"kept": "target", "deeper": map[string]interface{}{"kept": "target"}},
			"dictOverScalar": map[string]interface{}{"kept": "target"},
			"scalarOverDict": "target",
		})
	}
	source := collections.AsDictionary(map[string]interface{}{
		"added":          "source",
		"scalar":         "source",
		"nested":         map[string]interface{}{"kept": "source", "added": "source", "deeper": map[string]interface{}{"added": "source"}},
		"dictOverScalar": "source",
		"scalarOverDict": map[string]interface{}{"added": "source"},
	})

	expected := fmt.Sprint(newTarget().Merge(source))
	target := newTarget()
	untouched := fmt.Sprint(target)

	assert.Equal(t, expected, fmt.Sprint(mergeIntoCopy(target, source)))
	assert.Equal(t, untouched, fmt.Sprint(target))
}

func TestIncludeCallerContextErrorRecovery(t *testing.T) {
	t.Parallel()

	// Each failing line makes the error handler render the file again, so the include runs again
	template := MustNewTemplate(".", nil, "", nil)
	template.SetOption(StrictErrorCheck, true)
	template.Add("base", 1)
	keysBefore := fmt.Sprint(template.Context().GetKeys())
	_, err := template.ProcessContent(strings.Repeat("{{ include \"missing\" . }}\n", 8), "error recovery")
	assert.ErrorContains(t, err, `unable to find a template or a file named missing`)
	assert.Equal(t, keysBefore, fmt.Sprint(template.Context().GetKeys()))
}

func TestExec(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		args     []interface{}
		expected interface{}
	}{
		{
			name:     "should return bare text as is",
			script:   `echo -n 'hello world'`,
			args:     []interface{}{},
			expected: "hello world",
		},
		{
			name:   "should template in arguments",
			script: `echo -n 'hello @name'`,
			args: []interface{}{
				map[string]interface{}{
					"name": "bob",
				},
			},
			expected: "hello bob",
		},
		{
			name:     "should not template output",
			script:   `printf 'hello %s(2 + 2)' '@'`,
			args:     []interface{}{},
			expected: "hello @(2 + 2)",
		},
		{
			name: "should parse json output",
			script: `
				echo '{
					"foo": "bar",
					"num": 1,
					"bool": true,
					"nope": null,
					"with_at": "@@yay"
				}'`,
			args: []interface{}{},
			expected: json.Dictionary{
				"foo":  "bar",
				"num":  1,
				"bool": true,
				// Nils are converted to empty dictionaries for some reason
				"nope":    json.Dictionary{},
				"with_at": "@yay",
			},
		},
		{
			name: "should parse yaml output",
			script: `
				echo 'foo: bar'
				echo 'num: 1'
				echo 'bool: true'
				echo 'nope: null'
				echo 'with_at: "@@yay"'`,
			args: []interface{}{},
			expected: yaml.Dictionary{
				"foo":  "bar",
				"num":  1,
				"bool": true,
				// Nils are converted to empty dictionaries for some reason
				"nope":    yaml.Dictionary{},
				"with_at": "@yay",
			},
		},
		{
			name: "should parse hcl output",
			script: `
				echo 'foo = "bar"'
				echo 'num = 1'
				echo 'bool = true'
				echo 'with_at = "@@yay"'`,
			args: []interface{}{},
			expected: hcl.Dictionary{
				"foo":     "bar",
				"num":     1,
				"bool":    true,
				"with_at": "@yay",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			template := MustNewTemplate("", nil, "", nil)
			res, err := template.exec(test.script, test.args...)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
