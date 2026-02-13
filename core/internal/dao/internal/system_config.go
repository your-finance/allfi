// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemConfigDao is the data access object for the table system_config.
type SystemConfigDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SystemConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SystemConfigColumns defines and stores column names for the table system_config.
type SystemConfigColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	ConfigKey   string //
	ConfigValue string //
	Description string //
}

// systemConfigColumns holds the columns for the table system_config.
var systemConfigColumns = SystemConfigColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	ConfigKey:   "config_key",
	ConfigValue: "config_value",
	Description: "description",
}

// NewSystemConfigDao creates and returns a new DAO object for table data access.
func NewSystemConfigDao(handlers ...gdb.ModelHandler) *SystemConfigDao {
	return &SystemConfigDao{
		group:    "default",
		table:    "system_config",
		columns:  systemConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SystemConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SystemConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SystemConfigDao) Columns() SystemConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SystemConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SystemConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SystemConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
