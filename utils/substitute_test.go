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
		{"Simple regex", args{"This is a test", []string{`/\b(\w{2})\b/$1$1`, `/\b(\w)\b/-$1-`}}, "This isis -a- test", NoTiming},
		{"Only exec on no timing", args{"This is a cat256", []string{`/cat256/meow/b`, `/dummy/withtiming/e`}}, "This is a cat256", NoTiming},
		{"Only exec on b timing", args{"This is a dummy cat256", []string{`/cat256/meow/b`, `/dummy/withtiming`}}, "This is a dummy meow", BeginTiming},
		{"Only exec on e timing", args{"This is a dummy cat256", []string{`/cat256/meow`, `/dummy/smart/e`}}, "This is a smart cat256", EndTiming},
		{"Protection on begin", args{"Infamous @cat256", []string{`/@cat256/meow/p`}}, "Infamous _=!meow!=_", BeginTiming},
		{"Protection on end", args{"Infamous _=!meow!=_", []string{`/@cat256/meow/p`}}, "Infamous @cat256", EndTiming},
		{"Mix protect + begin", args{"Infamous @cat256", []string{`/@cat256/meow/p`, `/Infamous/famous/b`}}, "famous _=!meow!=_", BeginTiming},
		{"Mix protect + end", args{"famous _=!meow!=_", []string{`/@cat256/meow/p`, `/famous/The real/e`}}, "The real @cat256", EndTiming},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substitute(tt.args.content, tt.filter, InitReplacers(tt.args.replacers...)...); got != tt.want {
				t.Errorf("Substitute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitReplacers_panics(t *testing.T) {

	tests := []struct {
		name      string
		replacers []string
	}{
		{"Empty replacer", []string{""}},
		{"Not enough params", []string{"/potato"}},
		{"Too many params", []string{"/potato/banana/mango/anana"}},
		{"Bad timing information", []string{"/potato/banana/mango"}},
		{"Timing mix be", []string{"/potato/banana/be"}},
		{"Timing mix bp", []string{"/potato/banana/bp"}},
		{"Timing mix ep", []string{"/potato/banana/ep"}},
		{"Timing mix bep", []string{"/potato/banana/bep"}},
		{"Protected literal is regex", []string{`/\d+[potato]/banana/p`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { recover() }()
			res := InitReplacers(tt.replacers...)
			t.Log(tt.name, "should have panicked, instead returned", res)
			t.Fail()
		})
	}
}
