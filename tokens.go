package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// TokenCount represents token count information
type TokenCount struct {
	TotalTokens int `json:"total_tokens"`
	Details     struct {
		Prompt     int  `json:"prompt_tokens"`
		Completion int  `json:"completion_tokens,omitempty"`
		Truncated  bool `json:"truncated,omitempty"`
	} `json:"details"`
}

// CountTokens counts the number of tokens in the given text for a specific model
func (c *Client) CountTokens(ctx context.Context, model string, text string) (*TokenCount, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if model == "" {
		return nil, &errors.InvalidRequestError{
			Param: "model",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	if text == "" {
		return nil, &errors.InvalidRequestError{
			Param: "text",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	reqBody := struct {
		Model string `json:"model"`
		Text  string `json:"text"`
	}{
		Model: model,
		Text:  text,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/tokenizer/count", c.baseURL),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var count TokenCount
	if err := c.do(ctx, req, &count); err != nil {
		return nil, err
	}

	return &count, nil
}

// EstimateTokensFromMessages estimates the number of tokens in a chat completion request
func (c *Client) EstimateTokensFromMessages(ctx context.Context, model string, messages []Message) (*TokenCount, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if model == "" {
		return nil, &errors.InvalidRequestError{
			Param: "model",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	if len(messages) == 0 {
		return nil, &errors.InvalidRequestError{
			Param: "messages",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	// Convert messages to a single string for token counting
	var text string
	for _, msg := range messages {
		text += fmt.Sprintf("<%s>%s</s>\n", msg.Role, msg.Content)
	}

	return c.CountTokens(ctx, model, text)
}

// TokenizationResponse represents the response from the tokenization API
type TokenizationResponse struct {
	Tokens []string `json:"tokens"`
	IDs    []int    `json:"token_ids"`
}

// TokenizeText tokenizes the given text into tokens for a specific model
func (c *Client) TokenizeText(ctx context.Context, model string, text string) (*TokenizationResponse, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if model == "" {
		return nil, &errors.InvalidRequestError{
			Param: "model",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	if text == "" {
		return nil, &errors.InvalidRequestError{
			Param: "text",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	reqBody := struct {
		Model string `json:"model"`
		Text  string `json:"text"`
	}{
		Model: model,
		Text:  text,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/tokenizer/tokenize", c.baseURL),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var tokenization TokenizationResponse
	if err := c.do(ctx, req, &tokenization); err != nil {
		return nil, err
	}

	return &tokenization, nil
}
