package email

import (
	"context"

	"github.com/chokey2nv/go-notification/domain"
)

type EmailDispatcher struct {
	from       string
	smtp       SMTPClient
	templates  TemplateRenderer
	userLookup func(ctx context.Context, userID string) (email string, err error)
}

func NewEmailDispatcher(
	from string,
	smtp SMTPClient,
	templates TemplateRenderer,
	userLookup func(ctx context.Context, userID string) (string, error),
) *EmailDispatcher {
	return &EmailDispatcher{
		from:       from,
		smtp:       smtp,
		templates:  templates,
		userLookup: userLookup,
	}
}
func (d *EmailDispatcher) Channel() domain.Channel {
	return domain.ChannelEmail
}
func (d *EmailDispatcher) Send(
	ctx context.Context,
	n *domain.Notification,
) error {

	email, err := d.userLookup(ctx, n.UserID)
	if err != nil {
		return err
	}

	subject := n.Title
	htmlBody := n.Message

	if d.templates != nil {
		if tplID, ok := n.Metadata["template_id"]; ok {
			subject, htmlBody, err = d.templates.Render(ctx, tplID, n.Metadata)
			if err != nil {
				return err
			}
		}
	}

	return d.smtp.SendHTML(ctx, d.from, email, subject, htmlBody)
}
