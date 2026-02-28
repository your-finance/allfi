// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 08:41:17
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RiskMetrics is the golang structure of table risk_metrics for DAO operations like Where/Data.
type RiskMetrics struct {
	g.Meta              `orm:"table:risk_metrics, do:true"`
	Id                  any //
	CreatedAt           any //
	UpdatedAt           any //
	DeletedAt           any //
	UserId              any //
	MetricDate          any //
	PortfolioValue      any //
	Var95               any //
	Var99               any //
	SharpeRatio         any //
	SortinoRatio        any //
	MaxDrawdown         any //
	MaxDrawdownDuration any //
	Beta                any //
	Volatility          any //
	DownsideDeviation   any //
	CalculationPeriod   any //
}
