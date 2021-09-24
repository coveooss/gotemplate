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
		{"Only exec on no timing", args{"This is a chat256", []string{`/chat256/miaou/b`, `/dummy/withtiming/e`}}, "This is a chat256", NONE},
		{"Only exec on b timing", args{"This is a dummy chat256", []string{`/chat256/miaou/b`, `/dummy/withtiming`}}, "This is a dummy miaou", BEGIN},
		{"Only exec on e timing", args{"This is a dummy chat256", []string{`/chat256/miaou`, `/dummy/smart/e`}}, "This is a smart chat256", END},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substitute(tt.args.content, tt.filter, InitReplacers(tt.args.replacers...)...); got != tt.want {
				t.Errorf("Substitute() = %v, want %v", got, tt.want)
			}
		})
	}
}
