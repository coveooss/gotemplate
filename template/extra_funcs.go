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
		"concat":     utils.Concat,
		"formatList": utils.FormatList,
		"glob":       utils.GlobFunc,
		"joinLines":  utils.JoinLines,
		"mergeList":  utils.MergeLists,
		"pwd":        utils.Pwd,
		"splitLines": utils.SplitLines,
		"toYaml":     utils.ToYaml,
		"current":    func() string { return t.folder },

		"bool": func(str string) (bool, error) {
			return strconv.ParseBool(str)
		},
		"get": func(key string, dict map[string]interface{}) interface{} {
			return dict[key]
		},
		"set": func(key string, value interface{}, dict map[string]interface{}) string {
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
			return t.runTemplate(utils.Interface2string(source), context...)
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
	// We check if the supplied command is a template
	if command, err = t.runTemplate(command, args...); err != nil {
		return
	}

	// We expland the arguments
	cmdArgs := utils.GlobFunc(args...)

	executer, delegate, command := utils.ScriptParts(strings.TrimSpace(command))
	if executer != "" {
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
		command = executer
		cmdArgs = append([]string{temp.Name()}, cmdArgs...)
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
	cmd := exec.Command(command, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Dir = t.folder

	if result, err = cmd.CombinedOutput(); err == nil {
		result = string(result.([]byte))
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

func (t Template) runTemplate(source string, context ...interface{}) (string, error) {
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
			content, err := ioutil.ReadFile(tryFile)
			if err != nil {
				if _, ok := err.(*os.PathError); err != nil && !ok {
					return "", err
				}
			} else {
				source = string(content)
			}
		}
		// There is no file named <source>, so we consider that <source> is the content
		inline, err := t.New("inline").Parse(source)
		if err != nil {
			return "", err
		}
		internalTemplate = inline
	}

	// We execute the resulting t
	if err := internalTemplate.Execute(&out, utils.SingleContext(context...)); err != nil {
		return "", err
	}

	return out.String(), nil
}
func (t Template) runTemplateItf(source string, context ...interface{}) (interface{}, error) {
	return t.runTemplate(source, context...)
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
	if content, err := t.runTemplate(str, context...); err == nil {
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
	if content, err := t.runTemplate(source, context...); err == nil {
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
