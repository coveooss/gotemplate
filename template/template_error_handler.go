package template

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/coveooss/multilogger/reutils"
	"github.com/fatih/color"
)

// TemplateErrorHandler handle errors occuring during template evaluation and try to mitigate the error
// and continuing evaluation in order to return all potential errors instead of stopping after the first one
type errorHandler struct {
	*Template
	Filename string   // The filename associated with the current evaluation
	Source   string   // Original source code
	Code     string   // Modified code
	Lines    []string // Original source code as an array of string (one per line)
	Try      int      // The current evaluation try
}

func (t errorHandler) Handler(err error) (string, bool, error) {
	if t.Lines == nil {
		t.Lines = strings.Split(t.Source, "\n")
	}

	// We try to identify the exact position of the error. If the error occurred in a sub template, the
	// error will be the last one found in the error string
	errorsPosition := reError.FindAllStringIndex(err.Error(), -1)
	lastErrorBegin := 0
	if errorsPosition != nil {
		lastErrorBegin = errorsPosition[len(errorsPosition)-1][0]
	}

	if matches, _ := reutils.MultiMatch(err.Error()[lastErrorBegin:], errorParsingExpressions...); len(matches) > 0 {
		// We remove the faulty line and continue the processing to get all errors at once
		lines := strings.Split(t.Code, "\n")
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

		if matches[tagFile] != t.Filename {
			// An error occurred in an included external template file, we cannot try to recuperate
			// and try to find further errors, so we just return the error.

			if fileContent, err := ioutil.ReadFile(matches[tagFile]); err != nil {
				currentLine = String(fmt.Sprintf("Unable to read file: %v", err))
			} else {
				currentLine = String(fileContent).Lines()[toInt(matches[tagLine])-1]
			}
			return "", true, fmt.Errorf("%s %v in: %s", color.WhiteString(t.Filename), err, color.HiBlackString(currentLine.Str()))
		}
		if faultyColumn != 0 && strings.Contains(" (", currentLine[faultyColumn:faultyColumn+1].Str()) {
			// Sometime, the error is not reporting the exact column, we move 1 char forward to get the real problem
			faultyColumn++
		}

		errorLine := fmt.Sprintf(" in: %s", color.HiBlackString(t.Lines[faultyLine]))
		var logMessage string
		if key != "" {
			// Missing key and we disabled the <no value> mode
			context := String(currentLine).SelectContext(faultyColumn, t.LeftDelim(), t.RightDelim())
			if subContext := String(currentLine).SelectContext(faultyColumn, "(", ")"); subContext != "" {
				// There is an sub-context, so we replace it first
				context = subContext
			}
			current := String(currentLine).SelectWord(faultyColumn, ".", "_")
			newContext := context.Replace(current.Str(), undefError).Str()
			newLine := currentLine.Replace(context.Str(), newContext)

			left := fmt.Sprintf(`(?P<begin>(%s-?\s*(if|range|with)\s.*|\()\s*)?`, regexp.QuoteMeta(t.LeftDelim()))
			right := fmt.Sprintf(`(?P<end>\s*(-?%s|\)))`, regexp.QuoteMeta(t.RightDelim()))
			const (
				ifUndef = "ifUndef"
				isZero  = "isZero"
				assert  = "assert"
			)

			groupName := fmt.Sprintf("Undef%s", t.delimiters)
			expressions := reutils.GetRegexGroup(groupName)
			if expressions == nil {
				var errRegex error
				undefRegexDefinitions := []string{
					fmt.Sprintf(`%[1]s(undef|ifUndef|default)\s+(?P<%[3]s>.*?)\s+%[4]s%[2]s`, left, right, ifUndef, undefError),
					fmt.Sprintf(`%[1]s(?P<%[3]s>%[3]s|isNil|isNull|isEmpty|isSet)\s+%[4]s%[2]s`, left, right, isZero, undefError),
					fmt.Sprintf(`%[1]s%[3]s\s+(?P<%[3]s>%[4]s).*?%[2]s`, left, right, assert, undefError),
				}
				if expressions, errRegex = reutils.NewRegexGroup(fmt.Sprintf("Undef%s", t.delimiters), undefRegexDefinitions...); errRegex != nil {
					InternalLog.Error(errRegex)
				}
			}
			undefMatches, n := reutils.MultiMatch(newContext, expressions...)

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
			InternalLog.Tracef("Ignored error %s", logMessage)
		} else if logMessage != "" {
			InternalLog.Trace(logMessage)
		}

		if err != nil {
			err = fmt.Errorf("%s%s%s%s", color.WhiteString(matches[tagLocation]), errorText, errorLine, parserBug)
		}
		if lines[faultyLine] != currentLine.Str() || strings.Contains(err.Error(), noValueError) {
			// If we changed something in the current text, we try to continue the evaluation to get further errors
			newCode := strings.Join(lines, "\n")
			if err != nil {
				InternalLog.Infof("Retrying %d with:\n%s", t.Try, color.HiBlackString(String(newCode).AddLineNumber(0).Str()))
			}
			result, changed, err2 := t.processContentInternal(newCode, t.Filename, t.Lines, t.Try+1, false, nil)
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
			return result, changed, err
		}
	}
	return "", true, err
}

var (
	errorParsingExpressions []*regexp.Regexp
	linePrefix              = `template: ` + p(tagLocation, p(tagFile, `.*?`)+`:`+p(tagLine, `\d+`)+`(:`+p(tagCol, `\d+`)+`)?: `)
	reError                 = regexp.MustCompile(linePrefix)
)

func init() {
	execPrefix := "(?s)^" + linePrefix + `executing ".*" at <` + p(tagCode, `.*`) + `>: `
	templateErrors := []string{
		execPrefix + `map has no entry for key "` + p(tagKey, `.*`) + `"`,
		execPrefix + `(?s)error calling (raise|assert): ` + p(tagMsg, `.*`),
		execPrefix + p(tagErr, `.*`),
		linePrefix + p(tagErr, `.*`),
	}

	var err error
	if errorParsingExpressions, err = reutils.NewRegexGroup("errorParsingExpressions", templateErrors...); err != nil {
		panic(err)
	}
}
