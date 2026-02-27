// Package risk 风险管理模块 API 定义
package risk

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetMetricsReq 获取风险指标请求
type GetMetricsReq struct {
	g.Meta `path:"/risk/metrics" method:"get" tags:"风险管理" summary:"获取最新风险指标"`
}

// GetMetricsRes 获取风险指标响应
type GetMetricsRes struct {
	Metrics *RiskMetrics `json:"metrics" dc:"风险指标"`
}

// GetHistoryReq 获取历史风险指标请求
type GetHistoryReq struct {
	g.Meta `path:"/risk/history" method:"get" tags:"风险管理" summary:"获取历史风险指标"`
	Days   int `json:"days" v:"required|min:7|max:365" dc:"查询天数（7-365）"`
}

// GetHistoryRes 获取历史风险指标响应
type GetHistoryRes struct {
	History []*RiskMetrics `json:"history" dc:"历史风险指标列表"`
}

// CalculateReq 手动触发风险指标计算请求
type CalculateReq struct {
	g.Meta `path:"/risk/calculate" method:"post" tags:"风险管理" summary:"手动触发风险指标计算"`
	Period int `json:"period" v:"min:7|max:365" dc:"计算周期（天数，默认30）"`
}

// CalculateRes 手动触发风险指标计算响应
type CalculateRes struct {
	Metrics *RiskMetrics `json:"metrics" dc:"计算后的风险指标"`
	Message string       `json:"message" dc:"提示信息"`
}

// RiskMetrics 风险指标数据结构
type RiskMetrics struct {
	MetricDate           string  `json:"metric_date" dc:"指标计算日期"`
	PortfolioValue       float64 `json:"portfolio_value" dc:"组合总价值（USD）"`
	Var95                float64 `json:"var_95" dc:"95% 置信度 VaR"`
	Var99                float64 `json:"var_99" dc:"99% 置信度 VaR"`
	SharpeRatio          float64 `json:"sharpe_ratio" dc:"夏普比率"`
	SortinoRatio         float64 `json:"sortino_ratio" dc:"索提诺比率"`
	MaxDrawdown          float64 `json:"max_drawdown" dc:"最大回撤（百分比）"`
	MaxDrawdownDuration  int     `json:"max_drawdown_duration" dc:"最大回撤持续天数"`
	Beta                 float64 `json:"beta" dc:"Beta 系数（相对 BTC）"`
	Volatility           float64 `json:"volatility" dc:"波动率（年化）"`
	DownsideDeviation    float64 `json:"downside_deviation" dc:"下行偏差"`
	CalculationPeriod    int     `json:"calculation_period" dc:"计算周期（天数）"`
}
