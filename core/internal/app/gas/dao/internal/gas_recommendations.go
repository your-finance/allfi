// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GasRecommendationsDao is the data access object for the table gas_recommendations.
type GasRecommendationsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  GasRecommendationsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// GasRecommendationsColumns defines and stores column names for the table gas_recommendations.
type GasRecommendationsColumns struct {
	Id               string //
	CreatedAt        string //
	UpdatedAt        string //
	Chain            string //
	RecommendedTime  string //
	EstimatedSavings string //
	Confidence       string //
	ValidUntil       string //
}

// gasRecommendationsColumns holds the columns for the table gas_recommendations.
var gasRecommendationsColumns = GasRecommendationsColumns{
	Id:               "id",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	Chain:            "chain",
	RecommendedTime:  "recommended_time",
	EstimatedSavings: "estimated_savings",
	Confidence:       "confidence",
	ValidUntil:       "valid_until",
}

// NewGasRecommendationsDao creates and returns a new DAO object for table data access.
func NewGasRecommendationsDao(handlers ...gdb.ModelHandler) *GasRecommendationsDao {
	return &GasRecommendationsDao{
		group:    "default",
		table:    "gas_recommendations",
		columns:  gasRecommendationsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GasRecommendationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GasRecommendationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GasRecommendationsDao) Columns() GasRecommendationsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GasRecommendationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GasRecommendationsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GasRecommendationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
