// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GasPriceHistory is the golang structure of table gas_price_history for DAO operations like Where/Data.
type GasPriceHistory struct {
	g.Meta     `orm:"table:gas_price_history, do:true"`
	Id         any //
	CreatedAt  any //
	Chain      any //
	Low        any //
	Standard   any //
	Fast       any //
	Instant    any //
	BaseFee    any //
	RecordedAt any //
}
