package collections

import (
	"strings"
	"unicode"
)

// String is an enhanced class implementation of the standard go string library.
// This is convenient when manipulating go template string to have it considered as an object.
type String string

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (s String) Compare(b string) int { return strings.Compare(string(s), b) }

// Contains reports whether substr is within s.
func (s String) Contains(substr string) bool { return strings.Contains(string(s), substr) }

// ContainsAny reports whether any Unicode code points in chars are within s.
func (s String) ContainsAny(chars string) bool { return strings.ContainsAny(string(s), chars) }

// ContainsRune reports whether the Unicode code point r is within s.
func (s String) ContainsRune(r rune) bool { return strings.ContainsRune(string(s), r) }

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func (s String) Count(substr string) int { return strings.Count(string(s), substr) }

// EqualFold reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding.
func (s String) EqualFold(t string) bool { return strings.EqualFold(string(s), t) }

// Fields splits the string s around each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, returning an array of substrings of s or an
// empty list if s contains only white space.
func (s String) Fields() StringArray { return stringArray(strings.Fields(string(s))) }

// FieldsFunc splits the string s at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of s. If all code points in s satisfy f(c) or the
// string is empty, an empty slice is returned.
// FieldsFunc makes no guarantees about the order in which it calls f(c).
// If f does not return consistent results for a given c, FieldsFunc may crash.
func (s String) FieldsFunc(f func(rune) bool) StringArray {
	return stringArray(strings.FieldsFunc(string(s), f))
}

// HasPrefix tests whether the string s begins with prefix.
func (s String) HasPrefix(prefix string) bool { return strings.HasPrefix(string(s), prefix) }

// HasSuffix tests whether the string s ends with suffix.
func (s String) HasSuffix(suffix string) bool { return strings.HasSuffix(string(s), suffix) }

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func (s String) Index(substr string) int { return strings.Index(string(s), substr) }

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
func (s String) IndexAny(chars string) int { return strings.IndexAny(string(s), chars) }

// IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
func (s String) IndexByte(c byte) int { return strings.IndexByte(string(s), c) }

// IndexFunc returns the index into s of the first Unicode code point satisfying f(c), or -1 if none do.
func (s String) IndexFunc(f func(rune) bool) int { return strings.IndexFunc(string(s), f) }

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in s.
// If r is utf8.RuneError, it returns the first instance of any
// invalid UTF-8 byte sequence.
func (s String) IndexRune(r rune) int { return strings.IndexRune(string(s), r) }

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
func (s String) Join(array ...string) String { return String(strings.Join(array, string(s))) }

// LastIndex returns the index of the last instance of substr in s, or -1 if substr is not present in s.
func (s String) LastIndex(substr string) int { return strings.LastIndex(string(s), substr) }

// LastIndexAny returns the index of the last instance of any Unicode code point from chars in s, or -1
// if no Unicode code point from chars is present in s.
func (s String) LastIndexAny(chars string) int { return strings.LastIndexAny(string(s), chars) }

// LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
func (s String) LastIndexByte(c byte) int { return strings.LastIndexByte(string(s), c) }

// LastIndexFunc returns the index into s of the last
// Unicode code point satisfying f(c), or -1 if none do.
func (s String) LastIndexFunc(f func(rune) bool) int { return strings.LastIndexFunc(string(s), f) }

// Lines splits up s into a StringArray using the newline as character separator
func (s String) Lines() StringArray { return s.Split("\n") }

// Map returns a copy of the string s with all its characters modified
// according to the mapping function. If mapping returns a negative value, the character is
// dropped from the string with no replacement.
func (s String) Map(mapping func(rune) rune) String { return String(strings.Map(mapping, string(s))) }

// Repeat returns a new string consisting of count copies of the string s.
//
// It panics if count is negative or if the result of (len(s) * count) overflows.
func (s String) Repeat(count int) String { return String(strings.Repeat(string(s), count)) }

// Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators.
//
// If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
//
// If sep is empty, Split splits after each UTF-8 sequence. If both s and sep are empty, Split returns an empty slice.
//
// It is equivalent to SplitN with a count of -1.
func (s String) Split(sep string) StringArray { return stringArray(strings.Split(string(s), sep)) }

// SplitAfter slices s into all substrings after each instance of sep and returns a slice of those substrings.
//
// If s does not contain sep and sep is not empty, SplitAfter returns a slice of length 1 whose only element is s.
//
// If sep is empty, SplitAfter splits after each UTF-8 sequence. If both s and sep are empty, SplitAfter returns an empty slice.
//
// It is equivalent to SplitAfterN with a count of -1.
func (s String) SplitAfter(sep string) StringArray {
	return stringArray(strings.SplitAfter(string(s), sep))
}

// SplitAfterN slices s into substrings after each instance of sep and returns a slice of those substrings.
//
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings
//
// Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for SplitAfter.
func (s String) SplitAfterN(sep string, n int) StringArray {
	return stringArray(strings.SplitAfterN(string(s), sep, n))
}

// SplitN slices s into substrings separated by sep and returns a slice of the substrings between those separators.
//
// The count determines the number of substrings to return:
//   n > 0: at most n substrings; the last substring will be the unsplit remainder.
//   n == 0: the result is nil (zero substrings)
//   n < 0: all substrings
//
// Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for Split.
func (s String) SplitN(sep string, n int) StringArray {
	return stringArray(strings.SplitN(string(s), sep, n))
}

// Title returns a copy of the string s with all Unicode letters that begin words mapped to their title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
func (s String) Title() String { return String(strings.Title(string(s))) }

// ToLower returns a copy of the string s with all Unicode letters mapped to their lower case.
func (s String) ToLower() String { return String(strings.ToLower(string(s))) }

// ToTitle returns a copy of the string s with all Unicode letters mapped to their title case.
func (s String) ToTitle() String { return String(strings.ToTitle(string(s))) }

// ToUpper returns a copy of the string s with all Unicode letters mapped to their upper case.
func (s String) ToUpper() String { return String(strings.ToUpper(string(s))) }

// Trim returns a slice of the string s with all leading and // trailing Unicode code points contained in cutset removed.
func (s String) Trim(cutset string) String { return String(strings.Trim(string(s), cutset)) }

// TrimFunc returns a slice of the string s with all leading and trailing Unicode code points c satisfying f(c) removed.
func (s String) TrimFunc(f func(rune) bool) String { return String(strings.TrimFunc(string(s), f)) }

// TrimLeft returns a slice of the string s with all leading Unicode code points contained in cutset removed.
func (s String) TrimLeft(cutset string) String { return String(strings.TrimLeft(string(s), cutset)) }

// TrimLeftFunc returns a slice of the string s with all leading Unicode code points c satisfying f(c) removed.
func (s String) TrimLeftFunc(f func(rune) bool) String {
	return String(strings.TrimLeftFunc(string(s), f))
}

// TrimPrefix returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func (s String) TrimPrefix(prefix string) String { return String(strings.TrimPrefix(string(s), prefix)) }

// TrimRight returns a slice of the string s, with all trailing Unicode code points contained in cutset removed.
func (s String) TrimRight(cutset string) String { return String(strings.TrimRight(string(s), cutset)) }

// TrimRightFunc returns a slice of the string s with all trailing Unicode code points c satisfying f(c) removed.
func (s String) TrimRightFunc(f func(rune) bool) String {
	return String(strings.TrimRightFunc(string(s), f))
}

// TrimSpace returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode.
func (s String) TrimSpace() String { return String(strings.TrimSpace(string(s))) }

// TrimSuffix returns s without the provided trailing suffix string.
// If s doesn't end with suffix, s is returned unchanged.
func (s String) TrimSuffix(suffix string) String { return String(strings.TrimSuffix(string(s), suffix)) }

// ToLowerSpecial returns a copy of the string s with all Unicode letters mapped to their
// lower case, giving priority to the special casing rules.
func (s String) ToLowerSpecial(c unicode.SpecialCase) String {
	return String(strings.ToLowerSpecial(c, string(s)))
}

// ToTitleSpecial returns a copy of the string s with all Unicode letters mapped to their
// title case, giving priority to the special casing rules.
func (s String) ToTitleSpecial(c unicode.SpecialCase) String {
	return String(strings.ToTitleSpecial(c, string(s)))
}

// ToUpperSpecial returns a copy of the string s with all Unicode letters mapped to their
// upper case, giving priority to the special casing rules.
func (s String) ToUpperSpecial(c unicode.SpecialCase) String {
	return String(strings.ToUpperSpecial(c, string(s)))
}

// -------------------------------------------------------------------------------------------------------------------
// The following functions are extension or variation of the standard go string library

// String simply convert a String object into a regular string.
func (s String) String() string { return string(s) }

// Str is an abbreviation of String.
func (s String) Str() string { return string(s) }

// Len returns the length of the string.
func (s String) Len() int { return len(s) }

// FieldsID splits the string s at character that is not a valid identifier character (letter, digit or underscore).
func (s String) FieldsID() StringArray {
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}
	return s.FieldsFunc(f)
}

// Center returns the string centered within the specified width.
func (s String) Center(width int) String { return String(CenterString(string(s), width)) }

// Wrap returns the string wrapped with newline when exceeding the specified width.
func (s String) Wrap(width int) String { return String(WrapString(string(s), width)) }

// Replace returns a copy of the string s with the first n non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to
// k+1 replacements for a k-rune string.
func (s String) Replace(old, new string) String {
	return String(strings.Replace(string(s), old, new, -1))
}

// ReplaceN returns a copy of the string s with the first n non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to
// If n < 0, there is no limit on the number of replacements.
func (s String) ReplaceN(old, new string, n int) String {
	return String(strings.Replace(string(s), old, new, n))
}

// Indent returns the indented version of the supplied string (indent represents the string used to indent the lines).
func (s String) Indent(indent string) String {
	return String(Indent(string(s), indent))
}

// IndentN returns the indented version of the supplied string (indent represents the number of spaces used to indent the lines).
func (s String) IndentN(indent int) String {
	return String(IndentN(string(s), indent))
}

// UnIndent returns the string unindented.
func (s String) UnIndent() String {
	return String(UnIndent(string(s)))
}

// GetWordAtPosition returns the selected word and the start position from the specified position.
func (s String) GetWordAtPosition(pos int, accept ...string) (String, int) {
	if pos < 0 || pos >= len(s) {
		return "", -1
	}

	acceptChars := strings.Join(accept, "")
	isBreak := func(c rune) bool {
		return unicode.IsSpace(c) || unicode.IsPunct(c) && !strings.ContainsRune(acceptChars, c)
	}
	begin, end := pos, pos
	for begin >= 0 && !isBreak(rune(s[begin])) {
		begin--
	}
	for end < len(s) && !isBreak(rune(s[end])) {
		end++
	}
	if begin != end {
		begin++
	}
	return s[begin:end], begin
}

// SelectWord returns the selected word from the specified position.
func (s String) SelectWord(pos int, accept ...string) String {
	result, _ := s.GetWordAtPosition(pos, accept...)
	return result
}

// GetContextAtPosition tries to extend the context from the specified position within specified boundaries.
func (s String) GetContextAtPosition(pos int, left, right string) (String, int) {
	if pos >= len(s) {
		pos = len(s) - 1
	}
	if pos < 0 {
		return "", -1
	}
	begin, end := pos, pos+1
	if left != "" {
		begin = strings.LastIndex(s[:pos].Str(), left)
	}
	if right != "" {
		end = strings.Index(s[pos:].Str(), right) + pos + len(right)
	}
	if begin < 0 || end < 0 {
		return "", -1
	}
	return s[begin:end], begin
}

// SelectContext returns the selected word from the specified position.
func (s String) SelectContext(pos int, left, right string) String {
	result, _ := s.GetContextAtPosition(pos, left, right)
	return result
}
