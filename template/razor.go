package template

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

// Add additional functions to the go template context
func (t *Template) applyRazor(content []byte) []byte {
	if !bytes.Contains(content, []byte(t.delimiters[2])) {
		return content
	}
	t.ensureInit()

	for _, r := range replacements {
		printDebugInfo(r, string(content))
		if r.parser == nil {
			content = r.re.ReplaceAll(content, []byte(r.replace))
		} else {
			content = r.re.ReplaceAllFunc(content, func(match []byte) []byte {
				return []byte(r.parser(r, string(match)))
			})
		}
	}

	lines := strings.Split(string(content), "\n")
	n := int(math.Log10(float64(len(lines)))) + 1
	for i := range lines {
		lines[i] = fmt.Sprintf("%*d %s", n, i+1, lines[i])
	}
	log.Noticef("Generated content\n%s", color.GreenString(strings.Join(lines, "\n")))
	return content
}

var highlight = color.New(color.BgHiBlack, color.FgBlack).SprintFunc()
var iif = utils.IIf

// This is indented to simplify the following regular expression for patterns that are repeated several times
// Warning: The declaration order is important
var customMetaclass = [][2]string{
	{"function;", `@(?P<expr>[id]\([sp][expr][sp]\))`},
	{"endexpr;", `(?:[sp];)?`},             // Optional end expression (spaces + ;)
	{"[sem]", `(?:[[:blank:]]*;)?`},        // Optional semicolon at the end of a statement
	{"[sp]", `[[:blank:]]*`},               // Optional spaces
	{"[id]", `[\p{L}\d_]+`},                // Go language id
	{"[id2]", `[map_id;][map_id;\.]*`},     // Id with additional character that could be used to create variables in maps
	{"[idSel]", `[\p{L}\d_][\p{L}\d_\.]*`}, // Id with optional selection (object.selection.subselection)
	{"[expr]", `[^@;{\n]*`},                // Expression (any character that is not a new line, a start of razor expression or a semicolumn)
	{"map_id;", `\p{L}\d_\-\+\*%#!~`},      // Id with additional character that could be used to create variables in maps
}

// Warning: The declaration order is important
var expressions = [][]interface{}{
	{"Protect email", `(\W|^)([\w.!#$%&'*+/=?^_{|}~-])+@(\w(?:[\w-]{0,61}[\w])?(?:\.\w(?:[\w-]{0,61}\w))+)`, "", replacementFunc(protectEmail)},
	{"", `@@`, literalAt},
	{"Pseudo line comments - # @", `(?m)(?:^[sp](?:#|//)[sp])@`, "@"},
	{"Pseudo block comments - /*@  @*/", `(?s)/\*@\s*(?P<content>.*?)@\s*\*/`, "${content}"},
	{"Real comments - #|// .* @", `(?m)(?:^[sp](?:#|//).*)@.*$`, ""},
	{"Line comment - @// or @#", `(?m)@(#|//)[sp](?P<line_comment>.*)[sp]$`, "{{/* ${line_comment} */}}"},
	{"Block comment - @/* */", `(?s)@/\*(?P<block_comment>.*?)\*/`, "{{/*${block_comment}*/}}"},
	{"Single line command - @with (expr) action;", `@(?P<command>if|with|range)[sp]\([sp](?P<expr>[expr])[sp]\)[sp](?P<action>[^{\n]*)*[sp];`, `{{- ${command} ${expr} }}${action}{{- end }}`, replacementFunc(expressionParser)},
	{"Single line command - @range (expr) { action }", `(?m)@(?P<command>if|with|range)[sp]\([sp](?P<expr>[expr])[sp]\)[sp]{[sp](?P<action>[^;{\n]*)*}[sp]$`, `{{- ${command} ${expr} }}${action}{{- end }}`, replacementFunc(expressionParser)},
	{"Assign local - @var := value", `@\$(?P<id>[id])[sp]:=[sp](?P<expr>[expr])[sem]`, `{{- $$${id} := ${expr} }}`, replacementFunc(expressionParser)},
	{"Assign context - @.var = value", `@\.(?P<id>[id2])[sp]:=[sp](?P<expr>[expr])[sem]`, `{{- set . "${id}" (${expr}) }}`, replacementFunc(expressionParser)},
	{"Assign global - @var = value", `@(?P<id>[id2])[sp]:=[sp](?P<expr>[expr])[sem]`, `{{- set $ "${id}" (${expr}) }}`, replacementFunc(expressionParser)},
	{"various ends", `@(?P<command>end[sp](if|range|template|define|block|with))[sem]`, "{{- end }}"},
	{"else", `@else[sem]`, "{{- else }}"},
	{"Template", `@template\([sp](?P<args>.+)[sp]\)`, `{{- define ${args} -}}`},
	{"Command - @with (expr)", `@(?P<command>if|elseif|block|with|define|range)[sp]\([sp](?P<expr>[expr])[sp]\)[sp]`, `{{- ${command} ${expr} }}`, replacementFunc(expressionParser)},
	{"Slice function call - @func(args...)[...]", `function;(?P<index>\[[expr]\])endexpr;`, `@(${expr}${index})`},
	{"Function call - @func(args...)", `function;endexpr;`, `{{ ${expr} }}`, replacementFunc(expressionParserSkipError)},
	{"Function unmanaged - @func(value | func)", `@(?P<function>[id])\([sp](?P<args>[expr])[sp]\)[sem]`, `{{ ${function} ${args} }}`},
	{"Global variables - @var", `@(?P<name>[idSel])`, `{{ $$.${name} }}`},
	{"Context variables special - @var", `@\.(?P<name>[id2])`, `{{ get . "${name}" }}`},
	{"Global variables special - @var", `@(?P<name>[id2])`, `{{ get $$ "${name}" }}`},
	{"Local variables - @$var or @.var", `@(?P<name>[\$\.][\p{L}\d_\.]*)`, `{{ ${name} }}`},
	{"Expression var", `@\([sp](?P<name>[idSel])[sp]\)`, `{{ $$.${name} }}`},
	{"Expression", `@\([sp](?P<expr>[expr])[sp]\)[sem]`, `{{ ${expr} }}`, replacementFunc(expressionParser)},
	{"Global content", `@`, `{{ $$ }}`},
	{"Inline content", `"<<(?P<content>{{[sp].*[sp]}})"`, `${content}`},
	{"", literalAt, "@"},
}

const (
	literalAt = "_=!AT!=_"
	stringRep = "_STRING_"
	rangeExpr = "_range_"
	dotRep    = "_DOT_PREFIX_"
)

var dotPrefix = regexp.MustCompile(`(?P<prefix>^|\W)\.(?P<value>\w[\w\.]*)?`)

func expressionParser(repl replacement, match string) string {
	return expressionParserInternal(repl, match, false)
}

func expressionParserSkipError(repl replacement, match string) string {
	return expressionParserInternal(repl, match, true)
}

func findName(name string, values []string) (int, error) {
	for i, value := range values {
		if value == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%s not found in %s", name, values)
}

func expressionParserInternal(repl replacement, match string, skipError bool) (result string) {
	index := errors.Must(findName("expr", repl.re.SubexpNames())).(int)
	expression := repl.re.FindStringSubmatch(match)[index]
	if getLogLevel() >= logging.DEBUG {
		defer func() {
			if result != expression {
				log.Debug("Resulting expression =", result)
			}
		}()
	}

	expr := strings.Replace(expression, "$", stringRep, -1)
	expr = strings.Replace(expr, "range", rangeExpr, -1)
	expr = dotPrefix.ReplaceAllString(expr, fmt.Sprintf("${prefix}%s${value}", dotRep))
	expr = strings.Replace(expr, "<>", "!=", -1)
	expr = strings.Replace(expr, "÷", "/", -1)
	for key, val := range ops {
		expr = strings.Replace(expr, " "+val+" ", key, -1)
	}
	// We add support to partial slice
	expr = negativeSlice.ReplaceAllString(expr, "[${index}:0]")
	expr = strings.Replace(expr, "[:", "[0:", -1)
	expr = strings.Replace(expr, ":]", ":-1]", -1)

	tr, _ := parser.ParseExpr(expr)
	if tr != nil {
		result, err := nodeValue(tr)
		if err == nil {
			result = strings.Replace(result, stringRep, "$$", -1)
			result = strings.Replace(result, rangeExpr, "range", -1)
			result = strings.Replace(result, dotRep, ".", -1)
			return repl.re.ReplaceAllString(match, strings.Replace(repl.replace, "${expr}", result, -1))
		}
		log.Debug(color.CyanString(fmt.Sprintf("Invalid expression '%s' : %v", expression, err)))
	} else {
		log.Debug(color.CyanString(fmt.Sprintf("Invalid expression '%s'", expression)))
	}
	if skipError {
		return match
	}
	return repl.re.ReplaceAllString(match, strings.Replace(repl.replace, "${expr}", strings.Replace(expression, "$", "$$", -1), -1))
}

var negativeSlice = regexp.MustCompile(`\[(?P<index>-\d+):]`)

func protectEmail(repl replacement, match string) string {
	if match[0] == '@' {
		// If the first character is a @, this is not an email
		return match
	}
	return strings.Replace(match, "@", "@@", 1)
}

func nodeValue(node ast.Node) (result string, err error) {
	switch n := node.(type) {
	case *ast.UnaryExpr:
		var op, x string
		if op, err = opName(n.Op); err != nil {
			return
		}
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		if op == "sub" {
			result = iif(x[0] != '(', "-"+x, fmt.Sprintf("sub 0 %s", x)).(string)
			return
		}
		result = fmt.Sprintf("%s %s", op, x)
	case *ast.BinaryExpr:
		var op, x, y string
		if op, err = opName(n.Op); err != nil {
			return
		}
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		if y, err = nodeValueInternal(n.Y); err != nil {
			return
		}
		if op == "mul" && strings.Contains(y, "*") {
			// This is a special case where the expression contains 2 following *, meaning power instead of mul
			result = fmt.Sprintf("power %s %s", x, strings.Replace(y, "*", "", -1))
			return
		}
		result = fmt.Sprintf("%s %s %s", op, x, y)
	case *ast.Ident:
		result = n.Name
	case *ast.BasicLit:
		result = fmt.Sprint(n.Value)
	case *ast.SelectorExpr:
		var x, sel string
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		if sel, err = nodeValueInternal(n.Sel); err != nil {
			return
		}
		result = fmt.Sprintf("%s.%s", x, sel)
	case *ast.ParenExpr:
		var content string
		if content, err = nodeValue(n.X); err == nil {
			result = content
		}
	case *ast.CallExpr:
		if len(n.Args) == 0 {
			result = fmt.Sprint(n.Fun)
		} else {
			args := make([]string, len(n.Args))
			for i := range n.Args {
				s, err := nodeValueInternal(n.Args[i])
				if err != nil {
					return "", err
				}
				args[i] = s
			}
			result = fmt.Sprintf("%s %s", n.Fun, strings.Join(args, " "))
		}
	case *ast.StarExpr:
		var x string
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		// This is a special case where the expression contains 2 following *, meaning power instead of mul
		result = fmt.Sprintf("*%s", x)

	case *ast.IndexExpr:
		var x, index string
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		if index, err = nodeValueInternal(n.Index); err != nil {
			return
		}
		result = fmt.Sprintf("index %s %s", x, index)

	case *ast.SliceExpr:
		var x, low, high string
		if x, err = nodeValueInternal(n.X); err != nil {
			return
		}
		if low, err = nodeValueInternal(n.Low); err != nil {
			return
		}
		if high, err = nodeValueInternal(n.High); err != nil {
			return
		}
		result = fmt.Sprintf("slice %s %s %s", x, low, high)

	default:
		err = fmt.Errorf("Unknown: %v", reflect.TypeOf(node))
	}
	return
}

var ops = map[string]string{
	"==": "eq",
	"!=": "ne",
	"<":  "lt",
	"<=": "le",
	">":  "gt",
	">=": "ge",
	"+":  "add",
	"-":  "sub",
	"/":  "div",
	"*":  "mul",
	"%":  "mod",
	"||": "or",
	"&&": "and",
	"!":  "not",
	"<<": "lshift",
	">>": "rshift",
	"|":  "bor",
	"&":  "band",
	"^":  "bxor",
	"&^": "bclear",
}

func opName(token token.Token) (string, error) {
	if name, ok := ops[token.String()]; ok {
		return name, nil
	}
	return "", fmt.Errorf("Unknown operator %v", token)
}

func nodeValueInternal(node ast.Node) (result string, err error) {
	result, err = nodeValue(node)
	if !strings.HasPrefix(result, "\"") && strings.ContainsAny(result, " \t") {
		result = fmt.Sprintf("(%s)", result)
	}
	return
}

var replacements []replacement

type replacementFunc func(replacement, string) string
type replacement struct {
	name    string
	expr    string
	replace string
	re      *regexp.Regexp
	parser  replacementFunc
}

func (t *Template) ensureInit() {
	if replacements == nil {
		replacements = make([]replacement, 0, len(expressions))
		for _, expr := range expressions {
			var exprParser replacementFunc
			if len(expr) == 4 {
				exprParser = expr[3].(replacementFunc)
			}

			// We apply replacements in regular expression to make them regex compliant
			for i := range customMetaclass {
				key, value := customMetaclass[i][0], customMetaclass[i][1]
				expr[1] = strings.Replace(expr[1].(string), key, value, -1)
			}

			replacements = append(replacements, replacement{
				expr[0].(string),
				expr[1].(string),
				strings.Replace(strings.Replace(strings.Replace(expr[2].(string), "{{", t.delimiters[0], -1), "}}", t.delimiters[1], -1), "@", t.delimiters[2], -1),
				regexp.MustCompile(strings.Replace(expr[1].(string), "@", t.delimiters[2], -1)),
				exprParser,
			})
		}
	}
}

func printDebugInfo(r replacement, content string) {
	if r.name == "" || getLogLevel() < logging.INFO {
		return
	}

	// We only report each match once
	allUnique := make(map[string]int)
	for _, found := range r.re.FindAllString(content, -1) {
		if r.parser != nil {
			newContent := r.re.ReplaceAllStringFunc(found, func(match string) string {
				return r.parser(r, match)
			})
			if newContent == found {
				// There is no change
				continue
			}
		}
		lines := strings.Split(found, "\n")
		for i := range lines {
			lines[i] = fmt.Sprintf("%v%s", iif(i > 0, "  ", ""), highlight(lines[i]))
		}
		found = strings.Join(lines, "\n")
		allUnique[found] = allUnique[found] + 1
	}

	matches := make([]string, 0, len(allUnique)+1)
	if len(allUnique) > 0 {
		matches = append(matches, "")
	}
	for key, count := range allUnique {
		if count > 1 {
			key = fmt.Sprintf("%s (%d)", key, count)
		}
		matches = append(matches, key)
	}
	fmt.Fprintln(os.Stderr)
	log.Infof("%s: %s%s", color.YellowString(r.name), r.expr, strings.Join(matches, "\n- "))
}
