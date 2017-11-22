package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

// AddExtraFuncs adds yaml, json and hcl converter functions to the current context
func AddExtraFuncs(template *template.Template) {
	runTemplate := func(source string, context ...interface{}) (string, error) {
		var out bytes.Buffer
		var err error

		// We first try to find a template named <source>
		t := template.Lookup(source)
		if t == nil {
			// This is not a template, so we try to load file named <source>
			if !strings.Contains(source, "\n") {
				tryFile := source
				if !path.IsAbs(tryFile) {
					// We try to load the file from the relative path
					tryFile = path.Join(path.Dir(RunningTemplate.ParseName), tryFile)
				}
				content, err := ioutil.ReadFile(tryFile)
				if err != nil {
					if _, ok := err.(*os.PathError); err != nil && !ok {
						return "", err
					}
				} else {
					source = string(content)
				}
			}
			if t == nil {
				// There is no file named <source>, so we consider that <source> is the content
				t, err = template.Parse(source)
				if err != nil {
					return "", err
				}
			}
		}

		// We execute the resulting template
		if err := t.Execute(&out, getContext(context...)); err != nil {
			return "", err
		}

		return out.String(), nil
	}

	// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
	converter := func(source string, convFunc func([]byte, interface{}) error, context ...interface{}) (interface{}, error) {
		var out interface{}
		content, err := runTemplate(source, context...)
		if err != nil {
			return out, err
		}

		err = convFunc([]byte(content), &out)
		if err != nil {
			for i, line := range strings.Split(content, "\n") {
				fmt.Printf("%4d %s\n", i+1, line)
			}
			return out, err
		}
		return out, err
	}

	// converts the supplied string containing yaml/json to go map
	yamlConverter := func(str string, context ...interface{}) (interface{}, error) {
		result, err := converter(str, yaml.Unmarshal, context...)
		if err != nil {
			return nil, err
		}
		return interface2string(result), nil
	}

	// Converts the supplied string containing terraform/hcl to go map
	hclConverter := func(str string, context ...interface{}) (interface{}, error) {
		out, err := converter(str, hcl.Unmarshal, context...)
		if err != nil {
			return nil, err
		}
		return FlattenHCL(out.(map[string]interface{})), nil
	}

	*template = *template.Funcs(map[string]interface{}{
		"yaml": yamlConverter,
		"json": yamlConverter,
		"hcl":  hclConverter,
		"toYaml": func(v interface{}) (string, error) {
			if result, err := yaml.Marshal(v); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		},
	})

	*template = *template.Funcs(map[string]interface{}{
		"bool": func(str string) (bool, error) {
			return strconv.ParseBool(str)
		},
		"get": func(str string, dict map[string]interface{}) interface{} {
			return dict[str]
		},
		"concat": func(objects ...interface{}) string {
			var result string
			for _, object := range objects {
				result += fmt.Sprint(object)
			}
			return result
		},
		"lorem":      lorem,
		"formatList": formatList,
		"mergeList":  mergeLists,
		"exec":       execFunc,
		"glob":       globFunc,
	})
}

var toStrings = sprig.GenericFuncMap()["toStrings"].(func(interface{}) []string)

func formatList(format string, v interface{}) []string {
	source := toStrings(v)
	list := make([]string, 0, len(source))
	for _, val := range source {
		list = append(list, fmt.Sprintf(format, val))
	}
	return list
}

func mergeLists(lists ...[]interface{}) []interface{} {
	switch len(lists) {
	case 0:
		return nil
	case 1:
		return lists[0]
	}
	result := make([]interface{}, 0)
	for _, list := range lists {
		result = append(result, list...)
	}
	return result
}

func execFunc(command string, args ...interface{}) (string, error) {
	if _, err := exec.LookPath(command); err != nil && strings.Contains(command, " ") {
		split := strings.Split(command, " ")
		command = split[0]
		newArgs := make([]interface{}, len(split)-1)
		for i := range split[1:] {
			newArgs[i] = split[i+1]
		}
		args = append(newArgs, args...)
	}
	cmd := exec.Command(command, globFunc(args...)...)
	cmd.Stdin = os.Stdin
	if RunningTemplate != nil {
		cmd.Dir = path.Dir(RunningTemplate.ParseName)
	}
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func globFunc(args ...interface{}) (result []string) {
	for _, arg := range toStrings(args) {
		if strings.ContainsAny(arg, "*?[]") {
			if expanded, _ := filepath.Glob(arg); expanded != nil {
				result = append(result, expanded...)
				continue
			}
		}
		result = append(result, arg)
	}
	return
}
