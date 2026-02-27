// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:10:34
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta       `orm:"table:users, do:true"`
	Id           any //
	CreatedAt    any //
	UpdatedAt    any //
	DeletedAt    any //
	Username     any //
	Email        any //
	PasswordHash any //
	Nickname     any //
	Avatar       any //
	Status       any //
	LastLoginAt  any //
	LastLoginIp  any //
}
