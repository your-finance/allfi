// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// AssetSnapshots is the golang structure for table asset_snapshots.
type AssetSnapshots struct {
	Id                 int       `json:"id"                   orm:"id"                   description:""` //
	CreatedAt          time.Time `json:"created_at"           orm:"created_at"           description:""` //
	UpdatedAt          time.Time `json:"updated_at"           orm:"updated_at"           description:""` //
	DeletedAt          time.Time `json:"deleted_at"           orm:"deleted_at"           description:""` //
	UserId             int       `json:"user_id"              orm:"user_id"              description:""` //
	SnapshotTime       time.Time `json:"snapshot_time"        orm:"snapshot_time"        description:""` //
	TotalValueUsd      float64   `json:"total_value_usd"      orm:"total_value_usd"      description:""` //
	TotalValueCny      float64   `json:"total_value_cny"      orm:"total_value_cny"      description:""` //
	TotalValueBtc      float64   `json:"total_value_btc"      orm:"total_value_btc"      description:""` //
	CexValueUsd        float64   `json:"cex_value_usd"        orm:"cex_value_usd"        description:""` //
	BlockchainValueUsd float64   `json:"blockchain_value_usd" orm:"blockchain_value_usd" description:""` //
	ManualValueUsd     float64   `json:"manual_value_usd"     orm:"manual_value_usd"     description:""` //
	ExchangeRatesJson  string    `json:"exchange_rates_json"  orm:"exchange_rates_json"  description:""` //
}
