package template

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/sprig"
	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/utils"
)

const (
	dataBase       = "Data Manipulation"
	dataConversion = "Data Conversion"
)

var dataFuncsBase = map[string]interface{}{
	"array":     array,
	"bool":      strconv.ParseBool,
	"char":      toChar,
	"content":   content,
	"extract":   extract,
	"get":       get,
	"key":       key,
	"lenc":      utf8.RuneCountInString,
	"merge":     utils.MergeMaps,
	"omit":      omit,
	"pick":      pick,
	"pickv":     pickv,
	"safeIndex": safeIndex,
	"set":       set,
	"slice":     slice,
	"string":    toString,
	"String":    toStringClass,
	"undef":     utils.IfUndef,
}

var dataFuncsConversion = map[string]interface{}{
	"toBash":         utils.ToBash,
	"toHcl":          toHCL,
	"toInternalHcl":  toInternalHCL,
	"toJson":         toJSON,
	"toPrettyHcl":    toPrettyHCL,
	"toPrettyJson":   toPrettyJSON,
	"toPrettyTFVars": toPrettyTFVars,
	"toQuotedHcl":    toQuotedHCL,
	"toQuotedJson":   toQuotedJSON,
	"toQuotedTFVars": toQuotedTFVars,
	"toTFVars":       toTFVars,
	"toYaml":         utils.ToYaml,
}

var dataFuncsArgs = map[string][]string{
	"array":          {"value"},
	"bool":           {"str"},
	"char":           {"value"},
	"content":        {"keymap"},
	"data":           {"data", "context"},
	"extract":        {"source", "indexes"},
	"get":            {"map", "key"},
	"hcl":            {"hcl", "context"},
	"json":           {"json", "context"},
	"key":            {"value"},
	"lenc":           {"str"},
	"merge":          {"destination", "sources"},
	"omit":           {"dict", "keys"},
	"pick":           {"dict", "keys"},
	"pickv":          {"dict", "message", "keys"},
	"safeIndex":      {"value", "index", "default"},
	"set":            {"dict", "key", "value"},
	"slice":          {"value", "args"},
	"string":         {"value"},
	"String":         {"value"},
	"toBash":         {"value"},
	"toHcl":          {"value"},
	"toInternalHcl":  {"value"},
	"toJson":         {"value"},
	"toPrettyHcl":    {"value"},
	"toPrettyJson":   {"value"},
	"toPrettyTFVars": {"value"},
	"toQuotedHcl":    {"value"},
	"toQuotedJson":   {"value"},
	"toQuotedTFVars": {"value"},
	"toTFVars":       {"value"},
	"toYaml":         {"value"},
	"undef":          {"default", "values"},
	"yaml":           {"yaml", "context"},
}

var dataFuncsAliases = map[string][]string{
	"data":          {"DATA", "fromData", "fromDATA"},
	"hcl":           {"HCL", "fromHcl", "fromHCL", "tfvars", "fromTFVars", "TFVARS", "fromTFVARS"},
	"json":          {"JSON", "fromJson", "fromJSON"},
	"lenc":          {"nbChars"},
	"toHcl":         {"toHCL"},
	"toInternalHcl": {"toInternalHCL", "toIHCL", "toIHcl"},
	"toJson":        {"toJSON"},
	"toPrettyHcl":   {"toPrettyHCL"},
	"toPrettyJson":  {"toPrettyJSON"},
	"toQuotedHcl":   {"toQuotedHCL"},
	"toQuotedJson":  {"toQuotedJSON"},
	"toYaml":        {"toYAML"},
	"undef":         {"ifUndef"},
	"yaml":          {"YAML", "fromYaml", "fromYAML"},
}

var dataFuncsHelp = map[string]string{
	"array":          "Ensures that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).",
	"bool":           "Converts the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)",
	"char":           "Returns the character corresponging to the supplied integer value",
	"content":        "Returns the content of a single element map (used to retrieve content in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"data":           "Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.",
	"extract":        "Extracts values from a slice or a map, indexes could be either integers for slice or strings for maps",
	"get":            "Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)",
	"hcl":            "Converts the supplied hcl string into data structure (Go spec). If context is omitted, default context is used.",
	"json":           "Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.",
	"key":            "Returns the key name of a single element map (used to retrieve name in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"lenc":           "Returns the number of actual character in a string",
	"merge":          "",
	"omit":           "",
	"pick":           "",
	"pickv":          "",
	"safeIndex":      "Returns the element at index position or default if index is outside bounds.",
	"set":            "Adds the value to the supplied map using key as identifier.",
	"slice":          "Returns a slice of the supplied object (equivalent to object[from:to]).",
	"string":         "Converts the supplied value into its string representation.",
	"String":         "Returns a String class object that allows invoking standard string operations as method.",
	"toBash":         "Converts the supplied value to bash compatible representation.",
	"toHcl":          "Converts the supplied value to compact HCL representation.",
	"toInternalHcl":  "Converts the supplied value to compact HCL representation used inside outer HCL definition.",
	"toJson":         "Converts the supplied value to compact JSON representation.",
	"toPrettyHcl":    "Converts the supplied value to pretty HCL representation.",
	"toPrettyJson":   "Converts the supplied value to pretty JSON representation.",
	"toPrettyTFVars": "Converts the supplied value to pretty HCL representation (without multiple map declarations).",
	"toQuotedHcl":    "Converts the supplied value to compact quoted HCL representation.",
	"toQuotedJson":   "Converts the supplied value to compact quoted JSON representation.",
	"toQuotedTFVars": "Converts the supplied value to compact HCL representation (without multiple map declarations).",
	"toTFVars":       "Converts the supplied value to compact HCL representation (without multiple map declarations).",
	"toYaml":         "Converts the supplied value to YAML representation.",
	"undef":          "Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, \"\" are not considered as unset).",
	"yaml":           "Converts the supplied yaml string into data structure (Go spec). If context is omitted, default context is used.",
}

func (t *Template) addDataFuncs() {
	options := funcOptions{
		funcHelp:    dataFuncsHelp,
		funcArgs:    dataFuncsArgs,
		funcAliases: dataFuncsAliases,
	}
	t.AddFunctions(dataFuncsBase, dataBase, options)
	t.AddFunctions(dataFuncsConversion, dataConversion, options)
	t.AddFunctions(map[string]interface{}{
		"data": t.fromData,
		"hcl":  t.fromHCL,
		"json": t.fromJSON,
		"yaml": t.fromYAML,
	}, dataConversion, options)
	t.optionsEnabled[Data] = true
}

func (t Template) fromData(source interface{}, context ...interface{}) (interface{}, error) {
	return t.dataConverter(utils.Interface2string(source), context...)
}

func (t Template) fromHCL(source interface{}, context ...interface{}) (interface{}, error) {
	return t.hclConverter(utils.Interface2string(source), context...)
}

func (t Template) fromJSON(source interface{}, context ...interface{}) (interface{}, error) {
	return t.jsonConverter(utils.Interface2string(source), context...)
}

func (t Template) fromYAML(source interface{}, context ...interface{}) (interface{}, error) {
	return t.yamlConverter(utils.Interface2string(source), context...)
}

func toChar(value interface{}) (r interface{}, err error) {
	defer func() { err = trapError(err, recover()) }()
	return process(value, func(a interface{}) interface{} {
		return string(toInt(a))
	})
}

func toString(s interface{}) string            { return fmt.Sprint(s) }
func toStringClass(s interface{}) utils.String { return utils.String(toString(s)) }

func toHCL(v interface{}) (string, error) {
	output, err := hcl.Marshal(v)
	return string(output), err
}

func toInternalHCL(v interface{}) (string, error) {
	output, err := hcl.MarshalInternal(v)
	return string(output), err
}

func toPrettyHCL(v interface{}) (string, error) {
	output, err := hcl.MarshalIndent(v, "", "  ")
	return string(output), err
}

func toQuotedHCL(v interface{}) (string, error) {
	output, err := hcl.Marshal(v)
	result := fmt.Sprintf("%q", output)
	return result[1 : len(result)-1], err
}

func toTFVars(v interface{}) (string, error) {
	output, err := hcl.MarshalTFVars(v)
	return string(output), err
}

func toPrettyTFVars(v interface{}) (string, error) {
	output, err := hcl.MarshalTFVarsIndent(v, "", "  ")
	return string(output), err
}

func toQuotedTFVars(v interface{}) (string, error) {
	output, err := hcl.MarshalTFVars(v)
	result := fmt.Sprintf("%q", output)
	return result[1 : len(result)-1], err
}

func toJSON(v interface{}) (string, error) {
	output, err := json.Marshal(v)
	return string(output), err
}

func toPrettyJSON(v interface{}) (string, error) {
	output, err := json.MarshalIndent(v, "", "  ")
	return string(output), err
}

func toQuotedJSON(v interface{}) (string, error) {
	output, err := json.Marshal(v)
	result := fmt.Sprintf("%q", output)
	return result[1 : len(result)-1], err
}

func array(value interface{}) interface{} {
	if value == nil {
		return value
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		return value
	default:
		return []interface{}{value}
	}
}

func get(arg1, arg2 interface{}) (result interface{}, err error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter. So all following forms are supported:
	//    get map key
	//    get key map
	//    map | get key
	//    key | get map

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Cannot retrieve key from undefined map: %v", e)
		}
	}()

	var (
		dict map[string]interface{}
		key  string
	)
	if reflect.TypeOf(arg1).Kind() == reflect.Map {
		dict = arg1.(map[string]interface{})
		key = arg2.(string)
	} else {
		key = arg1.(string)
		dict = arg2.(map[string]interface{})
	}
	return dict[key], nil
}

func set(arg1, arg2, arg3 interface{}) (result string, err error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter. So all following forms are supported:
	//    set map key value
	//    set key value map
	//    map | set key value
	//    value | set map key
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Cannot set key from undefined map: %v", e)
		}
	}()

	var (
		dict  map[string]interface{}
		key   string
		value interface{}
	)
	if reflect.TypeOf(arg1).Kind() == reflect.Map {
		dict = arg1.(map[string]interface{})
		key = arg2.(string)
		value = arg3
	} else {
		key = arg1.(string)
		value = arg2
		dict = arg3.(map[string]interface{})
	}
	dict[key] = value
	return "", nil
}

func key(v interface{}) (interface{}, error) {
	key, _, err := getSingleMapElement(v)
	return key, err
}

func content(v interface{}) (interface{}, error) {
	_, value, err := getSingleMapElement(v)
	return value, err
}

type dataConverter func([]byte, interface{}) error

// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
func (t Template) converter(converter dataConverter, content string, sourceWithError bool, context ...interface{}) (result interface{}, err error) {
	if err = converter([]byte(content), &result); err != nil && sourceWithError {
		source := "\n"
		for i, line := range utils.SplitLines(content) {
			source += fmt.Sprintf("%4d %s\n", i+1, line)
		}
		err = fmt.Errorf("%s\n%v", source, err)
	}
	return
}

// Apply a converter to the result of the template execution of the supplied string
func (t Template) templateConverter(converter dataConverter, str string, context ...interface{}) (result interface{}, err error) {
	var content string
	if content, _, err = t.runTemplate(str, context...); err == nil {
		result, err = t.converter(converter, content, true, context...)
	}
	return
}

// converts the supplied string containing yaml to go map
func (t Template) yamlConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(utils.YamlUnmarshal, str, context...)
}

// converts the supplied string containing json to go map
func (t Template) jsonConverter(str string, context ...interface{}) (interface{}, error) {
	return t.templateConverter(json.Unmarshal, str, context...)
}

// Converts the supplied string containing terraform/hcl to go map
func (t Template) hclConverter(str string, context ...interface{}) (result interface{}, err error) {
	return t.templateConverter(hcl.Unmarshal, str, context...)
}

// Converts the supplied string containing yaml, json or terraform/hcl to go map
func (t Template) dataConverter(str string, context ...interface{}) (result interface{}, err error) {
	converter := func(bs []byte, out interface{}) (err error) {
		return utils.ConvertData(string(bs), out)
	}
	return t.templateConverter(converter, str, context...)
}

var sprigPick = sprig.GenericFuncMap()["pick"].(func(map[string]interface{}, ...string) map[string]interface{})
var sprigOmit = sprig.GenericFuncMap()["omit"].(func(map[string]interface{}, ...string) map[string]interface{})

func pick(dict map[string]interface{}, keys ...interface{}) map[string]interface{} {
	return sprigPick(dict, toStrings(convertArgs(nil, keys...))...)
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
	return sprigOmit(dict, toStrings(convertArgs(nil, keys...))...)
}
