package deepseek_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight-io/deepseek-go"
	"github.com/trustsight-io/deepseek-go/internal/testutil"
)

func TestNewClient(t *testing.T) {
	config := testutil.LoadTestConfig(t)

	tests := []struct {
		name    string
		apiKey  string
		opts    []deepseek.ClientOption
		wantErr bool
	}{
		{
			name:    "valid configuration",
			apiKey:  config.APIKey,
			opts:    nil,
			wantErr: false,
		},
		{
			name:   "with custom HTTP client",
			apiKey: config.APIKey,
			opts: []deepseek.ClientOption{
				deepseek.WithHTTPClient(&http.Client{
					Timeout: time.Minute,
				}),
			},
			wantErr: false,
		},
		{
			name:   "with max retries",
			apiKey: config.APIKey,
			opts: []deepseek.ClientOption{
				deepseek.WithMaxRetries(3),
			},
			wantErr: false,
		},
		{
			name:   "with retry wait time",
			apiKey: config.APIKey,
			opts: []deepseek.ClientOption{
				deepseek.WithRetryWaitTime(time.Second),
			},
			wantErr: false,
		},
		{
			name:   "with max request size",
			apiKey: config.APIKey,
			opts: []deepseek.ClientOption{
				deepseek.WithMaxRequestSize(1 << 20), // 1MB
			},
			wantErr: false,
		},
		{
			name:   "with debug mode",
			apiKey: config.APIKey,
			opts: []deepseek.ClientOption{
				deepseek.WithDebug(true),
			},
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			opts:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := deepseek.NewClient(tt.apiKey, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, client)

			// Test a simple API call to verify configuration
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			models, err := client.ListModels(ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, models)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, models)
			}
		})
	}
}

func TestClientTimeout(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)

	client, err := deepseek.NewClient(config.APIKey,
		deepseek.WithHTTPClient(&http.Client{
			Timeout: 1 * time.Millisecond, // Very short timeout
		}),
	)
	require.NoError(t, err)

	ctx := context.Background()
	_, respErr := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.Message{
			{
				Role:    deepseek.RoleUser,
				Content: "Hello!",
			},
		},
	})

	assert.Error(t, respErr)
	assert.Contains(t, respErr.Error(), "context deadline exceeded")
}

func TestClientContextCancellation(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, respErr := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.Message{
			{
				Role:    deepseek.RoleUser,
				Content: "Hello!",
			},
		},
	})

	assert.Error(t, respErr)
	assert.Contains(t, respErr.Error(), "context canceled")
}

func TestClientRetry(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)

	client, err := deepseek.NewClient(config.APIKey,
		deepseek.WithMaxRetries(2),
		deepseek.WithRetryWaitTime(time.Second),
	)
	require.NoError(t, err)

	ctx := context.Background()
	_, respErr := client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: "nonexistent-model", // This should trigger retries
		Messages: []deepseek.Message{
			{
				Role:    deepseek.RoleUser,
				Content: "Hello!",
			},
		},
	})

	assert.Error(t, respErr)
	// The error should indicate we tried multiple times
	assert.Contains(t, respErr.Error(), "400")
}
