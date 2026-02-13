// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AssetDetails is the golang structure of table asset_details for DAO operations like Where/Data.
type AssetDetails struct {
	g.Meta      `orm:"table:asset_details, do:true"`
	Id          any //
	CreatedAt   any //
	UpdatedAt   any //
	DeletedAt   any //
	UserId      any //
	SourceType  any //
	SourceId    any //
	AssetSymbol any //
	AssetName   any //
	Balance     any //
	PriceUsd    any //
	ValueUsd    any //
	LastUpdated any //
}
