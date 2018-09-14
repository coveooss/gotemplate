package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fatih/color"
)

func TestColor(t *testing.T) {
	t.Parallel()
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

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{"No argument", nil, ""},
		{"Empty arguments", []interface{}{}, ""},
		{"Single argument", []interface{}{"Hello"}, "Hello"},
		{"Two arguments", []interface{}{"Hello", "World"}, "Hello World"},
		{"Two arguments with format", []interface{}{"Hello %s! %d", "World", 100}, "Hello World! 100"},
		{"Bad format", []interface{}{"Hello %s! %d", "World"}, "Hello %s! %d World"},
		{"Escaped %", []interface{}{"You got %d%% off", 60}, "You got 60% off"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatMessage(tt.args...); got != tt.want {
				t.Errorf("FormatMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
