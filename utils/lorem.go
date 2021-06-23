package utils

import (
	"fmt"
	"strings"

	goLorem "github.com/drhodes/goLorem"
)

// LoremKind represents the various Lorem Ipsum generator type
type LoremKind int

// Constant used to describe the different kind of lorem generator
const (
	_ LoremKind = iota
	Word
	Sentence
	Paragraph
	Host
	EMail
	URL
)

// GetLoremKind converts a name to LoremKind
func GetLoremKind(name string) (kind LoremKind, err error) {
	switch strings.ToLower(name) {
	case "1", "word":
		kind = Word
	case "2", "", "words", "sentence":
		kind = Sentence
	case "3", "para", "paragraph", "sentences":
		kind = Paragraph
	case "4", "host":
		kind = Host
	case "5", "email":
		kind = EMail
	case "6", "url":
		kind = URL
	default:
		err = fmt.Errorf("undefined Lorem kind %s", name)
	}
	return
}

// Lorem generates random string using lorem ipsum generator
func Lorem(kind LoremKind, params ...int) (string, error) {
	min := 3
	max := 10
	if len(params) > 0 {
		min = params[0]
	}
	if len(params) > 1 {
		max = params[1]
	}
	switch kind {
	case Sentence:
		return goLorem.Sentence(min, max), nil
	case Paragraph:
		return goLorem.Paragraph(min, max), nil
	case Word:
		return goLorem.Word(min, max), nil
	case Host:
		return goLorem.Host(), nil
	case EMail:
		return goLorem.Email(), nil
	case URL:
		return goLorem.Url(), nil
	default:
		return "", fmt.Errorf("unknown lorem type %v", kind)
	}
}
