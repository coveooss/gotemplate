package collections

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/coveo/gotemplate/v3/errors"
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
			if _, e := TryAsList(out); e == nil && trySimplified() == nil {
				err = nil
			}
		}
	}()

	for _, key := range AsDictionary(TypeConverters).KeysAsString() {
		err = TypeConverters[key.Str()]([]byte(data), out)
		if err == nil {
			return
		}
		errs = append(errs, fmt.Errorf("Trying %s: %v", key, err))
	}

	return errs.AsError()
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
	return toBash(must(MarshalGo(value)), 0)
}

func toBash(value interface{}, level int) (result string) {
	if value, isString := value.(string); isString {
		result = value
		if strings.ContainsAny(value, " \t\n[]()") {
			result = fmt.Sprintf("%q", value)
		}
		return
	}

	if value, err := TryAsList(value); err == nil {
		results := value.Strings()
		for i := range results {
			results[i] = quote(results[i])
		}
		fmt.Println(results)
		switch level {
		case 2:
			result = strings.Join(results, ",")
		default:
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		}
		return
	}

	if value, err := TryAsDictionary(value); err == nil {
		results := make([]string, value.Len())
		vMap := value.AsMap()
		switch level {
		case 0:
			for i, key := range value.KeysAsString() {
				key := key.Str()
				val := toBash(vMap[key], level+1)
				if _, err := TryAsList(vMap[key]); err == nil {
					results[i] = fmt.Sprintf("declare -a %[1]s\n%[1]s=%[2]v", key, val)
				} else if _, err := TryAsDictionary(vMap[key]); err == nil {
					results[i] = fmt.Sprintf("declare -A %[1]s\n%[1]s=%[2]v", key, val)
				} else {
					results[i] = fmt.Sprintf("%s=%v", key, val)
				}
			}
			result = strings.Join(results, "\n")
		case 1:
			for i, key := range value.KeysAsString() {
				key := key.Str()
				val := toBash(vMap[key], level+1)
				val = strings.Replace(val, `$`, `\$`, -1)
				results[i] = fmt.Sprintf("[%s]=%s", key, val)
			}
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		default:
			for i, key := range value.KeysAsString() {
				key := key.Str()
				val := toBash(vMap[key], level+1)
				results[i] = fmt.Sprintf("%s=%s", key, quote(val))
			}
			result = strings.Join(results, ",")
		}
		return
	}
	return fmt.Sprint(value)
}

// GoMarshaler is the interface that could be implemented on object that want to customize
// the marshaling to a native go representation.
type GoMarshaler interface {
	MarshalGo(interface{}) (interface{}, error)
}

// MarshalGo converts any object to native go type (literals, maps, slices).
func MarshalGo(value interface{}) (result interface{}, err error) {
	if value == nil {
		return
	}

	typ, val := reflect.TypeOf(value), reflect.ValueOf(value)
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
		typ = val.Type()
	}

	if val.CanInterface() && val.Type().Implements(reflect.TypeOf((*GoMarshaler)(nil)).Elem()) {
		// The object implement a custom marshaller, so we let it generate its stuff
		return val.Interface().(GoMarshaler).MarshalGo(value)
	}

	switch typ.Kind() {
	case reflect.String:
		result = val.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		result = int(val.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		result = uint(val.Uint())

	case reflect.Int64:
		result = val.Int()

	case reflect.Uint64:
		result = val.Uint()

	case reflect.Float32, reflect.Float64:
		result = val.Float()

	case reflect.Bool:
		result = val.Bool()

	case reflect.Slice, reflect.Array:
		array := make([]interface{}, val.Len())
		for i := range array {
			if array[i], err = MarshalGo(val.Index(i).Interface()); err != nil {
				return
			}
		}
		if len(array) == 1 && reflect.TypeOf(array[0]).Kind() == reflect.Map {
			// If the result is an array of one map, we just return the inner element
			result = array[0]
		} else {
			result = array
		}

	case reflect.Map:
		m := make(map[string]interface{}, val.Len())
		for _, key := range val.MapKeys() {
			if m[fmt.Sprint(key)], err = MarshalGo(val.MapIndex(key).Interface()); err != nil {
				return
			}
		}
		result = m

	case reflect.Struct:
		m := make(map[string]interface{}, typ.NumField())

		info, key, err := getTags(typ)
		if err != nil {
			return nil, err
		}
		for i := 0; i < typ.NumField(); i++ {
			sf := typ.Field(i)
			if sf.Anonymous {
				t := sf.Type
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				// If embedded, StructField.PkgPath is not a reliable
				// indicator of whether the field is exported.
				// See https://golang.org/issue/21122
				if !IsExported(t.Name()) && t.Kind() != reflect.Struct {
					// Ignore embedded fields of unexported non-struct collections.
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
					continue
				}
			} else if sf.PkgPath != "" {
				// Ignore unexported non-embedded fields.
				continue
			}

			tag := info[i]
			if tag.name == "-" || !IsExported(sf.Name) || tag.options["omitempty"] && IsEmptyValue(val.Field(i)) {
				continue
			} else if tag.name == "" {
				tag.name = sf.Name
			}

			if tag.options["inline"] || tag.options["squash"] {
				if subMap, err := MarshalGo(val.Field(i).Interface()); err != nil {
					return nil, err
				} else if subMap, ok := subMap.(map[string]interface{}); ok {
					for key, value := range subMap {
						m[key] = value
					}
				} else {
					return nil, fmt.Errorf("Cannot apply inline or squash to non struct on field '%s'", sf.Name)
				}
			} else {
				if value, err = MarshalGo(val.Field(i).Interface()); err != nil {
					return nil, err
				}
				v := reflect.ValueOf(value)
				if IsEmptyValue(v) && (v.Kind() == reflect.Struct || v.Kind() == reflect.Map) && tag.options["omitempty"] {
					continue
				}
				if key >= 0 {
					m[val.Field(key).String()] = map[string]interface{}{tag.name: v.Interface()}
				} else {
					m[tag.name] = v.Interface()
				}
			}
		}
		result = m
	default:
		fmt.Fprintf(os.Stderr, "Unknown type %T %v : %v\n", value, typ.Kind(), value)
		result = fmt.Sprint(value)
	}
	return
}

type tagInfo struct {
	name    string
	options map[string]bool
}

func getTags(t reflect.Type) (result []tagInfo, key int, err error) {
	key = -1
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		var tag string
		for _, category := range []string{"hcl", "json", "yaml", "xml", "toml"} {
			if value := sf.Tag.Get(category); value != "" {
				tag = value
				break
			}
		}

		split := strings.Split(tag, ",")
		options := make(map[string]bool, len(split[1:]))
		for _, key := range split[1:] {
			options[key] = true
		}
		if options["key"] {
			if key != -1 {
				err = fmt.Errorf("Multiple keys defined on struct '%s' ('%s' and '%s')", t.Name(), t.Field(key).Name, sf.Name)
			}
			key = i
		}
		result = append(result, tagInfo{split[0], options})
	}
	return
}

// IsExported reports whether the identifier is exported.
func IsExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}

func quote(s string) string {
	if strings.ContainsAny(s, " \t,[]()") {
		s = fmt.Sprintf("%q", s)
	}
	return s
}
