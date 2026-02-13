// Package pnl 盈亏分析 API 定义
// 提供每日盈亏和盈亏汇总接口
package pnl

import "github.com/gogf/gf/v2/frame/g"

// GetDailyReq 获取每日盈亏请求
type GetDailyReq struct {
	g.Meta `path:"/analytics/pnl/daily" method:"get" summary:"获取每日盈亏" tags:"盈亏分析"`
	Days   int `json:"days" in:"query" d:"30" dc:"查询天数"`
}

// DailyPnLItem 每日盈亏条目
type DailyPnLItem struct {
	Date       string  `json:"date" dc:"日期"`
	PnL        float64 `json:"pnl" dc:"当日盈亏金额"`
	PnLPercent float64 `json:"pnl_percent" dc:"当日盈亏比例"`
	TotalValue float64 `json:"total_value" dc:"当日总资产价值"`
}

// GetDailyRes 获取每日盈亏响应
type GetDailyRes struct {
	Daily []DailyPnLItem `json:"daily" dc:"每日盈亏列表"`
}

// GetSummaryReq 获取盈亏汇总请求
type GetSummaryReq struct {
	g.Meta `path:"/analytics/pnl/summary" method:"get" summary:"获取盈亏汇总" tags:"盈亏分析"`
}

// GetSummaryRes 获取盈亏汇总响应
type GetSummaryRes struct {
	TotalPnL       float64 `json:"total_pnl" dc:"总盈亏金额"`
	TotalPnLPercent float64 `json:"total_pnl_percent" dc:"总盈亏比例"`
	PnL7d          float64 `json:"pnl_7d" dc:"7 日盈亏"`
	PnL30d         float64 `json:"pnl_30d" dc:"30 日盈亏"`
	PnL90d         float64 `json:"pnl_90d" dc:"90 日盈亏"`
	BestDay        string  `json:"best_day" dc:"最佳单日"`
	WorstDay       string  `json:"worst_day" dc:"最差单日"`
	BestDayPnL     float64 `json:"best_day_pnl" dc:"最佳单日盈亏"`
	WorstDayPnL    float64 `json:"worst_day_pnl" dc:"最差单日盈亏"`
}
