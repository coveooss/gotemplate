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
			sprigFuncs[key] = FuncInfo{f: value, group: info.group, aliases: info.aliases, args: info.args, desc: info.desc}
		}
	}
	t.AddFunctions(sprigFuncs)
}

var sprigFuncRef = map[string]struct {
	args, aliases []string
	group, desc   string
}{
	"hello": {group: sprigGen, desc: "Simple hello by Sprig"},
	// Date functions
	"date":           {group: sprigDate, desc: "The date function formats a dat (https://golang.org/pkg/time/#Time.Format).", args: []string{"fmt", "date"}},
	"now":            {group: sprigDate, desc: "The current date/time. Use this in conjunction with other date functions."},
	"htmlDate":       {group: sprigDate, desc: "The htmlDate function formates a date for inserting into an HTML date picker input field.", args: []string{"date"}},
	"htmlDateInZone": {group: sprigDate, desc: "Same as htmlDate, but with a timezone.", args: []string{"date", "zone"}},
	"dateInZone":     {group: sprigDate, desc: "Same as date, but with a timezone.", args: []string{"fmt", "date", "zone"}, aliases: []string{"date_in_zone"}},
	"dateModify":     {group: sprigDate, desc: "The dateModify takes a modification and a date and returns the timestamp.", args: []string{"fmt", "date"}, aliases: []string{"date_modify"}},
	"ago":            {group: sprigDate, desc: "The ago function returns duration from time.Now in seconds resolution.", args: []string{"date"}},
	"toDate":         {group: sprigDate, desc: "Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.", args: []string{"fmt", "str"}},
	// Strings functions
	"abbrev":       {group: sprigString, desc: "Truncates a string with ellipses (...).", args: []string{"width", "str"}},
	"abbrevboth":   {group: sprigString, desc: "Abbreviates both sides with ellipses (...).", args: []string{"left", "right", "str"}},
	"camelcase":    {group: sprigString, desc: "Converts string from snake_case to CamelCase.", args: []string{"str"}},
	"cat":          {group: sprigString, desc: "Concatenates multiple strings together into one, separating them with spaces."},
	"contains":     {group: sprigString, desc: "Tests to see if one string is contained inside of another.", args: []string{"substr", "str"}},
	"hasPrefix":    {group: sprigString, desc: "Tests whether a string has a given prefix.", args: []string{"prefix", "str"}},
	"hasSuffix":    {group: sprigString, desc: "Tests whether a string has a given suffix.", args: []string{"suffix", "str"}},
	"indent":       {group: sprigString, desc: "Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.", args: []string{"spaces", "str"}},
	"initials":     {group: sprigString, desc: "Given multiple words, takes the first letter of each word and combine.", args: []string{"str"}},
	"lower":        {group: sprigString, desc: "Converts the entire string to lowercase.", args: []string{"str"}},
	"nindent":      {group: sprigString, desc: "Same as the indent function, but prepends a new line to the beginning of the string.", args: []string{"spaces", "str"}},
	"nospace":      {group: sprigString, desc: "Removes all whitespace from a string.", args: []string{"str"}},
	"plural":       {group: sprigString, desc: "Pluralizes a string.", args: []string{"one", "many", "count"}},
	"quote":        {group: sprigString, desc: "Wraps each argument with double quotes.", args: []string{"str"}},
	"randAlpha":    {group: sprigString, desc: "Generates random string with letters.", args: []string{"count"}},
	"randAlphaNum": {group: sprigString, desc: "Generates random string with letters and digits.", args: []string{"count"}},
	"randAscii":    {group: sprigString, desc: "Generates random string with ASCII printable characters.", args: []string{"count"}},
	"randNumeric":  {group: sprigString, desc: "Generates random string with digits.", args: []string{"count"}},
	"repeat":       {group: sprigString, desc: "Repeats a string multiple times.", args: []string{"count", "str"}, aliases: []string{"repeatSprig"}},
	"replace":      {group: sprigString, desc: "Performs simple string replacement.", args: []string{"old", "new", "src"}},
	"shuffle":      {group: sprigString, desc: "Shuffle a string.", args: []string{"str"}},
	"snakecase":    {group: sprigString, desc: "Converts string from camelCase to snake_case.", args: []string{"str"}},
	"squote":       {group: sprigString, desc: "Wraps each argument with single quotes."},
	"substr":       {group: sprigString, desc: "Get a substring from a string.", args: []string{"start", "length", "str"}},
	"swapcase":     {group: sprigString, desc: "Swaps the uppercase to lowercase and lowercase to uppercase.", args: []string{"str"}},
	"title":        {group: sprigString, desc: "Converts to title case.", args: []string{"str"}},
	"toString":     {group: sprigString, desc: "Converts any value to string.", args: []string{"value"}},
	"trim":         {group: sprigString, desc: "Removes space from either side of a string.", args: []string{"str"}},
	"trimAll":      {group: sprigString, desc: "Removes given characters from the front or back of a string.", aliases: []string{"trimall"}, args: []string{"chars", "str"}},
	"trimPrefix":   {group: sprigString, desc: "Trims just the prefix from a string if present.", args: []string{"prefix", "str"}},
	"trimSuffix":   {group: sprigString, desc: "Trims just the suffix from a string if present.", args: []string{"suffix", "str"}},
	"trunc":        {group: sprigString, desc: "Truncates a string (and add no suffix).", args: []string{"length", "str"}, aliases: []string{"truncSprig"}},
	"untitle":      {group: sprigString, desc: `Removes title casing.`, args: []string{"str"}},
	"upper":        {group: sprigString, desc: "Converts the entire string to uppercase.", args: []string{"str"}},
	"wrap":         {group: sprigString, desc: "Wraps text at a given column count.", args: []string{"length", "str"}, aliases: []string{"wrapSprig"}},
	"wrapWith":     {group: sprigString, desc: "Works as wrap, but lets you specify the string to wrap with (wrap uses \\n).", args: []string{"length", "spe", "str"}},

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
	"sha256sum":         {group: sprigCrypto, desc: "", args: []string{"input"}},
	"genPrivateKey":     {group: sprigCrypto},
	"derivePassword":    {group: sprigCrypto},
	"genCA":             {group: sprigCrypto},
	"genSelfSignedCert": {group: sprigCrypto},
	"genSignedCert":     {group: sprigCrypto},

	// UUIDs:
	"uuidv4": {group: sprigGen, aliases: []string{"uuid", "guid", "GUID"}},

	// SemVer:
	"semver":        {group: sprigSemver, desc: "Parses a string into a Semantic Version.", args: []string{"version"}},
	"semverCompare": {group: sprigSemver, desc: "A more robust comparison function is provided as semverCompare. This version supports version ranges.", args: []string{"constraints", "version"}},

	// Flow Control:
	"fail": {group: sprigFlow, desc: "Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail."},

	// Regex
	"regexMatch":             {group: sprigRegex, desc: "Returns true if the input string matches the regular expression.", args: []string{"regex", "str"}},
	"regexFindAll":           {group: sprigRegex, desc: "Returns a slice of all matches of the regular expression in the input string.", args: []string{"regex", "str", "n"}},
	"regexFind":              {group: sprigRegex, desc: "Returns the first (left most) match of the regular expression in the input string.", args: []string{"regex", "str"}},
	"regexReplaceAll":        {group: sprigRegex, desc: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement. Inside string replacement, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.", args: []string{"regex", "str", "repl"}},
	"regexReplaceAllLiteral": {group: sprigRegex, desc: "Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement The replacement string is substituted directly, without using Expand.", args: []string{"regex", "str", "repl"}},
	"regexSplit":             {group: sprigRegex, desc: "Slices the input string into substrings separated by the expression and returns a slice of the substrings between those expression matches. The last parameter n determines the number of substrings to return, where -1 means return all matches.", args: []string{"regex", "str", "n"}},
}
