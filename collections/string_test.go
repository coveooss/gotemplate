package collections

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

func TestString_FindWord(t *testing.T) {
	type args struct {
		pos    int
		accept []string
	}
	tests := []struct {
		s       String
		args    args
		want    String
		wantPos int
	}{
		{"", args{0, nil}, "", -1},
		{"A single character", args{0, nil}, "A", 0},
		{"This a test", args{0, nil}, "This", 0},
		{"This is a secode test", args{3, nil}, "This", 0},
		{"Over", args{20, nil}, "", -1},
		{"Find the second word", args{5, nil}, "the", 5},
		{"Find the $third word", args{10, nil}, "$third", 9},
		{"Find the ($a.value) word", args{10, []string{"."}}, "$a.value", 10},
		{"Find the ($a.value[0]) word", args{10, []string{"."}}, "$a.value", 10},
		{"Find the ($a.value[10]) word", args{10, []string{".", "[]"}}, "$a.value[10]", 10},
		{"Match a space", args{5, nil}, "", 5},
	}
	for _, tt := range tests {
		t.Run(tt.s.Str(), func(t *testing.T) {
			got, pos := tt.s.FindWord(tt.args.pos, tt.args.accept...)
			if got != tt.want {
				t.Errorf("String.FindWord() got = %v, want %v", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("String.FindWord() pos = %d, want %d", pos, tt.wantPos)
			}
			if got2 := tt.s.SelectWord(tt.args.pos, tt.args.accept...); got != got2 {
				t.Errorf("String.SelectWord() returns %v while String.FindWord() returns %v", got, got2)
			}
		})
	}
}

func TestString_FindContext(t *testing.T) {
	type args struct {
		pos   int
		left  string
		right string
	}
	tests := []struct {
		s       String
		args    args
		want    String
		wantPos int
	}{
		{"", args{0, "", ""}, "", -1},
		{"", args{5, "", ""}, "", -1},
		{"A context (within parenthesis) should be returned", args{15, "(", ")"}, "(within parenthesis)", 10},
		{"A context [[within double bracket]] should be returned", args{15, "[[", "]]"}, "[[within double bracket]]", 10},
		{"A context [[from double bracket]] should be returned", args{24, "[[", ""}, "[[from double b", 10},
		{"A context [[to double bracket]] should be returned", args{22, "", "]]"}, "bracket]]", 22},
		{"A context (with (double level parenthesis))", args{22, "(", ")"}, "(double level parenthesis)", 16},
		{"A context (with no bracket)", args{19, "[", "]"}, "", -1},
		{"A context (with no enclosing context)", args{15, "", ""}, " ", 15},
		// 123456789012345678901234567890123456789012345678901234567890
		//          1         2         3         4         5         6
	}
	for _, tt := range tests {
		t.Run(tt.s.Str(), func(t *testing.T) {
			got, pos := tt.s.FindContext(tt.args.pos, tt.args.left, tt.args.right)
			if got != tt.want {
				t.Errorf("String.FindContext() got = %v, want %v", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("String.FindContext() pos = %v, want %v", pos, tt.wantPos)
			}
			if got2 := tt.s.SelectContext(tt.args.pos, tt.args.left, tt.args.right); got != got2 {
				t.Errorf("String.SelectWord() returns %v while String.SelectWord() returns %v", got, got2)
			}
		})
	}
}
