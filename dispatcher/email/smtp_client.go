package email

import "context"

type SMTPClient interface {
	SendHTML(
		ctx context.Context,
		from string,
		to string,
		subject string,
		htmlBody string,
	) error
}
