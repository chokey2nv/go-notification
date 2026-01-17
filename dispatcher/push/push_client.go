package push

import "context"

type PushClient interface {
	Send(
		ctx context.Context,
		deviceToken string,
		title string,
		body string,
		data map[string]string,
	) error
}
