// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExchangeRatesDao is the data access object for the table exchange_rates.
type ExchangeRatesDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  ExchangeRatesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// ExchangeRatesColumns defines and stores column names for the table exchange_rates.
type ExchangeRatesColumns struct {
	Id           string //
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
	FromCurrency string //
	ToCurrency   string //
	Rate         string //
	Source       string //
	FetchedAt    string //
}

// exchangeRatesColumns holds the columns for the table exchange_rates.
var exchangeRatesColumns = ExchangeRatesColumns{
	Id:           "id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
	FromCurrency: "from_currency",
	ToCurrency:   "to_currency",
	Rate:         "rate",
	Source:       "source",
	FetchedAt:    "fetched_at",
}

// NewExchangeRatesDao creates and returns a new DAO object for table data access.
func NewExchangeRatesDao(handlers ...gdb.ModelHandler) *ExchangeRatesDao {
	return &ExchangeRatesDao{
		group:    "default",
		table:    "exchange_rates",
		columns:  exchangeRatesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExchangeRatesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExchangeRatesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExchangeRatesDao) Columns() ExchangeRatesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExchangeRatesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExchangeRatesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExchangeRatesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
