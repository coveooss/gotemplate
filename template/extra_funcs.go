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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/utils"
)

// Add additional functions to the go template context
func (t *Template) addFuncs() {
	if t.include&Sprig != 0 {
		// Add functions from Sprig library https://github.com/Masterminds/sprig
		t.Funcs(sprig.GenericFuncMap())
	}

	if t.include&Math != 0 {
		t.addMathFuncs()
	}

	// Add utilities functions
	t.Funcs(map[string]interface{}{
		"string":     func(s interface{}) utils.String { return utils.String(fmt.Sprint(s)) },
		"concat":     utils.Concat,
		"formatList": utils.FormatList,
		"glob":       utils.GlobFunc,
		"joinLines":  utils.JoinLines,
		"mergeList":  utils.MergeLists,
		"pwd":        utils.Pwd,
		"splitLines": utils.SplitLines,
		"toYaml":     utils.ToYaml,
		"undef":      utils.IfUndef,
		"iif":        utils.IIf,
		"ifUndef":    utils.IfUndef,
		"slice":      slice,
		"safeIndex":  safeIndex,
		"extract":    extract,
		"id":         id,
		"center":     utils.CenterString,
		"current":    func() string { return t.folder },
		"lenc": func(s string) int {
			// Returns the actual length of a string
			return utf8.RuneCountInString(s)
		},
		"char": func(c interface{}) (r interface{}, err error) {
			defer func() { err = trapError(err, recover()) }()
			return process(c, func(a interface{}) interface{} {
				return string(toInt(a))
			})
		},
		"toHcl": func(v interface{}) (string, error) {
			output, err := hcl.Marshal(v)
			return string(output), err
		},
		"toPrettyHcl": func(v interface{}) (string, error) {
			output, err := hcl.MarshalIndent(v, "", "  ")
			return string(output), err
		},
		"toQuotedHcl": func(v interface{}) (string, error) {
			output, err := hcl.Marshal(v)
			result := fmt.Sprintf("%q", output)
			return result[1 : len(result)-1], err
		},
		"toTFVars": func(v interface{}) (string, error) {
			output, err := hcl.MarshalTFVars(v)
			return string(output), err
		},
		"toPrettyTFVars": func(v interface{}) (string, error) {
			output, err := hcl.MarshalTFVarsIndent(v, "", "  ")
			return string(output), err
		},
		"toQuotedTFVars": func(v interface{}) (string, error) {
			output, err := hcl.MarshalTFVars(v)
			result := fmt.Sprintf("%q", output)
			return result[1 : len(result)-1], err
		},
		"toJson": func(v interface{}) (string, error) {
			output, err := json.Marshal(v)
			return string(output), err
		},
		"toPrettyJson": func(v interface{}) (string, error) {
			output, err := json.MarshalIndent(v, "", "  ")
			return string(output), err
		},
		"toQuotedJson": func(v interface{}) (string, error) {
			output, err := json.Marshal(v)
			result := fmt.Sprintf("%q", output)
			return result[1 : len(result)-1], err
		},
		"toBash": utils.ToBash,
		// "toFile": func(file string, v interface{}) (string, error) {
		// },
		"bool": func(str string) (bool, error) {
			return strconv.ParseBool(str)
		},
		"get": func(arg1, arg2 interface{}) (result interface{}, err error) {
			// In pipe execution, the map is often the last parameter, but we also support to
			// put the map as the first parameter. So all following forms are supported:
			//    get map key
			//    get key map
			//    map | get key
			//    key | get map

			defer func() {
				if e := recover(); e != nil {
					err = fmt.Errorf("Cannot retrieve key from undefined map: %v", e)
				}
			}()

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
			return dict[key], nil
		},
		"set": func(arg1, arg2, arg3 interface{}) (result string, err error) {
			// In pipe execution, the map is often the last parameter, but we also support to
			// put the map as the first parameter. So all following forms are supported:
			//    set map key value
			//    set key value map
			//    map | set key value
			//    value | set map key
			defer func() {
				if e := recover(); e != nil {
					err = fmt.Errorf("Cannot set key from undefined map: %v", e)
				}
			}()

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
			return "", nil
		},
		"key": func(v interface{}) (interface{}, error) {
			key, _, err := getSingleMapElement(v)
			return key, err
		},
		"content": func(v interface{}) (interface{}, error) {
			_, value, err := getSingleMapElement(v)
			return value, err
		},
		"merge": utils.MergeMaps,
		"lorem": func(funcName interface{}, params ...int) (result string, err error) {
			kind, err := utils.GetLoremKind(fmt.Sprint(funcName))
			if err == nil {
				result, err = utils.Lorem(kind, params...)
			}
			return
		},
		"color": utils.SprintColor,
		"pick":  pick,
		"pickv": pickv,
		"omit":  omit,
	})

	// Add template related functions
	t.Funcs(map[string]interface{}{
		"functions":     t.getFunctions,
		"substitute":    t.substitute,
		"templateNames": t.getTemplateNames,
		"templates":     t.Templates,
		"alias": func(name, function string, source interface{}, args ...interface{}) (string, error) {
			return t.addAlias(name, function, source, false, false, args...)
		},
		"local_alias": func(name, function string, source interface{}, args ...interface{}) (string, error) {
			return t.addAlias(name, function, source, true, false, args...)
		},
		"func": func(name, function string, source, def, argNames interface{}) (string, error) {
			return t.addAlias(name, function, source, true, true, def, argNames)
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
		"ellipsis": t.ellipsis,
	})
}

// Define alias to an existing command
func (t *Template) addAlias(name, function string, source interface{}, local, context bool, defaultArgs ...interface{}) (result string, err error) {
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

	if !context {
		(*t.aliases)[name] = func(args ...interface{}) (result interface{}, err error) {
			return f(utils.Interface2string(source), append(defaultArgs, args...)...)
		}
		return
	}

	init := make(map[string]interface{})
	switch value := defaultArgs[0].(type) {
	case map[string]interface{}:
		init = value
	default:
		if err = utils.ConvertData(fmt.Sprint(value), &init); err != nil {
			return
		}
	}

	var argNames []string
	switch value := defaultArgs[1].(type) {
	case []string:
		argNames = value
	case []interface{}:
		argNames = toStrings(value)
	default:
		if err = utils.ConvertData(fmt.Sprint(value), &argNames); err != nil {
			return
		}
	}

	(*t.aliases)[name] = func(args ...interface{}) (result interface{}, err error) {
		context := make(map[string]interface{})
		parentContext, isMap := t.context.(map[string]interface{})
		if !isMap {
			context["DEFAULT"] = t.context
		}
		switch len(args) {
		case 1:
			if arg1, isMap := args[0].(map[string]interface{}); isMap {
				utils.MergeMaps(context, arg1, init, parentContext)
				break
			}
			if utils.ConvertData(fmt.Sprint(args[0]), &context) == nil {
				utils.MergeMaps(context, init, parentContext)
				break
			}
			fallthrough
		default:
			utils.MergeMaps(context, init, t.context.(map[string]interface{}))
			for i := range args {
				if i >= len(argNames) {
					context["ARGS"] = args[i:]
					break
				}
				context[argNames[i]] = args[i]
			}
		}
		return f(utils.Interface2string(source), context)
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

	var cmd *exec.Cmd
	if filename != "" {
		cmd, err = utils.GetCommandFromFile(filename, args...)
	} else {
		var tempFile string
		cmd, tempFile, err = utils.GetCommandFromString(command, args...)
		if tempFile != "" {
			defer func() { os.Remove(tempFile) }()
		}
	}

	if cmd == nil {
		return
	}

	var stdout, stderr bytes.Buffer
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
		inline, e := t.New("inline").Parse(source)
		if e != nil {
			err = e
			return
		}
		internalTemplate = inline
	}

	// We execute the resulting template
	if err = internalTemplate.Execute(&out, hcl.SingleContext(context...)); err != nil {
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

// This function is used to call a function that requires its last argument to be expanded ...
func (t Template) ellipsis(function string, args ...interface{}) (interface{}, error) {
	last := len(args) - 1
	if last < 0 || reflect.TypeOf(args[last]).Kind() != reflect.Slice {
		return nil, fmt.Errorf("The last argument must be a slice")
	}

	lastArg := reflect.ValueOf(args[last])
	argsStr := make([]string, 0, last+lastArg.Len())
	context := make(map[string]interface{})

	convertArg := func(arg interface{}) {
		argName := fmt.Sprintf("ARG%d", len(argsStr)+1)
		argsStr = append(argsStr, fmt.Sprintf(".%s", argName))
		context[argName] = arg
	}

	for i := range args[:last] {
		convertArg(args[i])
	}

	for i := 0; i < lastArg.Len(); i++ {
		convertArg(lastArg.Index(i).Interface())
	}

	template := fmt.Sprintf("%s %s %s %s", t.delimiters[0], function, strings.Join(argsStr, " "), t.delimiters[1])
	result, _, err := t.runTemplate(template, context)
	return result, err
}

type dataConverter func([]byte, interface{}) error

// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
func (t Template) converter(converter dataConverter, content string, sourceWithError bool, context ...interface{}) (result interface{}, err error) {
	if err = converter([]byte(content), &result); err != nil && sourceWithError {
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
	var content string
	if content, _, err = t.runTemplate(str, context...); err == nil {
		result, err = t.converter(converter, content, true, context...)
	}
	return
}

// converts the supplied string containing yaml to go map
func (t Template) yamlConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(utils.YamlUnmarshal, str, context...)
}

// converts the supplied string containing json to go map
func (t Template) jsonConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(json.Unmarshal, str, context...)
}

// Converts the supplied string containing terraform/hcl to go map
func (t Template) hclConverter(str string, context ...interface{}) (result interface{}, err error) {
	return t.templateConverter(hcl.Unmarshal, str, context...)
}

// Converts the supplied string containing yaml, json or terraform/hcl to go map
func (t Template) dataConverter(str string, context ...interface{}) (result interface{}, err error) {
	converter := func(bs []byte, out interface{}) (err error) {
		return utils.ConvertData(string(bs), out)
	}
	return t.templateConverter(converter, str, context...)
}

// Apply all regular expressions replacements to the supplied string
func (t Template) substitute(content string) string {
	return utils.Substitute(content, t.substitutes...)
}

// List the available functions in the template
func (t Template) getFunctions() []string {
	functions := []string{
		"and", "call", "html", "index", "js", "len", "not", "or", "print", "printf", "println", "urlquery",
		"eq", "ge", "gt", "le", "lt", "ne",
	}

	for name := range t.functions {
		functions = append(functions, name)
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

func id(id string, replace ...interface{}) string {
	// By default, replacement char for invalid chars would be _
	// but it is possible to specify an alternative string to act as the replacement
	replacement := fmt.Sprint(replace...)
	if replacement == "" {
		replacement = "_"
	}

	dup := duplicateUnderscore
	if replacement != "_" {
		// If the replacement string is not the default one, we generate a special substituter to remove duplicates
		// taking into account regex special chars such as +, ?, etc.
		dup = regexp.MustCompile(fmt.Sprintf(`(?:%s)+`, regexSpecial.ReplaceAllString(replacement, `\$0`)))
	}

	return dup.ReplaceAllString(validChars.ReplaceAllString(id, replacement), replacement)
}

var validChars = regexp.MustCompile(`[^\p{L}\d_]`)
var duplicateUnderscore = regexp.MustCompile(`_+`)
var regexSpecial = regexp.MustCompile(`[\+\.\?\(\)\[\]\{\}\\]`)
