// Package deepseek provides a Go client for the DeepSeek API.
// It supports chat completions, streaming, function calling, and more.
package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight-io/deepseek-go/internal/errors"
)

// BalanceInfo represents individual balance information for a currency
type BalanceInfo struct {
	Currency        string `json:"currency"`
	TotalBalance    string `json:"total_balance"`
	GrantedBalance  string `json:"granted_balance"`
	ToppedUpBalance string `json:"topped_up_balance"`
}

// Balance represents a user's balance information
type Balance struct {
	IsAvailable  bool          `json:"is_available"`
	BalanceInfos []BalanceInfo `json:"balance_infos"`
}

// GetBalance retrieves the current balance for the account
func (c *Client) GetBalance(ctx context.Context) (*Balance, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/user/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var balance Balance
	if err := c.do(ctx, req, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}
