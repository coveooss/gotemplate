package template

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

// FuncInfo contains the information related to a function made available to go template
type FuncInfo struct {
	f       interface{}
	group   string
	aliases []string
	args    []string
	desc    string
	in, out string
	alias   *FuncInfo
}

// Group returns the group name associated to the entry
func (fi FuncInfo) Group() string { return fi.group }

// Aliases returns the aliases related to the entry
func (fi FuncInfo) Aliases() []string { return fi.aliases }

// Description returns the description related to the entry
func (fi FuncInfo) Description() string { return fi.desc }

// String returns the presentation of the FuncInfo entry
func (fi FuncInfo) String() (result string) {
	var r []string
	if fi.alias != nil {
		fi = *fi.alias
	}
	if fi.group != "" {
		r = append(r, fmt.Sprint("Group = ", fi.group))
	}
	if fi.desc != "" {
		r = append(r, fmt.Sprint("Description = ", fi.desc))
	}
	if len(fi.aliases) > 0 {
		r = append(r, fmt.Sprint("Aliases = ", strings.Join(fi.aliases, ", ")))
	}
	r = append(r, fmt.Sprint("Arguments = ", strings.Join(fi.args, ", ")))
	return strings.Join(r, "\n")
}

type funcTableMap map[string]FuncInfo

var converted = make(map[uint]template.FuncMap)

func (ftm funcTableMap) convert() template.FuncMap {
	// The index is a combination of the map address & the length of the map,
	// if either of those change, the item will be updated in the converted
	// cache
	index := uint(reflect.ValueOf(ftm).Pointer()) + uint(len(ftm))
	if result := converted[index]; result != nil {
		return result
	}

	result := make(map[string]interface{}, len(ftm))
	for key, val := range ftm {
		if val.f == nil {
			continue
		}
		result[key] = val.f
		for i := range val.aliases {
			result[val.aliases[i]] = val.f
			// It is necessary here to take a distinct copy of the variable since
			// val will change over the iteration and we cannot rely on its address
			copy := val
			ftm[val.aliases[i]] = FuncInfo{alias: &copy}
		}
	}
	converted[index] = result
	return result
}

// AddFunctions add functions to the template, but keep a detailled definition of the function added for helping purpose
func (t *Template) AddFunctions(funcMap funcTableMap) *Template {
	if t.functions == nil {
		t.functions = make(funcTableMap)
	}
	for key, value := range funcMap {
		t.functions[key] = value
	}
	t.Funcs(funcMap.convert())
	return t
}

// List the available functions in the template
func (t Template) getFunctions() (result []string) {
	for name := range t.functions {
		result = append(result, name)
	}
	sort.Strings(result)
	return
}

// List the available functions in the template
func (t Template) getFunction(name string) FuncInfo {
	return t.functions[name]
}
