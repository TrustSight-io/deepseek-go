package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// CompletionRequest represents a request for text completion.
type CompletionRequest struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	MaxTokens        int      `json:"max_tokens,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	TopP             float32  `json:"top_p,omitempty"`
	N                int      `json:"n,omitempty"`
	Stream           bool     `json:"stream,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	JSONMode         bool     `json:"json_mode,omitempty"`
}

// CompletionChoice represents a completion choice.
type CompletionChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

// CompletionResponse represents a response from the completion API.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// CreateCompletion creates a completion.
func (c *Client) CreateCompletion(
	ctx context.Context,
	request *CompletionRequest,
) (*CompletionResponse, error) {
	if request == nil {
		return nil, &errors.InvalidRequestError{Param: "request", Err: fmt.Errorf("cannot be nil")}
	}

	if request.Prompt == "" {
		return nil, &errors.InvalidRequestError{Param: "prompt", Err: fmt.Errorf("cannot be empty")}
	}

	if request.Model == "" {
		request.Model = "deepseek-coder"
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/completions", request)
	if err != nil {
		return nil, err
	}

	var response CompletionResponse
	if err := c.do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateCompletionStream creates a streaming completion.
func (c *Client) CreateCompletionStream(
	ctx context.Context,
	request *CompletionRequest,
) (*CompletionStreamReader, error) {
	if request == nil {
		return nil, &errors.InvalidRequestError{Param: "request", Err: fmt.Errorf("cannot be nil")}
	}

	if request.Prompt == "" {
		return nil, &errors.InvalidRequestError{Param: "prompt", Err: fmt.Errorf("cannot be empty")}
	}

	request.Stream = true
	if request.Model == "" {
		request.Model = "deepseek-coder"
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/completions", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.HandleErrorResp(resp, &errors.APIError{StatusCode: resp.StatusCode})
	}

	return &CompletionStreamReader{
		ChatCompletionStreamReader: ChatCompletionStreamReader{
			reader:     resp.Body,
			response:   resp,
			delimiter:  []byte("\n"),
			remaining:  nil,
			finished:   false,
			httpClient: c.httpClient,
		},
	}, nil
}

// CompletionStreamReader handles streaming responses from the completion API.
type CompletionStreamReader struct {
	ChatCompletionStreamReader
}

// Recv receives the next message from the stream.
func (s *CompletionStreamReader) Recv() (*CompletionResponse, error) {
	chatResp, err := s.ChatCompletionStreamReader.Recv()
	if err != nil {
		return nil, err
	}

	resp := &CompletionResponse{
		ID:      chatResp.ID,
		Object:  chatResp.Object,
		Created: chatResp.Created,
		Model:   chatResp.Model,
		Choices: make([]CompletionChoice, len(chatResp.Choices)),
	}

	for i, choice := range chatResp.Choices {
		resp.Choices[i] = CompletionChoice{
			Text:         choice.Delta.Content,
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
		}
	}

	return resp, nil
}
