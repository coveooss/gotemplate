package utils

import (
	"reflect"
	"testing"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/collections/implementation"
)

func TestMergeDictionaries(t *testing.T) {
	collections.DictionaryHelper = implementation.DictionaryHelper
	collections.ListHelper = implementation.GenericListHelper
	map1 := map[string]interface{}{
		"int":        1000,
		"Add1Int":    1,
		"Add1String": "string",
	}
	map2 := map[string]interface{}{
		"int":        2000,
		"Add2Int":    2,
		"Add2String": "string",
		"map": map[string]interface{}{
			"sub1":   2,
			"newVal": "NewValue",
		},
	}
	tests := []struct {
		name    string
		args    []map[string]interface{}
		want    map[string]interface{}
		wantErr bool
	}{
		{"Empty", nil, map[string]interface{}{}, false},
		{"Add 1", []map[string]interface{}{map1}, map1, false},
		{"Add empty to 1", []map[string]interface{}{map1, map[string]interface{}{}}, map1, false},
		{"Add nil to 1", []map[string]interface{}{map1, nil}, map1, false},
		{"Add 2 to 1", []map[string]interface{}{map1, map2}, map[string]interface{}{
			"int":        1000,
			"Add1Int":    1,
			"Add2Int":    2,
			"Add1String": "string",
			"Add2String": "string",
			"map": map[string]interface{}{
				"sub1":   2,
				"newVal": "NewValue",
			},
		}, false},
		{"Add 1 to 2", []map[string]interface{}{map2, map1}, map[string]interface{}{
			"int":        2000,
			"Add1Int":    1,
			"Add2Int":    2,
			"Add1String": "string",
			"Add2String": "string",
			"map": map[string]interface{}{
				"sub1":   2,
				"newVal": "NewValue",
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeDictionaries(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeDictionaries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeDictionaries():\n got %v\nwant %v", got, tt.want)
			}
		})
	}
}
