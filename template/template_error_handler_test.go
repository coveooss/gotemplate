package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateErrorHandling(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		code     string
		err      string
		errCount int
	}{
		{
			"Undefined", `@value`,
			noValueError, 1,
		},
		{
			"Undefined with nil test", "@value@if(missing) whatever;",
			noValueError, 1,
		},
		{
			"2 Undefined with nil test", `
			   @value
			   @if(missing) whatever;
			   @otherValue
			`,
			noValueError, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", DefaultOptions().Set(StrictErrorCheck).Unset(AcceptNoValue))
			_, err := template.ProcessContent(tt.code, "")
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.errCount, strings.Count(err.Error(), tt.err))
			}
		})
	}
}
