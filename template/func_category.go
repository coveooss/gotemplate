package template

import (
	"sort"
)

// FuncCategory represents a group of functions of the same group.
type FuncCategory struct {
	name      string
	functions []string
}

// Name returns the name related to the entry.
func (fc FuncCategory) Name() string { return fc.name }

// Functions returns the list of functions associated with the category.
func (fc FuncCategory) Functions() []string { return fc.functions }

func (t *Template) getCategories() []FuncCategory {
	categories := make(map[string][]string)
	for name := range t.functions {
		fi := t.functions[name]
		if fi.alias != nil {
			fi = *fi.alias
		}
		categories[fi.group] = append(categories[fi.group], name)
	}

	categoryList := make([]string, 0, len(categories))
	for key := range categories {
		categoryList = append(categoryList, key)
	}
	sort.Strings(categoryList)

	result := make([]FuncCategory, len(categoryList))
	for i := range categoryList {
		sort.Strings(categories[categoryList[i]])
		result[i] = FuncCategory{
			name:      categoryList[i],
			functions: categories[categoryList[i]],
		}
	}
	return result
}
