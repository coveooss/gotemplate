package template

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/fatih/color"
)

// PrintTemplates output the list of templates available.
func (t *Template) PrintTemplates(all, long bool) {
	templates := t.getTemplateNames()
	var maxLen int
	for _, template := range templates {
		t := t.Lookup(template)
		if len(template) > maxLen && template != t.ParseName {
			maxLen = len(template)
		}
	}

	faint := color.New(color.Faint).SprintfFunc()

	for _, template := range templates {
		tpl := t.Lookup(template)
		if all || tpl.Name() != tpl.ParseName {
			name := tpl.Name()
			if tpl.Name() == tpl.ParseName {
				name = ""
			}
			folder := utils.Relative(t.folder, tpl.ParseName)
			if folder+name != "." {
				ErrPrintf("%-[3]*[1]s %[2]s\n", name, faint(folder), maxLen)
			}
		}
	}
	ErrPrintln()
}

// PrintFunctions outputs the list of functions available.
func (t *Template) PrintFunctions(all, long, groupByCategory bool, filters ...string) {
	functions := t.filterFunctions(all, groupByCategory, long, filters...)

	maxLength := 0
	categories := make(map[string][]string)
	for i := range functions {
		var group string
		if groupByCategory {
			funcInfo := t.functions[functions[i]]
			if funcInfo.alias != nil {
				funcInfo = funcInfo.alias
			}
			group = funcInfo.group
		}
		categories[group] = append(categories[group], functions[i])
		maxLength = int(math.Max(float64(len(functions[i])), float64(maxLength)))
	}

	keys := make([]string, 0, len(categories))
	for key := range categories {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	print := t.printFunctionsShort
	if long {
		print = t.printFunctionsDetailed
	}

	for _, key := range keys {
		if key != "" {
			title, link := collections.Split2(key, ", http")
			title = color.New(color.Underline, color.FgYellow).Sprint(title)
			if link != "" {
				link = color.BlackString(fmt.Sprintf(" http%s", link))
			}
			Printf("%s%s\n\n", title, link)
		}
		print(categories[key], maxLength, all)
		Println()
	}
}

func (t *Template) filterFunctions(all, category, detailed bool, filters ...string) []string {
	functions := t.getAllFunctions()
	if all && len(filters) == 0 {
		return functions
	}

	for i := range filters {
		filters[i] = strings.ToLower(filters[i])
	}

	filtered := make([]string, 0, len(functions))
	for i := range functions {
		funcInfo := t.functions[functions[i]]
		if funcInfo.alias != nil {
			if !all {
				continue
			}
			funcInfo = funcInfo.alias
		}

		if len(filters) == 0 {
			filtered = append(filtered, functions[i])
			continue
		}

		search := strings.ToLower(functions[i] + " " + strings.Join(funcInfo.aliases, " "))
		if category {
			search += " " + strings.ToLower(funcInfo.group)
		}
		if detailed {
			search += " " + strings.ToLower(funcInfo.description)
		}

		for f := range filters {
			if strings.Contains(search, filters[f]) {
				filtered = append(filtered, functions[i])
				break
			}
		}
	}
	return filtered
}

func (t *Template) printFunctionsShort(functions []string, maxLength int, alias bool) {
	const nbColumn = 5
	l := len(functions)
	colLength := int(math.Ceil(float64(l) / float64(nbColumn)))
	for i := 0; i < colLength*nbColumn; i += nbColumn {
		for j := 0; j < nbColumn; j++ {
			pos := j*colLength + i/nbColumn
			if pos >= l {
				continue
			}
			item, extraLen := functions[pos], 0

			if t.functions[item].alias != nil {
				ex := len(color.HiBlackString(""))
				Printf("%-[1]*[2]s", maxLength+2+ex, color.HiBlackString(item))
			} else {
				Printf("%-[1]*[2]s", maxLength+2+extraLen, item)
			}
		}
		Println()
	}
}

func (t *Template) printFunctionsDetailed(functions []string, maxLength int, alias bool) {
	t.options[Razor] = true
	t.completeExamples()

	// We only print entries that are not alias
	allFunc := make(map[string]int)
	for i := range functions {
		funcInfo := t.functions[functions[i]]
		if funcInfo.alias == nil {
			allFunc[functions[i]]++
		} else {
			allFunc[funcInfo.alias.name]++
		}
	}
	functions = make([]string, 0, len(functions))
	for f := range allFunc {
		functions = append(functions, f)
	}
	sort.Strings(functions)

	for i := range functions {
		fi := t.functions[functions[i]]
		if fi.description != "" {
			text := String(fi.description).Wrap(100).Indent("// ").Lines().TrimSuffix(" ").Join("\n").String()
			Println(color.GreenString(text))
		}
		Println(fi.Signature())

		if alias {
			sort.Strings(fi.aliases)
			for j := range fi.aliases {
				aliasFunc := t.functions[fi.aliases[j]]
				if !aliasFunc.IsAlias() || aliasFunc.Arguments() != fi.Arguments() || aliasFunc.Result() != fi.Result() {
					// The alias has been replaced
					continue
				}
				Println(aliasFunc.Signature())
			}
		}
		title := color.MagentaString("\nExample:")
		for _, ex := range fi.Examples() {
			Println(title)
			title = ""
			Print(toStringClass(ex).IndentN(4))
		}
		Println()
	}
}
