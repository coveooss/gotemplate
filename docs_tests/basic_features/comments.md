# Comments in gotemplate

It is important to notice that gotemplate doesn't know the language you are using and any identified gotemplate code is executed no matter where it is. Comments in the
host language mean nothing to gotemplate and will be evaluated.

## Pseudo comment

If you insert gotemplate code into file that contains another kind of code such as hcl, json, yaml, xml, java, c# or any other language, your code editor or linter may complains
because it will detect invalid characters.

To solve that problem, it is possible to inject pseudo comment into you code to hide the gotemplate code to your editor. The gotemplate is still interpretated, but obfuscated to the editor.

**It is important to render code that is valid for the host language.**

* `#!` is a pseudo comment that will removes the `#!` part but render everything after.
* `//!` is also used as pseudo comment and behave exactly as `#!`.
* `/*@ @*/` is used to specify pseudo comment in a multi line context, everything inside is rendered, but `/*@` and `@*/` are removed.

| Razor expression | Go Template           | Render    | Note
| ---------------- | -----------           | ------    | ----
| `# @(2+2)`       | `# {{ add 2 2 }}`     | `# 4`     | gotemplate code after # comment
| `// @(2+2)`      | `// {{ add 2 2 }}`    | `// 4`    | gotemplate code after // comment
| `/* @(2+2) */`   | `/* {{ add 2 2 }} */` | `/* 4 */` | gotemplate code within /* */ comment
| `#! @(2+2)`      | `{{ add 2 2 }}`       | `4`       | Pseudo comment #!
| `//! @(2+2)`     | `{{ add 2 2 }}`       | `4`       | Pseudo comment //!
| `/*@ @(2+2) @*/` | `{{ add 2 2 }}`       | `4`       | Pseudo block comment with /*@ @*/

## Real gotemplate comment

If you really want to add comment to your file and wish them to not be rendered, you must use the following syntax.

* `##@` removes every character from `##@` up to the end of line.
* `///@` removes every character from `//@` up to the end of line.
* `@#` generates gotemplate real comment `{{/* comment */}}` enclosing every character after the comment up to the end of line.
* `@//` acts exactly as `@#`.
* `@/* */` is used to generates gotemplate comment in a multi-lines context.

| Razor expression               | Go Template                                     | Render             | Note
| ----------------               | -----------                                     | ------             | ----
| `@(2+2) ##@ comment @(2*3)`    | `{{ add 2 2 }}`                                 | `4`                | Nothing is rendered after ##
| `@(2+2) ///@ comment @(2*3)`   | `{{ add 2 2 }}`                                 | `4`                | Nothing is rendered after ///
| `@(2+2) @# comment @(2*3)`     | `{{ add 2 2 }} {{/* comment {{ mul 2 3 }} */}}` | `4`                | @# generates a real gotemplate comment
| `@(2+2) @// comment @(2*3)`    | `{{ add 2 2 }} {{/* comment {{ mul 2 3 }} */}}` | `4`                | @// also generates a real gotemplate comment
| `@(2+2) @/* comment @(2*3) */` | `{{ add 2 2 }} {{/* comment {{ mul 2 3 }} */}}` | `4`                | @/* */ is used to generate multi-lines gotemplate comment

Like most of the gotemplate razor syntax, you can add the minus sign to your `@` command to render space eating gotemplate code.

| Razor expression | Go Template             | Note
| ---------------- | -----------             | ----
| `@// Comment`    | `{{/* Comment */}}`     | No space eater
| `@-// Comment`   | `{{- /* Comment */}}`   | Left space eater
| `@_-// Comment`  | `{{/* Comment */ -}}`   | Right space eaters
| `@--// Comment`  | `{{- /* Comment */ -}}` | Left and right space eaters

## Examples

### Example with JSON code

```go
/*@ @{value} := 2 + 8 * 15 @*/
{
    "Str": "string",
    "Int": 123,
    "Float": 1.23,
    "PiAsString": "@Math.Pi",
    "ComputedAsString": "@{value}",

    /* You can use the special << syntax to extract the value from the string delimiter */
    "Pi": "<<@Math.Pi",
    "Computed": "<<@{value}",
}
```

will give :

```go
{
    "Str": "string",
    "Int": 123,
    "Float": 1.23,
    "PiAsString": "3.141592653589793",
    "ComputedAsString": "122",

    /* You can use the special << syntax to extract the value from the string delimiter */
    "Pi": 3.141592653589793,
    "Computed": 122,
}
```

### Example with HCL code

```go
#! @value := 2 + 8 * 15
Str              = "string"
Int              = 123
Float            = 1.23
PiAsString       = "@Math.Pi"
ComputedAsString = "@value"

// You can use the special << syntax to extract the value from the string delimiter
Pi       = "<<@Math.Pi"
Computed = "<<@value"
```

will give:

```go
Str              = "string"
Int              = 123
Float            = 1.23
PiAsString       = "3.141592653589793"
ComputedAsString = "122"

// You can use the special << syntax to extract the value from the string delimiter
Pi       = 3.141592653589793
Computed = 122
```
