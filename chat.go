package deepseek

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// CreateChatCompletion creates a chat completion.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request *ChatCompletionRequest,
) (*ChatCompletionResponse, error) {
	if request == nil {
		return nil, &errors.InvalidRequestError{Param: "request", Err: fmt.Errorf("cannot be nil")}
	}

	if len(request.Messages) == 0 {
		return nil, &errors.InvalidRequestError{Param: "messages", Err: fmt.Errorf("cannot be empty")}
	}

	if request.Model == "" {
		request.Model = "deepseek-chat"
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/chat/completions", request)
	if err != nil {
		return nil, err
	}

	var response ChatCompletionResponse
	if err := c.do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateChatCompletionStream creates a streaming chat completion.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request *ChatCompletionRequest,
) (*ChatCompletionStreamReader, error) {
	if request == nil {
		return nil, &errors.InvalidRequestError{Param: "request", Err: fmt.Errorf("cannot be nil")}
	}

	if len(request.Messages) == 0 {
		return nil, &errors.InvalidRequestError{Param: "messages", Err: fmt.Errorf("cannot be empty")}
	}

	request.Stream = true
	if request.Model == "" {
		request.Model = "deepseek-chat"
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/chat/completions", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		var apiErr errors.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("deepseek: failed to decode error response: %v", err)
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, errors.HandleErrorResp(resp, &apiErr)
	}

	return &ChatCompletionStreamReader{
		reader:     resp.Body,
		response:   resp,
		delimiter:  []byte("\n"),
		remaining:  nil,
		finished:   false,
		httpClient: c.httpClient,
	}, nil
}
