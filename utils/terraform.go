package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coveooss/gotemplate/v3/errors"
)

// IsTerraformFile check if the file extension matches on of the terraform file extension
func IsTerraformFile(file string) bool {
	switch filepath.Ext(file) {
	case ".tf", ".tf.json", ".tfvars":
		return true
	default:
		return false
	}
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
			errs = append(errs, fmt.Errorf("%v: %s", errs, strings.TrimSpace(string(output))))
			continue
		}

		// Sometimes, terraform fmt needs to be executed twice to format its content properly...
		cmd = exec.Command("terraform", "fmt")
		cmd.Dir = folder
		errors.Must(cmd.Run())
	}
	return errs.AsError()
}
