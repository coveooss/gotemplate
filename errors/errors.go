package errors

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-errors/errors"
)

// Raise issues a panic with a formated message representing a managed error
func Raise(format string, args ...interface{}) {
	panic(Managed(fmt.Sprintf(format, args...)))
}

// Printf print a message to the stderr in red
func Printf(format string, args ...interface{}) {
	Print(fmt.Errorf(format, args...))
}

// Print print error to the stderr in red
func Print(err error) {
	fmt.Fprintln(color.Error, color.RedString(fmt.Sprintf("%v", err)))
}

// Must traps errors and return the remaining results to the caller
// If there is an error, a panic is issued
func Must(result ...interface{}) interface{} {
	if len(result) == 0 {
		return nil
	}
	last := len(result) - 1
	err := result[last]

	if err != nil {
		panic(errors.WrapPrefix(err, "errors.Must call failed", 2))
	}

	result = result[:last]
	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	default:
		return result
	}
}

// Trap catch any panic exception, and convert it to error
// It must be called with current error state and recover() as argument
func Trap(sourceErr error, recovered interface{}) (err error) {
	err = sourceErr
	var trap error
	switch rec := recovered.(type) {
	case nil:
	case error:
		trap = rec
	default:
		trap = fmt.Errorf("%v", rec)
	}
	if trap != nil {
		if err != nil {
			return Array{err, trap}
		}
		return trap
	}
	return
}

// TemplateNotFoundError is returned when a template does not exist
type TemplateNotFoundError struct {
	name string
}

func (e TemplateNotFoundError) Error() string {
	return fmt.Sprintf("Template %s not found", e.name)
}

// Array represent an array of error
type Array []error

func (errors Array) Error() string {
	errorsStr := make([]string, 0, len(errors))
	for i := range errors {
		if errors[i] != nil {
			errorsStr = append(errorsStr, errors[i].Error())
		}
	}
	return strings.Join(errorsStr, "\n")
}

// AsError returns nil if there is no error in the array
func (errors Array) AsError() error {
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}
	return errors
}

// Managed indicates an error that is properly handled (no need to print out stack)
type Managed string

func (err Managed) Error() string {
	return string(err)
}
