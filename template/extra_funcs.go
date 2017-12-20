package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/utils"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

// Add additional functions to the go template context
func (t *Template) addFuncs() {
	// Add functions form Sprig library https://github.com/Masterminds/sprig
	t.Funcs(sprig.GenericFuncMap())

	// Add utilities functions
	t.Funcs(map[string]interface{}{
		"concat":      utils.Concat,
		"formatList":  utils.FormatList,
		"glob":        utils.GlobFunc,
		"joinLines":   utils.JoinLines,
		"mergeList":   utils.MergeLists,
		"pwd":         utils.Pwd,
		"splitLines":  utils.SplitLines,
		"toYaml":      utils.ToYaml,
		"current":     func() string { return t.folder },
		"toHcl":       func(v interface{}) string { return string(utils.ToHCL(v)) },
		"toPrettyHcl": func(v interface{}) string { return string(utils.ToPrettyHCL(v)) },
		"toQuotedJson": func(v interface{}) string {
			output, _ := json.Marshal(v)
			result := fmt.Sprintf("%q", output)
			return result[1 : len(result)-1]
		},
		"bool": func(str string) (bool, error) {
			return strconv.ParseBool(str)
		},
		"get": func(arg1, arg2 interface{}) interface{} {
			// In pipe execution, the map is often the last parameter, but we also support to
			// put the map as the first parameter. So all following forms are supported:
			//    get map key
			//    get key map
			//    map | get key
			//    key | get map
			var (
				dict map[string]interface{}
				key  string
			)
			if reflect.TypeOf(arg1).Kind() == reflect.Map {
				dict = arg1.(map[string]interface{})
				key = arg2.(string)
			} else {
				key = arg1.(string)
				dict = arg2.(map[string]interface{})
			}
			return dict[key]
		},
		"set": func(arg1, arg2, arg3 interface{}) string {
			// In pipe execution, the map is often the last parameter, but we also support to
			// put the map as the first parameter. So all following forms are supported:
			//    set map key value
			//    set key value map
			//    map | set key value
			//    value | set map key
			var (
				dict  map[string]interface{}
				key   string
				value interface{}
			)
			if reflect.TypeOf(arg1).Kind() == reflect.Map {
				dict = arg1.(map[string]interface{})
				key = arg2.(string)
				value = arg3
			} else {
				key = arg1.(string)
				value = arg2
				dict = arg3.(map[string]interface{})
			}
			dict[key] = value
			return ""
		},
		"lorem": func(funcName string, params ...int) (result string, err error) {
			kind, err := utils.GetLoremKind(funcName)
			if err == nil {
				result, err = utils.Lorem(kind, params...)
			}
			return
		},
	})

	// Add template related functions
	t.Funcs(map[string]interface{}{
		"functions":     t.getFunctions,
		"substitute":    t.substitute,
		"templateNames": t.getTemplateNames,
		"templates":     t.Templates,
		"alias": func(name, function string, source interface{}, args ...interface{}) (result string, err error) {
			return t.addAlias(name, function, source, false, args...)
		},
		"local_alias": func(name, function string, source interface{}, args ...interface{}) (result string, err error) {
			return t.addAlias(name, function, source, true, args...)
		},
		"data": func(source interface{}, context ...interface{}) (interface{}, error) {
			return t.dataConverter(utils.Interface2string(source), context...)
		},
		"hcl": func(source interface{}, context ...interface{}) (interface{}, error) {
			return t.hclConverter(utils.Interface2string(source), context...)
		},
		"json": func(source interface{}, context ...interface{}) (interface{}, error) {
			return t.jsonConverter(utils.Interface2string(source), context...)
		},
		"yaml": func(source interface{}, context ...interface{}) (interface{}, error) {
			return t.yamlConverter(utils.Interface2string(source), context...)
		},
		"exec": func(command interface{}, args ...interface{}) (interface{}, error) {
			return t.exec(utils.Interface2string(command), args...)
		},
		"run": func(command interface{}, args ...interface{}) (interface{}, error) {
			return t.run(utils.Interface2string(command), args...)
		},
		"include": func(source interface{}, context ...interface{}) (interface{}, error) {
			content, _, err := t.runTemplate(utils.Interface2string(source), context...)
			return content, err
		},
	})
}

// Define alias to an existing command
func (t *Template) addAlias(name, function string, source interface{}, local bool, defaultArgs ...interface{}) (result string, err error) {
	for !local && t.parent != nil {
		// local specifies if the alias should be executed in the context of the template where it is
		// defined or in the context of the top parent
		t = t.parent
	}

	f := t.run

	switch function {
	case "run":
	case "exec":
		f = t.exec
	case "template", "include":
		f = t.runTemplateItf
	default:
		err = fmt.Errorf("%s unsupported for alias %s (only run, exec, template and include are supported)", function, name)
		return
	}

	(*t.aliases)[name] = func(args ...interface{}) (result interface{}, err error) {
		return f(utils.Interface2string(source), args...)
	}
	return
}

// Execute the command (command could be a file, a template or a script)
func (t *Template) run(command string, args ...interface{}) (result interface{}, err error) {
	var filename string

	// We check if the supplied command is a template
	if command, filename, err = t.runTemplate(command, args...); err != nil {
		return
	}

	// We expland the arguments
	cmdArgs := utils.GlobFunc(args...)

	executer, delegate, command := utils.ScriptParts(strings.TrimSpace(command))
	if executer != "" {
		if filename == "" {
			// The command is a shebang script, so we save the content as a temporary file
			var temp *os.File
			if temp, err = ioutil.TempFile(t.TempFolder, "exec_"); err != nil {
				return
			}
			defer func() { os.Remove(temp.Name()) }()

			if _, err = temp.WriteString(command); err != nil {
				return
			}
			temp.Close()
			filename = temp.Name()
		}
		command = executer
		cmdArgs = append([]string{filename}, cmdArgs...)
		if delegate != "" {
			cmdArgs = append([]string{delegate}, cmdArgs...)
		}
	} else if _, errPath := exec.LookPath(command); errPath != nil {
		if strings.Contains(command, " ") {
			// The command is a string that should be splitted up into several parts
			split := strings.Split(command, " ")
			command = split[0]
			cmdArgs = append(split[1:], cmdArgs...)
		} else {
			// The command does not exist
			return
		}
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = t.folder

	if err = cmd.Run(); err == nil {
		result = stdout.String()
	} else {
		err = fmt.Errorf("Error %v: %s", err, stderr.String())
	}
	return
}

// Execute the command (command could be a file, a template or a script) and convert its result as data if possible
func (t *Template) exec(command string, args ...interface{}) (result interface{}, err error) {
	if result, err = t.run(command, args...); err == nil {
		if result == nil {
			return
		}
		if data, errData := t.dataConverter(result.(string)); errData == nil {
			// The result of the command is a valid data structure
			if reflect.TypeOf(data).Kind() != reflect.String {
				result = data
			}
		}
	}
	return
}

func (t Template) runTemplate(source string, context ...interface{}) (resultContent, filename string, err error) {
	var out bytes.Buffer

	// We first try to find a template named <source>
	internalTemplate := t.Lookup(source)
	if internalTemplate == nil {
		// This is not a template, so we try to load file named <source>
		if !strings.Contains(source, "\n") {
			tryFile := source
			if !path.IsAbs(tryFile) {
				tryFile = path.Join(t.folder, tryFile)
			}
			if fileContent, e := ioutil.ReadFile(tryFile); e != nil {
				if _, ok := e.(*os.PathError); !ok {
					err = e
					return
				}
			} else {
				source = string(fileContent)
				filename = tryFile
			}
		}
		// There is no file named <source>, so we consider that <source> is the content
		if inline, e := t.New("inline").Parse(source); e != nil {
			err = e
			return
		} else {
			internalTemplate = inline
		}
	}

	// We execute the resulting template
	if err = internalTemplate.Execute(&out, utils.SingleContext(context...)); err != nil {
		return
	}

	resultContent = out.String()
	if resultContent != source {
		// If the content is different from the source, that is because the source contains
		// templating, In that case, we do not consider the original filename as unaltered source.
		filename = ""
	}
	return
}
func (t Template) runTemplateItf(source string, context ...interface{}) (interface{}, error) {
	content, _, err := t.runTemplate(source, context...)
	return content, err
}

type dataConverter func([]byte, interface{}) error

// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
func (t Template) converter(converter dataConverter, content string, context ...interface{}) (result interface{}, err error) {
	if err = converter([]byte(content), &result); err != nil {
		source := "\n"
		for i, line := range utils.SplitLines(content) {
			source += fmt.Sprintf("%4d %s\n", i+1, line)
		}
		err = fmt.Errorf("%s\n%v", source, err)
	}
	return
}

// Apply a converter to the result of the template execution of the supplied string
func (t Template) templateConverter(converter dataConverter, str string, context ...interface{}) (result interface{}, err error) {
	if content, _, err := t.runTemplate(str, context...); err == nil {
		if result, err = t.converter(converter, content, context...); err == nil {
			result = utils.MapKeyInterface2string(result)
		}
	}
	return
}

// converts the supplied string containing yaml to go map
func (t Template) yamlConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(yaml.Unmarshal, str, context...)
}

// converts the supplied string containing json to go map
func (t Template) jsonConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(json.Unmarshal, str, context...)
}

// Converts the supplied string containing terraform/hcl to go map
func (t Template) hclConverter(str string, context ...interface{}) (result interface{}, err error) {
	if result, err = t.templateConverter(hcl.Unmarshal, str, context...); err == nil && result != nil {
		result = utils.FlattenHCL(result.(map[string]interface{}))
	}
	return
}

// Converts the supplied string containing yaml, json or terraform/hcl to go map
func (t Template) dataConverter(source string, context ...interface{}) (result interface{}, err error) {
	if content, _, err := t.runTemplate(source, context...); err == nil {
		if result, err = t.converter(hcl.Unmarshal, content, context...); err == nil {
			result = utils.FlattenHCL(utils.MapKeyInterface2string(result).(map[string]interface{}))
		} else if result, err = t.converter(yaml.Unmarshal, content, context...); err == nil {
			result = utils.MapKeyInterface2string(result)
		}
	}
	return
}

// Apply all regular expressions replacements to the supplied string
func (t Template) substitute(content string) string {
	return utils.Substitute(content, t.substitutes...)
}

// List the available functions in the template
func (t Template) getFunctions() []string {
	functions := []string{
		"and", "call", "call", "html", "index", "js", "len", "not", "or", "print", "printf", "println", "urlquery",
		"eq", "ge", "gt", "le", "lt", "ne",
	}

	for _, k := range reflect.ValueOf(t).FieldByName("common").Elem().FieldByName("parseFuncs").MapKeys() {
		functions = append(functions, k.String())
	}
	sort.Strings(functions)
	return functions
}

// List the available template names
func (t Template) getTemplateNames() []string {
	templates := t.Templates()
	result := make([]string, len(templates))
	for i := range templates {
		result[i] = templates[i].Name()
	}
	sort.Strings(result)
	return result
}
