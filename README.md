# gotemplate

[![Build Status](https://github.com/coveooss/gotemplate/workflows/Build/badge.svg)](https://github.com/coveooss/gotemplate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/coveooss/gotemplate)](https://goreportcard.com/report/github.com/coveooss/gotemplate)
[![Documentation](https://img.shields.io/static/v1?label=doc&message=hugo&color=blue&logo=github)](https://coveooss.github.io/gotemplate/)

## Description

Apply template over files ending with `.template` in the current directory. Every matching `*.ext.template` file will render a file named `*.generated.ext`. It is also possible to overwrite the original files.

### Functions

Supports over a hundred functions:  

- [Go template](https://golang.org/pkg/text/template)
- [Sprig](https://github.com/Masterminds/sprig)
- Advanced serialization and deserialization in JSON, YAML, XML and HCL
- Looping and flow control functions of all kinds
- Plus a whole bunch implemented in this repository

### Syntax

Supports two distinct syntaxes (usable at the same time or individually)

Here are the statements to generate the following output:  

```text
Hello
World
```

Note: The `-` character trims whitespace. Otherwise, all lines are printed out as blank lines

#### Regular gotemplate

```go
{{- $test := list "Hello" "World" }}
{{- range $word := $test }}
{{ $word }}
{{- end }}
```

#### Razor

```go
@{test} := list("Hello", "World")
@-foreach($word := $test)
@{word}
@-end foreach
```

### Using variables

Variables can be imported from various formats (YAML, JSON and HCL) and set as CLI arguments and then used in templates. Here's an example:

`vars.json`

```json
{
  "my_var": "value"
}
```

Script:

```bash
gotemplate --var my_var2=value2 --import vars.json '{{ .my_var }} {{ .my_var2 }}'
  >>> value value2
```

More examples and statement in the [documentation](https://coveooss.github.io/gotemplate/)

Library usage [documentation](https://pkg.go.dev/github.com/coveooss/gotemplate/v3/template)
