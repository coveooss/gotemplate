package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
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
	var exitCode int

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf(color.RedString("Recovered %v\n"), rec)
			debug.PrintStack()
			exitCode = -1
		}
		cleanup()
		os.Exit(exitCode)
	}()

	var (
		app          = kingpin.New(os.Args[0], description)
		forceColor   = app.Flag("color", "Force rendering of colors event if output is redirected").Bool()
		forceNoColor = app.Flag("no-color", "Force rendering of colors event if output is redirected").Bool()
		getVersion   = app.Flag("version", "Get the current version of gotemplate").Short('v').Bool()

		run              = app.Command("run", "")
		delimiters       = run.Flag("delimiters", "Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)").PlaceHolder("{{,}},@").String()
		varFiles         = run.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').ExistingFiles()
		namedVars        = run.Flag("var", "Import named variables (if value is a file, the content is loaded)").PlaceHolder("values").Short('V').Strings()
		includePatterns  = run.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
		excludedPatterns = run.Flag("exclude", "Exclude file patterns (comma separated) when applying gotemplate recursively").PlaceHolder("pattern").Short('e').Strings()
		overwrite        = run.Flag("overwrite", "Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)").Short('o').Bool()
		substitutes      = run.Flag("substitute", "Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution, the first character acts as separator like in sed, see: Go regexp)").PlaceHolder("exp").Short('s').Strings()
		recursive        = run.Flag("recursive", "Process all template files recursively").Short('r').Bool()
		recursionDepth   = run.Flag("recursion-depth", "Process template files recursively specifying depth").Short('R').PlaceHolder("depth").Int()
		sourceFolder     = run.Flag("source", "Specify a source folder (default to the current folder)").PlaceHolder("folder").ExistingDir()
		targetFolder     = run.Flag("target", "Specify a target folder (default to source folder)").PlaceHolder("folder").String()
		forceStdin       = run.Flag("stdin", "Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)").Short('I').Bool()
		followSymLinks   = run.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
		print            = run.Flag("print", "Output the result directly to stdout").Short('P').Bool()
		disableRender    = run.Flag("disable", "Disable go template rendering (used to view razor conversion)").Short('d').Bool()
		debugLogLevel    = run.Flag("debug-log-level", "Set the debug logging level (0-9)").Default(utils.GetEnv(template.DebugEnvVar, "2")).Int8()
		logLevel         = run.Flag("log-level", "Set the logging level (0-9)").Short('L').Default("4").Int8()
		logSimple        = run.Flag("log-simple", "Disable the extended logging, i.e. no color, no date (--ls)").Bool()
		templates        = run.Arg("templates", "Template files or commands to process").Strings()

		list          = app.Command("list", "Get detailed help on gotemplate functions")
		listFunctions = list.Flag("functions", "Get detailed help on function").Short('f').Bool()
		listTemplates = list.Flag("templates", "List the available templates").Short('t').Bool()
		listLong      = list.Flag("long", "Get detailed list").Short('l').Bool()
		listAll       = list.Flag("all", "List all").Short('a').Bool()
		listCategory  = list.Flag("category", "Group functions by category").Short('c').Bool()
		listFilters   = list.Arg("filters", "List only functions that contains one of the filter").Strings()
	)

	quiet := app.Flag("quiet", "Deprecated").Hidden().Short('q').Bool()
	app.Flag("ll", "short version of --log-level").Hidden().Int8Var(logLevel)
	app.Flag("dl", "short version of --debug-log-level").Hidden().Int8Var(debugLogLevel)
	app.Flag("ls", "short version of --log-simple").Hidden().BoolVar(logSimple)
	app.Flag("del", "").Hidden().StringVar(delimiters)

	// Set the options for the available options (most of them are on by default)
	optionsOff := app.Flag("base", "Turn off all extensions (they could then be enabled explicitly)").Bool()
	options := make([]bool, template.OptionOnByDefaultCount)
	for i := range options {
		opt := template.Options(i)
		optName := strings.ToLower(fmt.Sprint(opt))
		options[i] = true
		app.Flag(optName, fmt.Sprintf("Option %v, on by default, --no-%s to disable", opt, optName)).BoolVar(&options[i])
	}
	app.Flag("ext", "short version of --extension").Hidden().BoolVar(&options[template.Extension])

	if len(os.Args) > 1 {
		// If no command is specified, we default to run
		switch os.Args[1] {
		case "run", "list":
			break
		default:
			os.Args = append([]string{os.Args[0]}, append([]string{"run"}, os.Args[1:]...)...)
		}
	}

	var changedArgs []int
	for i := range os.Args {
		// There is a problem with kingpin, it tries to interpret arguments beginning with @ as file
		if strings.HasPrefix(os.Args[i], "@") {
			os.Args[i] = "#!!" + os.Args[i]
			changedArgs = append(changedArgs, i)
		}
		if os.Args[i] == "--base" {
			for n := range options {
				options[n] = false
			}
		}
	}

	// Actually parse the parameters (do not evaluate args before that point since parameter values are not yet set)
	kingpin.CommandLine = app
	kingpin.CommandLine.HelpFlag.Short('h')
	command := kingpin.Parse()

	// We restore back the modified arguments
	for i := range *templates {
		if strings.HasPrefix((*templates)[i], "#!!") {
			(*templates)[i] = (*templates)[i][3:]
		}
	}

	// Build the optionsSet
	var optionsSet template.OptionsSet
	if *optionsOff {
		optionsSet = make(template.OptionsSet)
	} else {
		optionsSet = template.DefaultOptions()
	}

	optionsSet[template.RenderingDisabled] = *disableRender
	optionsSet[template.Overwrite] = *overwrite
	optionsSet[template.OutputStdout] = *print
	for i := range options {
		optionsSet[template.Options(i)] = options[i]
	}

	// Set the recursion level
	if *recursionDepth != 0 {
		template.ExtensionDepth = *recursionDepth
	}
	if *recursive && *recursionDepth == 0 {
		// If recursive is specified, but there is not recursion depth, we set it to a huge depth
		*recursionDepth = 1 << 16
	}

	if *forceNoColor {
		color.NoColor = true
	} else if *forceColor {
		color.NoColor = false
	}

	if *getVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	template.ConfigureLogging(logging.Level(*logLevel), logging.Level(*debugLogLevel), *logSimple)

	if *quiet {
		template.Log.Warning("--quiet options is deprecated, use --logginglevel instead")
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

	t := template.NewTemplate("", createContext(*varFiles, *namedVars), *delimiters, optionsSet, *substitutes...)
	t.TempFolder = tempFolder

	if command == list.FullCommand() {
		if !(*listFunctions || *listTemplates) {
			// If neither list functions or templates is selected, we default to list functions
			*listFunctions = true
		}
		t = t.GetNewContext("", false)
		if *listFunctions {
			t.PrintFunctions(*listAll, *listLong, *listCategory, *listFilters...)
		}
		if *listTemplates {
			t.PrintTemplates(*listAll, *listLong)
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
		exitCode = 1
	}

	if *forceStdin && stdinContent == "" {
		// If there is input in stdin and it has not already been consumed as data (---var)
		content := readStdin()
		if result, err := t.ProcessContent(content, "Piped input"); err == nil {
			fmt.Println(result)
		} else {
			errors.Print(err)
			exitCode = 2
		}
	}

	// Apply terraform fmt if some generated files are terraform files
	if !*print {
		utils.TerraformFormat(resultFiles...)
	}
}
