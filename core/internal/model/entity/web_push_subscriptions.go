// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// WebPushSubscriptions is the golang structure for table web_push_subscriptions.
type WebPushSubscriptions struct {
	Id        int       `json:"id"         orm:"id"         description:""` //
	CreatedAt time.Time `json:"created_at" orm:"created_at" description:""` //
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at" description:""` //
	DeletedAt time.Time `json:"deleted_at" orm:"deleted_at" description:""` //
	UserId    int       `json:"user_id"    orm:"user_id"    description:""` //
	Endpoint  string    `json:"endpoint"   orm:"endpoint"   description:""` //
	P256Dh    string    `json:"p_256_dh"   orm:"p256dh"     description:""` //
	Auth      string    `json:"auth"       orm:"auth"       description:""` //
}
