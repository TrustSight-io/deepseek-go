package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/trustsight-io/deepseek-go/internal/errors"
)

// Role represents a chat message role
type Role string

const (
	// RoleSystem represents a system message in a chat conversation
	RoleSystem Role = "system"
	// RoleUser represents a user message in a chat conversation
	RoleUser Role = "user"
	// RoleAssistant represents an assistant message in a chat conversation
	RoleAssistant Role = "assistant"
	// RoleFunction represents a function response message in a chat conversation
	RoleFunction Role = "function"
)

// Message represents a chat message
type Message struct {
	Role         Role          `json:"role"`
	Content      string        `json:"content"`
	Name         string        `json:"name,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// FunctionCall represents a function call in a chat message
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// Function represents a function that can be called by the model
type Function struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ChatCompletionRequest represents a request to the chat completions API
type ChatCompletionRequest struct {
	Model            string             `json:"model,omitempty"`
	Messages         []Message          `json:"messages"`
	Functions        []Function         `json:"functions,omitempty"`
	FunctionCall     string             `json:"function_call,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
	N                int                `json:"n,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             string             `json:"user,omitempty"`
	ResponseFormat   *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
	Seed     int64  `json:"seed,omitempty"`
	Tools    []Tool `json:"tools,omitempty"`
	JSONMode bool   `json:"json_mode,omitempty"`
}

// Tool represents a tool that can be used by the model
type Tool struct {
	Type     string    `json:"type"`
	Function *Function `json:"function,omitempty"`
}

// ChatCompletionResponse represents a response from the chat completions API
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CreateChatCompletion sends a chat completion request to the API
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	req *ChatCompletionRequest,
) (*ChatCompletionResponse, error) {
	if req == nil {
		return nil, &errors.InvalidRequestError{
			Param: "request",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if len(req.Messages) == 0 {
		return nil, &errors.InvalidRequestError{
			Param: "messages",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	if req.Model == "" {
		req.Model = "deepseek-chat"
	}

	resp, err := c.createChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %v", cerr)
		}
	}()

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// createChatCompletion handles the raw HTTP request to the chat completions API
func (c *Client) createChatCompletion(
	ctx context.Context,
	req *ChatCompletionRequest,
) (*http.Response, error) {
	httpReq, err := c.newRequest(ctx, http.MethodPost, "/chat/completions", req)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if cerr := resp.Body.Close(); cerr != nil {
			return nil, fmt.Errorf("error closing response body: %v", cerr)
		}

		var apiErr errors.APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
		}

		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	return resp, nil
}
