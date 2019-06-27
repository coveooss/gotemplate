# gotemplate

[![Build Status](https://travis-ci.org/coveooss/gotemplate.svg?branch=master)](https://travis-ci.org/coveooss/gotemplate)
[![Coverage Status](https://coveralls.io/repos/github/coveooss/gotemplate/badge.svg?branch=master)](https://coveralls.io/github/coveooss/gotemplate?branch=master)

## Description

Apply go template over files ending with .template in the current directory.

For more information on Go Template functionality, check this [link](https://golang.org/pkg/text/template).

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.ext.template` file will render a file named `*.generated.ext`. Other matched file (if --pattern is supplied) will replace the file with the rendered content and rename the original file `*.ext.original`.

[Documentation](https://coveo.github.io/gotemplate/)