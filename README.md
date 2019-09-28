# gotemplate

[![Build Status](https://github.com/coveooss/gotemplate/workflows/Build/badge.svg)](https://github.com/coveooss/gotemplate/actions)
[![codecov](https://codecov.io/gh/coveooss/gotemplate/branch/master/graph/badge.svg)](https://codecov.io/gh/coveooss/gotemplate)
[![Go Report Card](https://goreportcard.com/badge/github.com/coveooss/gotemplate)](https://goreportcard.com/report/github.com/coveooss/gotemplate)

## Description

Apply go template over files ending with .template in the current directory.

For more information on Go Template functionality, check this [link](https://golang.org/pkg/text/template).

This little utility just scan the current folder for `*.template` files and apply the [go template](https://golang.org/pkg/text/template) over them.

Every matching `*.ext.template` file will render a file named `*.generated.ext`. Other matched file (if --pattern is supplied) will replace the file with the rendered content and rename the original file `*.ext.original`.

[Documentation](https://coveooss.github.io/gotemplate/)
