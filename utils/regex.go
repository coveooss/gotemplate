package utils

import "regexp"

// MultiMatch returns a map of matching elements from a list of regular expressions (returning the first matching element)
func MultiMatch(s string, expressions ...*regexp.Regexp) (map[string]string, int) {
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

// GetRegexGroup cache compiled regex to avoid multiple interpretation of the same regex
func GetRegexGroup(key string, definitions []string) (result []*regexp.Regexp, err error) {
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
