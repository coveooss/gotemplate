package template

import (
	"github.com/Masterminds/sprig"
)

const (
	sprigFunctions = "Sprig functions http://masterminds.github.io/sprig"
)

var sprigFuncs funcTableMap

var sprigFuncMap = sprig.GenericFuncMap()

func (t *Template) addSprigFuncs() {
	if sprigFuncs == nil {
		sprigFuncs = make(funcTableMap)
		for key, value := range sprigFuncMap {
			sprigFuncs[key] = funcTable{value, sprigFunctions, nil, []string{}, ""}
		}
	}

	t.AddFunctions(sprigFuncs)
}
