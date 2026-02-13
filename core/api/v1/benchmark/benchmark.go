// Package benchmark 基准对比 API 定义
// 提供投资组合收益率与基准（BTC/ETH/S&P500）对比接口
package benchmark

import "github.com/gogf/gf/v2/frame/g"

// GetReq 获取基准对比数据请求
type GetReq struct {
	g.Meta `path:"/benchmark" method:"get" summary:"获取基准对比数据" tags:"基准对比"`
	Range  string `json:"range" in:"query" d:"30d" dc:"时间范围（7d/30d/90d/1y）"`
}

// BenchmarkSeries 基准时间序列
type BenchmarkSeries struct {
	Name   string       `json:"name" dc:"基准名称（Portfolio/BTC/ETH/S&P500）"`
	Points []DataPoint  `json:"points" dc:"数据点列表"`
	Return float64      `json:"return" dc:"区间收益率"`
}

// DataPoint 数据点
type DataPoint struct {
	Date  string  `json:"date" dc:"日期"`
	Value float64 `json:"value" dc:"归一化值（起始值为 100）"`
}

// GetRes 获取基准对比数据响应
type GetRes struct {
	Series    []BenchmarkSeries `json:"series" dc:"基准时间序列列表"`
	Range     string            `json:"range" dc:"时间范围"`
	StartDate string            `json:"start_date" dc:"开始日期"`
	EndDate   string            `json:"end_date" dc:"结束日期"`
}
