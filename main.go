package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New(os.Args[0], "A template processor for go.")
	varFiles = app.Flag("variables", "Variables file to import").Short('v').Strings()
	listFunc = app.Flag("list-functions", "List the available functions").Short('l').Bool()
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()
	processFiles(context(varFiles))
}

func context(varFiles *[]string) (context map[string]interface{}) {
	if varFiles == nil {
		return nil
	}

	context = map[string]interface{}{}

	for _, varFile := range *varFiles {
		var loader func(string) (map[string]interface{}, error)
		switch strings.ToLower(filepath.Ext(varFile)) {
		case ".hcl", ".tfvars":
			loader = loadHCL
		default:
			loader = loadYaml
		}

		var (
			content = map[string]interface{}{}
			err     error
		)
		if content, err = loader(varFile); err != nil {
			Must(fmt.Errorf("%s: %v", varFile, err))
		}
		for key, value := range content {
			context[key] = value
		}
	}
	return
}
