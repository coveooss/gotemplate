package template

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", nil)
			result, err := template.ProcessContent(tt.content, tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.result, result)
		})
	}
}
