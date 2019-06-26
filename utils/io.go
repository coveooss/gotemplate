package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/errors"
)

var must = errors.Must

// FindFiles returns the list of the files matching the array of patterns
func FindFiles(folder string, recursive, followLinks bool, patterns ...string) ([]string, error) {
	depth := 0
	if recursive {
		depth = 1 << 16
	}
	return FindFilesMaxDepth(folder, depth, followLinks, patterns...)
}

// FindFilesMaxDepth returns the list of the files matching the array of patterns
func FindFilesMaxDepth(folder string, maxDepth int, followLinks bool, patterns ...string) ([]string, error) {
	visited := map[string]bool{}
	var walker func(folder string) ([]string, error)
	walker = func(folder string) ([]string, error) {
		results := must(findFiles(folder, patterns...)).([]string)
		folder, _ = filepath.Abs(folder)
		if maxDepth == 0 {
			return results, nil
		}

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info == nil || path == folder {
				return nil
			}
			if info.IsDir() {
				visited[path] = true
				depth := strings.Count(must(filepath.Rel(path, folder)).(string), "..")
				if depth > maxDepth {
					return filepath.SkipDir
				}
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
	var matches []string
	for _, pattern := range patterns {
		files, err := filepath.Glob(filepath.Join(folder, pattern))
		if err != nil {
			return nil, err
		}
		matches = append(matches, files...)
	}
	return matches, nil
}

// MustFindFiles returns the list of the files matching the array of patterns with panic on error
func MustFindFiles(folder string, recursive, followLinks bool, patterns ...string) []string {
	return must(FindFiles(folder, recursive, followLinks, patterns...)).([]string)
}

// MustFindFilesMaxDepth returns the list of the files matching the array of patterns with panic on error
func MustFindFilesMaxDepth(folder string, maxDepth int, followLinks bool, patterns ...string) []string {
	return must(FindFilesMaxDepth(folder, maxDepth, followLinks, patterns...)).([]string)
}

func globFunc(trimUnmatch bool, args ...interface{}) (result []string) {
	for _, arg := range collections.ToStrings(args) {
		if strings.ContainsAny(arg, "*?[]") {
			if expanded, _ := filepath.Glob(arg); expanded != nil {
				result = append(result, expanded...)
				continue
			}
			if trimUnmatch {
				continue
			}
		}
		result = append(result, arg)
	}
	return
}

// GlobFunc returns an array of string representing the expansion of the supplied arguments using filepath.Glob function
func GlobFunc(args ...interface{}) []string { return globFunc(false, args...) }

// GlobFuncTrim returns an array of string representing the expansion of the supplied arguments using filepath.Glob function, it removes the unmatched arguments
func GlobFuncTrim(args ...interface{}) []string { return globFunc(true, args...) }

// Pwd returns the current folder
func Pwd() string {
	return must(os.Getwd()).(string)
}

// Relative returns the relative path of file from folder
func Relative(folder, file string) string {
	if !filepath.IsAbs(file) {
		return file
	}
	return must(filepath.Rel(folder, file)).(string)
}
