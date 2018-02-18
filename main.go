package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/template"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	goerrors "github.com/go-errors/errors"
	logging "github.com/op/go-logging"
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
		delimiters       = app.Flag("delimiters", "Define the default delimiters for go template (separate the left, right and razor delimiters by a comma) (--del)").PlaceHolder("{{,}},@").String()
		varFiles         = app.Flag("import", "Import variables files (could be any of YAML, JSON or HCL format)").PlaceHolder("file").Short('i').ExistingFiles()
		namedVars        = app.Flag("var", "Import named variables (if value is a file, the content is loaded)").PlaceHolder("values").Short('V').Strings()
		includePatterns  = app.Flag("patterns", "Additional patterns that should be processed by gotemplate").PlaceHolder("pattern").Short('p').Strings()
		excludedPatterns = app.Flag("exclude", "Exclude file patterns (comma separated) when applying gotemplate recursively").PlaceHolder("pattern").Short('e').Strings()
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
		razorSyntax      = app.Flag("razor", "Allow razor like expressions (@variable)").Short('R').Bool()
		disableRender    = app.Flag("disable", "Disable go template rendering (used to view razor conversion)").Short('d').Bool()
		forceColor       = app.Flag("color", "Force rendering of colors event if output is redirected").Short('c').Bool()
		logLevel         = app.Flag("log-level", "Set the logging level (0-5) (--ll)").Default("2").Int8()
		logSimple        = app.Flag("log-simple", "Disable the extended logging (--ls)").Bool()
		skipTemplate     = app.Flag("skip-templates", "Do not load the base template *.template files (--st)").Bool()
		files            = app.Arg("files", "Template files to process").ExistingFiles()
	)

	app.Flag("at", "short version of --all-templates").Hidden().BoolVar(listALlTemplates)
	app.Flag("ll", "short version of --log-level").Hidden().Int8Var(logLevel)
	app.Flag("ls", "short version of --log-simple").Hidden().BoolVar(logSimple)
	app.Flag("del", "").Hidden().StringVar(delimiters)
	app.Flag("st", "short version of --skip-templates").Hidden().BoolVar(skipTemplate)
	app.Flag("skip", "").Hidden().BoolVar(skipTemplate)
	app.Flag("skip-template", "").Hidden().BoolVar(skipTemplate)

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

	t := template.NewTemplate(createContext(*varFiles, *namedVars), *delimiters, *razorSyntax, *skipTemplate, *substitutes...)
	t.Quiet = *quiet
	t.Disabled = *disableRender
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
	*files = exclude(append(utils.MustFindFiles(*sourceFolder, *recursive, *followSymLinks, extend(*includePatterns)...), *files...), *excludedPatterns)

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

func createContext(varsFiles []string, namedVars []string) (context map[string]interface{}) {
	context = map[string]interface{}{}

	type fileDef struct {
		name    string
		value   interface{}
		unnamed bool
	}

	nameValuePairs := make([]fileDef, 0, len(varsFiles)+len(namedVars))
	for i := range varsFiles {
		nameValuePairs = append(nameValuePairs, fileDef{value: varsFiles[i]})
	}

	for i := range namedVars {
		data := make(map[string]interface{})
		if err := utils.ConvertData(namedVars[i], &data); err != nil {
			var fd fileDef
			fd.name, fd.value = utils.Split2(namedVars[i], "=")
			if fd.value == "" {
				fd = fileDef{"", fd.name, true}
			}
			nameValuePairs = append(nameValuePairs, fd)
			continue
		}
		for key, value := range utils.Flatten(data) {
			nameValuePairs = append(nameValuePairs, fileDef{key, value, false})
		}
	}

	for _, nv := range nameValuePairs {
		var loader func(string) (map[string]interface{}, error)
		filename, _ := reflect.ValueOf(nv.value).Interface().(string)
		if filename != "" {
			loader = func(filename string) (result map[string]interface{}, err error) {
				var content interface{}
				if err := utils.LoadData(filename, &content); err == nil {
					if content, isMap := content.(map[string]interface{}); isMap && nv.name == "" && !nv.unnamed {
						return content, nil
					}
					if nv.name == "" {
						nv.name = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
					}
					return map[string]interface{}{nv.name: content}, nil
				} else if _, isFileErr := err.(*os.PathError); !isFileErr {
					return nil, err
				}
				if nv.name == "" {
					nv.name = "DEFAULT"
				}
				return map[string]interface{}{nv.name: nv.value}, nil
			}
		}

		if loader == nil {
			context[nv.name] = nv.value
			continue
		}
		content, err := loader(filename)
		if err != nil {
			errors.Raise("Error %v while loading vars file %s", nv.value, err)
		}
		for key, value := range content {
			context[key] = value
		}

		// There is no content
		if len(content) == 0 && nv.unnamed {
			errors.Raise("--var parameter must be a file or assignation (name=value) %s", nv.value)
		}
	}
	return
}

func exclude(files []string, patterns []string) []string {
	if patterns = extend(patterns); len(patterns) == 0 {
		// There is no exclusion pattern, so we return the list of files as is
		return files
	}

	result := make([]string, 0, len(files))
	for i, file := range files {
		var excluded bool
		for _, pattern := range patterns {
			file = iif(filepath.IsAbs(pattern), file, filepath.Base(file)).(string)
			if excluded = errors.Must(filepath.Match(pattern, file)).(bool); excluded {
				template.Log.Noticef("%s ignored", files[i])
				break
			}
		}
		if !excluded {
			result = append(result, files[i])
		}
	}
	return result
}

func extend(values []string) []string {
	result := make([]string, 0, len(values))
	for i := range values {
		for _, sv := range strings.Split(values[i], ",") {
			sv = strings.TrimSpace(sv)
			if sv != "" {
				if strings.Contains(filepath.ToSlash(sv), "/") {
					// If the pattern contains /, we make it relative
					sv = errors.Must(filepath.Abs(sv)).(string)
				}
				result = append(result, sv)
			}
		}
	}
	return result
}

var iif = utils.IIf
