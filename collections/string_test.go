package collections

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnIndent(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

func TestString_GetWordAtPosition(t *testing.T) {
	t.Parallel()

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
			got, pos := tt.s.GetWordAtPosition(tt.args.pos, tt.args.accept...)
			if got != tt.want {
				t.Errorf("String.GetWordAtPosition() got = %v, want %v", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("String.GetWordAtPosition() pos = %d, want %d", pos, tt.wantPos)
			}
			if got2 := tt.s.SelectWord(tt.args.pos, tt.args.accept...); got != got2 {
				t.Errorf("String.SelectWord() returns %v while String.GetWordAtPosition() returns %v", got, got2)
			}
		})
	}
}

func TestString_GetContextAtPosition(t *testing.T) {
	t.Parallel()

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
		{"Before ()", args{-1, "(", ")"}, "", -1},
		{"After ()", args{100, "(", ")"}, "", -1},
		{"Function()", args{9, "(", ")"}, "()", 8},
		{"A context (within parenthesis) should be returned", args{15, "(", ")"}, "(within parenthesis)", 10},
		{"A context [[within double bracket]] should be returned", args{15, "[[", "]]"}, "[[within double bracket]]", 10},
		{"A context [[from double bracket]] should be returned", args{24, "[[", ""}, "[[from double b", 10},
		{"A context [[to double bracket]] should be returned", args{22, "", "]]"}, "bracket]]", 22},
		{"A context (with (double level parenthesis))", args{22, "(", ")"}, "(double level parenthesis)", 16},
		{"A context (with no bracket)", args{19, "[", "]"}, "", -1},
		{"A context (with no enclosing context)", args{15, "", ""}, " ", 15},
		{"A context (outside of context)", args{1, "(", ")"}, "", -1},
		{"(context) after", args{12, "(", ")"}, "", -1},
		{"Test (with (parenthesis inside) of the context)", args{7, "(", ")"}, "(with (parenthesis inside) of the context)", 5},
		{"Test (with (parenthesis inside) unclosed", args{7, "(", ")"}, "", -1},
		{"Test (with (parenthesis inside) (closed) many time)))", args{7, "(", ")"}, "(with (parenthesis inside) (closed) many time)", 5},
		{"Test (with ((((((a lot of non closed)", args{7, "(", ")"}, "", -1},
		{"Test (with (several) (parenthesis (inside)) of the context) (excluded)", args{7, "(", ")"}, "(with (several) (parenthesis (inside)) of the context)", 5},
		{"Test | with same | left and | right", args{7, "|", "|"}, "| with same |", 5},
		{"A context [[from [[double]] bracket]] [[with a little extra]]", args{12, "[[", "]]"}, "[[from [[double]] bracket]]", 10},
		{"A context [[from [[double]] bracket [[unclosed]]", args{12, "[[", "]]"}, "", -1},
		{"A context [[from [[double]] bracket [[extra]] closed]] many times]]]]", args{12, "[[", "]]"}, "[[from [[double]] bracket [[extra]] closed]]", 10},
	}
	for _, tt := range tests {
		t.Run(tt.s.Str(), func(t *testing.T) {
			got, pos := tt.s.GetContextAtPosition(tt.args.pos, tt.args.left, tt.args.right)
			if got != tt.want {
				t.Errorf("String.GetContextAtPosition() got = %v, want %v", got, tt.want)
			}
			if pos != tt.wantPos {
				t.Errorf("String.GetContextAtPosition() pos = %v, want %v", pos, tt.wantPos)
			}
			if got2 := tt.s.SelectContext(tt.args.pos, tt.args.left, tt.args.right); got != got2 {
				t.Errorf("String.SelectContext() returns %v while String.SelectWord() returns %v", got, got2)
			}
		})
	}
}

func TestString_Protect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      String
		want      String
		wantArray StringArray
	}{
		{"", "", nil},
		{`"This is a string"`, `"♠0"`, StringArray{`"This is a string"`}},
		{`A test with a "single string"`, `A test with a "♠0"`, StringArray{`"single string"`}},
		{"A test with `backtick string`", `A test with "♠0"`, StringArray{"`backtick string`"}},
		{`Non closed "string`, `Non closed "string`, nil},
		{`Non closed "with" escape "\"`, `Non closed "♠0" escape "\"`, StringArray{`"with"`}},
		{`This contains two "string1" and "string2"`, `This contains two "♠0" and "♠1"`, StringArray{`"string1"`, `"string2"`}},
		{"A mix of `backtick` and \"regular\" string", `A mix of "♠0" and "♠1" string`, StringArray{"`backtick`", `"regular"`}},
		{"A confused one of `backtick with \"` and \"regular with \\\" quoted and ` inside\" string", "A confused one of \"♠0\" and \"♠1\" string", StringArray{"`backtick with \"`", "\"regular with \\\" quoted and ` inside\""}},
		{`A string with "false \\\\\\" inside"`, `A string with "♠0" inside"`, StringArray{`"false \\\\\\"`}},
		{`A string with "true \\\\\\\" inside"`, `A string with "♠0"`, StringArray{`"true \\\\\\\" inside"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name.Str(), func(t *testing.T) {
			got, array := tt.name.Protect()
			if got != tt.want {
				t.Errorf("String.Protect() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(array, tt.wantArray) {
				t.Errorf("String.Protect() array = %v, want %v", array, tt.wantArray)
			}

			restored := got.RestoreProtected(array)
			if tt.name != restored {
				t.Errorf("String.RestoreProtected() got %v, want %v", restored, tt.name)
			}
		})
	}
}

func TestString_ParseBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value String
		want  bool
	}{
		{"", false},
		{"1", true},
		{"0", false},
		{"F", false},
		{"False", false},
		{"FALSE", false},
		{"No", false},
		{"N", false},
		{"T", true},
		{"true", true},
		{"on", true},
		{"OFF", false},
		{"Whatever", true},
		{"YES", true},
	}
	for _, tt := range tests {
		t.Run(tt.value.Str(), func(t *testing.T) {
			if got := tt.value.ParseBool(); got != tt.want {
				t.Errorf("ParseBoolFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_IndexAll(t *testing.T) {
	tests := []struct {
		name       string
		s          String
		substr     string
		wantResult []int
	}{
		{"Both empty", "", "", nil},
		{"Empty substr", "aaa", "", nil},
		{"Empty source", "", "aa", nil},
		{"Single", "abaabaaa", "a", []int{0, 2, 3, 5, 6, 7}},
		{"Double", "abaabaaabaaaa", "aa", []int{2, 5, 9, 11}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := tt.s.IndexAll(tt.substr); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("String.IndexAll() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestString_AddLineNumber(t *testing.T) {
	tests := []struct {
		name  string
		s     String
		space int
		want  String
	}{
		{"Empty", "", 0, "1 "},
		{"Just a newline", "\n", 0, "1 \n2 "},
		{"Several lines", "Line 1\nLine 2\nLine 3\n", 0, "1 Line 1\n2 Line 2\n3 Line 3\n4 "},
		{"Several lines", "Line 1\nLine 2\nLine 3\n", 4, "   1 Line 1\n   2 Line 2\n   3 Line 3\n   4 "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.AddLineNumber(tt.space); got != tt.want {
				t.Errorf("String.AddLineNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
