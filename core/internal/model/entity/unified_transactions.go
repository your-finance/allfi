// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// UnifiedTransactions is the golang structure for table unified_transactions.
type UnifiedTransactions struct {
	Id         int       `json:"id"          orm:"id"          description:""` //
	CreatedAt  time.Time `json:"created_at"  orm:"created_at"  description:""` //
	UpdatedAt  time.Time `json:"updated_at"  orm:"updated_at"  description:""` //
	DeletedAt  time.Time `json:"deleted_at"  orm:"deleted_at"  description:""` //
	UserId     int       `json:"user_id"     orm:"user_id"     description:""` //
	TxType     string    `json:"tx_type"     orm:"tx_type"     description:""` //
	Source     string    `json:"source"      orm:"source"      description:""` //
	SourceId   string    `json:"source_id"   orm:"source_id"   description:""` //
	FromAsset  string    `json:"from_asset"  orm:"from_asset"  description:""` //
	FromAmount float32   `json:"from_amount" orm:"from_amount" description:""` //
	ToAsset    string    `json:"to_asset"    orm:"to_asset"    description:""` //
	ToAmount   float32   `json:"to_amount"   orm:"to_amount"   description:""` //
	Fee        float32   `json:"fee"         orm:"fee"         description:""` //
	FeeCoin    string    `json:"fee_coin"    orm:"fee_coin"    description:""` //
	ValueUsd   float32   `json:"value_usd"   orm:"value_usd"   description:""` //
	TxHash     string    `json:"tx_hash"     orm:"tx_hash"     description:""` //
	Chain      string    `json:"chain"       orm:"chain"       description:""` //
	Timestamp  time.Time `json:"timestamp"   orm:"timestamp"   description:""` //
}
