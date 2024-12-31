package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// Balance represents a user's balance information
type Balance struct {
	Object         string  `json:"object"`
	TotalBalance   float64 `json:"total_balance"`
	Currency       string  `json:"currency"`
	GrantedQuota   float64 `json:"granted_quota"`
	UsedQuota      float64 `json:"used_quota"`
	RemainingQuota float64 `json:"remaining_quota"`
	QuotaResetTime string  `json:"quota_reset_time,omitempty"`
	ExpirationTime string  `json:"expiration_time,omitempty"`
	LastUpdated    string  `json:"last_updated"`
}

// GetBalance retrieves the current balance for the account
func (c *Client) GetBalance(ctx context.Context) (*Balance, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/billing/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var balance Balance
	if err := c.do(ctx, req, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// UsageRecord represents a usage record for a specific time period
type UsageRecord struct {
	Timestamp        string  `json:"timestamp"`
	RequestCount     int     `json:"request_count"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost"`
}

// UsageResponse represents the response from the usage API
type UsageResponse struct {
	Object    string        `json:"object"`
	Data      []UsageRecord `json:"data"`
	StartTime string        `json:"start_time"`
	EndTime   string        `json:"end_time"`
	Total     struct {
		RequestCount     int     `json:"request_count"`
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		TotalTokens      int     `json:"total_tokens"`
		TotalCost        float64 `json:"total_cost"`
	} `json:"total"`
}

// UsageParams represents parameters for retrieving usage history
type UsageParams struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Granularity string `json:"granularity,omitempty"` // hourly, daily, monthly
}

// GetUsage retrieves usage history for the account
func (c *Client) GetUsage(ctx context.Context, params *UsageParams) (*UsageResponse, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if params == nil {
		return nil, &errors.InvalidRequestError{
			Param: "params",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if params.StartTime == "" {
		return nil, &errors.InvalidRequestError{
			Param: "start_time",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	if params.EndTime == "" {
		return nil, &errors.InvalidRequestError{
			Param: "end_time",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/billing/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("start_time", params.StartTime)
	q.Add("end_time", params.EndTime)
	if params.Granularity != "" {
		q.Add("granularity", params.Granularity)
	}
	req.URL.RawQuery = q.Encode()

	var usage UsageResponse
	if err := c.do(ctx, req, &usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
