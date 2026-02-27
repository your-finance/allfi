// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:10:34
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LendingRateHistory is the golang structure of table lending_rate_history for DAO operations like Where/Data.
type LendingRateHistory struct {
	g.Meta            `orm:"table:lending_rate_history, do:true"`
	Id                any //
	Protocol          any //
	Chain             any //
	Token             any //
	SupplyApy         any //
	BorrowApyStable   any //
	BorrowApyVariable any //
	UtilizationRate   any //
	RecordedAt        any //
}
