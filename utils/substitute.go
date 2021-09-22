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
		if exprLen < 3 || expression[1] == "" {
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
			errors.Raise("Bad timing information %b, valid values are b(egin) or e(nd)", replacers[i][3])
		}
		result[i].regex = regexp.MustCompile(expression[1])
		result[i].replace = expression[2]

	}
	return result
}

// FilterReplacers will return only replacers that are marked with the right Timing
func filterReplacers(replacers []RegexReplacer, filter string) []RegexReplacer {
	var acc []RegexReplacer
	for _, r := range replacers {
		if r.timing == filter {
			acc = append(acc, r)
		}
	}
	return acc
}

func SubstituteFilteredReplacers(content string, replaceFilter string, replacers ...RegexReplacer) string {
	return Substitute(content, filterReplacers(replacers, replaceFilter)...)
}

// Substitute actually applies the configured substituter
func Substitute(content string, replacers ...RegexReplacer) string {
	// Check timing during replace phase
	for i := range replacers {
		content = replacers[i].regex.ReplaceAllString(content, replacers[i].replace)
	}
	return content
}
