// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:10:34
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LendingPositions is the golang structure of table lending_positions for DAO operations like Where/Data.
type LendingPositions struct {
	g.Meta               `orm:"table:lending_positions, do:true"`
	Id                   any //
	UserId               any //
	Protocol             any //
	Chain                any //
	WalletAddress        any //
	SupplyToken          any //
	SupplyAmount         any //
	SupplyValueUsd       any //
	SupplyApy            any //
	BorrowToken          any //
	BorrowAmount         any //
	BorrowValueUsd       any //
	BorrowApy            any //
	HealthFactor         any //
	LiquidationThreshold any //
	Ltv                  any //
	NetApy               any //
	CreatedAt            any //
	UpdatedAt            any //
}
