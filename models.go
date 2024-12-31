package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// Model represents a DeepSeek model
type Model struct {
	ID            string          `json:"id"`
	Object        string          `json:"object"`
	Created       int64           `json:"created"`
	OwnedBy       string          `json:"owned_by"`
	ContextWindow int             `json:"context_window"`
	Capabilities  map[string]bool `json:"capabilities"`
	MaxTokens     int             `json:"max_tokens"`
	PricingConfig *PricingConfig  `json:"pricing_config,omitempty"`
}

// PricingConfig represents the pricing configuration for a model
type PricingConfig struct {
	PromptTokenPrice     float64 `json:"prompt_token_price"`
	CompletionTokenPrice float64 `json:"completion_token_price"`
	Currency             string  `json:"currency"`
}

// ListModelsResponse represents the response from the list models API
type ListModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// ListModels returns a list of available models
func (c *Client) ListModels(ctx context.Context) (*ListModelsResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/models", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var response ListModelsResponse
	if err := c.do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetModel retrieves information about a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, error) {
	if modelID == "" {
		return nil, &errors.InvalidRequestError{Param: "modelID", Err: fmt.Errorf("cannot be empty")}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/models/%s", c.baseURL, modelID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var model Model
	if err := c.do(ctx, req, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// ModelConfig represents configuration options for a model
type ModelConfig struct {
	ContextWindow     int                    `json:"context_window"`
	MaxTokens         int                    `json:"max_tokens"`
	SupportedFeatures map[string]bool        `json:"supported_features"`
	DefaultParameters map[string]interface{} `json:"default_parameters"`
}

// GetModelConfig retrieves the configuration for a specific model
func (c *Client) GetModelConfig(ctx context.Context, modelID string) (*ModelConfig, error) {
	if modelID == "" {
		return nil, &errors.InvalidRequestError{Param: "modelID", Err: fmt.Errorf("cannot be empty")}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/models/%s/config", c.baseURL, modelID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var config ModelConfig
	if err := c.do(ctx, req, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
