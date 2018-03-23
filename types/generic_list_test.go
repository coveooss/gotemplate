package types

import (
	"reflect"
	"testing"
)

func TestAsGenericList(t *testing.T) {
	tests := []struct {
		name       string
		object     interface{}
		wantResult IGenericList
		wantErr    bool
	}{
		{"Simple ", GenericList{1, 2, 3}, GenericList{1, 2, 3}, false},
		{"Integer list", []int{1, 2, 3}, GenericList{1, 2, 3}, false},
		{"String list", []string{"1", "two"}, GenericList{"1", "two"}, false},
		{"Nil", nil, nil, false},
		{"Not working string", "From string", nil, true},
		{"Not working int", 1, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := AsGenericList(tt.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsGenericList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AsGenericList() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestGenericList_Get(t *testing.T) {
	list := GenericList{1, "2", "Three", GenericList{1, 2}, "last"}
	tests := []struct {
		name  string
		list  GenericList
		index int
		want  interface{}
	}{
		{"First", list, 0, 1},
		{"Last", list, list.Len() - 1, "last"},
		{"Negative index", list, -1, nil},
		{"Out of bound", list, 1000, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Get(tt.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenericList.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenericList_Set(t *testing.T) {
	list := GenericList{1, "2", "Three", GenericList{1, 2}, "last"}
	tests := []struct {
		name    string
		list    GenericList
		index   int
		value   interface{}
		want    IGenericList
		wantErr bool
	}{
		{"First", list, 0, 2, GenericList{2, "2", "Three", GenericList{1, 2}, "last"}, false},
		{"Last", list, list.Len() - 1, "changed", GenericList{1, "2", "Three", GenericList{1, 2}, "changed"}, false},
		{"Negative", list, -1, "negative", nil, true},
		{"Out of bound", list, 6, "Hello", GenericList{1, "2", "Three", GenericList{1, 2}, "last", nil, "Hello"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.list.Clone()
			got, err := l.Set(tt.index, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenericList.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenericList.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}
