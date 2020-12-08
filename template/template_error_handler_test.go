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
		name         string
		code         string
		err          string
		noValueCount int
	}{
		{"Undefined", `@value`, noValueError, 1},
		{"Undefined with nil test", "@value@if(missing) whatever;", noValueError, 1},
		{
			"2 Undefined with nil test", `
			@value
			@if(missing) whatever;
			@otherValue
			`,
			noValueError, 2,
		},
		{
			"Invalid assignation (undefined variable)", `
			@{var} := $value
			@{var}
			`,
			":2: undefined variable \"$value\" in: \t\t\t@{var} := $value", 0,
		},
		{
			"Invalid assignation (missing parameters)", `
			@{var} := 3 + default()
			@{var}
			`,
			":2: undefined variable \"$value\" in: \t\t\t@{var} := $value", 0,
		},
		{
			"Invalid assignation (bad function)", `
			@{var} := non_existing_func()
			@{var}
			`,
			":2: undefined variable \"$value\" in: \t\t\t@{var} := $value", 0,
		},
		{
			"Invalid if statement", `
			@if ($value)
				text
			@endif
			`,
			":2: undefined variable \"$value\" in: \t\t\t@if ($value)", 0,
		},
		{
			"Invalid with statement", `
			@with ($value)
				text
			@endif
			`,
			":2: undefined variable \"$value\" in: \t\t\t@with ($value)", 0,
		},
		{
			"Invalid foreach statement", `
			@for ($i := $value)
				text
			@end
			`,
			noValueError, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := MustNewTemplate(".", nil, "", DefaultOptions().Set(StrictErrorCheck).Unset(AcceptNoValue))
			_, err := template.ProcessContent(tt.code, "")
			if tt.err == "" {
				assert.NoError(t, err)
			} else if tt.noValueCount > 0 {
				assert.Contains(t, err.Error(), tt.err)
				assert.Equal(t, strings.Count(err.Error(), noValue), tt.noValueCount)
			} else {
				assert.EqualError(t, err, tt.err)
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
		{"Undefined value", "@value", fmt.Errorf("template: Undefined value:: contains undefined value(s)\n1 <no value>")},
		{"2 Undefined values", "@(value1 + value2)", fmt.Errorf("template: 2 Undefined values:: contains undefined value(s)\n1 <no value>")},
		{"Several errors", "@(value1)\n@non_Existing_Func()\n{{\n", fmt.Errorf("Several errors:2: function \"non_Existing_Func\" not defined in: @non_Existing_Func()\nSeveral errors:3: unexpected unclosed action in command in: {{")},
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
