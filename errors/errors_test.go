package errors

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	type i = interface{}
	tests := []struct {
		name    string
		result  []i
		want    i
		wantErr bool
	}{
		{"Nil", nil, nil, false},
		{"1 arg", []i{"Hello"}, nil, true},
		{"1 arg, nil", []i{nil}, nil, false},
		{"2 arg, second not nul", []i{"Hello", "World"}, nil, true},
		{"2 arg, second nul", []i{"Hello", nil}, "Hello", false},
		{"3 arg, last nul", []i{"Hello", "World", nil}, []i{"Hello", "World"}, false},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			defer func() { err = Trap(err, recover()) }()
			if got := Must(tt.result...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Must() = %v, want %v", got, tt.want)
			}
		})
		if (err != nil) != tt.wantErr {
			t.Errorf("CreateList() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestPrint(t *testing.T) {
	Printf("This is an error")
	var err error
	func() {
		defer func() { err = Trap(err, recover()) }()
		Raise("This is also an error")
	}()
	assert.NotNil(t, err, "Error should be not nil")
	Print(err)

	err = nil
	func() {
		defer func() { err = Trap(err, recover()) }()
		panic(911)
	}()
	assert.NotNil(t, err, "Error should be not nil")
	Print(err)

	// We left the err intact to test array creation
	func() {
		defer func() { err = Trap(err, recover()) }()
		panic(TemplateNotFoundError{"filename"})
	}()
	assert.NotNil(t, err, "Error should be not nil")
	Print(err)

	err = nil
	func() {
		defer func() { err = Trap(err, recover()) }()
		// No error
	}()
	assert.Nil(t, err, "Error should be nil")
}
