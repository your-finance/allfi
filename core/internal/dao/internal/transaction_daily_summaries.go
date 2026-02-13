// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TransactionDailySummariesDao is the data access object for the table transaction_daily_summaries.
type TransactionDailySummariesDao struct {
	table    string                           // table is the underlying table name of the DAO.
	group    string                           // group is the database configuration group name of the current DAO.
	columns  TransactionDailySummariesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler               // handlers for customized model modification.
}

// TransactionDailySummariesColumns defines and stores column names for the table transaction_daily_summaries.
type TransactionDailySummariesColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	UserId      string //
	Date        string //
	BuyCount    string //
	SellCount   string //
	TotalCount  string //
	TotalFeeUsd string //
	NetFlowUsd  string //
}

// transactionDailySummariesColumns holds the columns for the table transaction_daily_summaries.
var transactionDailySummariesColumns = TransactionDailySummariesColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserId:      "user_id",
	Date:        "date",
	BuyCount:    "buy_count",
	SellCount:   "sell_count",
	TotalCount:  "total_count",
	TotalFeeUsd: "total_fee_usd",
	NetFlowUsd:  "net_flow_usd",
}

// NewTransactionDailySummariesDao creates and returns a new DAO object for table data access.
func NewTransactionDailySummariesDao(handlers ...gdb.ModelHandler) *TransactionDailySummariesDao {
	return &TransactionDailySummariesDao{
		group:    "default",
		table:    "transaction_daily_summaries",
		columns:  transactionDailySummariesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TransactionDailySummariesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TransactionDailySummariesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TransactionDailySummariesDao) Columns() TransactionDailySummariesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TransactionDailySummariesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TransactionDailySummariesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TransactionDailySummariesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
