package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
)

const templateExt = ".template"

// Template let us extend the functionalities of base template
type Template struct {
	*template.Template
	Overwrite    bool
	Quiet        bool
	OutputStdout bool
	TempFolder   string
	RazorSyntax  bool
	substitutes  []utils.RegexReplacer
	context      interface{}
	delimiters   []string
	parent       *Template
	folder       string
	children     map[string]*Template
	aliases      *map[string]interface{}
}

// NewTemplate creates an Template object with default initialization
func NewTemplate(context interface{}, delimiters string, razor, skipTemplates bool, substitutes ...string) *Template {
	t := Template{Template: template.New("Main")}
	errors.Must(t.Parse(""))
	t.context = context
	t.aliases = &map[string]interface{}{}
	t.RazorSyntax = razor

	// Set the template delimiters
	if delimiters == "" {
		delimiters = "{{,}}"
		if t.RazorSyntax {
			delimiters += ",@"
		}
	}
	t.delimiters = strings.SplitN(delimiters, ",", 3)
	switch len(t.delimiters) {
	case 2:
		if t.RazorSyntax {
			t.delimiters = append(t.delimiters, "@")
		}
	case 3:
		break
	default:
		errors.Raise("Invalid delimiters '%s', must be two comma separated parts", delimiters)
	}

	t.init(utils.Pwd())

	// Set the regular expression replacements
	baseRegex := []string{`/(?m)^\s*#!\s*$/`}
	t.substitutes = utils.InitReplacers(append(baseRegex, substitutes...)...)

	if !skipTemplates {
		// Retrieve the template files
		for _, file := range utils.MustFindFiles(t.folder, true, true, "*.template") {
			// We just load all the template files available to ensure that all template definition are loaded
			// We do not use ParseFiles because it names the template with the base name of the file
			// which result in overriding templates with the same base name in different folders.
			content := string(errors.Must(ioutil.ReadFile(file)).([]byte))

			// We execute the content, but we ignore errors. The goal is only to register the sub templates and aliases properly
			t.ProcessContent(content, file)
		}

		// Add the children contexts to the main context
		for _, context := range t.children {
			t.importTemplates(*context)
		}

		// We reset the list of templates
		t.children = make(map[string]*Template)
	}

	return &t
}

// ProcessContent loads and runs the file template
func (t Template) ProcessContent(content, source string) (string, error) {
	log.Noticef("GoTemplate processing of %s", source)
	if t.RazorSyntax {
		content = string(t.applyRazor([]byte(content)))
	}

	context := t.getContext(filepath.Dir(source))
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

// ProcessFile loads and runs the file template
func (t Template) ProcessFile(file, sourceFolder, targetFolder string) (resultFile string, err error) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	content := t.substitute(string(fileContent))

	result, err := t.ProcessContent(content, file)
	if err != nil {
		return
	}

	resultFile = getTargetFile(strings.TrimSuffix(file, templateExt), sourceFolder, targetFolder)

	isTemplate := strings.HasSuffix(file, templateExt)
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

	if sourceFolder == targetFolder && result == content {
		return "", nil
	}

	if t.OutputStdout {
		err = t.printResult(file, resultFile, result)
		if err != nil {
			errors.Print(err)
		}
		return "", nil
	}

	mode := errors.Must(os.Stat(file)).(os.FileInfo).Mode()
	if !isTemplate && !t.Overwrite {
		newName := file + ".original"
		t.trace("%s => %s", utils.Relative(t.folder, file), utils.Relative(t.folder, newName))
		errors.Must(os.Rename(file, file+".original"))
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
		os.Remove(file)
	}
	return
}

// ProcessFiles loads and runs the file template
func (t Template) ProcessFiles(sourceFolder, targetFolder string, files ...string) (resultFiles []string, errors errors.Array) {
	resultFiles = make([]string, 0, len(files))

	for _, file := range files {
		resultFile, err := t.ProcessFile(file, sourceFolder, targetFolder)
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

// PrintFunctions output the list of functions available
func (t Template) PrintFunctions() {
	functions := t.getFunctions()

	const nbColumn = 5
	colLength := int(math.Ceil(float64(len(functions)) / float64(nbColumn)))
	maxLength := 0

	// Initialize the columns to sort function per column
	var list [nbColumn][]string
	for i := range list {
		list[i] = make([]string, colLength)
	}

	// Place functions into columns
	for i, function := range functions {
		column := list[i/colLength]
		column[i%colLength] = function
		maxLength = int(math.Max(float64(len(function)), float64(maxLength)))
	}

	// Print the columns
	for i := range list[0] {
		for _, column := range list {
			fmt.Printf("%-[1]*[2]s", maxLength+2, column[i])
		}
		fmt.Println()
	}
	fmt.Fprintln(os.Stderr)
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
			fmt.Fprintf(os.Stderr, "%-[3]*[1]s %[2]s\n", name, faint(utils.Relative(t.folder, tpl.ParseName)), maxLen)
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
}

// Import templates from another template
func (t *Template) importTemplates(source Template) {
	for _, subTemplate := range source.Templates() {
		if subTemplate.Name() != subTemplate.ParseName {
			t.AddParseTree(subTemplate.Name(), subTemplate.Tree)
		}
	}
}

// Get a distint context for each folder
func (t Template) getContext(folder string) *Template {
	if context, ok := t.children[folder]; ok {
		return context
	}

	newTemplate := Template(t)
	newTemplate.Template = template.New(folder)
	newTemplate.init(folder)
	newTemplate.parent = &t
	newTemplate.Funcs(*t.aliases)
	newTemplate.importTemplates(t)

	// We register the new template as a child of the main template
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

	if !strings.HasSuffix(source, templateExt) && !t.Overwrite {
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
	fmt.Printf("%s\n\n", result)
	return
}

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}
