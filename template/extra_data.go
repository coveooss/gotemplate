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

var dataFuncs = funcTableMap{
	// Base
	"array":     {f: array, group: dataBase, args: []string{"value"}, desc: "Ensure that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value)."},
	"bool":      {f: strconv.ParseBool, group: dataBase, args: []string{"str"}, desc: "Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)"},
	"char":      {f: toChar, group: dataBase, args: []string{"value"}, desc: "Returns the character corresponging to the supplied integer value"},
	"content":   {f: content, group: dataBase, args: []string{"keymap"}, desc: "Returns the content of a single element map (used to retrieve content in a declaration like `value \"name\" { a = 1 b = 3}`)"},
	"extract":   {f: extract, group: dataBase, args: []string{"source", "indexes"}, desc: "Extract values from a slice or a map, indexes could be either integers for slice or strings for maps"},
	"get":       {f: get, group: dataBase, args: []string{"map", "key"}, desc: "Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)"},
	"key":       {f: key, group: dataBase, args: []string{}, desc: ""},
	"lenc":      {f: utf8.RuneCountInString, group: dataBase, aliases: []string{"nbChars"}, desc: "Returns the number of actual character in a string"},
	"merge":     {f: utils.MergeMaps, group: dataBase, args: []string{}, desc: ""},
	"omit":      {f: omit, group: dataBase, args: []string{}, desc: ""},
	"pick":      {f: pick, group: dataBase, args: []string{}, desc: ""},
	"pickv":     {f: pickv, group: dataBase, args: []string{}, desc: ""},
	"safeIndex": {f: safeIndex, group: dataBase, args: []string{}, desc: ""},
	"set":       {f: set, group: dataBase, args: []string{}, desc: ""},
	"slice":     {f: slice, group: dataBase, args: []string{}, desc: ""},
	"string":    {f: toString, group: dataBase, args: []string{}, desc: ""},
	"undef":     {f: utils.IfUndef, group: dataBase, aliases: []string{"ifUndef"}, desc: ""},

	// Conversion to
	"toBash":         {f: utils.ToBash, group: dataConversion, desc: ""},
	"toHcl":          {f: toHCL, group: dataConversion, aliases: []string{"toHCL"}, desc: ""},
	"toJson":         {f: toJSON, group: dataConversion, aliases: []string{"toJSON"}, desc: ""},
	"toPrettyHcl":    {f: toPrettyHCL, group: dataConversion, aliases: []string{"toPrettyHCL"}, desc: ""},
	"toPrettyJson":   {f: toPrettyJSON, group: dataConversion, aliases: []string{"toPrettyJSON"}, desc: ""},
	"toPrettyTFVars": {f: toPrettyTFVars, group: dataConversion, desc: ""},
	"toQuotedHcl":    {f: toQuotedHCL, group: dataConversion, aliases: []string{"toQuotedHCL"}, desc: ""},
	"toQuotedJson":   {f: toQuotedJSON, group: dataConversion, args: []string{"toQuotedJSON"}, desc: ""},
	"toQuotedTFVars": {f: toQuotedTFVars, group: dataConversion, desc: ""},
	"toTFVars":       {f: toTFVars, group: dataConversion, desc: ""},
	"toYaml":         {f: utils.ToYaml, group: dataConversion, aliases: []string{"toYAML"}, desc: ""},
}

func (t *Template) addDataFuncs() {
	t.AddFunctions(dataFuncs)
	t.AddFunctions(funcTableMap{
		"data": {f: t.fromData, group: dataConversion, aliases: []string{"DATA", "fromData", "fromDATA"}, args: []string{}, desc: ""},
		"hcl":  {f: t.fromHCL, group: dataConversion, aliases: []string{"HCL", "fromHcl", "fromHCL", "tfvars", "fromTFVars", "TFVARS", "fromTFVARS"}, args: []string{}, desc: ""},
		"json": {f: t.fromJSON, group: dataConversion, aliases: []string{"JSON", "fromJson", "fromJSON"}, args: []string{}, desc: ""},
		"yaml": {f: t.fromYAML, group: dataConversion, aliases: []string{"YAML", "fromYaml", "fromYAML"}, args: []string{}, desc: ""},
	})
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

func toString(s interface{}) utils.String { return utils.String(fmt.Sprint(s)) }

func toHCL(v interface{}) (string, error) {
	output, err := hcl.Marshal(v)
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
