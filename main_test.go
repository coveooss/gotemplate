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

func TestCliNonExistingSource(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gotemplate", "--source", "/path/that/does/not/exit"}
	exitCode := runGotemplate()
	assert.Equal(t, 1, exitCode)

	os.Args = []string{"gotemplate", "--source", "/path/that/does/not/exit", "--ignore-missing-source"}
	exitCode = runGotemplate()
	assert.Equal(t, 0, exitCode)
}
