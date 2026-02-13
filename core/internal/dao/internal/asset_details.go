// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AssetDetailsDao is the data access object for the table asset_details.
type AssetDetailsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  AssetDetailsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// AssetDetailsColumns defines and stores column names for the table asset_details.
type AssetDetailsColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	UserId      string //
	SourceType  string //
	SourceId    string //
	AssetSymbol string //
	AssetName   string //
	Balance     string //
	PriceUsd    string //
	ValueUsd    string //
	LastUpdated string //
}

// assetDetailsColumns holds the columns for the table asset_details.
var assetDetailsColumns = AssetDetailsColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserId:      "user_id",
	SourceType:  "source_type",
	SourceId:    "source_id",
	AssetSymbol: "asset_symbol",
	AssetName:   "asset_name",
	Balance:     "balance",
	PriceUsd:    "price_usd",
	ValueUsd:    "value_usd",
	LastUpdated: "last_updated",
}

// NewAssetDetailsDao creates and returns a new DAO object for table data access.
func NewAssetDetailsDao(handlers ...gdb.ModelHandler) *AssetDetailsDao {
	return &AssetDetailsDao{
		group:    "default",
		table:    "asset_details",
		columns:  assetDetailsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AssetDetailsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AssetDetailsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AssetDetailsDao) Columns() AssetDetailsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AssetDetailsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AssetDetailsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AssetDetailsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
