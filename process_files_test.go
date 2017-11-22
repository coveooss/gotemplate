package main

import (
	"testing"
)

func Test_substitute(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name       string
		args       args
		parameters []string
		want       string
	}{
		{"Simple regex", args{"This is a test"}, []string{`\b(\w{2})\b/$1$1`, `\b(\w)\b/-$1-`}, "This isis -a- test"},
	}
	for _, tt := range tests {
		substitutes = &tt.parameters
		t.Run(tt.name, func(t *testing.T) {
			if got := substitute(tt.args.content); got != tt.want {
				t.Errorf("substitute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTargetFile(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTargetFile(tt.args.fileName, tt.args.sourcePath, tt.args.targetPath); got != tt.want {
				t.Errorf("targetFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
