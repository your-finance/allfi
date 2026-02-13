// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// TransactionDailySummaries is the golang structure of table transaction_daily_summaries for DAO operations like Where/Data.
type TransactionDailySummaries struct {
	g.Meta      `orm:"table:transaction_daily_summaries, do:true"`
	Id          any //
	CreatedAt   any //
	UpdatedAt   any //
	DeletedAt   any //
	UserId      any //
	Date        any //
	BuyCount    any //
	SellCount   any //
	TotalCount  any //
	TotalFeeUsd any //
	NetFlowUsd  any //
}
