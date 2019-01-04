package main

import (
	"reflect"
	"testing"
)

func TestExclude(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		files    []string
		patterns []string
		result   []string
	}{
		{
			name:     "Empty",
			files:    []string{},
			patterns: []string{},
			result:   []string{},
		},
		{
			name:     "No exclude pattern",
			files:    []string{"test.txt"},
			patterns: []string{},
			result:   []string{"test.txt"},
		},
		{
			name:     "Folder with wildcard",
			files:    []string{"./1.txt", "./f/1.txt", "./f/2.txt"},
			patterns: []string{"./f/*"},
			result:   []string{"./1.txt"},
		},
		{
			name:     "Nested folders with double star wildcard",
			files:    []string{"./f/f2/1.txt"},
			patterns: []string{"./f/**"},
			result:   []string{},
		},
		{
			name:     "Exclude only nested folders with wildcard",
			files:    []string{"./f/f2/1.txt", "./f/f3/1.txt"},
			patterns: []string{"./f/f2/*.txt"},
			result:   []string{"./f/f3/1.txt"},
		},
		{
			name:     "Files",
			files:    []string{"./1.txt", "./2.txt", "./f/1.txt", "./f/2.txt"},
			patterns: []string{"./1.txt", "./f/1.txt"},
			result:   []string{"./2.txt", "./f/2.txt"},
		},
		{
			name:     "Files with wildcard",
			files:    []string{"./1.txt", "./2.txt", "./1.data"},
			patterns: []string{"./*.txt"},
			result:   []string{"./1.data"},
		},
		{
			name:     "Relative path without .",
			files:    []string{"./1.txt", "2.txt", "./1.data", "2.data"},
			patterns: []string{"*.txt"},
			result:   []string{"./1.data", "2.data"},
		},
		{
			name:     "Absolute paths",
			files:    []string{"/home/user/file.txt", "/home/user/file.data"},
			patterns: []string{"/home/user/*.txt"},
			result:   []string{"/home/user/file.data"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exclude(tt.files, tt.patterns)
			if err != nil {
				t.Errorf("exclude(%v, %v), got err: %v", tt.files, tt.patterns, err)
			}
			if !reflect.DeepEqual(got, tt.result) {
				t.Errorf("exclude(%v, %v) = %v, want %v", tt.files, tt.patterns, got, tt.result)
			}
		})
	}
}
