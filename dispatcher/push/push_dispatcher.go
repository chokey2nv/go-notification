package push

import (
	"context"
	"errors"

	"github.com/chokey2nv/go-notification/dispatcher/email"
	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/helpers"
)
type PushDispatcher struct {
	client     PushClient
	templates  email.TemplateRenderer
	userLookup func(ctx context.Context, userID string) (deviceToken string, err error)
}
func NewPushDispatcher(
	client PushClient,
	templates email.TemplateRenderer,
	userLookup func(ctx context.Context, userID string) (string, error),
) *PushDispatcher {
	return &PushDispatcher{
		client:     client,
		templates:  templates,
		userLookup: userLookup,
	}
}
func (d *PushDispatcher) Channel() domain.Channel {
	return domain.ChannelPush
}
func (d *PushDispatcher) Send(
	ctx context.Context,
	n *domain.Notification,
) error {

	if d.userLookup == nil {
		return errors.New("push dispatcher: user lookup not configured")
	}

	token, err := d.userLookup(ctx, n.UserID)
	if err != nil {
		return err
	}

	title := n.Title
	body := n.Message

	if d.templates != nil {
		if tplID, ok := n.Metadata["template_id"]; ok {
			t, b, err := d.templates.Render(ctx, tplID, n.Metadata)
			if err != nil {
				return err
			}
			title = t
			body = helpers.StripHTML(b)
		}
	}

	return d.client.Send(ctx, token, title, body, n.Metadata)
}
