package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/coveooss/multilogger/errors"
	"github.com/fatih/color"
)

// CustomHandler allows caller to supply a custom handler during the evaluation of multiple template files
type CustomHandler func(name, original string, result *string, changed bool, status error) (bool, error)

// ProcessTemplatesWithHandler loads and runs the file template or execute the content if it is not a file and call the custom handler between after each template.
func (t *Template) ProcessTemplatesWithHandler(sourceFolder, targetFolder string, handler CustomHandler, templates ...string) (resultFiles []string, err error) {
	sourceFolder = iif(sourceFolder == "", t.folder, sourceFolder).(string)
	targetFolder = iif(targetFolder == "", t.folder, targetFolder).(string)
	resultFiles = make([]string, 0, len(templates))

	print := t.options[OutputStdout]

	var errors errors.Array
	for i := range templates {
		t.options[OutputStdout] = print // Some file may change this option at runtime, so we restore it back to its originalSourceLines value between each file
		resultFile, err := t.processTemplate(templates[i], sourceFolder, targetFolder, handler)
		if err == nil {
			if resultFile != "" {
				resultFiles = append(resultFiles, resultFile)
			}
		} else {
			errors = append(errors, err)
		}
	}
	return resultFiles, errors.AsError()
}

func (t *Template) processTemplate(template, sourceFolder, targetFolder string, handler CustomHandler) (resultFile string, err error) {
	isCode := t.IsCode(template)
	var content string

	if isCode {
		content = template
		template = "."
	} else if fileContent, fileError := ioutil.ReadFile(template); fileError == nil {
		content = string(fileContent)
	} else {
		err = fileError
		return
	}

	result, changed, err := t.processContentInternal(content, template, nil, 0, true, handler)

	if err != nil {
		return
	}

	if isCode {
		// This occurs when gotemplate code has been supplied as a filename. In that case, we simply render
		// the result to the stdout
		Println(result)
		return
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
			resultFile = ""
			return
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
		resultFile = ""
		return
	}

	if sourceFolder == targetFolder && !changed {
		resultFile = ""
		return
	}

	mode := must(os.Stat(template)).(os.FileInfo).Mode()
	if !isTemplate && !t.options[Overwrite] {
		newName := template + ".original"
		InternalLog.Infof("%s => %s", utils.Relative(t.folder, template), utils.Relative(t.folder, newName))
		must(os.Rename(template, template+".original"))
	}

	if sourceFolder != targetFolder {
		must(os.MkdirAll(filepath.Dir(resultFile), 0777))
	}
	InternalLog.Infoln("Writing file", utils.Relative(t.folder, resultFile))

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

func (t *Template) processContentInternal(originalContent, source string, originalSourceLines []string, retryCount int, cloneContext bool, handler CustomHandler) (result string, changed bool, err error) {
	th := errorHandler{
		Template: t,
		Filename: source,
		Source:   originalContent,
		Code:     originalContent,
		Lines:    originalSourceLines,
		Try:      retryCount,
	}

	topCall := th.Lines == nil

	pausingIsEnabled := strings.Contains(originalContent, pauseGoTemplate) || strings.Contains(originalContent, pauseRazor)

	// When pausing templating, we replace delimiters with dummy strings and then we revert the replacements when the processing is complete
	leftDelimReplacement, rightDelimReplacement, razorDelimReplacement := "$&paused-left&$", "$&paused-right&$", "$&paused-razor&$"
	revertReplacements := func(template string) string {
		if pausingIsEnabled {
			template = strings.ReplaceAll(template, leftDelimReplacement, t.LeftDelim())
			template = strings.ReplaceAll(template, rightDelimReplacement, t.RightDelim())
			template = strings.ReplaceAll(template, razorDelimReplacement, t.RazorDelim())
		}
		return template
	}

	if topCall {
		if handler != nil {
			defer func() {
				var changedByHandler bool
				changedByHandler, err = handler(th.Filename, th.Source, &result, changed, err)
				changed = changed || changedByHandler
			}()
		}

		th.Code = t.substitute(th.Code)

		if strings.HasPrefix(th.Code, "#!") {
			// If the content starts with a Shebang operator including gotemplate, we remove the first line
			lines := strings.Split(th.Code, "\n")
			if strings.Contains(lines[0], "gotemplate") {
				th.Code = strings.Join(lines[1:], "\n")
				t.options[OutputStdout] = true
			}
		}

		if pausingIsEnabled {
			splitLines := strings.Split(th.Code, "\n")
			isGoTemplatePaused, isRazorPaused := false, false
			for index, line := range splitLines {
				isGoTemplatePaused = (isGoTemplatePaused || strings.Contains(line, pauseGoTemplate)) && !strings.Contains(line, resumeGoTemplate)
				isRazorPaused = (isRazorPaused || strings.Contains(line, pauseRazor)) && !strings.Contains(line, resumeRazor)
				if isRazorPaused || isGoTemplatePaused {
					splitLines[index] = strings.ReplaceAll(splitLines[index], t.RazorDelim(), razorDelimReplacement)
				}
				if isGoTemplatePaused {
					splitLines[index] = strings.ReplaceAll(splitLines[index], t.LeftDelim(), leftDelimReplacement)
					splitLines[index] = strings.ReplaceAll(splitLines[index], t.RightDelim(), rightDelimReplacement)
				}
			}
			th.Code = strings.Join(splitLines, "\n")
		}

		var razor []byte
		razor, _ = t.applyRazor([]byte(th.Code))
		th.Code = string(razor)

		if t.options[RenderingDisabled] || !t.IsCode(th.Code) {
			// There is no template element to evaluate or the template rendering is off
			return revertReplacements(th.Code), false, nil
		}

		InternalLog.Debug("GoTemplate processing of ", th.Filename)

		if !t.options[AcceptNoValue] {
			// We replace any pre-existing no value to avoid false error detection
			th.Code = strings.Replace(th.Code, noValue, noValueRepl, -1)
			th.Code = strings.Replace(th.Code, nilValue, nilValueRepl, -1)
		}

		defer func() {
			// If we get errors and the file is not an explicit gotemplate file, or contains
			// gotemplate! or the strict error check mode is not enabled, we simply
			// add a trace with the error content and return the content unaltered
			if err != nil {
				strictMode := t.options[StrictErrorCheck]
				strictMode = strictMode || strings.Contains(th.Source, explicitGoTemplate)
				extension := filepath.Ext(th.Filename)
				strictMode = strictMode || (extension != "" && strings.Contains(".gt,.gte,.template", extension))
				if !(strictMode) {
					InternalLog.Errorf("Ignored gotemplate error in %s (file left unchanged):\n%s", color.CyanString(th.Filename), err.Error())
					result, err = th.Source, nil
				}
			}
		}()
	}

	context := t.GetNewContext(filepath.Dir(th.Filename), true)
	newTemplate := context.New(th.Filename)

	if topCall {
		newTemplate.Option("missingkey=default")
	} else if !t.options[AcceptNoValue] {
		// To help detect errors on second run, we enable the option to raise error on nil values
		newTemplate.Option("missingkey=error")
	}

	func() {
		// Here, we invoke the parser within a pseudo func because we cannot
		// call the parser without locking
		templateMutex.Lock()
		defer templateMutex.Unlock()
		newTemplate, err = newTemplate.Parse(th.Code)
	}()
	if err != nil {
		InternalLog.Debugf("%s(%d): Parsing error %v", th.Filename, th.Try, err)
		return th.Handler(err)
	}

	var out bytes.Buffer
	workingContext := t.context
	if cloneContext {
		workingContext = collections.AsDictionary(workingContext).Clone()
	}
	if err = newTemplate.Execute(&out, workingContext); err != nil {
		InternalLog.Debugf("%s(%d): Execution error %v", th.Filename, th.Try, err)
		return th.Handler(err)
	}
	result = revertReplacements(t.substitute(out.String()))

	changed = result != originalContent

	if topCall && !t.options[AcceptNoValue] {
		// Detect possible <no value> or <nil> that could be generated
		if pos := strings.Index(strings.Replace(result, nilValue, noValue, -1), noValue); pos >= 0 {
			line := len(strings.Split(result[:pos], "\n"))
			return th.Handler(fmt.Errorf("template: %s:%d: %s", th.Filename, line, noValueError))
		}
	}

	if !t.options[AcceptNoValue] {
		// We restore the existing no value if any
		result = strings.Replace(result, noValueRepl, noValue, -1)
		result = strings.Replace(result, nilValueRepl, nilValue, -1)
	}
	return
}
