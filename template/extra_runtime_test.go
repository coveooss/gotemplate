package template

import (
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/json"
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
