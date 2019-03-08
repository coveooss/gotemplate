{% include navigation.html %}

## Functions


| 
[Base go template functions](#base-go-template-functions) |
[Data Conversion](#data-conversion) |
[Data Manipulation](#data-manipulation) |
[Logging](#logging) |
[Mathematic Bit Operations](#mathematic-bit-operations) |
[Mathematic Fundamental](#mathematic-fundamental) |
[Mathematic Stats](#mathematic-stats) |
[Mathematic Trigonometry](#mathematic-trigonometry) |
[Mathematic Utilities](#mathematic-utilities) |
[Net](#net) |
[Operating systems functions](#operating-systems-functions) |
[Other utilities](#other-utilities) |
[Runtime](#runtime) |
[Sprig Cryptographic & Security](#sprig-cryptographic-&-security,-http://masterminds.github.io/sprig/crypto.html) |
[Sprig Date](#sprig-date,-http://masterminds.github.io/sprig/date.html) |
[Sprig Default](#sprig-default,-http://masterminds.github.io/sprig/defaults.html) |
[Sprig Dictionnary](#sprig-dictionnary,-http://masterminds.github.io/sprig/dicst.html) |
[Sprig Encoding](#sprig-encoding,-http://masterminds.github.io/sprig/encoding.html) |
[Sprig File Path](#sprig-file-path,-http://masterminds.github.io/sprig/paths.html) |
[Sprig Flow Control](#sprig-flow-control,-http://masterminds.github.io/sprig/flow_control.html) |
[Sprig General](#sprig-general,-http://masterminds.github.io/sprig/) |
[Sprig List](#sprig-list,-http://masterminds.github.io/sprig/lists.html) |
[Sprig Mathematics](#sprig-mathematics,-http://masterminds.github.io/sprig/math.html) |
[Sprig OS](#sprig-os,-http://masterminds.github.io/sprig/defaults.html) |
[Sprig Reflection](#sprig-reflection,-http://masterminds.github.io/sprig/reflection.html) |
[Sprig Regex](#sprig-regex,-http://masterminds.github.io/sprig/strings.html) |
[Sprig String Slice](#sprig-string-slice,-http://masterminds.github.io/sprig/string_slice.html) |
[Sprig Strings](#sprig-strings,-http://masterminds.github.io/sprig/strings.html) |
[Sprig Type Conversion](#sprig-type-conversion,-http://masterminds.github.io/sprig/conversion.html) |
[Sprig Version comparison](#sprig-version-comparison,-http://masterminds.github.io/sprig/semver.html) |


### Base go template functions

```go
// Returns the boolean AND of its arguments by returning the first empty argument or the last argument, that is, "and x y" behaves as "if x then y else x". All the arguments are evaluated.
func and(arg0 reflect.Value, args ...reflect.Value) reflect.Value
```

```go
// Returns the result of calling the first argument, which must be a function, with the remaining arguments as parameters. Thus "call .X.Y 1 2" is, in Go notation, dot.X.Y(1, 2) where Y is a func-valued field, map entry, or the like. The first argument must be the result of an evaluation that yields a value of function type (as distinct from a predefined function such as print). The function must return either one or two result values, the second of which is of type error. If the arguments don't match the function or the returned error value is non-nil, execution stops.
func call(fn reflect.Value, args ...reflect.Value) reflect.Value, error)
```

```go
// Returns the boolean truth of arg1 == arg2
func eq(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the boolean truth of arg1 >= arg2
func ge(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the boolean truth of arg1 > arg2
func gt(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the escaped HTML equivalent of the textual representation of its arguments. This function is unavailable in html/template, with a few exceptions.
func html(args ...interface{}) string
```

```go
// Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.
func index(item reflect.Value, indices ...reflect.Value) (reflect.Value, error)
```

```go
// Returns the escaped JavaScript equivalent of the textual representation of its arguments.
func js(args ...interface{}) string
```

```go
// Returns the boolean truth of arg1 <= arg2
func le(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the integer length of its argument.
func len(item interface{}) (int, error)
```

```go
// Returns the boolean truth of arg1 < arg2
func lt(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the boolean truth of arg1 != arg2
func ne(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error)
```

```go
// Returns the boolean negation of its single argument.
func not(not(arg reflect.Value) bool
```

```go
// Returns the boolean OR of its arguments by returning the first non-empty argument or the last argument, that is, "or x y" behaves as "if x then x else y". All the arguments are evaluated.
func or(or(arg0 reflect.Value, args ...reflect.Value) reflect.Value
```

```go
// An alias for fmt.Sprint
func print(args ...interface{}) string
```

```go
// An alias for fmt.Sprintf
func printf(format string, args ...interface{}) string
```

```go
// An alias for fmt.Sprintln
func println(args ...interface{}) string
```

```go
// Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query. This function is unavailable in html/template, with a few exceptions.
func urlquery(args ...interface{}) string
```
### Data Conversion

```go
// Tries to convert the supplied data string into data structure (Go spec). It will try to convert HCL, YAML and JSON format. If context is omitted, default context is used.
// Aliases: DATA, fromData, fromDATA
func data(data interface{}, context ...interface{}) interface{}, error
```

```go
// Converts the supplied hcl string into data structure (Go spec). If context is omitted, default context is used.
// Aliases: HCL, fromHcl, fromHCL, tfvars, fromTFVars, TFVARS, fromTFVARS
func hcl(hcl interface{}, context ...interface{}) interface{}, error
```

```go
// Converts the supplied json string into data structure (Go spec). If context is omitted, default context is used.
// Aliases: JSON, fromJson, fromJSON
func json(json interface{}, context ...interface{}) interface{}, error
```

```go
// Converts the supplied value to bash compatible representation.
func toBash(value interface{}) string
```

```go
// Converts the supplied value to compact HCL representation.
// Aliases: toHCL
func toHcl(value interface{}) string, error
```

```go
// Converts the supplied value to compact HCL representation used inside outer HCL definition.
// Aliases: toInternalHCL, toIHCL, toIHcl
func toInternalHcl(value interface{}) string, error
```

```go
// Converts the supplied value to compact JSON representation.
// Aliases: toJSON
func toJson(value interface{}) string, error
```

```go
// Converts the supplied value to pretty HCL representation.
// Aliases: toPrettyHCL
func toPrettyHcl(value interface{}) string, error
```

```go
// Converts the supplied value to pretty JSON representation.
// Aliases: toPrettyJSON
func toPrettyJson(value interface{}) string, error
```

```go
// Converts the supplied value to pretty HCL representation (without multiple map declarations).
func toPrettyTFVars(value interface{}) string, error
```

```go
// Converts the supplied value to compact quoted HCL representation.
// Aliases: toQuotedHCL
func toQuotedHcl(value interface{}) string, error
```

```go
// Converts the supplied value to compact quoted JSON representation.
// Aliases: toQuotedJSON
func toQuotedJson(value interface{}) string, error
```

```go
// Converts the supplied value to compact HCL representation (without multiple map declarations).
func toQuotedTFVars(value interface{}) string, error
```

```go
// Converts the supplied value to compact HCL representation (without multiple map declarations).
func toTFVars(value interface{}) string, error
```

```go
// Converts the supplied value to YAML representation.
// Aliases: toYAML
func toYaml(value interface{}) string, error
```

```go
// Converts the supplied yaml string into data structure (Go spec). If context is omitted, default context is used.
// Aliases: YAML, fromYaml, fromYAML
func yaml(yaml interface{}, context ...interface{}) interface{}, error
```
### Data Manipulation

```go
// Returns a String class object that allows invoking standard string operations as method.
func String(value interface{}) String
```

```go
// Append new items to an existing list, creating a new list.
// Aliases: push
func append(list interface{}, elements ...interface{}) IGenericList, error
```

```go
// Ensures that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the argument is replaced by []interface{} with a single value).
func array(value interface{}) interface{}
```

```go
// Converts the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)
func bool(str string) bool, error
```

```go
// Returns the character corresponging to the supplied integer value
func char(value interface{}) interface{}, error
```

```go
// Test to see if a list has a particular elements.
// Aliases: has
func contains(list interface{}, elements ...interface{}) bool, error
```

```go
// Returns the content of a single element map (used to retrieve content in a declaration like `value "name" { a = 1 b = 3}`)
func content(keymap interface{}) interface{}, error
```

```go
// Returns a new dictionary from a list of pairs (key, value).
// Aliases: dictionary
func dict(args ...interface{}) IDictionary, error
```

```go
// Extracts values from a slice or a map, indexes could be either integers for slice or strings for maps
func extract(source interface{}, indexes ...interface{}) interface{}, error
```

```go
// Returns the value associated with the supplied map, key and map could be inverted for convenience (i.e. when using piping mode)
func get(map interface{}, key interface{}, default ...interface{}) interface{}, error
```

```go
// Returns true if the dictionary contains the specified key.
func hasKey(dictionary interface{}, key interface{}) interface{}, error
```

```go
// Returns but the last element. 
func initial(list interface{}) interface{}, error
```

```go
// Returns a list that is the intersection of the list and all arguments (removing duplicates).
func intersect(list interface{}, elements ...interface{}) IGenericList, error
```

```go
// Returns true if the supplied value is nil.
// Aliases: isNull
func isNil(arg1 interface{}) bool
```

```go
// Returns true if the supplied value is not nil.
func isSet(arg1 interface{}) bool
```

```go
// Returns true if the supplied value is false, 0, nil or empty.
// Aliases: isEmpty
func isZero(arg1 interface{}) bool
```

```go
// Returns the key name of a single element map (used to retrieve name in a declaration like `value "name" { a = 1 b = 3}`)
func key(value interface{}) interface{}, error
```

```go
// Returns a list of all of the keys in a dict (in alphabetical order).
func keys(dictionary IDictionary) IGenericList
```

```go
// Returns the number of actual character in a string.
// Aliases: nbChars
func lenc(str string) int
```

```go
// Returns a generic list from the supplied arguments.
// Aliases: tuple
func list(args ...interface{}) IGenericList
```

```go
// Merges two or more dictionaries into one, giving precedence to the dest dictionary.
func merge(destination IDictionary, sources IDictionary, args ...IDictionary) IDictionary
```

```go
// Returns a new dict with all the keys that do not match the given keys.
func omit(dict IDictionary, keys interface{}, args ...interface{}) IDictionary
```

```go
// Selects just the given keys out of a dictionary, creating a new dict.
func pick(dict IDictionary, keys ...interface{}) IDictionary
```

```go
// Same as pick, but returns an error message if there are intruders in supplied dictionary.
func pickv(dict IDictionary, message string, keys interface{}, args ...interface{}) interface{}, error
```

```go
// Extracts a list of values matching the supplied key from a list of dictionary.
func pluck(key interface{}, dictionaries ...IDictionary) IGenericList
```

```go
// Push elements onto the front of a list, creating a new list.
func prepend(list interface{}, elements ...interface{}) IGenericList, error
```

```go
// Gets the tail of the list (everything but the first item)
func rest(list interface{}) interface{}, error
```

```go
// Produces a new list with the reversed elements of the given list.
func reverse(list interface{}) IGenericList, error
```

```go
// Returns the element at index position or default if index is outside bounds.
func safeIndex(value interface{}, index int, default interface{}) interface{}, error
```

```go
// Adds the value to the supplied map using key as identifier.
func set(dict interface{}, key interface{}, value interface{}) string, error
```

```go
// Returns a slice of the supplied object (equivalent to object[from:to]).
func slice(value interface{}, args ...interface{}) interface{}, error
```

```go
// Converts the supplied value into its string representation.
func string(value interface{}) string
```

```go
// Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, "" are not considered as unset).
// Aliases: ifUndef
func undef(default interface{}, values ...interface{}) interface{}
```

```go
// Returns a list that is the union of the list and all arguments (removing duplicates).
func union(list interface{}, elements ...interface{}) IGenericList, error
```

```go
// Generates a list with all of the duplicates removed.
// Aliases: uniq
func unique(list interface{}) IGenericList, error
```

```go
// Removes an element from a dictionary.
// Aliases: delete, remove
func unset(dictionary interface{}, key interface{}) string, error
```

```go
func values(arg1 IDictionary) IGenericList
```

```go
// Filters items out of a list.
func without(list interface{}, elements ...interface{}) IGenericList, error
```
### Logging

```go
// Logs a message using CRITICAL as log level (0).
// Aliases: criticalf
func critical(args ...interface{}) string
```

```go
// Logs a message using DEBUG as log level (5).
// Aliases: debugf
func debug(args ...interface{}) string
```

```go
// Logs a message using ERROR as log level (1).
// Aliases: errorf
func error(args ...interface{}) string
```

```go
// Equivalents to critical followed by a call to os.Exit(1).
// Aliases: fatalf
func fatal(args ...interface{}) string
```

```go
// Logs a message using INFO as log level (4).
// Aliases: infof
func info(args ...interface{}) string
```

```go
// Logs a message using NOTICE as log level (3).
// Aliases: noticef
func notice(args ...interface{}) string
```

```go
// Equivalents to critical followed by a call to panic.
// Aliases: panicf
func panic(args ...interface{}) string
```

```go
// Logs a message using WARNING as log level (2).
// Aliases: warn, warnf, warningf
func warning(args ...interface{}) string
```
### Mathematic Bit Operations

```go
// Aliases: bitwiseAND
func band(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error
```

```go
// Aliases: bitwiseClear
func bclear(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error
```

```go
// Aliases: bitwiseOR
func bor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error
```

```go
// Aliases: bitwiseXOR
func bxor(arg1 interface{}, arg2 interface{}, args ...interface{}) interface{}, error
```

```go
// Aliases: leftShift
func lshift(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
// Aliases: rightShift
func rshift(arg1 interface{}, arg2 interface{}) interface{}, error
```
### Mathematic Fundamental

```go
// Aliases: sum
func add(arg1 interface{}, args ...interface{}) interface{}, error
```

```go
// Returns the cube root of x.
// Special cases are:
//     cbrt(±0) = ±0
//     cbrt(±Inf) = ±Inf
//     cbrt(NaN) = NaN
func cbrt(x interface{}) interface{}, error
```

```go
// Returns the least integer value greater than or equal to x.
// Special cases are:
//     ceil(±0) = ±0
//     ceil(±Inf) = ±Inf
//     ceil(NaN) = NaN
// Aliases: roundUp, roundup
func ceil(x interface{}) interface{}, error
```

```go
// Returns the maximum of x-y or 0.
// Special cases are:
//     dim(+Inf, +Inf) = NaN
//     dim(-Inf, -Inf) = NaN
//     dim(x, NaN) = dim(NaN, x) = NaN
func dim(x interface{}, y interface{}) interface{}, error
```

```go
// Aliases: divide, quotient
func div(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
// Returns e**x, the base-e exponential of x.
// Special cases are:
//     exp(+Inf) = +Inf
//     exp(NaN) = NaN
// Very large values overflow to 0 or +Inf. Very small values underflow to 1.
// Aliases: exponent
func exp(x interface{}) interface{}, error
```

```go
// Returns 2**x, the base-2 exponential of x.
// Special cases are the same as exp.
// Aliases: exponent2
func exp2(x interface{}) interface{}, error
```

```go
// Returns e**x - 1, the base-e exponential of x minus 1. It is more
// accurate than exp(x) - 1 when x is near zero.
// Special cases are:
//     expm1(+Inf) = +Inf
//     expm1(-Inf) = -1
//     expm1(NaN) = NaN
// Very large values overflow to -1 or +Inf
func expm1(x interface{}) interface{}, error
```

```go
// Returns the greatest integer value less than or equal to x.
// Special cases are:
//     floor(±0) = ±0
//     floor(±Inf) = ±Inf
//     floor(NaN) = NaN
// Aliases: roundDown, rounddown, int, integer
func floor(x interface{}) interface{}, error
```

```go
// Returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with that of x.
// Special cases are:
//     mod(±Inf, y) = NaN
//     mod(NaN, y) = NaN
//     mod(x, 0) = NaN
//     mod(x, ±Inf) = x
//     mod(x, NaN) = NaN
// Aliases: modulo
func mod(x interface{}, y interface{}) interface{}, error
```

```go
// Returns integer and fractional floating-point numbers that sum to f. Both values have the same sign as f.
// Special cases are:
//     modf(±Inf) = ±Inf, NaN
//     modf(NaN) = NaN, NaN
func modf(f interface{}) interface{}, error
```

```go
// Aliases: multiply, prod, product
func mul(arg1 interface{}, args ...interface{}) interface{}, error
```

```go
// Returns x**y, the base-x exponential of y.
// Special cases are (in order):
//     pow(x, ±0) = 1 for any x
//     pow(1, y) = 1 for any y
//     pow(x, 1) = x for any x
//     pow(NaN, y) = NaN
//     pow(x, NaN) = NaN
//     pow(±0, y) = ±Inf for y an odd integer < 0
//     pow(±0, -Inf) = +Inf
//     pow(±0, +Inf) = +0
//     pow(±0, y) = +Inf for finite y < 0 and not an odd integer
//     pow(±0, y) = ±0 for y an odd integer > 0
//     pow(±0, y) = +0 for finite y > 0 and not an odd integer
//     pow(-1, ±Inf) = 1
//     pow(x, +Inf) = +Inf for |x| > 1
//     pow(x, -Inf) = +0 for |x| > 1
//     pow(x, +Inf) = +0 for |x| < 1
//     pow(x, -Inf) = +Inf for |x| < 1
//     pow(+Inf, y) = +Inf for y > 0
//     pow(+Inf, y) = +0 for y < 0
//     pow(-Inf, y) = Pow(-0, -y)
//     pow(x, y) = NaN for finite x < 0 and finite non-integer y
// Aliases: power
func pow(x interface{}, y interface{}) interface{}, error
```

```go
// Returns 10**n, the base-10 exponential of n.
// Special cases are:
//     pow10(n) =0 for n < -323
//     pow10(n) = +Inf for n > 308
// Aliases: power10
func pow10(n interface{}) interface{}, error
```

```go
// Returns the IEEE 754 floating-point remainder of x/y.
// Special cases are:
//     rem(±Inf, y) = NaN
//     rem(NaN, y) = NaN
//     rem(x, 0) = NaN
//     rem(x, ±Inf) = x
//     rem(x, NaN) = NaN
// Aliases: remainder
func rem(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
// Aliases: subtract
func sub(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
// Returns the integer value of x.
// Special cases are:
//     trunc(±0) = ±0
//     trunc(±Inf) = ±Inf
//     trunc(NaN) = NaN
// Aliases: truncate
func trunc(x interface{}) interface{}, error
```
### Mathematic Stats

```go
// Aliases: average
func avg(arg1 interface{}, args ...interface{}) interface{}, error
```

```go
// Returns the larger of x or y.
// Special cases are:
//     max(x, +Inf) = max(+Inf, x) = +Inf
//     max(x, NaN) = max(NaN, x) = NaN
//     max(+0, ±0) = max(±0, +0) = +0
//     max(-0, -0) = -0
// Aliases: maximum, biggest
func max(x ...interface{}) interface{}
```

```go
// Returns the smaller of x or y.
// Special cases are:
//     min(x, -Inf) = min(-Inf, x) = -Inf
//     min(x, NaN) = min(NaN, x) = NaN
//     min(-0, ±0) = min(±0, -0) = -0
// Aliases: minimum, smallest
func min(x ...interface{}) interface{}
```
### Mathematic Trigonometry

```go
// Returns the arccosine, in radians, of x.
// Special case is:
//     acos(x) = NaN if x < -1 or x > 1
// Aliases: arcCosine, arcCosinus
func acos(x interface{}) interface{}, error
```

```go
// Returns the inverse hyperbolic cosine of x.
// Special cases are:
//     acosh(+Inf) = +Inf
//     acosh(x) = NaN if x < 1
//     acosh(NaN) = NaN
// Aliases: arcHyperbolicCosine, arcHyperbolicCosinus
func acosh(x interface{}) interface{}, error
```

```go
// Returns the arcsine, in radians, of x.
// Special cases are:
//     asin(±0) = ±0
//     asin(x) = NaN if x < -1 or x > 1
// Aliases: arcSine, arcSinus
func asin(x interface{}) interface{}, error
```

```go
// Returns the inverse hyperbolic sine of x.
// Special cases are:
//     asinh(±0) = ±0
//     asinh(±Inf) = ±Inf
//     asinh(NaN) = NaN
// Aliases: arcHyperbolicSine, arcHyperbolicSinus
func asinh(x interface{}) interface{}, error
```

```go
// Returns the arctangent, in radians, of x.
// Special cases are:
//     atan(±0) = ±0
//     atan(±Inf) = ±Pi/2
// Aliases: arcTangent
func atan(x interface{}) interface{}, error
```

```go
// Returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
// Special cases are (in order):
//     atan2(y, NaN) = NaN
//     atan2(NaN, x) = NaN
//     atan2(+0, x>=0) = +0
//     atan2(-0, x>=0) = -0
//     atan2(+0, x<=-0) = +Pi
//     atan2(-0, x<=-0) = -Pi
//     atan2(y>0, 0) = +Pi/2
//     atan2(y<0, 0) = -Pi/2
//     atan2(+Inf, +Inf) = +Pi/4
//     atan2(-Inf, +Inf) = -Pi/4
//     atan2(+Inf, -Inf) = 3Pi/4
//     atan2(-Inf, -Inf) = -3Pi/4
//     atan2(y, +Inf) = 0
//     atan2(y>0, -Inf) = +Pi
//     atan2(y<0, -Inf) = -Pi
//     atan2(+Inf, x) = +Pi/2
//     atan2(-Inf, x) = -Pi/2
// Aliases: arcTangent2
func atan2(x interface{}, y interface{}) interface{}, error
```

```go
// Returns the inverse hyperbolic tangent of x.
// Special cases are:
//     atanh(1) = +Inf
//     atanh(±0) = ±0
//     atanh(-1) = -Inf
//     atanh(x) = NaN if x < -1 or x > 1
//     atanh(NaN) = NaN
// Aliases: arcHyperbolicTangent
func atanh(x interface{}) interface{}, error
```

```go
// Returns the cosine of the radian argument x.
// Special cases are:
//     cos(±Inf) = NaN
//     cos(NaN) = NaN
// Aliases: cosine, cosinus
func cos(x interface{}) interface{}, error
```

```go
// Returns the hyperbolic cosine of x.
// Special cases are:
//     cosh(±0) = 1
//     cosh(±Inf) = +Inf
//     cosh(NaN) = NaN
// Aliases: hyperbolicCosine, hyperbolicCosinus
func cosh(x interface{}) interface{}, error
```

```go
// Aliases: degree
func deg(arg1 interface{}) interface{}, error
```

```go
// Returns the binary exponent of x as an integer.
// Special cases are:
//     ilogb(±Inf) = MaxInt32
//     ilogb(0) = MinInt32
//     ilogb(NaN) = MaxInt32
func ilogb(x interface{}) interface{}, error
```

```go
// Returns the order-zero Bessel function of the first kind.
// Special cases are:
//     j0(±Inf) = 0
//     j0(0) = 1
//     j0(NaN) = NaN
// Aliases: firstBessel0
func j0(x interface{}) interface{}, error
```

```go
// Returns the order-one Bessel function of the first kind.
// Special cases are:
//     j1(±Inf) = 0
//     j1(NaN) = NaN
// Aliases: firstBessel1
func j1(x interface{}) interface{}, error
```

```go
// Returns the order-n Bessel function of the first kind.
// Special cases are:
//     jn(n, ±Inf) = 0
//     jn(n, NaN) = NaN
// Aliases: firstBesselN
func jn(n interface{}, x interface{}) interface{}, error
```

```go
// Returns the natural logarithm of x.
// Special cases are:
//     log(+Inf) = +Inf
//     log(0) = -Inf
//     log(x < 0) = NaN
//     log(NaN) = NaN
func log(x interface{}) interface{}, error
```

```go
// Returns the decimal logarithm of x. The special cases are the same as for log.
func log10(x interface{}) interface{}, error
```

```go
// Returns the natural logarithm of 1 plus its argument x. It is more accurate than log(1 + x) when x is near zero.
// Special cases are:
//     log1p(+Inf) = +Inf
//     log1p(±0) = ±0
//     log1p(-1) = -Inf
//     log1p(x < -1) = NaN
//     log1p(NaN) = NaN
func log1p(x interface{}) interface{}, error
```

```go
// Returns the binary logarithm of x. The special cases are the same as for log.
func log2(x interface{}) interface{}, error
```

```go
// Returns the binary exponent of x.
// Special cases are:
//     logb(±Inf) = +Inf
//     logb(0) = -Inf
//     logb(NaN) = NaN
func logb(x interface{}) interface{}, error
```

```go
// Aliases: radian
func rad(arg1 interface{}) interface{}, error
```

```go
// Returns the sine of the radian argument x.
// Special cases are:
//     sin(±0) = ±0
//     sin(±Inf) = NaN
//     sin(NaN) = NaN
// Aliases: sine, sinus
func sin(x interface{}) interface{}, error
```

```go
// Returns Sin(x), Cos(x).
// Special cases are:
//     sincos(±0) = ±0, 1
//     sincos(±Inf) = NaN, NaN
//     sincos(NaN) = NaN, NaN
// Aliases: sineCosine, sinusCosinus
func sincos(x interface{}) interface{}, error
```

```go
// Returns the hyperbolic sine of x.
// Special cases are:
//     sinh(±0) = ±0
//     sinh(±Inf) = ±Inf
//     sinh(NaN) = NaN
// Aliases: hyperbolicSine, hyperbolicSinus
func sinh(x interface{}) interface{}, error
```

```go
// Returns the tangent of the radian argument x.
// Special cases are:
//     tan(±0) = ±0
//     tan(±Inf) = NaN
//     tan(NaN) = NaN
// Aliases: tangent
func tan(x interface{}) interface{}, error
```

```go
// Returns the hyperbolic tangent of x.
// Special cases are:
//     tanh(±0) = ±0
//     tanh(±Inf) = ±1
//     tanh(NaN) = NaN
// Aliases: hyperbolicTangent
func tanh(x interface{}) interface{}, error
```

```go
// Returns the order-zero Bessel function of the second kind.
// Special cases are:
//     y0(+Inf) = 0
//     y0(0) = -Inf
//     y0(x < 0) = NaN
//     y0(NaN) = NaN
// Aliases: secondBessel0
func y0(x interface{}) interface{}, error
```

```go
// Returns the order-one Bessel function of the second kind.
// Special cases are:
//     y1(+Inf) = 0
//     y1(0) = -Inf
//     y1(x < 0) = NaN
//     y1(NaN) = NaN
// Aliases: secondBessel1
func y1(x interface{}) interface{}, error
```

```go
// Returns the order-n Bessel function of the second kind.
// Special cases are:
//     yn(n, +Inf) = 0
//     yn(n ≥ 0, 0) = -Inf
//     yn(n < 0, 0) = +Inf if n is odd, -Inf if n is even
//     yn(n, x < 0) = NaN
//     yn(n, NaN) = NaN
// Aliases: secondBesselN
func yn(n interface{}, x interface{}) interface{}, error
```
### Mathematic Utilities

```go
// Returns the absolute value of x.
// Special cases are:
//     abs(±Inf) = +Inf
//     abs(NaN) = NaN
// Aliases: absolute
func abs(x interface{}) interface{}, error
```

```go
// Aliases: decimal
func dec(arg1 interface{}) interface{}, error
```

```go
// Breaks f into a normalized fraction and an integral power of two. Returns frac and exp satisfying f == frac × 2**exp, with the absolute value of frac in the interval [½, 1).
// Special cases are:
//     frexp(±0) = ±0, 0
//     frexp(±Inf) = ±Inf, 0
//     frexp(NaN) = NaN, 0
func frexp(f interface{}) interface{}, error
```

```go
// Returns the Gamma function of x.
// Special cases are:
//     gamma(+Inf) = +Inf
//     gamma(+0) = +Inf
//     gamma(-0) = -Inf
//     gamma(x) = NaN for integer x < 0
//     gamma(-Inf) = NaN
//     gamma(NaN) = NaN
func gamma(x interface{}) interface{}, error
```

```go
// Aliases: hexa, hexaDecimal
func hex(arg1 interface{}) interface{}, error
```

```go
// Returns Sqrt(p*p + q*q), taking care to avoid unnecessary overflow and underflow.
// Special cases are:
//     hypot(±Inf, q) = +Inf
//     hypot(p, ±Inf) = +Inf
//     hypot(NaN, q) = NaN
//     hypot(p, NaN) = NaN
// Aliases: hypotenuse
func hypot(p interface{}, q interface{}) interface{}, error
```

```go
// Reports whether f is an infinity, according to sign. If sign > 0, isInf reports whether f is positive infinity. If sign < 0, IsInf reports whether f is negative infinity. If sign == 0, IsInf reports whether f is either infinity
// Aliases: isInfinity
func isInf(f interface{}, arg2 interface{}) interface{}, error
```

```go
// Reports whether f is an IEEE 754 'not-a-number' value
func isNaN(f interface{}) interface{}, error
```

```go
// Ldexp is the inverse of Frexp. Returns frac × 2**exp.
// Special cases are:
//     ldexp(±0, exp) = ±0
//     ldexp(±Inf, exp) = ±Inf
//     ldexp(NaN, exp) = NaN
func ldexp(frac interface{}, exp interface{}) interface{}, error
```

```go
// Returns the natural logarithm and sign (-1 or +1) of Gamma(x).
// Special cases are:
//     lgamma(+Inf) = +Inf
//     lgamma(0) = +Inf
//     lgamma(-integer) = +Inf
//     lgamma(-Inf) = -Inf
//     lgamma(NaN) = NaN
func lgamma(x interface{}) interface{}, error
```

```go
// Returns the next representable float64 value after x towards y.
// Special cases are:
//     Nextafter(x, x)   = x
// Nextafter(NaN, y) = NaN
// Nextafter(x, NaN) = NaN
func nextAfter(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
func signBit(arg1 interface{}, arg2 interface{}) interface{}, error
```

```go
// Returns the square root of x.
// Special cases are:
//     sqrt(+Inf) = +Inf
//     sqrt(±0) = ±0
//     sqrt(x < 0) = NaN
//     sqrt(NaN) = NaN
// Aliases: squareRoot
func sqrt(x interface{}) interface{}, error
```

```go
func to(args ...interface{}) interface{}, error
```

```go
func until(args ...interface{}) interface{}, error
```
### Net

```go
// Returns http document returned by supplied URL.
// Aliases: httpDocument, curl
func httpDoc(url interface{}) interface{}, error
```

```go
// Returns http get response from supplied URL.
func httpGet(url interface{}) *http.Response, error
```
### Operating systems functions

```go
// Returns a colored string that highlight differences between supplied texts.
// Aliases: difference
func diff(text1 interface{}, text2 interface{}) interface{}
```

```go
// Determines if a file exists or not.
// Aliases: fileExists, isExist
func exists(filename interface{}) bool, error
```

```go
// Returns the expanded list of supplied arguments (expand *[]? on filename).
// Aliases: expand
func glob(args ...interface{}) IGenericList
```

```go
// Returns the current user group information (user.Group object).
// Aliases: userGroup
func group() *user.Group, error
```

```go
// Returns the home directory of the current user.
// Aliases: homeDir, homeFolder
func home() string
```

```go
// Determines if the file is a directory.
// Aliases: isDirectory, isFolder
func isDir(filename interface{}) bool, error
```

```go
// Determines if the file is executable by the current user.
func isExecutable(filename interface{}) bool, error
```

```go
// Determines if the file is a file (i.e. not a directory).
func isFile(filename interface{}) bool, error
```

```go
// Determines if the file is readable by the current user.
func isReadable(filename interface{}) bool, error
```

```go
// Determines if the file is writeable by the current user.
func isWriteable(filename interface{}) bool, error
```

```go
// Returns the last modification time of the file.
// Aliases: lastModification, lastModificationTime
func lastMod(filename interface{}) time.Time, error
```

```go
// Returns the location of the specified executable (returns empty string if not found).
// Aliases: whereIs, look, which, type
func lookPath(arg1 interface{}) string
```

```go
// Returns the file mode.
// Aliases: fileMode
func mode(filename interface{}) os.FileMode, error
```

```go
// Returns the current working directory.
// Aliases: currentDir
func pwd() string
```

```go
// Save object to file.
// Aliases: write, writeTo
func save(filename string, object interface{}) string, error
```

```go
// Returns the file size.
// Aliases: fileSize
func size(filename interface{}) int64, error
```

```go
// Returns the file Stat information (os.Stat object).
// Aliases: fileStat
func stat(arg1 string) os.FileInfo, error
```

```go
// Returns the current user information (user.User object).
// Aliases: currentUser
func user() *user.User, error
```

```go
// Returns the current user name.
func username() string
```
### Other utilities

```go
// Returns the concatenation of supplied arguments centered within width.
// Aliases: centered
func center(width interface{}, args ...interface{}) string, error
```

```go
// Colors the rendered string.
// 
// The first arguments are interpretated as color attributes until the first non color attribute. Attributes are case insensitive.
// 
// Valid attributes are:
//     Reset, Bold, Faint, Italic, Underline, BlinkSlow, BlinkRapid, ReverseVideo, Concealed, CrossedOut
// 
// Valid color are:
//     Black, Red, Green, Yellow, Blue, Magenta, Cyan, White
// 
// Color can be prefixed by:
//     Fg:   Meaning foreground (Fg is assumed if not specified)
//     FgHi: Meaning high intensity forground
//     Bg:   Meaning background"
//     BgHi: Meaning high intensity background
// Aliases: colored, enhanced
func color(args ...interface{}) string, error
```

```go
// Returns the concatenation (without separator) of the string representation of objects.
func concat(args ...interface{}) string
```

```go
// Return a list of strings by applying the format to each element of the supplied list.
// 
// You can also use autoWrap as Razor expression if you don't want to specify the format.
// The format is then automatically induced by the context around the declaration).
// Valid aliases for autoWrap are: aWrap, awrap.
// 
// Ex:
//     Hello @<autoWrap(to(10)) World!
// Aliases: autoWrap, aWrap, awrap
func formatList(format string, list ...interface{}) IGenericList
```

```go
// Returns a valid go identifier from the supplied string (replacing any non compliant character by replacement, default _ ).
// Aliases: identifier
func id(identifier string, replaceChar ...interface{}) string
```

```go
// If testValue is empty, returns falseValue, otherwise returns trueValue.
//     WARNING: All arguments are evaluated and must by valid.
// Aliases: ternary
func iif(testValue interface{}, valueTrue interface{}, valueFalse interface{}) interface{}
```

```go
// Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.
func indent(nbSpace int, args ...interface{}) string
```

```go
// Merge the supplied objects into a newline separated string.
func joinLines(format ...interface{}) string
```

```go
// Returns a random string. Valid types are be word, words, sentence, para, paragraph, host, email, url.
// Aliases: loremIpsum
func lorem(loremType interface{}, params ...int) string, error
```

```go
// Return a single list containing all elements from the lists supplied.
func mergeList(lists ...IGenericList) IGenericList
```

```go
// Aliases: nindent
func nIndent(nbSpace int, args ...interface{}) string
```

```go
// Returns an array with the item repeated n times.
func repeat(n int, element interface{}) IGenericList, error
```

```go
// Indents the elements using the provided spacer.
// 
// You can also use autoIndent as Razor expression if you don't want to specify the spacer.
// Spacer will then be auto determined by the spaces that precede the expression.
// Valid aliases for autoIndent are: aIndent, aindent.
// Aliases: sindent, spaceIndent, autoIndent, aindent, aIndent
func sIndent(spacer string, args ...interface{}) string
```

```go
// Returns a list of strings from the supplied object with newline as the separator.
func splitLines(content interface{}) []interface{}
```

```go
// Wraps the rendered arguments within width.
// Aliases: wrapped
func wrap(width interface{}, args ...interface{}) string, error
```
### Runtime

```go
// Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the caller.
func alias(name string, function string, source interface{}, args ...interface{}) string, error
```

```go
// Returns the list of all functions that are simply an alias of another function.
func aliases() []string
```

```go
// Returns the list of all available functions.
func allFunctions() []string
```

```go
// Raises a formated error if the test condition is false.
// Aliases: assertion
func assert(test interface{}, message ...interface{}) string, error
```

```go
// Issues a formated warning if the test condition is false.
// Aliases: assertw
func assertWarning(test interface{}, message ...interface{}) string
```

```go
// Returns all functions group by categories.
// 
// The returned value has the following properties:
//     Name        string
//     Functions    []string
func categories() []template.FuncCategory
```

```go
// Returns the current folder (like pwd, but returns the folder of the currently running folder).
func current() string
```

```go
// Returns the result of the function by expanding its last argument that must be an array into values. It's like calling function(arg1, arg2, otherArgs...).
func ellipsis(function string, args ...interface{}) interface{}, error
```

```go
// Returns the result of the shell command as structured data (as string if no other conversion is possible).
// Aliases: execute
func exec(command interface{}, args ...interface{}) interface{}, error
```

```go
// Exits the current program execution.
func exit(exitValue int) int
```

```go
// Defines a function with the current context using the function (exec, run, include, template). Executed in the context of the caller.
func func(name string, function string, source interface{}, config interface{}) string, error
```

```go
// Returns the information relative to a specific function.
// 
// The returned value has the following properties:
//     Name        string
//     Description string
//     Signature   string
//     Group       string
//     Aliases     []string
//     Arguments   string
//     Result      string
func function(name string) template.FuncInfo
```

```go
// Returns the list of all available functions (excluding aliases).
func functions() []string
```

```go
// List all attributes accessible from the supplied object.
// Aliases: attr, attributes
func getAttributes(arg1 interface{}) string
```

```go
// List all methods signatures accessible from the supplied object.
// Aliases: methods
func getMethods(arg1 interface{}) string
```

```go
// List all attributes and methods signatures accessible from the supplied object.
// Aliases: sign, signature
func getSignature(arg1 interface{}) string
```

```go
// Returns the result of the named template rendering (like template but it is possible to capture the output).
func include(source interface{}, context ...interface{}) interface{}, error
```

```go
// Defines an alias (go template function) using the function (exec, run, include, template). Executed in the context of the function it maps to.
func localAlias(name string, function string, source interface{}, args ...interface{}) string, error
```

```go
// Raise a formated error.
// Aliases: raiseError
func raise(args ...interface{}) string, error
```

```go
// Returns the result of the shell command as string.
func run(command interface{}, args ...interface{}) interface{}, error
```

```go
// Applies the supplied regex substitute specified on the command line on the supplied string (see --substitute).
func substitute(content string) string
```

```go
// Returns the list of available templates names.
func templateNames() []string
```

```go
// Returns the list of available templates.
func templates() []*template.Template
```
### Sprig Cryptographic & Security, http://masterminds.github.io/sprig/crypto.html

```go
// Computes Adler-32 checksum.
func adler32sum(input string) string
```

```go
func buildCustomCert(arg1 string, arg2 string) sprig.certificate, error
```

```go
func derivePassword(arg1 uint32, arg2 string, arg3 string, arg4 string, arg5 string) string
```

```go
func genCA(arg1 string, arg2 int) sprig.certificate, error
```

```go
// Generates a new private key encoded into a PEM block. Type should be: ecdsa, dsa or rsa
func genPrivateKey(type string) string
```

```go
func genSelfSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int) sprig.certificate, error
```

```go
func genSignedCert(arg1 string, arg2 []interface{}, arg3 []interface{}, arg4 int, arg5 sprig.certificate) sprig.certificate, error
```

```go
// Computes SHA1 digest.
func sha1sum(input string) string
```

```go
// Computes SHA256 digest.
func sha256sum(input string) string
```
### Sprig Date, http://masterminds.github.io/sprig/date.html

```go
// The ago function returns duration from time.Now in seconds resolution.
func ago(date interface{}) string
```

```go
// The date function formats a dat (https://golang.org/pkg/time/#Time.Format).
func date(fmt string, date interface{}) string
```

```go
// Same as date, but with a timezone.
// Aliases: date_in_zone
func dateInZone(fmt string, date interface{}, zone string) string
```

```go
// The dateModify takes a modification and a date and returns the timestamp.
// Aliases: date_modify
func dateModify(fmt string, date time.Time) time.Time
```

```go
// The htmlDate function formates a date for inserting into an HTML date picker input field.
func htmlDate(date interface{}) string
```

```go
// Same as htmlDate, but with a timezone.
func htmlDateInZone(date interface{}, zone string) string
```

```go
// The current date/time. Use this in conjunction with other date functions.
func now() time.Time
```

```go
// Converts a string to a date. The first argument is the date layout and the second the date string. If the string can’t be convert it returns the zero value.
func toDate(fmt string, str string) time.Time
```
### Sprig Default, http://masterminds.github.io/sprig/defaults.html

```go
func coalesce(args ...interface{}) interface{}
```

```go
func compact(arg1 interface{}) []interface{}
```

```go
func default(arg1 interface{}, args ...interface{}) interface{}
```

```go
func empty(arg1 interface{}) bool
```

```go
// Aliases: ternary
func ternarySprig(arg1 interface{}, arg2 interface{}, arg3 bool) interface{}
```

```go
// Aliases: toJson
func toJsonSprig(arg1 interface{}) string
```

```go
// Aliases: toPrettyJson
func toPrettyJsonSprig(arg1 interface{}) string
```
### Sprig Dictionnary, http://masterminds.github.io/sprig/dicst.html

```go
// Aliases: dict
func dictSprig(args ...interface{}) map[string]interface{}
```

```go
// Aliases: hasKey
func hasKeySprig(arg1 map[string]interface{}, arg2 string) bool
```

```go
// Aliases: keys
func keysSprig(args ...map[string]interface{}) []string
```

```go
// Aliases: list, tuple, tupleSprig
func listSprig(args ...interface{}) []interface{}
```

```go
// Merge two or more dictionaries into one, giving precedence from **right to left**, effectively overwriting values in the dest dictionary
func mergeOverwrite(arg1 map[string]interface{}, args ...map[string]interface{}) interface{}
```

```go
// Aliases: merge
func mergeSprig(arg1 map[string]interface{}, args ...map[string]interface{}) interface{}
```

```go
// Aliases: omit
func omitSprig(arg1 map[string]interface{}, args ...string) map[string]interface{}
```

```go
// Aliases: pick
func pickSprig(arg1 map[string]interface{}, args ...string) map[string]interface{}
```

```go
// Aliases: pluck
func pluckSprig(arg1 string, args ...map[string]interface{}) []interface{}
```

```go
// Aliases: set
func setSprig(arg1 map[string]interface{}, arg2 string, arg3 interface{}) map[string]interface{}
```

```go
// Aliases: unset
func unsetSprig(arg1 map[string]interface{}, arg2 string) map[string]interface{}
```

```go
// Aliases: values
func valuesSprig(arg1 map[string]interface{}) []interface{}
```
### Sprig Encoding, http://masterminds.github.io/sprig/encoding.html

```go
func b32dec(arg1 string) string
```

```go
func b32enc(arg1 string) string
```

```go
func b64dec(arg1 string) string
```

```go
func b64enc(arg1 string) string
```
### Sprig File Path, http://masterminds.github.io/sprig/paths.html

```go
func base(arg1 string) string
```

```go
func clean(arg1 string) string
```

```go
func dir(arg1 string) string
```

```go
func ext(arg1 string) string
```

```go
func isAbs(arg1 string) bool
```
### Sprig Flow Control, http://masterminds.github.io/sprig/flow_control.html

```go
// Unconditionally returns an empty string and an error with the specified text. This is useful in scenarios where other conditionals have determined that template rendering should fail.
func fail(arg1 string) string, error
```
### Sprig General, http://masterminds.github.io/sprig/

```go
// Simple hello by Sprig
func hello() string
```

```go
// Aliases: uuid, guid, GUID
func uuidv4() string
```
### Sprig List, http://masterminds.github.io/sprig/lists.html

```go
// Aliases: append, push, pushSprig
func appendSprig(arg1 interface{}, arg2 interface{}) []interface{}
```

```go
func first(arg1 interface{}) interface{}
```

```go
// Aliases: has
func hasSprig(arg1 interface{}, arg2 interface{}) bool
```

```go
// Aliases: initial
func initialSprig(arg1 interface{}) []interface{}
```

```go
func last(arg1 interface{}) interface{}
```

```go
// Aliases: prepend
func prependSprig(arg1 interface{}, arg2 interface{}) []interface{}
```

```go
// Aliases: rest
func restSprig(arg1 interface{}) []interface{}
```

```go
// Aliases: reverse
func reverseSprig(arg1 interface{}) []interface{}
```

```go
// Aliases: slice
func sliceSprig(arg1 interface{}, args ...interface{}) interface{}
```

```go
// Aliases: uniq
func uniqSprig(arg1 interface{}) []interface{}
```

```go
// Aliases: without
func withoutSprig(arg1 interface{}, args ...interface{}) []interface{}
```
### Sprig Mathematics, http://masterminds.github.io/sprig/math.html

```go
func add1(arg1 interface{}) int64
```

```go
// Aliases: add
func addSprig(args ...interface{}) int64
```

```go
// Aliases: ceil
func ceilSprig(arg1 interface{}) float64
```

```go
// Aliases: div
func divSprig(arg1 interface{}, arg2 interface{}) int64
```

```go
// Aliases: floor
func floorSprig(arg1 interface{}) float64
```

```go
// Aliases: max, biggest, biggestSprig
func maxSprig(arg1 interface{}, args ...interface{}) int64
```

```go
// Aliases: min
func minSprig(arg1 interface{}, args ...interface{}) int64
```

```go
// Aliases: mod
func modSprig(arg1 interface{}, arg2 interface{}) int64
```

```go
// Aliases: mul
func mulSprig(arg1 interface{}, args ...interface{}) int64
```

```go
func round(arg1 interface{}, arg2 int, args ...float64) float64
```

```go
// Aliases: sub
func subSprig(arg1 interface{}, arg2 interface{}) int64
```

```go
func untilStep(arg1 int, arg2 int, arg3 int) []int
```
### Sprig OS, http://masterminds.github.io/sprig/defaults.html

```go
func env(arg1 string) string
```

```go
func expandenv(arg1 string) string
```
### Sprig Reflection, http://masterminds.github.io/sprig/reflection.html

```go
// Aliases: kindis
func kindIs(arg1 string, arg2 interface{}) bool
```

```go
// Aliases: kindof
func kindOf(arg1 interface{}) string
```

```go
// Aliases: typeis
func typeIs(arg1 string, arg2 interface{}) bool
```

```go
// Aliases: typeisLike
func typeIsLike(arg1 string, arg2 interface{}) bool
```

```go
// Aliases: typeof
func typeOf(arg1 interface{}) string
```
### Sprig Regex, http://masterminds.github.io/sprig/strings.html

```go
// Returns the first (left most) match of the regular expression in the input string.
func regexFind(regex string, str string) string
```

```go
// Returns a slice of all matches of the regular expression in the input string.
func regexFindAll(regex string, str string, n int) []string
```

```go
// Returns true if the input string matches the regular expression.
func regexMatch(regex string, str string) bool
```

```go
// Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement. Inside string replacement, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.
func regexReplaceAll(regex string, str string, repl string) string
```

```go
// Returns a copy of the input string, replacing matches of the Regexp with the replacement string replacement The replacement string is substituted directly, without using Expand.
func regexReplaceAllLiteral(regex string, str string, repl string) string
```

```go
// Slices the input string into substrings separated by the expression and returns a slice of the substrings between those expression matches. The last parameter n determines the number of substrings to return, where -1 means return all matches.
func regexSplit(regex string, str string, n int) []string
```
### Sprig String Slice, http://masterminds.github.io/sprig/string_slice.html

```go
func join(arg1 string, arg2 interface{}) string
```

```go
func sortAlpha(arg1 interface{}) []string
```

```go
func split(arg1 string, arg2 string) map[string]string
```

```go
func splitList(arg1 string, arg2 string) []string
```

```go
func splitn(arg1 string, arg2 int, arg3 string) map[string]string
```

```go
func toStrings(arg1 interface{}) []string
```
### Sprig Strings, http://masterminds.github.io/sprig/strings.html

```go
// Truncates a string with ellipses (...).
func abbrev(width int, str string) string
```

```go
// Abbreviates both sides with ellipses (...).
func abbrevboth(left int, right int, str string) string
```

```go
// Converts string from snake_case to CamelCase.
func camelcase(str string) string
```

```go
// Concatenates multiple strings together into one, separating them with spaces.
func cat(args ...interface{}) string
```

```go
// Tests to see if one string is contained inside of another.
// Aliases: contains
func containsSprig(substr string, str string) bool
```

```go
// Tests whether a string has a given prefix.
func hasPrefix(prefix string, str string) bool
```

```go
// Tests whether a string has a given suffix.
func hasSuffix(suffix string, str string) bool
```

```go
// Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.
// Aliases: indent
func indentSprig(spaces int, str string) string
```

```go
// Given multiple words, takes the first letter of each word and combine.
func initials(str string) string
```

```go
// Convert string from camelCase to kebab-case.
func kebabcase(str string) string
```

```go
// Converts the entire string to lowercase.
func lower(str string) string
```

```go
// Same as the indent function, but prepends a new line to the beginning of the string.
// Aliases: nindent
func nindentSprig(spaces int, str string) string
```

```go
// Removes all whitespace from a string.
func nospace(str string) string
```

```go
// Pluralizes a string.
func plural(one string, many string, count int) string
```

```go
// Wraps each argument with double quotes.
func quote(str ...interface{}) string
```

```go
// Generates random string with letters.
func randAlpha(count int) string
```

```go
// Generates random string with letters and digits.
func randAlphaNum(count int) string
```

```go
// Generates random string with ASCII printable characters.
func randAscii(count int) string
```

```go
// Generates random string with digits.
func randNumeric(count int) string
```

```go
// Repeats a string multiple times.
// Aliases: repeat
func repeatSprig(count int, str string) string
```

```go
// Performs simple string replacement.
func replace(old string, new string, src string) string
```

```go
// Shuffle a string.
func shuffle(str string) string
```

```go
// Converts string from camelCase to snake_case.
func snakecase(str string) string
```

```go
// Wraps each argument with single quotes.
func squote(args ...interface{}) string
```

```go
// Get a substring from a string.
func substr(start int, length int, str string) string
```

```go
// Swaps the uppercase to lowercase and lowercase to uppercase.
func swapcase(str string) string
```

```go
// Converts to title case.
func title(str string) string
```

```go
// Converts any value to string.
func toString(value interface{}) string
```

```go
// Removes space from either side of a string.
func trim(str string) string
```

```go
// Removes given characters from the front or back of a string.
// Aliases: trimall
func trimAll(chars string, str string) string
```

```go
// Trims just the prefix from a string if present.
func trimPrefix(prefix string, str string) string
```

```go
// Trims just the suffix from a string if present.
func trimSuffix(suffix string, str string) string
```

```go
// Truncates a string (and add no suffix).
// Aliases: trunc
func truncSprig(length int, str string) string
```

```go
// Removes title casing.
func untitle(str string) string
```

```go
// Converts the entire string to uppercase.
func upper(str string) string
```

```go
// Wraps text at a given column count.
// Aliases: wrap
func wrapSprig(length int, str string) string
```

```go
// Works as wrap, but lets you specify the string to wrap with (wrap uses \n).
func wrapWith(length int, spe string, str string) string
```
### Sprig Type Conversion, http://masterminds.github.io/sprig/conversion.html

```go
func atoi(arg1 string) int
```

```go
func float64(arg1 interface{}) float64
```

```go
func int64(arg1 interface{}) int64
```

```go
// Aliases: int
func intSprig(arg1 interface{}) int
```
### Sprig Version comparison, http://masterminds.github.io/sprig/semver.html

```go
// Parses a string into a Semantic Version.
func semver(version string) *semver.Version, error
```

```go
// A more robust comparison function is provided as semverCompare. This version supports version ranges.
func semverCompare(constraints string, version string) bool, error
```

