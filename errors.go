package deepseek

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the DeepSeek API.
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("deepseek: %s (code: %d, type: %s, status: %d)", e.Message, e.Code, e.Type, e.StatusCode)
}

// RequestError represents an error that occurred while making a request to the API.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	return fmt.Sprintf("deepseek: request failed with status %d: %v", e.HTTPStatusCode, e.Err)
}

// Unwrap returns the underlying error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// InvalidRequestError represents an error due to invalid request parameters.
type InvalidRequestError struct {
	Err error
}

// Error implements the error interface.
func (e *InvalidRequestError) Error() string {
	return fmt.Sprintf("deepseek: invalid request: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *InvalidRequestError) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication error.
type AuthenticationError struct {
	Err error
}

// Error implements the error interface.
func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("deepseek: authentication failed: %v", e.Err)
}

// Unwrap returns the underlying error.
func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

// handleErrorResp creates an appropriate error type based on the HTTP status code.
func handleErrorResp(resp *http.Response, apiErr *APIError) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &AuthenticationError{Err: fmt.Errorf("%s", apiErr.Message)}
	case http.StatusBadRequest:
		return &InvalidRequestError{Err: fmt.Errorf("%s", apiErr.Message)}
	default:
		return &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            fmt.Errorf("%s", apiErr.Message),
		}
	}
}
