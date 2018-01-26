package utils

import "testing"

func TestUnIndent(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Indented with tab", args{`
			Hello

			World
			end!
			`}, "\nHello\n\nWorld\nend!\n"},
		{"Indented with spaces", args{`
                Hello

                World
                end!
                `}, "\nHello\n\nWorld\nend!\n"},
		{"Normal string", args{"Hello World!"}, "Hello World!"},
		{"Normal string prefixed with spaces", args{"  Hello World!"}, "  Hello World!"},
		{"Indented with mixed spaces", args{`
			Hello

	        World
			end!
			`}, "\n\t\t\tHello\n\n\t        World\n\t\t\tend!\n\t\t\t"},
		{"One line indented", args{"\nHello\n   World\n"}, "\nHello\n   World\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnIndent(tt.args.s); got != tt.want {
				t.Errorf("UnIndent() = %v, want %v", got, tt.want)
			}
		})
	}
}
