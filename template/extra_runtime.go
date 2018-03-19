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
	runtimeFunc = "Runtime"
)

var runtimeFuncsArgs = map[string][]string{
	"alias":      {"name", "function", "source"},
	"ellipsis":   {"function"},
	"exec":       {"command"},
	"exit":       {"exitValue"},
	"func":       {"name", "function", "source", "config"},
	"function":   {"name"},
	"include":    {"source", "context"},
	"localAlias": {"name", "function", "source"},
	"run":        {"command"},
	"substitute": {"content"},
}

var runtimeFuncsAliases = map[string][]string{
	"exec": {"execute"},
}

var runtimeFuncsHelp = map[string]string{
	"alias":         "Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the caller.",
	"aliases":       "Returns the list of all functions that are simply an alias of another function.",
	"allFunctions":  "Returns the list of all available functions.",
	"current":       "Returns the current folder (like pwd, but returns the folder of the currently running folder).",
	"ellipsis":      "Returns the result of the function by expanding its last argument that must be an array into values. It's like calling function(arg1, arg2, otherArgs...).",
	"exec":          "Returns the result of the shell command as structured data (as string if no other conversion is possible).",
	"exit":          "Exits the current program execution.",
	"func":          "Defines a function with the current context using the function (exec, run, include, template). Executed in the context of the caller.",
	"function":      "Returns the information relative to a specific function.",
	"functions":     "Returns the list of all available functions (excluding aliases).",
	"include":       "Returns the result of the named template rendering (like template but it is possible to capture the output).",
	"localAlias":    "Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the function it maps to.",
	"run":           "Returns the result of the shell command as string.",
	"substitute":    "Applies the supplied regex substitute specified on the command line on the supplied string (see --substitute).",
	"templateNames": "Returns the list of available templates names.",
	"templates":     "Returns the list of available templates.",
}

func (t *Template) addRuntimeFuncs() {
	var funcs = map[string]interface{}{
		"alias":         t.alias,
		"aliases":       t.getAliases,
		"allFunctions":  t.getAllFunctions,
		"current":       t.current,
		"ellipsis":      t.ellipsis,
		"exec":          t.execCommand,
		"exit":          exit,
		"func":          t.defineFunc,
		"function":      t.getFunction,
		"functions":     t.getFunctions,
		"include":       t.include,
		"localAlias":    t.localAlias,
		"run":           t.runCommand,
		"substitute":    t.substitute,
		"templates":     t.Templates,
		"templateNames": t.getTemplateNames,
	}

	t.AddFunctions(funcs, runtimeFunc, funcOptions{
		funcHelp:    runtimeFuncsHelp,
		funcArgs:    runtimeFuncsArgs,
		funcAliases: runtimeFuncsAliases,
	})
}

func exit(exitValue int) int       { os.Exit(exitValue); return exitValue }
func (t Template) current() string { return t.folder }

func (t *Template) alias(name, function string, source interface{}, args ...interface{}) (string, error) {
	return t.addAlias(name, function, source, false, false, args...)
}

func (t *Template) localAlias(name, function string, source interface{}, args ...interface{}) (string, error) {
	return t.addAlias(name, function, source, true, false, args...)
}

func (t *Template) defineFunc(name, function string, source, config interface{}) (string, error) {
	return t.addAlias(name, function, source, true, true, config)
}

func (t *Template) execCommand(command interface{}, args ...interface{}) (interface{}, error) {
	return t.exec(utils.Interface2string(command), args...)
}

func (t *Template) runCommand(command interface{}, args ...interface{}) (interface{}, error) {
	return t.run(utils.Interface2string(command), args...)
}

func (t *Template) include(source interface{}, context ...interface{}) (interface{}, error) {
	content, _, err := t.runTemplate(utils.Interface2string(source), context...)
	if source == content {
		return nil, fmt.Errorf("Unable to find a template or a file named %s", source)
	}
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
		t.aliases[name] = FuncInfo{
			function: func(args ...interface{}) (result interface{}, err error) {
				return f(utils.Interface2string(source), append(defaultArgs, args...)...)
			},
			group: "User defined aliases",
		}
		return
	}

	var config map[string]interface{}
	var ok bool

	switch len(defaultArgs) {
	case 0:
		config = make(map[string]interface{})
	case 1:
		if config, ok = defaultArgs[0].(map[string]interface{}); !ok {
			if err = utils.ConvertData(fmt.Sprint(defaultArgs[0]), &config); err != nil {
				err = fmt.Errorf("Function configuration must be a valid map definition: %[1]T %[1]v", defaultArgs[0])
				return
			}
		}
	default:
		return "", fmt.Errorf("Too many parameter supplied")
	}

	for key, val := range config {
		switch strings.ToLower(key) {
		case "d", "desc", "description":
			config["description"] = val
		case "g", "group":
			config["group"] = val
		case "a", "args", "arguments":
			switch val := val.(type) {
			case []string:
				config["args"] = val
			case []interface{}:
				config["args"] = toStrings(val)
			default:
				err = fmt.Errorf("%[1]s must be a list of strings: %[2]T %[2]v", key, val)
				return
			}
		case "aliases":
			switch val := val.(type) {
			case []string:
				config["aliases"] = val
			case []interface{}:
				config["aliases"] = toStrings(val)
			default:
				err = fmt.Errorf("%[1]s must be a list of strings: %[2]T %[2]v", key, val)
				return
			}
		case "def", "default", "defaults":
			if _, ok = val.(map[string]interface{}); !ok {
				err = fmt.Errorf("%s must be a map", key)
				return
			}
			config["def"] = val
		default:
			return "", fmt.Errorf("Unknown configuration %s", key)
		}
	}

	fi := FuncInfo{
		name:        name,
		group:       defval(config["group"], "User defined functions").(string),
		description: defval(config["description"], "").(string),
		arguments:   defval(config["args"], []string{}).([]string),
		aliases:     defval(config["aliases"], []string{}).([]string),
	}

	defaultValues := defval(config["def"], make(map[string]interface{})).(map[string]interface{})

	fi.in = fmt.Sprintf("%s", strings.Join(fi.arguments, ", "))
	for i := range fi.arguments {
		// We only keep the arg name and get rid of any supplemental information (likely type)
		fi.arguments[i] = strings.Fields(fi.arguments[i])[0]
	}

	fi.function = func(args ...interface{}) (result interface{}, err error) {
		context := make(map[string]interface{})
		parentContext, isMap := t.context.(map[string]interface{})
		if !isMap {
			context["DEFAULT"] = t.context
		}

		switch len(args) {
		case 1:
			if len(fi.arguments) != 1 {
				if arg1, isMap := args[0].(map[string]interface{}); isMap {
					utils.MergeMaps(context, arg1, defaultValues, parentContext)
					break
				}
				if utils.ConvertData(fmt.Sprint(args[0]), &context) == nil {
					utils.MergeMaps(context, defaultValues, parentContext)
					break
				}
			}
			fallthrough
		default:
			utils.MergeMaps(context, defaultValues, t.context.(map[string]interface{}))
			for i := range args {
				if i >= len(fi.arguments) {
					context["ARGS"] = args[i:]
					break
				}
				context[fi.arguments[i]] = args[i]
			}
		}
		return f(utils.Interface2string(source), context)
	}

	t.aliases[name] = fi
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
	log.Notice("Launching", cmd.Args, "in", cmd.Dir)

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

	if len(context) == 0 {
		context = []interface{}{t.context}
	}
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
				source = string(t.applyRazor(fileContent))
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
	if last >= 0 && args[last] == nil {
		args[last] = []interface{}{}
	} else if last < 0 || reflect.TypeOf(args[last]).Kind() != reflect.Slice {
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
