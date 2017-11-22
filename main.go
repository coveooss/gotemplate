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
	delimiters      = app.Flag("delimiters", "Define the default delimiters for go template (separate the left and right delimiters by a comma)").PlaceHolder("{{,}}").String()
	varFiles        = app.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').Strings()
	includePatterns = app.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
	overwrite       = app.Flag("overwrite", "Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)").Short('o').Bool()
	substitutes     = app.Flag("substitute", "Substitute text in the processed files by applying the regex substitute expression (format: regex/substitution, see: Go regexp)").PlaceHolder("exp").Short('s').Strings()
	recursive       = app.Flag("recursive", "Process all template files recursively").Short('r').Bool()
	sourceFolder    = app.Flag("source", "Specify a source folder (default to the current folder)").PlaceHolder("folder").String()
	targetFolder    = app.Flag("target", "Specify a target folder (default to source folder)").PlaceHolder("folder").String()
	followSymLinks  = app.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
	dryRun          = app.Flag("dry-run", "Do not actually write files, just show the result").Short('d').Bool()
	listFunc        = app.Flag("list-functions", "List the available functions").Short('l').Bool()
	silent          = app.Flag("silent", "Don not print out the name of the generated files").Bool()
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

	if *targetFolder == "" {
		// Target folder default to source folder
		*targetFolder = *sourceFolder
	}
	*sourceFolder = PanicOnError(filepath.Abs(*sourceFolder)).(string)
	*targetFolder = PanicOnError(filepath.Abs(*targetFolder)).(string)

	// If target folder is not equal to source folder, we run in overwrite mode by default
	*overwrite = *overwrite || *sourceFolder != *targetFolder

	if *sourceFolder != "" {
		PanicOnError(os.Chdir(*sourceFolder))
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
