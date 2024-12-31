package deepseek_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustsight-io/deepseek-go"
	"github.com/trustsight-io/deepseek-go/internal/testutil"
)

func TestGetBalance(t *testing.T) {
	testutil.SkipIfShort(t)
	config := testutil.LoadTestConfig(t)
	client, err := deepseek.NewClient(config.APIKey)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	balance, err := client.GetBalance(ctx)
	require.NoError(t, err)
	assert.NotNil(t, balance)

	// Verify balance details
	assert.True(t, balance.IsAvailable)
	assert.NotEmpty(t, balance.BalanceInfos)

	// Verify balance info details
	for _, info := range balance.BalanceInfos {
		assert.NotEmpty(t, info.Currency)
		assert.NotEmpty(t, info.TotalBalance)
		assert.NotEmpty(t, info.GrantedBalance)
		assert.NotEmpty(t, info.ToppedUpBalance)
	}
}
