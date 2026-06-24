package operation_setting

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
)

// CheckinSetting 签到功能配置
type CheckinSetting struct {
	Enabled   bool    `json:"enabled"`    // 是否启用签到功能
	MinQuota  int     `json:"min_quota"`  // 签到最小额度奖励
	MaxQuota  int     `json:"max_quota"`  // 签到最大额度奖励
	MinAmount float64 `json:"min_amount"` // 签到最小金额奖励，单位 USD
	MaxAmount float64 `json:"max_amount"` // 签到最大金额奖励，单位 USD
}

// 默认配置
var checkinSetting = CheckinSetting{
	Enabled:  false, // 默认关闭
	MinQuota: 1000,  // 默认最小额度 1000 (约 0.002 USD)
	MaxQuota: 10000, // 默认最大额度 10000 (约 0.02 USD)
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("checkin_setting", &checkinSetting)
}

// GetCheckinSetting 获取签到配置
func GetCheckinSetting() *CheckinSetting {
	return &checkinSetting
}

// IsCheckinEnabled 是否启用签到功能
func IsCheckinEnabled() bool {
	return checkinSetting.Enabled
}

// GetCheckinQuotaRange 获取签到额度范围
func GetCheckinQuotaRange() (min, max int) {
	return checkinSetting.GetQuotaRange()
}

func (setting *CheckinSetting) GetQuotaRange() (min, max int) {
	if setting.MinAmount > 0 || setting.MaxAmount > 0 {
		min = common.AmountToQuota(setting.MinAmount)
		max = common.AmountToQuota(setting.MaxAmount)
		if max == 0 {
			max = min
		}
		if max < min {
			max = min
		}
		return min, max
	}
	min = setting.MinQuota
	max = setting.MaxQuota
	if max < min {
		max = min
	}
	return min, max
}

func (setting *CheckinSetting) GetAmountRange() (min, max float64) {
	if setting.MinAmount > 0 || setting.MaxAmount > 0 {
		min = setting.MinAmount
		max = setting.MaxAmount
		if max == 0 {
			max = min
		}
		if max < min {
			max = min
		}
		return min, max
	}
	min = common.QuotaToAmount(setting.MinQuota)
	max = common.QuotaToAmount(setting.MaxQuota)
	if max < min {
		max = min
	}
	return min, max
}
