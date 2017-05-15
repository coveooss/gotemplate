package main

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

func getContext(context ...interface{}) interface{} {
	switch len(context) {
	case 0:
		return nil
	case 1:
		return context[0]
	default:
		return context
	}
}

// FlattenHCL - Convert array of map to single map if there is only one element in the array
// By default, the hcl.Unmarshal returns array of map even if there is only a single map in the definition
func FlattenHCL(source map[string]interface{}) map[string]interface{} {
	for key, value := range source {
		switch value := value.(type) {
		case []map[string]interface{}:
			switch len(value) {
			case 1:
				source[key] = FlattenHCL(value[0])
			default:
				for i, subMap := range value {
					value[i] = FlattenHCL(subMap)
				}
			}
		}
	}
	return source
}

// Load HCL file into variable
func loadHCL(filename string) (result map[string]interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		if err = hcl.Unmarshal(content, &result); err == nil {
			result = FlattenHCL(result)
		}
	}
	return
}
