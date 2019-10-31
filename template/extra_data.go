package template

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/coveooss/gotemplate/v3/xml"
	"github.com/coveooss/gotemplate/v3/yaml"
	"github.com/coveooss/multilogger"
)

const (
	dataBase       = "Data Manipulation"
	dataConversion = "Data Conversion"
)

var dataFuncsBase = dictionary{
	"String":    toStringClass,
	"append":    addElements,
	"array":     array,
	"bool":      multilogger.ParseBool,
	"char":      toChar,
	"contains":  contains,
	"content":   content,
	"dict":      createDict,
	"extract":   extract,
	"get":       get,
	"hasKey":    hasKey,
	"initial":   initial,
	"intersect": intersect,
	"isNil":     func(value interface{}) bool { return value == nil },
	"isSet":     func(value interface{}) bool { return value != nil },
	"isZero":    isZero,
	"key":       key,
	"keys":      keys,
	"lenc":      utf8.RuneCountInString,
	"list":      collections.NewList,
	"merge":     merge,
	"omit":      omit,
	"pick":      pick,
	"pickv":     pickv,
	"pluck":     pluck,
	"prepend":   prepend,
	"rest":      rest,
	"reverse":   reverse,
	"safeIndex": safeIndex,
	"set":       set,
	"slice":     slice,
	"string":    toString,
	"undef":     collections.IfUndef,
	"unique":    unique,
	"union":     union,
	"unset":     unset,
	"values":    values,
	"without":   without,
}

var dataFuncsConversion = dictionary{
	"toBash":         collections.ToBash,
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
	"toYaml":         toYAML,
	//"toXml":          toXML,
}

var dataFuncsArgs = arguments{
	"append":         {"list", "elements"},
	"array":          {"value"},
	"bool":           {"str"},
	"char":           {"value"},
	"contains":       {"list", "elements"},
	"content":        {"keymap"},
	"data":           {"data", "context"},
	"extract":        {"source", "indexes"},
	"get":            {"map", "key", "default"},
	"hasKey":         {"dictionary", "key"},
	"hcl":            {"hcl", "context"},
	"initial":        {"list"},
	"intersect":      {"list", "elements"},
	"json":           {"json", "context"},
	"key":            {"value"},
	"keys":           {"dictionary"},
	"lenc":           {"str"},
	"merge":          {"destination", "sources"},
	"omit":           {"dict", "keys"},
	"pick":           {"dict", "keys"},
	"pickv":          {"dict", "message", "keys"},
	"pluck":          {"key", "dictionaries"},
	"prepend":        {"list", "elements"},
	"rest":           {"list"},
	"reverse":        {"list"},
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
	"unique":         {"list"},
	"union":          {"list", "elements"},
	"unset":          {"dictionary", "key"},
	"without":        {"list", "elements"},
	"xml":            {"yaml", "context"},
	"yaml":           {"yaml", "context"},
}

var dataFuncsAliases = aliases{
	"append":        {"push"},
	"contains":      {"has"},
	"data":          {"DATA", "fromData", "fromDATA"},
	"dict":          {"dictionary"},
	"hcl":           {"HCL", "fromHcl", "fromHCL", "tfvars", "fromTFVars", "TFVARS", "fromTFVARS"},
	"isNil":         {"isNull"},
	"isZero":        {"isEmpty"},
	"json":          {"JSON", "fromJson", "fromJSON"},
	"lenc":          {"nbChars"},
	"list":          {"tuple"},
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
	"unique":        {"uniq"},
	"unset":         {"delete", "remove"},
	"xml":           {"XML", "fromXml", "fromXML"},
	"yaml":          {"YAML", "fromYaml", "fromYAML"},
}

var dataFuncsHelp = descriptions{
	"String":         "Returns a String class object that allows invoking standard string operations as method.",
	"append":         "Append new items to an existing list, creating a new list.",
	"array":          "Ensures that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).",
	"bool":           "Converts the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)",
	"char":           "Returns the character corresponging to the supplied integer value",
	"contains":       "Test to see if a list has a particular elements.",
	"content":        "Returns the content of a single element map (used to retrieve content in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"data":           "Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.",
	"dict":           "Returns a new dictionary from a list of pairs (key, value).",
	"extract":        "Extracts values from a slice or a map, indexes could be either integers for slice or strings for maps",
	"get":            "Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)",
	"hasKey":         "Returns true if the dictionary contains the specified key.",
	"hcl":            "Converts the supplied hcl string into data structure (Go spec). If context is omitted, default context is used.",
	"initial":        "Returns but the last element.",
	"intersect":      "Returns a list that is the intersection of the list and all arguments (removing duplicates).",
	"isNil":          "Returns true if the supplied value is nil.",
	"isSet":          "Returns true if the supplied value is not nil.",
	"isZero":         "Returns true if the supplied value is false, 0, nil or empty.",
	"json":           "Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.",
	"key":            "Returns the key name of a single element map (used to retrieve name in a declaration like `value \"name\" { a = 1 b = 3}`)",
	"keys":           "Returns a list of all of the keys in a dict (in alphabetical order).",
	"lenc":           "Returns the number of actual character in a string.",
	"list":           "Returns a generic list from the supplied arguments.",
	"merge":          "Merges two or more dictionaries into one, giving precedence to the dest dictionary.",
	"omit":           "Returns a new dict with all the keys that do not match the given keys.",
	"pick":           "Selects just the given keys out of a dictionary, creating a new dict.",
	"pickv":          "Same as pick, but returns an error message if there are intruders in supplied dictionary.",
	"pluck":          "Extracts a list of values matching the supplied key from a list of dictionary.",
	"prepend":        "Push elements onto the front of a list, creating a new list.",
	"rest":           "Gets the tail of the list (everything but the first item)",
	"reverse":        "Produces a new list with the reversed elements of the given list.",
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
	"union":          "Returns a list that is the union of the list and all arguments (removing duplicates).",
	"unique":         "Generates a list with all of the duplicates removed.",
	"unset":          "Removes an element from a dictionary.",
	"values":         "Returns the list of values contained in a map.",
	"without":        "Filters items out of a list.",
	"xml":            "Converts the supplied xml string into data structure (Go spec). If context is omitted, default context is used.",
	"yaml":           "Converts the supplied yaml string into data structure (Go spec). If context is omitted, default context is used.",
}

var dataFuncsExamples = examples{
	"hasKey": {
		{`@hasKey(dict("key", "value"), "key")`, `{{ hasKey (dict "key" "value") "key" }}`, `true`},
		{`@hasKey("key", dict("key", "value"))`, ``, `true`},
		{`@hasKey(dict("key", "value"), "otherkey")`, ``, `false`},
	},
	"unset": {
		{`@{myDict} := dict("key", "value", "key2", "value2", "key3", "value3")
		@-unset($myDict, "key")
		@-unset("key2", $myDict)
		@-toJson($myDict)`, ``, `{"key3":"value3"}`},
	},
}

func (t *Template) addDataFuncs() {
	options := FuncOptions{
		FuncHelp:     dataFuncsHelp,
		FuncArgs:     dataFuncsArgs,
		FuncAliases:  dataFuncsAliases,
		FuncExamples: dataFuncsExamples,
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

func get(arg1, arg2 interface{}, defValue ...interface{}) (interface{}, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	var result interface{}
	if dict, err := collections.TryAsDictionary(arg1); err == nil {
		result = dict.Get(arg2)
	} else if dict, err = collections.TryAsDictionary(arg2); err == nil {
		result = dict.Get(arg1)
	} else {
		return nil, fmt.Errorf("Must supply dictionary object")
	}
	if result == nil {
		switch len(defValue) {
		case 0:
			break
		case 1:
			result = defValue[0]
		default:
			result = defValue
		}
	}
	return result, nil
}

func hasKey(arg1, arg2 interface{}) (interface{}, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := collections.TryAsDictionary(arg1); err == nil {
		return dict.Has(arg2), nil
	} else if dict, err = collections.TryAsDictionary(arg2); err == nil {
		return dict.Has(arg1), nil
	} else {
		return nil, fmt.Errorf("Must supply dictionary object")
	}
}

func set(arg1, arg2, arg3 interface{}) (string, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := collections.TryAsDictionary(arg1); err == nil {
		dict.Set(arg2, arg3)
	} else if dict, err = collections.TryAsDictionary(arg3); err == nil {
		dict.Set(arg1, arg2)
	} else {
		return "", fmt.Errorf("Must supply dictionary object")
	}
	return "", nil
}

func unset(arg1, arg2 interface{}) (string, error) {
	// In pipe execution, the map is often the last parameter, but we also support to
	// put the map as the first parameter.
	if dict, err := collections.TryAsDictionary(arg1); err == nil {
		dict.Delete(arg2)
	} else if dict, err = collections.TryAsDictionary(arg2); err == nil {
		dict.Delete(arg1)
	} else {
		return "", fmt.Errorf("Must supply dictionary object")
	}
	return "", nil
}

func merge(target iDictionary, dict iDictionary, otherDicts ...iDictionary) iDictionary {
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
		for i, line := range collections.SplitLines(content) {
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
		func(bs []byte, out interface{}) error { return collections.ConvertData(string(bs), out) },
		source, context...)
}

func pick(dict iDictionary, keys ...interface{}) iDictionary {
	return dict.Clone(keys...)
}

func omit(dict iDictionary, key interface{}, otherKeys ...interface{}) iDictionary {
	return dict.Omit(key, otherKeys...)
}

func pickv(dict iDictionary, message string, key interface{}, otherKeys ...interface{}) (interface{}, error) {
	o := dict.Omit(key, otherKeys...)

	if o.Len() > 0 {
		over := strings.Join(toStrings(o.GetKeys()), ", ")
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

func keys(dict iDictionary) iList   { return dict.GetKeys() }
func values(dict iDictionary) iList { return dict.GetValues() }

func createDict(v ...interface{}) (iDictionary, error) {
	if len(v)%2 != 0 {
		return nil, fmt.Errorf("Must supply even number of arguments (keypair)")
	}

	result := collections.CreateDictionary(len(v) / 2)
	for i := 0; i < len(v); i += 2 {
		result.Set(v[i], v[i+1])
	}
	return result, nil
}

func pluck(key interface{}, dicts ...iDictionary) iList {
	result := collections.CreateList(0, len(dicts))
	for i := range dicts {
		if dicts[i].Has(key) {
			result.Append(dicts[i].Get(key))
		}
	}
	return result
}

func rest(list interface{}) (interface{}, error)    { return slice(list, 1, -1) }
func initial(list interface{}) (interface{}, error) { return slice(list, 0, -2) }

func addElements(list interface{}, elements ...interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Append(elements...), nil
}

func prepend(list interface{}, elements ...interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Prepend(elements...), nil
}

func reverse(list interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Reverse(), nil
}

func unique(list interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Unique(), nil
}

func contains(list interface{}, elements ...interface{}) (r bool, err error) {
	// Then, the list argument must be a real list of elements
	defer func() { err = trapError(err, recover()) }()
	if _, err := collections.TryAsList(list); err != nil && len(elements) == 1 {
		if _, err2 := collections.TryAsList(elements[0]); err2 != nil {
			str, subStr := elements[0], list
			if s, isString := str.(collections.String); isString {
				// Check if the str argument is of type String
				str = string(s)
			}

			if s, isString := str.(string); isString {
				// Check if the list argument is of type string
				return strings.Contains(s, fmt.Sprint(subStr)), nil
			}
			return false, err
		}
		// Sprig has bad documentation and inverse the arguments, so we try to support both modes.
		list, elements = elements[0], []interface{}{list}
	}
	return collections.AsList(list).Contains(elements...), nil
}

func intersect(list interface{}, elements ...interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Intersect(elements...), nil
}

func union(list interface{}, elements ...interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Union(elements...), nil
}

func without(list interface{}, elements ...interface{}) (r iList, err error) {
	defer func() { err = trapError(err, recover()) }()
	return collections.AsList(list).Without(elements...), nil
}

func isZero(value interface{}) bool {
	return sprigDef(0, value) == 0
}
