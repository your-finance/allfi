// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// WalletAddresses is the golang structure for table wallet_addresses.
type WalletAddresses struct {
	Id         int       `json:"id"         orm:"id"         description:""` //
	CreatedAt  time.Time `json:"created_at" orm:"created_at" description:""` //
	UpdatedAt  time.Time `json:"updated_at" orm:"updated_at" description:""` //
	DeletedAt  time.Time `json:"deleted_at" orm:"deleted_at" description:""` //
	UserId     int       `json:"user_id"    orm:"user_id"    description:""` //
	Blockchain string    `json:"blockchain" orm:"blockchain" description:""` //
	Address    string    `json:"address"    orm:"address"    description:""` //
	Label      string    `json:"label"      orm:"label"      description:""` //
	IsActive   float64   `json:"is_active"  orm:"is_active"  description:""` //
}
