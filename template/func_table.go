package template

import (
	"reflect"
	"sort"
	"text/template"
)

type funcTable struct {
	function    interface{}
	group       string
	aliases     []string
	argNames    []string
	description string
}

type funcTableMap map[string]funcTable

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
		if val.function == nil {
			continue
		}
		result[key] = val.function
		for i := range val.aliases {
			result[val.aliases[i]] = val.function
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
func (t Template) getFunctions(all bool) []string {
	var functions []string
	for name := range t.functions {
		functions = append(functions, name)
		if all {
			aliases := t.functions[name].aliases
			for i := range aliases {
				functions = append(functions, aliases[i])
			}
		}
	}
	sort.Strings(functions)
	return functions
}
