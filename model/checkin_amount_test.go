package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withCheckinAmountSettings(t *testing.T, minAmount float64, maxAmount float64) {
	t.Helper()
	originalQuotaPerUnit := common.QuotaPerUnit
	setting := operation_setting.GetCheckinSetting()
	originalSetting := *setting

	t.Cleanup(func() {
		common.QuotaPerUnit = originalQuotaPerUnit
		*setting = originalSetting
	})

	common.QuotaPerUnit = 100
	setting.Enabled = true
	setting.MinQuota = 1
	setting.MaxQuota = 1
	setting.MinAmount = minAmount
	setting.MaxAmount = maxAmount
}

func TestUserCheckinUsesAmountRangeAndReturnsAmountFields(t *testing.T) {
	truncateTables(t)
	withCheckinAmountSettings(t, 1.25, 1.25)
	insertAffiliateRewardUser(t, 31, "checkin-user", 0, 0)

	checkin, err := UserCheckin(31)
	require.NoError(t, err)
	assert.Equal(t, 125, checkin.QuotaAwarded)
	assert.Equal(t, 1.25, checkin.QuotaAwardedAmount)

	stats, err := GetUserCheckinStats(31, checkin.CheckinDate[:7])
	require.NoError(t, err)
	assert.EqualValues(t, 125, stats["total_quota"])
	assert.Equal(t, 1.25, stats["total_quota_amount"])

	records, ok := stats["records"].([]CheckinRecord)
	require.True(t, ok)
	require.Len(t, records, 1)
	assert.Equal(t, 125, records[0].QuotaAwarded)
	assert.Equal(t, 1.25, records[0].QuotaAwardedAmount)
}
