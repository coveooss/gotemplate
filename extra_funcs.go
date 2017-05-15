package main

import (
	"bytes"
	"text/template"

	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

// AddExtraFuncs adds yaml, json and hcl converter functions to the current context
func AddExtraFuncs(template *template.Template) {
	runTemplate := func(templateName string, context ...interface{}) (string, error) {
		var out bytes.Buffer
		t := template.Lookup(templateName)
		if t == nil {
			return "", TemplateNotFoundError{templateName}
		}
		if err := t.Execute(&out, getContext(context)); err != nil {
			return out.String(), err
		}
		return out.String(), nil
	}

	// Internal function used to actually convert the supplied string and apply a conversion function over it to get a go map
	converter := func(str string, convFunc func([]byte, interface{}) error, context ...interface{}) (map[string]interface{}, error) {
		out := make(map[string]interface{})
		content, err := runTemplate(str, context...)
		if err != nil {
			if _, ok := err.(TemplateNotFoundError); !ok {
				return out, err
			}
			content = str
		}

		err = convFunc([]byte(content), &out)
		return out, err
	}

	// converts the supplied string containing yaml/json to go map
	yamlConverter := func(str string, context ...interface{}) (interface{}, error) {
		return converter(str, yaml.Unmarshal, context...)
	}

	// Converts the supplied string containing terraform/hcl to go map
	hclConverter := func(str string, context ...interface{}) (interface{}, error) {
		out, err := converter(str, hcl.Unmarshal, context...)
		return FlattenHCL(out), err
	}

	*template = *template.Funcs(map[string]interface{}{
		"yaml": yamlConverter,
		"json": yamlConverter,
		"hcl":  hclConverter,
	})
}
