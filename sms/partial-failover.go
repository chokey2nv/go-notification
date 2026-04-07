package sms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chokey2nv/go-notification/sms/providers"
)

type PartialFailoverProvider struct {
	providers []providers.SMSProvider
	timeout   time.Duration
}

func NewPartialFailoverProvider(timeout time.Duration, providers ...providers.SMSProvider) *PartialFailoverProvider {
	return &PartialFailoverProvider{
		providers: providers,
		timeout:   timeout,
	}
}

func (p *PartialFailoverProvider) Name() string {
	return "partial-failover"
}

func (p *PartialFailoverProvider) SendSMS(
	ctx context.Context,
	message, sender string,
	recipients []string,
) (*providers.SMSResult, error) {

	if len(p.providers) == 0 {
		return nil, fmt.Errorf("no sms providers configured")
	}

	remaining := make([]string, len(recipients))
	copy(remaining, recipients)

	finalResults := make(map[string]providers.RecipientResult)
	var errorsList []string

	for _, provider := range p.providers {

		if len(remaining) == 0 {
			break // all delivered
		}

		// timeout per provider
		pCtx := ctx
		if p.timeout > 0 {
			var cancel context.CancelFunc
			pCtx, cancel = context.WithTimeout(ctx, p.timeout)
			defer cancel()
		}

		result, err := provider.SendSMS(pCtx, message, sender, remaining)

		if err != nil && result == nil {
			errorsList = append(errorsList, fmt.Sprintf("[%s] %v", provider.Name(), err))
			continue
		}

		// Track next batch
		nextRemaining := []string{}

		for _, r := range result.Recipients {
			finalResults[r.Number] = r

			if !providers.IsSuccess(r.Status, r.StatusCode) {
				nextRemaining = append(nextRemaining, r.Number)
			}
		}

		remaining = nextRemaining
	}

	// Build final result
	final := &providers.SMSResult{
		Success:    len(remaining) == 0,
		Provider:   "multi",
		Recipients: make([]providers.RecipientResult, 0, len(finalResults)),
	}

	for _, r := range finalResults {
		final.Recipients = append(final.Recipients, r)
	}

	if len(remaining) > 0 {
		errorsList = append(errorsList, fmt.Sprintf("undelivered: %v", remaining))
		return final, fmt.Errorf("partial failure: %s", strings.Join(errorsList, " | "))
	}

	return final, nil
}
