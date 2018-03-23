package utils

import (
	"fmt"

	"github.com/coveo/gotemplate/types"
)

// MergeLists return a single list from all supplied lists
func MergeLists(lists ...[]interface{}) []interface{} {
	switch len(lists) {
	case 0:
		return nil
	case 1:
		return lists[0]
	}
	result := make([]interface{}, 0)
	for _, list := range lists {
		result = append(result, list...)
	}
	return result
}

// FormatList returns an array of string where format as been applied on every element of the supplied array
func FormatList(format string, v interface{}) []string {
	source := types.ToStrings(v)
	list := make([]string, 0, len(source))
	for _, val := range source {
		list = append(list, fmt.Sprintf(format, val))
	}
	return list
}
