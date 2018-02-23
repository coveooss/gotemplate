# gotemplate

## Description

Apply go template over files ending with .template in the current directory.

For more information on Go Template functionality, check this [link](https://golang.org/pkg/text/template).

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.ext.template` file will render a file named `*.generated.ext`. Other matched file (if --pattern is supplied) will replace the file with the rendered content and rename the original file `*.ext.original`.

## Usage

```text
usage: gotemplate [<flags>] [<files>...]

An extended template processor for go.

See: https://github.com/coveo/gotemplate/blob/master/README.md for complete documentation.

Flags:
  -h, --help                   Show context-sensitive help (also try --help-long and --help-man).
      --delimiters={{,}},@     Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)
  -i, --import=file ...        Import variables files (could be any of YAML, JSON or HCL format)
  -V, --var=values ...         Import named variables (if value is a file, the content is loaded)
  -p, --patterns=pattern ...   Additional patterns that should be processed by gotemplate
  -e, --exclude=pattern ...    Exclude file patterns (comma separated) when applying gotemplate recursively
  -o, --overwrite              Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
  -s, --substitute=exp ...     Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)
  -r, --recursive              Process all template files recursively
  -R, --recursion-depth=depth  Process template files recursively specifying depth
      --source=folder          Specify a source folder (default to the current folder)
      --target=folder          Specify a target folder (default to source folder)
  -I, --stdin                  Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)
  -f, --follow-symlinks        Follow the symbolic links while using the recursive option
  -P, --print                  Output the result directly to stdout
  -l, --list-functions         List the available functions
      --list-templates         List the available templates function (--lt)
      --all-templates          List all templates (--at)
  -q, --quiet                  Don not print out the name of the generated files
  -v, --version                Get the current version of gotemplate
  -d, --disable                Disable go template rendering (used to view razor conversion)
  -c, --color                  Force rendering of colors event if output is redirected
  -L, --log-level=2            Set the logging level (0-5)
      --log-simple             Disable the extended logging, i.e. no color, no date (--ls)
      --razor                  Allow razor like expressions (@variable), ON by default --no-razor to disable
      --math                   Include Math library, ON by default --no-math to disable
      --sprig                  Include Sprig library, ON by default --no-sprig to disable
      --extension              Activate gotemplate extension, ON by default, --no-extension or --no-ext to disable

Args:
  [<templates>]  Template files or command to process
```

## Available functions

### Base functions

* All functions included in [Go Template](https://golang.org/pkg/text/template)
* All function included in [Sprig](http://masterminds.github.io/sprig)

## Structured import functions

All following function allows you to import data in your template. The source could be either a template, a file path or
a direct data definition.

Function name | Argument(s() |Description
--- | --- | ---
data | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the result must be valid HCL, JSON or YAML)
hcl | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting HCL must be valid)
json | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting JSON must be valid)
yaml | source [context] | Convert the data referred by `source` with `context` as argument into a go interface (the resulting YAML must be valid)

## _Usage:_

```go
// JSON data definition
{{ $json_data := `"number": 10, "text": "Hello World!, "data": { "a": 1 }, "list": ["value 1, "value 2"]` | json }}

// YAML data definition
{{ $yaml_data := yaml `
number: 10
text: Hello world!
data:
    a: 1
    b: 2
list:
    - value 1
    - value 2
` }}

// HCL data definition from template
{{ define "source" }}
number = 10
text = "Hello world!"

data {
    a = 1
    b = 2
}

list = [
    "value 1",
    "value 2",
]
{{ end }}

{{ $hcl_data := hcl "source" }}

// YAML data definition for file
{{ $yaml_file := yaml "source.yml" }}

// Assign arbitrary data definition to a variable
// data can handle any HCL, JSON, YAML or string coming from the string directly or a template name or a file name
{{ $foo := `a = 1 b = 2 c = "Hello world!"` | data }}

```

## Other functions

Function name | Argument(s() |Description
--- | --- | ---
alias | name, function, source, args ... | Defines an alias (go template function) using the function (`exec`, `run`, `include`, `template`). Executed in the context of the caller.
array | arg | Ensure that the supplied argument is an array (if it is already an array/slice, there is no change, if not, the interface is replaced by []interface{} with a single value).
bool | string | Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`).
center | string, width int | Returns a string with both left and right padding to ensure that the string is centered.
concat | objects ... | Returns the concatenation (without separator) of the string representation of objects.
color | attributes ... objects ... | Prints a message like print, but the first parameters are used to set color attributes (see [fatih/color](https://github.com/fatih/color)).
current | | Returns the current folder (like `pwd`, but returns the folder of the currently running folder).
ellipsis | function string, args ... | Returns the result of the function by expanding its last argument that must be an array into values. It's like calling function(arg1, arg2, otherArgs...).
exec | command string, args ... | Returns the result of the shell command as structured data (as string if no other conversion is possible).
formatList | format string, list | Return a list of strings by applying the format to each element of the supplied list.
func | name, function, source, default map, args names | Defines a function with the current context using the function (`exec`, `run`, `include`, `template`). Executed in the context of the caller.
functions | | Return the list of available functions.
glob | args ... | Returns the expanded list of supplied arguments (expand *[]? on filename).
id | string [replacement] | Returns a valid go identifier from the supplied string (replacing any non compliant character by replacement, default `_` ).
include | template, args ... | Returns the result of the named template rendering (like template but it is possible to capture the output).
isUndef| default, value | Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, "" are not considered as unset).
joinLines | objects ... | Merge the supplied objects into a newline separated string.
lenc | string | Returns the actual number of characters in a string according to the width of unicode characters.
local_alias | name, function, source, args ... | Defines an alias (go template function) using the function (`exec`, `run`, `include`, `template`). Executed in the context of the function it maps to.
lorem | type string, min, max int | Returns a random string. Valid types are be `word`, `words`, `sentence`, `para`, `paragraph`, `host`, `email`, `url`.
mergeLists | lists | Return a single list containing all elements from the lists supplied.
pwd | | Return the current working directory.
run | command string, args ... | Returns the result of the shell command as string.
safeIndex | array, index int, default | Returns the element at index position or default if index is outside bounds.
slice | object start [end] | Returns a slice of the supplied object (equivalent to object[from:to]).
splitLines | object | Returns a list of strings from the supplied object with newline as the separator.
string | string | Returns a String class object that allows invoking standard string operations as method.
substitutes | string | Applies the supplied regex substitute specified on the command line on the supplied string (see `--substitute`).
templateNames | | Returns the list of available templates names.
templates | | Returns the list of available templates (including itself). The returned value is the actual Go Template object list.

### Data rendering functions

Function name | Argument(s() |Description
--- | --- | ---
toBash | interface | Returns the Bash 4 string representation of the supplied object.
toHcl | interface | Returns the HCL string representation of the supplied object.
toPrettyHcl | interface | Returns the indented HCL string representation of the supplied object.
toQuotedHcl | interface | Returns the HCL string representation of the supplied object with escaped quotes.
toTFVars | interface | Returns the HCL string representation of the supplied object (without hcl map special format).
toPrettyTFVars | interface | Returns the indented HCL string representation of the supplied object (without hcl map special format).
toQuotedTFVars | interface | Returns the HCL string representation of the supplied object with escaped quotes (without hcl map special format).
toQuotedJson | interface | Returns the JSON string representation of the supplied object with escaped quotes.
toYaml | interface | Returns the YAML string representation of the supplied object.

### Map functions

Function name | Argument(s() |Description
--- | --- | ---
get | key, map | Returns the value associated with the supplied map.
set | key string, value, map | Add the value to the supplied map using key as identifier.
key | single_key_map | Returns the key name of a single element map (used to retrieve name in a declaration like `value "name" { a = 1 b = 3}`)
content | single_key_map | Returns the content of a single element map (used to retrieve content in a declaration like `value "name" { a = 1 b = 3}`)

## _Some examples:_

```go
// Random text
{{ $text := lorem "paragraph" 5 10 }}

// Random email
{{ $email := lorem "email" }}

// Concatenation of various objects
{{ $string := concat "Hello" 1 " pi = " 3.14 " 3 x 5 = " (mul 3 5)}}

// Boolean conversion
{{ $bool := bool "True" }}

// Quote all elements of a list
{{ $quoted := list 1 2 3 | formatList "\"%v\"" }}

// Colored strings
{{ $colored := color "red" "I am red" }}
{{ $colored := color "red,bgwhite" "I am red with white background" }}
{{ $colored := color "red" "bgyellow" "I am %s with %s background" "red" "yellow" }}

// Define a script
{{ define "script" }}
#! /usr/bin/env python

for i in range({{ . | default 10 }}):
    print("- ".format(i))
{{ end }}

// Run the script (render the data as output by the python script)
{{ run "script" }}

// Run the script (returns a list of number from 0 to 24)
{{ $values := exec "script" 25 }}

// Get the content of a folder
{{ $values := run "ls" "-l" }}

// Get the content of the root folder (note that arguments could be either separated of embedded with the command)
{{ $values := run "ls -l /" }}

// Get the script in a variable and execute it later
{{ $script := include "script" 40 }}
{{ exec $script }}

// Add a function that could be called later (should be created in a .template file and used later as a function)
// Note that it is not possible to
{{ alias "generate_numbers" "exec" "script" }}
...
// Can be later used as a function in another template script
{{ $numbers := generate_numbers 200 }}
{{ $numbers := 1000 | generate_numbers }}
```

### A lot of template functions are available

They are coming from either gotemplate, Sprig or native Go Template.

```text
> gotemplate -l
abbrev                  cosine                  hyperbolicCosine        nospace                 snakecase
abbrevboth              cosinus                 hyperbolicCosinus       not                     sortAlpha
abs                     current                 hyperbolicSine          now                     split
acos                    data                    hyperbolicSinus         omit                    splitLines
acosh                   date                    hyperbolicTangent       or                      splitList
add                     dateInZone              hypot                   pick                    sqrt
add1                    dateModify              id                      pickv                   squote
ago                     date_in_zone            ifUndef                 pluck                   string
alias                   date_modify             iif                     plural                  sub
and                     dec                     ilogb                   pow                     substitute
append                  decimal                 include                 pow10                   substr
arcCosine               default                 indent                  power                   subtract
arcCosinus              deg                     index                   power10                 sum
arcHyperbolicCosine     degree                  initial                 prepend                 swapcase
arcHyperbolicCosinus    derivePassword          initials                print                   tan
arcHyperbolicSine       dict                    int                     printf                  tangent
arcHyperbolicSinus      dir                     int64                   println                 tanh
arcHyperbolicTangent    div                     isAbs                   prod                    templateNames
arcSine                 divide                  isInf                   product                 templates
arcSinus                ellipsis                isInfinity              push                    title
arcTangent              empty                   isNaN                   pwd                     to
arcTangent2             env                     join                    quote                   toBash
array                   eq                      joinLines               quotient                toDate
asin                    exec                    js                      rad                     toHcl
asinh                   exp                     json                    radian                  toJson
atan                    exp2                    key                     randAlpha               toPrettyHcl
atan2                   expandenv               keys                    randAlphaNum            toPrettyJson
atanh                   expm1                   kindIs                  randAscii               toPrettyTFVars
atoi                    exponent                kindOf                  randNumeric             toQuotedHcl
average                 exponent2               last                    regexFind               toQuotedJson
avg                     ext                     ldexp                   regexFindAll            toQuotedTFVars
b32dec                  extract                 le                      regexMatch              toString
b32enc                  fail                    leftShift               regexReplaceAll         toStrings
b64dec                  first                   len                     regexReplaceAllLiteral  toTFVars
b64enc                  float64                 lenc                    regexSplit              toYaml
band                    floor                   lgamma                  rem                     trim
base                    formatList              list                    remainder               trimAll
bclear                  frexp                   local_alias             repeat                  trimPrefix
biggest                 func                    log                     replace                 trimSuffix
bitwiseAND              functions               log10                   rest                    trimall
bitwiseClear            gamma                   log1p                   reverse                 trunc
bitwiseOR               ge                      log2                    rightShift              tuple
bitwiseXOR              genCA                   logb                    round                   typeIs
bool                    genPrivateKey           lorem                   rshift                  typeIsLike
bor                     genSelfSignedCert       lower                   run                     typeOf
bxor                    genSignedCert           lshift                  safeIndex               undef
call                    get                     lt                      semver                  uniq
camelcase               glob                    max                     semverCompare           unset
cat                     gt                      maximum                 set                     until
ceil                    has                     merge                   sha256sum               untilStep
center                  hasKey                  mergeList               shuffle                 untitle
char                    hasPrefix               min                     signBit                 upper
clean                   hasSuffix               minimum                 sin                     urlquery
coalesce                hcl                     mod                     sincos                  uuidv4
color                   hello                   modf                    sine                    without
compact                 hex                     modulo                  sineCosine              wrap
concat                  hexa                    mul                     sinh                    wrapWith
contains                hexaDecimal             multiply                sinus                   yaml
content                 html                    ne                      sinusCosinus
cos                     htmlDate                nextAfter               slice
cosh                    htmlDateInZone          nindent                 smallest
```

Links to documentations of foreign functions are in the section [base functions](#base-functions).
