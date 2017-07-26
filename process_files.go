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
)

func processFiles(context interface{}) {
	// Create the evaluation context
	root := template.New("context").Funcs(sprig.GenericFuncMap())
	AddExtraFuncs(root)

	if *listFunc {
		val := reflect.ValueOf(*root)
		common := val.FieldByName("common").Elem().FieldByName("funcMaps")
		keys := common.MapKeys()
		functions := make([]string, len(keys))
		for i, k := range keys {
			functions[i] = k.String()
		}
		sort.Strings(functions)
		fmt.Println(strings.Join(functions, "\n"))
		os.Exit(0)
	}

	var files []string
	if files = Must(filepath.Glob("*.template")).([]string); len(files) == 0 {
		// There is nothing to process
		return
	}

	// Parse the files
	templates := Must(root.ParseFiles(files...)).(*template.Template)

	// Process the files and generate resulting file if there is an output from the template
	for _, file := range files {
		template := templates.Lookup(file)
		var out bytes.Buffer
		Must(template.Execute(&out, context))

		if result := strings.TrimSpace(out.String()); result != "" {
			fileName := strings.TrimSuffix(file, path.Ext(file))
			ioutil.WriteFile(fileName, []byte(result), 0644)
		}
	}

	// If terraform is present, we apply a terraform fmt to the resulting templates
	if _, err := exec.LookPath("terraform"); err == nil {
		cmd := exec.Command("terraform", "fmt")
		cmd.Stderr = os.Stderr
		Must(cmd.Run())
	}
}
