# gotemplate
Apply go template over files ending with .template in the current directory

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

The template functions available are :
* All those supplied by [Sprig](http://masterminds.github.io/sprig)
* Conversion functions:
    * yaml(templateName)
    * json(templateName)
    * hcl(templateName)