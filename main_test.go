package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tempDir := must(ioutil.TempDir("", "gotemplate-test")).(string)
	templateFile := path.Join(tempDir, "test.template")
	templateFileContent := []byte("{{ $var := add 3 2 }}{{ $var }}")
	must(ioutil.WriteFile(templateFile, templateFileContent, 0644))

	os.Args = []string{"gotemplate", "--source", tempDir}
	exitCode := runGotemplate()

	assert.Equal(t, 0, exitCode)
	readTemplateFileContent := must(ioutil.ReadFile(templateFile)).([]byte)
	assert.Equal(t, templateFileContent, readTemplateFileContent)
	generatedFile := path.Join(tempDir, "test.generated")
	readGeneratedFileContent := must(ioutil.ReadFile(generatedFile)).([]byte)
	assert.Equal(t, []byte("5"), readGeneratedFileContent)
}

func TestErrorRecovery(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tempDir := must(ioutil.TempDir("", "gotemplate-test")).(string)
	templateFile := path.Join(tempDir, "test.template")
	templateFileContent := []byte("{{ panic `test string` }}")
	must(ioutil.WriteFile(templateFile, templateFileContent, 0644))

	os.Args = []string{"gotemplate", "--source", tempDir}
	exitCode := runGotemplate()

	assert.Equal(t, -1, exitCode)
}

func TestCliNonExistingSource(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gotemplate", "--source", "/path/that/does/not/exist"}
	exitCode := runGotemplate()
	assert.Equal(t, 1, exitCode)

	os.Args = []string{"gotemplate", "--source", "/path/that/does/not/exist", "--ignore-missing-source"}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)

	os.Args = []string{"gotemplate", "--source", "/path/that/does/not/exist", "--ignore-missing-paths"}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)
}

func TestCliNonExistingImport(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tempDir := must(ioutil.TempDir("", "gotemplate-test")).(string)
	templateFile := path.Join(tempDir, "test.template")
	templateFileContent := []byte("{{ $var := add 3 2 }}{{ $var }}")
	must(ioutil.WriteFile(templateFile, templateFileContent, 0644))

	// Regular variables
	os.Args = []string{"gotemplate", "--var", "test=123", "--var", "test=`/path/that/does/not/exist.json`", "--source", tempDir}
	exitCode := runGotemplate()
	assert.Equal(t, 0, exitCode)

	// Missing imports
	os.Args = []string{"gotemplate", "--import", "/path/that/does/not/exist.json", "--source", tempDir}
	exitCode = runGotemplate()
	assert.Equal(t, 1, exitCode)

	os.Args = []string{"gotemplate", "--var", "test=/path/that/does/not/exist.json", "--source", tempDir}
	exitCode = runGotemplate()
	assert.Equal(t, 1, exitCode)

	// Ignored missing imports
	os.Args = []string{"gotemplate", "--import", "/path/that/does/not/exist.json", "--ignore-missing-import", "--source", tempDir}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)

	os.Args = []string{"gotemplate", "--import", "/path/that/does/not/exist.json", "--ignore-missing-paths", "--source", tempDir}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)

	os.Args = []string{"gotemplate", "--var", "test=/path/that/does/not/exist.json", "--ignore-missing-import", "--source", tempDir}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)
}
