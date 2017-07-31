package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/coveo/terragrunt/util"
)

var basicFunctions = []string{
	"and", "call", "call", "html", "index", "js", "len", "not", "or", "print", "printf", "println", "urlquery",
	"eq", "ge", "gt", "le", "lt", "ne",
}

const templateExt = ".template"

func processFiles(context interface{}) {
	// Create the evaluation context
	mainTemplate := template.New("context").Funcs(sprig.GenericFuncMap())
	AddExtraFuncs(mainTemplate)

	*includePatterns = append(*includePatterns, "*"+templateExt)

	if *listFunc {
		printFunctions(mainTemplate)
		os.Exit(0)
	}

	pwd := PanicOnError(os.Getwd()).(string)
	var files []string
	if *recursive {
		files = PanicOnError(findFiles(pwd, *includePatterns...)).([]string)
	} else {
		files = PanicOnError(util.FindFiles(pwd, *includePatterns...)).([]string)
	}

	if len(files) == 0 {
		// There is nothing to process
		return
	}

	// Parse the template files
	templateFiles, files := isolateTemplateFiles(files)
	if len(templateFiles) > 0 {
		mainTemplate = PanicOnError(mainTemplate.ParseFiles(templateFiles...)).(*template.Template)
	}

	// Process the files and generate resulting file if there is an output from the template
	for _, file := range templateFiles {
		if *dryRun {
			fmt.Println("Processing file", file)
		}
		template := mainTemplate.Lookup(filepath.Base(file))
		out, err := executeTemplate(template, file, context)
		if err != nil {
			if *dryRun {
				fmt.Fprintln(os.Stderr, err, "while executing", file)
				continue
			} else {
				PanicOnError(err)
			}
		}

		if result := strings.TrimSpace(out.String()); result != "" {
			fileName := strings.TrimSuffix(file, templateExt)
			ext := path.Ext(fileName)
			fileName = fmt.Sprint(strings.TrimSuffix(fileName, ext), ".generated", ext)
			if *dryRun {
				fmt.Printf("  Would create file %s with content:\n%s\n", fileName, result)
			} else {
				info, _ := os.Stat(file)
				if err := ioutil.WriteFile(fileName, []byte(result), info.Mode()); err != nil {
					PanicOnError(err)
				}
			}
		}
	}

	for _, file := range files {
		if *dryRun {
			fmt.Println("Processing file", file)
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while reading", file)
			continue
		}
		template, err := mainTemplate.Parse(string(content))
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while parsing", file)
			continue
		}
		out, err := executeTemplate(template, file, context)
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "while executing", file)
			continue
		}

		if out.String() != string(content) {
			if *dryRun {
				fmt.Printf("  Would modify file %s with content:\n%s\n", file, out.String())
			} else {
				info, _ := os.Stat(file)
				if err := os.Rename(file, file+".original"); err != nil {
					PanicOnError(err)
				}
				ioutil.WriteFile(file, out.Bytes(), info.Mode())
			}
			fmt.Println()
		}
	}

	// If terraform is present, we apply a terraform fmt to the resulting templates
	// to ensure that the resulting files are valids.
	if _, err := exec.LookPath("terraform"); err == nil {
		cmd := exec.Command("terraform", "fmt")
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// RunningTemplate represents the current running template (this avoid concurrent processing)
var RunningTemplate *template.Template

func executeTemplate(template *template.Template, fileName string, context interface{}) (bytes.Buffer, error) {
	var out bytes.Buffer

	RunningTemplate = template
	RunningTemplate.ParseName = fileName

	err := RunningTemplate.Execute(&out, context)
	return out, err
}

func findFiles(folder string, patterns ...string) ([]string, error) {
	visited := map[string]bool{}
	var walker func(folder string) ([]string, error)
	walker = func(folder string) ([]string, error) {
		var results []string
		folder, _ = filepath.Abs(folder)

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				visited[path] = true
				files, err := util.FindFiles(path, *includePatterns...)
				if err != nil {
					return err
				}
				results = append(results, files...)
				return nil
			}

			if info.Mode()&os.ModeSymlink != 0 && *followSymLinks {
				link, err := os.Readlink(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, err, "while trying to follow link", path)
					return nil
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

func isolateTemplateFiles(files []string) (templates, others []string) {
	for _, file := range files {
		if filepath.Ext(file) == templateExt {
			templates = append(templates, file)
		} else {
			others = append(others, file)
		}
	}
	return
}

func printFunctions(template *template.Template) {
	common := reflect.ValueOf(template).FieldByName("common").Elem().FieldByName("parseFuncs")
	keys := common.MapKeys()
	functions := make([]string, len(keys))
	for i, k := range keys {
		functions[i] = k.String()
	}

	functions = append(basicFunctions, functions...)
	sort.Strings(functions)
	for i, function := range functions {
		fmt.Printf("%-15s", function)
		if (i+1)%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
