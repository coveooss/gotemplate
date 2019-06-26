package template

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"unicode"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/fatih/color"
)

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

	if strings.HasPrefix(result, `"`) || !strings.ContainsAny(result, " \t") {
		return
	}

	// if first, _ := collections.Split2(result, " "); !strings.HasPrefix(first, dotRep) || strings.Contains(first, funcCall) {
	if first, _ := collections.Split2(result, " "); !(strings.HasPrefix(first, dotRep) && strings.Contains(first, funcCall)) {
		result = fmt.Sprintf("(%s)", result)
	}
	return
}
