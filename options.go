package deepseek

import (
	"net/http"
	"time"
)

// ClientOption is a function that configures a Client.
type ClientOption func(*Client) error

// WithBaseURL sets the base URL for the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// WithHTTPClient sets the HTTP client for making requests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout sets the timeout for requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithUserAgent sets the User-Agent header for requests.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithRetry configures retry behavior for failed requests.
func WithRetry(maxRetries int, initialRetryDelay time.Duration) ClientOption {
	return func(c *Client) error {
		c.maxRetries = maxRetries
		c.retryDelay = initialRetryDelay
		return nil
	}
}

// WithOrganization sets the organization ID for requests.
func WithOrganization(orgID string) ClientOption {
	return func(c *Client) error {
		c.orgID = orgID
		return nil
	}
}

// defaultOptions returns the default configuration options.
func defaultOptions() []ClientOption {
	return []ClientOption{
		WithBaseURL(DefaultAPIEndpoint),
		WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
		}),
		WithUserAgent(DefaultUserAgent),
		WithRetry(3, time.Second),
	}
}
