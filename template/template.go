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
	"sort"
	"strings"
	"text/template"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	logging "github.com/op/go-logging"
	"golang.org/x/crypto/ssh/terminal"
)

var templateExt = []string{".gt", ".template"}

// Template let us extend the functionalities of base go template library.
type Template struct {
	*template.Template
	TempFolder           string
	substitutes          []utils.RegexReplacer
	context              interface{}
	delimiters           []string
	parent               *Template
	folder               string
	children             map[string]*Template
	aliases              funcTableMap
	functions            funcTableMap
	options              OptionsSet
	mathConstantInjected bool
}

// ExtensionDepth the depth level of search of gotemplate extension from the current directory (default = 2).
var ExtensionDepth = 2

var toStrings = utils.ToStrings

// NewTemplate creates an Template object with default initialization.
func NewTemplate(folder string, context interface{}, delimiters string, options OptionsSet, substitutes ...string) *Template {
	t := Template{Template: template.New("Main")}
	errors.Must(t.Parse(""))
	t.options = iif(options != nil, options, DefaultOptions()).(OptionsSet)
	t.folder, _ = filepath.Abs(iif(folder != "", folder, utils.Pwd()).(string))
	t.context = iif(context != nil, context, make(map[string]interface{}))
	t.aliases = make(funcTableMap)
	t.delimiters = []string{"{{", "}}", "@"}

	// Set the regular expression replacements
	baseRegex := []string{`/(?m)^\s*#!\s*$/`}
	t.substitutes = utils.InitReplacers(append(baseRegex, substitutes...)...)

	if t.options[Extension] {
		ext := t.GetNewContext("", false)
		ext.options = DefaultOptions()

		// We temporary set the logging level one grade lower
		logLevel := logging.GetLevel(logger)
		logging.SetLevel(logLevel-1, logger)
		defer func() { logging.SetLevel(logLevel, logger) }()

		// Retrieve the template extension files
		for _, file := range utils.MustFindFilesMaxDepth(ext.folder, ExtensionDepth, false, "*.gte") {
			// We just load all the template files available to ensure that all template definition are loaded
			// We do not use ParseFiles because it names the template with the base name of the file
			// which result in overriding templates with the same base name in different folders.
			content := string(errors.Must(ioutil.ReadFile(file)).([]byte))

			// We execute the content, but we ignore errors. The goal is only to register the sub templates and aliases properly
			if _, err := ext.ProcessContent(content, file); err != nil {
				Log.Noticef(color.RedString("Error while processing %v"), err)
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
	t.init("")
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

// IsCode determines if the supplied code appears to have gotemplate code.
func (t Template) IsCode(code string) bool {
	return t.IsRazor(code) || strings.Contains(code, t.LeftDelim()) || strings.Contains(code, t.RightDelim())
}

// IsRazor determines if the supplied code appears to have Razor code.
func (t Template) IsRazor(code string) bool {
	return strings.Contains(code, t.RazorDelim())
}

// LeftDelim returns the left delimiter.
func (t Template) LeftDelim() string { return t.delimiters[0] }

// RightDelim returns the right delimiter.
func (t Template) RightDelim() string { return t.delimiters[1] }

// RazorDelim returns the razor delimiter.
func (t Template) RazorDelim() string { return t.delimiters[2] }

// SetOption allows setting of template option after initialization.
func (t *Template) SetOption(option Options, value bool) { t.options[option] = value }

// ProcessContent loads and runs the file template.
func (t Template) ProcessContent(content, source string) (string, error) {
	content = t.substitute(content)

	if strings.HasPrefix(content, "#!") {
		// If the content starts with a Shebang operator including gotemplate, we remove the first line
		lines := strings.Split(content, "\n")
		if strings.Contains(lines[0], "gotemplate") {
			content = strings.Join(lines[1:], "\n")
			t.options[OutputStdout] = true
		}
	}

	content = string(t.applyRazor([]byte(content)))

	if t.options[RenderingDisabled] || !t.IsCode(content) {
		// There is no template element to evaluate or the template rendering is off
		return content, nil
	}

	Log.Notice("GoTemplate processing of", source)
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

// ProcessTemplate loads and runs the template if it is a file, otherwise, it simply process the content.
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
		if !t.options[Overwrite] {
			resultFile = fmt.Sprint(strings.TrimSuffix(resultFile, ext), ".generated", ext)
		}
	}

	if t.options[OutputStdout] {
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
	if !isTemplate && !t.options[Overwrite] {
		newName := template + ".original"
		Log.Noticef("%s => %s", utils.Relative(t.folder, template), utils.Relative(t.folder, newName))
		errors.Must(os.Rename(template, template+".original"))
	}

	if sourceFolder != targetFolder {
		errors.Must(os.MkdirAll(filepath.Dir(resultFile), 0777))
	}
	Log.Notice("Writing file", utils.Relative(t.folder, resultFile))

	if utils.IsShebangScript(result) {
		mode = 0755
	}

	if err = ioutil.WriteFile(resultFile, []byte(result), mode); err != nil {
		return
	}

	if isTemplate && t.options[Overwrite] && sourceFolder == targetFolder {
		os.Remove(template)
	}
	return
}

// ProcessTemplates loads and runs the file template or execute the content if it is not a file.
func (t Template) ProcessTemplates(sourceFolder, targetFolder string, templates ...string) (resultFiles []string, err error) {
	sourceFolder = iif(sourceFolder == "", t.folder, sourceFolder).(string)
	targetFolder = iif(targetFolder == "", t.folder, targetFolder).(string)
	resultFiles = make([]string, 0, len(templates))

	print := t.options[OutputStdout]

	var errors errors.Array
	for i := range templates {
		t.options[OutputStdout] = print // Some file may change this option at runtime, so we restore it back to its original value between each file
		resultFile, err := t.ProcessTemplate(templates[i], sourceFolder, targetFolder)
		if err == nil {
			if resultFile != "" {
				resultFiles = append(resultFiles, resultFile)
			}
		} else {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		err = errors
	}
	return
}

func (t Template) isTemplate(file string) bool {
	for i := range templateExt {
		if strings.HasSuffix(file, templateExt[i]) {
			return true
		}
	}
	return false
}

func (t Template) filterFunctions(all, category, detailed bool, filters ...string) []string {
	functions := t.getFunctions()
	if all && len(filters) == 0 {
		return functions
	}

	for i := range filters {
		filters[i] = strings.ToLower(filters[i])
	}

	filtered := make([]string, 0, len(functions))
	for i := range functions {
		funcInfo := t.functions[functions[i]]
		if funcInfo.alias != nil {
			if !all {
				continue
			}
			funcInfo = *funcInfo.alias
		}

		if len(filters) == 0 {
			filtered = append(filtered, functions[i])
			continue
		}

		search := strings.ToLower(functions[i])
		if category {
			search += strings.ToLower(funcInfo.group)
		}
		if detailed {
			search += strings.ToLower(funcInfo.desc)
		}

		for f := range filters {
			if strings.Contains(search, filters[f]) {
				filtered = append(filtered, functions[i])
				break
			}
		}
	}
	return filtered
}

// PrintFunctions outputs the list of functions available.
func (t Template) PrintFunctions(all, long, groupByCategory bool, filters ...string) {
	functions := t.filterFunctions(all, groupByCategory, long, filters...)

	maxLength := 0
	categories := make(map[string][]string)
	for i := range functions {
		var group string
		if groupByCategory {
			funcInfo := t.functions[functions[i]]
			if funcInfo.alias != nil {
				funcInfo = *funcInfo.alias
			}
			group = funcInfo.group
		}
		categories[group] = append(categories[group], functions[i])
		maxLength = int(math.Max(float64(len(functions[i])), float64(maxLength)))
	}

	keys := make([]string, 0, len(categories))
	for key := range categories {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	print := t.printFunctionsShort
	if long {
		print = t.printFunctionsDetailed
	}

	for _, key := range keys {
		if key != "" {
			title, link := utils.Split2(key, ", http")
			title = color.New(color.Underline, color.FgYellow).Sprint(title)
			if link != "" {
				link = color.BlackString(fmt.Sprintf(" http%s", link))
			}
			fmt.Printf("%s%s\n\n", title, link)
		}
		print(categories[key], maxLength)
		fmt.Println()
	}
}

func (t Template) printFunctionsShort(functions []string, maxLength int) {
	const nbColumn = 5
	l := len(functions)
	colLength := int(math.Ceil(float64(l) / float64(nbColumn)))
	for i := 0; i < colLength*nbColumn; i += nbColumn {
		for j := 0; j < nbColumn; j++ {
			pos := j*colLength + i/nbColumn
			if pos >= l {
				continue
			}
			item, extraLen := functions[pos], 0

			if t.functions[item].alias != nil {
				ex := len(color.HiBlackString(""))
				fmt.Printf("%-[1]*[2]s", maxLength+2+ex, color.HiBlackString(item))
			} else {
				fmt.Printf("%-[1]*[2]s", maxLength+2+extraLen, item)
			}
		}
		fmt.Println()
	}
}

func (t Template) printFunctionsDetailed(functions []string, maxLength int) {
	for i := range functions {
		col := color.HiBlueString
		funcInfo := t.functions[functions[i]]
		if funcInfo.alias != nil {
			funcInfo = *funcInfo.alias
			col = color.HiBlackString
		}
		var in, out string
		if funcInfo.f != nil {
			signature := reflect.ValueOf(funcInfo.f).Type()
			var parameters, outputs []string
			for i := 0; i < signature.NumIn(); i++ {
				arg := strings.Replace(fmt.Sprint(signature.In(i)), "interface {}", "interface{}", -1)
				var argName string
				if i < len(funcInfo.args) {
					argName = funcInfo.args[i]
				} else {
					if signature.IsVariadic() && i == signature.NumIn()-1 {
						argName = "args"
					} else {
						argName = fmt.Sprintf("arg%d", i+1)
					}
				}
				if signature.IsVariadic() && i == signature.NumIn()-1 {
					arg = "..." + arg[2:]
				}
				parameters = append(parameters, fmt.Sprintf("%s %s", argName, color.CyanString(arg)))
			}
			in = strings.Join(parameters, ", ")
			for i := 0; i < signature.NumOut(); i++ {
				outputs = append(outputs, strings.Replace(fmt.Sprint(signature.Out(i)), "interface {}", "interface{}", -1))
			}
			out = strings.Join(outputs, ", ")
		} else {
			in = "Check go template documentation"
		}
		if funcInfo.desc != "" {
			fmt.Printf(color.GreenString("%s\n"), funcInfo.desc)
		}
		fmt.Printf("%s(%s) %s\n", col(functions[i]), in, color.HiBlackString(out))
		fmt.Println()
	}
}

// PrintTemplates output the list of templates available.
func (t Template) PrintTemplates(all, long bool) {
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

// Initialize a new template with same attributes as the current context.
func (t *Template) init(folder string) {
	if folder != "" {
		t.folder, _ = filepath.Abs(folder)
	}
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

// Import templates from another template.
func (t *Template) importTemplates(source Template) {
	for _, subTemplate := range source.Templates() {
		if subTemplate.Name() != subTemplate.ParseName {
			t.AddParseTree(subTemplate.Name(), subTemplate.Tree)
		}
	}
}

// GetNewContext returns a distint context for each folder.
func (t Template) GetNewContext(folder string, useCache bool) *Template {
	folder = iif(folder != "", folder, t.folder).(string)
	if context, found := t.children[folder]; useCache && found {
		return context
	}

	newTemplate := Template(t)
	newTemplate.Template = template.New(folder)
	newTemplate.init(folder)
	newTemplate.parent = &t
	newTemplate.AddFunctions(t.aliases)
	newTemplate.importTemplates(t)
	newTemplate.options = make(OptionsSet)

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

	if !t.isTemplate(source) && !t.options[Overwrite] {
		source += ".original"
	}

	source = utils.Relative(t.folder, source)
	if relTarget := utils.Relative(t.folder, target); !strings.HasPrefix(relTarget, "../../../") {
		target = relTarget
	}
	if source != target {
		Log.Noticef("%s => %s", source, target)
	} else {
		Log.Notice(target)
	}
	fmt.Print(result)
	if result != "" && terminal.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println()
	}

	return
}

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}
