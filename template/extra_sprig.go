package template

import (
	"strings"

	"github.com/Masterminds/sprig/v3"
)

const (
	sprigCrypto         = "Sprig Cryptographic & Security, http://masterminds.github.io/sprig/crypto.html"
	sprigDate           = "Sprig Date, http://masterminds.github.io/sprig/date.html"
	sprigDict           = "Sprig Dictionary, http://masterminds.github.io/sprig/dicts.html"
	sprigEncoding       = "Sprig Encoding, http://masterminds.github.io/sprig/encoding.html"
	sprigFilePath       = "Sprig File Path, http://masterminds.github.io/sprig/paths.html"
	sprigFlow           = "Sprig Flow Control, http://masterminds.github.io/sprig/flow_control.html"
	sprigGen            = "Sprig General, http://masterminds.github.io/sprig/"
	sprigList           = "Sprig List, http://masterminds.github.io/sprig/lists.html"
	sprigMath           = "Sprig Mathematics, http://masterminds.github.io/sprig/math.html"
	sprigRegex          = "Sprig Regex, http://masterminds.github.io/sprig/strings.html"
	sprigSemver         = "Sprig Version comparison, http://masterminds.github.io/sprig/semver.html"
	sprigString         = "Sprig Strings, http://masterminds.github.io/sprig/strings.html"
	sprigReflect        = "Sprig Reflection, http://masterminds.github.io/sprig/reflection.html"
	sprigDefault        = "Sprig Default, http://masterminds.github.io/sprig/defaults.html"
	sprigOS             = "Sprig OS, http://masterminds.github.io/sprig/os.html"
	sprigTypeConversion = "Sprig Type Conversion, http://masterminds.github.io/sprig/conversion.html"
	sprigStringList     = "Sprig String Slice, http://masterminds.github.io/sprig/string_slice.html"
	sprigURL            = "Sprig URL functions, http://masterminds.github.io/sprig/url.html"
	sprigNetwork        = "Sprig Network functions, http://masterminds.github.io/sprig/network.html"
)

var sprigFuncs funcTableMap

var sprigFuncMap = sprig.GenericFuncMap()
var sprigDef = sprigFuncMap["default"].(func(interface{}, ...interface{}) interface{})
var sprigRound = sprigFuncMap["round"].(func(a interface{}, p int, r_opt ...float64) float64)

func (t *Template) addSprigFuncs() {
	if sprigFuncs == nil {
		// We get the list of aliases to avoid issuing warning a sprig unmapped functions.
		aliases := make(map[string]string)
		for key, value := range sprigFuncRef {
			for _, alias := range value.aliases {
				aliases[alias] = key
			}
		}

		sprigFuncs = make(funcTableMap)
		for key, value := range sprigFuncMap {
			info := sprigFuncRef[key]
			mustPrefix := "must"
			if info.group == "" && strings.HasPrefix(key, mustPrefix) {
				// Sprig also support functions prefixed by must or must_, but we don't add them since the
				// other version returning an error is enough.
				continue
			}

			if info.group == "" {
				if aliases[key] == "" {
					InternalLog.Warningf("function %s from Sprig not found in documentation", key)
					continue
				}
				key = aliases[key]
				info = sprigFuncRef[key]
			}
			sprigFuncs[key] = &FuncInfo{function: value, group: info.group, aliases: info.aliases, arguments: info.arguments, description: info.description}
		}
	}
	t.addFunctions(sprigFuncs)
}

var sprigFuncRef = map[string]struct {
	arguments, aliases []string
	group, description string
}{
	"hello": {group: sprigGen, description: "Simple hello by Sprig"},

	// Date functions
	"ago":            {group: sprigDate, description: "The ago function returns duration from time.Now in seconds resolution.", arguments: []string{"date"}},
	"date_modify":    {group: sprigDate, description: "The dateModify takes a modification and a date and returns the timestamp.", arguments: []string{"fmt", "date"}},
	"date":           {group: sprigDate, description: "The date function formats a date [format](https://golang.org/pkg/time/#Time.Format).", arguments: []string{"fmt", "date"}},
	"dateInZone":     {group: sprigDate, description: "Same as date, but with a timezone.", arguments: []string{"fmt", "date", "zone"}, aliases: []string{"date_in_zone"}},
	"dateModify":     {group: sprigDate, description: "The dateModify takes a modification and a date and returns the timestamp.", arguments: []string{"fmt", "date"}},
	"duration":       {group: sprigDate, description: "Formats a given amount of seconds as a time.Duration.", arguments: []string{"second"}},
	"durationRound":  {group: sprigDate, description: "Rounds a given duration to the most significant unit.", arguments: []string{"duration"}},
	"htmlDate":       {group: sprigDate, description: "The htmlDate function formats a date for inserting into an HTML date picker input field.", arguments: []string{"date"}},
	"htmlDateInZone": {group: sprigDate, description: "Same as htmlDate, but with a timezone.", arguments: []string{"date", "zone"}},
	"now":            {group: sprigDate, description: "The current date/time. Use this in conjunction with other date functions."},
	"toDate":         {group: sprigDate, description: "Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.", arguments: []string{"fmt", "str"}},
	"unixEpoch":      {group: sprigDate, description: "Returns the seconds since the unix epoch for a 'time.Time'.", arguments: []string{"date"}},

	// Strings functions
	"abbrev":        {group: sprigString, description: "Truncates a string with ellipses (...).", arguments: []string{"width", "str"}},
	"abbrevboth":    {group: sprigString, description: "Abbreviates both sides with ellipses (...).", arguments: []string{"left", "right", "str"}},
	"camelcase":     {group: sprigString, description: "Converts string from snake_case to CamelCase.", arguments: []string{"str"}},
	"cat":           {group: sprigString, description: "Concatenates multiple strings together into one, separating them with spaces."},
	"containsSprig": {group: sprigString, description: "Tests to see if one string is contained inside of another.", arguments: []string{"substr", "str"}, aliases: []string{"contains"}},
	"hasPrefix":     {group: sprigString, description: "Tests whether a string has a given prefix.", arguments: []string{"prefix", "str"}},
	"hasSuffix":     {group: sprigString, description: "Tests whether a string has a given suffix.", arguments: []string{"suffix", "str"}},
	"indentSprig":   {group: sprigString, description: "Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.", arguments: []string{"spaces", "str"}, aliases: []string{"indent"}},
	"initials":      {group: sprigString, description: "Given multiple words, takes the first letter of each word and combine.", arguments: []string{"str"}},
	"kebabcase":     {group: sprigString, description: "Convert string from camelCase to kebab-case.", arguments: []string{"str"}},
	"lower":         {group: sprigString, description: "Converts the entire string to lowercase.", arguments: []string{"str"}},
	"nindentSprig":  {group: sprigString, description: "Same as the indent function, but prepends a new line to the beginning of the string.", arguments: []string{"spaces", "str"}, aliases: []string{"nindent"}},
	"nospace":       {group: sprigString, description: "Removes all whitespace from a string.", arguments: []string{"str"}},
	"plural":        {group: sprigString, description: "Pluralizes a string.", arguments: []string{"one", "many", "count"}},
	"quote":         {group: sprigString, description: "Wraps each argument with double quotes.", arguments: []string{"str"}},
	"randAlpha":     {group: sprigString, description: "Generates random string with letters.", arguments: []string{"count"}},
	"randAlphaNum":  {group: sprigString, description: "Generates random string with letters and digits.", arguments: []string{"count"}},
	"randAscii":     {group: sprigString, description: "Generates random string with ASCII printable characters.", arguments: []string{"count"}},
	"randNumeric":   {group: sprigString, description: "Generates random string with digits.", arguments: []string{"count"}},
	"repeatSprig":   {group: sprigString, description: "Repeats a string multiple times.", arguments: []string{"count", "str"}, aliases: []string{"repeat"}},
	"replace":       {group: sprigString, description: "Performs simple string replacement.", arguments: []string{"old", "new", "src"}},
	"shuffle":       {group: sprigString, description: "Shuffle a string.", arguments: []string{"str"}},
	"snakecase":     {group: sprigString, description: "Converts string from camelCase to snake_case.", arguments: []string{"str"}},
	"squote":        {group: sprigString, description: "Wraps each argument with single quotes."},
	"substr":        {group: sprigString, description: "Get a substring from a string.", arguments: []string{"start", "length", "str"}},
	"swapcase":      {group: sprigString, description: "Swaps the uppercase to lowercase and lowercase to uppercase.", arguments: []string{"str"}},
	"title":         {group: sprigString, description: "Converts to title case.", arguments: []string{"str"}},
	"toString":      {group: sprigString, description: "Converts any value to string.", arguments: []string{"value"}},
	"trim":          {group: sprigString, description: "Removes space from either side of a string.", arguments: []string{"str"}},
	"trimAll":       {group: sprigString, description: "Removes given characters from the front or back of a string.", aliases: []string{"trimall"}, arguments: []string{"chars", "str"}},
	"trimPrefix":    {group: sprigString, description: "Trims just the prefix from a string if present.", arguments: []string{"prefix", "str"}},
	"trimSuffix":    {group: sprigString, description: "Trims just the suffix from a string if present.", arguments: []string{"suffix", "str"}},
	"truncSprig":    {group: sprigString, description: "Truncates a string (and add no suffix).", arguments: []string{"length", "str"}, aliases: []string{"trunc"}},
	"untitle":       {group: sprigString, description: `Removes title casing.`, arguments: []string{"str"}},
	"upper":         {group: sprigString, description: "Converts the entire string to uppercase.", arguments: []string{"str"}},
	"wrapSprig":     {group: sprigString, description: "Wraps text at a given column count.", arguments: []string{"length", "str"}, aliases: []string{"wrap"}},
	"wrapWith":      {group: sprigString, description: "Works as wrap, but lets you specify the string to wrap with (wrap uses \\n).", arguments: []string{"length", "spe", "str"}},

	"atoi":      {group: sprigTypeConversion, description: "Convert a string to an integer.", arguments: []string{"str"}},
	"float64":   {group: sprigTypeConversion, description: "Convert to a `float64`.", arguments: []string{"value"}},
	"int64":     {group: sprigTypeConversion, description: "Convert a numeric to an `int64`.", arguments: []string{"value"}},
	"intSprig":  {group: sprigTypeConversion, description: "Convert to an `int` at the system's width.", aliases: []string{"int"}, arguments: []string{"value"}},
	"toDecimal": {group: sprigTypeConversion, description: "Convert a unix octal to a 'int64'.", arguments: []string{"octalString"}},

	"join":      {group: sprigStringList, description: "Joins a list of strings into a single string, with the given separator.", arguments: []string{"separator", "list"}},
	"sortAlpha": {group: sprigStringList, description: "Sorts a list of strings into alphabetical (lexicographical) order.", arguments: []string{"list"}},
	"split":     {group: sprigStringList, description: "Splits a string into a `dict`. It is designed to make it easy to use template dot notation for accessing members", arguments: []string{"separator", "str"}},
	"splitList": {group: sprigStringList, description: "Splits a string into a list of strings.", arguments: []string{"separator", "str"}},
	"splitn":    {group: sprigStringList, description: "Splits a string into a `dict`. It is designed to make it easy to use template dot notation for accessing members.", arguments: []string{"separator", "count", "str"}},
	"toStrings": {group: sprigStringList, description: "Given a list-like collection, produce a slice of strings.", arguments: []string{"list"}},

	// VERY basic arithmetic.
	"add1":       {group: sprigMath, description: "Increments a value by 1."},
	"add1f":      {group: sprigMath, description: "Increments a float value by 1."},
	"addf":       {group: sprigMath, description: "Sum float numbers with `add`. Accepts two or more inputs."},
	"addSprig":   {group: sprigMath, description: "Sum numbers with `add`. Accepts two or more inputs.", aliases: []string{"add"}},
	"ceilSprig":  {group: sprigMath, description: "Returns the greatest float value greater than or equal to input value.", aliases: []string{"ceil"}},
	"divf":       {group: sprigMath, description: "Performs float division."},
	"divSprig":   {group: sprigMath, description: "Performs integer division.", aliases: []string{"div"}},
	"floorSprig": {group: sprigMath, description: "Returns the greatest float value less than or equal to input value", aliases: []string{"floor"}},
	"maxf":       {group: sprigMath, description: "Returns the largest of a series of floats."},
	"maxSprig":   {group: sprigMath, description: "Returns the largest of a series of integers.", aliases: []string{"max", "biggest", "biggestSprig"}},
	"minf":       {group: sprigMath, description: "Returns the smallest of a series of floats."},
	"minSprig":   {group: sprigMath, description: "Returns the smallest of a series of integers.", aliases: []string{"min"}},
	"modSprig":   {group: sprigMath, description: "Performs integer modulo.", aliases: []string{"mod"}},
	"mulf":       {group: sprigMath, description: "Multiplies floats. Accepts two or more inputs."},
	"mulSprig":   {group: sprigMath, description: "Multiplies numbers. Accepts two or more inputs.", aliases: []string{"mul"}},
	"randInt":    {group: sprigMath, description: "Returns a random integer value from min (inclusive) to max (exclusive).", arguments: []string{"min", "max"}},
	"round":      {group: sprigMath, description: "Returns a float value with the remainder rounded to the given number to digits after the decimal point."},
	"subf":       {group: sprigMath, description: "Subtracts a floats from another floats.`"},
	"subSprig":   {group: sprigMath, description: "Subtracts a number from another number.`", aliases: []string{"sub"}},

	"seq":       {group: sprigMath, description: "Works like the bash seq command."},
	"until":     {group: sprigMath, description: "Builds a range of integers.", aliases: []string{"until"}},
	"untilStep": {group: sprigMath, description: "Generates a list of counting integers. But it allows you to define a start, stop, and step"},

	// Defaults
	"coalesce":          {group: sprigDefault, description: "Takes a list of values and returns the first non-empty one."},
	"compact":           {group: sprigDefault, description: "Removes entries with empty values.", arguments: []string{"list"}},
	"default":           {group: sprigDefault, description: "Set a simple default value.", arguments: []string{"default", "value"}},
	"empty":             {group: sprigDefault, description: "Returns true if the given value is considered empty.", arguments: []string{"value"}},
	"fromJsonSprig":     {group: sprigDefault, description: "Decodes a JSON document into a structure. If the input cannot be decoded as JSON the function will return an empty string. mustFromJson will return an error in case the JSON is invalid.", aliases: []string{"fromJson"}},
	"ternarySprig":      {group: sprigDefault, description: "If the test value is true, the first value will be returned, otherwise, the second is returned.", arguments: []string{"true", "false", "condition"}, aliases: []string{"ternary"}},
	"toJsonSprig":       {group: sprigDefault, description: "Encodes an item into a JSON string.", aliases: []string{"toJson"}, arguments: []string{"object"}},
	"toPrettyJsonSprig": {group: sprigDefault, description: "Encodes an item into a pretty (indented) JSON string.", aliases: []string{"toPrettyJson"}, arguments: []string{"object"}},
	"toRawJson":         {group: sprigDefault, description: "Encodes an item into JSON string with HTML characters unescaped.", arguments: []string{"object"}},

	// Reflection
	"deepEqual":  {group: sprigReflect, description: "returns true if two values are deeply equal."},
	"kindIs":     {group: sprigReflect, description: "Let you verify that a value is a particular kind.", aliases: []string{"kindis"}},
	"kindOf":     {group: sprigReflect, description: "Returns the kind of an object.", aliases: []string{"kindof"}},
	"typeIs":     {group: sprigReflect, description: "Like `kindIs`, but for types.", aliases: []string{"typeis"}},
	"typeIsLike": {group: sprigReflect, description: "Works as `typeIs`, except that it also dereferences pointers.", aliases: []string{"typeisLike"}},
	"typeOf":     {group: sprigReflect, description: "Returns the underlying type of a value.", aliases: []string{"typeof"}},

	// OS:
	"env":       {group: sprigOS, description: "Reads an environment variable."},
	"expandenv": {group: sprigOS, description: "Substitute environment variables in a string."},

	// File Paths:
	"base":    {group: sprigFilePath, description: "Returns the last element of a path."},
	"clean":   {group: sprigFilePath, description: "Clean up a path."},
	"dir":     {group: sprigFilePath, description: "Returns the directory, stripping the last part of the path."},
	"ext":     {group: sprigFilePath, description: "Returns the file extension."},
	"isAbs":   {group: sprigFilePath, description: "Check whether a file path is absolute."},
	"osBase":  {group: sprigFilePath, description: "Return the last element of a filepath."},
	"osClean": {group: sprigFilePath, description: "Clean up a path."},
	"osDir":   {group: sprigFilePath, description: "Return the directory, stripping the last part of the path."},
	"osExt":   {group: sprigFilePath, description: "Return the file extension."},
	"osIsAbs": {group: sprigFilePath, description: "To check whether a file path is absolute."},

	// Encoding:
	"b32dec": {group: sprigEncoding, description: "Decode with Base32."},
	"b32enc": {group: sprigEncoding, description: "Encode with Base32."},
	"b64dec": {group: sprigEncoding, description: "Decode with Base64."},
	"b64enc": {group: sprigEncoding, description: "Encode with Base64."},

	// Data Structures:
	"deepCopy":       {group: sprigDict, description: "Takes a value and makes a deep copy of the value."},
	"dictSprig":      {group: sprigDict, description: "Create a dictionary.", aliases: []string{"dict"}},
	"dig":            {group: sprigDict, description: "Traverses a nested set of dicts, selecting keys from a list of values. It returns a default value if any of the keys are not found at the associated dict.", arguments: []string{"key...", "default", "dict"}},
	"getSprig":       {group: sprigDict, description: "Given a map and a key, get the value from the map.", aliases: []string{"get"}, arguments: []string{"dict", "key"}},
	"hasKeySprig":    {group: sprigDict, description: "Returns 'true' if the given dict contains the given key.", aliases: []string{"hasKey"}, arguments: []string{"dict", "key"}},
	"keysSprig":      {group: sprigDict, description: "Returns a list of all of the keys in one or more dict types.", aliases: []string{"keys"}},
	"listSprig":      {group: sprigDict, description: "Create a list of elements.", aliases: []string{"list", "tuple", "tupleSprig"}},
	"mergeOverwrite": {group: sprigDict, description: "Merge two or more dictionaries into one, giving precedence from **right to left**, effectively overwriting values in the dest dictionary."},
	"mergeSprig":     {group: sprigDict, description: "Merge two or more dictionaries into one, giving precedence to the dest dictionary.", aliases: []string{"merge"}},
	"omitSprig":      {group: sprigDict, description: "Is similar to 'pick', except it returns a new `dict` with all the keys that _do not_ match the given keys.", aliases: []string{"omit"}},
	"pickSprig":      {group: sprigDict, description: "Selects just the given keys out of a dictionary, creating a new `dict`.", aliases: []string{"pick"}},
	"pluckSprig":     {group: sprigDict, description: "Makes it possible to give one key and multiple maps, and get a list of all of the matches.", aliases: []string{"pluck"}, arguments: []string{"name", "dicts"}},
	"setSprig":       {group: sprigDict, description: "Add a new key/value pair to a dictionary.", aliases: []string{"set"}, arguments: []string{"dict", "key", "value"}},
	"unsetSprig":     {group: sprigDict, description: "Given a map and a key, delete the key from the map.", aliases: []string{"unset"}, arguments: []string{"dict", "key"}},
	"valuesSprig":    {group: sprigDict, description: "Is similar to 'keys', except it returns a new 'list' with all the values of the source 'dict' (only one dictionary is supported).", aliases: []string{"values"}},

	// Lists functions
	"all":          {group: sprigList, description: "Filters items out of a list.", aliases: []string{"without"}},
	"any":          {group: sprigList, description: "Takes a list of values and returns true if any values are non-empty."},
	"appendSprig":  {group: sprigList, description: "Append a new item to an existing list, creating a new list.", aliases: []string{"append"}},
	"chunk":        {group: sprigList, description: "Split a list into chunks of given size, use chunk size list. This is useful for pagination.", arguments: []string{"nbChunks", "list"}},
	"concatSprig":  {group: sprigList, description: "Concatenate arbitrary number of lists into one.", aliases: []string{"concat"}},
	"first":        {group: sprigList, description: "Get the head item on a list."},
	"hasSprig":     {group: sprigList, description: "Test to see if a list has a particular element.", aliases: []string{"has"}},
	"initialSprig": {group: sprigList, description: "Complements last by returning all but the last element.", aliases: []string{"initial"}},
	"last":         {group: sprigList, description: "Get the last item on a list."},
	"prependSprig": {group: sprigList, description: "Push an element onto the front of a list, creating a new list.", aliases: []string{"prepend"}},
	"pushSprig":    {group: sprigList, description: "Append a new item to an existing list, creating a new list.", aliases: []string{"push"}},
	"restSprig":    {group: sprigList, description: "Get the tail of the list (everything but the first item).", aliases: []string{"rest"}},
	"reverseSprig": {group: sprigList, description: "Produce a new list with the reversed elements of the given list.", aliases: []string{"reverse"}},
	"sliceSprig":   {group: sprigList, description: "Get partial elements of a list.", aliases: []string{"slice"}},
	"uniqSprig":    {group: sprigList, description: "Generate a list with all of the duplicates removed.", aliases: []string{"uniq"}},
	"withoutSprig": {group: sprigList, description: "Filters items out of a list.", aliases: []string{"without"}},

	// Cryptographic functions
	"adler32sum":               {group: sprigCrypto, description: "Computes Adler-32 checksum.", arguments: []string{"input"}},
	"bcrypt":                   {group: sprigCrypto, description: "Receives a string, and generates its bcrypt hash."},
	"buildCustomCert":          {group: sprigCrypto, description: "Allows customizing the certificate.", arguments: []string{"base64Cert", "base64PrivateKey"}},
	"decryptAES":               {group: sprigCrypto, description: "Receives a base64 string encoded by the AES-256 CBC algorithm and returns the decoded text.", arguments: []string{"encoded", "key"}},
	"derivePassword":           {group: sprigCrypto, description: "Derive a specific password based on some shared 'master password' constraints.", arguments: []string{"counter", "passwordType", "password", "user", "site"}},
	"encryptAES":               {group: sprigCrypto, description: "Encrypts text with AES-256 CBC and returns a base64 encoded string.", arguments: []string{"secret", "key"}},
	"genCA":                    {group: sprigCrypto, description: "Generates a new, self-signed x509 certificate authority.", arguments: []string{"cn", "nbDays"}},
	"genCAWithKey":             {group: sprigCrypto, description: "Generates a new, self-signed x509 certificate authority.", arguments: []string{"cn", "nbDays"}},
	"genPrivateKey":            {group: sprigCrypto, description: "Generates a new, self-signed x509 certificate authority using a given private key.", arguments: []string{"cn", "npDays", "privatePEM"}},
	"genSelfSignedCert":        {group: sprigCrypto, description: "Generates a new, self-signed x509 certificate.", arguments: []string{"cn", "ipList", "dnsList", "nbDays"}},
	"genSelfSignedCertWithKey": {group: sprigCrypto, description: "Generates a new, self-signed x509 certificate using a given private key.", arguments: []string{"cn", "ipList", "dnsList", "nbDays", "privatePEM"}},
	"genSignedCert":            {group: sprigCrypto, description: "Generates a new, x509 certificate signed by the specified CA.", arguments: []string{"cn", "ipList", "dnsList", "nbDays", "ca"}},
	"genSignedCertWithKey":     {group: sprigCrypto, description: "Generates a new, x509 certificate signed by the specified CA, using a given private key.", arguments: []string{"cn", "ipList", "dnsList", "nbDays", "ca", "privatePEM"}},
	"htpasswd":                 {group: sprigCrypto, description: "Takes a username and password and generates a bcrypt hash of the password. The result can be used for basic authentication on an Apache HTTP Server.", arguments: []string{"user", "password"}},
	"randBytes":                {group: sprigCrypto, description: "Accepts a count N and generates a cryptographically secure (uses crypto/rand) random sequence of N bytes. The sequence is returned as a base64 encoded string."},
	"sha1sum":                  {group: sprigCrypto, description: "Computes SHA1 digest.", arguments: []string{"input"}},
	"sha256sum":                {group: sprigCrypto, description: "Computes SHA256 digest.", arguments: []string{"input"}},
	"sha512sum":                {group: sprigCrypto, description: "Computes SHA512 digest.", arguments: []string{"input"}},

	// UUIDs:
	"uuidv4": {group: sprigGen, description: "Returns a new UUID of the v4 (randomly generated) type.", aliases: []string{"uuid", "guid", "GUID"}},

	// SemVer:
	"semver":        {group: sprigSemver, description: "Parses a string into a Semantic Version.", arguments: []string{"version"}},
	"semverCompare": {group: sprigSemver, description: "A more robust comparison function is provided as semverCompare. This version supports version ranges.", arguments: []string{"constraints", "version"}},

	// Flow Control:
	"fail": {group: sprigFlow, description: "Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail."},

	// Regex
	"regexFind":              {group: sprigRegex, description: "Returns the first (left most) match of the regular expression in the input string.", arguments: []string{"regex", "str"}},
	"regexFindAll":           {group: sprigRegex, description: "Returns a slice of all matches of the regular expression in the input string.", arguments: []string{"regex", "str", "n"}},
	"regexMatch":             {group: sprigRegex, description: "Returns true if the input string matches the regular expression.", arguments: []string{"regex", "str"}},
	"regexQuoteMeta":         {group: sprigRegex, description: "Returns a string that escapes all regular expression meta-characters inside the argument text; the returned string is a regular expression matching the literal text"},
	"regexReplaceAll":        {group: sprigRegex, description: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement. Inside string replacement, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first sub-match.", arguments: []string{"regex", "str", "repl"}},
	"regexReplaceAllLiteral": {group: sprigRegex, description: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement The replacement string is substituted directly, without using Expand.", arguments: []string{"regex", "str", "repl"}},
	"regexSplit":             {group: sprigRegex, description: "Slices the input string into substrings separated by the expression and returns a slice of the substrings between those expression matches. The last parameter n determines the number of substrings to return, where -1 means return all matches.", arguments: []string{"regex", "str", "n"}},

	// URL
	"urlJoin":  {group: sprigURL, description: "Joins map (produced by `urlParse`) to produce URL string.", arguments: []string{"dictionary"}},
	"urlParse": {group: sprigURL, description: "Parses string for URL and produces dict with URL parts.", arguments: []string{"uri"}},

	// Network
	"getHostByName": {group: sprigNetwork, description: "Receives a domain name and returns the ip address.", arguments: []string{"name"}},
}
