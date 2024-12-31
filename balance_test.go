package deepseek_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight/deepseek-go"
	"github.com/trustsight/deepseek-go/internal/testutil"
)

func TestGetBalance(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	balance, err := client.GetBalance(ctx)
	require.NoError(t, err)
	assert.NotNil(t, balance)

	// Verify balance details
	assert.Equal(t, "balance", balance.Object)
	assert.NotZero(t, balance.TotalBalance)
	assert.NotEmpty(t, balance.Currency)
	assert.NotZero(t, balance.GrantedQuota)
	assert.NotZero(t, balance.UsedQuota)
	assert.NotZero(t, balance.RemainingQuota)
	assert.NotEmpty(t, balance.LastUpdated)

	// Verify quota relationships
	assert.Equal(t, balance.GrantedQuota-balance.UsedQuota, balance.RemainingQuota)
}

func TestGetUsage(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client := deepseek.NewClient(config.APIKey)

	tests := []struct {
		name    string
		params  *deepseek.UsageParams
		wantErr bool
	}{
		{
			name: "get daily usage",
			params: &deepseek.UsageParams{
				StartTime:   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				EndTime:     time.Now().Format(time.RFC3339),
				Granularity: "daily",
			},
			wantErr: false,
		},
		{
			name: "get weekly usage",
			params: &deepseek.UsageParams{
				StartTime:   time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
				EndTime:     time.Now().Format(time.RFC3339),
				Granularity: "daily",
			},
			wantErr: false,
		},
		{
			name:    "nil params",
			params:  nil,
			wantErr: true,
		},
		{
			name: "missing start time",
			params: &deepseek.UsageParams{
				EndTime:     time.Now().Format(time.RFC3339),
				Granularity: "daily",
			},
			wantErr: true,
		},
		{
			name: "missing end time",
			params: &deepseek.UsageParams{
				StartTime:   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				Granularity: "daily",
			},
			wantErr: true,
		},
		{
			name: "invalid granularity",
			params: &deepseek.UsageParams{
				StartTime:   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				EndTime:     time.Now().Format(time.RFC3339),
				Granularity: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
			defer cancel()

			usage, err := client.GetUsage(ctx, tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, usage)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, usage)

			// Verify usage response structure
			assert.Equal(t, "usage", usage.Object)
			assert.NotEmpty(t, usage.StartTime)
			assert.NotEmpty(t, usage.EndTime)
			assert.NotNil(t, usage.Data)

			// Verify totals
			assert.NotZero(t, usage.Total.TotalTokens)
			assert.Equal(t, usage.Total.PromptTokens+usage.Total.CompletionTokens, usage.Total.TotalTokens)

			// Verify individual records
			for _, record := range usage.Data {
				assert.NotEmpty(t, record.Timestamp)
				assert.Equal(t, record.PromptTokens+record.CompletionTokens, record.TotalTokens)
			}
		})
	}
}
