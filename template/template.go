package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/types"
	"github.com/coveo/gotemplate/utils"
	logging "github.com/op/go-logging"
)

// Template let us extend the functionalities of base go template library.
type Template struct {
	*template.Template
	TempFolder     string
	substitutes    []utils.RegexReplacer
	context        interface{}
	delimiters     []string
	parent         *Template
	folder         string
	children       map[string]*Template
	aliases        funcTableMap
	functions      funcTableMap
	options        OptionsSet
	optionsEnabled OptionsSet
}

// ExtensionDepth the depth level of search of gotemplate extension from the current directory (default = 2).
var ExtensionDepth = 2

var toStrings = types.ToStrings

// NewTemplate creates an Template object with default initialization.
func NewTemplate(folder string, context interface{}, delimiters string, options OptionsSet, substitutes ...string) *Template {
	t := Template{Template: template.New("Main")}
	errors.Must(t.Parse(""))
	t.options = iif(options != nil, options, DefaultOptions()).(OptionsSet)
	t.optionsEnabled = make(OptionsSet)
	t.folder, _ = filepath.Abs(iif(folder != "", folder, utils.Pwd()).(string))
	t.context = iif(context != nil, context, make(dictionary))
	t.aliases = make(funcTableMap)
	t.delimiters = []string{"{{", "}}", "@"}

	// Set the regular expression replacements
	baseRegex := []string{`/(?m)^\s*#!\s*$/`}
	t.substitutes = utils.InitReplacers(append(baseRegex, substitutes...)...)

	if t.options[Extension] {
		t.initExtension()
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
	newTemplate.addFunctions(t.aliases)
	newTemplate.importTemplates(t)
	newTemplate.options = make(OptionsSet)

	// We register the new template as a child of the main template
	if !useCache {
		return &newTemplate
	}
	t.children[folder] = &newTemplate
	return t.children[folder]
}

// IsCode determines if the supplied code appears to have gotemplate code.
func (t Template) IsCode(code string) bool {
	return t.IsRazor(code) || strings.Contains(code, t.LeftDelim()) || strings.Contains(code, t.RightDelim())
}

// IsRazor determines if the supplied code appears to have Razor code.
func (t Template) IsRazor(code string) bool {
	return strings.Contains(code, t.RazorDelim()) || strings.Contains(code, ":=")
}

// LeftDelim returns the left delimiter.
func (t Template) LeftDelim() string { return t.delimiters[0] }

// RightDelim returns the right delimiter.
func (t Template) RightDelim() string { return t.delimiters[1] }

// RazorDelim returns the razor delimiter.
func (t Template) RazorDelim() string { return t.delimiters[2] }

// SetOption allows setting of template option after initialization.
func (t *Template) SetOption(option Options, value bool) { t.options[option] = value }

func (t Template) isTemplate(file string) bool {
	for i := range templateExt {
		if strings.HasSuffix(file, templateExt[i]) {
			return true
		}
	}
	return false
}

func (t *Template) initExtension() {
	ext := t.GetNewContext("", false)
	ext.options = DefaultOptions()

	// We temporary set the logging level one grade lower
	logLevel := logging.GetLevel(logger)
	logging.SetLevel(logLevel-1, logger)
	defer func() { logging.SetLevel(logLevel, logger) }()

	var extensionfiles []string
	if extensionFolders := strings.TrimSpace(os.Getenv("GOTEMPLATE_PATH")); extensionFolders != "" {
		for _, path := range strings.Split(extensionFolders, string(os.PathListSeparator)) {
			files, _ := utils.FindFilesMaxDepth(path, ExtensionDepth, false, "*.gte")
			extensionfiles = append(extensionfiles, files...)
		}
	}
	extensionfiles = append(extensionfiles, utils.MustFindFilesMaxDepth(ext.folder, ExtensionDepth, false, "*.gte")...)

	// Retrieve the template extension files
	for _, file := range extensionfiles {
		// We just load all the template files available to ensure that all template definition are loaded
		// We do not use ParseFiles because it names the template with the base name of the file
		// which result in overriding templates with the same base name in different folders.
		content := string(errors.Must(ioutil.ReadFile(file)).([]byte))

		// We execute the content, but we ignore errors. The goal is only to register the sub templates and aliases properly
		if _, err := ext.ProcessContent(content, file); err != nil {
			log.Error(err)
		}
	}

	// Add the children contexts to the main context
	for _, context := range ext.children {
		t.importTemplates(*context)
	}

	// We reset the list of templates
	t.children = make(map[string]*Template)
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
	c, err := types.AsDictionary(t.context)
	if err != nil {
		return
	}

	context := c.AsMap()
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
