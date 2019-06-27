package utils

import (
	"fmt"

	"github.com/coveooss/gotemplate/v3/collections"
)

// MergeLists return a single list from all supplied lists
func MergeLists(lists ...collections.IGenericList) collections.IGenericList {
	switch len(lists) {
	case 0:
		return nil
	case 1:
		return lists[0]
	}
	result := lists[0].Create(0, lists[0].Len()*len(lists))
	for _, list := range lists {
		result = result.Append(list.AsArray()...)
	}
	return result
}

// FormatList returns an array of string where format as been applied on every element of the supplied array
func FormatList(format string, args ...interface{}) collections.IGenericList {
	var list collections.IGenericList
	switch len(args) {
	case 1:
		if l, err := collections.TryAsList(args[0]); err == nil {
			list = l
		}
	default:
		list = collections.AsList(args)
	}
	result := list.Clone()
	for i, value := range list.AsArray() {
		result.Set(i, fmt.Sprintf(format, value))
	}
	return result
}
