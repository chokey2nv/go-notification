package sms

import "context"

type SMSClient interface {
	Send(
		ctx context.Context,
		to string,
		message string,
	) error
}
