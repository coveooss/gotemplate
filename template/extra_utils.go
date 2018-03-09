package template

import (
	"fmt"
	"regexp"

	"github.com/coveo/gotemplate/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	utilsBase = "Other utilities functions"
)

var utilsFuncs = funcTableMap{
	"concat":     {utils.Concat, utilsBase, nil, []string{}, ""},
	"formatList": {utils.FormatList, utilsBase, nil, []string{"format", "list"}, ""},
	"joinLines":  {utils.JoinLines, utilsBase, nil, nil, ""},
	"mergeList":  {utils.MergeLists, utilsBase, nil, []string{"lists"}, ""},
	"splitLines": {utils.SplitLines, utilsBase, nil, []string{}, ""},
	"id":         {id, utilsBase, nil, []string{"identifier", "replaceChar"}, ""},
	"center":     {utils.CenterString, utilsBase, nil, []string{}, ""},
	"glob":       {glob, utilsBase, nil, nil, ""},
	"pwd":        {utils.Pwd, utilsBase, nil, nil, "Returns the current working directory"},
	"iif":        {utils.IIf, utilsBase, nil, []string{"test", "valueIfTrue", "valueIfFalse"}, ""},
	"lorem":      {lorem, utilsBase, nil, []string{"funcName"}, ""},
	"color":      {utils.SprintColor, utilsBase, nil, nil, ""},
	"diff":       {diff, utilsBase, nil, nil, ""},
}

func (t *Template) addUtilsFuncs() {
	t.AddFunctions(utilsFuncs)
}

func glob(args ...interface{}) []string { return utils.GlobFuncTrim(args...) }

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

func diff(text1, text2 interface{}) interface{} {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(fmt.Sprint(text1), fmt.Sprint(text2), true)
	return dmp.DiffPrettyText(diffs)
}
