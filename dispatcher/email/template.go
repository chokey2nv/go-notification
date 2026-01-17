package email

import "context"

type TemplateRenderer interface {
	Render(
		ctx context.Context,
		templateID string,
		data map[string]string,
	) (subject string, htmlBody string, err error)
}
