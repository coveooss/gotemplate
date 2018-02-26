package template

import (
	"fmt"
	"regexp"

	"github.com/coveo/gotemplate/utils"
)

const (
	utilsBase = "Other utilities functions"
)

var utilsFuncs funcTableMap

func (t *Template) addUtilsFuncs() {
	if utilsFuncs == nil {
		utilsFuncs = funcTableMap{
			"concat":     {utils.Concat, utilsBase, nil, []string{}, ""},
			"formatList": {utils.FormatList, utilsBase, nil, []string{}, ""},
			"joinLines":  {utils.JoinLines, utilsBase, nil, []string{}, ""},
			"mergeList":  {utils.MergeLists, utilsBase, nil, []string{}, ""},
			"splitLines": {utils.SplitLines, utilsBase, nil, []string{}, ""},
			"id":         {id, utilsBase, nil, []string{}, ""},
			"center":     {utils.CenterString, utilsBase, nil, []string{}, ""},
			"glob":       {utils.GlobFunc, utilsBase, nil, []string{}, ""},
			"pwd":        {utils.Pwd, utilsBase, nil, []string{}, ""},
			"iif":        {utils.IIf, utilsBase, nil, []string{}, ""},
			"lorem":      {lorem, utilsBase, nil, []string{}, ""},
			"color":      {utils.SprintColor, utilsBase, nil, []string{}, ""},
		}
	}

	t.AddFunctions(utilsFuncs)
}

func lorem(funcName interface{}, params ...int) (result string, err error) {
	kind, err := utils.GetLoremKind(fmt.Sprint(funcName))
	if err == nil {
		result, err = utils.Lorem(kind, params...)
	}
	return
}

func id(id string, replace ...interface{}) string {
	// By default, replacement char for invalid chars would be _
	// but it is possible to specify an alternative string to act as the replacement
	replacement := fmt.Sprint(replace...)
	if replacement == "" {
		replacement = "_"
	}

	dup := duplicateUnderscore
	if replacement != "_" {
		// If the replacement string is not the default one, we generate a special substituter to remove duplicates
		// taking into account regex special chars such as +, ?, etc.
		dup = regexp.MustCompile(fmt.Sprintf(`(?:%s)+`, regexSpecial.ReplaceAllString(replacement, `\$0`)))
	}

	return dup.ReplaceAllString(validChars.ReplaceAllString(id, replacement), replacement)
}

var validChars = regexp.MustCompile(`[^\p{L}\d_]`)
var duplicateUnderscore = regexp.MustCompile(`_+`)
var regexSpecial = regexp.MustCompile(`[\+\.\?\(\)\[\]\{\}\\]`)
