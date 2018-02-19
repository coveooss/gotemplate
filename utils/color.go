package utils

import (
	"fmt"
	"strings"

	"github.com/coveo/gotemplate/errors"
	"github.com/fatih/color"
)

type Attribute color.Attribute

// The following constant are copied from the color package in order to get
// the actual names
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

//go:generate stringer -type=Attribute -output generated_colors.go

// Color returns a color attribute build from supplied attribute names
func Color(attributes ...string) (*color.Color, error) {
	if nameValues == nil {
		nameValues = make(map[string]color.Attribute, BgHiWhite)
		for i := Reset; i < BgHiWhite; i++ {
			name := strings.ToLower(Attribute(i).String())
			if strings.HasPrefix(name, "attribute(") {
				continue
			}
			nameValues[name] = color.Attribute(i)
			if strings.HasPrefix(name, "fg") {
				nameValues[name[2:]] = color.Attribute(i)
			}
		}
	}

	result := color.New()
	var containsColor bool
	var err errors.Array
	for _, attr := range attributes {
		for _, attr := range String(attr).FieldsId().Strings() {
			if a, match := nameValues[strings.ToLower(attr)]; match {
				result.Add(a)
				containsColor = true

			} else {
				err = append(err, fmt.Errorf("Attribute not found %s", attr))
			}
		}
	}

	if !containsColor {
		return result, fmt.Errorf("No color specified")
	}
	if len(err) > 0 {
		return result, err
	}
	return result, nil
}

// SprintColor returns a string formated with attributes that are supplied before
func SprintColor(args ...interface{}) (string, error) {
	var i int
	colorArgs := make([]string, len(args))
	for i = 0; i < len(args); i++ {
		if _, err := Color(fmt.Sprint(args[i])); err != nil {
			break
		}
		colorArgs[i] = fmt.Sprint(args[i])
	}

	c, _ := Color(colorArgs...)

	var format string
	if i < len(args) {
		format = fmt.Sprint(args[i])
	}

	if strings.Contains(format, "%") {
		return c.Sprintf(format, args[i+1:]...), nil
	}
	return c.Sprint(args[i:]...), nil
}

var nameValues map[string]color.Attribute
