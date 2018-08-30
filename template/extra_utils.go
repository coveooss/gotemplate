package template

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	utilsBase = "Other utilities"
)

var utilsFuncs = dictionary{
	"assert":     assert,
	"center":     center,
	"color":      utils.SprintColor,
	"concat":     collections.Concat,
	"diff":       diff,
	"formatList": utils.FormatList,
	"glob":       glob,
	"id":         id,
	"iif":        utils.IIf,
	"joinLines":  collections.JoinLines,
	"lorem":      lorem,
	"mergeList":  utils.MergeLists,
	"pwd":        utils.Pwd,
	"raise":      raise,
	"repeat":     repeat,
	"sIndent":    indent,
	"save":       saveToFile,
	"splitLines": collections.SplitLines,
	"wrap":       wrap,
}

var utilsFuncsArgs = arguments{
	"assert":     {"test", "message", "arguments"},
	"center":     {"width"},
	"diff":       {"text1", "text2"},
	"formatList": {"format", "list"},
	"id":         {"identifier", "replaceChar"},
	"iif":        {"testValue", "valueTrue", "valueFalse"},
	"joinLines":  {"format"},
	"lorem":      {"loremType", "params"},
	"mergeList":  {"lists"},
	"raise":      {"message", "arguments"},
	"repeat":     {"n", "element"},
	"sIndent":    {"spacer"},
	"save":       {"filename", "object"},
	"splitLines": {"content"},
	"wrap":       {"width"},
}

var utilsFuncsAliases = aliases{
	"assert":  {"assertion"},
	"center":  {"centered"},
	"color":   {"colored", "enhanced"},
	"diff":    {"difference"},
	"glob":    {"expand"},
	"id":      {"identifier"},
	"iif":     {"ternary"},
	"lorem":   {"loremIpsum"},
	"pwd":     {"currentDir"},
	"raise":   {"raiseError"},
	"sIndent": {"sindent", "spaceIndent"},
	"save":    {"write", "writeTo"},
	"wrap":    {"wrapped"},
}

var utilsFuncsHelp = descriptions{
	"assert": "Raise a formated error if the test condition is false.",
	"center": "Returns the concatenation of supplied arguments centered within width.",
	"color": strings.TrimSpace(collections.UnIndent(`
		Colors the rendered string.

		The first arguments are interpretated as color attributes until the first non color attribute. Attributes are case insensitive.

		Valid attributes are:
		    Reset, Bold, Faint, Italic, Underline, BlinkSlow, BlinkRapid, ReverseVideo, Concealed, CrossedOut

		Valid color are:
		    Black, Red, Green, Yellow, Blue, Magenta, Cyan, White

		Color can be prefixed by:
		    Fg:   Meaning foreground (Fg is assumed if not specified)
		    FgHi: Meaning high intensity forground
		    Bg:   Meaning background"
		    BgHi: Meaning high intensity background
	`)),
	"concat":     "Returns the concatenation (without separator) of the string representation of objects.",
	"diff":       "Returns a colored string that highlight differences between supplied texts.",
	"formatList": "Return a list of strings by applying the format to each element of the supplied list.",
	"glob":       "Returns the expanded list of supplied arguments (expand *[]? on filename).",
	"id":         "Returns a valid go identifier from the supplied string (replacing any non compliant character by replacement, default _ ).",
	"iif":        "If testValue is empty, returns falseValue, otherwise returns trueValue.\n    WARNING: All arguments are evaluated and must by valid.",
	"joinLines":  "Merge the supplied objects into a newline separated string.",
	"lorem":      "Returns a random string. Valid types are be word, words, sentence, para, paragraph, host, email, url.",
	"mergeList":  "Return a single list containing all elements from the lists supplied.",
	"pwd":        "Returns the current working directory.",
	"raise":      "Raise a formated error.",
	"repeat":     "Returns an array with the item repeated n times.",
	"save":       "Save object to file.",
	"sIndent": strings.TrimSpace(collections.UnIndent(`
		Intents the the elements using the provided spacer.
		
		You can also use autoIndent as Razor expression if you don't want to specify the spacer.
		Spacer will then be auto determined by the spaces that precede the expression.
		Valid aliases for autoIndent are: aIndent, aindent.
	`)),
	"splitLines": "Returns a list of strings from the supplied object with newline as the separator.",
	"wrap":       "Wraps the rendered arguments within width.",
}

func (t *Template) addUtilsFuncs() {
	t.AddFunctions(utilsFuncs, utilsBase, funcOptions{
		funcHelp:    utilsFuncsHelp,
		funcArgs:    utilsFuncsArgs,
		funcAliases: utilsFuncsAliases,
	})
}

func glob(args ...interface{}) collections.IGenericList {
	return collections.AsList(utils.GlobFuncTrim(args...))
}

func lorem(funcName interface{}, params ...int) (result string, err error) {
	kind, err := utils.GetLoremKind(fmt.Sprint(funcName))
	if err == nil {
		result, err = utils.Lorem(kind, params...)
	}
	return
}

func center(width interface{}, args ...interface{}) string {
	w := must(strconv.Atoi(fmt.Sprintf("%v", width))).(int)
	return collections.CenterString(fmt.Sprint(args...), w)
}

func wrap(width interface{}, args ...interface{}) string {
	w := must(strconv.Atoi(fmt.Sprintf("%v", width))).(int)
	return collections.WrapString(fmt.Sprint(args...), w)
}

func indent(spacer string, args ...interface{}) string {
	return collections.Indent(fmt.Sprint(args...), spacer)
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

func repeat(n int, a interface{}) (result collections.IGenericList, err error) {
	if n < 0 {
		err = fmt.Errorf("n must be greater or equal than 0")
		return
	}
	result = collections.CreateList(n)
	for i := 0; i < n; i++ {
		result.Set(i, a)
	}
	return
}

func saveToFile(filename string, object interface{}) (string, error) {
	return "", ioutil.WriteFile(filename, []byte(fmt.Sprint(object)), 0644)
}

func raise(format interface{}, args ...interface{}) (string, error) {
	if f := fmt.Sprint(format); strings.Contains(f, "%") {
		return "", fmt.Errorf(f, args...)
	}
	return "", fmt.Errorf(strings.TrimSpace(fmt.Sprintln(append([]interface{}{format}, args...)...)))
}

func assert(test interface{}, args ...interface{}) (string, error) {
	if isZero(test) {
		if len(args) == 0 {
			args = []interface{}{"Assertion failed"}
		}
		return raise(args[0], args[1:]...)
	}
	return "", nil
}
