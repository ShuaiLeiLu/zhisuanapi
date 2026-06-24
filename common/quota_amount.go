package common

import "github.com/shopspring/decimal"

func AmountToQuota(amount float64) int {
	if amount <= 0 || QuotaPerUnit <= 0 {
		return 0
	}
	quota := decimal.NewFromFloat(amount).
		Mul(decimal.NewFromFloat(QuotaPerUnit)).
		Round(0)
	return int(quota.IntPart())
}

func QuotaToAmount(quota int) float64 {
	if quota <= 0 || QuotaPerUnit <= 0 {
		return 0
	}
	return decimal.NewFromInt(int64(quota)).
		Div(decimal.NewFromFloat(QuotaPerUnit)).
		InexactFloat64()
}
