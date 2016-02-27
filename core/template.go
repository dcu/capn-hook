package core

import (
	"fmt"
	"strings"
)

// Vars is the variables map for the template
type Vars map[string]string

// Template takes a text and applies the variables to it.
type Template struct {
	Text string
}

// EscapeStringArray quotes the values included in the given array and returns the the values joined by a whitespace.
func EscapeStringArray(arr []string) string {
	return strings.Join(arr, " ")
}

// HasTemplateVariables returns true if the text has template variables.
func HasTemplateVariables(text string) bool {
	return strings.ContainsAny(text, "{}")
}

// Eval evaluates the template given the variables.
func (template *Template) Eval(variables Vars) string {
	output := template.Text
	for name, value := range variables {
		templateVar := fmt.Sprintf("{%s}", name)

		output = strings.Replace(output, templateVar, value, -1)
	}

	return output
}
