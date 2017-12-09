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

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --delimiters={{,}}      Define the default delimiters for go template (separate the left and right delimiters by a comma)
  -i, --import=file ...       Import variables files (could be any of YAML, JSON or HCL format)
  -V, --var=name:=file ...    Import named variables files (base file name is used as identifier if unspecified)
  -p, --patterns=pattern ...  Additional patterns that should be processed by gotemplate
  -o, --overwrite             Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)
  -s, --substitute=exp ...    Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)
  -r, --recursive             Process all template files recursively
      --source=folder         Specify a source folder (default to the current folder)
      --target=folder         Specify a target folder (default to source folder)
  -f, --follow-symlinks       Follow the symbolic links while using the recursive option
  -P, --print                 Output the result directly to stdout
  -l, --list-functions        List the available functions
  -L, --list-templates        List the available templates function
      --all-templates         List all templates (--at)
  -q, --quiet                 Don not print out the name of the generated files
  -v, --version               Get the current version of gotemplate

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

# Other functions

Function name | Argument(s() |Description
--- | --- | ---
alias | name, function, source, args ... | Defines an alias (go template function) using the function (`exec`, `run`, `include`, `template`). Executed in the context of the caller.
bool | string | Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`).
concat | objects ... | Returns the concatenation (without separator) of the string representation of objects.
current | | Returns the current folder (like `pwd`, but returns the folder of the currently running folder).
exec | command string, args ... | Returns the result of the shell command as structured data (as string if no other conversion is possible).
formatList | format string, list | Return a list of strings by applying the format to each element of the supplied list.
functions | | Return the list of available functions.
get | key, map | Returns the value associated with the supplied map.
glob | args ... | Returns the expanded list of supplied arguments (expand *[]? on filename).
include | template, args ... | Returns the result of the named template rendering (like template but it is possible to capture the output).
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

### A lot of template functions are available:

They are coming from either gotemplate, Sprig or native Go Template.

```
> gotemplate -l

abbrev                  dir                     int                     printf                  templateNames
abbrevboth              div                     int64                   println                 templates
add                     empty                   isAbs                   push                    title
add1                    env                     join                    pwd                     toDate
ago                     eq                      joinLines               quote                   toHcl
alias                   exec                    js                      randAlpha               toJson
and                     expandenv               json                    randAlphaNum            toPrettyHcl
append                  ext                     keys                    randAscii               toPrettyJson
atoi                    fail                    kindIs                  randNumeric             toQuotedJson
b32dec                  first                   kindOf                  regexFind               toString
b32enc                  float64                 last                    regexFindAll            toStrings
b64dec                  floor                   le                      regexMatch              toYaml
b64enc                  formatList              len                     regexReplaceAll         trim
base                    functions               list                    regexReplaceAllLiteral  trimAll
biggest                 ge                      local_alias             regexSplit              trimPrefix
bool                    genCA                   lorem                   repeat                  trimSuffix
call                    genPrivateKey           lower                   replace                 trimall
call                    genSelfSignedCert       lt                      rest                    trunc
camelcase               genSignedCert           max                     reverse                 tuple
cat                     get                     merge                   round                   typeIs
ceil                    glob                    mergeList               run                     typeIsLike
clean                   gt                      min                     semver                  typeOf
coalesce                has                     mod                     semverCompare           uniq
compact                 hasKey                  mul                     set                     unset
concat                  hasPrefix               ne                      sha256sum               until
contains                hasSuffix               nindent                 shuffle                 untilStep
current                 hcl                     nospace                 snakecase               untitle
data                    hello                   not                     sortAlpha               upper
date                    html                    now                     split                   urlquery
dateInZone              htmlDate                omit                    splitLines              uuidv4
dateModify              htmlDateInZone          or                      splitList               without
date_in_zone            include                 pick                    squote                  wrap
date_modify             indent                  pluck                   sub                     wrapWith
default                 index                   plural                  substitute              yaml
derivePassword          initial                 prepend                 substr
dict                    initials                print                   swapcase
```

Links to documentations of foreign fucntions are in the section [base functions](#base-functions).
