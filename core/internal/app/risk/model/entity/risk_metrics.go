// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-28 00:16:43
// =================================================================================

package entity

import (
	"time"
)

// RiskMetrics is the golang structure for table risk_metrics.
type RiskMetrics struct {
	Id                  int       `json:"id"                    orm:"id"                    description:""` //
	CreatedAt           time.Time `json:"created_at"            orm:"created_at"            description:""` //
	UpdatedAt           time.Time `json:"updated_at"            orm:"updated_at"            description:""` //
	DeletedAt           time.Time `json:"deleted_at"            orm:"deleted_at"            description:""` //
	UserId              int       `json:"user_id"               orm:"user_id"               description:""` //
	MetricDate          time.Time `json:"metric_date"           orm:"metric_date"           description:""` //
	PortfolioValue      float32   `json:"portfolio_value"       orm:"portfolio_value"       description:""` //
	Var95               float32   `json:"var_95"                orm:"var_95"                description:""` //
	Var99               float32   `json:"var_99"                orm:"var_99"                description:""` //
	SharpeRatio         float32   `json:"sharpe_ratio"          orm:"sharpe_ratio"          description:""` //
	SortinoRatio        float32   `json:"sortino_ratio"         orm:"sortino_ratio"         description:""` //
	MaxDrawdown         float32   `json:"max_drawdown"          orm:"max_drawdown"          description:""` //
	MaxDrawdownDuration int       `json:"max_drawdown_duration" orm:"max_drawdown_duration" description:""` //
	Beta                float32   `json:"beta"                  orm:"beta"                  description:""` //
	Volatility          float32   `json:"volatility"            orm:"volatility"            description:""` //
	DownsideDeviation   float32   `json:"downside_deviation"    orm:"downside_deviation"    description:""` //
	CalculationPeriod   int       `json:"calculation_period"    orm:"calculation_period"    description:""` //
}
