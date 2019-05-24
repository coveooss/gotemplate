package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/v3/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	templateExt    = []string{".gt", ".template"}
	linePrefix     = `template: ` + p(tagLocation, p(tagFile, `.*?`)+`:`+p(tagLine, `\d+`)+`(:`+p(tagCol, `\d+`)+`)?: `)
	reError        = regexp.MustCompile(linePrefix)
	execPrefix     = "(?s)^" + linePrefix + `executing ".*" at <` + p(tagCode, `.*`) + `>: `
	templateErrors = []string{
		execPrefix + `map has no entry for key "` + p(tagKey, `.*`) + `"`,
		execPrefix + `(?s)error calling (raise|assert): ` + p(tagMsg, `.*`),
		execPrefix + p(tagErr, `.*`),
		linePrefix + p(tagErr, `.*`),
	}
)

func p(name, expr string) string { return fmt.Sprintf("(?P<%s>%s)", name, expr) }

const (
	noValue      = "<no value>"
	noValueRepl  = "!NO_VALUE!"
	nilValue     = "<nil>"
	nilValueRepl = "!NIL_VALUE!"
	undefError   = `"` + noValue + `"`
	noValueError = "contains undefined value(s)"
	runError     = `"<RUN_ERROR>"`
	tagLine      = "line"
	tagCol       = "column"
	tagCode      = "code"
	tagMsg       = "message"
	tagLocation  = "location"
	tagFile      = "file"
	tagKey       = "key"
	tagErr       = "error"
)

// ProcessContent loads and runs the file template.
func (t Template) ProcessContent(content, source string) (result string, err error) {
	result, _, err = t.processContentInternal(content, source, nil, 0, true, nil)
	return
}

// ProcessTemplate loads and runs the template if it is a file, otherwise, it simply process the content.
func (t Template) ProcessTemplate(template, sourceFolder, targetFolder string) (resultFile string, err error) {
	return t.processTemplate(template, sourceFolder, targetFolder, nil)
}

// ProcessTemplates loads and runs the file template or execute the content if it is not a file.
func (t Template) ProcessTemplates(sourceFolder, targetFolder string, templates ...string) (resultFiles []string, err error) {
	return t.ProcessTemplatesWithHandler(sourceFolder, targetFolder, nil, templates...)
}

func (t Template) printResult(source, target, result string) (err error) {
	if utils.IsTerraformFile(target) {
		base := filepath.Base(target)
		tempFolder := must(ioutil.TempDir(t.TempFolder, base)).(string)
		tempFile := filepath.Join(tempFolder, base)
		err = ioutil.WriteFile(tempFile, []byte(result), 0644)
		if err != nil {
			return
		}
		err = utils.TerraformFormat(tempFile)
		bytes := must(ioutil.ReadFile(tempFile)).([]byte)
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
		log.Noticef("%s => %s", source, target)
	} else {
		log.Notice(target)
	}
	Print(result)
	if result != "" && terminal.IsTerminal(int(os.Stdout.Fd())) {
		Println()
	}

	return
}
