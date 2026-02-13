// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// SyncMetadata is the golang structure for table sync_metadata.
type SyncMetadata struct {
	Id           int       `json:"id"             orm:"id"             description:""` //
	CreatedAt    time.Time `json:"created_at"     orm:"created_at"     description:""` //
	UpdatedAt    time.Time `json:"updated_at"     orm:"updated_at"     description:""` //
	DeletedAt    time.Time `json:"deleted_at"     orm:"deleted_at"     description:""` //
	Source       string    `json:"source"         orm:"source"         description:""` //
	LastSyncTime time.Time `json:"last_sync_time" orm:"last_sync_time" description:""` //
	LastSyncId   string    `json:"last_sync_id"   orm:"last_sync_id"   description:""` //
	LastBlock    int       `json:"last_block"     orm:"last_block"     description:""` //
	TxCount      int       `json:"tx_count"       orm:"tx_count"       description:""` //
}
