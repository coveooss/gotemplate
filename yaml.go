package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func loadYaml(filename string) (result map[string]interface{}, err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		err = yaml.Unmarshal(content, &result)
	}
	return
}

func interface2string(source interface{}) interface{} {
	switch value := source.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(value))
		for key, val := range value {
			result[fmt.Sprintf("%v", key)] = interface2string(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(value), len(value))
		for i, val := range value {
			result[i] = interface2string(val)
		}
		return result
	}
	return source
}
