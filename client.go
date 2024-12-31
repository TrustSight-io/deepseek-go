package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/trustsight-io/deepseek-go/internal/errors"
	"github.com/trustsight-io/deepseek-go/internal/util"
)

// Version represents the current version of the client
const Version = "v0.1.0"

const (
	defaultBaseURL        = "https://api.deepseek.com"
	defaultTimeout        = 30 * time.Second
	defaultMaxRetries     = 3
	defaultRetryWaitTime  = 1 * time.Second
	defaultMaxRequestSize = 2 << 20 // 2MB
)

// Client represents a DeepSeek API client with all configuration options
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client

	// Configuration options
	maxRetries     int
	retryWaitTime  time.Duration
	maxRequestSize int64

	// Feature flags
	enableRetries bool
	debug         bool
}

// ClientOption represents a function that modifies the client configuration
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithMaxRetries sets the maximum number of retries for failed requests
func WithMaxRetries(retries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = retries
		c.enableRetries = retries > 0
	}
}

// WithRetryWaitTime sets the wait time between retries
func WithRetryWaitTime(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.retryWaitTime = duration
	}
}

// WithMaxRequestSize sets the maximum request size in bytes
func WithMaxRequestSize(size int64) ClientOption {
	return func(c *Client) {
		c.maxRequestSize = size
	}
}

// WithDebug enables debug logging
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient creates a new DeepSeek API client with the provided options
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	client := &Client{
		baseURL:        defaultBaseURL,
		apiKey:         apiKey,
		maxRetries:     defaultMaxRetries,
		retryWaitTime:  defaultRetryWaitTime,
		maxRequestSize: defaultMaxRequestSize,
		enableRetries:  true,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// Close closes any resources held by the client
func (c *Client) Close() error {
	// Currently no resources to clean up
	return nil
}

// newRequest creates a new HTTP request with the given method, path, and body
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode request body: %v", err)
		}
	}

	url := util.JoinURL(c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("User-Agent", "deepseek-go/"+Version)
	return req, nil
}

// do executes an HTTP request with retries and error handling
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		resp, body, err := c.executeRequest(req)
		if err != nil {
			lastErr = err
			if !c.shouldRetryRequest(attempt, err) {
				return err
			}
			time.Sleep(c.retryWaitTime * time.Duration(attempt+1))
			continue
		}

		if err := c.handleResponse(resp, body, v); err != nil {
			lastErr = err
			if !c.shouldRetryResponse(attempt, resp.StatusCode) {
				return err
			}
			time.Sleep(c.retryWaitTime * time.Duration(attempt+1))
			continue
		}

		return nil
	}

	return lastErr
}

// executeRequest executes a single HTTP request and returns the response and body
func (c *Client) executeRequest(req *http.Request) (*http.Response, []byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return resp, body, nil
}

// handleResponse processes the HTTP response and handles any errors
func (c *Client) handleResponse(resp *http.Response, body []byte, v interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if util.IsHTML(body) {
			return fmt.Errorf("received HTML response with status %d", resp.StatusCode)
		}

		var apiErr errors.APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
		}
		apiErr.StatusCode = resp.StatusCode
		return errors.HandleErrorResp(resp, &apiErr)
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return nil
}

// shouldRetryRequest determines if a request error should trigger a retry
func (c *Client) shouldRetryRequest(attempt int, _ error) bool {
	return c.enableRetries && attempt < c.maxRetries
}

// shouldRetryResponse determines if a response should trigger a retry based on status code
func (c *Client) shouldRetryResponse(attempt int, statusCode int) bool {
	return shouldRetry(statusCode) && c.enableRetries && attempt < c.maxRetries
}

// shouldRetry returns true if the status code indicates a retryable error
func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}
