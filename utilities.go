package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/gotemplate/template"
	"github.com/coveo/gotemplate/utils"
	goerrors "github.com/go-errors/errors"
)

var must = errors.Must

func cleanup() {
	os.RemoveAll(tempFolder)

	if err := recover(); err != nil {
		switch err := err.(type) {
		case errors.Managed:
			errors.Print(err)
		case *goerrors.Error:
			errors.Printf(err.ErrorStack())
		default:
			errors.Printf("%T: %[1]v\n%s", err, string(debug.Stack()))
		}
	}
}

func readStdin() string {
	if stdinContent != "" {
		return stdinContent
	}
	stdinContent = string(must(ioutil.ReadAll(os.Stdin)).([]byte))
	return stdinContent
}

var stdinContent string

func exclude(files []string, patterns []string) []string {
	if patterns = extend(patterns); len(patterns) == 0 {
		// There is no exclusion pattern, so we return the list of files as is
		return files
	}

	result := make([]string, 0, len(files))
	for i, file := range files {
		var excluded bool
		for _, pattern := range patterns {
			file = iif(strings.ContainsAny(pattern, `/\`), file, filepath.Base(file)).(string)
			if excluded = must(filepath.Match(pattern, file)).(bool); excluded {
				template.Log.Noticef("%s ignored", files[i])
				break
			}
		}
		if !excluded {
			result = append(result, files[i])
		}
	}
	return result
}

func extend(values []string) []string {
	result := make([]string, 0, len(values))
	for i := range values {
		for _, sv := range strings.Split(values[i], ",") {
			sv = strings.TrimSpace(sv)
			if sv != "" {
				result = append(result, sv)
			}
		}
	}
	return result
}

var iif = utils.IIf
