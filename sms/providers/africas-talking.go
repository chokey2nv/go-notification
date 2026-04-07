package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const africasTalkingDefaultURL = "https://api.africastalking.com/version1/messaging/bulk"

type ProviderConfig struct {
	BaseURL    string
	Username   string
	APIKey     string
	SenderID   string
	Timeout    int
	HTTPClient *http.Client
	Hooks      *Hooks
	Retries    int
}

// --------------------
// Provider
// --------------------

type AfricasTalkingProvider struct {
	baseURL  string
	username string
	apiKey   string
	senderID string
	client   *http.Client
	hooks    *Hooks
	retries  int
}

// --------------------
// Response Struct
// --------------------

type africasTalkingResponse struct {
	SMSMessageData *struct {
		Message    string `json:"Message"`
		Recipients []struct {
			StatusCode int    `json:"statusCode"`
			Number     string `json:"number"`
			Status     string `json:"status"`
			Cost       string `json:"cost"`
			MessageID  string `json:"messageId"`
		} `json:"Recipients"`
	} `json:"SMSMessageData,omitempty"`

	Error string `json:"error,omitempty"`
	ErrNo string `json:"errno,omitempty"`
}

// --------------------
// Constructor
// --------------------

func NewAfricasTalkingProvider(config ProviderConfig) *AfricasTalkingProvider {
	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		baseURL = africasTalkingDefaultURL
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	client := config.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	retries := config.Retries
	if retries <= 0 {
		retries = 3
	}

	return &AfricasTalkingProvider{
		baseURL:  baseURL,
		username: config.Username,
		apiKey:   config.APIKey,
		senderID: config.SenderID,
		client:   client,
		hooks:    config.Hooks,
		retries:  retries,
	}
}

// --------------------
// Public Methods
// --------------------

func (p *AfricasTalkingProvider) Name() string {
	return "africastalking"
}

func (p *AfricasTalkingProvider) SendSMS(
	ctx context.Context,
	message, sender string,
	recipients []string,
) (*SMSResult, error) {

	if len(recipients) == 0 {
		return nil, errors.New("no phone numbers provided")
	}

	senderID := strings.TrimSpace(sender)
	if senderID == "" {
		senderID = p.senderID
	}

	payload := struct {
		Username     string   `json:"username"`
		Message      string   `json:"message"`
		SenderID     string   `json:"senderId,omitempty"`
		PhoneNumbers []string `json:"phoneNumbers"`
	}{
		Username:     p.username,
		Message:      message,
		SenderID:     senderID,
		PhoneNumbers: recipients,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	var resp *http.Response
	var attemptErr error

	for attempt := 0; attempt < p.retries; attempt++ {

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		if p.apiKey != "" {
			req.Header.Set("apiKey", p.apiKey)
		}

		if p.hooks != nil && p.hooks.OnRequest != nil {
			p.hooks.OnRequest(req)
		}

		start := time.Now()
		resp, attemptErr = p.client.Do(req)

		if attemptErr != nil {
			p.handleError(attemptErr)
			p.backoff(attempt)
			continue
		}

		if p.hooks != nil && p.hooks.OnResponse != nil {
			p.hooks.OnResponse(resp, time.Since(start))
		}

		// Retry on 5xx
		if resp.StatusCode >= 500 {
			p.backoff(attempt)
			continue
		}

		break
	}

	if attemptErr != nil {
		return nil, fmt.Errorf("request failed after retries: %w", attemptErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &ProviderError{
			StatusCode: resp.StatusCode,
			Err:        errors.New("non-2xx response"),
		}
	}

	var atResp africasTalkingResponse
	if err := json.NewDecoder(resp.Body).Decode(&atResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if atResp.Error != "" {
		return nil, &ProviderError{
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("api error: %s (errno=%s)", atResp.Error, atResp.ErrNo),
		}
	}

	if atResp.SMSMessageData == nil {
		return nil, errors.New("invalid response: missing SMSMessageData")
	}

	return p.buildResult(atResp)
}

// --------------------
// Helpers
// --------------------

func (p *AfricasTalkingProvider) buildResult(atResp africasTalkingResponse) (*SMSResult, error) {

	result := &SMSResult{
		Success:     true,
		Provider:    p.Name(),
		Recipients:  make([]RecipientResult, 0, len(atResp.SMSMessageData.Recipients)),
		RawResponse: atResp,
	}

	var failures []string

	for _, r := range atResp.SMSMessageData.Recipients {

		res := RecipientResult{
			Number:     r.Number,
			Status:     r.Status,
			StatusCode: r.StatusCode,
			MessageID:  r.MessageID,
			Cost:       r.Cost,
		}

		result.Recipients = append(result.Recipients, res)

		if result.MessageID == "" && r.MessageID != "" {
			result.MessageID = r.MessageID
		}

		if !IsSuccess(r.Status, r.StatusCode) {
			failures = append(failures, fmt.Sprintf("%s:%s", r.Number, r.Status))
		}
	}

	if len(failures) > 0 {
		result.Success = false
		return result, fmt.Errorf("partial failure: %s", strings.Join(failures, ", "))
	}

	return result, nil
}

func (p *AfricasTalkingProvider) backoff(attempt int) {
	backoff(attempt)
}

func (p *AfricasTalkingProvider) handleError(err error) {
	if p.hooks != nil && p.hooks.OnError != nil {
		p.hooks.OnError(err)
	}
}
