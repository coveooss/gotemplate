package template

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/op/go-logging"
)

// Add additional functions to the go template context
func (t *Template) applyRazor(content []byte) []byte {
	if !bytes.Contains(content, []byte(t.delimiters[2])) {
		return content
	}
	t.ensureInit()
	for _, e := range replacements {
		if getLogLevel() >= logging.INFO && e.name != "" {
			all := e.re.FindAllString(string(content), -1)
			if len(all) > 0 {
				log.Infof("%s = %s", color.YellowString(e.name), strings.Join(all, " | "))
			}
		}
		if e.parser == nil {
			content = e.re.ReplaceAll(content, []byte(e.replace))
		} else {
			content = e.re.ReplaceAllFunc(content, func(match []byte) []byte {
				return []byte(e.parser(e, string(match)))
			})
		}
	}
	log.Noticef("Generated content\n%s", color.GreenString(string(content)))
	return content
}

var expressions = [][]interface{}{
	{"Protect email", `(\W|^)([\w.!#$%&'*+/=?^_{|}~-])+@(\w(?:[\w-]{0,61}[\w])?(?:\.\w(?:[\w-]{0,61}\w))+)`, "", replacementFunc(protectEmail)},
	{"", `@@`, literalAt},
	{"Pseudo comments", `(?:(?:#|//)\s*)@`, "@"},
	{"Line comment", `(?m)@(#|//)\s*(?P<line_comment>.*)\s*$`, "{{/* ${line_comment} */}}"},
	{"Block comment", `(?s)@/\*(?P<block_comment>.*?)\*/`, "{{/*${block_comment}*/}}"},
	{"else if", `@(?P<command>elseif)\(\s*(?P<arg>.*?)\s*\)`, "{{- else if ${arg} }}"},
	{"if", `@if\(\s*(?P<expr>[^@{]*)\s*\)`, "{{- if ${expr} }}", replacementFunc(expressionParser)},
	{"various ends", `@(?P<command>end(if|range|template|define|block|with))`, "{{- end }}"},
	{"else", `@(?P<command>else)`, "{{- else }}"},
	{"Assign local", `(?m)@\$(?P<id>\w[\w-]*)\s*:=\s*(?P<value>.*)\s*$`, `{{- $$${id} := ${value} }}`},
	{"Assing global", `(?m)@(?P<id>\w[\w-]*)\s*:=\s*(?P<value>.*)\s*$`, `{{- set $ "${id}" (${value}) }}`},
	{"Range", `(?m)@range\((?P<args>.+)\).*@endrange`, ``},
	{"Template", `@template\(\s*(?P<args>.+)\w*\)`, `{{- define ${args} -}}`},
	{"Section", `@(?P<id>block|with|define)\(\s*(?P<args>.+)\w*\)`, `{{- ${id} ${args} }}`},
	{"Function call", `@(?P<function>\w+)\(\s*(?P<args>.+)\w*\)`, `{{ ${function} ${args} }}`},
	{"Global variables", `@(?P<name>\w[\w\.]*)`, `{{ $$.${name} }}`},
	{"Local variables", `@(?P<name>[\$\.][\w\.]*)`, `{{ ${name} }}`},
	{"Expresion var", `@\(\s*(?P<expr>\w[\w\.]*)\s*\)`, `{{ $$.${expr} }}`},
	{"Expresion", `@\(\s*(?P<expr>[^@{]*)\s*\)`, `{{ ${expr} }}`, replacementFunc(expressionParser)},
	{"Global content", `@`, `{{ $$ }}`},
	{"Inline content", `"<<(?P<content>{{\s*.*\s*}})"`, `${content}`},
	{"", literalAt, "@"},
}

const (
	literalAt = "#!AT#!"
	stringRep = "_STRING_"
	dotRep    = "_DOT_PREFIX_"
)

var dotPrefix = regexp.MustCompile(`(?P<prefix>^|\W)\.(?P<value>\w[\w\.]*)?`)

func expressionParser(repl replacement, match string) string {
	expression := repl.re.FindStringSubmatch(match)[1]
	expr := strings.Replace(expression, "$", stringRep, -1)
	expr = dotPrefix.ReplaceAllString(expr, fmt.Sprintf("${prefix}%s${value}", dotRep))
	expr = strings.Replace(expr, "<>", "!=", -1)
	expr = strings.Replace(expr, " and ", "&&", -1)
	expr = strings.Replace(expr, " or ", "||", -1)
	tr, _ := parser.ParseExpr(expr)
	if tr != nil {
		result, err := nodeValue(tr)
		if err == nil {
			result := repl.re.ReplaceAllString(match, strings.Replace(repl.replace, "${expr}", result, -1))
			result = strings.Replace(result, stringRep, "$", -1)
			result = strings.Replace(result, dotRep, ".", -1)
			return result
		}
		log.Debug(color.YellowString(fmt.Sprintf("Invalid expression '%s' : %v", expression, err)))
	} else {
		log.Debug(color.YellowString(fmt.Sprintf("Invalid expression '%s'", expression)))
	}
	return repl.re.ReplaceAllString(match, strings.Replace(repl.replace, "${expr}", expression, -1))
}

func protectEmail(repl replacement, match string) string {
	if match[0] == '@' {
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
		var x string
		if x, err = nodeValue(n.X); err != nil {
			return
		}
		result = x
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
}

func opName(token token.Token) (string, error) {
	if name, ok := ops[token.String()]; ok {
		return name, nil
	}
	return "", fmt.Errorf("Unknown operator %v", token)
}

func nodeValueInternal(node ast.Node) (result string, err error) {
	result, err = nodeValue(node)
	if strings.ContainsAny(result, " \t") {
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
