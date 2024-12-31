// Package errors provides custom error types and error handling for the DeepSeek client.
// It includes API errors, request errors, authentication errors, and rate limit errors.
package errors

import (
	"fmt"
	"net/http"
)

// Error types
const (
	ErrorTypeInvalidRequest        = "invalid_request"
	ErrorTypeAuthentication        = "authentication"
	ErrorTypePermission            = "permission"
	ErrorTypeRateLimit             = "rate_limit"
	ErrorTypeServer                = "server_error"
	ErrorTypeModelNotFound         = "model_not_found"
	ErrorTypeContextLengthExceeded = "context_length_exceeded"
	ErrorTypeTokenLimitExceeded    = "token_limit_exceeded"
)

// APIError represents an error returned by the DeepSeek API
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Param      string `json:"param,omitempty"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("deepseek: %s (code: %d, type: %s, status: %d)", e.Message, e.Code, e.Type, e.StatusCode)
}

// RequestError represents an error that occurred while making a request
type RequestError struct {
	StatusCode int
	Err        error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("deepseek: request failed with status %d: %v", e.StatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// InvalidRequestError represents an error due to invalid request parameters
type InvalidRequestError struct {
	Param string
	Err   error
}

func (e *InvalidRequestError) Error() string {
	if e.Param != "" {
		return fmt.Sprintf("deepseek: invalid request parameter '%s': %v", e.Param, e.Err)
	}
	return fmt.Sprintf("deepseek: invalid request: %v", e.Err)
}

func (e *InvalidRequestError) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication error
type AuthenticationError struct {
	Err error
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("deepseek: authentication failed: %v", e.Err)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	RetryAfter int
	Err        error
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("deepseek: rate limit exceeded, retry after %d seconds: %v", e.RetryAfter, e.Err)
}

func (e *RateLimitError) Unwrap() error {
	return e.Err
}

// ModelNotFoundError represents a model not found error
type ModelNotFoundError struct {
	Model string
	Err   error
}

func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("deepseek: model '%s' not found: %v", e.Model, e.Err)
}

func (e *ModelNotFoundError) Unwrap() error {
	return e.Err
}

// HandleErrorResp creates an appropriate error type based on the response
func HandleErrorResp(resp *http.Response, apiErr *APIError) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{Err: fmt.Errorf("%s", apiErr.Message)}
	case http.StatusForbidden:
		return &InvalidRequestError{Param: apiErr.Param, Err: fmt.Errorf("%s", apiErr.Message)}
	case http.StatusNotFound:
		if apiErr.Type == ErrorTypeModelNotFound {
			return &ModelNotFoundError{Model: apiErr.Param, Err: fmt.Errorf("%s", apiErr.Message)}
		}
		return &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", apiErr.Message)}
	case http.StatusTooManyRequests:
		retryAfter := 0
		if s := resp.Header.Get("Retry-After"); s != "" {
			if _, err := fmt.Sscanf(s, "%d", &retryAfter); err != nil {
				// If parsing fails, use default retry after value
				retryAfter = 60
			}
		}
		return &RateLimitError{RetryAfter: retryAfter, Err: fmt.Errorf("%s", apiErr.Message)}
	default:
		return &RequestError{
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("%s", apiErr.Message),
		}
	}
}
