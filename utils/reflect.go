package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/coveo/gotemplate/errors"
	"github.com/imdario/mergo"
)

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

// IfUndef returns default value if nothing is supplied as values
func IfUndef(def interface{}, values ...interface{}) interface{} {
	switch len(values) {
	case 0:
		return def
	case 1:
		if values[0] == nil {
			return def
		}
		return values[0]
	default:
		return values
	}
}

// IIf acts as a generic ternary operator. It returns valueTrue if testValue is not empty,
// otherwise, it returns valueFalse
func IIf(testValue, valueTrue, valueFalse interface{}) interface{} {
	if IsEmptyValue(reflect.ValueOf(testValue)) {
		return valueFalse
	}
	return valueTrue
}

// Default returns the value if it is not empty or default value.
func Default(value, defaultValue interface{}) interface{} {
	return IIf(value, value, defaultValue)
}

// MergeMaps merges several maps into one privileging the leftmost
func MergeMaps(destination map[string]interface{}, sources ...map[string]interface{}) (map[string]interface{}, error) {
	for i := range sources {
		if err := mergo.Merge(&destination, sources[i]); err != nil {
			return destination, err
		}
	}
	return destination, nil
}

// ToNativeRepresentation converts any object to native (literals, maps, slices)
func ToNativeRepresentation(value interface{}) interface{} {
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
			result[i] = ToNativeRepresentation(reflect.ValueOf(value).Index(i).Interface())
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
					// Ignore embedded fields of unexported non-struct types.
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

// ConvertData returns a go representation of the supplied string (YAML, JSON or HCL)
func ConvertData(data string, out interface{}) (err error) {
	defer func() {
		// YAML converter returns a string if it encounter invalid data, so we
		// check the result to ensure that is is different from the input.
		if out, isItf := out.(*interface{}); isItf && data == fmt.Sprint(*out) {
			err = fmt.Errorf("Unable to find template of file named %s", data)
			*out = nil
		}
	}()

	var errs errors.Array
	if HCLConvert != nil {
		if err = HCLConvert([]byte(data), out); err == nil {
			return
		}
		errs = append(errs, err)
	}

	trySimplified := func() error {
		if strings.Count(data, "=") == 0 {
			return fmt.Errorf("Not simplifiable")
		}
		// Special case where we want to have a map and the supplied string is simplified such as "a = 10 b = string"
		// so we try transform the supplied string in valid YAML
		simplified := regexp.MustCompile(`[ \t]*=[ \t]*`).ReplaceAllString(data, ":")
		simplified = regexp.MustCompile(`[ \t]+`).ReplaceAllString(simplified, "\n")
		simplified = strings.Replace(simplified, ":", ": ", -1) + "\n"
		return YamlUnmarshal([]byte(simplified), out)
	}

	if _, isInterface := out.(*interface{}); isInterface && trySimplified() == nil {
		return nil
	}

	if err := YamlUnmarshal([]byte(data), out); err != nil {
		if _, isMap := out.(*map[string]interface{}); isMap && trySimplified() == nil {
			return nil
		}
		if len(errs) > 0 {
			return append(errs, err)
		}
		return err
	}
	return nil
}

// LoadData returns a go representation of the supplied file name (YAML, JSON or HCL)
func LoadData(filename string, out interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err == nil {
		return ConvertData(string(content), out)
	}
	return
}

// HCLConvert is used to avoid circular reference
var HCLConvert func([]byte, interface{}) error

// ToBash returns the bash 4 variable representation of value
func ToBash(value interface{}) string {
	return toBash(ToNativeRepresentation(value), 0)
}

func toBash(value interface{}, level int) (result string) {
	switch value := value.(type) {
	case string:
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) && !strings.ContainsAny(value, " \t\n[]()") {
			value = value[1 : len(value)-1]
		}
		result = value

	case []interface{}:
		results := ToStrings(value)
		for i := range results {
			results[i] = quote(results[i])
		}
		switch level {
		case 2:
			result = strings.Join(results, ",")
		default:
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		}

	case map[string]interface{}:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		results := make([]string, len(value))
		switch level {
		case 0:
			for i, key := range keys {
				val := toBash(value[key], level+1)
				switch value[key].(type) {
				case []interface{}:
					results[i] = fmt.Sprintf("declare -a %[1]s\n%[1]s=%[2]v", key, val)
				case map[string]interface{}:
					results[i] = fmt.Sprintf("declare -A %[1]s\n%[1]s=%[2]v", key, val)
				default:
					results[i] = fmt.Sprintf("%s=%v", key, val)
				}
			}
			result = strings.Join(results, "\n")
		case 1:
			for i := range keys {
				val := toBash(value[keys[i]], level+1)
				val = strings.Replace(val, `$`, `\$`, -1)
				results[i] = fmt.Sprintf("[%s]=%s", keys[i], val)
			}
			result = fmt.Sprintf("(%s)", strings.Join(results, " "))
		default:
			for i := range keys {
				val := toBash(value[keys[i]], level+1)
				results[i] = fmt.Sprintf("%s=%s", keys[i], quote(val))
			}
			result = strings.Join(results, ",")
		}
	}
	return
}

func quote(s string) string {
	if strings.ContainsAny(s, " \t,[]()") {
		s = fmt.Sprintf("%q", s)
	}
	return s
}
