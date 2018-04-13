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

	"github.com/coveo/gotemplate/errors"
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

	if value, err := TryAsList(value); err == nil {
		results := value.Strings()
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

	if value, err := TryAsDictionary(value); err == nil {
		results := make([]string, value.Len())
		vMap := value.AsMap()
		switch level {
		case 0:
			for i, key := range value.KeysAsString() {
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

// ToNativeRepresentation converts any object to native (literals, maps, slices)
func ToNativeRepresentation(value interface{}) (x interface{}) {
	if value == nil {
		return nil
	}

	typ, val := reflect.TypeOf(value), reflect.ValueOf(value)
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = val.Type()
	}
	switch typ.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return fmt.Sprintf("%v", value)

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, val.Len())
		for i := range result {
			result[i] = ToNativeRepresentation(val.Index(i).Interface())
		}
		if len(result) == 1 && reflect.TypeOf(result[0]).Kind() == reflect.Map {
			// If the result is an array of one map, we just return the inner element
			return result[0]
		}
		return result

	case reflect.Map:
		result := make(map[string]interface{}, val.Len())
		for _, key := range val.MapKeys() {
			result[fmt.Sprintf("%v", key)] = ToNativeRepresentation(val.MapIndex(key).Interface())
		}
		return result

	case reflect.Struct:
		result := make(map[string]interface{}, typ.NumField())
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
			tag := sf.Tag.Get("hcl")
			if tag == "" {
				// If there is no hcl specific tag, we rely on json tag if there is
				tag = sf.Tag.Get("json")
			}
			if tag == "-" {
				continue
			}

			split := strings.Split(tag, ",")
			name := split[0]
			if name == "" {
				name = sf.Name
			}
			options := make(map[string]bool, len(split[1:]))
			for i := range split[1:] {
				options[split[i+1]] = true
			}

			if options["omitempty"] && IsEmptyValue(val.Field(i)) {
				continue
			}

			if options["inline"] {
				for key, value := range ToNativeRepresentation(val.Field(i).Interface()).(map[string]interface{}) {
					result[key] = value
				}
			} else {
				result[name] = ToNativeRepresentation(val.Field(i).Interface())
			}
		}
		return result
	default:
		fmt.Fprintf(os.Stderr, "Unknown type %T %v : %v\n", value, typ.Kind(), value)
		return fmt.Sprintf("%v", value)
	}
}

// IsEmptyValue determines if a value is a zero value
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
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
