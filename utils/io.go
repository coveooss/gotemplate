package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coveo/gotemplate/errors"
)

// FindFiles returns the list of the files matching the array of patterns
func FindFiles(folder string, recursive, followLinks bool, patterns ...string) ([]string, error) {
	visited := map[string]bool{}
	var walker func(folder string) ([]string, error)
	walker = func(folder string) ([]string, error) {
		var results []string
		folder, _ = filepath.Abs(folder)

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if !recursive && folder != path {
					return filepath.SkipDir
				}
				visited[path] = true
				files, err := findFiles(path, patterns...)
				if err != nil {
					return err
				}
				results = append(results, files...)
				return nil
			}

			if info.Mode()&os.ModeSymlink != 0 && followLinks {
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}

				if !filepath.IsAbs(link) {
					link = filepath.Join(filepath.Dir(path), link)
				}
				link, _ = filepath.Abs(link)
				if !visited[link] {
					// Check if we already visited that link to avoid recursive loop
					linkFiles, err := walker(link)
					if err != nil {
						return err
					}
					results = append(results, linkFiles...)
				}
			}
			return nil
		})
		return results, nil
	}
	return walker(folder)
}

// FindFiles returns the list of files in the specified folder that match one of the supplied patterns
func findFiles(folder string, patterns ...string) ([]string, error) {
	var tfFiles []string
	for _, ext := range patterns {
		files, err := filepath.Glob(filepath.Join(folder, ext))
		if err != nil {
			return nil, err
		}
		tfFiles = append(tfFiles, files...)
	}
	return tfFiles, nil
}

// MustFindFiles returns the list of the files matching the array of patterns with panic on error
func MustFindFiles(folder string, recursive, followLinks bool, patterns ...string) []string {
	return errors.Must(FindFiles(folder, recursive, followLinks, patterns...)).([]string)
}

// GlobFunc returns an array of string representing the expansion of the supplied arguments using filepath.Glob function
func GlobFunc(args ...interface{}) (result []string) {
	for _, arg := range toStrings(args) {
		if strings.ContainsAny(arg, "*?[]") {
			if expanded, _ := filepath.Glob(arg); expanded != nil {
				result = append(result, expanded...)
				continue
			}
		}
		result = append(result, arg)
	}
	return
}

// Pwd returns the current folder
func Pwd() string {
	return errors.Must(os.Getwd()).(string)
}

// Relative returns the relative path of file from folder
func Relative(folder, file string) string {
	if !filepath.IsAbs(file) {
		return file
	}
	return errors.Must(filepath.Rel(folder, file)).(string)
}

// https://regex101.com/r/ykVKPt/3
var shebang = regexp.MustCompile(`(?sm)^#!\s*(?P<program>[^\s]*)\s*(?P<app>[^\s]*)?\s*$\s*(?P<source>.*)`)

// IsShebangScript determines if the supplied code has a Shebang definition #! program subprogram
func IsShebangScript(content string) bool {
	return shebang.MatchString(strings.TrimSpace(content))
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
