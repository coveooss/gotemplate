package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	logging "github.com/op/go-logging"
)

var templateExt = []string{".gt", ".template"}

// Template let us extend the functionalities of base template
type Template struct {
	*template.Template
	Overwrite    bool
	Quiet        bool
	OutputStdout bool
	TempFolder   string
	RazorSyntax  bool
	Disabled     bool
	substitutes  []utils.RegexReplacer
	context      interface{}
	delimiters   []string
	parent       *Template
	folder       string
	children     map[string]*Template
	aliases      *map[string]interface{}
	functions    template.FuncMap
	include      Libraries
}

type Libraries int

const (
	_ Libraries = iota << 2
	Math
	Sprig
	AllLibraries = 0xFFFF
)

var toStrings = utils.ToStrings

// NewTemplate creates an Template object with default initialization
func NewTemplate(context interface{}, delimiters string, libraries Libraries, razor, skipTemplates bool, substitutes ...string) *Template {
	t := Template{Template: template.New("Main")}
	errors.Must(t.Parse(""))
	t.context = context
	t.aliases = &map[string]interface{}{}
	t.delimiters = []string{"{{", "}}", "@"}

	// Set the regular expression replacements
	baseRegex := []string{`/(?m)^\s*#!\s*$/`}
	t.substitutes = utils.InitReplacers(append(baseRegex, substitutes...)...)

	if !skipTemplates {
		ext := t.GetNewContext(utils.Pwd(), false)
		ext.RazorSyntax = true
		ext.include = AllLibraries
		ext.init(utils.Pwd())

		// We temporary set the logging level one grade lower
		logLevel := logging.GetLevel(logger)
		logging.SetLevel(logLevel-1, logger)
		defer func() { logging.SetLevel(logLevel, logger) }()

		// Retrieve the template extension files
		for _, file := range utils.MustFindFiles(t.folder, true, true, "*.gte") {
			// We just load all the template files available to ensure that all template definition are loaded
			// We do not use ParseFiles because it names the template with the base name of the file
			// which result in overriding templates with the same base name in different folders.
			content := string(errors.Must(ioutil.ReadFile(file)).([]byte))

			// We execute the content, but we ignore errors. The goal is only to register the sub templates and aliases properly
			if _, err := ext.ProcessContent(content, file); err != nil {
				Log.Notice(color.New(color.FgRed).Sprintf("Error while processing %v", err))
			}
		}

		// Add the children contexts to the main context
		for _, context := range ext.children {
			t.importTemplates(*context)
		}

		// We reset the list of templates
		t.children = make(map[string]*Template)
	}

	// Set the options supplied by caller
	t.include = libraries
	t.RazorSyntax = razor
	t.init(utils.Pwd())
	if delimiters != "" {
		for i, delimiter := range strings.Split(delimiters, ",") {
			if i == len(t.delimiters) {
				errors.Raise("Invalid delimiters '%s', must be two comma separated parts", delimiters)
			}
			t.delimiters[i] = delimiter
		}
	}
	return &t
}

// IsCode determines if the supplied code appears to have gotemplate code
func (t Template) IsCode(code string) bool {
	return t.IsRazor(code) || strings.Contains(code, t.LeftDelim()) || strings.Contains(code, t.RightDelim())
}

// IsRazor determines if the supplied code appears to have razor code
func (t Template) IsRazor(code string) bool {
	return strings.Contains(code, t.RazorDelim())
}

// LeftDelim returns the left delimiter
func (t Template) LeftDelim() string { return t.delimiters[0] }

// RightDelim returns the right delimiter
func (t Template) RightDelim() string { return t.delimiters[1] }

// RazorDelim returns the razor delimiter
func (t Template) RazorDelim() string { return t.delimiters[2] }

// Funcs overrides the base Funcs to get a copy of the functins added to the template
func (t *Template) Funcs(funcMap template.FuncMap) *Template {
	if t.functions == nil {
		t.functions = make(template.FuncMap)
	}
	for key, value := range funcMap {
		t.functions[key] = value
	}
	t.Template.Funcs(funcMap)
	return t
}

// ProcessContent loads and runs the file template
func (t Template) ProcessContent(content, source string) (string, error) {
	Log.Notice("GoTemplate processing of", source)
	content = t.substitute(content)

	if strings.HasPrefix(content, "#!") {
		// If the content starts with a Shebang operator including gotemplate, we remove the first line
		lines := strings.Split(content, "\n")
		if strings.Contains(lines[0], "gotemplate") {
			content = strings.Join(lines[1:], "\n")
		}
	}

	if t.RazorSyntax && t.IsRazor(content) {
		content = string(t.applyRazor([]byte(content)))
	}

	if t.Disabled || !t.IsCode(content) {
		// There is no template element to evaluate or the template rendering is off
		return content, nil
	}

	context := t.GetNewContext(filepath.Dir(source), true)
	newTemplate, err := context.New(source).Parse(content)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = newTemplate.Execute(&out, t.context)
	if err != nil {
		switch err.(type) {
		case template.ExecError:
			return "", err
		default:
			errors.Must(err)
		}
	}
	return t.substitute(out.String()), nil
}

// ProcessTemplate loads and runs the template if it is a file, otherwise, it simply process the content
func (t Template) ProcessTemplate(template, sourceFolder, targetFolder string) (resultFile string, err error) {
	isCode := t.IsCode(template)
	var content string

	if isCode {
		content = template
		template = "."
	} else if fileContent, err := ioutil.ReadFile(template); err == nil {
		content = string(fileContent)
	} else {
		return "", err
	}

	result, err := t.ProcessContent(content, template)
	if err != nil {
		return
	}

	if isCode {
		fmt.Println(result)
		return "", nil
	}
	resultFile = template
	for i := range templateExt {
		resultFile = strings.TrimSuffix(resultFile, templateExt[i])
	}
	resultFile = getTargetFile(resultFile, sourceFolder, targetFolder)
	isTemplate := t.isTemplate(template)
	if isTemplate {
		ext := path.Ext(resultFile)
		if strings.TrimSpace(result)+ext == "" {
			// We do not save anything for an empty resulting template that has no extension
			return "", nil
		}
		if !t.Overwrite {
			resultFile = fmt.Sprint(strings.TrimSuffix(resultFile, ext), ".generated", ext)
		}
	}

	if t.OutputStdout {
		err = t.printResult(template, resultFile, result)
		if err != nil {
			errors.Print(err)
		}
		return "", nil
	}

	if sourceFolder == targetFolder && result == content {
		return "", nil
	}

	mode := errors.Must(os.Stat(template)).(os.FileInfo).Mode()
	if !isTemplate && !t.Overwrite {
		newName := template + ".original"
		t.trace("%s => %s", utils.Relative(t.folder, template), utils.Relative(t.folder, newName))
		errors.Must(os.Rename(template, template+".original"))
	}

	if sourceFolder != targetFolder {
		errors.Must(os.MkdirAll(filepath.Dir(resultFile), 0777))
	}
	t.trace(utils.Relative(t.folder, resultFile))

	if utils.IsShebangScript(result) {
		mode = 0755
	}

	if err = ioutil.WriteFile(resultFile, []byte(result), mode); err != nil {
		return
	}

	if isTemplate && t.Overwrite && sourceFolder == targetFolder {
		os.Remove(template)
	}
	return
}

// ProcessTemplates loads and runs the file template or execute the content if it is not a file
func (t Template) ProcessTemplates(sourceFolder, targetFolder string, templates ...string) (resultFiles []string, errors errors.Array) {
	resultFiles = make([]string, 0, len(templates))

	for i := range templates {
		resultFile, err := t.ProcessTemplate(templates[i], sourceFolder, targetFolder)
		if err == nil {
			if resultFile != "" {
				resultFiles = append(resultFiles, resultFile)
			}
		} else {
			errors = append(errors, err)
		}
	}
	return
}

func (t Template) trace(format string, args ...interface{}) {
	if t.Quiet {
		return
	}
	fmt.Fprintln(os.Stderr, color.GreenString(fmt.Sprintf(format, args...)))
}

func (t Template) isTemplate(file string) bool {
	for i := range templateExt {
		if strings.HasSuffix(file, templateExt[i]) {
			return true
		}
	}
	return false
}

func (t Template) filterFunctions(filters ...string) []string {
	functions := t.getFunctions()
	if len(filters) == 0 {
		return functions
	}

	for i := range filters {
		filters[i] = strings.ToLower(filters[i])
	}

	filtered := make([]string, 0, len(functions))
	for i := range functions {
		for f := range filters {
			if strings.Contains(strings.ToLower(functions[i]), filters[f]) {
				filtered = append(filtered, functions[i])
				break
			}
		}
	}
	return filtered
}

func (t Template) printFunctionsDetailed(functions []string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Function\tParameters\tOutputs")
	fmt.Fprintln(w, "--------\t----------\t-------")
	for i := range functions {
		function := t.functions[functions[i]]
		var in, out string
		if function != nil {
			signature := reflect.ValueOf(function).Type()
			var parameters, outputs []string
			for i := 0; i < signature.NumIn(); i++ {
				parameters = append(parameters, strings.Replace(fmt.Sprint(signature.In(i)), "interface {}", "generic", -1))
			}
			in = strings.Join(parameters, ", ")
			for i := 0; i < signature.NumOut(); i++ {
				outputs = append(outputs, strings.Replace(fmt.Sprint(signature.Out(i)), "interface {}", "generic", -1))
			}
			out = strings.Join(outputs, ", ")
		} else {
			in = "Check go template documentation"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", functions[i], in, out)
	}
	w.Flush()
}

// PrintFunctions output the list of functions available
func (t Template) PrintFunctions(filters ...string) {
	const nbColumn = 5
	functions := t.filterFunctions(filters...)
	if getLogLevel() >= 3 {
		t.printFunctionsDetailed(functions)
		return
	}

	colLength := int(math.Ceil(float64(len(functions)) / float64(nbColumn)))

	// Initialize the columns to sort function per column
	var list [nbColumn][]string
	for i := range list {
		list[i] = make([]string, colLength)
	}

	// Place functions into columns
	maxLength := 0
	for i := range functions {
		column := list[i/colLength]
		column[i%colLength] = functions[i]
		maxLength = int(math.Max(float64(len(functions[i])), float64(maxLength)))
	}

	// Print the columns
	for i := range list[0] {
		for _, column := range list {
			fmt.Printf("%-[1]*[2]s", maxLength+2, column[i])
		}
		fmt.Println()
	}
	fmt.Println()
}

// PrintTemplates output the list of templates available
func (t Template) PrintTemplates(all bool) {
	templates := t.getTemplateNames()
	var maxLen int
	for _, template := range templates {
		t := t.Lookup(template)
		if len(template) > maxLen && template != t.ParseName {
			maxLen = len(template)
		}
	}

	faint := color.New(color.Faint).SprintfFunc()

	for _, template := range templates {
		tpl := t.Lookup(template)
		if all || tpl.Name() != tpl.ParseName {
			name := tpl.Name()
			if tpl.Name() == tpl.ParseName {
				name = ""
			}
			folder := utils.Relative(t.folder, tpl.ParseName)
			if folder+name != "." {
				fmt.Fprintf(os.Stderr, "%-[3]*[1]s %[2]s\n", name, faint(folder), maxLen)
			}
		}
	}
	fmt.Fprintln(os.Stderr)
}

// Initialize a new template with same attributes as the current context
func (t *Template) init(folder string) {
	t.folder = folder
	t.addFuncs()
	t.Parse("")
	t.children = make(map[string]*Template)
	t.Delims(t.delimiters[0], t.delimiters[1])
	t.setConstant(false, "\n", "NL", "CR", "NEWLINE")
}

func (t *Template) setConstant(stopOnFirst bool, value interface{}, names ...string) {
	context, isMap := t.context.(map[string]interface{})
	if !isMap {
		return
	}
	for i := range names {
		if val, isSet := context[names[i]]; !isSet {
			context[names[i]] = value
			if stopOnFirst {
				return
			}
		} else if isSet && reflect.DeepEqual(value, val) {
			return
		}
	}
}

// Import templates from another template
func (t *Template) importTemplates(source Template) {
	for _, subTemplate := range source.Templates() {
		if subTemplate.Name() != subTemplate.ParseName {
			t.AddParseTree(subTemplate.Name(), subTemplate.Tree)
		}
	}
}

// GetNewContext returns a distint context for each folder
func (t Template) GetNewContext(folder string, useCache bool) *Template {
	folder, _ = filepath.Abs(folder)
	if context, found := t.children[folder]; useCache && found {
		return context
	}

	newTemplate := Template(t)
	newTemplate.Template = template.New(folder)
	newTemplate.init(folder)
	newTemplate.parent = &t
	newTemplate.Funcs(*t.aliases)
	newTemplate.importTemplates(t)

	// We register the new template as a child of the main template
	if !useCache {
		return &newTemplate
	}
	t.children[folder] = &newTemplate
	return t.children[folder]
}

func (t Template) printResult(source, target, result string) (err error) {
	if utils.IsTerraformFile(target) {
		base := filepath.Base(target)
		tempFolder := errors.Must(ioutil.TempDir(t.TempFolder, base)).(string)
		tempFile := filepath.Join(tempFolder, base)
		err = ioutil.WriteFile(tempFile, []byte(result), 0644)
		if err != nil {
			return
		}
		err = utils.TerraformFormat(tempFile)
		bytes := errors.Must(ioutil.ReadFile(tempFile)).([]byte)
		result = string(bytes)
	}

	if !t.isTemplate(source) && !t.Overwrite {
		source += ".original"
	}

	source = utils.Relative(t.folder, source)
	if relTarget := utils.Relative(t.folder, target); !strings.HasPrefix(relTarget, "../../../") {
		target = relTarget
	}
	if source != target {
		t.trace("# %s => %s", source, target)
	} else {
		t.trace("# %s", target)
	}
	fmt.Printf(result)
	fmt.Fprintln(os.Stderr)
	return
}

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}
