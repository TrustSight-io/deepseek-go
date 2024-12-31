package deepseek

import (
	"context"
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
