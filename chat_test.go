package deepseek_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight/deepseek-go"
	"github.com/trustsight/deepseek-go/internal/testutil"
)

func TestCreateChatCompletion(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *deepseek.ChatCompletionRequest
		wantErr bool
	}{
		{
			name: "successful completion",
			req: &deepseek.ChatCompletionRequest{
				Model: "deepseek-chat",
				Messages: []deepseek.Message{
					{
						Role:    deepseek.RoleUser,
						Content: "Say hello!",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with function calling",
			req: &deepseek.ChatCompletionRequest{
				Model: "deepseek-chat",
				Messages: []deepseek.Message{
					{
						Role:    deepseek.RoleUser,
						Content: "What's the weather in New York?",
					},
				},
				Functions: []deepseek.Function{
					{
						Name:        "get_weather",
						Description: "Get the current weather in a location",
						Parameters: []byte(`{
							"type": "object",
							"properties": {
								"location": {
									"type": "string",
									"description": "The city and state"
								}
							},
							"required": ["location"]
						}`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty messages",
			req: &deepseek.ChatCompletionRequest{
				Model:    "deepseek-chat",
				Messages: []deepseek.Message{},
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			resp, err := client.CreateChatCompletion(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.ID)
			assert.NotEmpty(t, resp.Created)
			assert.NotEmpty(t, resp.Model)
			assert.NotEmpty(t, resp.Choices)
			assert.Greater(t, len(resp.Choices), 0)
			assert.NotEmpty(t, resp.Choices[0].Message.Content)
			assert.NotEmpty(t, resp.Usage.TotalTokens)
		})
	}
}

func TestCreateChatCompletionStream(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, &deepseek.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.Message{
			{
				Role:    deepseek.RoleUser,
				Content: "Write a short story about a robot learning to paint.",
			},
		},
		Stream: true,
	})
	require.NoError(t, err)
	defer stream.Close()

	var collected []string
	var receivedRole bool
	var receivedContent bool

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		if len(resp.Choices) > 0 {
			if resp.Choices[0].Delta.Role != "" {
				receivedRole = true
			}
			if resp.Choices[0].Delta.Content != "" {
				receivedContent = true
				collected = append(collected, resp.Choices[0].Delta.Content)
			}
		}
	}

	assert.True(t, receivedRole, "should receive role in stream")
	assert.True(t, receivedContent, "should receive content in stream")
	assert.NotEmpty(t, collected, "should collect streamed content")

	// Join the collected content and verify it forms a coherent story
	story := ""
	for _, chunk := range collected {
		story += chunk
	}
	assert.NotEmpty(t, story)
	assert.Contains(t, story, "robot")
	assert.Contains(t, story, "paint")
}

func TestCreateChatCompletionStreamErrors(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *deepseek.ChatCompletionRequest
		wantErr bool
	}{
		{
			name: "empty messages",
			req: &deepseek.ChatCompletionRequest{
				Model:    "deepseek-chat",
				Messages: []deepseek.Message{},
				Stream:   true,
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid model",
			req: &deepseek.ChatCompletionRequest{
				Model: "nonexistent-model",
				Messages: []deepseek.Message{
					{
						Role:    deepseek.RoleUser,
						Content: "Hello!",
					},
				},
				Stream: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			stream, err := client.CreateChatCompletionStream(ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				if stream != nil {
					stream.Close()
				}
				return
			}

			require.NoError(t, err)
			defer stream.Close()

			resp, err := stream.Recv()
			require.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
