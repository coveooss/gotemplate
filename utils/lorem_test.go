package utils

import (
	"math/rand"
	"testing"
)

func TestLorem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"Word", "araneae", false},
		{"1", "araneae", false},
		{"Sentence", "Fac vel varias vim, cor munerum traiecta.", false},
		{"Url", "http://www.diuvi.com/curam/beata.html", false},
		{"EMail", "ardentius@curainest.net", false},
		{"Anything", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0) // This function is deprecated, but there is no alternative yet as the lorem package is using the default global rand generator
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
