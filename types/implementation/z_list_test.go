package implementation

/*
import (
	"reflect"
	"testing"
)

func Test_list_String(t *testing.T) {
	tests := []struct {
		name string
		l    genList
		want string
	}{
		{"Nil", nil, "[]"},
		{"Empty List", genList{}, "[]"},
		{"List of int", genList{1, 2, 3}, "[1 2 3]"},
		{"List of string", *strFixture, `[Hello World, I'm Foo Bar!]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("genList.String():\ngot %v\nwant %v", got, tt.want)
			}
		})
	}
}
func TestAsGenericList(t *testing.T) {
	tests := []struct {
		name       string
		object     interface{}
		wantResult iList
		wantErr    bool
	}{
		{"Simple ", genList{1, 2, 3}, (*listImpl)(al(genList{1, 2, 3})), false},
		{"Integer genList", []int{1, 2, 3}, (*listImpl)(al(genList{1, 2, 3})), false},
		{"String genList", []string{"1", "two"}, (*listImpl)(al(genList{"1", "two"})), false},
		{"Nil", nil, nil, false},
		{"Not working string", "From string", nil, true},
		{"Not working int", 1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := asGenericList(tt.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsGenericList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AsGenericList = %[1]v (%[1]T), want %[2]v (%[2]T)", gotResult, tt.wantResult)
			}
		})
	}
}
*/
