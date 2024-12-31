package deepseek_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight/deepseek-go"
	"github.com/trustsight/deepseek-go/internal/testutil"
)

func TestListModels(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	resp, err := client.ListModels(ctx)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify response structure
	assert.NotEmpty(t, resp.Object)
	assert.NotEmpty(t, resp.Data)

	// Verify model details
	for _, model := range resp.Data {
		assert.NotEmpty(t, model.ID)
		assert.NotEmpty(t, model.Object)
		assert.NotZero(t, model.Created)
		assert.NotEmpty(t, model.OwnedBy)
		assert.NotZero(t, model.ContextWindow)
		assert.NotEmpty(t, model.Capabilities)
		assert.NotZero(t, model.MaxTokens)

		// Verify pricing if available
		if model.PricingConfig != nil {
			assert.NotZero(t, model.PricingConfig.PromptTokenPrice)
			assert.NotZero(t, model.PricingConfig.CompletionTokenPrice)
			assert.NotEmpty(t, model.PricingConfig.Currency)
		}
	}
}

func TestGetModel(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name    string
		modelID string
		wantErr bool
	}{
		{
			name:    "get deepseek-chat",
			modelID: "deepseek-chat",
			wantErr: false,
		},
		{
			name:    "empty model ID",
			modelID: "",
			wantErr: true,
		},
		{
			name:    "nonexistent model",
			modelID: "nonexistent-model",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			model, err := client.GetModel(ctx, tt.modelID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, model)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, model)
			assert.Equal(t, tt.modelID, model.ID)
			assert.NotEmpty(t, model.Object)
			assert.NotZero(t, model.Created)
			assert.NotEmpty(t, model.OwnedBy)
			assert.NotZero(t, model.ContextWindow)
			assert.NotEmpty(t, model.Capabilities)
			assert.NotZero(t, model.MaxTokens)
		})
	}
}

func TestGetModelConfig(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	tests := []struct {
		name    string
		modelID string
		wantErr bool
	}{
		{
			name:    "get deepseek-chat config",
			modelID: "deepseek-chat",
			wantErr: false,
		},
		{
			name:    "empty model ID",
			modelID: "",
			wantErr: true,
		},
		{
			name:    "nonexistent model",
			modelID: "nonexistent-model",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			config, err := client.GetModelConfig(ctx, tt.modelID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, config)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, config)
			assert.NotZero(t, config.ContextWindow)
			assert.NotZero(t, config.MaxTokens)
			assert.NotEmpty(t, config.SupportedFeatures)
			assert.NotEmpty(t, config.DefaultParameters)
		})
	}
}
