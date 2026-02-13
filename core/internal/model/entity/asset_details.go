// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// AssetDetails is the golang structure for table asset_details.
type AssetDetails struct {
	Id          int       `json:"id"           orm:"id"           description:""` //
	CreatedAt   time.Time `json:"created_at"   orm:"created_at"   description:""` //
	UpdatedAt   time.Time `json:"updated_at"   orm:"updated_at"   description:""` //
	DeletedAt   time.Time `json:"deleted_at"   orm:"deleted_at"   description:""` //
	UserId      int       `json:"user_id"      orm:"user_id"      description:""` //
	SourceType  string    `json:"source_type"  orm:"source_type"  description:""` //
	SourceId    int       `json:"source_id"    orm:"source_id"    description:""` //
	AssetSymbol string    `json:"asset_symbol" orm:"asset_symbol" description:""` //
	AssetName   string    `json:"asset_name"   orm:"asset_name"   description:""` //
	Balance     float64   `json:"balance"      orm:"balance"      description:""` //
	PriceUsd    float64   `json:"price_usd"    orm:"price_usd"    description:""` //
	ValueUsd    float64   `json:"value_usd"    orm:"value_usd"    description:""` //
	LastUpdated time.Time `json:"last_updated" orm:"last_updated" description:""` //
}
