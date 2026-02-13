// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReportsDao is the data access object for the table reports.
type ReportsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ReportsColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ReportsColumns defines and stores column names for the table reports.
type ReportsColumns struct {
	Id            string //
	CreatedAt     string //
	UpdatedAt     string //
	DeletedAt     string //
	UserId        string //
	Type          string //
	Period        string //
	TotalValue    string //
	Change        string //
	ChangePercent string //
	TopGainers    string //
	TopLosers     string //
	BtcBenchmark  string //
	EthBenchmark  string //
	Content       string //
	GeneratedAt   string //
}

// reportsColumns holds the columns for the table reports.
var reportsColumns = ReportsColumns{
	Id:            "id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
	UserId:        "user_id",
	Type:          "type",
	Period:        "period",
	TotalValue:    "total_value",
	Change:        "change",
	ChangePercent: "change_percent",
	TopGainers:    "top_gainers",
	TopLosers:     "top_losers",
	BtcBenchmark:  "btc_benchmark",
	EthBenchmark:  "eth_benchmark",
	Content:       "content",
	GeneratedAt:   "generated_at",
}

// NewReportsDao creates and returns a new DAO object for table data access.
func NewReportsDao(handlers ...gdb.ModelHandler) *ReportsDao {
	return &ReportsDao{
		group:    "default",
		table:    "reports",
		columns:  reportsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReportsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReportsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReportsDao) Columns() ReportsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReportsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReportsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReportsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
