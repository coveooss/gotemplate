package template

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

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
	content = []byte(strings.Replace(string(content), funcCall, "", -1))
	log.Noticef("Generated content\n\n%s\n", color.HiCyanString(String(content).AddLineNumber(0).Str()))
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
	{"assign;", `(?P<assign>(?:\$[id][ \t,]*){1,2}:=[sp])`},       // Optional assignment
	{"index;", `(?P<index>\[[expr]+\])`},                          // Extended index operator that support picky selection using ',' as element separator
	{"selector;", `(?P<selection>\.[expr]+)`},                     // Optional selector following expression indicating that the expression must include the content after the closing ) (i.e. @function(args).selection)
	{"reduce;", `(?P<reduce>((?P<reduce1>-)|_)?(?P<reduce2>-?))`}, // Optional reduces sign (-) indicating that the generated code must start with {{- (and end with -}} if two dashes are specified @--)
	{"endexpr;", `(?:[sp];)?`},                                    // End expression (spaces + ; or end of line)
	{"[sp]", `[[:blank:]]*`},                                      // Optional spaces
	{"[id]", `[\p{L}\d_]+`},                                       // Go language id
	{"[flexible_id]", `[map_id;][map_id;\.\-]*`},                  // Id with additional character that could be used to create variables in maps
	{"[idSel]", `[\p{L}_][\p{L}\d_\.]*`},                          // Id with optional selection (object.selection.subselection)
	{"map_id;", `\p{L}\d_\+\*%#!~`},                               // Id with additional character that could be used to create variables in maps
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
	{`Inline content "@<...>"`, `"@reduce;<(?P<content>.*?)>"`, `"<<@${reduce}(${content})"`},
	{"Newline expression", `@<`, `{{- $.NEWLINE }}@`},

	// Comments
	{"Pseudo line comments - #! @", `(?m)(?:[sp](?:#|//)![sp])@`, "@"},
	{"Pseudo block comments - /*@  @*/", `(?s)/\*@(?P<content>.*?)@\*/`, "${content}"},
	{"Real comments - ##|/// @ comment", `(?m)[sp](?:##|///)[sp]@.*$`, `{{- "" }}`},
	{"Line comment - @// or @#", `(?m)@reduce;(#|//)[sp](?P<line_comment>.*)[sp]$`, "{{${reduce1} /* ${line_comment} */ ${reduce2}}}"},
	{"Block comment - @/* */", `(?s)@reduce;/\*(?P<block_comment>.*?)\*/`, "{{${reduce1} /*${block_comment}*/ ${reduce2}}}"},
	{"", `{{ /\*`, `{{/*`}, {"", `\*/ }}`, `*/}}`}, // Gotemplate is picky about spaces around comment {{- /* comment */ -}} and {{/* comment */}} are valid, but {{-/* comment */-}} and {{ /* comment */ }} are not.

	// Commands
	{"Foreach", `@reduce;for(?:[sp]each)?[sp]\(`, "@${reduce}range("},
	{"Single line command - @command (expr) action;", `@reduce;(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp](?P<action>[^\n]+?)[sp];`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}${action}{{${reduce1} end ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Single line command - @command (expr) { action }", `(?m)@reduce;(?P<command>if|with|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]{[sp](?P<action>[^\n]+?)}[sp]$`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}${action}{{${reduce1} end ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"Command(expr)", `@reduce;(?P<command>if|else[sp]if|block|with|define|range)[sp]\([sp]assign;?[sp](?P<expr>[expr]+)[sp]\)[sp]`, `{{${reduce1} ${command} ${assign}${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},
	{"else", `@reduce;else`, "{{${reduce1} else ${reduce2}}}"},
	{"various ends", `@reduce;(?P<command>end[sp](if|range|define|block|with|for[sp]each|for|))endexpr;`, "{{${reduce1} end ${reduce2}}}"},

	// Assignations
	{"Assign - @var := value", `(?P<type>@[\$\.]?)(?P<id>[flexible_id])[sp](?P<assign>:=|=)[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},
	{"Assign - @{var} := value", `(?P<type>@{)(?P<id>[id])}[sp](?P<assign>:=|=)[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},
	{"Assign - @{var := expr}", `(?P<type>@{)(?P<id>[id])[sp](?P<assign>:=|=)[sp](?P<expr>[expr]+?)}endexpr;`, ``, replacementFunc(assignExpressionAcceptError)},
	// TODO Remove in future version
	{"DEPRECATED Assign - $var := value", `(?:{{-?[sp](?:if|range|with)?[sp](\$[id],)?[sp]|@\(\s*)?(?P<type>\$)(?P<id>[flexible_id])[sp](?P<assign>:=|=)[sp](?P<expr>[expr]+)endexpr;`, ``, replacementFunc(assignExpression)},

	// Function calls
	{"Function call followed by expression - @func(args...).args", `function;selector;endexpr;`, `@${reduce}((${expr})${selection});`, replacementFunc(expressionParserSkipError)},
	{"Function call with slice - @func(args...)[...]", `function;index;endexpr;`, `{{${reduce1} ${slicer} (${expr}) ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Function call - @func(args...)", `function;endexpr;`, `{{${reduce1} ${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Function unmanaged - @func(value | func)", `@reduce;(?P<function>[id])\([sp](?P<args>[expr]+)[sp]\)endexpr;`, `{{${reduce1} ${function} ${args} ${reduce2}}}`},

	// Variables
	{"Local variables - @{var}", `@reduce;{[sp](?P<name>[\p{L}\d_\.]*)[sp]}(?P<end>endexpr;)`, `@${reduce}($$${name});`},
	{"Global variables followed by expression", `@reduce;(?P<expr>[idSel]selector;index;?)(?P<end>endexpr;)`, `@${reduce}(${expr});`, replacementFunc(expressionParserSkipError)},
	{"Context variables - @.var", `@reduce;\.(?P<name>[idSel])endexpr;`, `@${reduce}(.${name})`},
	{"Global variables with slice - @var[...]", `@reduce;(?P<name>[idSel])index;endexpr;`, `{{${reduce1} ${slicer} $$.${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Context variables special with slice", `@reduce;\.(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce1} ${slicer} (get . "${name}") ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Global variables special with slice", `@reduce;(?P<expr>(?P<name>[flexible_id])index;)endexpr;`, `{{${reduce1} ${slicer} (get $$ "${name}") ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Local variables with slice", `@reduce;(?P<expr>(?P<name>[\$\.][\p{L}\d_\.]*)index;)endexpr;`, `{{${reduce1} ${slicer} ${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Global variables - @var", `@reduce;(?P<name>[idSel])endexpr;`, `{{${reduce1} $$.${name} ${reduce2}}}`},
	{"Context variables special - @.var", `@reduce;\.(?P<name>[flexible_id])endexpr;`, `{{${reduce1} get . "${name}" ${reduce2}}}`},
	{"Global variables special - @var", `@reduce;(?P<name>[flexible_id])endexpr;`, `{{${reduce1} get $$ "${name}" ${reduce2}}}`},
	{"Local variables - @$var or @.var", `@reduce;(?P<name>[\$\.][\p{L}\d_\.]*)endexpr;`, `{{${reduce1} ${name} ${reduce2}}}`},

	// Expressions
	{"Expression @(var)[...]", `@reduce;(?P<expr>\([sp](?P<name>[idSel])[sp]\)index;)endexpr;`, `{{${reduce1} ${slicer} $$.${name} ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(var)", `@reduce;\([sp](?P<expr>[idSel])[sp]\)endexpr;`, `{{${reduce1} ${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)[...]", `@reduce;\([sp](?P<expr>[expr]+)[sp]\)index;endexpr;`, `{{${reduce1} ${slicer} (${expr}) ${index} ${reduce2}}}`, replacementFunc(expressionParserSkipError)},
	{"Expression @(expr)", `@reduce;\([sp](?P<assign>.*?:= ?)?[sp](?P<expr>[expr]+)[sp]\)endexpr;`, `{{${reduce1} ${assign}${expr} ${reduce2}}}`, replacementFunc(expressionParserSkipError), replacementFunc(expressionParser)},

	{"Space eater", `@-`, `{{- "" -}}`},

	// Inline contents: Render the content without its enclosing quotes
	{`Inline content "<<..."`, `"<<(?P<content>{{[sp].*[sp]}})"`, `${content}`},

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
	funcCall      = "__FUNCCALL__"
	dotRep        = "_DOT_PREFIX_"
	ellipsisRep   = "_ELLIPSIS_"
	globalRep     = "_GLOBAL_"
)

var dotPrefix = regexp.MustCompile(`(?P<prefix>^|[^\w\)\]])\.(?P<value>\w[\w\.]*)?`)
var idRegex = regexp.MustCompile(`^[\p{L}\d_]+$`)

func assignExpression(repl replacement, match string) string {
	return assignExpressionInternal(repl, match, false)
}

func assignExpressionAcceptError(repl replacement, match string) string {
	return assignExpressionInternal(repl, match, true)
}

// TODO: Deprecated, to remove in future version
var deprecatedAssign = String(os.Getenv(EnvDeprecatedAssign)).ParseBool()

func assignExpressionInternal(repl replacement, match string, acceptError bool) string {
	// TODO: Deprecated, to remove in future version
	if strings.HasPrefix(match, repl.delimiters[0]) || strings.HasPrefix(match, repl.delimiters[2]+"(") {
		// This is an already go template assignation
		return match
	}
	if strings.HasPrefix(match, "$") {
		if deprecatedAssign {
			return match
		}
	}

	subExp := repl.re.SubexpNames()
	subMatches := repl.re.FindStringSubmatch(match)
	tp := valueOf("type", subExp, subMatches)
	id := valueOf("id", subExp, subMatches)
	expr := valueOf("expr", subExp, subMatches)
	assign := valueOf("assign", subExp, subMatches)
	if tp == "" || id == "" || expr == "" || assign == "" {
		log.Errorf("Invalid assign regex %s: %s, must contains type, id and expr", repl.name, repl.expr)
		return match
	}

	local := (tp == "$" || tp == "@{") && idRegex.MatchString(id)
	var err error
	if expr, err = expressionParserInternal(exprRepl, expr, true, !local); err != nil && !acceptError {
		return match
	}

	if local {
		if strings.HasPrefix(match, "$") {
			// TODO: Deprecated, to remove in future version
			Log.Warningf("$var := value assignation is deprecated, use @{var} := value instead. In: %s", color.HiBlackString(match))
		}

		return fmt.Sprintf("%s- $%s %s %s %s", repl.delimiters[0], id, assign, expr, repl.delimiters[1])
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

	// To avoid breaking change, we issue a warning instead of assertion if the variable has not been declared before being set
	// or declared more than once and the feature flag GOTEMPLATE_DEPRECATED_ASSIGN is not set
	validateFunction := iif(deprecatedAssign, "assert", "assertWarning")
	validateCode := fmt.Sprintf(map[bool]string{
		true:  `%[2]s (not (isNil %[1]s)) "%[1]s does not exist, use := to declare new variable"`,
		false: `%[2]s (isNil %[1]s) "%[1]s has already been declared, use = to overwrite existing value"`,
	}[assign == "="], fmt.Sprintf("%s%s", iif(strings.HasSuffix(object, "."), object, object+"."), id), validateFunction)

	return fmt.Sprintf(`%[1]s- %[3]s %[2]s%[1]s- set %[4]s "%[5]s" %s %[2]s`, repl.delimiters[0], repl.delimiters[1], validateCode, object, id, expr, repl.delimiters[1])
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

		// We first protect strings declared in the expression
		protected, includedStrings := String(expression).Protect()

		// We transform the expression into a valid go statement
		protected = protected.Replace("$", stringRep)
		protected = protected.Replace("range", rangeExpr)
		protected = protected.Replace("default", defaultExpr)
		protected = protected.Replace("func", funcExpr)
		protected = protected.Replace("...", ellipsisRep)
		protected = String(dotPrefix.ReplaceAllString(protected.Str(), fmt.Sprintf("${prefix}%s${value}", dotRep)))
		protected = protected.Replace(ellipsisRep, "...")
		protected = protected.Replace("<>", "!=")
		protected = protected.Replace("÷", "/")
		for key, val := range ops {
			protected = protected.Replace(" "+val+" ", key)
		}
		// We add support to partial slice
		protected = String(indexExpression(protected.Str()))

		// We restore the strings into the expression
		expr = protected.RestoreProtected(includedStrings).Str()
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
				result = strings.Replace(result, globalRep, "$$.", -1)
				repl.replace = strings.Replace(repl.replace, "${expr}", result, -1)
				result = repl.re.ReplaceAllString(match, repl.replace)
				result = strings.Replace(result, "$.slice ", "slice $.", -1)
				return result, nil
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
		index := must(strconv.Atoi(repl.re.FindStringSubmatch(match)[1])).(int)
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
			result = iif(unicode.IsDigit(rune(x[0])), "-"+x, fmt.Sprintf("sub 0 %s", x)).(string)
			break
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
		if !strings.HasPrefix(result, dotRep) && !strings.HasPrefix(result, stringRep) && !strings.Contains(result, funcCall) {
			result = globalRep + result
		}
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
		result = fmt.Sprintf("%s.%s", x, strings.TrimPrefix(sel, globalRep))
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
		if !strings.ContainsRune(fun, '.') {
			fun = strings.TrimPrefix(fun, globalRep)
		}
		fun = fun + funcCall
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
	name       string
	expr       string
	replace    string
	re         *regexp.Regexp
	parser     replacementFunc
	delimiters []string
}

func (t *Template) ensureInit() {
	delimiters := fmt.Sprint(t.delimiters)
	if _, ok := replacementsInit[delimiters]; !ok {
		// We must ensure that search and replacement expression are compatible with the set of delimiters
		replacements := make([]replacement, 0, len(expressions))
		for _, expr := range expressions {
			comment := expr[0].(string)
			re := strings.Replace(expr[1].(string), "@", regexp.QuoteMeta(t.delimiters[2]), -1)
			re = strings.Replace(re, "{{", regexp.QuoteMeta(t.delimiters[0]), -1)
			re = strings.Replace(re, "}}", regexp.QuoteMeta(t.delimiters[1]), -1)
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
				replacements = append(replacements, replacement{comment, subExpressions[i], replace, re, exprParser, t.delimiters})
			}

			if len(subExpressions) > 1 && len(expr) == 5 {
				// If there is a fallback expression evaluator, we apply it on the first replacement alternative
				re := regexp.MustCompile(subExpressions[0])
				replacements = append(replacements, replacement{comment, subExpressions[0], replace, re, expr[4].(replacementFunc), t.delimiters})
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
		ErrPrintln()
	}
}

var debugMode bool
