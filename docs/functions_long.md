{% include navigation.html %}
## Functions
```text
Truncates a string with ellipses (...).
abbrev(width int, str string) string

Abbreviates both sides with ellipses (...).
abbrevboth(left int, right int, str string) string

abs(arg1 interface{}) interface{}, error

acos(arg1 interface{}) interface{}, error

acosh(arg1 interface{}) interface{}, error

add(arg1 interface{}, args ...interface{}) interface{}, error

add1(arg1 interface{}) int64

The ago function returns duration from time.Now in seconds resolution.
ago(date interface{}) string

alias(arg1 string, arg2 string, arg3 interface{}, args ...interface{}) string, error

Returns the boolean AND of its arguments by returning the first empty argument or the last argument, that is, "and x y" behaves as "if x then y else x". All the arguments are evaluated.
and(arg0 reflect.Value, args ...reflect.Value) reflect.Value

append(arg1 interface{}, arg2 interface{}) []interface{}

Ensure that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).
array(value interface{}) interface{}

asin(arg1 interface{}) interface{}, error

asinh(arg1 interface{}) interface{}, error

atan(arg1 interface{}) interface{}, error

atan2(arg1 interface{}, arg2 interface{}) interface{}, error

atanh(arg1 interface{}) interface{}, error

atoi(arg1 string) int

avg(arg1 interface{}, args ...interface{}) interface{}, error

b32dec(arg1 string) string

b32enc(arg1 string) string

b64dec(arg1 string) string

b64enc(arg1 string) string

band(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

base(arg1 string) string

bclear(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)
bool(str string) bool, error

bor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

bxor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters. Thus "call .X.Y 1 2" is, in Go notation, dot.X.Y(1, 2) where Y is a func-valued field, map entry, or the like. The first argument must be the result of an evaluation that yields a value of function type (as distinct from a predefined function such as print). The function must return either one or two result values, the second of which is of type error. If the arguments don't match the function or the returned error value is non-nil, execution stops.
call(fn reflect.Value, args ...reflect.Value) reflect.Value, error)

Converts string from snake_case to CamelCase.
camelcase(str string) string

Concatenates multiple strings together into one, separating them with spaces.
cat(args ...interface{}) string

ceil(arg1 interface{}) interface{}, error

center(width interface{}, str interface{}) string

Returns the character corresponging to the supplied integer value
char(value interface{}) interface{}, error

clean(arg1 string) string

coalesce(args ...interface{}) interface{}

color(args ...interface{}) string, error

compact(arg1 interface{}) []interface{}

concat(args ...interface{}) string

Tests to see if one string is contained inside of another.
contains(substr string, str string) bool

Returns the content of a single element map (used to retrieve content in a declaration like `value "name" { a = 1 b = 3}`)
content(keymap interface{}) interface{}, error

cos(arg1 interface{}) interface{}, error

cosh(arg1 interface{}) interface{}, error

current() string

Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.
data(data interface{}, context ...interface{}) interface{}, error

The date function formats a dat (https://golang.org/pkg/time/#Time.Format).
date(fmt string, date interface{}) string

Same as date, but with a timezone.
dateInZone(fmt string, date interface{}, zone string) string

The dateModify takes a modification and a date and returns the timestamp.
dateModify(fmt string, date time.Time) time.Time

debug(args ...interface{}) string

debugf(format string, args ...interface{}) string

dec(arg1 interface{}) interface{}, error

default(arg1 interface{}, args ...interface{}) interface{}

deg(arg1 interface{}) interface{}, error

derivePassword(arg1 uint32, arg2 string, arg3 string, arg4 string, arg5 string) string

dict(args ...interface{}) map[string]interface{}

diff(arg1 interface{}, arg2 interface{}) interface{}

dim(arg1 interface{}, arg2 interface{}) interface{}, error

dir(arg1 string) string

div(arg1 interface{}, arg2 interface{}) interface{}, error

ellipsis(arg1 string, args ...interface{}) interface{}, error

empty(arg1 interface{}) bool

env(arg1 string) string

Returns the boolean truth of arg1 == arg2
eq(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

error(args ...interface{}) string

errorf(format string, args ...interface{}) string

exec(arg1 interface{}, args ...interface{}) interface{}, error

exit(arg1 int) int

exp(arg1 interface{}) interface{}, error

exp2(arg1 interface{}) interface{}, error

expandenv(arg1 string) string

expm1(arg1 interface{}) interface{}, error

ext(arg1 string) string

Extract values from a slice or a map, indexes could be either integers for slice or strings for maps
extract(source interface{}, indexes ...interface{}) interface{}, error

Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail.
fail(arg1 string) string, error

fatal(args ...interface{}) string

fatalf(format string, args ...interface{}) string

first(arg1 interface{}) interface{}

float64(arg1 interface{}) float64

floor(arg1 interface{}) interface{}, error

formatList(format string, list interface{}) []string

frexp(arg1 interface{}) interface{}, error

func(arg1 string, arg2 string, arg3 interface{}, arg4 interface{}, arg5 interface{}) string, error

Returns the information relative to a specific function
function(name string) template.FuncInfo

functions() []string

gamma(arg1 interface{}) interface{}, error

Returns the boolean truth of arg1 >= arg2
ge(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

genCA(arg1 string, arg2 int) sprig.certificate, error

genPrivateKey(arg1 string) string

genSelfSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int) sprig.certificate, error

genSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int, arg5 sprig.certificate) sprig.certificate, error

Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)
get(map interface{}, key interface{}) interface{}, error

glob(args ...interface{}) []string

Returns the boolean truth of arg1 > arg2
gt(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

has(arg1 interface{}, arg2 interface{}) bool

hasKey(arg1 map[string]interface{}, arg2 string) bool

Tests whether a string has a given prefix.
hasPrefix(prefix string, str string) bool

Tests whether a string has a given suffix.
hasSuffix(suffix string, str string) bool

Converts the supplied hcl string into data structure (Go spec). If context is omitted, default context is used.
hcl(hcl interface{}, context ...interface{}) interface{}, error

Simple hello by Sprig
hello() string

hex(arg1 interface{}) interface{}, error

Returns the escaped HTML equivalent of the textual representation of its arguments. This function is unavailable in html/template, with a few exceptions.
html(args ...interface{}) string

The htmlDate function formates a date for inserting into an HTML date picker input field.
htmlDate(date interface{}) string

Same as htmlDate, but with a timezone.
htmlDateInZone(date interface{}, zone string) string

hypot(arg1 interface{}, arg2 interface{}) interface{}, error

id(identifier string, replaceChar ...interface{}) string

iif(test interface{}, valueIfTrue interface{}, valueIfFalse interface{}) interface{}

ilogb(arg1 interface{}) interface{}, error

include(arg1 interface{}, args ...interface{}) interface{}, error

Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.
indent(spaces int, str string) string

Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.
index(item reflect.Value, indices ...reflect.Value) (reflect.Value, error)

info(args ...interface{}) string

infof(format string, args ...interface{}) string

initial(arg1 interface{}) []interface{}

Given multiple words, takes the first letter of each word and combine.
initials(str string) string

int64(arg1 interface{}) int64

isAbs(arg1 string) bool

isInf(arg1 interface{}, arg2 interface{}) interface{}, error

isNaN(arg1 interface{}) interface{}, error

join(arg1 string, arg2 interface{}) string

joinLines(args ...interface{}) string

Returns the escaped JavaScript equivalent of the textual representation of its arguments.
js(args ...interface{}) string

Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.
json(json interface{}, context ...interface{}) interface{}, error

key(arg1 interface{}) interface{}, error

keys(arg1 map[string]interface{}) []string

kindIs(arg1 string, arg2 interface{}) bool

kindOf(arg1 interface{}) string

last(arg1 interface{}) interface{}

ldexp(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the boolean truth of arg1 <= arg2
le(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

Returns the integer length of its argument.
len(item interface{}) (int, error)

Returns the number of actual character in a string
lenc(arg1 string) int

lgamma(arg1 interface{}) interface{}, error

list(args ...interface{}) []interface{}

localAlias(arg1 string, arg2 string, arg3 interface{}, args ...interface{}) string, error

log(arg1 interface{}) interface{}, error

log10(arg1 interface{}) interface{}, error

log1p(arg1 interface{}) interface{}, error

log2(arg1 interface{}) interface{}, error

logb(arg1 interface{}) interface{}, error

lorem(funcName interface{}, args ...int) string, error

Converts the entire string to lowercase.
lower(str string) string

lshift(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the boolean truth of arg1 < arg2
lt(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

max(args ...interface{}) interface{}

merge(arg1 map[string]interface{}, args ...map[string]interface{}) map[string]interface{}, error

mergeList(lists ...[]interface{}) []interface{}

min(args ...interface{}) interface{}

mod(arg1 interface{}, arg2 interface{}) interface{}, error

modf(arg1 interface{}) interface{}, error

mul(arg1 interface{}, args ...interface{}) interface{}, error

Returns the boolean truth of arg1 != arg2
ne(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

nextAfter(arg1 interface{}, arg2 interface{}) interface{}, error

Same as the indent function, but prepends a new line to the beginning of the string.
nindent(spaces int, str string) string

Removes all whitespace from a string.
nospace(str string) string

Returns the boolean negation of its single argument.
not(not(arg reflect.Value) bool

notice(args ...interface{}) string

noticef(format string, args ...interface{}) string

The current date/time. Use this in conjunction with other date functions.
now() time.Time

omit(arg1 map[string]interface{}, args ...interface{}) map[string]interface{}

Returns the boolean OR of its arguments by returning the first non-empty argument or the last argument, that is, "or x y" behaves as "if x then x else y". All the arguments are evaluated.
or(or(arg0 reflect.Value, args ...reflect.Value) reflect.Value

pick(arg1 map[string]interface{}, args ...interface{}) map[string]interface{}

pickv(arg1 map[string]interface{}, arg2 string, args ...interface{}) map[string]interface{}, error

pluck(arg1 string, args ...map[string]interface{}) []interface{}

Pluralizes a string.
plural(one string, many string, count int) string

pow(arg1 interface{}, arg2 interface{}) interface{}, error

pow10(arg1 interface{}) interface{}, error

prepend(arg1 interface{}, arg2 interface{}) []interface{}

An alias for fmt.Sprint
print(args ...interface{}) string

An alias for fmt.Sprintf
printf(format string, args ...interface{}) string

An alias for fmt.Sprintln
println(args ...interface{}) string

Returns the current working directory
pwd() string

Wraps each argument with double quotes.
quote(str ...interface{}) string

rad(arg1 interface{}) interface{}, error

Generates random string with letters.
randAlpha(count int) string

Generates random string with letters and digits.
randAlphaNum(count int) string

Generates random string with ASCII printable characters.
randAscii(count int) string

Generates random string with digits.
randNumeric(count int) string

Returns the first (left most) match of the regular expression in the input string.
regexFind(regex string, str string) string

Returns a slice of all matches of the regular expression in the input string.
regexFindAll(regex string, str string, n int) []string

Returns true if the input string matches the regular expression.
regexMatch(regex string, str string) bool

Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement. Inside string replacement, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.
regexReplaceAll(regex string, str string, repl string) string

Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement The replacement string is substituted directly, without using Expand.
regexReplaceAllLiteral(regex string, str string, repl string) string

Slices the input string into substrings separated by the expression and returns a slice of the substrings between those expression matches. The last parameter n determines the number of substrings to return, where -1 means return all matches.
regexSplit(regex string, str string, n int) []string

rem(arg1 interface{}, arg2 interface{}) interface{}, error

Returns an array with the item repeated n times.
repeat(n int, item interface{}) []interface{}, error

Performs simple string replacement.
replace(old string, new string, src string) string

rest(arg1 interface{}) []interface{}

reverse(arg1 interface{}) []interface{}

round(arg1 interface{}, arg2 int, args ...float64) float64

rshift(arg1 interface{}, arg2 interface{}) interface{}, error

run(arg1 interface{}, args ...interface{}) interface{}, error

safeIndex(arg1 interface{}, arg2 int, arg3 interface{}) interface{}, error

Parses a string into a Semantic Version.
semver(version string) *semver.Version, error

A more robust comparison function is provided as semverCompare. This version supports version ranges.
semverCompare(constraints string, version string) bool, error

set(arg1 interface{}, arg2 interface{}, arg3 interface{}) string, error

sha256sum(input string) string

Shuffle a string.
shuffle(str string) string

signBit(arg1 interface{}, arg2 interface{}) interface{}, error

sin(arg1 interface{}) interface{}, error

sincos(arg1 interface{}) interface{}, error

sinh(arg1 interface{}) interface{}, error

slice(arg1 interface{}, args ...interface{}) interface{}, error

Converts string from camelCase to snake_case.
snakecase(str string) string

sortAlpha(arg1 interface{}) []string

split(arg1 string, arg2 string) map[string]string

splitLines(arg1 interface{}) []interface{}

splitList(arg1 string, arg2 string) []string

sqrt(arg1 interface{}) interface{}, error

Wraps each argument with single quotes.
squote(args ...interface{}) string

string(arg1 interface{}) utils.String

sub(arg1 interface{}, arg2 interface{}) interface{}, error

substitute(arg1 string) string

Get a substring from a string.
substr(start int, length int, str string) string

Swaps the uppercase to lowercase and lowercase to uppercase.
swapcase(str string) string

tan(arg1 interface{}) interface{}, error

tanh(arg1 interface{}) interface{}, error

templateNames() []*template.Template

Converts to title case.
title(str string) string

to(args ...interface{}) interface{}, error

Convert the supplied value to bash compatible representation.
toBash(value interface{}) string

Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.
toDate(fmt string, str string) time.Time

Convert the supplied value to compact HCL representation.
toHcl(value interface{}) string, error

Convert the supplied value to compact JSON representation.
toJson(value interface{}) string, error

Convert the supplied value to pretty HCL representation.
toPrettyHcl(value interface{}) string, error

Convert the supplied value to pretty JSON representation.
toPrettyJson(value interface{}) string, error

Convert the supplied value to pretty HCL representation (without multiple map declarations).
toPrettyTFVars(value interface{}) string, error

Convert the supplied value to compact quoted HCL representation.
toQuotedHcl(value interface{}) string, error

Convert the supplied value to compact quoted JSON representation.
toQuotedJson(value interface{}) string, error

Convert the supplied value to compact HCL representation (without multiple map declarations).
toQuotedTFVars(value interface{}) string, error

Converts any value to string.
toString(value interface{}) string

toStrings(arg1 interface{}) []string

Convert the supplied value to compact HCL representation (without multiple map declarations).
toTFVars(value interface{}) string, error

Convert the supplied value to YAML representation.
toYaml(value interface{}) string, error

Removes space from either side of a string.
trim(str string) string

Removes given characters from the front or back of a string.
trimAll(chars string, str string) string

Trims just the prefix from a string if present.
trimPrefix(prefix string, str string) string

Trims just the suffix from a string if present.
trimSuffix(suffix string, str string) string

trunc(arg1 interface{}) interface{}, error

typeIs(arg1 string, arg2 interface{}) bool

typeIsLike(arg1 string, arg2 interface{}) bool

typeOf(arg1 interface{}) string

undef(arg1 interface{}, args ...interface{}) interface{}

uniq(arg1 interface{}) []interface{}

unset(arg1 map[string]interface{}, arg2 string) map[string]interface{}

until(args ...interface{}) interface{}, error

untilStep(arg1 int, arg2 int, arg3 int) []int

Removes title casing.
untitle(str string) string

Converts the entire string to uppercase.
upper(str string) string

Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query. This function is unavailable in html/template, with a few exceptions.
urlquery(args ...interface{}) string

uuidv4() string

warning(args ...interface{}) string

warningf(format string, args ...interface{}) string

without(arg1 interface{}, args ...interface{}) []interface{}

wrap(width interface{}, s interface{}) string

Works as wrap, but lets you specify the string to wrap with (wrap uses \n).
wrapWith(length int, spe string, str string) string

Converts the supplied yaml string into data structure (Go spec). If context is omitted, default context is used.
yaml(yaml interface{}, context ...interface{}) interface{}, error
```
