// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Notifications is the golang structure of table notifications for DAO operations like Where/Data.
type Notifications struct {
	g.Meta    `orm:"table:notifications, do:true"`
	Id        any //
	CreatedAt any //
	UpdatedAt any //
	DeletedAt any //
	UserId    any //
	Type      any //
	Title     any //
	Content   any //
	IsRead    any //
}
