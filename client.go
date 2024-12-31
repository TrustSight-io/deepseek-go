package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the DeepSeek API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string
	orgID      string
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a new DeepSeek API client.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("deepseek: API key is required")
	}

	c := &Client{
		apiKey: apiKey,
	}

	// Apply default options
	for _, opt := range defaultOptions() {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Apply user-provided options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// newRequest creates a new HTTP request.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	if c.orgID != "" {
		req.Header.Set("DeepSeek-Organization", c.orgID)
	}

	return req, nil
}

// do executes an HTTP request and decodes the response.
func (c *Client) do(req *http.Request, v interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * c.retryDelay)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if v != nil {
				if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
					return fmt.Errorf("deepseek: failed to decode response: %v", err)
				}
			}
			return nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("deepseek: failed to read error response: %v", err)
			continue
		}

		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			lastErr = fmt.Errorf("deepseek: failed to decode error response: %v", err)
			continue
		}

		apiErr.StatusCode = resp.StatusCode
		return handleErrorResp(resp, &apiErr)
	}

	return fmt.Errorf("deepseek: request failed after %d retries: %v", c.maxRetries, lastErr)
}
