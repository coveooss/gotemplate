package template

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/json"
	"github.com/coveo/gotemplate/types"
	"github.com/coveo/gotemplate/utils"
	"github.com/coveo/gotemplate/xml"
	"github.com/coveo/gotemplate/yaml"
)

const (
	dataBase       = "Data Manipulation"
	dataConversion = "Data Conversion"
)

var dataFuncsBase = dictionary{
	"String":    toStringClass,
	"array":     array,
	"bool":      strconv.ParseBool,
	"char":      toChar,
	"content":   content,
	"dict":      createDict,
	"extract":   extract,
	"get":       get,
	"hasKey":    hasKey,
	"key":       key,
	"keys":      keys,
	"lenc":      utf8.RuneCountInString,
	"merge":     merge,
	"omit":      omit,
	"pick":      pick,
	"pickv":     pickv,
	"pluck":     pluck,
	"safeIndex": safeIndex,
	"set":       set,
	"slice":     slice,
	"string":    toString,
	"undef":     utils.IfUndef,
	"unset":     unset,
}

var dataFuncsConversion = dictionary{
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
	//"toXml":          toXML,
	"toYaml": toYAML,
}

var dataFuncsArgs = arguments{
	"String":         {"value"},
	"array":          {"value"},
	"bool":           {"str"},
	"char":           {"value"},
	"content":        {"keymap"},
	"data":           {"data", "context"},
	"extract":        {"source", "indexes"},
	"get":            {"map", "key"},
	"hasKey":         {"dictionary", "key"},
	"hcl":            {"hcl", "context"},
	"json":           {"json", "context"},
	"key":            {"value"},
	"keys":           {"dictionary"},
	"lenc":           {"str"},
	"merge":          {"destination", "sources"},
	"omit":           {"dict", "keys"},
	"pick":           {"dict", "keys"},
	"pickv":          {"dict", "message", "keys"},
	"pluck":          {"key", "dictionaries"},
	"safeIndex":      {"value", "index", "default"},
	"set":            {"dict", "key", "value"},
	"slice":          {"value", "args"},
	"string":         {"value"},
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
	"unset":          {"dictionary", "key"},
	"xml":            {"yaml", "context"},
	"yaml":           {"yaml", "context"},
}

var dataFuncsAliases = aliases{
	"data":          {"DATA", "fromData", "fromDATA"},
	"dict":          {"dictionary"},
	"hasKey":        {"has"},
	"hcl":           {"HCL", "fromHcl", "fromHCL", "tfvars", "fromTFVars", "TFVARS", "fromTFVARS"},
	"json":          {"JSON", "fromJson", "fromJSON"},
	"lenc":          {"nbChars"},
	"toHcl":         {"toHCL"},
	"toInternalHcl": {"toInternalHCL", "toIHCL", "toIHcl"},
	"toJson":        {"toJSON"},
	"toPrettyHcl":   {"toPrettyHCL"},
	"toPrettyJson":  {"toPrettyJSON"},
	"toPrettyXml":   {"toPrettyXML"},
	"toQuotedHcl":   {"toQuotedHCL"},
	"toQuotedJson":  {"toQuotedJSON"},
	"toXml":         {"toXML"},
	"toYaml":        {"toYAML"},
	"undef":         {"ifUndef"},
	"unset":         {"delete", "remove"},
	"xml":           {"XML", "fromXml", "fromXML"},
	"yaml":          {"YAML", "fromYaml", "fromYAML"},
}

var dataFuncsHelp = descriptions{
	"String":         "Returns a String class object that allows invoking standard string operations as method.",
	"array":          "Ensures that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).",
	"bool":           "Converts the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)",
	"char":           "Returns the character corresponging to the supplied integer value",
	"content":        "Returns the content of a single element map (used to retrieve content in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"data":           "Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.",
	"dict":           "Returns a new dictionary from a list of pairs (key, value).",
	"extract":        "Extracts values from a slice or a map, indexes could be either integers for slice or strings for maps",
	"get":            "Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)",
	"hasKey":         "Returns true if the dictionary contains the specified key.",
	"hcl":            "Converts the supplied hcl string into data structure (Go spec). If context is omitted, default context is used.",
	"json":           "Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.",
	"key":            "Returns the key name of a single element map (used to retrieve name in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"keys":           "Returns a list of all of the keys in a dict (in alphabetical order).",
	"lenc":           "Returns the number of actual character in a string",
	"merge":          "Merges two or more dictionaries into one, giving precedence to the dest dictionary.",
	"omit":           "Returns a new dict with all the keys that do not match the given keys.",
	"pick":           "Selects just the given keys out of a dictionary, creating a new dict.",
	"pickv":          "Same as pick, but returns an error message if there are intruders in supplied dictionary.",
	"pluck":          "Extracts a list of values matching the supplied key from a list of dictionary.",
	"safeIndex":      "Returns the element at index position or default if index is outside bounds.",
	"set":            "Adds the value to the supplied map using key as identifier.",
	"slice":          "Returns a slice of the supplied object (equivalent to object[from:to]).",
	"string":         "Converts the supplied value into its string representation.",
	"toBash":         "Converts the supplied value to bash compatible representation.",
	"toHcl":          "Converts the supplied value to compact HCL representation.",
	"toInternalHcl":  "Converts the supplied value to compact HCL representation used inside outer HCL definition.",
	"toJson":         "Converts the supplied value to compact JSON representation.",
	"toPrettyHcl":    "Converts the supplied value to pretty HCL representation.",
	"toPrettyJson":   "Converts the supplied value to pretty JSON representation.",
	"toPrettyTFVars": "Converts the supplied value to pretty HCL representation (without multiple map declarations).",
	"toPrettyXml":    "Converts the supplied value to pretty XML representation.",
	"toQuotedHcl":    "Converts the supplied value to compact quoted HCL representation.",
	"toQuotedJson":   "Converts the supplied value to compact quoted JSON representation.",
	"toQuotedTFVars": "Converts the supplied value to compact HCL representation (without multiple map declarations).",
	"toTFVars":       "Converts the supplied value to compact HCL representation (without multiple map declarations).",
	"toXml":          "Converts the supplied value to XML representation.",
	"toYaml":         "Converts the supplied value to YAML representation.",
	"undef":          "Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, \"\" are not considered as unset).",
	"unset":          "Removes an element from a dictionary.",
	"xml":            "Converts the supplied xml string into data structure (Go spec). If context is omitted, default context is used.",
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
	t.AddFunctions(dictionary{
		"data": t.dataConverter,
		"hcl":  t.hclConverter,
		"json": t.jsonConverter,
		//"xml":  t.xmlConverter,
		"yaml": t.yamlConverter,
	}, dataConversion, options)
	t.optionsEnabled[Data] = true
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

func toXML(v interface{}) (string, error) {
	output, err := xml.Marshal(v)
	return string(output), err
}

func toYAML(v interface{}) (string, error) {
	output, err := yaml.Marshal(v)
	return string(output), err
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

func get(arg1, arg2 interface{}) (interface{}, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := types.TryAsDictionary(arg1); err == nil {
		return dict.Get(arg2), nil
	} else if dict, err = types.TryAsDictionary(arg2); err == nil {
		return dict.Get(arg1), nil
	} else {
		return nil, fmt.Errorf("Must supply dictionary object")
	}
}

func hasKey(arg1, arg2 interface{}) (interface{}, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := types.TryAsDictionary(arg1); err == nil {
		return dict.Has(arg2), nil
	} else if dict, err = types.TryAsDictionary(arg2); err == nil {
		return dict.Has(arg1), nil
	} else {
		return nil, fmt.Errorf("Must supply dictionary object")
	}
}

func set(arg1, arg2, arg3 interface{}) (string, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := types.TryAsDictionary(arg1); err == nil {
		dict.Set(arg2, arg3)
	} else if dict, err = types.TryAsDictionary(arg3); err == nil {
		dict.Set(arg1, arg2)
	} else {
		return "", fmt.Errorf("Must supply dictionary object")
	}
	return "", nil
}

func unset(arg1, arg2 interface{}) (string, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := types.TryAsDictionary(arg1); err == nil {
		dict.Delete(arg2)
	} else if dict, err = types.TryAsDictionary(arg2); err == nil {
		dict.Delete(arg1)
	} else {
		return "", fmt.Errorf("Must supply dictionary object")
	}
	return "", nil
}

func merge(target Dictionary, dict Dictionary, otherDicts ...Dictionary) Dictionary {
	return target.Merge(dict, otherDicts...)
}

func key(v interface{}) (interface{}, error) {
	key, _, err := getSingleMapElement(v)
	return key, err
}

func content(v interface{}) (interface{}, error) {
	_, value, err := getSingleMapElement(v)
	return value, err
}

type marshaler func(interface{}) ([]byte, error)
type unMarshaler func([]byte, interface{}) error

// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
func (t Template) converter(from unMarshaler, content string, sourceWithError bool, context ...interface{}) (result interface{}, err error) {
	if err = from([]byte(content), &result); err != nil && sourceWithError {
		source := "\n"
		for i, line := range types.SplitLines(content) {
			source += fmt.Sprintf("%4d %s\n", i+1, line)
		}
		err = fmt.Errorf("%s\n%v", source, err)
	}
	return
}

// Apply a converter to the result of the template execution of the supplied string
func (t Template) templateConverter(to marshaler, from unMarshaler, source interface{}, context ...interface{}) (result interface{}, err error) {
	if source == nil {
		return nil, nil
	}
	if reflect.TypeOf(source).Kind() != reflect.String {
		if source, err = to(source); err != nil {
			return
		}
		source = string(source.([]byte))
	}

	var content string
	if content, _, err = t.runTemplate(fmt.Sprint(source), context...); err == nil {
		result, err = t.converter(from, content, true, context...)
	}
	return
}

func (t Template) xmlConverter(source interface{}, context ...interface{}) (interface{}, error) {
	return t.templateConverter(xml.Marshal, xml.Unmarshal, source, context...)
}

func (t Template) yamlConverter(source interface{}, context ...interface{}) (interface{}, error) {
	return t.templateConverter(yaml.Marshal, yaml.Unmarshal, source, context...)
}

func (t Template) jsonConverter(source interface{}, context ...interface{}) (interface{}, error) {
	return t.templateConverter(json.Marshal, json.Unmarshal, source, context...)
}

func (t Template) hclConverter(source interface{}, context ...interface{}) (result interface{}, err error) {
	return t.templateConverter(hcl.Marshal, hcl.Unmarshal, source, context...)
}

func (t Template) dataConverter(source interface{}, context ...interface{}) (result interface{}, err error) {
	return t.templateConverter(
		func(in interface{}) ([]byte, error) { return []byte(fmt.Sprint(in)), nil },
		func(bs []byte, out interface{}) error { return utils.ConvertData(string(bs), out) },
		source, context...)
}

// Dictionary represents an implementation of IDictionary
type Dictionary = types.IDictionary

// List represents an implementation of IGenericList
type List = types.IGenericList

func pick(dict Dictionary, keys ...interface{}) Dictionary {
	return dict.Clone(keys...)
}

func omit(dict Dictionary, key interface{}, otherKeys ...interface{}) Dictionary {
	return dict.Omit(key, otherKeys...)
}

func pickv(dict Dictionary, message string, key interface{}, otherKeys ...interface{}) (interface{}, error) {
	o := dict.Omit(key, otherKeys...)

	if o.Len() > 0 {
		over := strings.Join(toStrings(o.Keys()), ", ")
		if strings.Contains(message, "%v") {
			message = fmt.Sprintf(message, over)
		} else {
			message = iif(message == "", "Unwanted values", message).(string)
			message = fmt.Sprintf("%s %s", message, over)
		}
		return nil, fmt.Errorf(message)
	}
	return pick(dict, append(otherKeys, key)), nil
}

func keys(dict Dictionary) List { return dict.Keys() }

func createDict(v ...interface{}) (Dictionary, error) {
	if len(v)%2 != 0 {
		return nil, fmt.Errorf("Must supply even number of arguments (keypair)")
	}

	result := types.CreateDictionary(len(v) / 2)
	for i := 0; i < len(v); i += 2 {
		result.Set(v[i], v[i+1])
	}
	return result, nil
}

func pluck(key interface{}, dicts ...Dictionary) List {
	result := types.CreateList(0, len(dicts))
	for i := range dicts {
		if dicts[i].Has(key) {
			result.Append(dicts[i].Get(key))
		}
	}
	return result
}
