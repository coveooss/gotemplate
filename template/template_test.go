package template

import (
	"testing"
)

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
