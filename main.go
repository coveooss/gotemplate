package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/hcl"
	"github.com/coveo/gotemplate/template"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	goerrors "github.com/go-errors/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version is initialized at build time through -ldflags "-X main.Version=<version number>"
var version = "master"
var tempFolder = errors.Must(ioutil.TempDir("", "gotemplate-")).(string)

const disableStdinCheck = "GOTEMPLATE_NO_STDIN"

func cleanup() {
	os.RemoveAll(tempFolder)

	if err := recover(); err != nil {
		switch err := err.(type) {
		case errors.Managed:
			errors.Print(err)
		case *goerrors.Error:
			errors.Printf(err.ErrorStack())
		default:
			errors.Printf("%T: %[1]v\n%s", err, string(debug.Stack()))
		}
	}
}

const description = `
A template processor for go.

See: https://github.com/coveo/gotemplate/blob/master/README.md for complete documentation.
`

func main() {
	defer cleanup()

	var (
		app              = kingpin.New(os.Args[0], description)
		delimiters       = app.Flag("delimiters", "Define the default delimiters for go template (separate the left and right delimiters by a comma)").PlaceHolder("{{,}}").String()
		varFiles         = app.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').ExistingFiles()
		namedVarFiles    = app.Flag("var", "Import named variables files (base file name is used as identifier if unspecified)").PlaceHolder("name:=file").Short('V').Strings()
		includePatterns  = app.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
		overwrite        = app.Flag("overwrite", "Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)").Short('o').Bool()
		substitutes      = app.Flag("substitute", "Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)").PlaceHolder("exp").Short('s').Strings()
		recursive        = app.Flag("recursive", "Process all template files recursively").Short('r').Bool()
		sourceFolder     = app.Flag("source", "Specify a source folder (default to the current folder)").PlaceHolder("folder").ExistingDir()
		targetFolder     = app.Flag("target", "Specify a target folder (default to source folder)").PlaceHolder("folder").String()
		forceStdin       = app.Flag("stdin", "Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)").Short('I').Bool()
		followSymLinks   = app.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
		print            = app.Flag("print", "Output the result directly to stdout").Short('P').Bool()
		listFunctions    = app.Flag("list-functions", "List the available functions").Short('l').Bool()
		listTemplates    = app.Flag("list-templates", "List the available templates function").Short('L').Bool()
		listALlTemplates = app.Flag("all-templates", "List all templates (--at)").Bool()
		quiet            = app.Flag("quiet", "Don not print out the name of the generated files").Short('q').Bool()
		getVersion       = app.Flag("version", "Get the current version of gotemplate").Short('v').Bool()
		forceColor       = app.Flag("color", "Force rendering of colors event if output is redirected").Short('c').Bool()
		files            = app.Arg("files", "Template files to process").ExistingFiles()
	)

	app.Flag("at", "short version of --all-templates").Hidden().BoolVar(listALlTemplates)

	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *forceColor {
		color.NoColor = false
	}

	if *getVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *targetFolder == "" {
		// Target folder default to source folder
		*targetFolder = *sourceFolder
	}
	*sourceFolder = errors.Must(filepath.Abs(*sourceFolder)).(string)
	*targetFolder = errors.Must(filepath.Abs(*targetFolder)).(string)

	// If target folder is not equal to source folder, we run in overwrite mode by default
	*overwrite = *overwrite || *sourceFolder != *targetFolder

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) == 0 && os.Getenv(disableStdinCheck) == "" {
		*forceStdin = true
	}

	t := template.NewTemplate(createContext(*varFiles, *namedVarFiles), *delimiters, *substitutes...)
	t.Quiet = *quiet
	t.Overwrite = *overwrite
	t.OutputStdout = *print
	t.TempFolder = tempFolder

	if *listFunctions || *listTemplates || *listALlTemplates {
		if *listFunctions {
			t.PrintFunctions()
		}
		if *listTemplates {
			t.PrintTemplates(false)
		}
		if *listALlTemplates {
			t.PrintTemplates(true)
		}
		os.Exit(0)
	}

	errors.Must(os.Chdir(*sourceFolder))
	if !*forceStdin {
		// We only process template files if go template has not been called with piped input
		*files = append(utils.MustFindFiles(*sourceFolder, *recursive, *followSymLinks, "*.template"), *files...)
	}
	*files = append(utils.MustFindFiles(*sourceFolder, *recursive, *followSymLinks, *includePatterns...), *files...)
	resultFiles, err := t.ProcessFiles(*sourceFolder, *targetFolder, *files...)
	if err != nil {
		errors.Print(err)
	}

	if *forceStdin {
		content := string(errors.Must(ioutil.ReadAll(os.Stdin)).([]byte))
		if result, err := t.ProcessContent(content, "Piped input"); err == nil {
			fmt.Print(result)
		} else {
			errors.Print(err)
		}
	}

	// Apply terraform fmt if some generated files are terraform files
	if !*print {
		utils.TerraformFormat(resultFiles...)
	}
}

func createContext(varsFiles []string, namedVarsFiles []string) (context map[string]interface{}) {
	context = map[string]interface{}{}

	// We validate the named configuration
	for i := range namedVarsFiles {
		split := strings.Split(namedVarsFiles[i], ":=")
		switch len(split) {
		case 1:
			base := strings.Split(filepath.Base(split[0]), ".")[0]
			namedVarsFiles[i] = fmt.Sprintf("%s==%s", namedVarsFiles[i], base)
		case 2:
			namedVarsFiles[i] = fmt.Sprintf("%s==%s", split[1], split[0])
		default:
			errors.Raise("Invalid named vars file %s", namedVarsFiles[i])
		}
	}

	for _, varsFile := range append(varsFiles, namedVarsFiles...) {
		varsFile, varName := utils.Split2(varsFile, "==")

		loader := utils.LoadYaml
		switch strings.ToLower(filepath.Ext(varsFile)) {
		case ".hcl", ".tfvars":
			loader = hcl.Load
		}

		var content = map[string]interface{}{}
		var err error

		if content, err = loader(varsFile); err != nil {
			errors.Raise("Error %v while loading vars file %s", varsFile, err)
		}

		if varName == "" {
			for key, value := range content {
				context[key] = value
			}
		} else {
			var varContent = map[string]interface{}{}
			for key, value := range content {
				varContent[key] = value
			}
			context[varName] = varContent
		}
	}
	return
}
