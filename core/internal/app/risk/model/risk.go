// Package model 风险管理模块业务 DTO
// 定义风险管理相关的数据传输对象和常量
package model

// RiskLevel 风险等级
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"    // 低风险
	RiskLevelMedium RiskLevel = "medium" // 中等风险
	RiskLevelHigh   RiskLevel = "high"   // 高风险
)

// RiskMetrics 风险指标（业务层）
type RiskMetrics struct {
	MetricDate          string  `json:"metric_date"`           // 指标计算日期
	PortfolioValue      float64 `json:"portfolio_value"`       // 组合总价值（USD）
	Var95               float64 `json:"var_95"`                // 95% 置信度 VaR
	Var99               float64 `json:"var_99"`                // 99% 置信度 VaR
	SharpeRatio         float64 `json:"sharpe_ratio"`          // 夏普比率
	SortinoRatio        float64 `json:"sortino_ratio"`         // 索提诺比率
	MaxDrawdown         float64 `json:"max_drawdown"`          // 最大回撤（百分比）
	MaxDrawdownDuration int     `json:"max_drawdown_duration"` // 最大回撤持续天数
	Beta                float64 `json:"beta"`                  // Beta 系数（相对 BTC）
	Volatility          float64 `json:"volatility"`            // 波动率（年化）
	DownsideDeviation   float64 `json:"downside_deviation"`    // 下行偏差
	CalculationPeriod   int     `json:"calculation_period"`    // 计算周期（天数）
	RiskLevel           string  `json:"risk_level"`            // 风险等级
}

// CalculationResult 风险指标计算结果
type CalculationResult struct {
	Var95             float64 // 95% 置信度 VaR
	Var99             float64 // 99% 置信度 VaR
	SharpeRatio       float64 // 夏普比率
	SortinoRatio      float64 // 索提诺比率
	MaxDrawdown       float64 // 最大回撤（百分比）
	MaxDrawdownDays   int     // 最大回撤持续天数
	Beta              float64 // Beta 系数
	Volatility        float64 // 波动率（年化）
	DownsideDeviation float64 // 下行偏差
}

// GetRiskLevel 根据风险指标判断风险等级
func GetRiskLevel(volatility, maxDrawdown float64) RiskLevel {
	// 综合波动率和最大回撤判断风险等级
	// 波动率 > 30% 或 最大回撤 > 20% → 高风险
	// 波动率 > 15% 或 最大回撤 > 10% → 中等风险
	// 其他 → 低风险
	if volatility > 30 || maxDrawdown > 20 {
		return RiskLevelHigh
	}
	if volatility > 15 || maxDrawdown > 10 {
		return RiskLevelMedium
	}
	return RiskLevelLow
}
