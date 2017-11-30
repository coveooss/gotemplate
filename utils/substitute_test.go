package utils

import "testing"

func TestSubstitute(t *testing.T) {
	type args struct {
		content   string
		replacers []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple regex", args{"This is a test", []string{`/\b(\w{2})\b/$1$1`, `/\b(\w)\b/-$1-`}}, "This isis -a- test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substitute(tt.args.content, InitReplacers(tt.args.replacers...)...); got != tt.want {
				t.Errorf("Substitute() = %v, want %v", got, tt.want)
			}
		})
	}
}
