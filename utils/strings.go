package utils

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/Masterminds/sprig"
)

// Interface2string returns the string representation of any interface
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

var ToStrings = sprig.GenericFuncMap()["toStrings"].(func(interface{}) []string)

// SplitLines return a list of interface object for each line in the supplied content
func SplitLines(content interface{}) []interface{} {
	content = Interface2string(content)
	split := strings.Split(strings.TrimSuffix(content.(string), "\n"), "\n")
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

// UnIndent returns the unindented version of the supplied string only if all lines are prefixed
// with the same pattern of spaces
func UnIndent(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= 1 {
		return s
	}

	var spaces *string
	for i, line := range lines {
		if spaces == nil {
			if strings.TrimSpace(line) == "" {
				// We do not consider empty lines
				continue
			}
			trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
			trimmed = string(lines[i][:len(lines[i])-len(trimmed)])
			spaces = &trimmed
		}
		if !strings.HasPrefix(line, *spaces) && strings.TrimSpace(line) != "" {
			return s
		}
		lines[i] = strings.TrimPrefix(line, *spaces)
	}

	return strings.Join(lines, "\n")
}

// Indent returns the indented version of the supplied string
func Indent(s, indent string) string {
	split := strings.Split(s, "\n")
	for i := range split {
		split[i] = indent + split[i]
	}
	return strings.Join(split, "\n")
}

// IndentN returns the indented version (indent as a number of spaces) of the supplied string
func IndentN(s string, indent int) string { return Indent(s, strings.Repeat(" ", indent)) }

// PrettyPrintStruct returns a readable version of an object
func PrettyPrintStruct(object interface{}) string {
	var out string
	isZero := func(x interface{}) bool {
		return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
	}

	val := reflect.ValueOf(object)
	switch val.Kind() {
	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		val = val.Elem()
	}

	result := make([][2]string, 0, val.NumField())
	max := 0
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		if !field.CanInterface() {
			continue
		}

		itf := val.Field(i).Interface()
		if isZero(itf) {
			continue
		}

		itf = reflect.Indirect(val.Field(i)).Interface()
		value := strings.Split(strings.TrimSpace(UnIndent(fmt.Sprint(itf))), "\n")
		if val.Field(i).Kind() == reflect.Struct {
			value[0] = "\n" + IndentN(strings.Join(value, "\n"), 4)
		} else if len(value) > 1 {
			value[0] += " ..."
		}
		result = append(result, [2]string{fieldType.Name, value[0]})
		if len(fieldType.Name) > max {
			max = len(fieldType.Name)
		}
	}

	for _, entry := range result {
		out += fmt.Sprintf("%*s = %v\n", -max, entry[0], entry[1])
	}

	return out
}
