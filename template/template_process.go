package template

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/collections"
	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/utils"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	templateExt    = []string{".gt", ".template"}
	linePrefix     = `^template: ` + p(tagLocation, p(tagFile, `.*?`)+`:`+p(tagLine, `\d+`)+`(:`+p(tagCol, `\d+`)+`)?: `)
	execPrefix     = linePrefix + `executing ".*" at <` + p(tagCode, `.*`) + `>: `
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
func (t Template) ProcessContent(content, source string) (string, error) {
	return t.processContentInternal(content, source, nil, 0, true)
}

func (t Template) processContentInternal(originalContent, source string, originalSourceLines []string, retryCount int, cloneContext bool) (result string, err error) {
	topCall := originalSourceLines == nil
	content := originalContent
	if topCall {
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

		if !t.options[AcceptNoValue] {
			// We replace any pre-existing no value to avoid false error detection
			content = strings.Replace(content, noValue, noValueRepl, -1)
			content = strings.Replace(content, nilValue, nilValueRepl, -1)
		}

		defer func() {
			// If we get errors and the file is not an explicit gotemplate file, or contains
			// gotemplate! or the strict error check mode is not enabled, we simply
			// add a trace with the error content and return the content unaltered
			if err != nil {
				strictMode := t.options[StrictErrorCheck]
				strictMode = strictMode || strings.Contains(originalContent, explicitGoTemplate)
				extension := filepath.Ext(source)
				strictMode = strictMode || (extension != "" && strings.Contains(".gt,.gte,.template", extension))
				if !(strictMode) {
					Log.Noticef("Ignored gotemplate error in %s (file left unchanged):\n%s", color.CyanString(source), err.Error())
					result, err = originalContent, nil
				}
			}
		}()
	}

	// This local functions handle all errors from Parse or Execute and tries to fix the template to allow discovering
	// of all errors in a template instead of stopping after the first one encountered
	handleError := func(err error) (string, error) {
		if originalSourceLines == nil {
			originalSourceLines = strings.Split(originalContent, "\n")
		}

		regexGroup := must(utils.GetRegexGroup("Parse", templateErrors)).([]*regexp.Regexp)

		if matches, _ := utils.MultiMatch(err.Error(), regexGroup...); len(matches) > 0 {
			// We remove the faulty line and continue the processing to get all errors at once
			lines := strings.Split(content, "\n")
			faultyLine := toInt(matches[tagLine]) - 1
			faultyColumn := 0
			key, message, errText, code := matches[tagKey], matches[tagMsg], matches[tagErr], matches[tagCode]

			if matches[tagCol] != "" {
				faultyColumn = toInt(matches[tagCol]) - 1
			}

			errorText, parserBug := color.RedString(errText), ""

			if faultyLine >= len(lines) {
				faultyLine = len(lines) - 1
				// TODO: This code can be removed once issue has been fixed
				parserBug = color.HiRedString("\nBad error line reported: check: https://github.com/golang/go/issues/27319")
			}

			currentLine := String(lines[faultyLine])

			if matches[tagFile] != source {
				// An error occurred in an included external template file, we cannot try to recuperate
				// and try to find further errors, so we just return the error.

				if fileContent, err := ioutil.ReadFile(matches[tagFile]); err != nil {
					currentLine = String(fmt.Sprintf("Unable to read file: %v", err))
				} else {
					currentLine = String(fileContent).Lines()[toInt(matches[tagLine])-1]
				}
				return "", fmt.Errorf("%s %v in: %s", color.WhiteString(source), err, color.HiBlackString(currentLine.Str()))
			}
			if faultyColumn != 0 && strings.Contains(" (", currentLine[faultyColumn:faultyColumn+1].Str()) {
				// Sometime, the error is not reporting the exact column, we move 1 char forward to get the real problem
				faultyColumn++
			}

			errorLine := fmt.Sprintf(" in: %s", color.HiBlackString(originalSourceLines[faultyLine]))
			var logMessage string
			if key != "" {
				// Missing key and we disabled the <no value> mode
				context := String(currentLine).SelectContext(faultyColumn, t.LeftDelim(), t.RightDelim())
				if subContext := String(currentLine).SelectContext(faultyColumn, "(", ")"); subContext != "" {
					// There is an sub-context, so we replace it first
					context = subContext
				}
				current := String(currentLine).SelectWord(faultyColumn, ".")
				newContext := context.Replace(current.Str(), undefError).Str()
				newLine := currentLine.Replace(context.Str(), newContext)

				left := fmt.Sprintf(`(?P<begin>(%s-?\s*(if|range|with)\s.*|\()\s*)?`, regexp.QuoteMeta(t.LeftDelim()))
				right := fmt.Sprintf(`(?P<end>\s*(-?%s|\)))`, regexp.QuoteMeta(t.RightDelim()))
				const (
					ifUndef = "ifUndef"
					isZero  = "isZero"
					assert  = "assert"
				)
				undefRegexDefintions := []string{
					fmt.Sprintf(`%[1]s(undef|ifUndef|default)\s+(?P<%[3]s>.*?)\s+%[4]s%[2]s`, left, right, ifUndef, undefError),
					fmt.Sprintf(`%[1]s(?P<%[3]s>%[3]s|isNil|isNull|isEmpty|isSet)\s+%[4]s%[2]s`, left, right, isZero, undefError),
					fmt.Sprintf(`%[1]s%[3]s\s+(?P<%[3]s>%[4]s).*?%[2]s`, left, right, assert, undefError),
				}
				expressions, errRegex := utils.GetRegexGroup(fmt.Sprintf("Undef%s", t.delimiters), undefRegexDefintions)
				if errRegex != nil {
					log.Error(errRegex)
				}
				undefMatches, n := utils.MultiMatch(newContext, expressions...)

				if undefMatches[ifUndef] != "" {
					logMessage = fmt.Sprintf("Managed undefined value %s: %s", key, context)
					err = nil
					lines[faultyLine] = newLine.Replace(newContext, expressions[n].ReplaceAllString(newContext, fmt.Sprintf("${begin}${%s}${end}", ifUndef))).Str()
				} else if undefMatches[isZero] != "" {
					logMessage = fmt.Sprintf("Managed undefined value %s: %s", key, context)
					err = nil
					value := fmt.Sprintf("%s%v%s", undefMatches["begin"], undefMatches[isZero] != "isSet", undefMatches["end"])
					lines[faultyLine] = newLine.Replace(newContext, value).Str()
				} else if undefMatches[assert] != "" {
					logMessage = fmt.Sprintf("Managed assertion on %s: %s", key, context)
					err = nil
					lines[faultyLine] = newLine.Replace(newContext, strings.Replace(newContext, undefError, "0", 1)).Str()
				} else {
					logMessage = fmt.Sprintf("Unmanaged undefined value %s: %s", key, context)
					errorText = color.RedString("Undefined value ") + color.YellowString(key)
					lines[faultyLine] = newLine.Str()
				}
			} else if message != "" {
				logMessage = fmt.Sprintf("User defined error: %s", message)
				errorText = color.RedString(message)
				lines[faultyLine] = fmt.Sprintf("ERROR %s", errText)
			} else if code != "" {
				logMessage = fmt.Sprintf("Execution error: %s", err)
				context := String(currentLine).SelectContext(faultyColumn, t.LeftDelim(), t.RightDelim())
				errorText = fmt.Sprintf(color.RedString("%s (%s)", errText, code))
				if context == "" {
					// We have not been able to find the current context, we wipe the erroneous line
					lines[faultyLine] = fmt.Sprintf("ERROR %s", errText)
				} else {
					lines[faultyLine] = currentLine.Replace(context.Str(), runError).Str()
				}
			} else if errText != noValueError {
				logMessage = fmt.Sprintf("Parsing error: %s", err)
				lines[faultyLine] = fmt.Sprintf("ERROR %s", errText)
			}
			if currentLine.Contains(runError) || strings.Contains(code, undefError) {
				// The erroneous line has already been replaced, we do not report the error again
				err, errorText = nil, ""
				log.Debugf("Ignored error %s", logMessage)
			} else if logMessage != "" {
				log.Debug(logMessage)
			}

			if err != nil {
				err = fmt.Errorf("%s%s%s%s", color.WhiteString(matches[tagLocation]), errorText, errorLine, parserBug)
			}
			if lines[faultyLine] != currentLine.Str() || strings.Contains(err.Error(), noValueError) {
				// If we changed something in the current text, we try to continue the evaluation to get further errors
				result, err2 := t.processContentInternal(strings.Join(lines, "\n"), source, originalSourceLines, retryCount+1, false)
				if err2 != nil {
					if err != nil && errText != noValueError {
						if err.Error() == err2.Error() {
							// TODO See: https://github.com/golang/go/issues/27319
							err = fmt.Errorf("%v\n%s", err, color.HiRedString("Unable to continue processing to check for further errors"))
						} else {
							err = fmt.Errorf("%v\n%v", err, err2)
						}
					} else {
						err = err2
					}
				}
				return result, err
			}
		}
		return "", err
	}

	context := t.GetNewContext(filepath.Dir(source), true)
	newTemplate := context.New(source)

	if topCall {
		newTemplate.Option("missingkey=default")
	} else if !t.options[AcceptNoValue] {
		// To help detect errors on second run, we enable the option to raise error on nil values
		// log.Infof("%s(%d): Activating the missing key error option", source, retryCount)
		newTemplate.Option("missingkey=error")
	}

	func() {
		// Here, we invoke the parser within a pseudo func because we cannot
		// call the parser without locking
		templateMutex.Lock()
		defer templateMutex.Unlock()
		newTemplate, err = newTemplate.Parse(content)
	}()
	if err != nil {
		log.Infof("%s(%d): Parsing error %v", source, retryCount, err)
		return handleError(err)
	}

	var out bytes.Buffer
	workingContext := t.context
	if cloneContext {
		workingContext = collections.AsDictionary(workingContext).Clone()
	}
	if err = newTemplate.Execute(&out, workingContext); err != nil {
		log.Infof("%s(%d): Execution error %v", source, retryCount, err)
		return handleError(err)
	}
	result = t.substitute(out.String())

	if topCall && !t.options[AcceptNoValue] {
		// Detect possible <no value> or <nil> that could be generated
		if pos := strings.Index(strings.Replace(result, nilValue, noValue, -1), noValue); pos >= 0 {
			line := len(strings.Split(result[:pos], "\n"))
			return handleError(fmt.Errorf("template: %s:%d: %s", source, line, noValueError))
		}
	}

	if !t.options[AcceptNoValue] {
		// We restore the existing no value if any
		result = strings.Replace(result, noValueRepl, noValue, -1)
		result = strings.Replace(result, nilValueRepl, nilValue, -1)
	}
	return
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
		// This occurs when gotemplate code has been supplied as a filename. In that case, we simply render
		// the result to the stdout
		Println(result)
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

	mode := must(os.Stat(template)).(os.FileInfo).Mode()
	if !isTemplate && !t.options[Overwrite] {
		newName := template + ".original"
		log.Noticef("%s => %s", utils.Relative(t.folder, template), utils.Relative(t.folder, newName))
		must(os.Rename(template, template+".original"))
	}

	if sourceFolder != targetFolder {
		must(os.MkdirAll(filepath.Dir(resultFile), 0777))
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
		t.options[OutputStdout] = print // Some file may change this option at runtime, so we restore it back to its originalSourceLines value between each file
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
