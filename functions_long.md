```text
Returns a String class object that allows invoking standard string operations as method.
String(value interface{}) utils.String

Truncates a string with ellipses (...).
abbrev(width int, str string) string

Abbreviates both sides with ellipses (...).
abbrevboth(left int, right int, str string) string

Returns the absolute value of x.
Special cases are:
    abs(±Inf) = +Inf
    abs(NaN) = NaN
abs(x interface{}) interface{}, error

Returns the arccosine, in radians, of x.
Special case is:
    acos(x) = NaN if x < -1 or x > 1
acos(x interface{}) interface{}, error

Returns the inverse hyperbolic cosine of x.
Special cases are:
    acosh(+Inf) = +Inf
    acosh(x) = NaN if x < 1
    acosh(NaN) = NaN
acosh(x interface{}) interface{}, error

add(arg1 interface{}, args ...interface{}) interface{}, error

add1(arg1 interface{}) int64

The ago function returns duration from time.Now in seconds resolution.
ago(date interface{}) string

Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the caller.
alias(name string, function string, source interface{}, args ...interface{}) string, error

Returns the list of all functions that are simply an alias of another function.
aliases() []string

Returns the list of all available functions.
allFunctions() []string

Returns the boolean AND of its arguments by returning the first empty argument or the last argument, that is, "and x y" behaves as "if x then y else x". All the arguments are evaluated.
and(arg0 reflect.Value, args ...reflect.Value) reflect.Value

append(arg1 interface{}, arg2 interface{}) []interface{}

Ensures that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).
array(value interface{}) interface{}

Returns the arcsine, in radians, of x.
Special cases are:
    asin(±0) = ±0
    asin(x) = NaN if x < -1 or x > 1
asin(x interface{}) interface{}, error

Returns the inverse hyperbolic sine of x.
Special cases are:
    asinh(±0) = ±0
    asinh(±Inf) = ±Inf
    asinh(NaN) = NaN
asinh(x interface{}) interface{}, error

Returns the arctangent, in radians, of x.
Special cases are:
    atan(±0) = ±0
    atan(±Inf) = ±Pi/2
atan(x interface{}) interface{}, error

Returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
Special cases are (in order):
    atan2(y, NaN) = NaN
    atan2(NaN, x) = NaN
    atan2(+0, x>=0) = +0
    atan2(-0, x>=0) = -0
    atan2(+0, x<=-0) = +Pi
    atan2(-0, x<=-0) = -Pi
    atan2(y>0, 0) = +Pi/2
    atan2(y<0, 0) = -Pi/2
    atan2(+Inf, +Inf) = +Pi/4
    atan2(-Inf, +Inf) = -Pi/4
    atan2(+Inf, -Inf) = 3Pi/4
    atan2(-Inf, -Inf) = -3Pi/4
    atan2(y, +Inf) = 0
    atan2(y>0, -Inf) = +Pi
    atan2(y<0, -Inf) = -Pi
    atan2(+Inf, x) = +Pi/2
    atan2(-Inf, x) = -Pi/2
atan2(x interface{}, y interface{}) interface{}, error

Returns the inverse hyperbolic tangent of x.
Special cases are:
    atanh(1) = +Inf
    atanh(±0) = ±0
    atanh(-1) = -Inf
    atanh(x) = NaN if x < -1 or x > 1
    atanh(NaN) = NaN
atanh(x interface{}) interface{}, error

atoi(arg1 string) int

avg(arg1 interface{}, args ...interface{}) interface{}, error

b32dec(arg1 string) string

b32enc(arg1 string) string

b64dec(arg1 string) string

b64enc(arg1 string) string

band(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

base(arg1 string) string

bclear(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

Converts the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)
bool(str string) bool, error

bor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

bxor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error

Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters. Thus "call .X.Y 1 2" is, in Go notation, dot.X.Y(1, 2) where Y is a func-valued field, map entry, or the like. The first argument must be the result of an evaluation that yields a value of function type (as distinct from a predefined function such as print). The function must return either one or two result values, the second of which is of type error. If the arguments don't match the function or the returned error value is non-nil, execution stops.
call(fn reflect.Value, args ...reflect.Value) reflect.Value, error)

Converts string from snake_case to CamelCase.
camelcase(str string) string

Concatenates multiple strings together into one, separating them with spaces.
cat(args ...interface{}) string

Returns the cube root of x.
Special cases are:
    cbrt(±0) = ±0
    cbrt(±Inf) = ±Inf
    cbrt(NaN) = NaN
cbrt(x interface{}) interface{}, error

Returns the least integer value greater than or equal to x.
Special cases are:
    ceil(±0) = ±0
    ceil(±Inf) = ±Inf
    ceil(NaN) = NaN
ceil(x interface{}) interface{}, error

Returns the concatenation of supplied arguments centered within width.
center(width interface{}, args ...interface{}) string

Returns the character corresponging to the supplied integer value
char(value interface{}) interface{}, error

clean(arg1 string) string

coalesce(args ...interface{}) interface{}

Colors the rendered string.

The first arguments are interpretated as color attributes until the first non color attribute. Attributes are case insensitive.

Valid attributes are:
    Reset, Bold, Faint, Italic, Underline, BlinkSlow, BlinkRapid, ReverseVideo, Concealed, CrossedOut

Valid color are:
    Black, Red, Green, Yellow, Blue, Magenta, Cyan, White

Color can be prefixed by:
    Fg:   Meaning foreground (Fg is assumed if not specified)
    FgHi: Meaning high intensity forground
    Bg:   Meaning background"
    BgHi: Meaning high intensity background
color(args ...interface{}) string, error

compact(arg1 interface{}) []interface{}

Returns the concatenation (without separator) of the string representation of objects.
concat(args ...interface{}) string

Tests to see if one string is contained inside of another.
contains(substr string, str string) bool

Returns the content of a single element map (used to retrieve content in a declaration like `value "name" { a = 1 b = 3}`)
content(keymap interface{}) interface{}, error

Returns the cosine of the radian argument x.
Special cases are:
    cos(±Inf) = NaN
    cos(NaN) = NaN
cos(x interface{}) interface{}, error

Returns the hyperbolic cosine of x.
Special cases are:
    cosh(±0) = 1
    cosh(±Inf) = +Inf
    cosh(NaN) = NaN
cosh(x interface{}) interface{}, error

logs a message using CRITICAL as log level (0).
critical(args ...interface{}) string

logs a message with format string using CRITICAL as log level (0).
criticalf(format string, args ...interface{}) string

Returns the current folder (like pwd, but returns the folder of the currently running folder).
current() string

Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.
data(data interface{}, context ...interface{}) interface{}, error

The date function formats a dat (https://golang.org/pkg/time/#Time.Format).
date(fmt string, date interface{}) string

Same as date, but with a timezone.
dateInZone(fmt string, date interface{}, zone string) string

The dateModify takes a modification and a date and returns the timestamp.
dateModify(fmt string, date time.Time) time.Time

logs a message using DEBUG as log level (5).
debug(args ...interface{}) string

logs a message with format using DEBUG as log level (5).
debugf(format string, args ...interface{}) string

dec(arg1 interface{}) interface{}, error

default(arg1 interface{}, args ...interface{}) interface{}

deg(arg1 interface{}) interface{}, error

derivePassword(arg1 uint32, arg2 string, arg3 string, arg4 string, arg5 string) string

dict(args ...interface{}) map[string]interface{}

Returns a colored string that highlight differences between supplied texts.
diff(text1 interface{}, text2 interface{}) interface{}

Returns the maximum of x-y or 0.
Special cases are:
    dim(+Inf, +Inf) = NaN
    dim(-Inf, -Inf) = NaN
    dim(x, NaN) = dim(NaN, x) = NaN
dim(x interface{}, y interface{}) interface{}, error

dir(arg1 string) string

div(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the result of the function by expanding its last argument that must be an array into values. It's like calling function(arg1, arg2, otherArgs...).
ellipsis(function string, args ...interface{}) interface{}, error

empty(arg1 interface{}) bool

env(arg1 string) string

Returns the boolean truth of arg1 == arg2
eq(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

logs a message using ERROR as log level (1).
error(args ...interface{}) string

logs a message with format using ERROR as log level (1).
errorf(format string, args ...interface{}) string

Returns the result of the shell command as structured data (as string if no other conversion is possible).
exec(command interface{}, args ...interface{}) interface{}, error

Exits the current program execution.
exit(exitValue int) int

Returns e**x, the base-e exponential of x.
Special cases are:
    exp(+Inf) = +Inf
    exp(NaN) = NaN
Very large values overflow to 0 or +Inf. Very small values underflow to 1.
exp(x interface{}) interface{}, error

Returns 2**x, the base-2 exponential of x.
Special cases are the same as exp.
exp2(x interface{}) interface{}, error

expandenv(arg1 string) string

Returns e**x - 1, the base-e exponential of x minus 1. It is more
accurate than exp(x) - 1 when x is near zero.
Special cases are:
    expm1(+Inf) = +Inf
    expm1(-Inf) = -1
    expm1(NaN) = NaN
Very large values overflow to -1 or +Inf
expm1(x interface{}) interface{}, error

ext(arg1 string) string

Extracts values from a slice or a map, indexes could be either integers for slice or strings for maps
extract(source interface{}, indexes ...interface{}) interface{}, error

Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail.
fail(arg1 string) string, error

Equivalents to critical followed by a call to os.Exit(1).
fatal(args ...interface{}) string

Equivalents to criticalf followed by a call to os.Exit(1).
fatalf(format string, args ...interface{}) string

first(arg1 interface{}) interface{}

float64(arg1 interface{}) float64

Returns the greatest integer value less than or equal to x.
Special cases are:
    floor(±0) = ±0
    floor(±Inf) = ±Inf
    floor(NaN) = NaN
floor(x interface{}) interface{}, error

Return a list of strings by applying the format to each element of the supplied list.
formatList(format string, list interface{}) []string

Breaks f into a normalized fraction and an integral power of two. Returns frac and exp satisfying f == frac × 2**exp, with the absolute value of frac in the interval [½, 1).
Special cases are:
    frexp(±0) = ±0, 0
    frexp(±Inf) = ±Inf, 0
    frexp(NaN) = NaN, 0
frexp(f interface{}) interface{}, error

Defines a function with the current context using the function (exec, run, include, template). Executed in the context of the caller.
func(name string, function string, source interface{}, config interface{}) string, error

Returns the information relative to a specific function.
function(name string) template.FuncInfo

Returns the list of all available functions (excluding aliases).
functions() []string

Returns the Gamma function of x.
Special cases are:
    gamma(+Inf) = +Inf
    gamma(+0) = +Inf
    gamma(-0) = -Inf
    gamma(x) = NaN for integer x < 0
    gamma(-Inf) = NaN
    gamma(NaN) = NaN
gamma(x interface{}) interface{}, error

Returns the boolean truth of arg1 >= arg2
ge(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

genCA(arg1 string, arg2 int) sprig.certificate, error

genPrivateKey(arg1 string) string

genSelfSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int) sprig.certificate, error

genSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int, arg5 sprig.certificate) sprig.certificate, error

Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)
get(map interface{}, key interface{}) interface{}, error

Returns the expanded list of supplied arguments (expand *[]? on filename).
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

Returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.
Special cases are:
    hypot(±Inf, q) = +Inf
    hypot(p, ±Inf) = +Inf
    hypot(NaN, q) = NaN
    hypot(p, NaN) = NaN
hypot(p interface{}, q interface{}) interface{}, error

Returns a valid go identifier from the supplied string (replacing any non compliant character by replacement, default _ ).
id(identifier string, replaceChar ...interface{}) string

If testValue is empty, returns falseValue, otherwise returns trueValue.
    WARNING: All arguments are evaluated and must by valid.
iif(testValue interface{}, valueTrue interface{}, valueFalse interface{}) interface{}

Returns the binary exponent of x as an integer.
Special cases are:
    ilogb(±Inf) = MaxInt32
    ilogb(0) = MinInt32
    ilogb(NaN) = MaxInt32
ilogb(x interface{}) interface{}, error

Returns the result of the named template rendering (like template but it is possible to capture the output).
include(source interface{}, context ...interface{}) interface{}, error

Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.
indent(spaces int, str string) string

Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.
index(item reflect.Value, indices ...reflect.Value) (reflect.Value, error)

logs a message using INFO as log level (4).
info(args ...interface{}) string

logs a message with format using INFO as log level (4).
infof(format string, args ...interface{}) string

initial(arg1 interface{}) []interface{}

Given multiple words, takes the first letter of each word and combine.
initials(str string) string

int64(arg1 interface{}) int64

isAbs(arg1 string) bool

Reports whether f is an infinity, according to sign. If sign > 0, isInf reports whether f is positive infinity. If sign < 0, IsInf reports whether f is negative infinity. If sign == 0, IsInf reports whether f is either infinity
isInf(f interface{}, arg2 interface{}) interface{}, error

Reports whether f is an IEEE 754 'not-a-number' value
isNaN(f interface{}) interface{}, error

Returns the order-zero Bessel function of the first kind.
Special cases are:
    j0(±Inf) = 0
    j0(0) = 1
    j0(NaN) = NaN
j0(x interface{}) interface{}, error

Returns the order-one Bessel function of the first kind.
Special cases are:
    j1(±Inf) = 0
    j1(NaN) = NaN
j1(x interface{}) interface{}, error

Returns the order-n Bessel function of the first kind.
Special cases are:
    jn(n, ±Inf) = 0
    jn(n, NaN) = NaN
jn(n interface{}, x interface{}) interface{}, error

join(arg1 string, arg2 interface{}) string

Merge the supplied objects into a newline separated string.
joinLines(format ...interface{}) string

Returns the escaped JavaScript equivalent of the textual representation of its arguments.
js(args ...interface{}) string

Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.
json(json interface{}, context ...interface{}) interface{}, error

Returns the key name of a single element map (used to retrieve name in a declaration like `value "name" { a = 1 b = 3}`)
key(value interface{}) interface{}, error

keys(arg1 map[string]interface{}) []string

kindIs(arg1 string, arg2 interface{}) bool

kindOf(arg1 interface{}) string

last(arg1 interface{}) interface{}

Ldexp is the inverse of Frexp. Returns frac × 2**exp.
Special cases are:
    ldexp(±0, exp) = ±0
    ldexp(±Inf, exp) = ±Inf
    ldexp(NaN, exp) = NaN
ldexp(frac interface{}, exp interface{}) interface{}, error

Returns the boolean truth of arg1 <= arg2
le(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

Returns the integer length of its argument.
len(item interface{}) (int, error)

Returns the number of actual character in a string
lenc(str string) int

Returns the natural logarithm and sign (-1 or +1) of Gamma(x).
Special cases are:
    lgamma(+Inf) = +Inf
    lgamma(0) = +Inf
    lgamma(-integer) = +Inf
    lgamma(-Inf) = -Inf
    lgamma(NaN) = NaN
lgamma(x interface{}) interface{}, error

list(args ...interface{}) []interface{}

Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the function it maps to.
localAlias(name string, function string, source interface{}, args ...interface{}) string, error

Returns the natural logarithm of x.
Special cases are:
    log(+Inf) = +Inf
    log(0) = -Inf
    log(x < 0) = NaN
    log(NaN) = NaN
log(x interface{}) interface{}, error

Returns the decimal logarithm of x. The special cases are the same as for log.
log10(x interface{}) interface{}, error

Returns the natural logarithm of 1 plus its argument x. It is more accurate than log(1 + x) when x is near zero.
Special cases are:
    log1p(+Inf) = +Inf
    log1p(±0) = ±0
    log1p(-1) = -Inf
    log1p(x < -1) = NaN
    log1p(NaN) = NaN
log1p(x interface{}) interface{}, error

Returns the binary logarithm of x. The special cases are the same as for log.
log2(x interface{}) interface{}, error

Returns the binary exponent of x.
Special cases are:
    logb(±Inf) = +Inf
    logb(0) = -Inf
    logb(NaN) = NaN
logb(x interface{}) interface{}, error

Returns a random string. Valid types are be word, words, sentence, para, paragraph, host, email, url.
lorem(loremType interface{}, params ...int) string, error

Converts the entire string to lowercase.
lower(str string) string

lshift(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the boolean truth of arg1 < arg2
lt(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

Returns the larger of x or y.
Special cases are:
    max(x, +Inf) = max(+Inf, x) = +Inf
    max(x, NaN) = max(NaN, x) = NaN
    max(+0, ±0) = max(±0, +0) = +0
    max(-0, -0) = -0
max(x ...interface{}) interface{}

merge(destination map[string]interface{}, sources ...map[string]interface{}) map[string]interface{}, error

Return a single list containing all elements from the lists supplied.
mergeList(lists ...[]interface{}) []interface{}

Returns the smaller of x or y.
Special cases are:
    min(x, -Inf) = min(-Inf, x) = -Inf
    min(x, NaN) = min(NaN, x) = NaN
    min(-0, ±0) = min(±0, -0) = -0
min(x ...interface{}) interface{}

Returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with that of x.
Special cases are:
    mod(±Inf, y) = NaN
    mod(NaN, y) = NaN
    mod(x, 0) = NaN
    mod(x, ±Inf) = x
    mod(x, NaN) = NaN
mod(x interface{}, y interface{}) interface{}, error

Returns integer and fractional floating-point numbers that sum to f. Both values have the same sign as f.
Special cases are:
    modf(±Inf) = ±Inf, NaN
    modf(NaN) = NaN, NaN
modf(f interface{}) interface{}, error

mul(arg1 interface{}, args ...interface{}) interface{}, error

Returns the boolean truth of arg1 != arg2
ne(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)

Returns the next representable float64 value after x towards y.
Special cases are:
    Nextafter(x, x)   = x
Nextafter(NaN, y) = NaN
Nextafter(x, NaN) = NaN
nextAfter(arg1 interface{}, arg2 interface{}) interface{}, error

Same as the indent function, but prepends a new line to the beginning of the string.
nindent(spaces int, str string) string

Removes all whitespace from a string.
nospace(str string) string

Returns the boolean negation of its single argument.
not(not(arg reflect.Value) bool

logs a message using NOTICE as log level (3).
notice(args ...interface{}) string

logs a message with format using NOTICE as log level (3).
noticef(format string, args ...interface{}) string

The current date/time. Use this in conjunction with other date functions.
now() time.Time

omit(dict map[string]interface{}, keys ...interface{}) map[string]interface{}

Returns the boolean OR of its arguments by returning the first non-empty argument or the last argument, that is, "or x y" behaves as "if x then x else y". All the arguments are evaluated.
or(or(arg0 reflect.Value, args ...reflect.Value) reflect.Value

Equivalents to critical followed by a call to panic.
panic(args ...interface{}) string

Equivalents to criticalf followed by a call to panic.
panicf(format string, args ...interface{}) string

pick(dict map[string]interface{}, keys ...interface{}) map[string]interface{}

pickv(dict map[string]interface{}, message string, keys ...interface{}) map[string]interface{}, error

pluck(arg1 string, args ...map[string]interface{}) []interface{}

Pluralizes a string.
plural(one string, many string, count int) string

Returns x**y, the base-x exponential of y.
Special cases are (in order):
    pow(x, ±0) = 1 for any x
    pow(1, y) = 1 for any y
    pow(x, 1) = x for any x
    pow(NaN, y) = NaN
    pow(x, NaN) = NaN
    pow(±0, y) = ±Inf for y an odd integer < 0
    pow(±0, -Inf) = +Inf
    pow(±0, +Inf) = +0
    pow(±0, y) = +Inf for finite y < 0 and not an odd integer
    pow(±0, y) = ±0 for y an odd integer > 0
    pow(±0, y) = +0 for finite y > 0 and not an odd integer
    pow(-1, ±Inf) = 1
    pow(x, +Inf) = +Inf for |x| > 1
    pow(x, -Inf) = +0 for |x| > 1
    pow(x, +Inf) = +0 for |x| < 1
    pow(x, -Inf) = +Inf for |x| < 1
    pow(+Inf, y) = +Inf for y > 0
    pow(+Inf, y) = +0 for y < 0
    pow(-Inf, y) = Pow(-0, -y)
    pow(x, y) = NaN for finite x < 0 and finite non-integer y
pow(x interface{}, y interface{}) interface{}, error

Returns 10**n, the base-10 exponential of n.
Special cases are:
    pow10(n) =0 for n < -323
    pow10(n) = +Inf for n > 308
pow10(n interface{}) interface{}, error

prepend(arg1 interface{}, arg2 interface{}) []interface{}

An alias for fmt.Sprint
print(args ...interface{}) string

An alias for fmt.Sprintf
printf(format string, args ...interface{}) string

An alias for fmt.Sprintln
println(args ...interface{}) string

Returns the current working directory.
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

Returns the IEEE 754 floating-point remainder of x/y.
Special cases are:
    rem(±Inf, y) = NaN
    rem(NaN, y) = NaN
    rem(x, 0) = NaN
    rem(x, ±Inf) = x
    rem(x, NaN) = NaN
rem(arg1 interface{}, arg2 interface{}) interface{}, error

Returns an array with the item repeated n times.
repeat(n int, element interface{}) []interface{}, error

Performs simple string replacement.
replace(old string, new string, src string) string

rest(arg1 interface{}) []interface{}

reverse(arg1 interface{}) []interface{}

round(arg1 interface{}, arg2 int, args ...float64) float64

rshift(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the result of the shell command as string.
run(command interface{}, args ...interface{}) interface{}, error

Intents the the elements using the provided spacer.

You can also use autoIndent as Razor expression if you don't want to specify the spacer.
Spacer will then be auto determined by the spaces that precede the expression.
Valid aliases for autoIndent are: aIndent, aindent.
sIndent(spacer string, args ...interface{}) string

Returns the element at index position or default if index is outside bounds.
safeIndex(value interface{}, index int, default interface{}) interface{}, error

Parses a string into a Semantic Version.
semver(version string) *semver.Version, error

A more robust comparison function is provided as semverCompare. This version supports version ranges.
semverCompare(constraints string, version string) bool, error

Adds the value to the supplied map using key as identifier.
set(dict interface{}, key interface{}, value interface{}) string, error

sha256sum(input string) string

Shuffle a string.
shuffle(str string) string

signBit(arg1 interface{}, arg2 interface{}) interface{}, error

Returns the sine of the radian argument x.
Special cases are:
    sin(±0) = ±0
    sin(±Inf) = NaN
    sin(NaN) = NaN
sin(x interface{}) interface{}, error

Returns Sin(x), Cos(x).
Special cases are:
    sincos(±0) = ±0, 1
    sincos(±Inf) = NaN, NaN
    sincos(NaN) = NaN, NaN
sincos(x interface{}) interface{}, error

Returns the hyperbolic sine of x.
Special cases are:
    sinh(±0) = ±0
    sinh(±Inf) = ±Inf
    sinh(NaN) = NaN
sinh(x interface{}) interface{}, error

Returns a slice of the supplied object (equivalent to object[from:to]).
slice(value interface{}, args ...interface{}) interface{}, error

Converts string from camelCase to snake_case.
snakecase(str string) string

sortAlpha(arg1 interface{}) []string

split(arg1 string, arg2 string) map[string]string

Returns a list of strings from the supplied object with newline as the separator.
splitLines(content interface{}) []interface{}

splitList(arg1 string, arg2 string) []string

Returns the square root of x.
Special cases are:
    sqrt(+Inf) = +Inf
    sqrt(±0) = ±0
    sqrt(x < 0) = NaN
    sqrt(NaN) = NaN
sqrt(x interface{}) interface{}, error

Wraps each argument with single quotes.
squote(args ...interface{}) string

Converts the supplied value into its string representation.
string(value interface{}) string

sub(arg1 interface{}, arg2 interface{}) interface{}, error

Applies the supplied regex substitute specified on the command line on the supplied string (see --substitute).
substitute(content string) string

Get a substring from a string.
substr(start int, length int, str string) string

Swaps the uppercase to lowercase and lowercase to uppercase.
swapcase(str string) string

Returns the tangent of the radian argument x.
Special cases are:
    tan(±0) = ±0
    tan(±Inf) = NaN
    tan(NaN) = NaN
tan(x interface{}) interface{}, error

Returns the hyperbolic tangent of x.
Special cases are:
    tanh(±0) = ±0
    tanh(±Inf) = ±1
    tanh(NaN) = NaN
tanh(x interface{}) interface{}, error

Returns the list of available templates names.
templateNames() []string

Returns the list of available templates.
templates() []*template.Template

Converts to title case.
title(str string) string

to(args ...interface{}) interface{}, error

Converts the supplied value to bash compatible representation.
toBash(value interface{}) string

Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.
toDate(fmt string, str string) time.Time

Converts the supplied value to compact HCL representation.
toHcl(value interface{}) string, error

Converts the supplied value to compact HCL representation used inside outer HCL definition.
toInternalHcl(value interface{}) string, error

Converts the supplied value to compact JSON representation.
toJson(value interface{}) string, error

Converts the supplied value to pretty HCL representation.
toPrettyHcl(value interface{}) string, error

Converts the supplied value to pretty JSON representation.
toPrettyJson(value interface{}) string, error

Converts the supplied value to pretty HCL representation (without multiple map declarations).
toPrettyTFVars(value interface{}) string, error

Converts the supplied value to compact quoted HCL representation.
toQuotedHcl(value interface{}) string, error

Converts the supplied value to compact quoted JSON representation.
toQuotedJson(value interface{}) string, error

Converts the supplied value to compact HCL representation (without multiple map declarations).
toQuotedTFVars(value interface{}) string, error

Converts any value to string.
toString(value interface{}) string

toStrings(arg1 interface{}) []string

Converts the supplied value to compact HCL representation (without multiple map declarations).
toTFVars(value interface{}) string, error

Converts the supplied value to YAML representation.
toYaml(value interface{}) string, error

Removes space from either side of a string.
trim(str string) string

Removes given characters from the front or back of a string.
trimAll(chars string, str string) string

Trims just the prefix from a string if present.
trimPrefix(prefix string, str string) string

Trims just the suffix from a string if present.
trimSuffix(suffix string, str string) string

Returns the integer value of x.
Special cases are:
    trunc(±0) = ±0
    trunc(±Inf) = ±Inf
    trunc(NaN) = NaN
trunc(x interface{}) interface{}, error

typeIs(arg1 string, arg2 interface{}) bool

typeIsLike(arg1 string, arg2 interface{}) bool

typeOf(arg1 interface{}) string

Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, "" are not considered as unset).
undef(default interface{}, values ...interface{}) interface{}

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

logs a message using WARNING as log level (2).
warning(args ...interface{}) string

logs a message with format using WARNING as log level (2).
warningf(format string, args ...interface{}) string

without(arg1 interface{}, args ...interface{}) []interface{}

Wraps the rendered arguments within width.
wrap(width interface{}, args ...interface{}) string

Works as wrap, but lets you specify the string to wrap with (wrap uses \n).
wrapWith(length int, spe string, str string) string

Returns the order-zero Bessel function of the second kind.
Special cases are:
    y0(+Inf) = 0
    y0(0) = -Inf
    y0(x < 0) = NaN
    y0(NaN) = NaN
y0(x interface{}) interface{}, error

Returns the order-one Bessel function of the second kind.
Special cases are:
    y1(+Inf) = 0
    y1(0) = -Inf
    y1(x < 0) = NaN
    y1(NaN) = NaN
y1(x interface{}) interface{}, error

Converts the supplied yaml string into data structure (Go spec). If context is omitted, default context is used.
yaml(yaml interface{}, context ...interface{}) interface{}, error

Returns the order-n Bessel function of the second kind.
Special cases are:
    yn(n, +Inf) = 0
    yn(n ≥ 0, 0) = -Inf
    yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even
    yn(n, x < 0) = NaN
    yn(n, NaN) = NaN
yn(n interface{}, x interface{}) interface{}, error
```
