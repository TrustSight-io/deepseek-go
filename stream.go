package deepseek

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/trustsight-io/deepseek-go/internal/errors"
)

// Stream represents a streaming response from the API
type Stream struct {
	reader    *bufio.Reader
	response  *http.Response
	errChan   chan error
	done      bool
	closeOnce chan struct{}
}

// StreamChoice represents a choice in a streaming response
type StreamChoice struct {
	Delta struct {
		Content string `json:"content,omitempty"`
		Role    string `json:"role,omitempty"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// StreamResponse represents a streamed response chunk
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// newStream creates a new Stream from an HTTP response
func newStream(resp *http.Response) *Stream {
	return &Stream{
		reader:    bufio.NewReader(resp.Body),
		response:  resp,
		errChan:   make(chan error, 1),
		closeOnce: make(chan struct{}),
	}
}

// Recv receives the next chunk of data from the stream
func (s *Stream) Recv() (*StreamResponse, error) {
	if s.done {
		return nil, io.EOF
	}

	line, err := s.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			s.done = true
			return nil, io.EOF
		}
		return nil, fmt.Errorf("error reading from stream: %w", err)
	}

	// Remove "data: " prefix
	line = bytes.TrimPrefix(line, []byte("data: "))
	line = bytes.TrimSpace(line)

	// Skip empty lines
	if len(line) == 0 {
		return s.Recv()
	}

	// Check for stream end
	if bytes.Equal(line, []byte("[DONE]")) {
		s.done = true
		return nil, io.EOF
	}

	var response StreamResponse
	if err := json.Unmarshal(line, &response); err != nil {
		return nil, &errors.InvalidRequestError{
			Err: fmt.Errorf("failed to decode stream response: %w", err),
		}
	}

	return &response, nil
}

// Close closes the stream
func (s *Stream) Close() error {
	select {
	case <-s.closeOnce:
		return nil
	default:
		close(s.closeOnce)
		if s.response != nil && s.response.Body != nil {
			return s.response.Body.Close()
		}
		return nil
	}
}

// CreateChatCompletionStream creates a streaming chat completion request
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	req *ChatCompletionRequest,
) (*Stream, error) {
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

	req.Stream = true
	if req.Model == "" {
		req.Model = "deepseek-chat"
	}

	httpReq, err := c.newRequest(ctx, http.MethodPost, "/chat/completions", req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				err = fmt.Errorf("error closing response body: %v", cerr)
			}
		}()
		var apiErr errors.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, errors.HandleErrorResp(resp, &apiErr)
	}

	return newStream(resp), nil
}

// ContentAccumulator helps accumulate streamed content
type ContentAccumulator struct {
	buffer strings.Builder
}

// Add adds content to the accumulator
func (ca *ContentAccumulator) Add(content string) {
	ca.buffer.WriteString(content)
}

// String returns the accumulated content
func (ca *ContentAccumulator) String() string {
	return ca.buffer.String()
}

// Reset clears the accumulated content
func (ca *ContentAccumulator) Reset() {
	ca.buffer.Reset()
}

// CollectFullResponse collects a complete response from a stream
func CollectFullResponse(stream *Stream) (response string, err error) {
	defer func() {
		if cerr := stream.Close(); cerr != nil {
			err = fmt.Errorf("error closing stream: %v", cerr)
		}
	}()

	var accumulator ContentAccumulator

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if len(response.Choices) > 0 {
			accumulator.Add(response.Choices[0].Delta.Content)
		}
	}

	return accumulator.String(), nil
}
