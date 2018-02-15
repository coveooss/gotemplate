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

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/utils"
)

// Add additional functions to the go template context
func (t *Template) addFuncs() {
	// Add functions form Sprig library https://github.com/Masterminds/sprig
	t.Funcs(sprig.GenericFuncMap())

	t.addMathFuncs()

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
		"undef":      utils.IfUndef,
		"iif":        utils.IIf,
		"ifUndef":    utils.IfUndef,
		"slice":      slice,
		"id":         id,
		"current":    func() string { return t.folder },

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
		"toQuotedJson": func(v interface{}) (string, error) {
			output, err := json.Marshal(v)
			result := fmt.Sprintf("%q", output)
			return result[1 : len(result)-1], err
		},
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
		if result, err = t.converter(converter, content, true, context...); err == nil {
			result = utils.MapKeyInterface2string(result)
		}
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
	if result, err = t.templateConverter(hcl.Unmarshal, str, context...); err == nil && result != nil {
		result = hcl.Flatten(result.(map[string]interface{}))
	}
	return
}

// Converts the supplied string containing yaml, json or terraform/hcl to go map
func (t Template) dataConverter(source string, context ...interface{}) (result interface{}, err error) {
	var content string
	if content, _, err = t.runTemplate(source, context...); err == nil {
		var errs errors.Array
		if result, err = t.converter(hcl.Unmarshal, content, true, context...); err == nil {
			result = hcl.Flatten(utils.MapKeyInterface2string(result).(map[string]interface{}))
		} else {
			errs = append(errs, err)
			if result, err = t.converter(utils.YamlUnmarshal, content, false, context...); err == nil {
				result = utils.MapKeyInterface2string(result)
				err = nil
			} else {
				// If there is still an error, we return both the hcl and yaml error
				err = append(errs, err)
			}
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
		"and", "call", "html", "index", "js", "len", "not", "or", "print", "printf", "println", "urlquery",
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

func slice(list interface{}, args ...interface{}) (interface{}, error) {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array, reflect.String:
		l2 := reflect.ValueOf(list)
		l := l2.Len()
		if l == 0 {
			return nil, nil
		}

		begin := 0
		end := l
		reverse := false

		switch len(args) {
		case 2:
			end = toInt(args[1])
			if end < 0 {
				end = l + end + 1
			}
			fallthrough
		case 1:
			begin = toInt(args[0])
			if begin < 0 {
				begin = l + begin + 1
			}
		}
		if end < begin {
			end, begin = begin, end
			reverse = true
		}
		end = int(min(end, l).(int64))
		begin = int(max(begin, 0).(int64))

		if tp == reflect.String {
			result := list.(string)[begin:end]
			if reverse {
				return reverseString(result), nil
			}
			return result, nil
		}
		if begin > l {
			return []interface{}{}, nil
		}
		nl := make([]interface{}, end-begin)
		for i := range nl {
			nl[i] = l2.Index(i + begin).Interface()
		}
		if reverse {
			return reverseArray(nl), nil
		}

		return nl, nil
	default:
		return nil, fmt.Errorf("Cannot apply slice on type %s", tp)
	}
}

func getSingleMapElement(m interface{}) (key, value interface{}, err error) {
	err = fmt.Errorf("Argument must be a map with a single key")
	if m == nil {
		return
	}
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	switch t.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) != 1 {
			return
		}
		return keys[0], v.MapIndex(keys[0]).Interface(), nil
	case reflect.Slice:
		l := v.Len()
		keys := make([]interface{}, l)
		values := make([]interface{}, l)
		for i := range keys {
			if keys[i], values[i], err = getSingleMapElement(v.Index(i).Interface()); err != nil {
				return
			}
		}

		results := make(map[string]interface{})
		for i := range keys {
			results[fmt.Sprint(keys[i])] = values[i]
		}
		return keys, results, nil

	default:
		return
	}
}

var reverseArray = sprig.GenericFuncMap()["reverse"].(func(v interface{}) []interface{})

// Reverse returns its argument string reversed rune-wise left to right.
func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
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
