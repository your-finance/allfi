// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// ManualAssets is the golang structure for table manual_assets.
type ManualAssets struct {
	Id        int       `json:"id"         orm:"id"         description:""` //
	CreatedAt time.Time `json:"created_at" orm:"created_at" description:""` //
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at" description:""` //
	DeletedAt time.Time `json:"deleted_at" orm:"deleted_at" description:""` //
	UserId    int       `json:"user_id"    orm:"user_id"    description:""` //
	AssetType string    `json:"asset_type" orm:"asset_type" description:""` //
	AssetName string    `json:"asset_name" orm:"asset_name" description:""` //
	Amount    float64   `json:"amount"     orm:"amount"     description:""` //
	AmountUsd float64   `json:"amount_usd" orm:"amount_usd" description:""` //
	Currency  string    `json:"currency"   orm:"currency"   description:""` //
	Notes     string    `json:"notes"      orm:"notes"      description:""` //
	IsActive  float64   `json:"is_active"  orm:"is_active"  description:""` //
}
