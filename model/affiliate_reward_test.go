package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withAffiliateRewardSettings(t *testing.T, trigger string, inviterPercent float64, inviteePercent float64) {
	t.Helper()
	originalQuotaPerUnit := common.QuotaPerUnit
	originalQuotaForInviter := common.QuotaForInviter
	originalQuotaForInvitee := common.QuotaForInvitee
	quotaSetting := operation_setting.GetQuotaSetting()
	originalQuotaSetting := *quotaSetting
	paymentSetting := operation_setting.GetPaymentSetting()
	originalPaymentSetting := *paymentSetting

	t.Cleanup(func() {
		common.QuotaPerUnit = originalQuotaPerUnit
		common.QuotaForInviter = originalQuotaForInviter
		common.QuotaForInvitee = originalQuotaForInvitee
		*quotaSetting = originalQuotaSetting
		*paymentSetting = originalPaymentSetting
	})

	common.QuotaPerUnit = 100
	common.QuotaForInviter = 0
	common.QuotaForInvitee = 0
	quotaSetting.AffiliateRewardTrigger = trigger
	quotaSetting.InviterRegistrationRewardAmount = 0
	quotaSetting.InviteeRegistrationRewardAmount = 0
	quotaSetting.InviterTopupRebatePercent = inviterPercent
	quotaSetting.InviteeTopupRebatePercent = inviteePercent
	paymentSetting.ComplianceConfirmed = true
	paymentSetting.ComplianceTermsVersion = operation_setting.CurrentComplianceTermsVersion
}

func withRegistrationAffiliateRewardSettings(t *testing.T, trigger string, inviterAmount float64, inviteeAmount float64) {
	t.Helper()
	withAffiliateRewardSettings(t, trigger, 0, 0)
	quotaSetting := operation_setting.GetQuotaSetting()
	quotaSetting.InviterRegistrationRewardAmount = inviterAmount
	quotaSetting.InviteeRegistrationRewardAmount = inviteeAmount
}

func insertAffiliateRewardUser(t *testing.T, id int, username string, quota int, inviterId int) {
	t.Helper()
	require.NoError(t, DB.Create(&User{
		Id:        id,
		Username:  username,
		Status:    common.UserStatusEnabled,
		Quota:     quota,
		AffCode:   fmt.Sprintf("aff-%d", id),
		InviterId: inviterId,
		AffCount:  7,
	}).Error)
}

func insertSuccessfulTopUpForAffiliateTest(t *testing.T, tradeNo string, userId int, originalPayAmountUSD float64) *TopUp {
	t.Helper()
	topUp := &TopUp{
		UserId:               userId,
		Amount:               100,
		Money:                999,
		OriginalPayAmountUSD: originalPayAmountUSD,
		TradeNo:              tradeNo,
		PaymentMethod:        PaymentMethodStripe,
		PaymentProvider:      PaymentProviderStripe,
		Status:               common.TopUpStatusSuccess,
		CreateTime:           time.Now().Unix(),
		CompleteTime:         time.Now().Unix(),
	}
	require.NoError(t, topUp.Insert())
	return topUp
}

func requireAffiliateRewardUserState(t *testing.T, id int) User {
	t.Helper()
	var user User
	require.NoError(t, DB.Where("id = ?", id).First(&user).Error)
	return user
}

func countTopUpRebateLogs(t *testing.T) int64 {
	t.Helper()
	var count int64
	require.NoError(t, DB.Model(&TopUpRebateLog{}).Count(&count).Error)
	return count
}

func TestGrantRegistrationAffiliateRewards_RegistrationModeUsesUSDAmounts(t *testing.T) {
	truncateTables(t)
	withRegistrationAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerRegistration, 1.5, 0.75)
	insertAffiliateRewardUser(t, 41, "inviter-registration", 0, 0)
	insertAffiliateRewardUser(t, 42, "invitee-registration", 0, 41)

	grantRegistrationAffiliateRewards(42, 41)

	inviter := requireAffiliateRewardUserState(t, 41)
	invitee := requireAffiliateRewardUserState(t, 42)
	assert.Equal(t, 150, inviter.AffQuota)
	assert.Equal(t, 150, inviter.AffHistoryQuota)
	assert.Equal(t, 8, inviter.AffCount)
	assert.Equal(t, 0, inviter.AffRebateCount)
	assert.Equal(t, 75, invitee.Quota)
}

func TestGrantRegistrationAffiliateRewards_TopupModesOnlyIncrementAffCount(t *testing.T) {
	truncateTables(t)
	withRegistrationAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerFirstTopup, 1.5, 0.75)
	insertAffiliateRewardUser(t, 51, "inviter-first-topup-registration", 0, 0)
	insertAffiliateRewardUser(t, 52, "invitee-first-topup-registration", 0, 51)

	grantRegistrationAffiliateRewards(52, 51)

	inviter := requireAffiliateRewardUserState(t, 51)
	invitee := requireAffiliateRewardUserState(t, 52)
	assert.Equal(t, 0, inviter.AffQuota)
	assert.Equal(t, 0, inviter.AffHistoryQuota)
	assert.Equal(t, 8, inviter.AffCount)
	assert.Equal(t, 0, inviter.AffRebateCount)
	assert.Equal(t, 0, invitee.Quota)
}

func TestAfterTopUpSuccessHook_FirstTopupIsIdempotentAndUsesOriginalPayAmount(t *testing.T) {
	truncateTables(t)
	withAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerFirstTopup, 10, 5)
	insertAffiliateRewardUser(t, 1, "inviter-first", 0, 0)
	insertAffiliateRewardUser(t, 2, "invitee-first", 0, 1)

	firstTopUp := insertSuccessfulTopUpForAffiliateTest(t, "first-topup-rebate", 2, 20)
	require.NoError(t, AfterTopUpSuccessHook(firstTopUp))
	require.NoError(t, AfterTopUpSuccessHook(firstTopUp))

	inviter := requireAffiliateRewardUserState(t, 1)
	invitee := requireAffiliateRewardUserState(t, 2)
	assert.Equal(t, 200, inviter.AffQuota)
	assert.Equal(t, 200, inviter.AffHistoryQuota)
	assert.Equal(t, 1, inviter.AffRebateCount)
	assert.Equal(t, 7, inviter.AffCount)
	assert.Equal(t, 100, invitee.Quota)
	assert.True(t, invitee.HasFirstTopupRebate)
	assert.EqualValues(t, 2, countTopUpRebateLogs(t))

	secondTopUp := insertSuccessfulTopUpForAffiliateTest(t, "second-topup-no-rebate", 2, 30)
	require.NoError(t, AfterTopUpSuccessHook(secondTopUp))

	inviter = requireAffiliateRewardUserState(t, 1)
	invitee = requireAffiliateRewardUserState(t, 2)
	assert.Equal(t, 200, inviter.AffQuota)
	assert.Equal(t, 1, inviter.AffRebateCount)
	assert.Equal(t, 100, invitee.Quota)
	assert.EqualValues(t, 2, countTopUpRebateLogs(t))
}

func TestAfterTopUpSuccessHook_FirstTopupWithZeroPercentConsumesOpportunity(t *testing.T) {
	truncateTables(t)
	withAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerFirstTopup, 0, 0)
	insertAffiliateRewardUser(t, 61, "inviter-zero-first", 0, 0)
	insertAffiliateRewardUser(t, 62, "invitee-zero-first", 0, 61)

	firstTopUp := insertSuccessfulTopUpForAffiliateTest(t, "first-topup-zero-rebate", 62, 20)
	require.NoError(t, AfterTopUpSuccessHook(firstTopUp))

	invitee := requireAffiliateRewardUserState(t, 62)
	assert.True(t, invitee.HasFirstTopupRebate)
	assert.EqualValues(t, 0, countTopUpRebateLogs(t))

	operation_setting.GetQuotaSetting().InviterTopupRebatePercent = 10
	secondTopUp := insertSuccessfulTopUpForAffiliateTest(t, "first-topup-zero-second", 62, 20)
	require.NoError(t, AfterTopUpSuccessHook(secondTopUp))

	inviter := requireAffiliateRewardUserState(t, 61)
	assert.Equal(t, 0, inviter.AffQuota)
	assert.Equal(t, 0, inviter.AffRebateCount)
	assert.EqualValues(t, 0, countTopUpRebateLogs(t))
}

func TestAfterTopUpSuccessHook_EveryTopupRewardsEachOrderWithoutChangingAffCount(t *testing.T) {
	truncateTables(t)
	withAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerEveryTopup, 10, 0)
	insertAffiliateRewardUser(t, 11, "inviter-every", 0, 0)
	insertAffiliateRewardUser(t, 12, "invitee-every", 0, 11)

	firstTopUp := insertSuccessfulTopUpForAffiliateTest(t, "every-topup-one", 12, 10)
	secondTopUp := insertSuccessfulTopUpForAffiliateTest(t, "every-topup-two", 12, 30)
	require.NoError(t, AfterTopUpSuccessHook(firstTopUp))
	require.NoError(t, AfterTopUpSuccessHook(firstTopUp))
	require.NoError(t, AfterTopUpSuccessHook(secondTopUp))

	inviter := requireAffiliateRewardUserState(t, 11)
	assert.Equal(t, 400, inviter.AffQuota)
	assert.Equal(t, 400, inviter.AffHistoryQuota)
	assert.Equal(t, 2, inviter.AffRebateCount)
	assert.Equal(t, 7, inviter.AffCount)
	assert.EqualValues(t, 2, countTopUpRebateLogs(t))
}

func TestTransferAffAmountToQuotaUsesUSDInput(t *testing.T) {
	truncateTables(t)
	withAffiliateRewardSettings(t, operation_setting.AffiliateRewardTriggerEveryTopup, 0, 0)
	insertAffiliateRewardUser(t, 21, "transfer-user", 50, 0)
	require.NoError(t, DB.Model(&User{}).Where("id = ?", 21).Update("aff_quota", 250).Error)

	user := requireAffiliateRewardUserState(t, 21)
	require.NoError(t, user.TransferAffAmountToQuota(1.25))

	user = requireAffiliateRewardUserState(t, 21)
	assert.Equal(t, 175, user.Quota)
	assert.Equal(t, 125, user.AffQuota)
}
