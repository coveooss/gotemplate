package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version is initialized at build time through -ldflags "-X main.Version=<version number>"
var version = "master"

var (
	app             = kingpin.New(os.Args[0], "A template processor for go.")
	varFiles        = app.Flag("import", "Import variables files (could be YAML, JSON or HCL format)").PlaceHolder("file").Short('i').Strings()
	includePatterns = app.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
	listFunc        = app.Flag("list-functions", "List the available functions").Short('l').Bool()
	recursive       = app.Flag("recursive", "Process all template files recursively").Short('r').Bool()
	followSymLinks  = app.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
	dryRun          = app.Flag("dry-run", "Do not actually overwrite files, just show the result").Short('d').Bool()
	getVersion      = app.Flag("version", "Get the current version of gotemplate").Short('v').Bool()
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			switch err := err.(type) {
			case *errors.Error:
				fmt.Fprintln(os.Stderr, err.ErrorStack())
			default:
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}()

	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *getVersion {
		fmt.Println(version)
		os.Exit(0)
	}
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
			PanicOnError(fmt.Errorf("%s: %v", varFile, err))
		}
		for key, value := range content {
			context[key] = value
		}
	}
	return
}
