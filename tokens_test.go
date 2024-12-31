package deepseek_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight/deepseek-go"
	"github.com/trustsight/deepseek-go/internal/testutil"
)

func TestEstimateTokenCount(t *testing.T) {
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name     string
		text     string
		minCount int // Since these are estimates, we test for minimum expected tokens
	}{
		{
			name:     "english text",
			text:     "Hello, world!",
			minCount: 3, // At least 3 tokens for "Hello", ",", "world!"
		},
		{
			name:     "chinese text",
			text:     "你好世界",
			minCount: 2, // At least 2 tokens for 4 Chinese characters
		},
		{
			name:     "mixed text",
			text:     "Hello 世界!",
			minCount: 2, // At least 2 tokens
		},
		{
			name:     "empty text",
			text:     "",
			minCount: 1, // Minimum 1 token
		},
		{
			name:     "numbers and symbols",
			text:     "123 !@#",
			minCount: 2, // At least 2 tokens
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := client.EstimateTokenCount(tt.text)
			assert.NotNil(t, estimate)
			assert.GreaterOrEqual(t, estimate.EstimatedTokens, tt.minCount)
		})
	}
}

func TestEstimateTokensFromMessages(t *testing.T) {
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name     string
		messages []deepseek.Message
		minCount int // Since these are estimates, we test for minimum expected tokens
	}{
		{
			name: "single message",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleUser,
					Content: "Hello!",
				},
			},
			minCount: 4, // At least 4 tokens (3 for role + 1 for content)
		},
		{
			name: "multiple messages",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    deepseek.RoleUser,
					Content: "Hi!",
				},
			},
			minCount: 10, // At least 10 tokens (3 per role + content tokens)
		},
		{
			name: "with function call",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleUser,
					Content: "Get weather",
					FunctionCall: &deepseek.FunctionCall{
						Name:      "get_weather",
						Arguments: []byte(`{"location":"New York"}`),
					},
				},
			},
			minCount: 8, // At least 8 tokens (3 for role + content + function name + args)
		},
		{
			name:     "empty messages",
			messages: []deepseek.Message{},
			minCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate := client.EstimateTokensFromMessages(tt.messages)
			assert.NotNil(t, estimate)
			assert.GreaterOrEqual(t, estimate.EstimatedTokens, tt.minCount)
		})
	}
}
