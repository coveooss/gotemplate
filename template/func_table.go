package template

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/fatih/color"
)

// FuncInfo contains the information related to a function made available to go template.
type FuncInfo struct {
	name        string
	function    interface{}
	group       string
	aliases     []string
	arguments   []string
	description string
	in, out     string
	alias       *FuncInfo
	examples    []Example
}

// Example can be added to a function to describe how to use it.
type Example struct {
	Razor    string
	Template string
	Result   string
}

func (e Example) String() (result string) {
	section := color.HiBlackString
	if e.Razor != "" {
		result += fmt.Sprintln(section("Razor:       "), String(e.Razor).IndentN(14).TrimSpace())
	}
	if e.Template != "" {
		result += fmt.Sprintln(section("Go Template: "), String(e.Template).IndentN(14).TrimSpace())
	}
	if e.Result != "" {
		result += fmt.Sprintln(section("Result:      "), String(e.Result).IndentN(14).TrimSpace())
	}
	return result
}

// Aliases returns the aliases related to the entry.
func (fi FuncInfo) Aliases() []string { return ifUndef(&fi, fi.alias).(*FuncInfo).aliases }

// Arguments returns the list of arguments that must be supplied to the function.
func (fi FuncInfo) Arguments() string { return fi.getArguments(false) }

// Description returns the description related to the entry.
func (fi FuncInfo) Description() string { return ifUndef(&fi, fi.alias).(*FuncInfo).description }

// Group returns the group name associated to the entry.
func (fi FuncInfo) Group() string { return ifUndef(&fi, fi.alias).(*FuncInfo).group }

// Name returns the name related to the entry.
func (fi FuncInfo) Name() string { return fi.name }

// Signature returns the function signature.
func (fi FuncInfo) Signature() string { return fi.getSignature(false) }

// IsAlias indicates if the current function is an alias.
func (fi FuncInfo) IsAlias() bool { return fi.alias != nil }

// RealName returns the real name of a function wether it's an alias or not.
func (fi FuncInfo) RealName() string {
	if fi.alias == nil {
		return fi.name
	}
	return fi.alias.name
}

// String returns the presentation of the FuncInfo entry.
func (fi FuncInfo) String() (result string) {
	signature := fi.Signature()
	if fi.alias != nil {
		fi = *fi.alias
	}

	if fi.group != "" {
		result += fmt.Sprintf(color.GreenString("Group: %s\n"), fi.group)
	}
	if fi.description != "" {
		result += fmt.Sprintf(color.GreenString("%s\n"), fi.description)
	}
	return result + signature
}

// Examples returns the list of examples associated to the function.
func (fi FuncInfo) Examples() []Example { return fi.examples }

func (fi FuncInfo) getSignature(isMethod bool) string {
	col := color.HiBlueString
	name := fi.name
	if fi.alias != nil {
		fi = *fi.alias
		col = color.HiBlackString
	}

	return fmt.Sprintf("%s(%s) %s", col(name), fi.getArguments(isMethod), color.HiBlackString(fi.Result()))
}

func (fi FuncInfo) getArguments(isMethod bool) string {
	if fi.in != "" || fi.function == nil {
		return fi.in
	}

	if fi.alias != nil {
		fi = *fi.alias
	}

	signature := reflect.ValueOf(fi.function).Type()
	var parameters []string
	for i := iif(isMethod, 1, 0).(int); i < signature.NumIn(); i++ {
		arg := strings.Replace(fmt.Sprint(signature.In(i)), "interface {}", "interface{}", -1)
		arg = strings.Replace(arg, "collections.", "", -1)
		var argName string
		if i < len(fi.arguments) {
			argName = fi.arguments[i]
		} else {
			if signature.IsVariadic() && i == signature.NumIn()-1 {
				argName = "args"
			} else {
				argName = fmt.Sprintf("arg%d", i+1)
			}
		}
		if signature.IsVariadic() && i == signature.NumIn()-1 {
			arg = "..." + arg[2:]
		}
		if isMethod {
			parameters = append(parameters, color.CyanString(arg))
		} else {
			parameters = append(parameters, fmt.Sprintf("%s %s", argName, color.CyanString(arg)))
		}
	}
	return strings.Join(parameters, ", ")
}

// Result returns the list of output produced by the function.
func (fi FuncInfo) Result() string {
	if fi.out != "" || fi.function == nil {
		return fi.out
	}
	signature := reflect.ValueOf(fi.function).Type()
	var outputs []string
	for i, nb := 0, signature.NumOut(); i < nb; i++ {
		r := fmt.Sprint(signature.Out(i))
		if r == "error" && i == nb-1 {
			break
		}
		r = strings.Replace(r, "interface {}", "interface{}", -1)
		r = strings.Replace(r, "collections.", "", -1)
		outputs = append(outputs, r)
	}
	if len(outputs) == 1 {
		return outputs[0]

	}
	return fmt.Sprintf("(%s)", strings.Join(outputs, ", "))
}

type funcTableMap map[string]*FuncInfo

func (ftm funcTableMap) convert() template.FuncMap {
	result := collections.CreateDictionary(len(ftm))
	for key, val := range ftm {
		if val.function == nil {
			continue
		}
		result.Set(key, val.function)
	}
	return result.AsMap()
}

// FuncOptionsSet defines categories that could be used to define elaborated help context when adding functions to go template.
type FuncOptionsSet int8

const (
	// FuncHelp is used to associate help to functions added to go templates.
	FuncHelp FuncOptionsSet = iota
	// FuncArgs is used to associate arguments name to functions added to go templates.
	FuncArgs
	// FuncAliases is used to associate aliases to functions added to go templates.
	FuncAliases
	// FuncGroup is used to associate a group to functions added to go templates.
	FuncGroup
	// FuncExamples is used to associate examples (from razor to template to result) to functions added to go templates.
	FuncExamples
)

// FuncOptions is a map of FuncOptionsSet that is used to associates help, aliases, arguments and groups to functions added to go template.
type FuncOptions map[FuncOptionsSet]interface{}

type aliases = map[string][]string
type arguments = map[string][]string
type descriptions = map[string]string
type dictionary = map[string]interface{}
type examples = map[string][]Example
type groups = map[string]string

// AddFunctions add functions to the template, but keep a detailed definition of the function added for helping purpose
func (t *Template) AddFunctions(funcs dictionary, group string, options FuncOptions) *Template {
	ft := make(funcTableMap, len(funcs))
	help := defval(options[FuncHelp], descriptions{}).(descriptions)
	aliases := defval(options[FuncAliases], aliases{}).(aliases)
	arguments := defval(options[FuncArgs], arguments{}).(arguments)
	groups := defval(options[FuncGroup], groups{}).(groups)
	examples := defval(options[FuncExamples], examples{}).(examples)
	for key, val := range funcs {
		funcInfo := &FuncInfo{
			function:    val,
			group:       defval(group, groups[key]).(string),
			aliases:     aliases[key],
			arguments:   arguments[key],
			description: help[key],
			examples:    examples[key],
		}
		ft[key] = funcInfo
	}

	return t.addFunctions(ft)
}

func (t *Template) addFunctions(funcMap funcTableMap) *Template {
	templateMutex.Lock()
	defer templateMutex.Unlock()

	if t.functions == nil {
		t.functions = make(funcTableMap)
	}
	for key, value := range funcMap {
		value.name = key
		t.functions[key] = value
		for i := range value.aliases {
			funcMap[value.aliases[i]] = &FuncInfo{alias: value, function: value.function, name: value.aliases[i]}
			t.functions[value.aliases[i]] = funcMap[value.aliases[i]]
		}
	}
	t.Funcs(funcMap.convert())
	return t
}

// Completes example that are not fully described.
func (t *Template) completeExamples() string {
	for _, function := range t.functions {
		for i := range function.examples {
			example := &function.examples[i]
			if example.Razor != "" && example.Template == "" {
				result, _ := t.applyRazor([]byte(example.Razor))
				example.Template = string(result)
			}
			if example.Template != "" && example.Result == "" {
				var err error
				if example.Result, err = t.ProcessContent(example.Template, "."); err != nil {
					InternalLog.Errorf("Error while ")
				}
			}
		}
	}
	return ""
}

// List the available functions in the template
func (t *Template) getFunctionsInternal(original, alias bool) (result []string) {
	for name := range t.functions {
		fi := t.functions[name]
		if original && fi.alias == nil {
			result = append(result, name)
		}
		if alias && fi.alias != nil {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return
}

func (t *Template) getAliases() []string      { return t.getFunctionsInternal(false, true) }
func (t *Template) getAllFunctions() []string { return t.getFunctionsInternal(true, true) }
func (t *Template) getFunctions() []string    { return t.getFunctionsInternal(true, false) }

// List the available functions in the template
func (t *Template) getFunction(name string) *FuncInfo {
	return t.functions[name]
}
