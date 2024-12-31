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
	req, err := c.newRequest(ctx, http.MethodGet, "/models", nil)
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

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/models/%s", modelID), nil)
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

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/models/%s/config", modelID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var config ModelConfig
	if err := c.do(ctx, req, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role         string      `json:"role"`
	Content      string      `json:"content"`
	Name         string      `json:"name,omitempty"`
	FunctionCall interface{} `json:"function_call,omitempty"`
}

// ChatCompletionRequest represents a request to create a chat completion
type ChatCompletionRequest struct {
	Model            string               `json:"model"`
	Messages         []ChatMessage        `json:"messages"`
	MaxTokens        int                  `json:"max_tokens,omitempty"`
	Temperature      float32              `json:"temperature,omitempty"`
	TopP             float32              `json:"top_p,omitempty"`
	N                int                  `json:"n,omitempty"`
	Stream           bool                 `json:"stream,omitempty"`
	Stop             []string             `json:"stop,omitempty"`
	PresencePenalty  float32              `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32              `json:"frequency_penalty,omitempty"`
	Functions        []FunctionDefinition `json:"functions,omitempty"`
	FunctionCall     interface{}          `json:"function_call,omitempty"`
	JSONMode         bool                 `json:"json_mode,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completion API
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage provides token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FunctionDefinition represents a function that can be called by the model
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters"`
}

// StreamChoice represents a choice in a streaming response
type StreamChoice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// Delta represents the change in a streaming response
type Delta struct {
	Role         string `json:"role,omitempty"`
	Content      string `json:"content,omitempty"`
	FunctionCall *struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function_call,omitempty"`
}

// ChatCompletionStream represents a streaming response
type ChatCompletionStream struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}
