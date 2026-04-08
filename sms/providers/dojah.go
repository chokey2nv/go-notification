package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const dojahDefaultURL = "https://api.dojah.io"

// --------------------
// Provider
// --------------------

type DojahProvider struct {
	baseURL   string
	appID     string
	secretKey string
	senderID  string
	priority  bool
	client    *http.Client
	hooks     *Hooks
	retries   int
	workers   int // concurrency control
}

// --------------------
// Config
// --------------------

type DojahConfig struct {
	BaseURL    string
	AppID      string
	SecretKey  string
	SenderID   string
	Priority   bool
	Timeout    int
	HTTPClient *http.Client
	Hooks      *Hooks
	Retries    int
	Workers    int
}

// --------------------
// Constructor
// --------------------
func NewDojahProvider(config DojahConfig) *DojahProvider {
	baseURL := strings.TrimSuffix(strings.TrimSpace(config.BaseURL), "/")
	if baseURL == "" {
		baseURL = dojahDefaultURL
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

	workers := config.Workers
	if workers <= 0 {
		workers = 5 // safe default
	}

	return &DojahProvider{
		baseURL:   baseURL,
		appID:     config.AppID,
		secretKey: config.SecretKey,
		senderID:  config.SenderID,
		priority:  config.Priority,
		client:    client,
		hooks:     config.Hooks,
		retries:   retries,
		workers:   workers,
	}
}

func (p *DojahProvider) Name() string {
	return "dojah"
}

// --------------------
// Public API
// --------------------

func (p *DojahProvider) SendSMS(
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

	result := &SMSResult{
		Success:    true,
		Provider:   p.Name(),
		Recipients: make([]RecipientResult, len(recipients)),
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, p.workers)

	failures := make([]string, 0)
	mu := sync.Mutex{}

	for i, recipient := range recipients {
		wg.Add(1)

		go func(i int, recipient string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			res, err := p.sendWithRetry(ctx, message, senderID, recipient)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failures = append(failures, fmt.Sprintf("%s:%v", recipient, err))
				result.Recipients[i] = RecipientResult{
					Number: recipient,
					Status: "Failed",
				}
				result.Success = false
				return
			}

			result.Recipients[i] = res

			if result.MessageID == "" && res.MessageID != "" {
				result.MessageID = res.MessageID
			}

		}(i, recipient)
	}

	wg.Wait()

	if len(failures) > 0 {
		return result, fmt.Errorf("partial failure: %s", strings.Join(failures, ", "))
	}

	return result, nil
}

// --------------------
// Retry Wrapper
// --------------------

func (p *DojahProvider) sendWithRetry(
	ctx context.Context,
	message, senderID, recipient string,
) (RecipientResult, error) {

	var lastErr error

	for attempt := 0; attempt < p.retries; attempt++ {

		res, err := p.sendToRecipient(ctx, message, senderID, recipient)
		if err == nil {
			return res, nil
		}

		lastErr = err
		p.handleError(err)
		p.backoff(attempt)
	}

	return RecipientResult{}, fmt.Errorf("failed after retries: %w", lastErr)
}

// --------------------
// Core Request
// --------------------

func (p *DojahProvider) sendToRecipient(
	ctx context.Context,
	message, senderID, recipient string,
) (RecipientResult, error) {

	payload := map[string]interface{}{
		"destination": recipient,
		"message":     message,
		"channel":     "sms",
		"sender_id":   senderID,
		"priority":    p.priority,
	}

	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v1/messaging/sms", p.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return RecipientResult{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AppId", p.appID)
	req.Header.Set("Authorization", p.secretKey)

	if p.hooks != nil && p.hooks.OnRequest != nil {
		p.hooks.OnRequest(req)
	}

	start := time.Now()
	resp, err := p.client.Do(req)
	if err != nil {
		return RecipientResult{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if p.hooks != nil && p.hooks.OnResponse != nil {
		p.hooks.OnResponse(bodyBytes, time.Since(start))
	}

	if resp.StatusCode >= 300 {
		return RecipientResult{}, &ProviderError{
			StatusCode: resp.StatusCode,
			Err:        errors.New("non-2xx response"),
		}
	}

	var r struct {
		Entity struct {
			MessageID string `json:"message_id"`
			Status    string `json:"status"`
		}
		Error interface{} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return RecipientResult{}, err
	}

	if r.Error != nil {
		return RecipientResult{}, fmt.Errorf("dojah error: %v", r.Error)
	}

	return RecipientResult{
		Number:    recipient,
		Status:    r.Entity.Status,
		MessageID: r.Entity.MessageID,
	}, nil
}

// --------------------
// Utils
// --------------------

func (p *DojahProvider) backoff(attempt int) {
	base := 200 * time.Millisecond
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(time.Duration(attempt+1)*base + jitter)
}

func (p *DojahProvider) handleError(err error) {
	if p.hooks != nil && p.hooks.OnError != nil {
		p.hooks.OnError(err)
	}
}
