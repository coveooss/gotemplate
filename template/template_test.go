package template

import (
	"fmt"
	"testing"
)

func Test_getTargetFile(t *testing.T) {
	t.Parallel()

	type args struct {
		fileName   string
		sourcePath string
		targetPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple move", args{"/source/file", "/source", "/target"}, "/target/file"},
		{"Relative", args{"source/file", "/source", "/target"}, "/target/source/file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTargetFile(tt.args.fileName, tt.args.sourcePath, tt.args.targetPath); got != tt.want {
				t.Errorf("targetFile() = %v, want %v", got, tt.want)
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
		{"Undefined value", "@value", fmt.Errorf("Undefined value:1:4: Undefined value value in: @value")},
		{"2 Undefined values", "@(value1 + value2)", fmt.Errorf("2 Undefined values:1:8: Undefined value value1 in: @(value1 + value2)\n2 Undefined values:1:21: Undefined value value2 in: @(value1 + value2)")},
		{"Several errors", "@(value1)\n@nonExistingFunc()\n{{\n", fmt.Errorf("Several errors:2: function \"nonExistingFunc\" not defined in: @nonExistingFunc()\nSeveral errors:3: unexpected unclosed action in command in: {{\nSeveral errors:1:4: Undefined value value1 in: @(value1)")},
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
