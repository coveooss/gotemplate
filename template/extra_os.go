package template

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"time"

	"github.com/coveooss/gotemplate/v3/collections"
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	osBase = "Operating systems functions"
)

var osFuncs = dictionary{
	"diff":         diff,
	"exists":       fileExists,
	"glob":         glob,
	"group":        userGroup,
	"home":         userHome,
	"isDir":        isDir,
	"isExecutable": isExecutable,
	"isFile":       isFile,
	"isReadable":   isReadable,
	"isWriteable":  isWriteable,
	"joinPath":     path.Join,
	"lastMod":      lastMod,
	"load":         loadFromFile,
	"lookPath":     lookPath,
	"mode":         fileMode,
	"pwd":          utils.Pwd,
	"save":         saveToFile,
	"size":         fileSize,
	"stat":         os.Stat,
	"user":         user.Current,
	"username":     username,
}

var osFuncsArgs = arguments{
	"diff":         {"text1", "text2"},
	"exists":       {"filename"},
	"isDir":        {"filename"},
	"isExecutable": {"filename"},
	"isFile":       {"filename"},
	"isReadable":   {"filename"},
	"isWriteable":  {"filename"},
	"lastMod":      {"filename"},
	"load":         {"filename", "binary"},
	"mode":         {"filename"},
	"save":         {"filename", "object"},
	"size":         {"filename"},
}

var osFuncsAliases = aliases{
	"diff":     {"difference"},
	"exists":   {"fileExists", "isExist"},
	"glob":     {"expand"},
	"group":    {"userGroup"},
	"home":     {"homeDir", "homeFolder"},
	"isDir":    {"isDirectory", "isFolder"},
	"lastMod":  {"lastModification", "lastModificationTime"},
	"save":     {"write", "writeTo"},
	"lookPath": {"whereIs", "look", "which", "type"},
	"mode":     {"fileMode"},
	"pwd":      {"currentDir"},
	"load":     {"read", "readFrom"},
	"size":     {"fileSize"},
	"stat":     {"fileStat"},
	"user":     {"currentUser"},
}

var osFuncsHelp = descriptions{
	"diff":         "Returns a colored string that highlight differences between supplied texts.",
	"exists":       "Determines if a file exists or not.",
	"glob":         "Returns the expanded list of supplied arguments (expand *[]? on filename).",
	"group":        "Returns the current user group information (user.Group object).",
	"home":         "Returns the home directory of the current user.",
	"isDir":        "Determines if the file is a directory.",
	"isExecutable": "Determines if the file is executable by the current user.",
	"isFile":       "Determines if the file is a file (i.e. not a directory).",
	"isReadable":   "Determines if the file is readable by the current user.",
	"isWriteable":  "Determines if the file is writeable by the current user.",
	"joinPath":     "Joins any number of path elements into a single path, adding a separating slash if necessary. The result is Cleaned; in particular all empty strings are ignored.",
	"lastMod":      "Returns the last modification time of the file.",
	"load":         "Read object from file. If binary is true, the object is returned as a byte array.",
	"lookPath":     "Returns the location of the specified executable (returns empty string if not found).",
	"mode":         "Returns the file mode.",
	"pwd":          "Returns the current working directory.",
	"save":         "Save object to file. If the object is not an array of byte, then the object is converted to a string.",
	"size":         "Returns the file size.",
	"stat":         "Returns the file Stat information (os.Stat object).",
	"user":         "Returns the current user information (user.User object).",
	"username":     "Returns the current user name.",
}

func (t *Template) addOSFuncs() {
	t.AddFunctions(osFuncs, osBase, FuncOptions{
		FuncHelp:    osFuncsHelp,
		FuncArgs:    osFuncsArgs,
		FuncAliases: osFuncsAliases,
	})
}

func fileExists(file interface{}) (bool, error) {
	if _, err := os.Stat(fmt.Sprint(file)); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	return true, nil
}

func fileMode(file interface{}) (os.FileMode, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return 0, err
	}
	return stat.Mode(), nil
}

func fileSize(file interface{}) (int64, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func isDir(file interface{}) (bool, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

func isFile(file interface{}) (bool, error) {
	isDir, err := isDir(file)
	if err != nil {
		return false, err
	}
	return !isDir, nil
}

func isReadable(file interface{}) (bool, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return false, err
	}
	return stat.Mode()&0400 != 0, nil
}

func isWriteable(file interface{}) (bool, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return false, err
	}
	return stat.Mode()&0200 != 0, nil
}

func isExecutable(file interface{}) (bool, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return false, err
	}
	return stat.Mode()&0100 != 0 && !stat.IsDir(), nil
}

func lastMod(file interface{}) (time.Time, error) {
	stat, err := os.Stat(fmt.Sprint(file))
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

func glob(args ...interface{}) collections.IGenericList {
	return collections.AsList(utils.GlobFuncTrim(args...))
}

func diff(text1, text2 interface{}) interface{} {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(fmt.Sprint(text1), fmt.Sprint(text2), true)
	return dmp.DiffPrettyText(diffs)
}

func loadFromFile(filename string, binary ...bool) (interface{}, error) {
	isBinary := false
	switch len(binary) {
	case 0:
		break
	case 1:
		isBinary = binary[0]
	default:
		return "", fmt.Errorf("invalid number of arguments")
	}
	content, err := os.ReadFile(filename)
	if isBinary {
		return content, err
	}
	return string(content), err
}

func saveToFile(filename string, object interface{}) (string, error) {
	folder := path.Dir(filename)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err = os.Mkdir(folder, 0777); err != nil {
			return "", err
		}
	}
	if list, err := collections.TryAsList(object); err != nil {
		// If object is not a []byte, convert it to a string
		object = []byte(fmt.Sprint(object))
	} else {
		// Convert each element of the list to byte
		byteArray := make([]byte, len(list.AsArray()))
		for i, v := range list.AsArray() {
			byteValue, ok := v.(byte)
			if !ok {
				// Not all elements are byte, convert the whole list to a string
				byteArray = []byte(fmt.Sprint(object))
				break
			}
			byteArray[i] = byteValue
		}
		object = byteArray
	}

	return "", os.WriteFile(filename, object.([]byte), 0644)
}

func username() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.Username
}

func userGroup() (*user.Group, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	return user.LookupGroupId(u.Gid)
}

func userHome() string {
	u, err := user.Current()
	if err != nil {
		return os.Getenv("HOME")
	}
	return u.HomeDir
}

func lookPath(file interface{}) string {
	path, _ := exec.LookPath(fmt.Sprint(file))
	return path
}
