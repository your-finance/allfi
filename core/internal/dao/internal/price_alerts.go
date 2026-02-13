// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PriceAlertsDao is the data access object for the table price_alerts.
type PriceAlertsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PriceAlertsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PriceAlertsColumns defines and stores column names for the table price_alerts.
type PriceAlertsColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	UserId      string //
	Symbol      string //
	Condition   string //
	TargetPrice string //
	IsActive    string //
	Triggered   string //
	TriggeredAt string //
	Note        string //
}

// priceAlertsColumns holds the columns for the table price_alerts.
var priceAlertsColumns = PriceAlertsColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserId:      "user_id",
	Symbol:      "symbol",
	Condition:   "condition",
	TargetPrice: "target_price",
	IsActive:    "is_active",
	Triggered:   "triggered",
	TriggeredAt: "triggered_at",
	Note:        "note",
}

// NewPriceAlertsDao creates and returns a new DAO object for table data access.
func NewPriceAlertsDao(handlers ...gdb.ModelHandler) *PriceAlertsDao {
	return &PriceAlertsDao{
		group:    "default",
		table:    "price_alerts",
		columns:  priceAlertsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PriceAlertsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PriceAlertsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PriceAlertsDao) Columns() PriceAlertsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PriceAlertsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PriceAlertsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PriceAlertsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
