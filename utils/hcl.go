package utils

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

// SingleContext converts array of 1 to single object otherwise, let the context unchanged
func SingleContext(context ...interface{}) interface{} {
	if len(context) == 1 {
		return context[0]
	}
	return context
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

// LoadHCL loads hcl file into variable
func LoadHCL(filename string) (result map[string]interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		if err = hcl.Unmarshal(content, &result); err == nil {
			result = FlattenHCL(result)
		}
	}
	return
}
