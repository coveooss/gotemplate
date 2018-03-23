package types

import (
	"fmt"
	"reflect"
	"testing"
)

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

func TestString_ToTitle(t *testing.T) {
	tests := []struct {
		s    String
		want string
	}{
		{"Hello world", "HELLO WORLD"},
	}
	for _, tt := range tests {
		t.Run(string(tt.s), func(t *testing.T) {
			if got := tt.s.ToTitle(); string(got) != tt.want {
				t.Errorf("String.ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapString(t *testing.T) {
	sample := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
	tests := []struct {
		s          string
		width      int
		wantResult string
	}{
		{sample, 1, "Lorem\nipsum\ndolor\nsit\namet,\nconsectetur\nadipiscing\nelit,\nsed\ndo\neiusmod\ntempor\nincididunt\nut\nlabore\net\ndolore\nmagna\naliqua."},
		{sample, 5, "Lorem\nipsum\ndolor\nsit\namet,\nconsectetur\nadipiscing\nelit,\nsed do\neiusmod\ntempor\nincididunt\nut\nlabore\net\ndolore\nmagna\naliqua."},
		{sample, 10, "Lorem ipsum\ndolor sit\namet,\nconsectetur\nadipiscing\nelit, sed\ndo eiusmod\ntempor\nincididunt\nut labore\net dolore\nmagna\naliqua."},
		{sample, 20, "Lorem ipsum dolor sit\namet, consectetur\nadipiscing elit, sed\ndo eiusmod tempor\nincididunt ut labore\net dolore magna\naliqua."},
		{sample, 30, "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor\nincididunt ut labore et dolore\nmagna aliqua."},
		{sample, 40, "Lorem ipsum dolor sit amet, consectetur\nadipiscing elit, sed do eiusmod tempor\nincididunt ut labore et dolore magna\naliqua."},
		{sample, 50, "Lorem ipsum dolor sit amet, consectetur adipiscing\nelit, sed do eiusmod tempor incididunt ut labore et\ndolore magna aliqua."},
		{sample, 75, "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod\ntempor incididunt ut labore et dolore magna aliqua."},
		{sample, 100, "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore\net dolore magna aliqua."},
		{sample, 125, sample},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test with %d", tt.width), func(t *testing.T) {
			if gotResult := WrapString(tt.s, tt.width); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("WrapString() =\n%q, want\n%q", gotResult, tt.wantResult)
			}
		})
	}
}
