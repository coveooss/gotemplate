package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coveooss/gotemplate/v3/utils"
	"golang.org/x/term"
)

var templateExt = []string{".gt", ".template"}

func p(name, expr string) string { return fmt.Sprintf("(?P<%s>%s)", name, expr) }

const (
	noValue       = "<no value>"
	noValueRepl   = "!NO_VALUE!"
	nilValue      = "<nil>"
	nilValueRepl  = "!NIL_VALUE!"
	undefError    = `"` + noValue + `"`
	noValueError  = "contains undefined value(s)"
	runError      = `"<RUN_ERROR>"`
	tagLine       = "line"
	tagActualLine = "actual_line"
	tagCol        = "column"
	tagCode       = "code"
	tagMsg        = "message"
	tagLocation   = "location"
	tagFile       = "file"
	tagKey        = "key"
	tagErr        = "error"
	tagValue      = "value"
)

// ProcessContent loads and runs the file template.
func (t *Template) ProcessContent(content, source string) (result string, err error) {
	result, _, err = t.processContentInternal(content, source, nil, 0, true, nil)
	return
}

// ProcessTemplate loads and runs the template if it is a file, otherwise, it simply process the content.
func (t *Template) ProcessTemplate(template, sourceFolder, targetFolder string) (resultFile string, err error) {
	return t.processTemplate(template, sourceFolder, targetFolder, nil)
}

// ProcessTemplates loads and runs the file template or execute the content if it is not a file.
func (t *Template) ProcessTemplates(sourceFolder, targetFolder string, templates ...string) (resultFiles []string, err error) {
	return t.ProcessTemplatesWithHandler(sourceFolder, targetFolder, nil, templates...)
}

func (t *Template) printResult(source, target, result string, changed bool) (err error) {
	if utils.IsTerraformFile(target) {
		base := filepath.Base(target)
		tempFolder := must(os.CreateTemp(t.tempFolder, base)).(string)
		tempFile := filepath.Join(tempFolder, base)
		err = os.WriteFile(tempFile, []byte(result), 0644)
		if err != nil {
			return
		}
		err = utils.TerraformFormat(tempFile)
		bytes := must(os.ReadFile(tempFile)).([]byte)
		result = string(bytes)
	}

	if changed && !t.isTemplate(source) && !t.options[Overwrite] {
		source += ".original"
	}

	source = utils.Relative(t.folder, source)
	if relTarget := utils.Relative(t.folder, target); !strings.HasPrefix(relTarget, "../../../") {
		target = relTarget
	}
	if source != target {
		InternalLog.Infof("%s => %s", source, target)
	} else {
		InternalLog.Info(target)
	}
	Print(result)
	if result != "" && term.IsTerminal(int(os.Stdout.Fd())) {
		Println()
	}

	return
}
