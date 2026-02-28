// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GasPriceHistoryDao is the data access object for the table gas_price_history.
type GasPriceHistoryDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  GasPriceHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// GasPriceHistoryColumns defines and stores column names for the table gas_price_history.
type GasPriceHistoryColumns struct {
	Id         string //
	CreatedAt  string //
	Chain      string //
	Low        string //
	Standard   string //
	Fast       string //
	Instant    string //
	BaseFee    string //
	RecordedAt string //
}

// gasPriceHistoryColumns holds the columns for the table gas_price_history.
var gasPriceHistoryColumns = GasPriceHistoryColumns{
	Id:         "id",
	CreatedAt:  "created_at",
	Chain:      "chain",
	Low:        "low",
	Standard:   "standard",
	Fast:       "fast",
	Instant:    "instant",
	BaseFee:    "base_fee",
	RecordedAt: "recorded_at",
}

// NewGasPriceHistoryDao creates and returns a new DAO object for table data access.
func NewGasPriceHistoryDao(handlers ...gdb.ModelHandler) *GasPriceHistoryDao {
	return &GasPriceHistoryDao{
		group:    "default",
		table:    "gas_price_history",
		columns:  gasPriceHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GasPriceHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GasPriceHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GasPriceHistoryDao) Columns() GasPriceHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GasPriceHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GasPriceHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GasPriceHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
