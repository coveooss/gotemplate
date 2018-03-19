package template

import (
	"github.com/Masterminds/sprig"
)

const (
	sprigCrypto         = "Sprig Cryptographic & Security, http://masterminds.github.io/sprig/crypto.html"
	sprigDate           = "Sprig Date, http://masterminds.github.io/sprig/date.html"
	sprigDict           = "Sprig Dictionnary, http://masterminds.github.io/sprig/dicst.html"
	sprigEncoding       = "Sprig Encoding, http://masterminds.github.io/sprig/encoding.html"
	sprigFilePath       = "Sprig File Path, http://masterminds.github.io/sprig/paths.html"
	sprigFlow           = "Sprig Flow Control, http://masterminds.github.io/sprig/flow_control.html"
	sprigGen            = "Sprig General, http://masterminds.github.io/sprig/"
	sprigList           = "Sprig List, http://masterminds.github.io/sprig/lists.html"
	sprigMath           = "Sprig Sprig Mathematics, http://masterminds.github.io/sprig/math.html"
	sprigRegex          = "Sprig Regex, http://masterminds.github.io/sprig/strings.html"
	sprigSemver         = "Sprig Version comparison, http://masterminds.github.io/sprig/semver.html"
	sprigString         = "Sprig Strings, http://masterminds.github.io/sprig/strings.html"
	sprigReflect        = "Sprig Reflection, http://masterminds.github.io/sprig/reflection.html"
	sprigDefault        = "Sprig Default, http://masterminds.github.io/sprig/defaults.html"
	sprigOS             = "Sprig OS, http://masterminds.github.io/sprig/defaults.html"
	sprigTypeConversion = "Sprig Type Conversion, http://masterminds.github.io/sprig/conversion.html"
	sprigStringList     = "Sprig String Slice, http://masterminds.github.io/sprig/string_slice.html"
)

var sprigFuncs funcTableMap

var sprigFuncMap = sprig.GenericFuncMap()

func (t *Template) addSprigFuncs() {
	if sprigFuncs == nil {
		sprigFuncs = make(funcTableMap)
		for key, value := range sprigFuncMap {
			info := sprigFuncRef[key]
			if info.group == "" {
				continue
			}
			sprigFuncs[key] = FuncInfo{function: value, group: info.group, aliases: info.aliases, arguments: info.arguments, description: info.description}
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
	"date":           {group: sprigDate, description: "The date function formats a dat (https://golang.org/pkg/time/#Time.Format).", arguments: []string{"fmt", "date"}},
	"now":            {group: sprigDate, description: "The current date/time. Use this in conjunction with other date functions."},
	"htmlDate":       {group: sprigDate, description: "The htmlDate function formates a date for inserting into an HTML date picker input field.", arguments: []string{"date"}},
	"htmlDateInZone": {group: sprigDate, description: "Same as htmlDate, but with a timezone.", arguments: []string{"date", "zone"}},
	"dateInZone":     {group: sprigDate, description: "Same as date, but with a timezone.", arguments: []string{"fmt", "date", "zone"}, aliases: []string{"date_in_zone"}},
	"dateModify":     {group: sprigDate, description: "The dateModify takes a modification and a date and returns the timestamp.", arguments: []string{"fmt", "date"}, aliases: []string{"date_modify"}},
	"ago":            {group: sprigDate, description: "The ago function returns duration from time.Now in seconds resolution.", arguments: []string{"date"}},
	"toDate":         {group: sprigDate, description: "Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.", arguments: []string{"fmt", "str"}},
	// Strings functions
	"abbrev":       {group: sprigString, description: "Truncates a string with ellipses (...).", arguments: []string{"width", "str"}},
	"abbrevboth":   {group: sprigString, description: "Abbreviates both sides with ellipses (...).", arguments: []string{"left", "right", "str"}},
	"camelcase":    {group: sprigString, description: "Converts string from snake_case to CamelCase.", arguments: []string{"str"}},
	"cat":          {group: sprigString, description: "Concatenates multiple strings together into one, separating them with spaces."},
	"contains":     {group: sprigString, description: "Tests to see if one string is contained inside of another.", arguments: []string{"substr", "str"}},
	"hasPrefix":    {group: sprigString, description: "Tests whether a string has a given prefix.", arguments: []string{"prefix", "str"}},
	"hasSuffix":    {group: sprigString, description: "Tests whether a string has a given suffix.", arguments: []string{"suffix", "str"}},
	"indent":       {group: sprigString, description: "Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.", arguments: []string{"spaces", "str"}},
	"initials":     {group: sprigString, description: "Given multiple words, takes the first letter of each word and combine.", arguments: []string{"str"}},
	"lower":        {group: sprigString, description: "Converts the entire string to lowercase.", arguments: []string{"str"}},
	"nindent":      {group: sprigString, description: "Same as the indent function, but prepends a new line to the beginning of the string.", arguments: []string{"spaces", "str"}},
	"nospace":      {group: sprigString, description: "Removes all whitespace from a string.", arguments: []string{"str"}},
	"plural":       {group: sprigString, description: "Pluralizes a string.", arguments: []string{"one", "many", "count"}},
	"quote":        {group: sprigString, description: "Wraps each argument with double quotes.", arguments: []string{"str"}},
	"randAlpha":    {group: sprigString, description: "Generates random string with letters.", arguments: []string{"count"}},
	"randAlphaNum": {group: sprigString, description: "Generates random string with letters and digits.", arguments: []string{"count"}},
	"randAscii":    {group: sprigString, description: "Generates random string with ASCII printable characters.", arguments: []string{"count"}},
	"randNumeric":  {group: sprigString, description: "Generates random string with digits.", arguments: []string{"count"}},
	"repeat":       {group: sprigString, description: "Repeats a string multiple times.", arguments: []string{"count", "str"}, aliases: []string{"repeatSprig"}},
	"replace":      {group: sprigString, description: "Performs simple string replacement.", arguments: []string{"old", "new", "src"}},
	"shuffle":      {group: sprigString, description: "Shuffle a string.", arguments: []string{"str"}},
	"snakecase":    {group: sprigString, description: "Converts string from camelCase to snake_case.", arguments: []string{"str"}},
	"squote":       {group: sprigString, description: "Wraps each argument with single quotes."},
	"substr":       {group: sprigString, description: "Get a substring from a string.", arguments: []string{"start", "length", "str"}},
	"swapcase":     {group: sprigString, description: "Swaps the uppercase to lowercase and lowercase to uppercase.", arguments: []string{"str"}},
	"title":        {group: sprigString, description: "Converts to title case.", arguments: []string{"str"}},
	"toString":     {group: sprigString, description: "Converts any value to string.", arguments: []string{"value"}},
	"trim":         {group: sprigString, description: "Removes space from either side of a string.", arguments: []string{"str"}},
	"trimAll":      {group: sprigString, description: "Removes given characters from the front or back of a string.", aliases: []string{"trimall"}, arguments: []string{"chars", "str"}},
	"trimPrefix":   {group: sprigString, description: "Trims just the prefix from a string if present.", arguments: []string{"prefix", "str"}},
	"trimSuffix":   {group: sprigString, description: "Trims just the suffix from a string if present.", arguments: []string{"suffix", "str"}},
	"trunc":        {group: sprigString, description: "Truncates a string (and add no suffix).", arguments: []string{"length", "str"}, aliases: []string{"truncSprig"}},
	"untitle":      {group: sprigString, description: `Removes title casing.`, arguments: []string{"str"}},
	"upper":        {group: sprigString, description: "Converts the entire string to uppercase.", arguments: []string{"str"}},
	"wrap":         {group: sprigString, description: "Wraps text at a given column count.", arguments: []string{"length", "str"}, aliases: []string{"wrapSprig"}},
	"wrapWith":     {group: sprigString, description: "Works as wrap, but lets you specify the string to wrap with (wrap uses \\n).", arguments: []string{"length", "spe", "str"}},

	"atoi":    {group: sprigTypeConversion},
	"int64":   {group: sprigTypeConversion},
	"int":     {group: sprigTypeConversion, aliases: []string{"intSprig"}},
	"float64": {group: sprigTypeConversion},

	"split":     {group: sprigStringList},
	"splitList": {group: sprigStringList},
	"toStrings": {group: sprigStringList},
	"join":      {group: sprigStringList},
	"sortAlpha": {group: sprigStringList},

	// VERY basic arithmetic.
	"add1":  {group: sprigMath},
	"add":   {group: sprigMath, aliases: []string{"addSprig"}},
	"sub":   {group: sprigMath, aliases: []string{"subSprig"}},
	"div":   {group: sprigMath, aliases: []string{"divSprig"}},
	"mod":   {group: sprigMath, aliases: []string{"modSprig"}},
	"mul":   {group: sprigMath, aliases: []string{"mulSprig"}},
	"max":   {group: sprigMath, aliases: []string{"maxSprig", "biggest"}},
	"min":   {group: sprigMath, aliases: []string{"minSprig"}},
	"ceil":  {group: sprigMath, aliases: []string{"ceilSprig"}},
	"floor": {group: sprigMath, aliases: []string{"floorSprig"}},
	"round": {group: sprigMath},

	"until":     {group: sprigMath, aliases: []string{"untilSprig"}},
	"untilStep": {group: sprigMath},

	// Defaults
	"default":      {group: sprigDefault},
	"empty":        {group: sprigDefault},
	"coalesce":     {group: sprigDefault},
	"compact":      {group: sprigDefault},
	"toJson":       {group: sprigDefault, aliases: []string{"toJsonSprig"}},
	"toPrettyJson": {group: sprigDefault, aliases: []string{"toJsonSprig"}},

	// Reflection
	"typeOf":     {group: sprigReflect},
	"typeIs":     {group: sprigReflect},
	"typeIsLike": {group: sprigReflect},
	"kindOf":     {group: sprigReflect},
	"kindIs":     {group: sprigReflect},

	// OS:
	"env":       {group: sprigOS},
	"expandenv": {group: sprigOS},

	// File Paths:
	"base":  {group: sprigFilePath},
	"dir":   {group: sprigFilePath},
	"clean": {group: sprigFilePath},
	"ext":   {group: sprigFilePath},
	"isAbs": {group: sprigFilePath},

	// Encoding:
	"b64enc": {group: sprigEncoding},
	"b64dec": {group: sprigEncoding},
	"b32enc": {group: sprigEncoding},
	"b32dec": {group: sprigEncoding},

	// Data Structures:
	"list":   {group: sprigDict, aliases: []string{"tuple"}},
	"dict":   {group: sprigDict},
	"set":    {group: sprigDict, aliases: []string{"setSprig"}},
	"unset":  {group: sprigDict},
	"hasKey": {group: sprigDict},
	"pluck":  {group: sprigDict},
	"keys":   {group: sprigDict},
	"pick":   {group: sprigDict, aliases: []string{"pickSprig"}},
	"omit":   {group: sprigDict, aliases: []string{"omitSprig"}},
	"merge":  {group: sprigDict, aliases: []string{"mergeSprig"}},

	// Lists functions
	"append":  {group: sprigList, aliases: []string{"push"}},
	"prepend": {group: sprigList},
	"first":   {group: sprigList},
	"rest":    {group: sprigList},
	"last":    {group: sprigList},
	"initial": {group: sprigList},
	"reverse": {group: sprigList},
	"uniq":    {group: sprigList},
	"without": {group: sprigList},
	"has":     {group: sprigList},

	// Cryptographics functions
	"sha256sum":         {group: sprigCrypto, description: "", arguments: []string{"input"}},
	"genPrivateKey":     {group: sprigCrypto},
	"derivePassword":    {group: sprigCrypto},
	"genCA":             {group: sprigCrypto},
	"genSelfSignedCert": {group: sprigCrypto},
	"genSignedCert":     {group: sprigCrypto},

	// UUIDs:
	"uuidv4": {group: sprigGen, aliases: []string{"uuid", "guid", "GUID"}},

	// SemVer:
	"semver":        {group: sprigSemver, description: "Parses a string into a Semantic Version.", arguments: []string{"version"}},
	"semverCompare": {group: sprigSemver, description: "A more robust comparison function is provided as semverCompare. This version supports version ranges.", arguments: []string{"constraints", "version"}},

	// Flow Control:
	"fail": {group: sprigFlow, description: "Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail."},

	// Regex
	"regexMatch":             {group: sprigRegex, description: "Returns true if the input string matches the regular expression.", arguments: []string{"regex", "str"}},
	"regexFindAll":           {group: sprigRegex, description: "Returns a slice of all matches of the regular expression in the input string.", arguments: []string{"regex", "str", "n"}},
	"regexFind":              {group: sprigRegex, description: "Returns the first (left most) match of the regular expression in the input string.", arguments: []string{"regex", "str"}},
	"regexReplaceAll":        {group: sprigRegex, description: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement. Inside string replacement, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.", arguments: []string{"regex", "str", "repl"}},
	"regexReplaceAllLiteral": {group: sprigRegex, description: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement The replacement string is substituted directly, without using Expand.", arguments: []string{"regex", "str", "repl"}},
	"regexSplit":             {group: sprigRegex, description: "Slices the input string into substrings separated by the expression and returns a slice of the substrings between those expression matches. The last parameter n determines the number of substrings to return, where -1 means return all matches.", arguments: []string{"regex", "str", "n"}},
}
