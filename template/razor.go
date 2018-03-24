package template

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

// Add additional functions to the go template context
func (t *Template) applyRazor(content []byte) []byte {
	if !t.options[Razor] || !t.IsRazor(string(content)) {
		return content
	}
	t.ensureInit()

	for _, r := range replacementsInit[fmt.Sprint(t.delimiters)] {
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
	log.Noticef("Generated content\n\n%s\n", color.HiCyanString(strings.Join(lines, "\n")))
	return content
}

var highlight = color.New(color.BgHiBlack, color.FgBlack).SprintFunc()
var iif = utils.IIf
var ifUndef = utils.IfUndef
var defval = utils.Default

// This is indented to simplify the following regular expression for patterns that are repeated several times
// Warning: The declaration order is important
var customMetaclass = [][2]string{
	{"function;", `@reduce;(?P<expr>[id]\([sp][expr]*[sp]\))`},
	{"assign;", `(?P<assign>(?:\$[id][ \t,]*){1,2}:=[sp])`}, // Optional assignment
	{"index;", `(?P<index>\[[expr]+\])`},                    // Extended index operator that support picky selection using ',' as element separator
	{"selector;", `(?P<selection>\.[expr]+)`},               // Optional selector following expression indicating that the expression must include the content after the closing ) (i.e. @function(args).selection)
	{"reduce;", `(?P<reduce>-?)`},                           // Optional reduce sign (-) indicating that the generated code must start with {{-
	{"endexpr;", `(?:[sp];)?`},                              // End expression (spaces + ; or end of line)
	{"[sp]", `[[:blank:]]*`},                                // Optional spaces
	{"[id]", `[\p{L}\d_]+`},                                 // Go language id
	{"[flexible_id]", `[map_id;][map_id;\.\-]*`},            // Id with additional character that could be used to create variables in maps
	{"[idSel]", `[\p{L}_][\p{L}\d_\.]*`},                    // Id with optional selection (object.selection.subselection)
	{"map_id;", `\p{L}\d_\+\*%#!~`},                         // Id with additional character that could be used to create variables in maps
}

// Expression (any character that is not a new line, a start of razor expression or a semicolumn)
var expressionList = []string{`[^\n]`, `[^@\n]`, `[^@;\n]`, `[^@{\n]`, `[^@;{\n]`, fmt.Sprintf(`[^@;{\n'%s]`, "`"), `[\p{L}_\.]`}

const expressionKey = "[expr]"

// Warning: The declaration order is important
var expressions = [][]interface{}{
	// Literals
	{"Protect email", `(\W|^)[\w.!#$%&'*+/=?^_{|}~-]+@[\w-]{1,61}(?:\.[\w-]{1,61})+`, "", replacementFunc(protectEmail)},
	{"", `@@`, literalAt},
	{"", `@{{`, literalStart},
	{"", "(?s)`+.*?`+", "", replacementFunc(protectMultiLineStrings)},
	{"", `@<;`, `{{- $.NEWLINE }}`},
	{"Auto indent", `(?m)^(?P<spaces>.*)@(?:autoIndent|aindent|aIndent)\(`, "@<-sIndent(\"${spaces}\", "},
	{"Newline expression", `@<`, `{{- $.NEWLINE }}@`},

	// Comments
	{"Pseudo line comments - # @", `(?m)(?:^[sp](?:#|//)[sp])@`, "@"},
	{"Pseudo block comments - /*@  @*/", `(?s)/\*@\s*(?P<content>.*?)@\s*\*/`, "${content}"},
	{"Real comments - ##|/// @ comment", `(?m)^[sp](?:##|///)[sp]@.*$`, ""},
	{"Line comment - @// or @#", `(?m)@(#|//)[sp](?P<line_comment>.*)[sp]$`, "{{/* ${line_comment} */}}"},
	{"Block comment - @/* */", `(?s)@/\*(?P<block_comment>.*?)\*/`, "{{/*${block_comment}*/}}"},

	// Commands
	{"Foreach", `@for(?:[sp]each)?[sp]\(`, "@range("},
	{"Single line command - @command (expr) action;", `@(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp](?P<action>[^\n]+?)[sp];`, `{{- ${command} ${assign}${expr} }}${action}{{- end }}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Single line command - @command (expr) { action }", `(?m)@(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]{[sp](?P<action>[^\n]+?)}[sp]$`, `{{- ${command} ${assign}${expr} }}${action}{{- end }}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Command(expr)", `@(?P<command>if|else[sp]if|block|with|define|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]`, `{{- ${command} ${assign}${expr} }}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"else", `@else`, "{{- else }}"},
	{"various ends", `@(?P<command>end[sp](if|range|define|block|with|for[sp]each|for|))endexpr;`, "{{- end }}"},

	// Assignations
	{"Assign local flexible - @$var := value", `(?mU)@\$(?P<id>[id])[sp]:=[sp]?(?P<expr>[expr]+)(?:;|$)`, `{{- $$${id} := ${expr} }}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Assign - @var := value", `(?P<type>@[\$\.]?)(?P<id>[flexible_id])[sp]:=[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},
	{"Assign - $var := value", `(?:{{-?[sp](?:if|range|with)?[sp](\$[id],)?[sp])?(?P<type>\$)(?P<id>[flexible_id])[sp]:=[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},

	// Function calls
	{"Function call followed by expression - @func(args...).args", `function;selector;endexpr;`, `@${reduce}((${expr})${selection});`, replacementFunc(expressionParserSkipError)},
	{"Function call with slice - @func(args...)[...]", `function;index;endexpr;`, `{{${reduce} ${slicer} (${expr}) ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Function call - @func(args...)", `function;endexpr;`, `{{${reduce} ${expr} }}`, replacementFunc(expressionParserSkipError)},
	{"Function unmanaged - @func(value | func)", `@reduce;(?P<function>[id])\([sp](?P<args>[expr]+)[sp]\)endexpr;`, `{{${reduce} ${function} ${args} }}`},

	// Variables
	{"Local variables - @{var}", `@reduce;{[sp](?P<name>[\p{L}\d_\.]*)[sp]}(?P<end>endexpr;)`, `@${reduce}$$${name}${end}`},
	{"Global variables followed by expression", `@reduce;(?P<expr>[idSel]selector;)endexpr;`, `@${reduce}($$.${expr});`, replacementFunc(expressionParserSkipError)},
	{"Global variables with slice - @var[...]", `@reduce;(?P<name>[idSel])index;endexpr;`, `{{${reduce} ${slicer} $$.${name} ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Context variables special with slice", `@reduce;\.(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce} ${slicer} (get . "${name}") ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Global variables special with slice", `@reduce;(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce} ${slicer} (get $$ "${name}") ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Local variables with slice", `@reduce;(?P<expr>(?P<name>[\$\.][\p{L}\d_\.]*)index;)endexpr;`, `{{${reduce} ${slicer} ${name} ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Global variables - @var", `@reduce;(?P<name>[idSel])endexpr;`, `{{${reduce} $$.${name} }}`},
	{"Context variables special - @.var", `@reduce;\.(?P<name>[flexible_id])endexpr;`, `{{${reduce} get . "${name}" }}`},
	{"Global variables special - @var", `@reduce;(?P<name>[flexible_id])endexpr;`, `{{${reduce} get $$ "${name}" }}`},
	{"Local variables - @$var or @.var", `@reduce;(?P<name>[\$\.][\p{L}\d_\.]*)endexpr;`, `{{${reduce} ${name} }}`},

	// Expressions
	{"Expression @(var).selector", `@\([sp](?P<name>[idSel])[sp]\)selector;endexpr;`, `@($$.${name}${selection});`, replacementFunc(expressionParserSkipError)},
	{"Expression @(var)[...]", `@reduce;(?P<expr>\([sp](?P<name>[idSel])[sp]\)index;)endexpr;`, `{{${reduce} ${slicer} $$.${name} ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(var)", `@reduce;\([sp](?P<expr>[idSel])[sp]\)endexpr;`, `{{${reduce} $$.${expr} }}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr).selector", `@\([sp](?P<expr>[expr]+)[sp]\)selector;endexpr;`, `@(${expr}${selection});`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)[...]", `@reduce;\([sp](?P<expr>[expr]+)[sp]\)index;endexpr;`, `{{${reduce} ${slicer} (${expr}) ${index} }}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)", `@reduce;\([sp](?P<expr>[expr]+)[sp]\)endexpr;`, `{{${reduce} ${expr} }}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},

	// Inline contents: Render the content without its enclosing quotes
	{"Inline content", `"<<(?P<content>{{[sp].*[sp]}})"`, `${content}`},

	// Restoring literals
	{"", `}}\\\.`, "}}."},
	{"", literalAt, "@"},
	{"", fmt.Sprintf(`\x60%s(?P<num>\d+)\x60`, protectString), "", replacementFunc(protectMultiLineStrings)},
}

const (
	protectString = "_=LONG_STRING="
	literalAt     = "_=!AT!=_"
	literalStart  = `{{ "{{" }}`
	stringRep     = "_STRING_"
	rangeExpr     = "_range_"
	defaultExpr   = "_default_"
	funcExpr      = "_func_"
	dotRep        = "_DOT_PREFIX_"
	ellipsisRep   = "_ELLIPSIS_"
)

var dotPrefix = regexp.MustCompile(`(?P<prefix>^|[^\w\)\]])\.(?P<value>\w[\w\.]*)?`)
var idRegex = regexp.MustCompile(`^[\p{L}\d_]+$`)

func assignExpression(repl replacement, match string) string {
	if strings.HasPrefix(match, "{{") {
		// This is an already go template assignation
		return match
	}

	subExp := repl.re.SubexpNames()
	subMatches := repl.re.FindStringSubmatch(match)
	tp := valueOf("type", subExp, subMatches)
	id := valueOf("id", subExp, subMatches)
	ex := valueOf("expr", subExp, subMatches)
	if tp == "" || id == "" || ex == "" {
		log.Errorf("Invalid asssign regex %s: %s, must contains type, id and expr", repl.name, repl.expr)
		return match
	}

	local := tp == "$" && idRegex.MatchString(id)
	var err error
	if ex, err = expressionParserInternal(exprRepl, ex, true, !local); err != nil {
		return match
	}

	if local {
		return fmt.Sprintf("{{- $%s := %s }}", id, ex)
	}

	parts := strings.Split(id, ".")
	object := strings.Join(parts[:len(parts)-1], ".")
	id = parts[len(parts)-1]

	if tp == "$" {
		if len(parts) < 2 {
			if alreadyIssued[match] == 0 {
				log.Errorf("Invalid local assignment: %s", match)
				alreadyIssued[match]++
			}
			return match
		}
		object = "$" + object
	} else if strings.HasSuffix(tp, ".") {
		object = "." + object
	} else {
		object = iif(object == "", "$", "$."+object).(string)
	}

	return fmt.Sprintf(`{{- set %s "%s" %s }}`, object, id, ex)
}

var alreadyIssued = make(map[string]int)

func expressionParser(repl replacement, match string) string {
	expr, _ := expressionParserInternal(repl, match, false, false)
	return expr
}

func expressionParserSkipError(repl replacement, match string) string {
	expr, _ := expressionParserInternal(repl, match, true, false)
	return expr
}

func indexOf(name string, names []string) int {
	for i := range names {
		if name == names[i] {
			return i
		}
	}
	return -1
}

func valueOf(name string, names, values []string) string {
	index := indexOf(name, names)
	if index < 0 {
		return ""
	}
	return values[index]
}

func expressionParserInternal(repl replacement, match string, skipError, internal bool) (result string, err error) {
	var expr, expression string
	subNames := repl.re.SubexpNames()
	subMatches := repl.re.FindStringSubmatch(match)
	if expression = valueOf("expr", subNames, subMatches); expression != "" {
		if getLogLevelInternal() >= logging.DEBUG {
			defer func() {
				if !debugMode && result != match {
					log.Debug("Resulting expression =", result)
				}
			}()
		}
		expr = strings.Replace(expression, "$", stringRep, -1)
		expr = strings.Replace(expr, "range", rangeExpr, -1)
		expr = strings.Replace(expr, "default", defaultExpr, -1)
		expr = strings.Replace(expr, "func", funcExpr, -1)
		expr = strings.Replace(expr, "...", ellipsisRep, -1)
		expr = dotPrefix.ReplaceAllString(expr, fmt.Sprintf("${prefix}%s${value}", dotRep))
		expr = strings.Replace(expr, ellipsisRep, "...", -1)
		expr = strings.Replace(expr, "<>", "!=", -1)
		expr = strings.Replace(expr, "÷", "/", -1)
		for key, val := range ops {
			expr = strings.Replace(expr, " "+val+" ", key, -1)
		}
		// We add support to partial slice
		expr = indexExpression(expr)
	}

	if indexExpr := valueOf("index", subNames, subMatches); indexExpr != "" {
		indexExpr = indexExpression(indexExpr)
		indexExpr = indexExpr[1 : len(indexExpr)-1]

		sep, slicer, limit2 := ",", "extract", false
		if strings.Contains(indexExpr, ":") {
			sep, slicer, limit2 = ":", "slice", true
		}
		values := strings.Split(indexExpr, sep)
		if !debugMode && limit2 && len(values) > 2 {
			log.Errorf("Only one : character is allowed in slice expression: %s", match)
		}
		for i := range values {
			if values[i], err = expressionParserInternal(exprRepl, values[i], true, true); err != nil {
				return match, err
			}
		}
		indexExpr = strings.Replace(strings.Join(values, " "), `$`, `$$`, -1)
		repl.replace = strings.Replace(repl.replace, "${index}", indexExpr, -1)
		repl.replace = strings.Replace(repl.replace, "${slicer}", slicer, -1)
	}

	if selectExpr := valueOf("selection", subNames, subMatches); selectExpr != "" {
		if selectExpr, err = expressionParserInternal(exprRepl, selectExpr, true, true); err != nil {
			return match, err
		}
		repl.replace = strings.Replace(repl.replace, "${selection}", selectExpr, -1)
	}

	if expr != "" {
		node := nodeValue
		if internal {
			node = nodeValueInternal
		}
		tr, err := parser.ParseExpr(expr)
		if err == nil {
			result, err := node(tr)
			if err == nil {
				result = strings.Replace(result, stringRep, "$$", -1)
				result = strings.Replace(result, rangeExpr, "range", -1)
				result = strings.Replace(result, defaultExpr, "default", -1)
				result = strings.Replace(result, funcExpr, "func", -1)
				result = strings.Replace(result, dotRep, ".", -1)
				repl.replace = strings.Replace(repl.replace, "${expr}", result, -1)
				return repl.re.ReplaceAllString(match, repl.replace), nil
			}
		}
		if !debugMode && err != nil && getLogLevelInternal() >= 6 {
			log.Debug(color.CyanString(fmt.Sprintf("Invalid expression '%s' : %v", expression, err)))
		}
		if skipError {
			return match, err
		}
		repl.replace = strings.Replace(repl.replace, "${expr}", strings.Replace(expression, "$", "$$", -1), -1)
	}

	return repl.re.ReplaceAllString(match, repl.replace), nil
}

var exprRepl = replacement{
	name:    "Expression",
	re:      regexp.MustCompile(`(?P<expr>.*)`),
	replace: `${expr}`,
}

func indexExpression(expr string) string {
	expr = negativeSlice.ReplaceAllString(expr, "[${index}:0]")
	expr = strings.Replace(expr, "[]", "[0:-1]", -1)
	expr = strings.Replace(expr, "[:", "[0:", -1)
	expr = strings.Replace(expr, ":]", ":-1]", -1)
	return expr
}

var negativeSlice = regexp.MustCompile(`\[(?P<index>-\d+):]`)

func protectMultiLineStrings(repl replacement, match string) string {
	if strings.HasPrefix(match[1:], protectString) {
		// We restore back the long string
		index := errors.Must(strconv.Atoi(repl.re.FindStringSubmatch(match)[1])).(int)
		restore := longStrings[index]
		longStrings[index] = ""
		return restore
	}
	if !strings.Contains(match, "\n") || strings.Contains(match, "``") {
		// We do not have to protect lines without newline or non real multiline string
		return match
	}
	// We save the long string in a buffer, they will be restored at the end of razor preprocessing
	longStrings = append(longStrings, match)
	return fmt.Sprintf("`%s%d`", protectString, len(longStrings)-1)
}

var longStrings []string

func protectEmail(repl replacement, match string) string {
	if match[0] == '@' || match[0] == '#' {
		// This is not an email
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
		var fun string
		if fun, err = nodeValue(n.Fun); err != nil {
			return
		}
		if len(n.Args) == 0 {
			result = fmt.Sprint(fun)
		} else {
			args := make([]string, len(n.Args))
			for i := range n.Args {
				s, err := nodeValueInternal(n.Args[i])
				if err != nil {
					return "", err
				}
				args[i] = s
			}
			result = fmt.Sprintf("%s %s", fun, strings.Join(args, " "))

			if n.Ellipsis != token.NoPos {
				result = fmt.Sprintf("ellipsis %q %s", fun, strings.Join(args, " "))
			}
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
		result = fmt.Sprintf("slice %s %s", x, index)

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
	if !debugMode && getLogLevelInternal() >= 6 {
		log.Debugf(color.HiBlueString("%T => %s"), node, result)
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

var replacementsInit = make(map[string][]replacement)

type replacementFunc func(replacement, string) string
type replacement struct {
	name    string
	expr    string
	replace string
	re      *regexp.Regexp
	parser  replacementFunc
}

func (t *Template) ensureInit() {
	delimiters := fmt.Sprint(t.delimiters)
	if _, ok := replacementsInit[delimiters]; !ok {
		// We must ensure that search and replacement expression are compatible with the set of delimiters
		replacements := make([]replacement, 0, len(expressions))
		for _, expr := range expressions {
			comment := expr[0].(string)
			re := strings.Replace(expr[1].(string), "@", t.delimiters[2], -1)
			replace := strings.Replace(strings.Replace(strings.Replace(expr[2].(string), "{{", t.delimiters[0], -1), "}}", t.delimiters[1], -1), "@", t.delimiters[2], -1)
			var exprParser replacementFunc
			if len(expr) >= 4 {
				exprParser = expr[3].(replacementFunc)
			}

			// We apply replacements in regular expression to make them regex compliant
			for i := range customMetaclass {
				key, value := customMetaclass[i][0], customMetaclass[i][1]
				re = strings.Replace(re, key, value, -1)
			}

			subExpressions := []string{re}
			if strings.Contains(re, expressionKey) {
				// If regex contains the generic expression token [expr], we generate several expression evaluator
				// that go from the most generic expression to the most specific one
				subExpressions = make([]string, len(expressionList))
				for i := range expressionList {
					subExpressions[i] = strings.Replace(re, expressionKey, expressionList[i], -1)
				}
			}

			for i := range subExpressions {
				re := regexp.MustCompile(subExpressions[i])
				replacements = append(replacements, replacement{comment, subExpressions[i], replace, re, exprParser})
			}

			if len(subExpressions) > 1 && len(expr) == 5 {
				// If there is a fallback expression evaluator, we apply it on the first replacement alternative
				re := regexp.MustCompile(subExpressions[0])
				replacements = append(replacements, replacement{comment, subExpressions[0], replace, re, expr[4].(replacementFunc)})
			}
		}
		replacementsInit[delimiters] = replacements
	}
}

func printDebugInfo(r replacement, content string) {
	if r.name == "" || getLogLevelInternal() < logging.INFO {
		return
	}

	debugMode = true
	defer func() { debugMode = false }()

	// We only report each match once
	allUnique := make(map[string]int)
	for _, found := range r.re.FindAllString(content, -1) {
		if r.parser != nil {
			newContent := r.re.ReplaceAllStringFunc(found, func(match string) string {
				return r.parser(r, match)
			})
			if newContent == found && getLogLevelInternal() < 6 {
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

	if len(allUnique) == 0 && getLogLevelInternal() < 7 {
		return
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

	log.Infof("%s: %s%s", color.YellowString(r.name), r.expr, strings.Join(matches, "\n- "))
	if len(matches) > 0 {
		fmt.Fprintln(os.Stderr)
	}
}

var debugMode bool
