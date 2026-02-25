// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-25 10:57:34
// =================================================================================

package entity

import (
	"time"
)

// Users is the golang structure for table users.
type Users struct {
	Id           int       `json:"id"            orm:"id"            description:""` //
	CreatedAt    time.Time `json:"created_at"    orm:"created_at"    description:""` //
	UpdatedAt    time.Time `json:"updated_at"    orm:"updated_at"    description:""` //
	DeletedAt    time.Time `json:"deleted_at"    orm:"deleted_at"    description:""` //
	Username     string    `json:"username"      orm:"username"      description:""` //
	Email        string    `json:"email"         orm:"email"         description:""` //
	PasswordHash string    `json:"password_hash" orm:"password_hash" description:""` //
	Nickname     string    `json:"nickname"      orm:"nickname"      description:""` //
	Avatar       string    `json:"avatar"        orm:"avatar"        description:""` //
	Status       string    `json:"status"        orm:"status"        description:""` //
	LastLoginAt  time.Time `json:"last_login_at" orm:"last_login_at" description:""` //
	LastLoginIp  string    `json:"last_login_ip" orm:"last_login_ip" description:""` //
}
