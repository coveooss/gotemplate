package template

import "regexp"

// Add additional functions to the go template context
func (t *Template) applyRazor(content []byte) []byte {
	variables := regexp.MustCompile(`@[\w\.]+`)

	for _, match := range variables.FindAllIndex(content, -1) {
		start, length := match[0], match[1]
		content = append(content[:start], content[start+length:]...)
	}
	return content
}
