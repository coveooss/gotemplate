package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/utils"
)

var sprigPick = sprig.GenericFuncMap()["pick"].(func(map[string]interface{}, ...string) map[string]interface{})
var sprigOmit = sprig.GenericFuncMap()["omit"].(func(map[string]interface{}, ...string) map[string]interface{})

func pick(dict map[string]interface{}, keys ...interface{}) map[string]interface{} {
	return sprigPick(dict, utils.ToStrings(convertArgs(nil, keys...))...)
}

func pickv(dict map[string]interface{}, message string, keys ...interface{}) (map[string]interface{}, error) {
	omit := omit(dict, keys...)
	if len(omit) > 0 {
		over := make([]string, 0, len(omit))
		for key := range omit {
			over = append(over, key)
		}
		sort.Strings(over)

		if strings.Contains(message, "%v") {
			message = fmt.Sprintf(message, strings.Join(over, ", "))
		} else {
			message = iif(message == "", "Unwanted values", message).(string)
			message = fmt.Sprintf("%s %s", message, strings.Join(over, ", "))
		}
		return nil, fmt.Errorf(message)
	}
	return pick(dict, keys...), nil
}

func omit(dict map[string]interface{}, keys ...interface{}) map[string]interface{} {
	return sprigOmit(dict, utils.ToStrings(convertArgs(nil, keys...))...)
}
