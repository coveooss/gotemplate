package template

import (
	"fmt"
	"strings"

	"github.com/coveooss/multilogger/reutils"
	"github.com/fatih/color"
)

var alreadyIssued = make(map[string]int)

func assignExpression(repl replacement, match string) string {
	return assignExpressionInternal(repl, match, false)
}

func assignExpressionAcceptError(repl replacement, match string) string {
	return assignExpressionInternal(repl, match, true)
}

func assignExpressionInternal(repl replacement, match string, acceptError bool) string {
	matches, _ := reutils.MultiMatch(match, repl.re)
	tp := matches["type"]
	id := matches["id"]
	expr := matches["expr"]
	assign := matches["assign"]
	if tp == "" || id == "" || expr == "" || assign == "" {
		InternalLog.Errorf("Invalid assign regex %s: %s, must contains type, id and expr", repl.name, repl.expr)
		return match
	}

	local := (tp == "$" || tp == "@{" || tp == "@$") && idRegex.MatchString(id)
	var err error
	if expr, err = expressionParserInternal(exprRepl, expr, true, !local); err != nil && !acceptError {
		return match
	}

	if local {
		if assign == "~=" {
			if alreadyIssued[match] == 0 {
				InternalLog.Error("~= assignment is not supported on local variables in", color.HiBlackString(match))
				alreadyIssued[match]++
			}
			return match
		}
		return fmt.Sprintf("%s- $%s %s %s %s", repl.delimiters[0], id, assign, expr, repl.delimiters[1])
	}

	parts := strings.Split(id, ".")
	object := strings.Join(parts[:len(parts)-1], ".")
	id = parts[len(parts)-1]

	if tp == "$" {
		if len(parts) < 2 {
			if alreadyIssued[match] == 0 {
				InternalLog.Errorf("Invalid local assignment: %s", match)
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

	if assign == "~=" {
		// This is a flexible assign, we do not check if the variable already exist (or not)
		return fmt.Sprintf(`%[1]s- set %[3]s "%[4]s" %[5]s %[2]s`, repl.delimiters[0], repl.delimiters[1], object, id, expr)
	}
	validateCode := fmt.Sprintf(map[bool]string{
		true:  `assert (not (isNil %[1]s)) "%[1]s does not exist, use := to declare new variable"`,
		false: `assert (isNil %[1]s) "%[1]s has already been declared, use = to overwrite existing value"`,
	}[assign == "="], fmt.Sprintf("%s%s", iif(strings.HasSuffix(object, "."), object, object+"."), id))

	return fmt.Sprintf(`%[1]s- %[3]s %[2]s%[1]s- set %[4]s "%[5]s" %[6]s %[2]s`, repl.delimiters[0], repl.delimiters[1], validateCode, object, id, expr)
}
