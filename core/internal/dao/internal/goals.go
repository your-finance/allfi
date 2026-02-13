// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GoalsDao is the data access object for the table goals.
type GoalsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GoalsColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GoalsColumns defines and stores column names for the table goals.
type GoalsColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	Title       string //
	Type        string //
	TargetValue string //
	Currency    string //
	Deadline    string //
}

// goalsColumns holds the columns for the table goals.
var goalsColumns = GoalsColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	Title:       "title",
	Type:        "type",
	TargetValue: "target_value",
	Currency:    "currency",
	Deadline:    "deadline",
}

// NewGoalsDao creates and returns a new DAO object for table data access.
func NewGoalsDao(handlers ...gdb.ModelHandler) *GoalsDao {
	return &GoalsDao{
		group:    "default",
		table:    "goals",
		columns:  goalsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GoalsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GoalsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GoalsDao) Columns() GoalsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GoalsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GoalsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GoalsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
