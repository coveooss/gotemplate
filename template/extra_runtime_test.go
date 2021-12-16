package template

import (
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
			content: `@define("func")
			@-context()
			@-end
			@-include("func", 1, 2, 3)`,
			result: `{"ARGS":[1,2,3],"_":{"base":1},"base":1}`,
		},
		{
			name: "Override parent value",
			content: `@define("func")
			@-println("base =", base)
			@-println("_.base =", _.base)
			@-end
			@-include("func", data("base=over"))`,
			result: "base = over\n_.base = 1\n",
		}}
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
