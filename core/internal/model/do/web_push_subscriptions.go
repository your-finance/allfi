// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// WebPushSubscriptions is the golang structure of table web_push_subscriptions for DAO operations like Where/Data.
type WebPushSubscriptions struct {
	g.Meta    `orm:"table:web_push_subscriptions, do:true"`
	Id        any //
	CreatedAt any //
	UpdatedAt any //
	DeletedAt any //
	UserId    any //
	Endpoint  any //
	P256Dh    any //
	Auth      any //
}
