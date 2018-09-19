package collections

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// StringArray is an implementation of a array of class String objects.
type StringArray []String

// Strings convert array of String objects to an array of regular go string.
func (as StringArray) Strings() []string { return ToStrings(as) }

// Str is just an abbreviation of Strings.
func (as StringArray) Str() []string { return as.Strings() }

// Title applies the corresponding String method on each string in the array.
func (as StringArray) Title() StringArray { return as.apply(String.Title) }

// ToTitle applies the corresponding String method on each string in the array.
func (as StringArray) ToTitle() StringArray { return as.apply(String.ToTitle) }

// ToLower applies the corresponding String method on each string in the array.
func (as StringArray) ToLower() StringArray { return as.apply(String.ToLower) }

// ToUpper applies the corresponding String method on each string in the array.
func (as StringArray) ToUpper() StringArray { return as.apply(String.ToUpper) }

// Trim applies the corresponding String method on each string in the array.
func (as StringArray) Trim(cutset string) StringArray { return as.applyStr(String.Trim, cutset) }

// TrimFunc applies the corresponding String method on each string in the array.
func (as StringArray) TrimFunc(f func(rune) bool) StringArray { return as.applyRF(String.TrimFunc, f) }

// TrimLeft applies the corresponding String method on each string in the array.
func (as StringArray) TrimLeft(cutset string) StringArray { return as.applyStr(String.TrimLeft, cutset) }

// TrimLeftFunc applies the corresponding String method on each string in the array.
func (as StringArray) TrimLeftFunc(f func(rune) bool) StringArray {
	return as.applyRF(String.TrimLeftFunc, f)
}

// TrimPrefix applies the corresponding String method on each string in the array.
func (as StringArray) TrimPrefix(prefix string) StringArray {
	return as.applyStr(String.TrimPrefix, prefix)
}

// TrimRight applies the corresponding String method on each string in the array.
func (as StringArray) TrimRight(cutset string) StringArray {
	return as.applyStr(String.TrimRight, cutset)
}

// TrimRightFunc applies the corresponding String method on each string in the array.
func (as StringArray) TrimRightFunc(f func(rune) bool) StringArray {
	return as.applyRF(String.TrimRightFunc, f)
}

// TrimSpace applies the corresponding String method on each string in the array.
func (as StringArray) TrimSpace() StringArray { return as.apply(String.TrimSpace) }

// TrimSuffix applies the corresponding String method on each string in the array.
func (as StringArray) TrimSuffix(suffix string) StringArray {
	return as.applyStr(String.TrimSuffix, suffix)
}

// ToLowerSpecial applies the corresponding String method on each string in the array.
func (as StringArray) ToLowerSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToLowerSpecial, c)
}

// ToTitleSpecial applies the corresponding String method on each string in the array.
func (as StringArray) ToTitleSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToTitleSpecial, c)
}

// ToUpperSpecial applies the corresponding String method on each string in the array.
func (as StringArray) ToUpperSpecial(c unicode.SpecialCase) StringArray {
	return as.applySC(String.ToUpperSpecial, c)
}

// Center applies the corresponding String method on each string in the array.
func (as StringArray) Center(width int) StringArray { return as.applyInt(String.Center, width) }

// Wrap applies the corresponding String method on each string in the array.
func (as StringArray) Wrap(width int) StringArray { return as.applyInt(String.Wrap, width) }

// Indent applies the corresponding String method on each string in the array.
func (as StringArray) Indent(indent string) StringArray { return as.applyStr(String.Indent, indent) }

// IndentN applies the corresponding String method on each string in the array.
func (as StringArray) IndentN(indent int) StringArray { return as.applyInt(String.IndentN, indent) }

// UnIndent applies the corresponding String method on each string in the array.
func (as StringArray) UnIndent() StringArray { return as.apply(String.UnIndent) }

// Join joins all lines in the array using sep as the join string.
func (as StringArray) Join(sep interface{}) String {
	return String(strings.Join(as.Strings(), fmt.Sprint(sep)))
}

// Sorted apply a sort on the array.
func (as StringArray) Sorted() StringArray {
	sorted := as.Strings()
	sort.Strings(sorted)
	return stringArray(sorted)
}

// Utility functions to apply String methods on each element of an array

func (as StringArray) applyRF(f func(String, func(rune) bool) String, fr func(rune) bool) (result StringArray) {
	return as.apply(func(s String) String { return f(s, fr) })
}

func (as StringArray) applySC(f func(String, unicode.SpecialCase) String, c unicode.SpecialCase) (result StringArray) {
	return as.apply(func(s String) String { return f(s, c) })
}

func (as StringArray) applyStr(f func(String, string) String, a string) (result StringArray) {
	return as.apply(func(s String) String { return f(s, a) })
}

func (as StringArray) applyInt(f func(String, int) String, a int) (result StringArray) {
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
