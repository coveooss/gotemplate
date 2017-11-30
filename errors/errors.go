package errors

import (
	"fmt"
	"os"
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
	fmt.Fprintln(os.Stderr, color.RedString(fmt.Sprintf("%v", err)))
}

// Must traps errors and return the remaining results to the caller
// If there is an error, a panic is issued
func Must(result ...interface{}) interface{} {
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

// Managed indicates an error that is properly handled (no need to print out stack)
type Managed string

func (err Managed) Error() string {
	return string(err)
}
