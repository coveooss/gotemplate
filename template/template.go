package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/utils"
	multicolor "github.com/coveooss/multilogger/color"
)

// String is an alias to collections.String
type String = collections.String

var templateMutex sync.Mutex

// Template let us extend the functionalities of base go template library.
type Template struct {
	*template.Template
	tempFolder     string
	substitutes    []utils.RegexReplacer
	context        interface{}
	constantKeys   []interface{}
	delimiters     []string
	parent         *Template
	folder         string
	children       map[string]*Template
	aliases        funcTableMap
	functions      funcTableMap
	options        OptionsSet
	optionsEnabled OptionsSet
}

// Environment variables that could be defined to override default behaviors.
const (
	EnvAcceptNoValue    = "GOTEMPLATE_NO_VALUE"
	EnvStrictErrorCheck = "GOTEMPLATE_STRICT_ERROR"
	EnvSubstitutes      = "GOTEMPLATE_SUBSTITUTES"
	EnvDebug            = "GOTEMPLATE_DEBUG"
	EnvExtensionPath    = "GOTEMPLATE_PATH"
)

const (
	noGoTemplate       = "no-gotemplate!"
	noRazor            = "no-razor!"
	explicitGoTemplate = "force-gotemplate!"

	pauseGoTemplate  = "gotemplate-pause!"
	resumeGoTemplate = "gotemplate-resume!"

	pauseRazor  = "razor-pause!"
	resumeRazor = "razor-resume!"
)

// Common variables
var (
	// ExtensionDepth the depth level of search of gotemplate extension from the current directory (default = 2).
	ExtensionDepth = 2
	toStrings      = collections.ToStrings
	acceptNoValue  = String(os.Getenv(EnvAcceptNoValue)).ParseBool()
	strictError    = String(os.Getenv(EnvStrictErrorCheck)).ParseBool()
	Print          = multicolor.Print
	Printf         = multicolor.Printf
	Println        = multicolor.Println
	ErrPrintf      = multicolor.ErrorPrintf
	ErrPrintln     = multicolor.ErrorPrintln
	ErrPrint       = multicolor.ErrorPrint
)

// IsRazor determines if the supplied code appears to have Razor code (using default delimiters).
func IsRazor(code string) bool { return strings.Contains(code, "@") }

// IsCode determines if the supplied code appears to have gotemplate code (using default delimiters).
func IsCode(code string) bool {
	return IsRazor(code) || strings.Contains(code, "{{") || strings.Contains(code, "}}")
}

// NewTemplate creates an Template object with default initialization.
func NewTemplate(folder string, context interface{}, delimiters string, options OptionsSet, substitutes ...string) (result *Template, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			result, err = nil, fmt.Errorf("%v", rec)
		}
	}()
	t := Template{Template: template.New("Main")}
	must(t.Parse(""))
	t.options = iif(options != nil, options, DefaultOptions()).(OptionsSet)
	if acceptNoValue {
		t.options[AcceptNoValue] = true
	}
	if strictError {
		t.options[StrictErrorCheck] = true
	}
	t.optionsEnabled = make(OptionsSet)
	t.folder, _ = filepath.Abs(iif(folder != "", folder, utils.Pwd()).(string))
	t.context = iif(context != nil, context, collections.CreateDictionary())
	t.aliases = make(funcTableMap)
	t.delimiters = []string{"{{", "}}", "@"}

	// Set the regular expression replacements
	baseSubstitutesRegex := []string{
		`/(?m)^\s*#!\s*$/`,
		`/"!Q!(?P<content>[\S]+?)!Q!"/${content}`,
		`/"!Q!(?P<content>.*?)!Q!"/"${content}"`,
	}
	if substitutesFromEnv := os.Getenv(EnvSubstitutes); substitutesFromEnv != "" {
		baseSubstitutesRegex = append(baseSubstitutesRegex, strings.Split(substitutesFromEnv, "\n")...)
	}
	t.substitutes = utils.InitReplacers(append(baseSubstitutesRegex, substitutes...)...)

	if t.options[Extension] {
		t.initExtension()
	}

	// Set the options supplied by caller
	t.init("")
	if delimiters != "" {
		for i, delimiter := range strings.Split(delimiters, ",") {
			if i == len(t.delimiters) {
				return nil, fmt.Errorf("Invalid delimiters '%s', must be a maximum of three comma separated parts", delimiters)
			}
			if delimiter != "" {
				t.delimiters[i] = delimiter
			}
		}
	}
	return &t, nil
}

// MustNewTemplate creates an Template object with default initialization.
// It panics if an error occurs.
func MustNewTemplate(folder string, context interface{}, delimiters string, options OptionsSet, substitutes ...string) *Template {
	return must(NewTemplate(folder, context, delimiters, options, substitutes...)).(*Template)
}

// TempFolder set temporary folder used by this template.
func (t *Template) TempFolder(folder string) *Template {
	t.tempFolder = folder
	return t
}

// GetNewContext returns a distinct context for each folder.
func (t *Template) GetNewContext(folder string, useCache bool) *Template {
	folder = iif(folder != "", folder, t.folder).(string)
	if context, found := t.children[folder]; useCache && found {
		return context
	}

	newTemplate := Template(*t)
	newTemplate.Template = template.New(folder)
	newTemplate.addFunctions(t.functions)
	newTemplate.addFunctions(t.aliases)
	newTemplate.init(folder)
	newTemplate.parent = t
	newTemplate.importTemplates(t)
	newTemplate.options = make(OptionsSet)
	if dict := t.Context(); dict.Len() > 0 {
		newTemplate.context = dict.Clone()
	}

	// We duplicate the options because the new context may alter them afterwhile and
	// it should not modify the original values.
	for k, v := range t.options {
		newTemplate.options[k] = v
	}

	if !useCache {
		return &newTemplate
	}
	// We register the new template as a child of the main template
	t.children[folder] = &newTemplate
	return t.children[folder]
}

// IsCode determines if the supplied code appears to have gotemplate code.
func (t *Template) IsCode(code string) bool {
	return !strings.Contains(code, noGoTemplate) && (t.IsRazor(code) || strings.Contains(code, t.LeftDelim()) || strings.Contains(code, t.RightDelim()))
}

// IsRazor determines if the supplied code appears to have Razor code.
func (t *Template) IsRazor(code string) bool {
	return strings.Contains(code, t.RazorDelim()) && !strings.Contains(code, noGoTemplate) && !strings.Contains(code, noRazor)
}

// LeftDelim returns the left delimiter.
func (t *Template) LeftDelim() string { return t.delimiters[0] }

// RightDelim returns the right delimiter.
func (t *Template) RightDelim() string { return t.delimiters[1] }

// RazorDelim returns the razor delimiter.
func (t *Template) RazorDelim() string { return t.delimiters[2] }

// SetOption allows setting of template option after initialization.
func (t *Template) SetOption(option Options, value bool) { t.options[option] = value }

func (t *Template) isTemplate(file string) bool {
	for i := range templateExt {
		if strings.HasSuffix(file, templateExt[i]) {
			return true
		}
	}
	return false
}

func (t *Template) initExtension() {
	ext := t.GetNewContext("", false)
	t.constantKeys = ext.constantKeys
	ext.options = DefaultOptions()

	var extensionfiles []string
	if extensionFolders := strings.TrimSpace(os.Getenv(EnvExtensionPath)); extensionFolders != "" {
		for _, path := range strings.Split(extensionFolders, string(os.PathListSeparator)) {
			if path != "" {
				files, _ := utils.FindFilesMaxDepth(path, ExtensionDepth, false, "*.gte")
				extensionfiles = append(extensionfiles, files...)
			}
		}
	}
	extensionfiles = append(extensionfiles, utils.MustFindFilesMaxDepth(ext.folder, ExtensionDepth, false, "*.gte")...)

	// Retrieve the template extension files
	for _, file := range extensionfiles {
		// We just load all the template files available to ensure that all template definition are loaded
		// We do not use ParseFiles because it names the template with the base name of the file
		// which result in overriding templates with the same base name in different folders.
		content := string(must(ioutil.ReadFile(file)).([]byte))

		// We execute the content, but we ignore errors. The goal is only to register the sub templates and aliases properly
		// We also do not ask to clone the context as we wish to let extension to be able to alter the supplied context
		if _, _, err := ext.processContentInternal(content, file, nil, 0, false, nil); err != nil {
			InternalLog.Error(err)
		}
	}

	// Add the children contexts to the main context
	for _, context := range ext.children {
		t.importTemplates(context)
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
	t.children = make(map[string]*Template)
	t.Delims(t.delimiters[0], t.delimiters[1])
	t.setConstant(false, "\n", "NL", "CR", "NEWLINE")
	t.setConstant(false, true, "true")
	t.setConstant(false, false, "false")
	t.setConstant(false, nil, "null")
}

func (t *Template) setConstant(stopOnFirst bool, value interface{}, names ...string) {
	c, err := collections.TryAsDictionary(t.context)
	if err != nil {
		return
	}

	context := c.AsMap()
	for i := range names {
		if val, isSet := context[names[i]]; !isSet {
			context[names[i]] = value
			t.constantKeys = append(t.constantKeys, names[i])
			if stopOnFirst {
				return
			}
		} else if isSet && reflect.DeepEqual(value, val) {
			return
		}
	}
}

// Import templates from another template.
func (t *Template) importTemplates(source *Template) {
	for _, subTemplate := range source.Templates() {
		if subTemplate.Name() != subTemplate.ParseName {
			t.AddParseTree(subTemplate.Name(), subTemplate.Tree)
		}
	}
}

// Add allows adding a value to the template context.
// The context must be a dictionnary to use that function, otherwise, it will panic.
func (t *Template) Add(key string, value interface{}) {
	collections.AsDictionary(t.context).Add(key, value)
}

// Merge allows adding multiple values to the template context.
// The context and values must both be dictionnary to use that function, otherwise, it will panic.
func (t *Template) Merge(values interface{}) {
	collections.AsDictionary(t.context).Add(key, collections.AsDictionary(values))
}

// Context returns the template context as a dictionnary if possible, otherwise, it returns null.
func (t *Template) Context() (result collections.IDictionary) {
	if result, _ = collections.TryAsDictionary(t.context); result == nil {
		result = collections.CreateDictionary()
	}
	return
}
