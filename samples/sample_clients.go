package samples

import (
	"context"
)

type sampleSmsClient struct{}

func NewSampleSmsClient() *sampleSmsClient { return &sampleSmsClient{} }

func (client *sampleSmsClient) Send(ctx context.Context, to string, message string) error { return nil }

type samplePushClient struct{}

func NewSamplePushClient() *samplePushClient { return &samplePushClient{} }

func (client *samplePushClient) Send(
	ctx context.Context,
	deviceToken string,
	title string,
	body string,
	data map[string]string,
) error {
	return nil
}
