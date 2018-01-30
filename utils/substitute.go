package utils

import (
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/errors"
)

// RegexReplacer defines struct composed of one regular expression and its replacement string
type RegexReplacer struct {
	regex   *regexp.Regexp
	replace string
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
		if len(expression) != 3 || expression[1] == "" {
			errors.Raise("Bad replacer %s", replacers[i])
		}
		result[i].regex = regexp.MustCompile(expression[1])
		result[i].replace = expression[2]
	}
	return result
}

// Substitute actually applies the configured substituter
func Substitute(content string, replacers ...RegexReplacer) string {
	for i := range replacers {
		content = replacers[i].regex.ReplaceAllString(content, replacers[i].replace)
	}
	return content
}
