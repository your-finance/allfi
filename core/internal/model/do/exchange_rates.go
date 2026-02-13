// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ExchangeRates is the golang structure of table exchange_rates for DAO operations like Where/Data.
type ExchangeRates struct {
	g.Meta       `orm:"table:exchange_rates, do:true"`
	Id           any //
	CreatedAt    any //
	UpdatedAt    any //
	DeletedAt    any //
	FromCurrency any //
	ToCurrency   any //
	Rate         any //
	Source       any //
	FetchedAt    any //
}
