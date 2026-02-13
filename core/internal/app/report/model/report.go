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
