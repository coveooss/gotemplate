package template

import (
	"fmt"
	"go/parser"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/v3/utils"
	"github.com/fatih/color"
	"github.com/op/go-logging"
)

const (
	protectString = "_=LONG_STRING="
	literalAt     = "_=!AT!=_"
	literalStart  = `{{ "{{" }}`
	stringRep     = "__StRiNg__"
	rangeExpr     = "__RaNgE__"
	defaultExpr   = "__DeFaUlT__"
	funcExpr      = "__FuNc__"
	funcCall      = "__FuNcAlL__"
	typeExpr      = "__TyPe__"
	mapExpr       = "__MaP__"
	structExpr    = "__StRuCt__"
	dotRep        = "__DoT_PrEfIx__"
	ellipsisRep   = "__ElLiPsIs__"
	globalRep     = "__GlObAl__"
)

var dotPrefix = regexp.MustCompile(`(?P<prefix>^|[^\w\)\]])\.(?P<value>\w[\w\.]*)?`)
var idRegex = regexp.MustCompile(`^[\p{L}\d_]+$`)

func expressionParser(repl replacement, match string) string {
	expr, _ := expressionParserInternal(repl, match, false, false)
	return expr
}

func expressionParserSkipError(repl replacement, match string) string {
	expr, _ := expressionParserInternal(repl, match, true, false)
	return expr
}

func expressionParserInternal(repl replacement, match string, skipError, internal bool) (result string, err error) {
	matches, _ := utils.MultiMatch(match, repl.re)
	var expr, expression string
	if expression = matches["expr"]; expression != "" {
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
		for k, v := range map[string]string{"$": stringRep, "range": rangeExpr, "default": defaultExpr, "func": funcExpr, "...": ellipsisRep, "type": typeExpr, "struct": structExpr, "map": mapExpr} {
			protected = protected.Replace(k, v)
		}
		protected = String(dotPrefix.ReplaceAllString(protected.Str(), fmt.Sprintf("${prefix}%s${value}", dotRep)))
		for k, v := range map[string]string{ellipsisRep: "...", "<>": "!=", "÷": "/", "≠": "!=", "≦": "<=", "≧": ">=", "«": "<<", "»": ">>"} {
			protected = protected.Replace(k, v)
		}

		for key, val := range ops {
			protected = protected.Replace(" "+val+" ", key)
		}
		// We add support to partial slice
		protected = String(indexExpression(protected.Str()))

		// We restore the strings into the expression
		expr = protected.RestoreProtected(includedStrings).Str()
	}

	if indexExpr := matches["index"]; indexExpr != "" {
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

	if selectExpr := matches["selection"]; selectExpr != "" {
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
				result = strings.Replace(result, typeExpr, "type", -1)
				result = strings.Replace(result, mapExpr, "map", -1)
				result = strings.Replace(result, structExpr, "struct", -1)
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
