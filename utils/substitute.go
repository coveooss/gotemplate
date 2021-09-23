package utils

import (
	"regexp"
	"strings"

	"github.com/coveooss/multilogger/errors"
)

// RegexReplacer defines struct composed of one regular expression and its replacement string
type RegexReplacer struct {
	regex   *regexp.Regexp
	replace string
	timing  string
}

// InitReplacers configures the list of substitution that should be applied on each document
func InitReplacers(replacers ...string) []RegexReplacer {
	result := make([]RegexReplacer, len(replacers))
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
		// last part is optional, but specifies the timing. By default, we assume it's not important
		if exprLen == 4 && (strings.Contains("be", strings.ToLower(expression[3]))) {
			result[i].timing = expression[3]
		} else if exprLen == 4 && !strings.Contains("be", strings.ToLower(expression[3])) {
			errors.Raise("Bad timing information %b, valid values are b(egin) or e(nd). e.g. /regex/replacer[/b | /e]", replacers[i][3])
		}
		result[i].regex = regexp.MustCompile(expression[1])
		result[i].replace = expression[2]

	}
	return result
}

// FilterReplacers will return only replacers that are marked with the right timing
func filterReplacers(replacers []RegexReplacer, timingFilter string) []RegexReplacer {
	var acc []RegexReplacer
	for _, r := range replacers {
		if r.timing == timingFilter {
			acc = append(acc, r)
		}
	}
	return acc
}

// Substitute actually applies the configured substituter
func Substitute(content string, replacerFilter string, replacers ...RegexReplacer) string {
	filteredReplacers := filterReplacers(replacers, replacerFilter)
	for i := range filteredReplacers {
		content = filteredReplacers[i].regex.ReplaceAllString(content, filteredReplacers[i].replace)
	}
	return content
}
