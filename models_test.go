package deepseek_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight-io/deepseek-go"
	"github.com/trustsight-io/deepseek-go/internal/testutil"
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
	assert.Equal(t, "list", resp.Object)
	assert.NotEmpty(t, resp.Data)

	// Verify model details
	for _, model := range resp.Data {
		assert.NotEmpty(t, model.ID)
		assert.Equal(t, "model", model.Object)
		assert.Equal(t, "deepseek", model.OwnedBy)
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
			assert.Equal(t, "model", model.Object)
			assert.Equal(t, "deepseek", model.OwnedBy)
		})
	}
}
