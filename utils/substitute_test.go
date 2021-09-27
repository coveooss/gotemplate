package utils

import "testing"

func TestSubstitute(t *testing.T) {
	type args struct {
		content   string
		replacers []string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		filter SubstituteTiming
	}{
		{"Simple regex", args{"This is a test", []string{`/\b(\w{2})\b/$1$1`, `/\b(\w)\b/-$1-`}}, "This isis -a- test", NONE},
		{"Only exec on no timing", args{"This is a cat256", []string{`/cat256/meow/b`, `/dummy/withtiming/e`}}, "This is a cat256", NONE},
		{"Only exec on b timing", args{"This is a dummy cat256", []string{`/cat256/meow/b`, `/dummy/withtiming`}}, "This is a dummy meow", BEGIN},
		{"Only exec on e timing", args{"This is a dummy cat256", []string{`/cat256/meow`, `/dummy/smart/e`}}, "This is a smart cat256", END},
		{"Protection on begin", args{"Infamous @sha256", []string{`/@sha256/meow/p`}}, "Infamous _=!meow!=_", BEGIN},
		{"Protection on end", args{"Infamous _=!meow!=_", []string{`/@sha256/meow/p`}}, "Infamous @sha256", END},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substitute(tt.args.content, tt.filter, InitReplacers(tt.args.replacers...)...); got != tt.want {
				t.Errorf("Substitute() = %v, want %v", got, tt.want)
			}
		})
	}
}
