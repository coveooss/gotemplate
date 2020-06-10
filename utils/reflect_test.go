package utils

import (
	"reflect"
	"testing"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/collections/implementation"
)

func TestMergeDictionaries(t *testing.T) {
	collections.SetDictionaryHelper(implementation.DictionaryHelper)
	collections.SetListHelper(implementation.GenericListHelper)
	map1 := map[string]interface{}{
		"int":         1000,
		"Add1Int":     1,
		"Add1String":  "string",
		"Add1Boolean": true,
		"Boolean":     true,
		"map": map[string]interface{}{
			"AddBoolean1": false,
			"Boolean":     true,
			"Boolean2":    false,
		},
	}
	map2 := map[string]interface{}{
		"int":         2000,
		"Add2Int":     2,
		"Add2String":  "string",
		"Add2Boolean": true,
		"Boolean":     false,
		"map": map[string]interface{}{
			"AddBoolean2": true,
			"Boolean":     false,
			"Boolean2":    true,
			"sub1":        2,
			"newVal":      "NewValue",
		},
	}
	map3 := map[string]interface{}{
		"Add3Int": 2,
		"Boolean": false,
		"map": map[string]interface{}{
			"AddBoolean3": true,
			"Boolean":     false,
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
		{"Add empty to 1", []map[string]interface{}{map1, {}}, map1, false},
		{"Add nil to 1", []map[string]interface{}{map1, nil}, map1, false},
		{"Add 2 to 1", []map[string]interface{}{map1, map2}, map[string]interface{}{
			"int":         1000,
			"Add1Int":     1,
			"Add2Int":     2,
			"Add1String":  "string",
			"Add2String":  "string",
			"Add1Boolean": true,
			"Add2Boolean": true,
			"Boolean":     true,
			"map": map[string]interface{}{
				"AddBoolean1": false,
				"AddBoolean2": true,
				"Boolean":     true,
				"Boolean2":    false,
				"sub1":        2,
				"newVal":      "NewValue",
			},
		}, false},
		{"Add 1 to 2", []map[string]interface{}{map2, map1}, map[string]interface{}{
			"int":         2000,
			"Add1Int":     1,
			"Add2Int":     2,
			"Add1String":  "string",
			"Add2String":  "string",
			"Add1Boolean": true,
			"Add2Boolean": true,
			"Boolean":     false,
			"map": map[string]interface{}{
				"AddBoolean1": false,
				"AddBoolean2": true,
				"Boolean":     false,
				"Boolean2":    true,
				"sub1":        2,
				"newVal":      "NewValue",
			},
		}, false},
		{"Add 1 to 2 to 3", []map[string]interface{}{map3, map2, map1}, map[string]interface{}{
			"int":         2000,
			"Boolean":     false,
			"Add1Boolean": true,
			"Add1Int":     1,
			"Add1String":  "string",
			"Add2Boolean": true,
			"Add2Int":     2,
			"Add2String":  "string",
			"Add3Int":     2,
			"map": map[string]interface{}{
				"AddBoolean1": false,
				"AddBoolean2": true,
				"AddBoolean3": true,
				"Boolean":     false,
				"Boolean2":    true,
				"sub1":        2,
				"newVal":      "NewValue",
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
