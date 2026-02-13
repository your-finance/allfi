// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id        int       `json:"id"         orm:"id"         description:""` //
	CreatedAt time.Time `json:"created_at" orm:"created_at" description:""` //
	UpdatedAt time.Time `json:"updated_at" orm:"updated_at" description:""` //
	DeletedAt time.Time `json:"deleted_at" orm:"deleted_at" description:""` //
	UserId    int       `json:"user_id"    orm:"user_id"    description:""` //
	Type      string    `json:"type"       orm:"type"       description:""` //
	Title     string    `json:"title"      orm:"title"      description:""` //
	Content   string    `json:"content"    orm:"content"    description:""` //
	IsRead    float64   `json:"is_read"    orm:"is_read"    description:""` //
}
