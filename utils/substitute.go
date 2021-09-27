package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coveooss/multilogger/errors"
)

// RegexReplacer defines struct composed of one regular expression and its replacement string
type RegexReplacer struct {
	regex   *regexp.Regexp
	replace string
	timing  substituteTiming
}

func (r *RegexReplacer) fromExpressionArray(source []string) RegexReplacer {
	r.timing = extractTiming(source)
	r.regex = regexp.MustCompile(source[1])
	r.replace = source[2]
	return *r
}

// Typed string to represent the different timings for the replacers. Private, so no other values can be instatiated
type substituteTiming string

// Interface to access the private type substituteTiming used to define a simili enum
type SubstituteTiming interface {
	Get() substituteTiming
}

func (sub substituteTiming) Get() substituteTiming {
	return sub
}

const (
	BEGIN    substituteTiming = "b"
	END      substituteTiming = "e"
	_PROTECT substituteTiming = "p"
	NONE     substituteTiming = ""
)

// InitReplacers configures the list of substitution that should be applied on each document
func InitReplacers(replacers ...string) []RegexReplacer {
	result := make([]RegexReplacer, len(replacers))
	// Static allocation would add empty replacers, so we start with no protectors
	var protectors []RegexReplacer
	for i := range replacers {
		replacers[i] = strings.TrimSpace(replacers[i])
		if replacers[i] == "" {
			errors.Raise("Bad replacer %s", replacers[i])
		}
		expression := strings.Split(replacers[i], string(replacers[i][0]))
		exprLen := len(expression)
		if 3 > exprLen || exprLen > 4 || expression[1] == "" {
			errors.Raise("Bad replacer %s", replacers[i])
		}

		if expression[2] == "d" {
			// If the replace expression is a single d (as in delete), we replace the
			// expression by nothing
			if strings.HasSuffix(expression[1], "$") {
				// If we really want to delete lines, we must add \n explicitly
				expression[1] += `\n`
				if !strings.HasPrefix(expression[1], "(?m)") {
					// If the search expression doesn't enable multi line
					// we enable it
					expression[1] = "(?m)" + expression[1]
				}
			}
			expression[2] = ""
		}
		timing := extractTiming(expression)
		if timing == _PROTECT {
			var protectExpr []string
			expression, protectExpr = genProtectExpressions(expression)
			protectors = append(protectors, (&RegexReplacer{}).fromExpressionArray(protectExpr))
		}
		result[i].fromExpressionArray(expression)

	}
	return append(result, protectors...)
}

func extractTiming(expression []string) substituteTiming {
	exprLen := len(expression)
	// the exprLen is repeated, but with boolean algebra magic, we can prove it disapears everytime and makes the program not crash (^,^)
	isValidTiming := exprLen == 4 && strings.Contains("bep", strings.ToLower(expression[3]))
	if exprLen == 4 && isValidTiming {
		return substituteTiming(expression[3])
	} else if exprLen == 4 && !isValidTiming {
		errors.Raise("Bad timing information %s, valid values are b(egin) or e(nd) or p(rotect) for both e.g. /regex/replacer[/b | /e | /p]", expression[3])
	}
	return NONE
}

// generate a "begin" replacer and an "end" replacer and put the extra replacer in the protectors slice.
// It is arbitrary, but place the "end" replacer in the protectors slice.
//
// Return a modified version of the expression passed in params.
func genProtectExpressions(expression []string) (begin []string, end []string) {
	protectionString := fmt.Sprintf("_=!%s!=_", expression[2])
	beginExpr := []string{"", expression[1], protectionString, string(BEGIN)}
	endExpr := []string{"", protectionString, expression[1], string(END)}
	return beginExpr, endExpr
}

// filterReplacers will return only replacers that are marked with the right timing
func filterReplacers(replacers []RegexReplacer, timingFilter substituteTiming) []RegexReplacer {
	acc := make([]RegexReplacer, 0, len(replacers))
	for _, r := range replacers {
		if r.timing == timingFilter.Get() {
			acc = append(acc, r)
		}
	}
	return acc
}

// Substitute actually applies the configured substituter
func Substitute(content string, replacerFilter SubstituteTiming, replacers ...RegexReplacer) string {
	filteredReplacers := filterReplacers(replacers, replacerFilter.Get())
	for i := range filteredReplacers {
		content = filteredReplacers[i].regex.ReplaceAllString(content, filteredReplacers[i].replace)
	}
	return content
}
