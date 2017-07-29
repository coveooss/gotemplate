package main

import (
	"fmt"

	"github.com/go-errors/errors"
)

// PanicOnError traps errors and return the remaining results to the caller
// If there is an error, a panic is issued
func PanicOnError(result ...interface{}) interface{} {
	last := len(result) - 1
	err := result[last]

	if err != nil {
		panic(errors.WrapPrefix(err, "PanicOnError call failed", 2))
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
