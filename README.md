# gotemplate
Apply go template over files ending with .template in the current directory

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.template` file will render the file without the `.template` extension.

## The template functions available are :

### Base functions

* All functions included in [Go Template](https://golang.org/pkg/text/template) 
* All function included in [Sprig](http://masterminds.github.io/sprig)

### Other functions

Function name | Argument(s() |Description
--- | --- | ---
yaml | templateName [context] | Convert the template referred by `templateName` with `context` as argument into a go interface (the resulting YAML must be valid)
json | templateName [context] | Convert the template referred by `templateName` with `context` as argument into a go interface (the resulting YAML must be valid)
hcl | templateName [context] | Convert the template referred by `templateName` with `context` as argument into a go interface (the resulting YAML must be valid)
bool | string | Convert the `string` into boolean value (`string` must be `True`, `true`, `TRUE`, `1` or `False`, `false`, `FALSE`, `0`)
