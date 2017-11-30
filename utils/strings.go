package utils

import (
	"fmt"
	"strings"

	"github.com/Masterminds/sprig"
)

func Interface2string(str interface{}) string {
	switch str := str.(type) {
	case string:
		return str
	default:
		return fmt.Sprintf("%v", str)
	}
}

// Concat returns a string with all string representation of object concatenated without space
func Concat(objects ...interface{}) string {
	var result string
	for _, object := range objects {
		result += fmt.Sprint(object)
	}
	return result
}

var toStrings = sprig.GenericFuncMap()["toStrings"].(func(interface{}) []string)

// SplitLines return a list of interface object for each line in the supplied content
func SplitLines(content interface{}) []interface{} {
	content = Interface2string(content)
	split := strings.Split(content.(string), "\n")
	result := make([]interface{}, len(split))
	for i := range split {
		result[i] = split[i]
	}
	return result
}

// JoinLines concatenate the representation of supplied arguments as a string separated by newlines
func JoinLines(objects ...interface{}) string {
	result := make([]string, len(objects))
	for i := range objects {
		result[i] = fmt.Sprintf("%v", objects[i])
	}
	return strings.Join(result, "\n")
}

// Split2 returns left and right part of a split
func Split2(source, sep string) (left, right string) {
	split := strings.SplitN(source, sep, 2)
	left = split[0]
	if len(split) > 1 {
		right = split[1]
	}
	return
}
