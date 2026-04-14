package utils

import (
	"strings"
	"testing"
)

func TestLorem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"Word", false},
		{"1", false},
		{"Sentence", false},
		{"Url", false},
		{"EMail", false},
		{"Anything", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, _ := GetLoremKind(tt.name)
			got, err := Lorem(kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lorem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && strings.TrimSpace(got) == "" {
				t.Errorf("Lorem() returned empty string for kind %q", tt.name)
			}
		})
	}
}
