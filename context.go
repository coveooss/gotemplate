package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/types"
	"github.com/coveo/gotemplate/utils"
)

type iDict = types.IDictionary
type dict = types.Dictionary

func createContext(varsFiles []string, namedVars []string) (context dict) {
	context = make(dict)

	type fileDef struct {
		name    string
		value   interface{}
		unnamed bool
	}

	nameValuePairs := make([]fileDef, 0, len(varsFiles)+len(namedVars))
	for i := range varsFiles {
		nameValuePairs = append(nameValuePairs, fileDef{value: varsFiles[i]})
	}

	for i := range namedVars {
		data := make(dict)
		if err := utils.ConvertData(namedVars[i], &data); err != nil {
			var fd fileDef
			fd.name, fd.value = types.Split2(namedVars[i], "=")
			if fd.value == "" {
				fd = fileDef{"", fd.name, true}
			}
			nameValuePairs = append(nameValuePairs, fd)
			continue
		}
		if len(data) == 0 && strings.Contains(namedVars[i], "=") {
			// The hcl converter consider "value=" as an empty map instead of empty value in a map
			// we handle it
			name, value := types.Split2(namedVars[i], "=")
			data[name] = value
		}
		for key, value := range data {
			nameValuePairs = append(nameValuePairs, fileDef{key, value, false})
		}
	}

	var unnamed []interface{}
	for _, nv := range nameValuePairs {
		var loader func(string) (iDict, error)
		filename, _ := reflect.ValueOf(nv.value).Interface().(string)
		if filename != "" {
			loader = func(filename string) (iDict, error) {
				var content interface{}
				if err := utils.LoadData(filename, &content); err == nil {
					if nv.name == "" && !nv.unnamed {
						if content, err := types.AsDictionary(content); err == nil {
							return content, nil
						}
					}
					if nv.name == "" {
						nv.name = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
					}
					return dict{nv.name: content}, nil
				} else if _, isFileErr := err.(*os.PathError); !isFileErr {
					return nil, err
				}

				// Finally, we just try to convert the data with the converted value
				if err := utils.ConvertData(filename, &content); err != nil {
					content = nv.value
				}

				if nv.name == "" {
					unnamed = append(unnamed, content)
					return nil, nil
				}

				// If it does not work, we just set the value as is
				return dict{nv.name: content}, nil
			}

			if filename == "-" {
				loader = func(filename string) (result iDict, err error) {
					var content interface{}
					if err = utils.ConvertData(readStdin(), &content); err != nil {
						return nil, err
					}
					if nv.name == "" {
						if content, err := types.AsDictionary(content); err == nil {
							return content, nil
						}
					}
					if nv.name == "" {
						nv.name = "STDIN"
					}
					return dict{nv.name: content}, nil
				}
			}
		}

		if loader == nil {
			context[nv.name] = nv.value
			continue
		}
		content, err := loader(filename)
		if err != nil {
			errors.Raise("Error %v while loading vars file %s", nv.value, err)
		}
		for key, value := range content.AsMap() {
			context[key] = value
		}
	}

	if len(unnamed) > 0 {
		context["ARGS"] = unnamed
	}
	return
}
