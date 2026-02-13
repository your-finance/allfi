// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UnifiedTransactions is the golang structure of table unified_transactions for DAO operations like Where/Data.
type UnifiedTransactions struct {
	g.Meta     `orm:"table:unified_transactions, do:true"`
	Id         any //
	CreatedAt  any //
	UpdatedAt  any //
	DeletedAt  any //
	UserId     any //
	TxType     any //
	Source     any //
	SourceId   any //
	FromAsset  any //
	FromAmount any //
	ToAsset    any //
	ToAmount   any //
	Fee        any //
	FeeCoin    any //
	ValueUsd   any //
	TxHash     any //
	Chain      any //
	Timestamp  any //
}
