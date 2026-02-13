// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AssetSnapshotsDao is the data access object for the table asset_snapshots.
type AssetSnapshotsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  AssetSnapshotsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// AssetSnapshotsColumns defines and stores column names for the table asset_snapshots.
type AssetSnapshotsColumns struct {
	Id                 string //
	CreatedAt          string //
	UpdatedAt          string //
	DeletedAt          string //
	UserId             string //
	SnapshotTime       string //
	TotalValueUsd      string //
	TotalValueCny      string //
	TotalValueBtc      string //
	CexValueUsd        string //
	BlockchainValueUsd string //
	ManualValueUsd     string //
	ExchangeRatesJson  string //
}

// assetSnapshotsColumns holds the columns for the table asset_snapshots.
var assetSnapshotsColumns = AssetSnapshotsColumns{
	Id:                 "id",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
	UserId:             "user_id",
	SnapshotTime:       "snapshot_time",
	TotalValueUsd:      "total_value_usd",
	TotalValueCny:      "total_value_cny",
	TotalValueBtc:      "total_value_btc",
	CexValueUsd:        "cex_value_usd",
	BlockchainValueUsd: "blockchain_value_usd",
	ManualValueUsd:     "manual_value_usd",
	ExchangeRatesJson:  "exchange_rates_json",
}

// NewAssetSnapshotsDao creates and returns a new DAO object for table data access.
func NewAssetSnapshotsDao(handlers ...gdb.ModelHandler) *AssetSnapshotsDao {
	return &AssetSnapshotsDao{
		group:    "default",
		table:    "asset_snapshots",
		columns:  assetSnapshotsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AssetSnapshotsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AssetSnapshotsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AssetSnapshotsDao) Columns() AssetSnapshotsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AssetSnapshotsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AssetSnapshotsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AssetSnapshotsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
