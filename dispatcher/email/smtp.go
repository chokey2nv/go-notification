package email

import (
	"context"
	"fmt"
	"net/smtp"
)

type DefaultSMTPClient struct {
	host string
	port string
	auth smtp.Auth
}

func NewDefaultSMTPClient(
	host string,
	port string,
	username string,
	password string,
) *DefaultSMTPClient {
	return &DefaultSMTPClient{
		host: host,
		port: port,
		auth: smtp.PlainAuth("", username, password, host),
	}
}

func (c *DefaultSMTPClient) SendHTML(
	ctx context.Context,
	from string,
	to string,
	subject string,
	htmlBody string,
) error {

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n"+
			"%s",
		from,
		to,
		subject,
		htmlBody,
	))

	addr := fmt.Sprintf("%s:%s", c.host, c.port)
	return smtp.SendMail(addr, c.auth, from, []string{to}, msg)
}
