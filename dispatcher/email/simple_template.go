package email

import (
	"bytes"
	"context"
	"errors"
	"text/template"
)

type SimpleTemplateRenderer struct {
	templates map[string]*template.Template
}

func NewSimpleTemplateRenderer(
	raw map[string]string,
) (*SimpleTemplateRenderer, error) {

	tpls := make(map[string]*template.Template)

	for id, tpl := range raw {
		parsed, err := template.New(id).Parse(tpl)
		if err != nil {
			return nil, err
		}
		tpls[id] = parsed
	}

	return &SimpleTemplateRenderer{templates: tpls}, nil
}

func (r *SimpleTemplateRenderer) Render(
	ctx context.Context,
	templateID string,
	data map[string]string,
) (string, string, error) {

	tpl, ok := r.templates[templateID]
	if !ok {
		return "", "", errors.New("template not found")
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", "", err
	}

	return data["subject"], buf.String(), nil
}
