// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WebPushSubscriptionsDao is the data access object for the table web_push_subscriptions.
type WebPushSubscriptionsDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  WebPushSubscriptionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// WebPushSubscriptionsColumns defines and stores column names for the table web_push_subscriptions.
type WebPushSubscriptionsColumns struct {
	Id        string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
	UserId    string //
	Endpoint  string //
	P256Dh    string //
	Auth      string //
}

// webPushSubscriptionsColumns holds the columns for the table web_push_subscriptions.
var webPushSubscriptionsColumns = WebPushSubscriptionsColumns{
	Id:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
	UserId:    "user_id",
	Endpoint:  "endpoint",
	P256Dh:    "p256dh",
	Auth:      "auth",
}

// NewWebPushSubscriptionsDao creates and returns a new DAO object for table data access.
func NewWebPushSubscriptionsDao(handlers ...gdb.ModelHandler) *WebPushSubscriptionsDao {
	return &WebPushSubscriptionsDao{
		group:    "default",
		table:    "web_push_subscriptions",
		columns:  webPushSubscriptionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *WebPushSubscriptionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *WebPushSubscriptionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *WebPushSubscriptionsDao) Columns() WebPushSubscriptionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *WebPushSubscriptionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *WebPushSubscriptionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *WebPushSubscriptionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
