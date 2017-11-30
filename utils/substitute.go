package utils

import (
	"regexp"
	"strings"
)

// RegexReplacer defines struct composed of one regular expression and its replacement string
type RegexReplacer struct {
	regex   *regexp.Regexp
	replace string
}

func InitReplacers(replacers ...string) []RegexReplacer {
	result := make([]RegexReplacer, len(replacers))
	for i := range replacers {
		replacers[i] = strings.TrimSpace(replacers[i])
		if replacers[i] == "" {
			continue
		}
		expression := strings.SplitN(replacers[i], string(replacers[i][0]), 3)
		if expression[1] == "" {
			continue
		}
		result[i].regex = regexp.MustCompile(expression[1])
		result[i].replace = expression[2]
	}
	return result
}

func Substitute(content string, replacers ...RegexReplacer) string {
	for i := range replacers {
		content = replacers[i].regex.ReplaceAllString(content, replacers[i].replace)
	}
	return content
}
