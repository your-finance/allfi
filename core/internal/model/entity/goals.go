// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// Goals is the golang structure for table goals.
type Goals struct {
	Id          int       `json:"id"           orm:"id"           description:""` //
	CreatedAt   time.Time `json:"created_at"   orm:"created_at"   description:""` //
	UpdatedAt   time.Time `json:"updated_at"   orm:"updated_at"   description:""` //
	DeletedAt   time.Time `json:"deleted_at"   orm:"deleted_at"   description:""` //
	Title       string    `json:"title"        orm:"title"        description:""` //
	Type        string    `json:"type"         orm:"type"         description:""` //
	TargetValue float32   `json:"target_value" orm:"target_value" description:""` //
	Currency    string    `json:"currency"     orm:"currency"     description:""` //
	Deadline    time.Time `json:"deadline"     orm:"deadline"     description:""` //
}
