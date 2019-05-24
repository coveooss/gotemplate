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

	"github.com/coveo/gotemplate/v3/collections"
	"github.com/coveo/gotemplate/v3/hcl"
	"github.com/coveo/gotemplate/v3/utils"
	"github.com/fatih/color"
)

const (
	runtimeFunc = "Runtime"
)

var runtimeFuncsArgs = arguments{
	"alias":         {"name", "function", "source"},
	"assert":        {"test", "message", "arguments"},
	"assertWarning": {"test", "message", "arguments"},
	"categories":    {"functionsGroups"},
	"ellipsis":      {"function"},
	"exec":          {"command"},
	"exit":          {"exitValue"},
	"func":          {"name", "function", "source", "config"},
	"function":      {"name"},
	"include":       {"source", "context"},
	"localAlias":    {"name", "function", "source"},
	"run":           {"command"},
	"substitute":    {"content"},
}

var runtimeFuncsAliases = aliases{
	"assert":        {"assertion"},
	"assertWarning": {"assertw"},
	"exec":          {"execute"},
	"getAttributes": {"attr", "attributes"},
	"getMethods":    {"methods"},
	"getSignature":  {"sign", "signature"},
	"raise":         {"raiseError"},
}

var runtimeFuncsHelp = descriptions{
	"alias":         "Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the caller.",
	"aliases":       "Returns the list of all functions that are simply an alias of another function.",
	"allFunctions":  "Returns the list of all available functions.",
	"assert":        "Raises a formated error if the test condition is false.",
	"assertWarning": "Issues a formated warning if the test condition is false.",
	"categories": strings.TrimSpace(collections.UnIndent(`
		Returns all functions group by categories.

		The returned value has the following properties:
		    Name        string
		    Functions    []string
	`)),
	"current":  "Returns the current folder (like pwd, but returns the folder of the currently running folder).",
	"ellipsis": "Returns the result of the function by expanding its last argument that must be an array into values. It's like calling function(arg1, arg2, otherArgs...).",
	"exec":     "Returns the result of the shell command as structured data (as string if no other conversion is possible).",
	"exit":     "Exits the current program execution.",
	"func":     "Defines a function with the current context using the function (exec, run, include, template). Executed in the context of the caller.",
	"function": strings.TrimSpace(collections.UnIndent(`
		Returns the information relative to a specific function.

		The returned value has the following properties:
		    Name        string
		    Description string
		    Signature   string
		    Group       string
		    Aliases     []string
		    Arguments   string
		    Result      string
	`)),
	"functions":     "Returns the list of all available functions (excluding aliases).",
	"getAttributes": "List all attributes accessible from the supplied object.",
	"getMethods":    "List all methods signatures accessible from the supplied object.",
	"getSignature":  "List all attributes and methods signatures accessible from the supplied object.",
	"include":       "Returns the result of the named template rendering (like template but it is possible to capture the output).",
	"localAlias":    "Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the function it maps to.",
	"raise":         "Raise a formated error.",
	"run":           "Returns the result of the shell command as string.",
	"substitute":    "Applies the supplied regex substitute specified on the command line on the supplied string (see --substitute).",
	"templateNames": "Returns the list of available templates names.",
	"templates":     "Returns the list of available templates.",
}

func (t *Template) addRuntimeFuncs() {
	var funcs = dictionary{
		"alias":         t.alias,
		"aliases":       t.getAliases,
		"allFunctions":  t.getAllFunctions,
		"assert":        assertError,
		"assertWarning": assertWarning,
		"categories":    t.getCategories,
		"current":       t.current,
		"ellipsis":      t.ellipsis,
		"exec":          t.execCommand,
		"exit":          exit,
		"func":          t.defineFunc,
		"function":      t.getFunction,
		"functions":     t.getFunctions,
		"getAttributes": getAttributes,
		"getMethods":    getMethods,
		"getSignature":  getSignature,
		"include":       t.include,
		"localAlias":    t.localAlias,
		"raise":         raise,
		"run":           t.runCommand,
		"substitute":    t.substitute,
		"templateNames": t.getTemplateNames,
		"templates":     t.Templates,
	}

	t.AddFunctions(funcs, runtimeFunc, FuncOptions{
		FuncHelp:    runtimeFuncsHelp,
		FuncArgs:    runtimeFuncsArgs,
		FuncAliases: runtimeFuncsAliases,
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
	return t.exec(collections.Interface2string(command), args...)
}

func (t *Template) runCommand(command interface{}, args ...interface{}) (interface{}, error) {
	return t.run(collections.Interface2string(command), args...)
}

func (t *Template) include(source interface{}, context ...interface{}) (interface{}, error) {
	content, _, err := t.runTemplate(collections.Interface2string(source), context...)
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
				return f(collections.Interface2string(source), append(defaultArgs, args...)...)
			},
			group: "User defined aliases",
		}
		return
	}

	var config iDictionary

	switch len(defaultArgs) {
	case 0:
		config = collections.CreateDictionary()
	case 1:
		if defaultArgs[0] == nil {
			err = fmt.Errorf("Default configuration is nil")
			return
		}
		if reflect.TypeOf(defaultArgs[0]).Kind() == reflect.String {
			var configFromString interface{}
			if err = collections.ConvertData(fmt.Sprint(defaultArgs[0]), &configFromString); err != nil {
				err = fmt.Errorf("Function configuration must be valid type: %v\n%v", defaultArgs[0], err)
				return
			}
			defaultArgs[0] = configFromString
		}
		if config, err = collections.TryAsDictionary(defaultArgs[0]); err != nil {
			err = fmt.Errorf("Function configuration must be valid dictionary: %[1]T %[1]v", defaultArgs[0])
			return
		}
	default:
		return "", fmt.Errorf("Too many parameters supplied")
	}

	for key, val := range config.AsMap() {
		switch strings.ToLower(key) {
		case "d", "desc", "description":
			config.Set("description", val)
		case "g", "group":
			config.Set("group", val)
		case "a", "args", "arguments":
			switch val := val.(type) {
			case iList:
				config.Set("args", val)
			default:
				err = fmt.Errorf("%[1]s must be a list of strings: %[2]T %[2]v", key, val)
				return
			}
		case "aliases":
			switch val := val.(type) {
			case iList:
				config.Set("aliases", val)
			default:
				err = fmt.Errorf("%[1]s must be a list of strings: %[2]T %[2]v", key, val)
				return
			}
		case "def", "default", "defaults":
			switch val := val.(type) {
			case iDictionary:
				config.Set("def", val)
			default:
				err = fmt.Errorf("%s must be a dictionary: %T", key, val)
				return
			}
		default:
			return "", fmt.Errorf("Unknown configuration %s", key)
		}
	}

	emptyList := collections.CreateList()
	fi := FuncInfo{
		name:        name,
		group:       defval(config.Get("group"), "User defined functions").(string),
		description: defval(config.Get("description"), "").(string),
		arguments:   defval(config.Get("args"), emptyList).(iList).Strings(),
		aliases:     defval(config.Get("aliases"), emptyList).(iList).Strings(),
	}

	defaultValues := defval(config.Get("def"), collections.CreateDictionary()).(iDictionary)

	fi.in = fmt.Sprintf("%s", strings.Join(fi.arguments, ", "))
	for i := range fi.arguments {
		// We only keep the arg name and get rid of any supplemental information (likely type)
		fi.arguments[i] = strings.Fields(fi.arguments[i])[0]
	}

	fi.function = func(args ...interface{}) (result interface{}, err error) {
		context := collections.CreateDictionary()
		parentContext, err := collections.TryAsDictionary(t.context)
		if err != nil {
			context.Set("DEFAULT", t.context)
		}

		switch len(args) {
		case 1:
			if len(fi.arguments) != 1 {
				switch arg := args[0].(type) {
				case string:
					var out interface{}
					if collections.ConvertData(arg, &out) == nil {
						args[0] = out
					}
				}

				if arg, err := collections.TryAsDictionary(args[0]); err == nil {
					context.Merge(arg, defaultValues, parentContext)
					break
				}
			}
			fallthrough
		default:
			templateContext, err := collections.TryAsDictionary(t.context)
			if err != nil {
				return nil, err
			}

			context.Merge(defaultValues, templateContext)
			for i := range args {
				if i >= len(fi.arguments) {
					context.Set("ARGS", args[i:])
					break
				}
				context.Set(fi.arguments[i], args[i])
			}
		}
		return f(collections.Interface2string(source), context)
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

func (t *Template) tryConvert(value string) interface{} {
	if data, err := t.dataConverter(value); err == nil {
		// The result of the command is a valid data structure
		if reflect.TypeOf(data).Kind() != reflect.String {
			return data
		}
	}
	return value
}

// Execute the command (command could be a file, a template or a script) and convert its result as data if possible
func (t *Template) exec(command string, args ...interface{}) (result interface{}, err error) {
	if result, err = t.run(command, args...); err == nil {
		if result == nil {
			return
		}
		result = t.tryConvert(result.(string))
	}
	return
}

func (t Template) runTemplate(source string, context ...interface{}) (resultContent, filename string, err error) {
	if source == "" {
		return
	}
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
				razor, _ := t.applyRazor(fileContent)
				source = string(razor)
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

	if !t.options[AcceptNoValue] {
		internalTemplate.Option("missingkey=error")
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
	context := make(dictionary)

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
	if err != nil {
		split := strings.SplitN(err.Error(), ">: ", 2)
		if len(split) == 2 {
			// For internal evaluation, we do not want the file/position details on the error
			err = fmt.Errorf(split[1])
		}
	}
	return t.tryConvert(result), err
}

func getAttributes(object interface{}) string {
	if object == nil {
		return ""
	}

	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	numField := 0
	if t.Kind() == reflect.Struct {
		numField = t.NumField()
	}
	result := make([]string, 0, numField)
	for i := 0; i < numField; i++ {
		name := t.Field(i).Name
		if !collections.IsExported(name) {
			continue
		}
		typeName := color.HiBlackString(fmt.Sprint(t.Field(i).Type))
		attrName := color.HiGreenString(name)
		result = append(result, fmt.Sprintf("%s %s", attrName, typeName))
	}
	return strings.Join(result, "\n")
}

func getMethods(object interface{}) string {
	if object == nil {
		return ""
	}

	t := reflect.TypeOf(object)
	result := make([]string, 0, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		result = append(result, FuncInfo{
			name:     t.Method(i).Name,
			function: t.Method(i).Func.Interface(),
		}.getSignature(true))
	}
	return strings.Join(result, "\n")
}

func getSignature(object interface{}) string {
	attributes := getAttributes(object)
	methods := getMethods(object)
	if attributes != "" && methods != "" {
		return attributes + "\n\n" + methods
	}
	return attributes + methods
}

func raise(args ...interface{}) (string, error) {
	return "", fmt.Errorf(utils.FormatMessage(args...))
}

func assertError(test interface{}, args ...interface{}) (string, error) {
	if isZero(test) {
		if len(args) == 0 {
			args = []interface{}{"Assertion failed"}
		}
		return raise(args...)
	}
	return "", nil
}

func assertWarning(test interface{}, args ...interface{}) string {
	if isZero(test) {
		if len(args) == 0 {
			args = []interface{}{"Assertion failed"}
		}
		Log.Warning(utils.FormatMessage(args...))
	}
	return ""
}
