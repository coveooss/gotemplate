package template

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	utilsBase = "Other utilities"
)

var utilsFuncs = funcTableMap{
	"concat":     {function: utils.Concat, group: utilsBase, description: ""},
	"formatList": {function: utils.FormatList, group: utilsBase, arguments: []string{"format", "list"}, description: ""},
	"joinLines":  {function: utils.JoinLines, group: utilsBase, description: ""},
	"mergeList":  {function: utils.MergeLists, group: utilsBase, arguments: []string{"lists"}, description: ""},
	"splitLines": {function: utils.SplitLines, group: utilsBase, arguments: []string{}, description: ""},
	"id":         {function: id, group: utilsBase, arguments: []string{"identifier", "replaceChar"}, description: ""},
	"center":     {function: center, group: utilsBase, arguments: []string{"width", "str"}, description: ""},
	"glob":       {function: glob, group: utilsBase, description: ""},
	"wrap":       {function: wrap, group: utilsBase, arguments: []string{"width", "s"}, description: ""},
	"pwd":        {function: utils.Pwd, group: utilsBase, description: "Returns the current working directory"},
	"iif":        {function: utils.IIf, group: utilsBase, arguments: []string{"test", "valueIfTrue", "valueIfFalse"}, description: ""},
	"lorem":      {function: lorem, group: utilsBase, arguments: []string{"funcName"}, description: ""},
	"color":      {function: utils.SprintColor, group: utilsBase, description: ""},
	"diff":       {function: diff, group: utilsBase, description: ""},
	"repeat":     {function: repeat, group: utilsBase, arguments: []string{"n", "item"}, description: "Returns an array with the item repeated n times."},
}

func (t *Template) addUtilsFuncs() {
	t.addFunctions(utilsFuncs)
}

func glob(args ...interface{}) []string { return utils.GlobFuncTrim(args...) }

func lorem(funcName interface{}, params ...int) (result string, err error) {
	kind, err := utils.GetLoremKind(fmt.Sprint(funcName))
	if err == nil {
		result, err = utils.Lorem(kind, params...)
	}
	return
}

func center(width, s interface{}) string {
	w := errors.Must(strconv.Atoi(fmt.Sprintf("%v", width))).(int)
	return utils.CenterString(fmt.Sprint(s), w)
}

func wrap(width, s interface{}) string {
	w := errors.Must(strconv.Atoi(fmt.Sprintf("%v", width))).(int)
	return utils.WrapString(fmt.Sprint(s), w)
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

func repeat(n int, a interface{}) (result []interface{}, err error) {
	if n < 0 {
		err = fmt.Errorf("n must be greater or equal than 0")
		return
	}
	result = make([]interface{}, n)
	for i := range result {
		result[i] = a
	}
	return
}
