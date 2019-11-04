package template

import (
	"fmt"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/stretchr/testify/assert"
)

func TestFunctionTemplating(t *testing.T) {
	// Remove the run in parallel because we want to be sure that examples are generated using the right
	// output format
	template := MustNewTemplate(".", nil, "", nil)
	collections.SetListHelper(json.GenericListHelper)
	collections.SetDictionaryHelper(json.DictionaryHelper)
	template.completeExamples()

	for _, functionName := range template.getFunctions() {
		funcInfo := template.getFunction(functionName)
		for i, test := range funcInfo.examples {
			example := test
			t.Run(fmt.Sprintf("%s_#%d", funcInfo.name, i), func(t *testing.T) {
				if example.Razor != "" {
					appliedRazor, changed := template.applyRazor([]byte(example.Razor))
					assert.Equal(t, example.Template, string(appliedRazor), "Razor wasn't resolved correctly")
					assert.True(t, changed)
				}
				if example.Template != "" {
					appliedTemplate, err := template.ProcessContent(example.Template, ".")
					assert.Equal(t, example.Result, appliedTemplate, "Template wasn't resolved correctly")
					assert.NoError(t, err)
				}
			})
		}
	}
}
