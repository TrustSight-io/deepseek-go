package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// Balance represents the account's credit balance information.
type Balance struct {
	TotalCredits     float64 `json:"total_credits"`
	UsedCredits      float64 `json:"used_credits"`
	RemainingCredits float64 `json:"remaining_credits"`
	LastUpdated      string  `json:"last_updated"`
}

// GetBalance retrieves the current balance information for the account.
func (c *Client) GetBalance(ctx context.Context) (*Balance, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{Param: "context", Err: fmt.Errorf("cannot be nil")}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/balance", nil)
	if err != nil {
		return nil, err
	}

	var balance Balance
	if err := c.do(ctx, req, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// APIUsage represents detailed API usage information.
type APIUsage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Models    []struct {
		Model     string  `json:"model"`
		Requests  int64   `json:"requests"`
		Tokens    int64   `json:"tokens"`
		Credits   float64 `json:"credits"`
		TokenCost float64 `json:"token_cost"`
		TotalCost float64 `json:"total_cost"`
	} `json:"models"`
	TotalRequests int64   `json:"total_requests"`
	TotalTokens   int64   `json:"total_tokens"`
	TotalCredits  float64 `json:"total_credits"`
}

// GetUsage retrieves the API usage information for a specified time period.
func (c *Client) GetUsage(ctx context.Context, startDate, endDate string) (*APIUsage, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{Param: "context", Err: fmt.Errorf("cannot be nil")}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/usage", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if startDate != "" {
		q.Add("start_date", startDate)
	}
	if endDate != "" {
		q.Add("end_date", endDate)
	}
	req.URL.RawQuery = q.Encode()

	var usage APIUsage
	if err := c.do(ctx, req, &usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
