package template

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/coveo/gotemplate/types"
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
}

// Name returns the name related to the entry.
func (fi FuncInfo) Name() string { return fi.name }

// Group returns the group name associated to the entry.
func (fi FuncInfo) Group() string { return ifUndef(&fi, fi.alias).(*FuncInfo).group }

// Aliases returns the aliases related to the entry.
func (fi FuncInfo) Aliases() []string { return ifUndef(&fi, fi.alias).(*FuncInfo).aliases }

// Description returns the description related to the entry.
func (fi FuncInfo) Description() string { return ifUndef(&fi, fi.alias).(*FuncInfo).description }

// Signature returns the function signature
func (fi FuncInfo) Signature() string {
	col := color.HiBlueString
	name := fi.name
	if fi.alias != nil {
		fi = *fi.alias
		col = color.HiBlackString
	}

	return fmt.Sprintf("%s(%s) %s", col(name), fi.Arguments(), color.HiBlackString(fi.Result()))
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

// Arguments returns the list of arguments that must be supplied to the function.
func (fi FuncInfo) Arguments() string {
	if fi.in != "" {
		return fi.in
	}

	signature := reflect.ValueOf(fi.function).Type()
	var parameters []string
	for i := 0; i < signature.NumIn(); i++ {
		arg := strings.Replace(fmt.Sprint(signature.In(i)), "interface {}", "interface{}", -1)
		arg = strings.Replace(arg, "types.", "", -1)
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
		parameters = append(parameters, fmt.Sprintf("%s %s", argName, color.CyanString(arg)))
	}
	return strings.Join(parameters, ", ")
}

// Result returns the list of output produced by the function.
func (fi FuncInfo) Result() string {
	if fi.out != "" {
		return fi.out
	}
	signature := reflect.ValueOf(fi.function).Type()
	var outputs []string
	for i := 0; i < signature.NumOut(); i++ {
		r := strings.Replace(fmt.Sprint(signature.Out(i)), "interface {}", "interface{}", -1)
		r = strings.Replace(r, "types.", "", -1)
		outputs = append(outputs, r)
	}
	return strings.Join(outputs, ", ")
}

type funcTableMap map[string]FuncInfo

func (ftm funcTableMap) convert() template.FuncMap {
	// The index is a combination of the map address & the length of the map,
	// if either of those change, the item will be updated in the converted
	// cache
	index := uint(reflect.ValueOf(ftm).Pointer()) + uint(len(ftm))
	if result := converted[index]; result != nil {
		return result
	}

	result := types.CreateDictionary(len(ftm))
	for key, val := range ftm {
		if val.function == nil {
			continue
		}
		result.Set(key, val.function)
	}
	converted[index] = result.AsMap()
	return result.AsMap()
}

var converted = make(map[uint]template.FuncMap)

type funcOptionsSet int8

const (
	funcHelp funcOptionsSet = iota
	funcArgs
	funcAliases
	funcGroup
)

type funcOptions map[funcOptionsSet]interface{}
type aliases map[string][]string
type arguments map[string][]string
type descriptions map[string]string
type dictionary map[string]interface{}
type groups map[string]string

// AddFunctions add functions to the template, but keep a detailled definition of the function added for helping purpose
func (t *Template) AddFunctions(funcs dictionary, group string, options funcOptions) *Template {
	ft := make(funcTableMap, len(funcs))
	help := defval(options[funcHelp], descriptions{}).(descriptions)
	aliases := defval(options[funcAliases], aliases{}).(aliases)
	arguments := defval(options[funcArgs], arguments{}).(arguments)
	groups := defval(options[funcGroup], groups{}).(groups)
	for key, val := range funcs {
		ft[key] = FuncInfo{
			function:    val,
			group:       defval(group, groups[key]).(string),
			aliases:     aliases[key],
			arguments:   arguments[key],
			description: help[key],
		}
	}

	return t.addFunctions(ft)
}

func (t *Template) addFunctions(funcMap funcTableMap) *Template {
	if t.functions == nil {
		t.functions = make(funcTableMap)
	}
	for key, value := range funcMap {
		value.name = key
		t.functions[key] = value
		for i := range value.aliases {
			// It is necessary here to take a distinct copy of the variable since
			// val will change over the iteration and we cannot rely on its address
			copy := value
			funcMap[value.aliases[i]] = FuncInfo{alias: &copy, function: value.function, name: value.aliases[i]}
			t.functions[value.aliases[i]] = funcMap[value.aliases[i]]
		}
	}
	t.Funcs(funcMap.convert())
	return t
}

// List the available functions in the template
func (t Template) getFunctionsInternal(original, alias bool) (result []string) {
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

func (t Template) getAliases() []string      { return t.getFunctionsInternal(false, true) }
func (t Template) getAllFunctions() []string { return t.getFunctionsInternal(true, true) }
func (t Template) getFunctions() []string    { return t.getFunctionsInternal(true, false) }

// List the available functions in the template
func (t Template) getFunction(name string) FuncInfo {
	return t.functions[name]
}
