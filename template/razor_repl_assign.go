package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/fatih/color"
)

var alreadyIssued = make(map[string]int)

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

	matches, _ := utils.MultiMatch(match, repl.re)
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
		if strings.HasPrefix(match, "$") {
			// TODO: Deprecated, to remove in future version
			InternalLog.Warningf("$var := value assignation is deprecated, use @{var} := value instead. In: %s", color.HiBlackString(match))
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

	// To avoid breaking change, we issue a warning instead of assertion if the variable has not been declared before being set
	// or declared more than once and the feature flag GOTEMPLATE_DEPRECATED_ASSIGN is not set
	validateFunction := iif(deprecatedAssign, "assert", "assertWarning")
	validateCode := fmt.Sprintf(map[bool]string{
		true:  `%[2]s (not (isNil %[1]s)) "%[1]s does not exist, use := to declare new variable"`,
		false: `%[2]s (isNil %[1]s) "%[1]s has already been declared, use = to overwrite existing value"`,
	}[assign == "="], fmt.Sprintf("%s%s", iif(strings.HasSuffix(object, "."), object, object+"."), id), validateFunction)

	return fmt.Sprintf(`%[1]s- %[3]s %[2]s%[1]s- set %[4]s "%[5]s" %s %[2]s`, repl.delimiters[0], repl.delimiters[1], validateCode, object, id, expr, repl.delimiters[1])
}
