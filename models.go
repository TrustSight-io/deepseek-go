package deepseek

import (
	"context"
	"fmt"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// Model represents a model's information
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// ModelList represents the response from the models endpoint
type ModelList struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// ListModels lists all available models
func (c *Client) ListModels(ctx context.Context) (*ModelList, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/models", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var models ModelList
	if err := c.do(ctx, req, &models); err != nil {
		return nil, err
	}

	return &models, nil
}

// GetModel retrieves information about a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*Model, error) {
	if ctx == nil {
		return nil, &errors.InvalidRequestError{
			Param: "context",
			Err:   fmt.Errorf("cannot be nil"),
		}
	}

	if modelID == "" {
		return nil, &errors.InvalidRequestError{
			Param: "model_id",
			Err:   fmt.Errorf("cannot be empty"),
		}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/models/"+modelID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	var model Model
	if err := c.do(ctx, req, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
