package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/coveo/terragrunt/util"
)

// IsTerraformFile check if the file extension matches on of the terraform file extension
func IsTerraformFile(file string) bool {
	return util.ListContainsElement([]string{".tf", ".tf.json", ".tfvars"}, filepath.Ext(file))
}

// TerraformFormat applies terraform fmt on
func TerraformFormat(files ...string) error {
	if _, err := exec.LookPath("terraform"); err != nil {
		return nil
	}

	tfFolders := make(map[string]bool)
	for _, file := range files {
		if IsTerraformFile(file) {
			tfFolders[filepath.Dir(file)] = true
		}
	}

	var errs errors.Array
	for folder := range tfFolders {
		cmd := exec.Command("terraform", "fmt")
		cmd.Dir = folder
		if output, err := cmd.CombinedOutput(); err != nil {
			errs = append(errs, fmt.Errorf("%v: %s", err, strings.TrimSpace(string(output))))
			continue
		}

		// Sometimes, terraform fmt needs to be executed twice to format its content properly...
		cmd = exec.Command("terraform", "fmt")
		cmd.Dir = folder
		errors.Must(cmd.Run())
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
