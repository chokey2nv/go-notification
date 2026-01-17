package email

import (
	"bytes"
	"context"
	"errors"
	"html/template"
)

type HTMLTemplate struct {
	Subject *template.Template
	Body    *template.Template
}

type HTMLTemplateRenderer struct {
	templates map[string]HTMLTemplate
}
func NewHTMLTemplateRenderer(
	raw map[string]struct {
		Subject string
		Body    string
	},
) (*HTMLTemplateRenderer, error) {

	tpls := make(map[string]HTMLTemplate)

	for id, tpl := range raw {
		subjectTpl, err := template.New(id + "_subject").Parse(tpl.Subject)
		if err != nil {
			return nil, err
		}

		bodyTpl, err := template.New(id + "_body").Parse(tpl.Body)
		if err != nil {
			return nil, err
		}

		tpls[id] = HTMLTemplate{
			Subject: subjectTpl,
			Body:    bodyTpl,
		}
	}

	return &HTMLTemplateRenderer{templates: tpls}, nil
}
func (r *HTMLTemplateRenderer) Render(
	ctx context.Context,
	templateID string,
	data map[string]string,
) (string, string, error) {

	tpl, ok := r.templates[templateID]
	if !ok {
		return "", "", errors.New("email template not found")
	}

	var subjectBuf bytes.Buffer
	var bodyBuf bytes.Buffer

	if err := tpl.Subject.Execute(&subjectBuf, data); err != nil {
		return "", "", err
	}

	if err := tpl.Body.Execute(&bodyBuf, data); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}
