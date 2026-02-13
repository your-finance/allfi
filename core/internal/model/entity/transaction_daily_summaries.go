// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// TransactionDailySummaries is the golang structure for table transaction_daily_summaries.
type TransactionDailySummaries struct {
	Id          int       `json:"id"            orm:"id"            description:""` //
	CreatedAt   time.Time `json:"created_at"    orm:"created_at"    description:""` //
	UpdatedAt   time.Time `json:"updated_at"    orm:"updated_at"    description:""` //
	DeletedAt   time.Time `json:"deleted_at"    orm:"deleted_at"    description:""` //
	UserId      int       `json:"user_id"       orm:"user_id"       description:""` //
	Date        time.Time `json:"date"          orm:"date"          description:""` //
	BuyCount    int       `json:"buy_count"     orm:"buy_count"     description:""` //
	SellCount   int       `json:"sell_count"    orm:"sell_count"    description:""` //
	TotalCount  int       `json:"total_count"   orm:"total_count"   description:""` //
	TotalFeeUsd float32   `json:"total_fee_usd" orm:"total_fee_usd" description:""` //
	NetFlowUsd  float32   `json:"net_flow_usd"  orm:"net_flow_usd"  description:""` //
}
