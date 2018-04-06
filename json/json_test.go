package json

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_list_String(t *testing.T) {
	tests := []struct {
		name string
		l    jsonList
		want string
	}{
		{"Nil", nil, "null"},
		{"Empty list", jsonList{}, "[]"},
		{"List of int", jsonList{1, 2, 3}, "[1,2,3]"},
		{"List of string", strFixture, `["Hello","World,","I'm","Foo","Bar!"]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("jsonList.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_String(t *testing.T) {
	tests := []struct {
		name string
		d    jsonDict
		want string
	}{
		{"nil", nil, "null"},
		{"Empty dict", jsonDict{}, "{}"},
		{"Map", dictFixture, `{"float":1.23,"int":123,"list":[1,"two"],"listInt":[1,2,3],"map":{"sub1":1,"sub2":"two"},"mapInt":{"1":1,"2":"two"},"string":"Foo bar"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("jsonList.String():\n  %v\n  %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    interface{}
		wantErr bool
	}{
		{"Empty", "", nil, true},
		{"Empty list", "[]", al(jsonList{}), false},
		{"List of int", "[1,2,3]", al(jsonList{1, 2, 3}), false},
		{"Map", fmt.Sprint(dictFixture), dictFixture, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out interface{}
			err := Unmarshal([]byte(tt.json), &out)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(out, tt.want) {
				t.Errorf("Unmarshal:\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", out, tt.want)
			}
		})
	}
}

func TestUnmarshalToMap(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    interface{}
		wantErr bool
	}{
		{"Empty", "", nil, true},
		{"Empty list", "[]", nil, true},
		{"List of int", "[1,2,3]", nil, true},
		{"Map", fmt.Sprint(dictFixture), dictFixture.AsMap(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := make(map[string]interface{})
			err := Unmarshal([]byte(tt.json), &out)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(out, tt.want) {
				t.Errorf("Unmarshal:\n got %[1]v (%[1]T)\nwant %[2]v (%[2]T)", out, tt.want)
			}
		})
	}
}
