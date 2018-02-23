package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/template"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	logging "github.com/op/go-logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version is initialized at build time through -ldflags "-X main.Version=<version number>"
var version = "master"
var tempFolder = errors.Must(ioutil.TempDir("", "gotemplate-")).(string)

const disableStdinCheck = "GOTEMPLATE_NO_STDIN"
const description = `
An extended template processor for go.

See: https://github.com/coveo/gotemplate/blob/master/README.md for complete documentation.
`

func main() {
	defer cleanup()

	var (
		app          = kingpin.New(os.Args[0], description)
		includeMath  = true
		includeSprig = true
	)

	var (
		delimiters       = app.Flag("delimiters", "Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)").PlaceHolder("{{,}},@").String()
		varFiles         = app.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').ExistingFiles()
		namedVars        = app.Flag("var", "Import named variables (if value is a file, the content is loaded)").PlaceHolder("values").Short('V').Strings()
		includePatterns  = app.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
		excludedPatterns = app.Flag("exclude", "Exclude file patterns (comma separated) when applying gotemplate recursively").PlaceHolder("pattern").Short('e').Strings()
		overwrite        = app.Flag("overwrite", "Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)").Short('o').Bool()
		substitutes      = app.Flag("substitute", "Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)").PlaceHolder("exp").Short('s').Strings()
		recursive        = app.Flag("recursive", "Process all template files recursively").Short('r').Bool()
		recursionDepth   = app.Flag("recursion-depth", "Process template files recursively specifying depth").Short('R').PlaceHolder("depth").Int()
		sourceFolder     = app.Flag("source", "Specify a source folder (default to the current folder)").PlaceHolder("folder").ExistingDir()
		targetFolder     = app.Flag("target", "Specify a target folder (default to source folder)").PlaceHolder("folder").String()
		forceStdin       = app.Flag("stdin", "Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)").Short('I').Bool()
		followSymLinks   = app.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
		print            = app.Flag("print", "Output the result directly to stdout").Short('P').Bool()
		listFunctions    = app.Flag("list-functions", "List the available functions").Short('l').Bool()
		listTemplates    = app.Flag("list-templates", "List the available templates function (--lt)").Bool()
		listALlTemplates = app.Flag("all-templates", "List all templates (--at)").Bool()
		quiet            = app.Flag("quiet", "Don not print out the name of the generated files").Short('q').Bool()
		getVersion       = app.Flag("version", "Get the current version of gotemplate").Short('v').Bool()
		disableRender    = app.Flag("disable", "Disable go template rendering (used to view razor conversion)").Short('d').Bool()
		forceColor       = app.Flag("color", "Force rendering of colors event if output is redirected").Short('c').Bool()
		logLevel         = app.Flag("log-level", "Set the logging level (0-5)").Short('L').Default("2").Int8()
		logSimple        = app.Flag("log-simple", "Disable the extended logging, i.e. no color, no date (--ls)").Bool()
		templates        = app.Arg("templates", "Template files or command to process").Strings()
	)
	app.Flag("razor", "Allow razor like expressions (@variable), ON by default --no-razor to disable").BoolVar(&template.RazorMode)
	app.Flag("math", "Include Math library, ON by default --no-math to disable").BoolVar(&includeMath)
	app.Flag("sprig", "Include Sprig library, ON by default --no-sprig to disable").BoolVar(&includeSprig)
	app.Flag("extension", "Activate gotemplate extension, ON by default, --no-extension or --no-ext to disable").BoolVar(&template.Extension)

	app.Flag("lt", "short version of --list-template").Hidden().BoolVar(listTemplates)
	app.Flag("at", "short version of --all-templates").Hidden().BoolVar(listALlTemplates)
	app.Flag("ll", "short version of --log-level").Hidden().Int8Var(logLevel)
	app.Flag("ls", "short version of --log-simple").Hidden().BoolVar(logSimple)
	app.Flag("del", "").Hidden().StringVar(delimiters)
	app.Flag("ext", "short version of --extension").Hidden().BoolVar(&template.Extension)

	if *recursionDepth != 0 {
		template.ExtensionDepth = *recursionDepth
	}
	if *recursive && *recursionDepth == 0 {
		// If recursive is specified, but there is not recursion depth, we set it to a huge depth
		*recursionDepth = 1 << 16
	}

	for i := range os.Args {
		// There is a problem with kingpin, it tries to interpret arguments beginning with @ as file
		if strings.HasPrefix(os.Args[i], "@") {
			os.Args[i] = "#" + os.Args[i]
		}
	}

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

	template.InitLogging(logging.Level(*logLevel), *logSimple)

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

	template.Libraries = 0
	if includeMath {
		template.Libraries |= template.Math
	}
	if includeSprig {
		template.Libraries |= template.Sprig
	}

	t := template.NewTemplate(createContext(*varFiles, *namedVars), *delimiters, *substitutes...)
	t.Quiet = *quiet
	t.Disabled = *disableRender
	t.Overwrite = *overwrite
	t.OutputStdout = *print
	t.TempFolder = tempFolder

	if *listFunctions || *listTemplates || *listALlTemplates {
		t = t.GetNewContext("", false)
		if *listFunctions {
			t.PrintFunctions(*templates...)
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
	if !*forceStdin && len(*templates) == 0 {
		// We only process template files if go template has not been called with piped input or explicit files
		*includePatterns = append(*includePatterns, "*.gt,*.template")
	}

	*templates = append(utils.MustFindFilesMaxDepth(*sourceFolder, *recursionDepth, *followSymLinks, extend(*includePatterns)...), *templates...)
	for i, template := range *templates {
		if !t.IsCode(template) {
			if rel, err := filepath.Rel(*sourceFolder, (*templates)[i]); err == nil {
				(*templates)[i] = rel
			}
		}
	}
	*templates = exclude(*templates, *excludedPatterns)

	resultFiles, err := t.ProcessTemplates(*sourceFolder, *targetFolder, *templates...)
	if err != nil {
		errors.Print(err)
	}

	if *forceStdin && stdinContent == "" {
		// If there is input in stdin and it has not already been consumed as data (---var)
		content := readStdin()
		if result, err := t.ProcessContent(content, "Piped input"); err == nil {
			fmt.Println(result)
		} else {
			errors.Print(err)
		}
	}

	// Apply terraform fmt if some generated files are terraform files
	if !*print {
		utils.TerraformFormat(resultFiles...)
	}
}
