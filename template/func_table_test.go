package template

import (
	"testing"

	assertion "github.com/stretchr/testify/assert"
)

func TestFunctionTemplating(t *testing.T) {
	t.Parallel()

	template := MustNewTemplate(".", nil, "", nil)
	for _, functionName := range template.getFunctions() {
		funcInfo := template.getFunction(functionName)
		t.Run(funcInfo.name, func(t *testing.T) {
			t.Parallel()
			if funcInfo.exampleRazor == "" {
				t.Skipf("%s skipped because it has no example", funcInfo.name)
			}
			appliedRazor, _ := template.applyRazor([]byte(funcInfo.exampleRazor))
			assertion.Equal(t, funcInfo.exampleTemplate, string(appliedRazor), "Razor wasn't resolved correctly")
			appliedTemplate, _ := template.ProcessContent(funcInfo.exampleTemplate, ".")
			assertion.Equal(t, funcInfo.exampleResult, appliedTemplate, "Template wasn't resolved correctly")
		})

	}
}
