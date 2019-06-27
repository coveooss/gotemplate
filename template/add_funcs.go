package template

import (
	"sort"
	"sync"

	"github.com/coveooss/gotemplate/v3/utils"
)

const (
	goTemplateBase = "Base go template functions"
)

var addFuncMutex sync.Mutex

// Add additional functions to the go template context
func (t *Template) addFuncs() {
	// We cannot create functions table in multiple threads concurrently
	addFuncMutex.Lock()
	defer addFuncMutex.Unlock()

	if baseGoTemplateFuncs == nil {
		baseGoTemplateFuncs = make(funcTableMap, len(baseGoTemplate))
		for key, val := range baseGoTemplate {
			baseGoTemplateFuncs[key] = FuncInfo{
				group:       goTemplateBase,
				description: val.description,
				in:          val.args,
				out:         val.out,
			}
		}
	}
	t.addFunctions(baseGoTemplateFuncs)

	add := func(o Options, f func()) {
		if t.options[o] {
			f()
		}
	}

	add(Sprig, t.addSprigFuncs)
	add(Math, t.addMathFuncs)
	add(Data, t.addDataFuncs)
	add(Logging, t.addLoggingFuncs)
	add(Runtime, t.addRuntimeFuncs)
	add(Utils, t.addUtilsFuncs)
	add(Net, t.addNetFuncs)
	add(OS, t.addOSFuncs)
}

// Apply all regular expressions replacements to the supplied string
func (t Template) substitute(content string) string {
	return utils.Substitute(content, t.substitutes...)
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

var baseGoTemplateFuncs funcTableMap

var baseGoTemplate = map[string]struct {
	description, args, out string
}{
	"and": {
		`Returns the boolean AND of its arguments by returning the first empty argument or the last argument, that is, "and x y" behaves as "if x then y else x". All the arguments are evaluated.`,
		`arg0 reflect.Value, args ...reflect.Value`, `reflect.Value`,
	},
	"call": {
		`Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters. Thus "call .X.Y 1 2" is, in Go notation, dot.X.Y(1, 2) where Y is a func-valued field, map entry, or the like. The first argument must be the result of an evaluation that yields a value of function type (as distinct from a predefined function such as print). The function must return either one or two result values, the second of which is of type error. If the arguments don't match the function or the returned error value is non-nil, execution stops.`,
		`fn reflect.Value, args ...reflect.Value`, `reflect.Value, error)`,
	},
	"html": {
		`Returns the escaped HTML equivalent of the textual representation of its arguments. This function is unavailable in html/template, with a few exceptions.`,
		`args ...interface{}`, `string`,
	},
	"index": {
		`Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.`,
		`item reflect.Value, indices ...reflect.Value`, `(reflect.Value, error)`,
	},
	"js": {
		`Returns the escaped JavaScript equivalent of the textual representation of its arguments.`,
		`args ...interface{}`, `string`,
	},
	"len": {
		`Returns the integer length of its argument.`,
		`item interface{}`, `(int, error)`,
	},
	"not": {
		`Returns the boolean negation of its single argument.`,
		`not(arg reflect.Value`, `bool`,
	},
	"or": {
		`Returns the boolean OR of its arguments by returning the first non-empty argument or the last argument, that is, "or x y" behaves as "if x then x else y". All the arguments are evaluated.`,
		`or(arg0 reflect.Value, args ...reflect.Value`, `reflect.Value`,
	},
	"print": {
		`An alias for fmt.Sprint`,
		`args ...interface{}`, `string`,
	},
	"printf": {
		`An alias for fmt.Sprintf`,
		`format string, args ...interface{}`, `string`,
	},
	"println": {
		`An alias for fmt.Sprintln`,
		`args ...interface{}`, `string`,
	},
	"urlquery": {
		`Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query. This function is unavailable in html/template, with a few exceptions.`,
		`args ...interface{}`, `string`,
	},

	"eq": {
		`Returns the boolean truth of arg1 == arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
	"ge": {
		`Returns the boolean truth of arg1 >= arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
	"gt": {
		`Returns the boolean truth of arg1 > arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
	"le": {
		`Returns the boolean truth of arg1 <= arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
	"lt": {
		`Returns the boolean truth of arg1 < arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
	"ne": {
		`Returns the boolean truth of arg1 != arg2`,
		`arg1 reflect.Value, arg2 ...reflect.Value`, `(bool, error)`,
	},
}
