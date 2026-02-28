// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskMetricsDao is the data access object for the table risk_metrics.
type RiskMetricsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RiskMetricsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RiskMetricsColumns defines and stores column names for the table risk_metrics.
type RiskMetricsColumns struct {
	Id                  string //
	CreatedAt           string //
	UpdatedAt           string //
	DeletedAt           string //
	UserId              string //
	MetricDate          string //
	PortfolioValue      string //
	Var95               string //
	Var99               string //
	SharpeRatio         string //
	SortinoRatio        string //
	MaxDrawdown         string //
	MaxDrawdownDuration string //
	Beta                string //
	Volatility          string //
	DownsideDeviation   string //
	CalculationPeriod   string //
}

// riskMetricsColumns holds the columns for the table risk_metrics.
var riskMetricsColumns = RiskMetricsColumns{
	Id:                  "id",
	CreatedAt:           "created_at",
	UpdatedAt:           "updated_at",
	DeletedAt:           "deleted_at",
	UserId:              "user_id",
	MetricDate:          "metric_date",
	PortfolioValue:      "portfolio_value",
	Var95:               "var_95",
	Var99:               "var_99",
	SharpeRatio:         "sharpe_ratio",
	SortinoRatio:        "sortino_ratio",
	MaxDrawdown:         "max_drawdown",
	MaxDrawdownDuration: "max_drawdown_duration",
	Beta:                "beta",
	Volatility:          "volatility",
	DownsideDeviation:   "downside_deviation",
	CalculationPeriod:   "calculation_period",
}

// NewRiskMetricsDao creates and returns a new DAO object for table data access.
func NewRiskMetricsDao(handlers ...gdb.ModelHandler) *RiskMetricsDao {
	return &RiskMetricsDao{
		group:    "default",
		table:    "risk_metrics",
		columns:  riskMetricsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskMetricsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskMetricsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskMetricsDao) Columns() RiskMetricsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskMetricsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskMetricsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RiskMetricsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
