package main

import (
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
