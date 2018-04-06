package utils

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/types"
)

// TypeConverters is used to register the available converters
var TypeConverters = make(map[string]func([]byte, interface{}) error)

// ConvertData returns a go representation of the supplied string (YAML, JSON or HCL)
func ConvertData(data string, out interface{}) (err error) {
	trySimplified := func() error {
		if strings.Count(data, "=") == 0 {
			return fmt.Errorf("Not simplifiable")
		}
		// Special case where we want to have a map and the supplied string is simplified such as "a = 10 b = string"
		// so we try transform the supplied string in valid YAML
		simplified := regexp.MustCompile(`[ \t]*=[ \t]*`).ReplaceAllString(data, ":")
		simplified = regexp.MustCompile(`[ \t]+`).ReplaceAllString(simplified, "\n")
		simplified = strings.Replace(simplified, ":", ": ", -1) + "\n"
		return ConvertData(simplified, out)
	}
	var errs errors.Array

	defer func() {
		if err == nil {
			// YAML converter returns a string if it encounter invalid data, so we check the result to ensure that is is different from the input.
			if out, isItf := out.(*interface{}); isItf && data == fmt.Sprint(*out) && strings.ContainsAny(data, "=:{}") {
				if _, isString := (*out).(string); isString {
					if trySimplified() == nil && data != fmt.Sprint(*out) {
						err = nil
						return
					}

					err = errs
					*out = nil
				}
			}
		} else {
			if _, e := types.TryAsList(out); e == nil && trySimplified() == nil {
				err = nil
			}
		}
	}()

	for _, key := range types.AsDictionary(TypeConverters).KeysAsString() {
		err = TypeConverters[key]([]byte(data), out)
		if err == nil {
			return
		}
		errs = append(errs, err)
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errs
	}
}

// LoadData returns a go representation of the supplied file name (YAML, JSON or HCL)
func LoadData(filename string, out interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		return ConvertData(string(content), out)
	}
	return
}

// ToBash returns the bash 4 variable representation of value
func ToBash(value interface{}) string {
	return toBash(ToNativeRepresentation(value), 0)
}

func toBash(value interface{}, level int) (result string) {
	if value, isString := value.(string); isString {
		result = value
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) && !strings.ContainsAny(value, " \t\n[]()") {
			result = value[1 : len(value)-1]
		}
		return
	}

	if value, err := types.TryAsList(value); err == nil {
		results := types.ToStrings(value.AsArray())
		for i := range results {
			results[i] = quote(results[i])
		}
		switch level {
		case 2:
			result = strings.Join(results, ",")
		default:
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		}
		return
	}

	if value, err := types.TryAsDictionary(value); err == nil {
		results := make([]string, value.Len())
		vMap := value.AsMap()
		switch level {
		case 0:
			for i, key := range value.KeysAsString() {
				val := toBash(vMap[key], level+1)
				if _, err := types.TryAsList(vMap[key]); err == nil {
					results[i] = fmt.Sprintf("declare -a %[1]s\n%[1]s=%[2]v", key, val)
				} else if _, err := types.TryAsDictionary(vMap[key]); err == nil {
					results[i] = fmt.Sprintf("declare -A %[1]s\n%[1]s=%[2]v", key, val)
				} else {
					results[i] = fmt.Sprintf("%s=%v", key, val)
				}
			}
			result = strings.Join(results, "\n")
		case 1:
			for i, key := range value.KeysAsString() {
				val := toBash(vMap[key], level+1)
				val = strings.Replace(val, `$`, `\$`, -1)
				results[i] = fmt.Sprintf("[%s]=%s", key, val)
			}
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		default:
			for i, key := range value.KeysAsString() {
				val := toBash(vMap[key], level+1)
				results[i] = fmt.Sprintf("%s=%s", key, quote(val))
			}
			result = strings.Join(results, ",")
		}
		return
	}
	return fmt.Sprint(value)
}

func quote(s string) string {
	if strings.ContainsAny(s, " \t,[]()") {
		s = fmt.Sprintf("%q", s)
	}
	return s
}
