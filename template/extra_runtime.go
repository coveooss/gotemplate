package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"

	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/utils"
)

const (
	runtimeFunc = "Runtime functions"
)

var runtimeFuncs funcTableMap

func (t *Template) addRuntimeFuncs() {
	if runtimeFuncs == nil {
		runtimeFuncs = funcTableMap{
			"functions":     {t.getFunctions, runtimeFunc, nil, []string{}, ""},
			"substitute":    {t.substitute, runtimeFunc, nil, []string{}, ""},
			"templateNames": {t.Templates, runtimeFunc, nil, []string{}, ""},
			"ellipsis":      {t.ellipsis, runtimeFunc, nil, []string{}, ""},
			"alias":         {t.alias, runtimeFunc, nil, []string{}, ""},
			"localAlias":    {t.localAlias, runtimeFunc, nil, []string{}, ""},
			"func":          {t.defineFunc, runtimeFunc, nil, []string{}, ""},
			"exec":          {t.execCommand, runtimeFunc, nil, []string{}, ""},
			"run":           {t.runCommand, runtimeFunc, nil, []string{}, ""},
			"include":       {t.include, runtimeFunc, nil, []string{}, ""},
			"current":       {t.current, runtimeFunc, nil, []string{}, ""},
			"exit":          {exit, runtimeFunc, nil, []string{}, ""},
		}
	}

	t.AddFunctions(runtimeFuncs)
}

func exit(exitValue int) int       { os.Exit(exitValue); return exitValue }
func (t Template) current() string { return t.folder }

func (t *Template) alias(name, function string, source interface{}, args ...interface{}) (string, error) {
	return t.addAlias(name, function, source, false, false, args...)
}

func (t *Template) localAlias(name, function string, source interface{}, args ...interface{}) (string, error) {
	return t.addAlias(name, function, source, true, false, args...)
}

func (t *Template) defineFunc(name, function string, source, def, argNames interface{}) (string, error) {
	return t.addAlias(name, function, source, true, true, def, argNames)
}

func (t *Template) execCommand(command interface{}, args ...interface{}) (interface{}, error) {
	return t.exec(utils.Interface2string(command), args...)
}

func (t *Template) runCommand(command interface{}, args ...interface{}) (interface{}, error) {
	return t.run(utils.Interface2string(command), args...)
}

func (t *Template) include(source interface{}, context ...interface{}) (interface{}, error) {
	content, _, err := t.runTemplate(utils.Interface2string(source), context...)
	return content, err
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
