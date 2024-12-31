package deepseek

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/trustsight/deepseek-go/internal/errors"
)

// ChatCompletionStreamReader handles streaming responses from the chat completion API.
type ChatCompletionStreamReader struct {
	reader     io.ReadCloser
	response   *http.Response
	delimiter  []byte
	remaining  []byte
	finished   bool
	httpClient *http.Client
}

// Recv receives the next message from the stream.
func (s *ChatCompletionStreamReader) Recv() (*ChatCompletionStream, error) {
	if s.finished {
		return nil, io.EOF
	}

	chunk, err := s.readNext()
	if err != nil {
		return nil, err
	}

	var response ChatCompletionStream
	if err := json.Unmarshal(chunk, &response); err != nil {
		return nil, &errors.InvalidRequestError{Err: fmt.Errorf("failed to decode stream response: %v", err)}
	}

	if len(response.Choices) > 0 && response.Choices[0].FinishReason != "" {
		s.finished = true
	}

	return &response, nil
}

// Close closes the response body.
func (s *ChatCompletionStreamReader) Close() error {
	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
}

// readNext reads the next chunk from the stream.
func (s *ChatCompletionStreamReader) readNext() ([]byte, error) {
	reader := bufio.NewReader(s.reader)
	var buffer bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				s.finished = true
				if buffer.Len() > 0 {
					return buffer.Bytes(), nil
				}
				return nil, io.EOF
			}
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("data: ")) {
			line = bytes.TrimPrefix(line, []byte("data: "))
			if bytes.Equal(line, []byte("[DONE]")) {
				s.finished = true
				if buffer.Len() > 0 {
					return buffer.Bytes(), nil
				}
				return nil, io.EOF
			}
		}

		buffer.Write(line)
		if bytes.HasSuffix(line, s.delimiter) {
			return buffer.Bytes(), nil
		}
	}
}
