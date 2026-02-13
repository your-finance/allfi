// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AssetSnapshots is the golang structure of table asset_snapshots for DAO operations like Where/Data.
type AssetSnapshots struct {
	g.Meta             `orm:"table:asset_snapshots, do:true"`
	Id                 any //
	CreatedAt          any //
	UpdatedAt          any //
	DeletedAt          any //
	UserId             any //
	SnapshotTime       any //
	TotalValueUsd      any //
	TotalValueCny      any //
	TotalValueBtc      any //
	CexValueUsd        any //
	BlockchainValueUsd any //
	ManualValueUsd     any //
	ExchangeRatesJson  any //
}
