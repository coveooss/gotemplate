package main

import (
	"fmt"
	"strings"

	goLorem "github.com/drhodes/goLorem"
)

func lorem(function string, params ...int) (string, error) {
	min := 3
	max := 10
	if len(params) > 0 {
		min = params[0]
	}
	if len(params) > 1 {
		max = params[1]
	}
	switch strings.ToLower(function) {
	case "":
		fallthrough
	case "words":
		fallthrough
	case "sentence":
		return goLorem.Sentence(min, max), nil
	case "para":
		fallthrough
	case "paragraph":
		return goLorem.Paragraph(min, max), nil
	case "word":
		return goLorem.Word(min, max), nil
	case "host":
		return goLorem.Host(), nil
	case "email":
		return goLorem.Email(), nil
	case "url":
		return goLorem.Url(), nil
	default:
		return "", fmt.Errorf("Unknown lorem type %s", function)
	}
}
