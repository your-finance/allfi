// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-25 10:57:34
// =================================================================================

package entity

import (
	"time"
)

// PriceAlerts is the golang structure for table price_alerts.
type PriceAlerts struct {
	Id          int       `json:"id"           orm:"id"           description:""` //
	CreatedAt   time.Time `json:"created_at"   orm:"created_at"   description:""` //
	UpdatedAt   time.Time `json:"updated_at"   orm:"updated_at"   description:""` //
	DeletedAt   time.Time `json:"deleted_at"   orm:"deleted_at"   description:""` //
	UserId      int       `json:"user_id"      orm:"user_id"      description:""` //
	Symbol      string    `json:"symbol"       orm:"symbol"       description:""` //
	Condition   string    `json:"condition"    orm:"condition"    description:""` //
	TargetPrice float64   `json:"target_price" orm:"target_price" description:""` //
	IsActive    int       `json:"is_active"    orm:"is_active"    description:""` //
	Triggered   int       `json:"triggered"    orm:"triggered"    description:""` //
	TriggeredAt time.Time `json:"triggered_at" orm:"triggered_at" description:""` //
	Note        string    `json:"note"         orm:"note"         description:""` //
}
