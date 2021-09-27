package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/hcl"
	"github.com/coveooss/gotemplate/v3/json"
	"github.com/coveooss/gotemplate/v3/template"
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/coveooss/gotemplate/v3/yaml"
	"github.com/coveooss/multilogger"
	multicolor "github.com/coveooss/multilogger/color"
	"github.com/coveooss/multilogger/errors"
	"github.com/coveord/kingpin/v2"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Version is initialized at build time through -ldflags "-X main.Version=<version number>"
var version = "2.7.4"
var tempFolder = errors.Must(ioutil.TempDir("", "gotemplate-")).(string)

const (
	envDisableStdinCheck = "GOTEMPLATE_NO_STDIN"
)
const description = `
An extended template processor for go.

See: https://coveo.github.io/gotemplate for complete documentation.
`

func runGotemplate() (exitCode int) {
	defer func() {
		if rec := recover(); rec != nil {
			errPrintf(color.RedString("Recovered %v\n"), rec)
			debug.PrintStack()
			exitCode = -1
		}
		cleanup()
	}()

	app := kingpin.New(path.Base(os.Args[0]), description).AutoShortcut().DefaultEnvars().InitOnlyOnce().UsageWriter(os.Stdout)
	app.DeleteFlag("help")
	app.DeleteFlag("help-long")
	app.HelpFlag = app.Flag("help", "Show context-sensitive help (also try --help-man).").Short('h').PreAction(func(c *kingpin.ParseContext) error { return helpAction(app, c) })
	app.HelpFlag.Bool()
	kingpin.CommandLine = app

	var (
		colorIsSet           bool
		colorEnabled         = app.Flag("color", "Force rendering of colors event if output is redirected").IsSetByUser(&colorIsSet).Bool()
		getVersion           = app.Flag("version", "Get the current version of gotemplate").Short('v').NoEnvar().Bool()
		templateLogLevel     = app.Flag("template-log-level", "Set the template logging level. Accepted values: "+multilogger.AcceptedLevelsString()).Default(logrus.InfoLevel).PlaceHolder("level").String()
		internalLogLevel     = app.Flag("internal-log-level", "Set the internal logging level. Accepted values: "+multilogger.AcceptedLevelsString()).Default(logrus.WarnLevel).PlaceHolder("level").Short('L').Alias("log-level").String()
		logFilePath          = app.Flag("internal-log-file-path", "Set a file where verbose logs should be written").PlaceHolder("path").Short('F').String()
		templateLogFileLevel = app.Flag("template-log-file-level", "Set the template logging level for the verbose logs file").Default(logrus.TraceLevel).PlaceHolder("level").String()
		internalLogFileLevel = app.Flag("internal-log-file-level", "Set the internal logging level for the verbose logs file").Default(logrus.DebugLevel).PlaceHolder("level").String()

		run                 = app.Command("run", "").Default()
		delimiters          = run.Flag("delimiters", "Define the default delimiters for go template (separate the left, right and razor delimiters by a comma)").Alias("del").PlaceHolder("{{,}},@").String()
		varFiles            = run.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').Strings()
		varFilesIfExist     = run.Flag("import-if-exist", "Import variables files (do not consider missing file as an error)").PlaceHolder("file").Strings()
		namedVars           = run.Flag("var", "Import named variables (if value is a file, the content is loaded)").PlaceHolder("values").Short('V').Strings()
		typeMode            = run.Flag("type", "Force the type used for the main context (Json, Yaml, Hcl)").Short('t').Enum("Hcl", "h", "hcl", "H", "HCL", "Json", "j", "json", "J", "JSON", "Yaml", "Yml", "y", "yml", "yaml", "Y", "YML", "YAML")
		includePatterns     = run.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
		excludedPatterns    = run.Flag("exclude", "Exclude file patterns (comma separated) when applying gotemplate recursively").PlaceHolder("pattern").Short('e').Strings()
		overwrite           = run.Flag("overwrite", "Overwrite file instead of renaming them if they exist (required only if source folder is the same as the target folder)").Short('o').Bool()
		substitutes         = run.Flag("substitute", "Substitute text in the processed files by applying the regex substitute expression (format: /regex/substitution[/timing][:/other/expr], the first character acts as separator like in sed, timing indicates when to run the substitution (b(egin), e(nd)), see: Go regexp)").PlaceHolder("exp").Short('s').Strings()
		removeEmptyLines    = run.Flag("remove-empty-lines", "Remove empty lines from the result").Alias("remove-empty").Short('E').Bool()
		recursive           = run.Flag("recursive", "Process all template files recursively").Short('r').Bool()
		recursionDepth      = run.Flag("recursion-depth", "Process template files recursively specifying depth").Short('R').PlaceHolder("depth").Int()
		sourceFolder        = run.Flag("source", "Specify a source folder (default to the current folder)").PlaceHolder("folder").String()
		targetFolder        = run.Flag("target", "Specify a target folder (default to source folder)").PlaceHolder("folder").String()
		forceStdin          = run.Flag("stdin", "Force read of the standard input to get a template definition (useful only if GOTEMPLATE_NO_STDIN is set)").Short('I').Bool()
		followSymLinks      = run.Flag("follow-symlinks", "Follow the symbolic links while using the recursive option").Short('f').Bool()
		printOutput         = run.Flag("print", "Output the result directly to stdout").Short('P').Bool()
		disableRender       = run.Flag("disable", "Disable go template rendering (used to view razor conversion)").Short('d').Bool()
		acceptNoValue       = run.Flag("accept-no-value", "Do not consider rendering <no value> as an error").Alias("no-value").Envar(template.EnvAcceptNoValue).Bool()
		strictError         = run.Flag("strict-error-validation", "Consider error encountered in any file as real error").Alias("strict").Envar(template.EnvStrictErrorCheck).Short('S').Bool()
		strictAssignations  = run.Flag("strict-assignations-validation", "Enforce strict assignation validation on global variables").Default("warning").Enum("on", "off", "warning")
		ignoreMissingImport = run.Flag("ignore-missing-import", "Exit with code 0 even if import does not exist").Bool()
		ignoreMissingSource = run.Flag("ignore-missing-source", "Exit with code 0 even if source does not exist").Bool()
		ignoreMissingPaths  = run.Flag("ignore-missing-paths", "Exit with code 0 even if import or source do not exist").Bool()
		templates           = run.Arg("templates", "Template files or commands to process").Strings()

		list          = app.Command("list", "Get detailed help on gotemplate functions").NoAutoShortcut()
		listFunctions = list.Flag("functions", "Get detailed help on function").Short('f').NoEnvar().Bool()
		listTemplates = list.Flag("templates", "List the available templates").Short('t').NoEnvar().Bool()
		listLong      = list.Flag("long", "Get detailed list").Short('l').NoEnvar().Bool()
		listAll       = list.Flag("all", "List all").Short('a').NoEnvar().Bool()
		listCategory  = list.Flag("category", "Group functions by category").Short('c').NoEnvar().Bool()
		listFilters   = list.Arg("filters", "List only functions that contains one of the filter").Strings()
	)

	loadAllAddins := true
	for i := range os.Args {
		// There is a problem with kingpin, it tries to interpret arguments beginning with @ as file
		if strings.HasPrefix(os.Args[i], "@") {
			os.Args[i] = "#!!" + os.Args[i]
		}
		if os.Args[i] == "--base" {
			loadAllAddins = false
		}
	}

	// Set the options for the available options (most of them are on by default)
	optionsOff := app.Flag("base", "Turn off all addons (they could then be enabled explicitly)").NoAutoShortcut().Bool()
	options := make([]bool, template.OptionOnByDefaultCount)
	for i := range options {
		opt := template.Options(i)
		optName := strings.ToLower(fmt.Sprint(opt))
		app.Flag(optName, fmt.Sprintf("%v Addon (ON by default)", opt)).NoAutoShortcut().Default(loadAllAddins).BoolVar(&options[i])
	}
	app.GetFlag("extension").Alias("ext")

	// Actually parse the arguments
	command, err := kingpin.CommandLine.Parse(os.Args[1:])
	if err != nil {
		// There is a problem with kingpin. If there is no command, the flags associated with the default command
		// must be specified after the templates. So, we just give it another try by injecting the default command.
		if command, err = kingpin.CommandLine.Parse(append([]string{"run"}, os.Args[1:]...)); err != nil {
			kingpin.Usage()
			os.Exit(0)
		}
	}

	// We restore back the modified arguments
	for i := range *templates {
		(*templates)[i] = strings.TrimPrefix((*templates)[i], "#!!")
	}

	// Build the optionsSet
	var optionsSet template.OptionsSet
	if *optionsOff {
		optionsSet = make(template.OptionsSet)
	} else {
		optionsSet = template.DefaultOptions()
	}

	// By default, we generate JSON list and dictionary
	if mode := *typeMode; mode != "" {
		switch strings.ToUpper(mode[:1]) {
		case "Y":
			collections.SetListHelper(yaml.GenericListHelper)
			collections.SetDictionaryHelper(yaml.DictionaryHelper)
		case "H":
			collections.SetListHelper(hcl.GenericListHelper)
			collections.SetDictionaryHelper(hcl.DictionaryHelper)
		case "J":
			collections.SetListHelper(json.GenericListHelper)
			collections.SetDictionaryHelper(json.DictionaryHelper)
		}
	} else {
		collections.SetListHelper(json.GenericListHelper)
		collections.SetDictionaryHelper(json.DictionaryHelper)
	}

	optionsSet[template.RenderingDisabled] = *disableRender
	optionsSet[template.Overwrite] = *overwrite
	optionsSet[template.OutputStdout] = *printOutput
	optionsSet[template.AcceptNoValue] = *acceptNoValue
	optionsSet[template.StrictErrorCheck] = *strictError
	for i := range options {
		optionsSet[template.Options(i)] = options[i]
	}

	switch *strictAssignations {
	case "on":
		template.StrictAssignationMode = template.AssignationValidationStrict
	case "warning":
		template.StrictAssignationMode = template.AssignationValidationWarning
	default:
		template.StrictAssignationMode = template.AssignationValidationDisabled
	}

	// Set the recursion level
	if *recursionDepth != 0 {
		template.ExtensionDepth = *recursionDepth
	}
	if *recursive && *recursionDepth == 0 {
		// If recursive is specified, but there is not recursion depth, we set it to a huge depth
		*recursionDepth = 1 << 16
	}

	if colorIsSet {
		color.NoColor = !*colorEnabled
	}

	if *getVersion {
		fmt.Println(version)
		return 0
	}

	if *ignoreMissingPaths {
		*ignoreMissingImport = true
		*ignoreMissingSource = true
	}

	if err := template.TemplateLog.SetHookLevel("", *templateLogLevel); err != nil {
		errors.Printf("Unable to set logging level for templates: %v", err)
	}
	if err := template.InternalLog.SetHookLevel("", *internalLogLevel); err != nil {
		errors.Printf("Unable to set logging level for internal logs: %v", err)
	}
	if path := *logFilePath; path != "" {
		template.TemplateLog.AddFile(path, false, *templateLogFileLevel)
		template.InternalLog.AddFile(path, false, *internalLogFileLevel)
	}

	if *targetFolder == "" {
		// Target folder default to source folder
		*targetFolder = *sourceFolder
	}
	*sourceFolder = errors.Must(filepath.Abs(*sourceFolder)).(string)
	if _, err := os.Stat(*sourceFolder); os.IsNotExist(err) {
		if !*ignoreMissingSource {
			errors.Printf("Source folder: %s does not exist", *sourceFolder)
			return 1
		}
		template.InternalLog.Debugf("Source folder: %s does not exist, skipping gotemplate", *sourceFolder)
		return 0
	}

	*targetFolder = errors.Must(filepath.Abs(*targetFolder)).(string)

	// If target folder is not equal to source folder, we run in overwrite mode by default
	*overwrite = *overwrite || *sourceFolder != *targetFolder

	stat, _ := os.Stdin.Stat()
	if (stat.Mode()&os.ModeCharDevice) == 0 && os.Getenv(envDisableStdinCheck) == "" {
		*forceStdin = true
	}

	if *removeEmptyLines {
		// Options to remove empty lines
		*substitutes = append(*substitutes, `/^\s*$/d`)
	}

	context, err := createContext(*varFiles, *varFilesIfExist, *namedVars, *typeMode, *ignoreMissingImport)
	if err != nil {
		errors.Print(err)
		return 1
	}

	t, err := template.NewTemplate("", context, *delimiters, optionsSet, *substitutes...)
	if err != nil {
		errors.Print(err)
		return 3
	}
	t.TempFolder(tempFolder)

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
		return 0
	}

	errors.Must(os.Chdir(*sourceFolder))
	if !*forceStdin && len(*templates) == 0 {
		// We only process template files if go template has not been called with piped input or explicit files
		*includePatterns = append(*includePatterns, "*.gt,*.template")
	}

	*templates = append(utils.MustFindFilesMaxDepth(*sourceFolder, *recursionDepth, *followSymLinks, extend(*includePatterns)...), *templates...)
	for i, template := range *templates {
		if file, err := filepath.Rel(*sourceFolder, template); err == nil {
			if _, err = os.Stat(file); err == nil {
				(*templates)[i] = file
			}
		}
	}
	*templates, err = exclude(*templates, *excludedPatterns)
	if err != nil {
		errors.Print(err)
		exitCode = 1
	}

	if *forceStdin && stdinContent == "" {
		*templates = append(*templates, readStdin())
	}

	resultFiles, err := t.ProcessTemplates(*sourceFolder, *targetFolder, *templates...)
	if err != nil {
		errors.Print(err)
		exitCode = 1
	}

	// Apply terraform fmt if some generated files are terraform files
	if !*printOutput {
		utils.TerraformFormat(resultFiles...)
	}
	return
}

type flag interface {
	GetFlag(string) *kingpin.FlagClause
}

func enhanceHelp(object flag, flags []*kingpin.FlagModel) {
	for _, flag := range flags {
		if flag.Name == "version" {
			continue
		}
		if len(flag.Aliases) > 0 {
			flag.Help = fmt.Sprintf("%s (alias: --%s)", flag.Help, strings.Join(flag.Aliases, ", --"))
		}
		if flag.IsBoolFlag() && len(flag.Default) > 0 {
			flag.Help = fmt.Sprintf("%s (off: --%s)", flag.Help, strings.Join(flag.NegativeAliases, ", --"))
		}
		f := object.GetFlag(flag.Name)
		f.Help(flag.Help)
	}
}

func helpAction(app *kingpin.Application, c *kingpin.ParseContext) error {
	app.Parse([]string{"run"})
	enhanceHelp(app, app.Model().Flags)
	for _, cmd := range app.Model().Commands {
		cmd := app.GetCommand(cmd.Name)
		enhanceHelp(cmd, cmd.Model().Flags)
	}

	// We manipulate the help to avoid repeating the common section of run command (since it is the default command)
	var buffer bytes.Buffer
	app.UsageWriter(&buffer)
	if os.Getenv("COLUMNS") == "" {
		os.Setenv("COLUMNS", "160")
	}
	if err := app.UsageForContextWithTemplate(c, 2, kingpin.LongHelpTemplate); err != nil {
		return err
	}
	fmt.Println(regexp.MustCompile(`(?ms)^\s+--delimiters.*?-imp\)\n(.*)$`).ReplaceAllString(buffer.String(), "$1"))

	os.Exit(0) // Force exit on display help
	return nil
}

func main() {
	os.Exit(runGotemplate())
}

var errPrintf = multicolor.ErrorPrintf
