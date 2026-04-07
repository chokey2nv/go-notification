package providers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type SMSProvider interface {
	Name() string
	SendSMS(ctx context.Context, message, sender string, recipients []string) (*SMSResult, error)
}

type SMSResult struct {
	// Success indicates if the overall send operation succeeded
	Success bool `json:"success"`

	// Provider is the name of the provider that handled this request
	Provider string `json:"provider"`

	// MessageID is the primary message identifier (if available)
	MessageID string `json:"messageId,omitempty"`

	// Recipients contains the delivery status for each recipient
	Recipients []RecipientResult `json:"recipients,omitempty"`

	// RawResponse stores the original provider response for debugging
	RawResponse interface{} `json:"-"`
}

// RecipientResult contains the delivery status for a single recipient.
type RecipientResult struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode,omitempty"`
	MessageID  string `json:"messageId,omitempty"`
	Cost       string `json:"cost,omitempty"`
}

// --------------------
// Hooks (for observability)
// --------------------

type Hooks struct {
	OnRequest  func(*http.Request)
	OnResponse func(*http.Response, time.Duration)
	OnError    func(error)
}

func DefaultHooks() *Hooks {
	return &Hooks{}
}

// --------------------
// Error Types
// --------------------

type ProviderError struct {
	StatusCode int
	Body       string
	Err        error
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("sms provider error: status=%d err=%v body=%s", e.StatusCode, e.Err, e.Body)
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}

// --------------------
// Utils
// --------------------
func backoff(attempt int) {
	base := 200 * time.Millisecond
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(time.Duration(attempt+1)*base + jitter)
}
func IsSuccess(status string, code int) bool {
	return strings.EqualFold(status, "Success") || code == 101
}
