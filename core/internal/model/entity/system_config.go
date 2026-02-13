// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// SystemConfig is the golang structure for table system_config.
type SystemConfig struct {
	Id          int       `json:"id"           orm:"id"           description:""` //
	CreatedAt   time.Time `json:"created_at"   orm:"created_at"   description:""` //
	UpdatedAt   time.Time `json:"updated_at"   orm:"updated_at"   description:""` //
	DeletedAt   time.Time `json:"deleted_at"   orm:"deleted_at"   description:""` //
	ConfigKey   string    `json:"config_key"   orm:"config_key"   description:""` //
	ConfigValue string    `json:"config_value" orm:"config_value" description:""` //
	Description string    `json:"description"  orm:"description"  description:""` //
}
