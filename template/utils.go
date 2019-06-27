package template

import (
	"path/filepath"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/errors"
	"github.com/coveooss/gotemplate/v3/utils"
)

var must = errors.Must
var trapError = errors.Trap
var iif = collections.IIf
var ifUndef = collections.IfUndef
var defval = collections.Default

type (
	iList       = collections.IGenericList
	iDictionary = collections.IDictionary
)

func getTargetFile(targetFile, sourcePath, targetPath string) string {
	if targetPath != "" {
		targetFile = filepath.Join(targetPath, utils.Relative(sourcePath, targetFile))
	}
	return targetFile
}
