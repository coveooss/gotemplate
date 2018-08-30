package template

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
)

var must = errors.Must

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}

// multiMatch returns a map of matching elements from templateErrors list of regular expressions
func multiMatch(s string, expressions ...*regexp.Regexp) (map[string]string, int) {
	for i, re := range expressions {
		if matches := re.FindStringSubmatch(s); len(matches) != 0 {
			results := make(map[string]string)
			for i, key := range re.SubexpNames() {
				if key != "" {
					results[key] = matches[i]
				}
			}
			return results, i
		}
	}
	return nil, -1
}

// getRegexGroup cache compiled regex to avoid multiple interpretation of the same regex
func getRegexGroup(key string, definitions []string) (result []*regexp.Regexp, err error) {
	if result, ok := cachedRegex[key]; ok {
		return result, nil
	}

	result = make([]*regexp.Regexp, len(definitions))
	for i := range definitions {
		regex, err := regexp.Compile(definitions[i])
		if err != nil {
			return nil, err
		}
		result[i] = regex
	}

	cachedRegex[key] = result
	return
}

var cachedRegex = map[string][]*regexp.Regexp{}

// ProtectString returns a string with all included strings replaced by a token and an array of all replaced strings.
// The function is able to detect strings enclosed between quotes "" or backtick `` and it detects escaped characters.
func ProtectString(s string) (result string, array []string) {
	defer func() { result += s }()

	escaped := func(from int) bool {
		// Determine if the previous characters are escaping the current value
		count := 0
		for ; s[from-count] == '\\'; count++ {
		}
		return count%2 == 1
	}

	for end := 0; end >= 0; {
		pos := strings.IndexAny(s, "`\"")
		if pos < 0 {
			break
		}

		for end = strings.IndexRune(s[pos+1:], rune(s[pos])); end >= 0; {
			end += pos + 1
			if s[end] == '"' && escaped(end-1) {
				// The quote has been escaped, so we find the next one
				if newEnd := strings.IndexRune(s[end+1:], '"'); newEnd >= 0 {
					end += newEnd - pos
				}
			} else {
				array = append(array, s[pos:end+1])
				result += s[:pos] + fmt.Sprintf(replacementFormat, len(array)-1)
				s = s[end+1:]
				break
			}
		}
	}
	return
}

// RestoreProtectedString restores a string transformed by ProtectString to its original value.
func RestoreProtectedString(s string, array []string) string {
	return replacementRegex.ReplaceAllStringFunc(s, func(match string) string {
		index := toInt(replacementRegex.FindStringSubmatch(match)[1])
		return array[index]
	})
}

const replacementFormat = `"♠%d"`

var replacementRegex = regexp.MustCompile(`"♠(\d+)"`)
