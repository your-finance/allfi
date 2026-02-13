// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SyncMetadata is the golang structure of table sync_metadata for DAO operations like Where/Data.
type SyncMetadata struct {
	g.Meta       `orm:"table:sync_metadata, do:true"`
	Id           any //
	CreatedAt    any //
	UpdatedAt    any //
	DeletedAt    any //
	Source       any //
	LastSyncTime any //
	LastSyncId   any //
	LastBlock    any //
	TxCount      any //
}
