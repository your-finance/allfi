// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PriceAlerts is the golang structure of table price_alerts for DAO operations like Where/Data.
type PriceAlerts struct {
	g.Meta      `orm:"table:price_alerts, do:true"`
	Id          any //
	CreatedAt   any //
	UpdatedAt   any //
	DeletedAt   any //
	UserId      any //
	Symbol      any //
	Condition   any //
	TargetPrice any //
	IsActive    any //
	Triggered   any //
	TriggeredAt any //
	Note        any //
}
