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
	"concat":     {f: utils.Concat, group: utilsBase, desc: ""},
	"formatList": {f: utils.FormatList, group: utilsBase, args: []string{"format", "list"}, desc: ""},
	"joinLines":  {f: utils.JoinLines, group: utilsBase, desc: ""},
	"mergeList":  {f: utils.MergeLists, group: utilsBase, args: []string{"lists"}, desc: ""},
	"splitLines": {f: utils.SplitLines, group: utilsBase, args: []string{}, desc: ""},
	"id":         {f: id, group: utilsBase, args: []string{"identifier", "replaceChar"}, desc: ""},
	"center":     {f: center, group: utilsBase, args: []string{"width", "str"}, desc: ""},
	"glob":       {f: glob, group: utilsBase, desc: ""},
	"wrap":       {f: wrap, group: utilsBase, args: []string{"width", "s"}, desc: ""},
	"pwd":        {f: utils.Pwd, group: utilsBase, desc: "Returns the current working directory"},
	"iif":        {f: utils.IIf, group: utilsBase, args: []string{"test", "valueIfTrue", "valueIfFalse"}, desc: ""},
	"lorem":      {f: lorem, group: utilsBase, args: []string{"funcName"}, desc: ""},
	"color":      {f: utils.SprintColor, group: utilsBase, desc: ""},
	"diff":       {f: diff, group: utilsBase, desc: ""},
	"repeat":     {f: repeat, group: utilsBase, args: []string{"n", "item"}, desc: "Returns an array with the item repeated n times."},
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
