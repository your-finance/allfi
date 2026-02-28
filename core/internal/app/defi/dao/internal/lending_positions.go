// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:16:43
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LendingPositionsDao is the data access object for the table lending_positions.
type LendingPositionsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  LendingPositionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// LendingPositionsColumns defines and stores column names for the table lending_positions.
type LendingPositionsColumns struct {
	Id                   string //
	UserId               string //
	Protocol             string //
	Chain                string //
	WalletAddress        string //
	SupplyToken          string //
	SupplyAmount         string //
	SupplyValueUsd       string //
	SupplyApy            string //
	BorrowToken          string //
	BorrowAmount         string //
	BorrowValueUsd       string //
	BorrowApy            string //
	HealthFactor         string //
	LiquidationThreshold string //
	Ltv                  string //
	NetApy               string //
	CreatedAt            string //
	UpdatedAt            string //
}

// lendingPositionsColumns holds the columns for the table lending_positions.
var lendingPositionsColumns = LendingPositionsColumns{
	Id:                   "id",
	UserId:               "user_id",
	Protocol:             "protocol",
	Chain:                "chain",
	WalletAddress:        "wallet_address",
	SupplyToken:          "supply_token",
	SupplyAmount:         "supply_amount",
	SupplyValueUsd:       "supply_value_usd",
	SupplyApy:            "supply_apy",
	BorrowToken:          "borrow_token",
	BorrowAmount:         "borrow_amount",
	BorrowValueUsd:       "borrow_value_usd",
	BorrowApy:            "borrow_apy",
	HealthFactor:         "health_factor",
	LiquidationThreshold: "liquidation_threshold",
	Ltv:                  "ltv",
	NetApy:               "net_apy",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewLendingPositionsDao creates and returns a new DAO object for table data access.
func NewLendingPositionsDao(handlers ...gdb.ModelHandler) *LendingPositionsDao {
	return &LendingPositionsDao{
		group:    "default",
		table:    "lending_positions",
		columns:  lendingPositionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LendingPositionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LendingPositionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LendingPositionsDao) Columns() LendingPositionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LendingPositionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LendingPositionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LendingPositionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
