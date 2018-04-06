package implementation

/*
import (
	"reflect"
	"testing"
)

var simpleDict = dictImpl{
	"a": 1,
	"b": "2",
	"c": []int{1, 2, 3},
}

func Test_dict_String(t *testing.T) {
	tests := []struct {
		name string
		d    genDict
		want string
	}{
		{"nil", nil, "dict[]"},
		{"Empty List", genDict{}, "dict[]"},
		{"Map", genDict(dictFixture).Clone().AsMap(), "dict[float:1.23 int:123 list:[1 two] listInt:[1 2 3] map:dict[sub1:1 sub2:two] mapInt:dict[1:1 2:two] string:Foo bar]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("dict.String():\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}

func Test_dictImpl_String(t *testing.T) {
	tests := []struct {
		name string
		d    dictImpl
		want string
	}{
		{"nil", nil, "gen_dict[]"},
		{"Empty List", dictImpl{}, "gen_dict[]"},
		{"Map", dictImpl(dictFixture), "gen_dict[float:1.23 int:123 list:[1 two] listInt:[1 2 3] map:gen_dict[sub1:1 sub2:two] mapInt:gen_dict[1:1 2:two] string:Foo bar]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("dict.String():\ngot  %v\nwant %v", got, tt.want)
			}
		})
	}
}
func TestAsDictionary(t *testing.T) {
	tests := []struct {
		name    string
		object  interface{}
		want    iDict
		wantErr bool
	}{
		{"Simple a", map[string]interface{}{"a": 1}, simpleDict.Clone("a"), false},
		{"Nil", nil, dictImpl{}, false},
		{"Not working string", "From string", nil, true},
		{"Not working int", 1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := asDictionary(tt.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsDictionary():\ngot  %[1]v (%[1]T)\nwant %[2]v (%[2]T)", got, tt.want)
			}
		})
	}
}
*/
