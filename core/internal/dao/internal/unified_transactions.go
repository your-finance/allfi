// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UnifiedTransactionsDao is the data access object for the table unified_transactions.
type UnifiedTransactionsDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  UnifiedTransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// UnifiedTransactionsColumns defines and stores column names for the table unified_transactions.
type UnifiedTransactionsColumns struct {
	Id         string //
	CreatedAt  string //
	UpdatedAt  string //
	DeletedAt  string //
	UserId     string //
	TxType     string //
	Source     string //
	SourceId   string //
	FromAsset  string //
	FromAmount string //
	ToAsset    string //
	ToAmount   string //
	Fee        string //
	FeeCoin    string //
	ValueUsd   string //
	TxHash     string //
	Chain      string //
	Timestamp  string //
}

// unifiedTransactionsColumns holds the columns for the table unified_transactions.
var unifiedTransactionsColumns = UnifiedTransactionsColumns{
	Id:         "id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
	UserId:     "user_id",
	TxType:     "tx_type",
	Source:     "source",
	SourceId:   "source_id",
	FromAsset:  "from_asset",
	FromAmount: "from_amount",
	ToAsset:    "to_asset",
	ToAmount:   "to_amount",
	Fee:        "fee",
	FeeCoin:    "fee_coin",
	ValueUsd:   "value_usd",
	TxHash:     "tx_hash",
	Chain:      "chain",
	Timestamp:  "timestamp",
}

// NewUnifiedTransactionsDao creates and returns a new DAO object for table data access.
func NewUnifiedTransactionsDao(handlers ...gdb.ModelHandler) *UnifiedTransactionsDao {
	return &UnifiedTransactionsDao{
		group:    "default",
		table:    "unified_transactions",
		columns:  unifiedTransactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UnifiedTransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UnifiedTransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UnifiedTransactionsDao) Columns() UnifiedTransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UnifiedTransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UnifiedTransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UnifiedTransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
