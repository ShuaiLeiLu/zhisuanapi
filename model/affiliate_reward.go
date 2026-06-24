package model

import (
	"errors"
	"fmt"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	topupRebateRecipientRoleInviter = "inviter"
	topupRebateRecipientRoleInvitee = "invitee"
)

type creditedTopUpRebate struct {
	userId    int
	quota     int
	amountUSD float64
	role      string
}

func formatUSDAmount(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

func creditInviterAffQuotaTx(tx *gorm.DB, inviterId int, quota int, incrementRebateCount bool) error {
	if quota <= 0 {
		return nil
	}
	updates := map[string]interface{}{
		"aff_quota":   gorm.Expr("aff_quota + ?", quota),
		"aff_history": gorm.Expr("aff_history + ?", quota),
	}
	if incrementRebateCount {
		updates["aff_rebate_count"] = gorm.Expr("aff_rebate_count + ?", 1)
	}
	result := tx.Model(&User{}).Where("id = ?", inviterId).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func incrementInviterAffCount(inviterId int) error {
	if inviterId == 0 {
		return nil
	}
	result := DB.Model(&User{}).Where("id = ?", inviterId).
		Update("aff_count", gorm.Expr("aff_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func grantRegistrationAffiliateRewards(userId int, inviterId int) {
	if inviterId == 0 {
		return
	}
	if err := incrementInviterAffCount(inviterId); err != nil {
		common.SysLog(fmt.Sprintf("failed to increment inviter aff_count: inviter_id=%d invitee_id=%d error=%v", inviterId, userId, err))
		return
	}
	if !operation_setting.IsPaymentComplianceConfirmed() {
		return
	}
	if operation_setting.GetAffiliateRewardTrigger() != operation_setting.AffiliateRewardTriggerRegistration {
		return
	}

	inviteeQuota := operation_setting.GetInviteeRegistrationRewardQuota()
	if inviteeQuota > 0 {
		if err := IncreaseUserQuota(userId, inviteeQuota, true); err != nil {
			common.SysLog(fmt.Sprintf("failed to grant invitee registration reward: user_id=%d quota=%d error=%v", userId, inviteeQuota, err))
		} else {
			RecordLog(userId, LogTypeSystem, fmt.Sprintf("使用邀请码赠送 %s", formatUSDAmount(common.QuotaToAmount(inviteeQuota))))
		}
	}

	inviterQuota := operation_setting.GetInviterRegistrationRewardQuota()
	if inviterQuota <= 0 {
		return
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		return creditInviterAffQuotaTx(tx, inviterId, inviterQuota, false)
	})
	if err != nil {
		common.SysLog(fmt.Sprintf("failed to grant inviter registration reward: inviter_id=%d invitee_id=%d quota=%d error=%v", inviterId, userId, inviterQuota, err))
		return
	}
	RecordLog(inviterId, LogTypeSystem, fmt.Sprintf("邀请用户赠送 %s", formatUSDAmount(common.QuotaToAmount(inviterQuota))))
}

func rebateAmountToQuota(originalPayAmountUSD float64, percent float64) (float64, int) {
	if originalPayAmountUSD <= 0 || percent <= 0 || common.QuotaPerUnit <= 0 {
		return 0, 0
	}
	amount := decimal.NewFromFloat(originalPayAmountUSD).
		Mul(decimal.NewFromFloat(percent)).
		Div(decimal.NewFromInt(100))
	quota := amount.
		Mul(decimal.NewFromFloat(common.QuotaPerUnit)).
		Round(0).
		IntPart()
	return amount.InexactFloat64(), int(quota)
}

func createTopUpRebateTx(
	tx *gorm.DB,
	topUp *TopUp,
	inviteeId int,
	recipientUserId int,
	recipientRole string,
	triggerMode string,
	percent float64,
) (*creditedTopUpRebate, error) {
	amountUSD, quota := rebateAmountToQuota(topUp.OriginalPayAmountUSD, percent)
	if quota <= 0 {
		return nil, nil
	}

	var recipientCount int64
	if err := tx.Model(&User{}).Where("id = ?", recipientUserId).Count(&recipientCount).Error; err != nil {
		return nil, err
	}
	if recipientCount == 0 {
		return nil, nil
	}

	log := &TopUpRebateLog{
		TradeNo:              topUp.TradeNo,
		TopUpId:              topUp.Id,
		InviteeId:            inviteeId,
		RecipientUserId:      recipientUserId,
		RecipientRole:        recipientRole,
		TriggerMode:          triggerMode,
		OriginalPayAmountUSD: topUp.OriginalPayAmountUSD,
		RebatePercent:        percent,
		RebateAmountUSD:      amountUSD,
		RebateQuota:          quota,
	}
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(log)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	switch recipientRole {
	case topupRebateRecipientRoleInviter:
		if err := creditInviterAffQuotaTx(tx, recipientUserId, quota, true); err != nil {
			return nil, err
		}
	case topupRebateRecipientRoleInvitee:
		result := tx.Model(&User{}).Where("id = ?", recipientUserId).
			Update("quota", gorm.Expr("quota + ?", quota))
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}
	default:
		return nil, errors.New("invalid topup rebate recipient role")
	}

	return &creditedTopUpRebate{
		userId:    recipientUserId,
		quota:     quota,
		amountUSD: amountUSD,
		role:      recipientRole,
	}, nil
}

func AfterTopUpSuccessHook(topUp *TopUp) error {
	if topUp == nil || topUp.Status != common.TopUpStatusSuccess {
		return nil
	}
	if topUp.OriginalPayAmountUSD <= 0 {
		return nil
	}
	if !operation_setting.IsPaymentComplianceConfirmed() {
		return nil
	}

	triggerMode := operation_setting.GetAffiliateRewardTrigger()
	if triggerMode != operation_setting.AffiliateRewardTriggerFirstTopup &&
		triggerMode != operation_setting.AffiliateRewardTriggerEveryTopup {
		return nil
	}

	credited := make([]creditedTopUpRebate, 0, 2)
	err := DB.Transaction(func(tx *gorm.DB) error {
		var invitee User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", topUp.UserId).First(&invitee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if invitee.InviterId == 0 {
			return nil
		}
		if triggerMode == operation_setting.AffiliateRewardTriggerFirstTopup && invitee.HasFirstTopupRebate {
			return nil
		}

		inviterCredit, err := createTopUpRebateTx(
			tx,
			topUp,
			invitee.Id,
			invitee.InviterId,
			topupRebateRecipientRoleInviter,
			triggerMode,
			operation_setting.GetInviterTopupRebatePercent(),
		)
		if err != nil {
			return err
		}
		if inviterCredit != nil {
			credited = append(credited, *inviterCredit)
		}

		inviteeCredit, err := createTopUpRebateTx(
			tx,
			topUp,
			invitee.Id,
			invitee.Id,
			topupRebateRecipientRoleInvitee,
			triggerMode,
			operation_setting.GetInviteeTopupRebatePercent(),
		)
		if err != nil {
			return err
		}
		if inviteeCredit != nil {
			credited = append(credited, *inviteeCredit)
		}

		if triggerMode == operation_setting.AffiliateRewardTriggerFirstTopup {
			if err := tx.Model(&User{}).Where("id = ?", invitee.Id).
				Update("has_first_topup_rebate", true).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, item := range credited {
		if item.role == topupRebateRecipientRoleInvitee {
			_ = cacheIncrUserQuota(item.userId, int64(item.quota))
			RecordLog(item.userId, LogTypeSystem, fmt.Sprintf("邀请充值返利到账 %s", formatUSDAmount(item.amountUSD)))
			continue
		}
		RecordLog(item.userId, LogTypeSystem, fmt.Sprintf("邀请用户充值返利 %s", formatUSDAmount(item.amountUSD)))
	}
	return nil
}
