package operation_setting

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
)

const (
	AffiliateRewardTriggerRegistration = "registration"
	AffiliateRewardTriggerFirstTopup   = "first_topup"
	AffiliateRewardTriggerEveryTopup   = "every_topup"
)

type QuotaSetting struct {
	EnableFreeModelPreConsume       bool    `json:"enable_free_model_pre_consume"` // 是否对免费模型启用预消耗
	AffiliateRewardTrigger          string  `json:"affiliate_reward_trigger"`
	InviterRegistrationRewardAmount float64 `json:"inviter_registration_reward_amount"`
	InviteeRegistrationRewardAmount float64 `json:"invitee_registration_reward_amount"`
	InviterTopupRebatePercent       float64 `json:"inviter_topup_rebate_percent"`
	InviteeTopupRebatePercent       float64 `json:"invitee_topup_rebate_percent"`
}

// 默认配置
var quotaSetting = QuotaSetting{
	EnableFreeModelPreConsume: true,
	AffiliateRewardTrigger:    AffiliateRewardTriggerRegistration,
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("quota_setting", &quotaSetting)
}

func GetQuotaSetting() *QuotaSetting {
	return &quotaSetting
}

func NormalizeAffiliateRewardTrigger(trigger string) string {
	switch trigger {
	case AffiliateRewardTriggerRegistration, AffiliateRewardTriggerFirstTopup, AffiliateRewardTriggerEveryTopup:
		return trigger
	default:
		return AffiliateRewardTriggerRegistration
	}
}

func GetAffiliateRewardTrigger() string {
	return NormalizeAffiliateRewardTrigger(quotaSetting.AffiliateRewardTrigger)
}

func registrationRewardQuota(amount float64, legacyQuota int) int {
	if amount > 0 {
		return common.AmountToQuota(amount)
	}
	if legacyQuota > 0 {
		return legacyQuota
	}
	return 0
}

func GetInviterRegistrationRewardQuota() int {
	return registrationRewardQuota(quotaSetting.InviterRegistrationRewardAmount, common.QuotaForInviter)
}

func GetInviteeRegistrationRewardQuota() int {
	return registrationRewardQuota(quotaSetting.InviteeRegistrationRewardAmount, common.QuotaForInvitee)
}

func GetInviterTopupRebatePercent() float64 {
	if quotaSetting.InviterTopupRebatePercent < 0 {
		return 0
	}
	return quotaSetting.InviterTopupRebatePercent
}

func GetInviteeTopupRebatePercent() float64 {
	if quotaSetting.InviteeTopupRebatePercent < 0 {
		return 0
	}
	return quotaSetting.InviteeTopupRebatePercent
}
