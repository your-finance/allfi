// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ManualAssets is the golang structure of table manual_assets for DAO operations like Where/Data.
type ManualAssets struct {
	g.Meta    `orm:"table:manual_assets, do:true"`
	Id        any //
	CreatedAt any //
	UpdatedAt any //
	DeletedAt any //
	UserId    any //
	AssetType any //
	AssetName any //
	Amount    any //
	AmountUsd any //
	Currency  any //
	Notes     any //
	IsActive  any //
}
