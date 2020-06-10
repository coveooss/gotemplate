package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/utils"
	multicolor "github.com/coveooss/multilogger/color"
)

const (
	utilsBase = "Other utilities"
)

var utilsFuncs = dictionary{
	"center":     center,
	"color":      multicolor.Sprint,
	"colorln":    multicolor.Sprintln,
	"concat":     collections.Concat,
	"formatList": utils.FormatList,
	"id":         id,
	"iif":        collections.IIf,
	"joinLines":  collections.JoinLines,
	"lorem":      lorem,
	"mergeList":  utils.MergeLists,
	"repeat":     repeat,
	"indent":     indent,
	"nIndent":    nIndent,
	"raw":        rawPrint,
	"reCompile":  regexp.Compile,
	"sIndent":    sIndent,
	"splitLines": collections.SplitLines,
	"stripColor": striptColor,
	"wrap":       wrap,
}

var utilsFuncsArgs = arguments{
	"center":     {"width"},
	"formatList": {"format", "list"},
	"id":         {"identifier", "replaceChar"},
	"iif":        {"testValue", "valueTrue", "valueFalse"},
	"joinLines":  {"format"},
	"lorem":      {"loremType", "params"},
	"mergeList":  {"lists"},
	"repeat":     {"n", "element"},
	"indent":     {"nbSpace"},
	"nIndent":    {"nbSpace"},
	"sIndent":    {"spacer"},
	"splitLines": {"content"},
	"wrap":       {"width"},
}

var utilsFuncsAliases = aliases{
	"center":     {"centered"},
	"color":      {"colored", "enhanced"},
	"id":         {"identifier"},
	"iif":        {"ternary"},
	"lorem":      {"loremIpsum"},
	"nIndent":    {"nindent"},
	"formatList": {"autoWrap", "aWrap", "awrap"},
	"raw":        {"printRaw"},
	"sIndent":    {"sindent", "spaceIndent", "autoIndent", "aindent", "aIndent"},
	"stripColor": {"stripansi", "stripANSI", "striptcolor"},
	"wrap":       {"wrapped"},
}

var utilsFuncsHelp = descriptions{
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
		    FgHi: Meaning high intensity forgeround
		    Bg:   Meaning background"
		    BgHi: Meaning high intensity background
	`)),
	"colorln": "Same as color, but using sprintln instead of sprint to format arguments",
	"concat":  "Returns the concatenation (without separator) of the string representation of objects.",
	"formatList": strings.TrimSpace(collections.UnIndent(`
		Return a list of strings by applying the format to each element of the supplied list.

		You can also use autoWrap as Razor expression if you don't want to specify the format.
		The format is then automatically induced by the context around the declaration).
		Valid aliases for autoWrap are: aWrap, awrap.

		Ex:
		    Hello @<autoWrap(to(10)) World!
	`)),
	"id":        "Returns a valid go identifier from the supplied string (replacing any non compliant character by replacement, default _ ).",
	"iif":       "If testValue is empty, returns falseValue, otherwise returns trueValue.\n    WARNING: All arguments are evaluated and must by valid.",
	"indent":    "Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.",
	"joinLines": "Merge the supplied objects into a newline separated string.",
	"lorem":     "Returns a random string. Valid types are be word, words, sentence, para, paragraph, host, email, url.",
	"mergeList": "Return a single list containing all elements from the lists supplied.",
	"nIndent":   "Work as indent but add a newline before.",
	"raw":       "Print the arguments outside of their enclosing quotes",
	"repeat":    "Returns an array with the item repeated n times.",
	"reCompile": "Parses a regular expression and returns Regexp object that can be used to match against text.",
	"sIndent": strings.TrimSpace(collections.UnIndent(`
		Indents the elements using the provided spacer.
		
		You can also use autoIndent as Razor expression if you don't want to specify the spacer.
		Spacer will then be auto determined by the spaces that precede the expression.
		Valid aliases for autoIndent are: aIndent, aindent.
	`)),
	"splitLines": "Returns a list of strings from the supplied object with newline as the separator.",
	"stripColor": "Remove all ANSI colors & attributes from a string.",
	"wrap":       "Wraps the rendered arguments within width.",
}

func (t *Template) addUtilsFuncs() {
	t.AddFunctions(utilsFuncs, utilsBase, FuncOptions{
		FuncHelp:    utilsFuncsHelp,
		FuncArgs:    utilsFuncsArgs,
		FuncAliases: utilsFuncsAliases,
	})
}

func lorem(funcName interface{}, params ...int) (result string, err error) {
	kind, err := utils.GetLoremKind(fmt.Sprint(funcName))
	if err == nil {
		result, err = utils.Lorem(kind, params...)
	}
	return
}

func center(width interface{}, args ...interface{}) (string, error) {
	w, err := strconv.Atoi(fmt.Sprintf("%v", width))
	if err != nil {
		return "", fmt.Errorf("width must be integer")
	}
	return collections.CenterString(multicolor.FormatMessage(args...), w), nil
}

func wrap(width interface{}, args ...interface{}) (string, error) {
	w, err := strconv.Atoi(fmt.Sprintf("%v", width))
	if err != nil {
		return "", fmt.Errorf("width must be integer")
	}
	return collections.WrapString(multicolor.FormatMessage(args...), w), nil
}

func striptColor(value interface{}) string {
	return stripansi.Strip(fmt.Sprint(value))
}

func indent(space int, args ...interface{}) string {
	args = convertArgs(nil, args...).AsArray()
	return collections.Indent(strings.Join(toStrings(args), "\n"), strings.Repeat(" ", space))
}

func nIndent(space int, args ...interface{}) string {
	return "\n" + indent(space, args...)
}

func sIndent(spacer string, args ...interface{}) string {
	args = convertArgs(nil, args...).AsArray()
	return collections.Indent(strings.Join(toStrings(args), "\n"), spacer)
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

func rawPrint(args ...interface{}) interface{} {
	if len(args) <= 1 {
		return fmt.Sprint(args...)
	}
	return utils.FormatList("!Q!%v!Q!", args...)
}
