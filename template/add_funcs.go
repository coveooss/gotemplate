package template

import (
	"sort"

	"github.com/coveo/gotemplate/utils"
)

// Add additional functions to the go template context
func (t *Template) addFuncs() {
	if t.options[Sprig] {
		t.addSprigFuncs()
	}

	if t.options[Math] {
		t.addMathFuncs()
	}

	if t.options[Data] {
		t.addDataFuncs()
	}

	if t.options[Logging] {
		t.addLoggingFuncs()
	}

	if t.options[Runtime] {
		t.addRuntimeFuncs()
	}

	if t.options[Utils] {
		t.addUtilsFuncs()
	}
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
