// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SyncMetadataDao is the data access object for the table sync_metadata.
type SyncMetadataDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SyncMetadataColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SyncMetadataColumns defines and stores column names for the table sync_metadata.
type SyncMetadataColumns struct {
	Id           string //
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
	Source       string //
	LastSyncTime string //
	LastSyncId   string //
	LastBlock    string //
	TxCount      string //
}

// syncMetadataColumns holds the columns for the table sync_metadata.
var syncMetadataColumns = SyncMetadataColumns{
	Id:           "id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
	Source:       "source",
	LastSyncTime: "last_sync_time",
	LastSyncId:   "last_sync_id",
	LastBlock:    "last_block",
	TxCount:      "tx_count",
}

// NewSyncMetadataDao creates and returns a new DAO object for table data access.
func NewSyncMetadataDao(handlers ...gdb.ModelHandler) *SyncMetadataDao {
	return &SyncMetadataDao{
		group:    "default",
		table:    "sync_metadata",
		columns:  syncMetadataColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SyncMetadataDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SyncMetadataDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SyncMetadataDao) Columns() SyncMetadataColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SyncMetadataDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SyncMetadataDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SyncMetadataDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
