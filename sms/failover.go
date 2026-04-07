package sms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chokey2nv/go-notification/sms/providers"
)

type FailoverProvider struct {
	providers []providers.SMSProvider
	timeout   time.Duration
}

func NewFailoverProvider(timeout time.Duration, providers ...providers.SMSProvider) *FailoverProvider {
	return &FailoverProvider{
		providers: providers,
		timeout:   timeout,
	}
}

func (f *FailoverProvider) Name() string {
	return "failover"
}

func (f *FailoverProvider) SendSMS(
	ctx context.Context,
	message, sender string,
	recipients []string,
) (*providers.SMSResult, error) {

	if len(f.providers) == 0 {
		return nil, fmt.Errorf("no sms providers configured")
	}

	var errorsList []string

	for _, provider := range f.providers {

		// Apply per-provider timeout
		pCtx := ctx
		if f.timeout > 0 {
			var cancel context.CancelFunc
			pCtx, cancel = context.WithTimeout(ctx, f.timeout)
			defer cancel()
		}

		result, err := provider.SendSMS(pCtx, message, sender, recipients)

		// SUCCESS (even partial success counts — your call)
		if err == nil && result != nil && result.Success {
			result.Provider = provider.Name()
			return result, nil
		}

		// Collect error
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("[%s] %v", provider.Name(), err))
		} else {
			errorsList = append(errorsList, fmt.Sprintf("[%s] unknown failure", provider.Name()))
		}
	}

	return nil, fmt.Errorf("all providers failed: %s", strings.Join(errorsList, " | "))
}
