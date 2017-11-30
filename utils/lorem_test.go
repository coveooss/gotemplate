package utils

import (
	"testing"
)

func TestLorem(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// WARNING: If you add tests, you must add them after those that are already there, otherwise, all results will be changed.
		{"Word", "expedita", false},
		{"1", "hac", false},
		{"Sentence", "Laudes en sequatur aer deo vos.", false},
		{"Url", "http://www.dicamfactis.net/integer", false},
		{"EMail", "dicentium@montiumita.org", false},
		{"Anything", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, _ := GetLoremKind(tt.name)
			got, err := Lorem(kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lorem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Lorem() = %v, want %v", got, tt.want)
			}
		})
	}
}
