package template

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/errors"
)

func safeIndex(value interface{}, index int, def interface{}) (result interface{}, err error) {
	defer func() { err = errors.Trap(err, recover()) }()
	valueOf := reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		if index < 0 || index >= valueOf.Len() {
			return def, nil
		}
		return valueOf.Index(index).Interface(), nil
	default:
		return nil, fmt.Errorf("First argument is not indexable %T", value)
	}
}

func slice(value interface{}, args ...interface{}) (result interface{}, err error) {
	return sliceInternal(value, false, args...)
}

func extract(value interface{}, args ...interface{}) (result interface{}, err error) {
	return sliceInternal(value, true, args...)
}

func sliceInternal(value interface{}, extract bool, args ...interface{}) (result interface{}, err error) {
	defer func() { err = errors.Trap(err, recover()) }()

	args = convertArgs(nil, args...)

	valueOf := reflect.ValueOf(value)
	switch valueOf.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		switch len(args) {
		case 0:
			return valueOf.Interface(), nil
		case 1:
			return selectElement(valueOf, toInt(args[0])), nil
		case 2:
			if !extract {
				return sliceList(valueOf, args...)
			}
			fallthrough
		default:
			if !extract {
				return nil, fmt.Errorf("To many parameters")
			}
			result := make([]interface{}, len(args))
			for i := range args {
				result[i] = selectElement(valueOf, toInt(args[i]))
			}
			if valueOf.Kind() == reflect.String {
				return fmt.Sprint(result...), nil
			}
			return result, nil
		}

	case reflect.Map:
		if !extract {
			return sliceMap(valueOf, args...)
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("Cannot apply slice on type %s", reflect.TypeOf(value))
	}
}

func sliceMap(data reflect.Value, args ...interface{}) (interface{}, error) {
	results := []interface{}{}
	switch len(args) {
	case 0:
		return results, nil
	case 1:
		return data.MapIndex(reflect.ValueOf(args[0])), nil
	case 2:
		keys := data.MapKeys()
		keyStrings := make([]string, len(keys))
		mapStrings := make(map[string]reflect.Value)
		for i := range keys {
			keyStrings[i] = fmt.Sprint(keys[i].Interface())
			mapStrings[keyStrings[i]] = data.MapIndex(keys[i])
		}
		sort.Strings(keyStrings)
		argsStr := toStrings(args)
		for i := range keyStrings {
			if keyStrings[i] >= argsStr[0] && keyStrings[i] <= argsStr[1] {
				results = append(results, mapStrings[keyStrings[i]].Interface())
			}
		}
		return results, nil
	default:
		return nil, nil
	}
}

func sliceList(value reflect.Value, args ...interface{}) (result interface{}, err error) {
	length := value.Len()
	begin := toInt(args[0])
	begin = iif(begin < 0, length+begin+1, begin).(int)
	end := toInt(args[1])
	end = iif(end < 0, length+end+1, end).(int)

	// Check if we should reverse the section
	reverse := end < begin
	if reverse {
		end, begin = begin, end
	}

	// For slice operation, there is no error if the index are of limit
	end = int(min(end, length).(int64))
	begin = int(max(begin, 0).(int64))

	if value.Kind() == reflect.String {
		// String slices are returned as string instead of array of runes
		result := value.String()[begin:end]
		if reverse {
			return reverseString(result), nil
		}
		return result, nil
	}

	if begin > length {
		// Begin is after the end
		return []interface{}{}, nil
	}
	results := make([]interface{}, end-begin)
	for i := range results {
		results[i] = value.Index(i + begin).Interface()
	}
	if reverse {
		return reverseArray(results), nil
	}
	return results, nil
}

func selectElement(value reflect.Value, index int) interface{} {
	index = iif(index < 0, value.Len()+index, index).(int)
	if value.Kind() == reflect.String {
		return value.String()[index : index+1]
	}
	return value.Index(index).Interface()
}

func getSingleMapElement(m interface{}) (key, value interface{}, err error) {
	err = fmt.Errorf("Argument must be a map with a single key")
	if m == nil {
		return
	}
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	switch t.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) != 1 {
			return
		}
		return keys[0].Interface(), v.MapIndex(keys[0]).Interface(), nil
	case reflect.Slice:
		length := v.Len()
		keys := make([]interface{}, length)
		values := make([]interface{}, length)
		for i := range keys {
			if keys[i], values[i], err = getSingleMapElement(v.Index(i).Interface()); err != nil {
				return
			}
		}

		results := make(map[string]interface{})
		for i := range keys {
			results[fmt.Sprint(keys[i])] = values[i]
		}
		return keys, results, nil

	default:
		return
	}
}

var reverseArray = sprig.GenericFuncMap()["reverse"].(func(v interface{}) []interface{})

// Reverse returns its argument string reversed rune-wise left to right.
func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
