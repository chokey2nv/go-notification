package sms

import (
	"context"
	"errors"

	"github.com/chokey2nv/go-notification/dispatcher/email"
	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/helpers"
)

type SMSDispatcher struct {
	client     SMSClient
	templates  email.TemplateRenderer
	userLookup func(ctx context.Context, userID string) (phone string, err error)
}

func NewSMSDispatcher(
	client SMSClient,
	templates email.TemplateRenderer,
	userLookup func(ctx context.Context, userID string) (string, error),
) *SMSDispatcher {
	return &SMSDispatcher{
		client:     client,
		templates:  templates,
		userLookup: userLookup,
	}
}
func (d *SMSDispatcher) Channel() domain.Channel {
	return domain.ChannelSMS
}
func (d *SMSDispatcher) Send(
	ctx context.Context,
	n *domain.Notification,
) error {

	if d.userLookup == nil {
		return errors.New("sms dispatcher: user lookup not configured")
	}

	phone, err := d.userLookup(ctx, n.UserID)
	if err != nil {
		return err
	}

	message := n.Message

	if d.templates != nil {
		if tplID, ok := n.Metadata["template_id"]; ok {
			_, body, err := d.templates.Render(ctx, tplID, n.Metadata)
			if err != nil {
				return err
			}
			message = helpers.StripHTML(body)
		}
	}

	return d.client.Send(ctx, phone, message)
}
