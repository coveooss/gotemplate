package implementation

import "testing"

func Test_baseDict_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    baseDict
		want string
	}{
		{"nil", nil, "dict[]"},
		{"Empty List", baseDict{}, "dict[]"},
		{"Map", dictFixture, "dict[float:1.23 int:123 list:[1 two] listInt:[1 2 3] map:dict[sub1:1 sub2:two] mapInt:dict[1:1 2:two] string:Foo bar]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("dict.String():\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}
