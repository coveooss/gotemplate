package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"golang.org/x/crypto/ssh/terminal"
)

var templateExt = []string{".gt", ".template"}

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

	log.Notice("GoTemplate processing of", source)
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
		log.Noticef("%s => %s", utils.Relative(t.folder, template), utils.Relative(t.folder, newName))
		errors.Must(os.Rename(template, template+".original"))
	}

	if sourceFolder != targetFolder {
		errors.Must(os.MkdirAll(filepath.Dir(resultFile), 0777))
	}
	log.Notice("Writing file", utils.Relative(t.folder, resultFile))

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
		log.Noticef("%s => %s", source, target)
	} else {
		log.Notice(target)
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
