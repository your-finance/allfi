// Package model 基准对比模块业务 DTO
// 定义基准对比相关的数据传输对象
package model

// BenchmarkResult 基准对比结果
type BenchmarkResult struct {
	Period     string           `json:"period"`      // 时间范围（7d/30d/90d/1y）
	UserReturn float64          `json:"user_return"` // 用户收益率（百分比）
	Benchmarks []BenchmarkIndex `json:"benchmarks"`  // 各基准指数收益率
	StartDate  string           `json:"start_date"`  // 起始日期
	EndDate    string           `json:"end_date"`    // 结束日期
}

// BenchmarkIndex 基准指数收益率
type BenchmarkIndex struct {
	Name   string  `json:"name"`   // 基准名称（如 Bitcoin/Ethereum）
	Symbol string  `json:"symbol"` // 基准标识（如 BTC/ETH）
	Return float64 `json:"return"` // 基准收益率（百分比）
	Diff   float64 `json:"diff"`   // 用户超额收益（百分比）
}

// BenchmarkDataPoint 基准数据点（时间序列）
type BenchmarkDataPoint struct {
	Date  string  `json:"date"`  // 日期（YYYY-MM-DD）
	Value float64 `json:"value"` // 归一化值（起始值为 100）
}

// BenchmarkSeries 基准时间序列
type BenchmarkSeries struct {
	Name   string               `json:"name"`   // 序列名称
	Points []BenchmarkDataPoint `json:"points"` // 数据点列表
	Return float64              `json:"return"` // 区间收益率
}
