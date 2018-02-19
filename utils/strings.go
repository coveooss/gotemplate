package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Masterminds/sprig"
)

// -------------------------- String Type --------------------------
type String string

func (s String) Str() string                   { return string(s) }
func (s String) String() string                { return string(s) }
func (s String) Compare(b string) int          { return strings.Compare(string(s), b) }
func (s String) Contains(substr string) bool   { return strings.Contains(string(s), substr) }
func (s String) ContainsAny(chars string) bool { return strings.ContainsAny(string(s), chars) }
func (s String) ContainsRune(r rune) bool      { return strings.ContainsRune(string(s), r) }
func (s String) Count(substr string) int       { return strings.Count(string(s), substr) }
func (s String) EqualFold(t string) bool       { return strings.EqualFold(string(s), t) }
func (s String) Fields() StringArray           { return stringArray(strings.Fields(string(s))) }
func (s String) FieldsFunc(f func(rune) bool) StringArray {
	return stringArray(strings.FieldsFunc(string(s), f))
}
func (s String) HasPrefix(prefix string) bool        { return strings.HasPrefix(string(s), prefix) }
func (s String) HasSuffix(suffix string) bool        { return strings.HasSuffix(string(s), suffix) }
func (s String) Index(substr string) int             { return strings.Index(string(s), substr) }
func (s String) IndexAny(chars string) int           { return strings.IndexAny(string(s), chars) }
func (s String) IndexByte(c byte) int                { return strings.IndexByte(string(s), c) }
func (s String) IndexFunc(f func(rune) bool) int     { return strings.IndexFunc(string(s), f) }
func (s String) IndexRune(r rune) int                { return strings.IndexRune(string(s), r) }
func (s String) Join(array ...string) String         { return String(strings.Join(array, string(s))) }
func (s String) LastIndex(substr string) int         { return strings.LastIndex(string(s), substr) }
func (s String) LastIndexAny(chars string) int       { return strings.LastIndexAny(string(s), chars) }
func (s String) LastIndexByte(c byte) int            { return strings.LastIndexByte(string(s), c) }
func (s String) LastIndexFunc(f func(rune) bool) int { return strings.LastIndexFunc(string(s), f) }
func (s String) Map(mapping func(rune) rune) String  { return String(strings.Map(mapping, string(s))) }
func (s String) Repeat(count int) String             { return String(strings.Repeat(string(s), count)) }
func (s String) Replace(old, new string, n int) String {
	return String(strings.Replace(string(s), old, new, n))
}
func (s String) Split(sep string) StringArray { return stringArray(strings.Split(string(s), sep)) }
func (s String) SplitAfter(sep string) StringArray {
	return stringArray(strings.SplitAfter(string(s), sep))
}
func (s String) SplitAfterN(sep string, n int) StringArray {
	return stringArray(strings.SplitAfterN(string(s), sep, n))
}
func (s String) SplitN(sep string, n int) StringArray {
	return stringArray(strings.SplitN(string(s), sep, n))
}
func (s String) Title() String                     { return String(strings.Title(string(s))) }
func (s String) ToLower() String                   { return String(strings.ToLower(string(s))) }
func (s String) ToTitle() String                   { return String(strings.ToTitle(string(s))) }
func (s String) ToUpper() String                   { return String(strings.ToUpper(string(s))) }
func (s String) Trim(cutset string) String         { return String(strings.Trim(string(s), cutset)) }
func (s String) TrimFunc(f func(rune) bool) String { return String(strings.TrimFunc(string(s), f)) }
func (s String) TrimLeft(cutset string) String     { return String(strings.TrimLeft(string(s), cutset)) }
func (s String) TrimLeftFunc(f func(rune) bool) String {
	return String(strings.TrimLeftFunc(string(s), f))
}
func (s String) TrimPrefix(prefix string) String { return String(strings.TrimPrefix(string(s), prefix)) }
func (s String) TrimRight(cutset string) String  { return String(strings.TrimRight(string(s), cutset)) }
func (s String) TrimRightFunc(f func(rune) bool) String {
	return String(strings.TrimRightFunc(string(s), f))
}
func (s String) TrimSpace() String               { return String(strings.TrimSpace(string(s))) }
func (s String) TrimSuffix(suffix string) String { return String(strings.TrimSuffix(string(s), suffix)) }
func (s String) ToLowerSpecial(c unicode.SpecialCase) String {
	return String(strings.ToLowerSpecial(c, string(s)))
}
func (s String) ToTitleSpecial(c unicode.SpecialCase) String {
	return String(strings.ToTitleSpecial(c, string(s)))
}
func (s String) ToUpperSpecial(c unicode.SpecialCase) String {
	return String(strings.ToUpperSpecial(c, string(s)))
}
func (s String) FieldsId() StringArray {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	return s.FieldsFunc(f)
}
func (s String) Center(width interface{}) String { return String(CenterString(string(s), width)) }

// -------------------------- StringArray Type --------------------------

type StringArray []String

func (as StringArray) Str() []string        { return as.Strings() }
func (as StringArray) Strings() []string    { return ToStrings(as) }
func (as StringArray) Title() StringArray   { return as.apply(String.Title) }
func (as StringArray) ToLower() StringArray { return as.apply(String.ToLower) }
func (as StringArray) ToUpper() StringArray { return as.apply(String.ToUpper) }
func (as StringArray) ToTitle() StringArray { return as.apply(String.ToTitle) }
func (as StringArray) Trim(cutset string) StringArray {
	return as.applyStr(String.Trim, cutset)
}
func (as StringArray) TrimFunc(f func(rune) bool) StringArray {
	return as.applyRF(String.TrimFunc, f)
}
func (as StringArray) TrimLeft(cutset string) StringArray {
	return as.applyStr(String.TrimLeft, cutset)
}
func (as StringArray) TrimLeftFunc(f func(rune) bool) StringArray {
	return as.applyRF(String.TrimLeftFunc, f)
}
func (as StringArray) TrimPrefix(prefix string) StringArray {
	return as.applyStr(String.TrimPrefix, prefix)
}
func (as StringArray) TrimRight(cutset string) StringArray {
	return as.applyStr(String.TrimRight, cutset)
}
func (as StringArray) TrimRightFunc(f func(rune) bool) StringArray {
	return as.applyRF(String.TrimRightFunc, f)
}
func (as StringArray) TrimSpace() StringArray { return as.apply(String.TrimSpace) }
func (as StringArray) TrimSuffix(suffix string) StringArray {
	return as.applyStr(String.TrimSuffix, suffix)
}
func (as StringArray) ToLowerSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToLowerSpecial, c)
}
func (as StringArray) ToTitleSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToTitleSpecial, c)
}
func (as StringArray) ToUpperSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToUpperSpecial, c)
}

func (as StringArray) Join(sep interface{}) String {
	return String(reflect.ValueOf(sep).String()).Join(ToStrings(as)...)
}

func (as StringArray) applyRF(f func(String, func(rune) bool) String, fr func(rune) bool) (result StringArray) {
	return as.apply(func(s String) String { return f(s, fr) })
}

func (as StringArray) applySC(f func(String, unicode.SpecialCase) String, c unicode.SpecialCase) (result StringArray) {
	return as.apply(func(s String) String { return f(s, c) })
}

func (as StringArray) applyStr(f func(String, string) String, a string) (result StringArray) {
	return as.apply(func(s String) String { return f(s, a) })
}

func (as StringArray) apply(f func(s String) String) (result StringArray) {
	result = make(StringArray, len(as))
	for i := range result {
		result[i] = f(as[i])
	}
	return
}

func stringArray(a []string) (result StringArray) {
	result = make(StringArray, len(a))
	for i := range result {
		result[i] = String(a[i])
	}
	return
}

// CenterString returns string s centered within width
func CenterString(s string, width interface{}) string {
	var w int
	w, err := strconv.Atoi(fmt.Sprintf("%v", width))
	l := utf8.RuneCountInString(s)
	if l > w || err != nil {
		return s
	}
	left := (w - l) / 2
	right := w - left - l
	return fmt.Sprintf("%[2]*[1]s%[4]s%[3]*[1]s", "", left, right, s)
}

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

// ToStrings converts the supplied parameter into an array of string
var ToStrings = sprig.GenericFuncMap()["toStrings"].(func(interface{}) []string)

// ToInterfaces converts an array of strings into an array of interfaces
func ToInterfaces(values ...string) []interface{} {
	result := make([]interface{}, len(values))
	for i := range values {
		result[i] = values[i]
	}
	return result
}

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
