// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// Strategies is the golang structure for table strategies.
type Strategies struct {
	Id            int       `json:"id"             orm:"id"             description:""` //
	CreatedAt     time.Time `json:"created_at"     orm:"created_at"     description:""` //
	UpdatedAt     time.Time `json:"updated_at"     orm:"updated_at"     description:""` //
	DeletedAt     time.Time `json:"deleted_at"     orm:"deleted_at"     description:""` //
	UserId        int       `json:"user_id"        orm:"user_id"        description:""` //
	Name          string    `json:"name"           orm:"name"           description:""` //
	Type          string    `json:"type"           orm:"type"           description:""` //
	Config        string    `json:"config"         orm:"config"         description:""` //
	IsActive      float64   `json:"is_active"      orm:"is_active"      description:""` //
	LastChecked   time.Time `json:"last_checked"   orm:"last_checked"   description:""` //
	LastTriggered time.Time `json:"last_triggered" orm:"last_triggered" description:""` //
	TriggerCount  int       `json:"trigger_count"  orm:"trigger_count"  description:""` //
}
