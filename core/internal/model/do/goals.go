// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Goals is the golang structure of table goals for DAO operations like Where/Data.
type Goals struct {
	g.Meta      `orm:"table:goals, do:true"`
	Id          any //
	CreatedAt   any //
	UpdatedAt   any //
	DeletedAt   any //
	Title       any //
	Type        any //
	TargetValue any //
	Currency    any //
	Deadline    any //
}
