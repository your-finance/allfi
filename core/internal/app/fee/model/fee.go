// Package model 费用分析模块业务 DTO
// 定义费用分析相关的数据传输对象
package model

// FeeAnalytics 费用分析结果
type FeeAnalytics struct {
	TotalFees      float64        `json:"total_fees"`       // 总费用
	TradingFees    float64        `json:"trading_fees"`     // 交易手续费总计
	GasFees        float64        `json:"gas_fees"`         // Gas 费用总计
	Currency       string         `json:"currency"`         // 计价货币
	Breakdown      []FeeBreakdown `json:"breakdown"`        // 费用明细（按来源和类型）
	DailyTrend     []DailyFee     `json:"daily_trend"`      // 每日费用趋势
	ComparePercent float64        `json:"compare_percent"`  // 与上期对比百分比
	Suggestions    []string       `json:"suggestions"`      // 优化建议
}

// FeeBreakdown 费用明细（按来源和类型分组）
type FeeBreakdown struct {
	Source   string  `json:"source"`   // 来源（交易所名/链名）
	Type     string  `json:"type"`     // 费用类型（trading_fee/gas_fee）
	Amount   float64 `json:"amount"`   // 费用金额
	Currency string  `json:"currency"` // 计价货币
	Count    int     `json:"count"`    // 交易次数
}

// DailyFee 每日费用
type DailyFee struct {
	Date   string  `json:"date"`   // 日期（YYYY-MM-DD）
	Amount float64 `json:"amount"` // 当日费用总额
}

// MonthlyFee 月度费用
type MonthlyFee struct {
	Month    string  `json:"month"`     // 格式 "2026-01"
	TotalFee float64 `json:"total_fee"` // 当月总费用
	Trading  float64 `json:"trading"`   // 交易手续费
	Gas      float64 `json:"gas"`       // Gas 费
	Withdraw float64 `json:"withdraw"`  // 提现费
}

// 费用类型常量
const (
	FeeTypeTrade    = "trading_fee"    // 交易手续费
	FeeTypeGas      = "gas_fee"        // 链上 Gas 费
	FeeTypeWithdraw = "withdrawal_fee" // 提现手续费
)
