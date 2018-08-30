package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// https://regex101.com/r/ykVKPt/6
var shebang = regexp.MustCompile(`(?sm)^(?:\s*#!\s*(?P<program>[^\s]*)[ \t]*(?P<app>[^\s]*)?\s*$)?\s*(?P<source>.*)`)

// IsShebangScript determines if the supplied code has a Shebang definition #! program subprogram
func IsShebangScript(content string) bool {
	matches := shebang.FindStringSubmatch(strings.TrimSpace(content))
	return len(matches) > 0 && matches[1] != ""
}

// ScriptParts splits up the supplied content into program, subprogram and source if the content matches Shebang defintion
func ScriptParts(content string) (program, subprogram, source string) {
	matches := shebang.FindStringSubmatch(strings.TrimSpace(content))
	if len(matches) > 0 {
		program = matches[1]
		subprogram = matches[2]
		source = matches[3]
	} else {
		source = content
	}
	return
}

// GetCommandFromFile returns an exec.Cmd structure to run the supplied script file
func GetCommandFromFile(filename string, args ...interface{}) (cmd *exec.Cmd, err error) {
	script, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	executer, delegate, command := ScriptParts(strings.TrimSpace(string(script)))

	strArgs := GlobFunc(args...)

	if executer != "" {
		command = executer
		strArgs = append([]string{filename}, strArgs...)
		if delegate != "" {
			strArgs = append([]string{delegate}, strArgs...)
		}
	} else if _, errPath := exec.LookPath(command); errPath != nil {
		if strings.Contains(command, " ") {
			// The command is a string that should be splitted up into several parts
			split := strings.Split(command, " ")
			command = split[0]
			strArgs = append(split[1:], strArgs...)
		} else {
			// The command does not exist
			return
		}
	}

	cmd = exec.Command(command, strArgs...)
	return
}

// GetCommandFromString returns an exec.Cmd structure to run the supplied command
func GetCommandFromString(script string, args ...interface{}) (cmd *exec.Cmd, tempFile string, err error) {
	if executer, delegate, command := ScriptParts(strings.TrimSpace(script)); executer != "" {
		cmd, tempFile, err = saveTempFile(fmt.Sprintf("#! %s %s\n%s", executer, delegate, command), args...)
	} else {
		strArgs := GlobFunc(args...)
		if _, err = exec.LookPath(command); err == nil {
			cmd = exec.Command(command, strArgs...)
			return
		}
		if IsCommand(command) {
			// The command does just not exist, we return the error
			return
		}

		if cmd, tempFile, err = saveTempFile(fmt.Sprintf("#! %s\n%s", getDefaultShell(), command), args...); err != nil {
			return
		}
	}
	return
}

// IsCommand ensures that the supplied command does not contain any shell specific characters
func IsCommand(command string) bool {
	return !strings.ContainsAny(command, " \t|&$,;(){}<>[]")
}

func saveTempFile(content string, args ...interface{}) (cmd *exec.Cmd, fileName string, err error) {
	var temp *os.File
	if temp, err = ioutil.TempFile("", "exec_"); err != nil {
		return
	}

	if _, err = temp.WriteString(content); err != nil {
		return
	}
	temp.Close()
	fileName = temp.Name()
	cmd, err = GetCommandFromFile(fileName, args...)
	return
}

func getDefaultShell() string {
	if defaultShell != "" {
		return defaultShell
	}
	defShells := []string{"bash", "sh", "zsh", "ksh"}
	if runtime.GOOS == "windows" {
		defShells = []string{"powershell", "pwsh", "cmd"}
	}
	for _, sh := range defShells {
		if sh, err := exec.LookPath(sh); err == nil {
			defaultShell = sh
			break
		}
	}
	return defaultShell
}

var defaultShell string

// GetEnv returns the environment variable value or defaultValue if it does not exist.
func GetEnv(varName, defaultValue string) string {
	if value, exist := os.LookupEnv(varName); exist {
		return value
	}
	return defaultValue
}
