package template

import (
	"fmt"
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

func Test_templateWithErrors(t *testing.T) {
	t.Parallel()

	template, _ := NewTemplate(".", nil, "", nil)
	template.SetOption(StrictErrorCheck, true)
	tests := []struct {
		name    string
		content string
		err     error
	}{
		{"Empty template", "", nil},
		{"Non closed brace", "{{", fmt.Errorf("Non closed brace:1: unexpected unclosed action in command in: {{")},
		{"Non opened brace", "}}", nil},
		{"Undefined value", "@value", fmt.Errorf("Undefined value:1: contains undefined value(s) in: @value")},
		{"2 Undefined values", "@(value1 + value2)", fmt.Errorf("2 Undefined values:1: contains undefined value(s) in: @(value1 + value2)")},
		{"Several errors", "@(value1)\n@non_Existing_Func()\n{{\n", fmt.Errorf("Several errors:2: function \"non_Existing_Func\" not defined in: @non_Existing_Func()\nSeveral errors:3: unexpected unclosed action in command in: {{\nSeveral errors:1: contains undefined value(s) in: @(value1)")},
		{"undefined variable", "@(value_non_existing)", fmt.Errorf("undefined variable:1: contains undefined value(s) in: @(value_non_existing)")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := template.ProcessContent(tt.content, tt.name); err != tt.err {
				if err != nil && tt.err != nil && err.Error() == tt.err.Error() {
					return
				}
				t.Errorf("ProcessContent()=\n%v\n\nWant:\n%v", err, tt.err)
			}
		})
	}
}
