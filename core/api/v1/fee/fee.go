// Package fee 费用分析 API 定义
// 提供交易费用（手续费 + Gas）分析接口
package fee

import "github.com/gogf/gf/v2/frame/g"

// GetAnalyticsReq 获取费用分析请求
type GetAnalyticsReq struct {
	g.Meta   `path:"/analytics/fees" method:"get" summary:"获取费用分析" tags:"费用分析"`
	Range    string `json:"range" in:"query" d:"30d" dc:"时间范围（7d/30d/90d/1y）"`
	Currency string `json:"currency" in:"query" d:"USD" dc:"计价货币"`
}

// FeeBreakdown 费用明细
type FeeBreakdown struct {
	Source   string  `json:"source" dc:"来源（交易所名/链名）"`
	Type     string  `json:"type" dc:"费用类型（trading_fee/gas_fee）"`
	Amount   float64 `json:"amount" dc:"费用金额"`
	Currency string  `json:"currency" dc:"计价货币"`
	Count    int     `json:"count" dc:"交易次数"`
}

// DailyFee 每日费用
type DailyFee struct {
	Date   string  `json:"date" dc:"日期"`
	Amount float64 `json:"amount" dc:"当日费用总额"`
}

// GetAnalyticsRes 获取费用分析响应
type GetAnalyticsRes struct {
	TotalFees    float64        `json:"total_fees" dc:"费用总计"`
	TradingFees  float64        `json:"trading_fees" dc:"交易手续费总计"`
	GasFees      float64        `json:"gas_fees" dc:"Gas 费用总计"`
	Currency     string         `json:"currency" dc:"计价货币"`
	Breakdown    []FeeBreakdown `json:"breakdown" dc:"费用明细"`
	DailyTrend   []DailyFee     `json:"daily_trend" dc:"每日费用趋势"`
}
