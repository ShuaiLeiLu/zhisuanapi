package model

type TopUpRebateLog struct {
	Id                   int     `json:"id" gorm:"primaryKey;autoIncrement"`
	TradeNo              string  `json:"trade_no" gorm:"type:varchar(255);not null;uniqueIndex:idx_topup_rebate_recipient"`
	TopUpId              int     `json:"topup_id" gorm:"not null;index"`
	InviteeId            int     `json:"invitee_id" gorm:"not null;index"`
	RecipientUserId      int     `json:"recipient_user_id" gorm:"not null;uniqueIndex:idx_topup_rebate_recipient"`
	RecipientRole        string  `json:"recipient_role" gorm:"type:varchar(16);not null;uniqueIndex:idx_topup_rebate_recipient"`
	TriggerMode          string  `json:"trigger_mode" gorm:"type:varchar(32);not null"`
	OriginalPayAmountUSD float64 `json:"original_pay_amount_usd" gorm:"type:decimal(12,6);not null;default:0"`
	RebatePercent        float64 `json:"rebate_percent" gorm:"type:decimal(8,4);not null;default:0"`
	RebateAmountUSD      float64 `json:"rebate_amount_usd" gorm:"type:decimal(12,6);not null;default:0"`
	RebateQuota          int     `json:"rebate_quota" gorm:"not null;default:0"`
	CreatedAt            int64   `json:"created_at" gorm:"autoCreateTime;column:created_at"`
}

func (TopUpRebateLog) TableName() string {
	return "topup_rebate_logs"
}
