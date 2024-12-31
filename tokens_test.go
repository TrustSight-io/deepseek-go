package deepseek_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight/deepseek-go"
	"github.com/trustsight/deepseek-go/internal/testutil"
)

func TestCountTokens(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name    string
		model   string
		text    string
		wantErr bool
	}{
		{
			name:    "basic text",
			model:   "deepseek-chat",
			text:    "Hello, world!",
			wantErr: false,
		},
		{
			name:    "empty text",
			model:   "deepseek-chat",
			text:    "",
			wantErr: true,
		},
		{
			name:    "empty model",
			model:   "",
			text:    "Hello, world!",
			wantErr: true,
		},
		{
			name:    "long text",
			model:   "deepseek-chat",
			text:    string(make([]byte, 10000)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			count, err := client.CountTokens(ctx, tt.model, tt.text)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, count)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, count)
			assert.NotZero(t, count.TotalTokens)
			assert.NotZero(t, count.Details.Prompt)
			if tt.name == "long text" {
				assert.True(t, count.Details.Truncated)
			}
		})
	}
}

func TestEstimateTokensFromMessages(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name     string
		model    string
		messages []deepseek.Message
		wantErr  bool
	}{
		{
			name:  "single message",
			model: "deepseek-chat",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleUser,
					Content: "Hello!",
				},
			},
			wantErr: false,
		},
		{
			name:     "empty messages",
			model:    "deepseek-chat",
			messages: []deepseek.Message{},
			wantErr:  true,
		},
		{
			name:  "empty model",
			model: "",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleUser,
					Content: "Hello!",
				},
			},
			wantErr: true,
		},
		{
			name:  "multiple messages",
			model: "deepseek-chat",
			messages: []deepseek.Message{
				{
					Role:    deepseek.RoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    deepseek.RoleUser,
					Content: "Hello!",
				},
				{
					Role:    deepseek.RoleAssistant,
					Content: "Hi there! How can I help you today?",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			count, err := client.EstimateTokensFromMessages(ctx, tt.model, tt.messages)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, count)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, count)
			assert.NotZero(t, count.TotalTokens)
			assert.NotZero(t, count.Details.Prompt)
		})
	}
}

func TestTokenizeText(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name    string
		model   string
		text    string
		wantErr bool
	}{
		{
			name:    "basic text",
			model:   "deepseek-chat",
			text:    "Hello, world!",
			wantErr: false,
		},
		{
			name:    "empty text",
			model:   "deepseek-chat",
			text:    "",
			wantErr: true,
		},
		{
			name:    "empty model",
			model:   "",
			text:    "Hello, world!",
			wantErr: true,
		},
		{
			name:    "special characters",
			model:   "deepseek-chat",
			text:    "Hello 👋 World! 🌍",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			tokenization, err := client.TokenizeText(ctx, tt.model, tt.text)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tokenization)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, tokenization)
			assert.NotEmpty(t, tokenization.Tokens)
			assert.NotEmpty(t, tokenization.IDs)
			assert.Equal(t, len(tokenization.Tokens), len(tokenization.IDs))

			// Verify token IDs are valid
			for _, id := range tokenization.IDs {
				assert.NotZero(t, id)
			}
		})
	}
}
