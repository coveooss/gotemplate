package template

import (
	"fmt"
	"strconv"
	"strings"
)

func protectMultiLineStrings(repl replacement, match string) string {
	if strings.HasPrefix(match[1:], protectString) {
		// We restore back the long string
		index := must(strconv.Atoi(repl.re.FindStringSubmatch(match)[1])).(int)
		restore := longStrings[index]
		longStrings[index] = ""
		return restore
	}
	if !strings.Contains(match, "\n") {
		// We do not have to protect lines without newline
		return match
	}
	// We save the long string in a buffer, they will be restored at the end of razor preprocessing
	longStrings = append(longStrings, match)
	return fmt.Sprintf("`%s%d`", protectString, len(longStrings)-1)
}

var longStrings []string

func protectEmail(repl replacement, match string) string {
	if match[0] == '@' || match[0] == '#' {
		// This is not an email
		return match
	}
	return strings.Replace(match, "@", "@@", 1)
}
