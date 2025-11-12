package helpers

import (
	"bytes"
	"text/template"
)

type TemplateRenderer struct{}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{}
}

func (tr *TemplateRenderer) Render(templateStr string, variables map[string]any) (string, error) {
	if variables == nil {
		return templateStr, nil
	}

	tmpl, err := template.New("push").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
