package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/errors"
	"github.com/coveooss/gotemplate/v3/template"
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

func exclude(files []string, patterns []string) (result []string, err error) {
	if patterns = extend(patterns); len(patterns) == 0 {
		// There is no exclusion pattern, so we return the list of files as is
		result = files
		return
	}
	result = make([]string, 0, len(files))

	for i, file := range files {
		file, err = filepath.Abs(file)
		if err != nil {
			return
		}
		var excluded bool
		for _, pattern := range patterns {
			pattern, err = filepath.Abs(pattern)
			if err != nil {
				return
			}
			if excluded, err = doublestar.Match(pattern, file); err != nil {
				return
			} else if excluded {
				template.Log.Debugf("%s ignored", files[i])
				break
			}
		}
		if !excluded {
			result = append(result, files[i])
		}
	}
	return
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

var iif = collections.IIf
