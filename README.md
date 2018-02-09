# gotemplate

## Description

Apply go template over files ending with .template in the current directory.

For more information on Go Template functionality, check this [link](https://golang.org/pkg/text/template).

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.ext.template` file will render a file named `*.generated.ext`. Other matched file (if --pattern is supplied) will replace the file with the rendered content and rename the original file `*.ext.original`.

## Usage

```text
usage: gotemplate [<flags>] [<files>...]

A template processor for go.

See: https://github.com/coveo/gotemplate/blob/master/README.md for complete documentation.

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --delimiters={{,}},@    Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)
  -i, --import=file ...       Import variables files (could be any of YAML, JSON or HCL format)
  -V, --var=name=file ...     Import named variables (if value is a file, the content is loaded)
  -p, --patterns=pattern ...  Additional patterns that should be processed by gotemplate
  -o, --overwrite             Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
  -s, --substitute=exp ...    Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)
  -r, --recursive             Process all template files recursively
      --source=folder         Specify a source folder (default to the current folder)
      --target=folder         Specify a target folder (default to source folder)
  -I, --stdin                 Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)
  -f, --follow-symlinks       Follow the symbolic links while using the recursive option
  -P, --print                 Output the result directly to stdout
  -l, --list-functions        List the available functions
  -L, --list-templates        List the available templates function
      --all-templates         List all templates (--at)
  -q, --quiet                 Don not print out the name of the generated files
  -v, --version               Get the current version of gotemplate
  -R, --razor                 Allow razor like expressions (@variable)
  -d, --disable               Disable go template rendering (used to view razor conversion)
  -c, --color                 Force rendering of colors event if output is redirected
      --log-level=2           Set the logging level (0-5) (--ll)
      --log-simple            Disable the extended logging (--ls)
      --skip-templates        Do not load the base template *.template files, (--st)

Args:
  [<files>]  Template files to process
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
bool | string | Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`).
concat | objects ... | Returns the concatenation (without separator) of the string representation of objects.
color | attributes ... objects ... | Prints a message like print, but the first parameters are used to set color attributes (see [fatih/color](https://github.com/fatih/color)).
current | | Returns the current folder (like `pwd`, but returns the folder of the currently running folder).
exec | command string, args ... | Returns the result of the shell command as structured data (as string if no other conversion is possible).
formatList | format string, list | Return a list of strings by applying the format to each element of the supplied list.
functions | | Return the list of available functions.
get | key, map | Returns the value associated with the supplied map.
glob | args ... | Returns the expanded list of supplied arguments (expand *[]? on filename).
include | template, args ... | Returns the result of the named template rendering (like template but it is possible to capture the output).
isUndef| default, value | Returns the default value if value is not set, alias `undef` (differs from Sprig `default` function as empty value such as 0, false, "" are not considered as unset).
joinLines | objects ... | Merge the supplied objects into a newline separated string.
local_alias | name, function, source, args ... | Defines an alias (go template function) using the function (`exec`, `run`, `include`, `template`). Executed in the context of the function it maps to.
lorem | type string, min, max int | Returns a random string. Valid types are be `word`, `words`, `sentence`, `para`, `paragraph`, `host`, `email`, `url`.
mergeLists | lists | Return a single list containing all elements from the lists supplied.
pwd | | Return the current working directory.
run | command string, args ... | Returns the result of the shell command as string.
set | key string, value, map | Add the value to the supplied map using key as identifier.
splitLines | object | Returns a list of strings from the supplied object with newline as the separator.
substitutes | string | Applies the supplied regex substitute specified on the command line on the supplied string (see `--substitute`).
templateNames | | Returns the list of available templates names.
templates | | Returns the list of available templates (including itself). The returned value is the actual Go Template object list.
toHcl | interface | Returns the HCL string representation of the supplied object.
toPrettyHcl | interface | Returns the indented HCL string representation of the supplied object.
toQuotedJson | interface | Returns the JSON string representation of the supplied object with escaped quote.
toYaml | interface | Returns the YAML string representation of the supplied object.

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

E                       call                    hasPrefix               minimum                 sine
Ln10                    camelcase               hasSuffix               mod                     sineCosine
Ln2                     cat                     hcl                     modf                    sinh
Log10E                  ceil                    hello                   modulo                  sinus
Log2E                   char                    hex                     mul                     sinusCosinus
MaxFloat                clean                   hexa                    multiply                slice
MaxInt                  coalesce                hexaDecimal             ne                      smallest
MaxUInt                 color                   html                    nextAfter               snakecase
MinNonZeroFloat         compact                 htmlDate                nindent                 sortAlpha
Nan                     concat                  htmlDateInZone          nospace                 split
Phi                     contains                hyperbolicCosine        not                     splitLines
Pi                      cos                     hyperbolicCosinus       now                     splitlist
Sqrt2                   cosh                    hyperbolicSine          omit                    sqrt
SqrtE                   cosine                  hyperbolicSinus         or                      squote
SqrtPhi                 cosinus                 hyperbolicTangent       pick                    sub
SqrtPi                  current                 hypot                   pluck                   substitute
abbrev                  data                    ifUndef                 plural                  substr
abbrevboth              date                    iif                     pow                     subtract
abs                     dateInZone              ilogb                   pow10                   sum
acos                    dateModify              include                 power                   swapcase
acosh                   date_in_zone            indent                  power10                 tan
add                     date_modify             index                   prepend                 tangent
add1                    dec                     inf                     print                   tanh
ago                     decimal                 infinity                printf                  templateNames
alias                   default                 initial                 println                 templates
and                     deg                     initials                prod                    title
append                  degree                  int                     product                 to
arcCosine               derivePassword          int64                   push                    toDate
arcCosinus              dict                    isAbs                   pwd                     toHcl
arcHyperbolicCosine     dir                     isInf                   quote                   toJson
arcHyperbolicCosinus    div                     isInfinity              quotient                toPrettyHcl
arcHyperbolicSine       divide                  isNaN                   rad                     toPrettyJson
arcHyperbolicSinus      empty                   join                    radian                  toQuotedHcl
arcHyperbolicTangent    env                     joinLines               randAlpha               toQuotedJson
arcSine                 eq                      js                      randAlphaNum            toString
arcSinus                exec                    json                    randAscii               toStrings
arcTangent              exp                     keys                    randNumeric             toYaml
arcTangent2             exp2                    kindIs                  regexFind               trim
asin                    expandenv               kindOf                  regexFindAll            trimAll
asinh                   expm1                   last                    regexMatch              trimPrefix
atan                    exponent                ldexp                   regexReplaceAll         trimSuffix
atan2                   exponent2               le                      regexReplaceAllLiteral  trimall
atanh                   ext                     leftShift               regexSplit              trunc
atoi                    fail                    len                     rem                     tuple
average                 first                   lgamma                  remainder               typeIs
avg                     float64                 list                    repeat                  typeIsLike
b32dec                  floor                   local_alias             replace                 typeOf
b32enc                  formatList              log                     rest                    undef
b64dec                  frexp                   log10                   reverse                 uniq
b64enc                  functions               log1p                   rightShift              unset
band                    gamma                   log2                    round                   until
base                    ge                      logb                    rshift                  untilStep
bclear                  genCA                   lorem                   run                     untitle
biggest                 genPrivateKey           lower                   semver                  upper
bitwiseAND              genSelfSignedCert       lshift                  semverCompare           urlquery
bitwiseClear            genSignedCert           lt                      set                     uuidv4
bitwiseOR               get                     max                     sha256sum               without
bitwiseXOR              glob                    maximum                 shuffle                 wrap
bool                    gt                      merge                   signBit                 wrapWith
bor                     has                     mergeList               sin                     yaml
bxor                    hasKey                  min                     sincos
```

Links to documentations of foreign fucntions are in the section [base functions](#base-functions).
