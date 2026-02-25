// Package model 报告模块 - 数据传输对象 (DTO)
package model

// 报告类型常量
const (
	ReportTypeDaily   = "daily"
	ReportTypeWeekly  = "weekly"
	ReportTypeMonthly = "monthly"
	ReportTypeAnnual  = "annual"
)

// GainerLoser 涨跌幅资产
type GainerLoser struct {
	Symbol string  `json:"symbol"` // 币种符号
	Change float64 `json:"change"` // 变化百分比
	Value  float64 `json:"value"`  // 当前价值
}

// MonthReturn 月度收益
type MonthReturn struct {
	Month  int     `json:"month"`  // 月份（1-12）
	Return float64 `json:"return"` // 收益率百分比
	Value  float64 `json:"value"`  // 月末资产价值
}

// BestWorstMonth 最佳/最差月份
type BestWorstMonth struct {
	Month  int     `json:"month"`  // 月份（1-12）
	Return float64 `json:"return"` // 收益率百分比
}

// AnnualSummary 年度总览
type AnnualSummary struct {
	TotalReturn      float64         `json:"totalReturn"`      // 全年收益率
	TotalReturnValue float64         `json:"totalReturnValue"` // 全年收益金额
	StartValue       float64         `json:"startValue"`       // 年初资产
	EndValue         float64         `json:"endValue"`         // 年末资产
	BestMonth        *BestWorstMonth `json:"bestMonth"`        // 最佳月份
	WorstMonth       *BestWorstMonth `json:"worstMonth"`       // 最差月份
	TotalTransactions int            `json:"totalTransactions"` // 总交易次数
	TotalFeesPaid    float64         `json:"totalFeesPaid"`    // 总费用
}

// AnnualBenchmarks 基准对比
type AnnualBenchmarks struct {
	Btc   float64 `json:"btc"`   // BTC 年度收益率
	Eth   float64 `json:"eth"`   // ETH 年度收益率
	Sp500 float64 `json:"sp500"` // S&P 500 年度收益率
}

// AnnualReportContent 年度报告完整内容（前端期望的结构）
type AnnualReportContent struct {
	Year           int              `json:"year"`           // 年份
	Summary        AnnualSummary    `json:"summary"`        // 年度总览
	MonthlyReturns []MonthReturn    `json:"monthlyReturns"` // 月度收益列表
	Benchmarks     AnnualBenchmarks `json:"benchmarks"`     // 基准对比
	TopGainers     []GainerLoser    `json:"topGainers"`     // 涨幅 Top
	TopLosers      []GainerLoser    `json:"topLosers"`      // 跌幅 Top
}

// ReportContent 报告详情内容
type ReportContent struct {
	Type          string        `json:"type"`           // 报告类型
	TotalValue    float64       `json:"total_value"`    // 总资产价值
	Change        float64       `json:"change"`         // 变化金额
	ChangePercent float64       `json:"change_percent"` // 变化百分比
	SnapshotCount int           `json:"snapshot_count"` // 快照数量
	TopGainers    []GainerLoser `json:"top_gainers"`    // 涨幅 Top
	TopLosers     []GainerLoser `json:"top_losers"`     // 跌幅 Top
	BtcBenchmark  float64       `json:"btc_benchmark"`  // BTC 基准
	EthBenchmark  float64       `json:"eth_benchmark"`  // ETH 基准
	Summary       string        `json:"summary"`        // 摘要文本
}
