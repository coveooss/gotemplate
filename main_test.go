package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	variableTempDir := t.TempDir()

	variableFile := path.Join(variableTempDir, "test.variables")
	err := os.WriteFile(variableFile, []byte("testInt = 5\ntestString = \"hello\"\ntestBool = true"), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		pipe           string
		args           []string
		template       string
		expectedCode   int
		expectedResult string
	}{
		{
			name:           "No arguments",
			args:           []string{},
			template:       "{{ add 3 2 }}",
			expectedCode:   0,
			expectedResult: "5",
		},
		{
			name:           "Integers",
			args:           []string{"--var", "test1=3", "--var", "test2=7"},
			template:       "{{ $var := add .test1 .test2 }}{{ $var }}",
			expectedCode:   0,
			expectedResult: "10",
		},
		{
			name:           "Import variables",
			args:           []string{"--import", variableFile, "--var", "test1=3"},
			template:       "{{ $var := add .test1 .testInt }}{{ $var }}",
			expectedCode:   0,
			expectedResult: "8",
		},
		{
			name:           "STDIN variables",
			pipe:           "testInt=3",
			args:           []string{"--var=-", "--var", "test2=4", "test.template"},
			template:       "{{ $var := add .testInt .test2 }}{{ $var }}",
			expectedCode:   0,
			expectedResult: "7",
		},

		// Ambiguous variables
		{
			name:           "Variables that may look like filenames",
			args:           []string{"--var", "test1=2.3.4", "--var", "test2=/var/path/literal"},
			template:       `{{ print .test1 "-" .test2 }}`,
			expectedCode:   0,
			expectedResult: "2.3.4-/var/path/literal",
		},
		{
			name:           "Variable file as var",
			args:           []string{"--var", "testFile=" + variableFile},
			template:       `{{ print .testFile.testInt "-" .testFile.testString }}`,
			expectedCode:   0,
			expectedResult: "5-hello",
		},
		{
			name:           "Literal existing file name",
			args:           []string{"--var", fmt.Sprintf("testFile='%s'", variableFile)},
			template:       "{{ print .testFile }}",
			expectedCode:   0,
			expectedResult: variableFile,
		},

		// Missing source
		{
			name:         "Non-existing source",
			args:         []string{"--source", "/path/that/does/not/exist"},
			expectedCode: 1,
		},
		{
			name:         "Non-existing source (Ignored)",
			args:         []string{"--source", "/path/that/does/not/exist", "--ignore-missing-source"},
			expectedCode: 0,
		},
		{
			name:         "Non-existing source (Ignored2)",
			args:         []string{"--source", "/path/that/does/not/exist", "--ignore-missing-paths"},
			expectedCode: 0,
		},

		// Missing import
		{
			name:         "Non-existing import",
			args:         []string{"--import", "/path/that/does/not/exist"},
			expectedCode: 1,
		},
		{
			name:         "Non-existing import (Ignored)",
			args:         []string{"--import", "/path/that/does/not/exist", "--ignore-missing-import"},
			expectedCode: 0,
		},
		{
			name:         "Non-existing import (Ignored2)",
			args:         []string{"--import", "/path/that/does/not/exist", "--ignore-missing-paths"},
			expectedCode: 0,
		},
		{
			name:         "Non-existing import",
			args:         []string{"--import-if-exist", "/path/that/does/not/exist"},
			expectedCode: 0,
		},

		// Errors
		{
			name:         "Error recovery",
			args:         []string{},
			template:     "{{ panic `test string` }}",
			expectedCode: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tempDir string

			// Change and revert (defer) os values
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			oldCw := must(os.Getwd()).(string)
			defer os.Chdir(oldCw)

			if tt.pipe != "" {
				content := []byte(tt.pipe)
				tmpStdin, err := os.CreateTemp("", "example")
				assert.NoError(t, err)
				defer os.Remove(tmpStdin.Name())
				must(tmpStdin.Write(content))
				must(tmpStdin.Seek(0, 0))

				oldStdin := os.Stdin
				defer func() {
					tmpStdin.Close()
					os.Stdin = oldStdin
				}()

				os.Stdin = tmpStdin
			}

			if tt.template != "" {
				// Create template
				tempDir = t.TempDir()
				os.Chdir(tempDir)
				templateFile := path.Join(tempDir, "test.template")
				err = os.WriteFile(templateFile, []byte(tt.template), 0644)
				assert.NoError(t, err)
				os.Args = append([]string{"gotemplate", "--source", tempDir}, tt.args...)
			} else {
				os.Args = append([]string{"gotemplate"}, tt.args...)
			}

			// Run gotemplate
			exitCode := runGotemplate()

			assert.Equal(t, tt.expectedCode, exitCode, "Bad exit code")

			if tt.expectedResult != "" {
				generatedFile := path.Join(tempDir, "test.generated")
				readGeneratedFileContent, err := os.ReadFile(generatedFile)
				assert.NoError(t, err)
				assert.Equal(t, []byte(tt.expectedResult), readGeneratedFileContent, "Bad generated content")
			}
		})
	}
}
