package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestColor(t *testing.T) {
	tests := []struct {
		args    []string
		want    *color.Color
		wantErr bool
	}{
		{[]string{"red   ;green"}, color.New(color.FgRed, color.FgGreen), false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.args), func(t *testing.T) {
			got, err := Color(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color() = %v, want %v", got, tt.want)
			}
		})
	}
}
