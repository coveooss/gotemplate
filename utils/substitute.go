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
type substituteTiming int

// Interface to access the private type substituteTiming used to define a simili enum
type SubstituteTiming interface {
	Get() substituteTiming
}

func (sub substituteTiming) Get() substituteTiming {
	return sub
}

func (sub substituteTiming) String() string {
	return []string{"b", "e", "p", ""}[sub]
}

func substituteTiming_from_string(value string) substituteTiming {
	switch value {
	case "b":
		return BeginTiming
	case "e":
		return EndTiming
	case "p":
		return _ProtectTiming
	case "":
		return NoTiming
	}
	errors.Raise("Bad timing information %s, valid values are b(egin) or e(nd) or p(rotect) for both e.g. /regex/replacer[/b | /e | /p]", value)
	return NoTiming // this is useless, but the linter was complaining
}

const (
	BeginTiming substituteTiming = iota
	EndTiming
	_ProtectTiming
	NoTiming
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
		if timing == _ProtectTiming {
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
	if exprLen == 4 {
		return substituteTiming_from_string(expression[3])
	}
	return NoTiming
}

// generate a "begin" replacer and an "end" replacer and put the extra replacer in the protectors slice.
// It is arbitrary, but place the "end" replacer in the protectors slice.
//
// Return the beginExpr and the endExpr
func genProtectExpressions(expression []string) (begin []string, end []string) {
	protectedLiteral := regexp.QuoteMeta(expression[1])
	if protectedLiteral != expression[1] {
		errors.Raise("Can't have regex metacharacters in the protected literal: %s", protectedLiteral)
	}
	protectionString := fmt.Sprintf("_=!%s!=_", expression[2])
	beginExpr := []string{"", protectedLiteral, protectionString, BeginTiming.String()}
	endExpr := []string{"", protectionString, protectedLiteral, EndTiming.String()}
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
