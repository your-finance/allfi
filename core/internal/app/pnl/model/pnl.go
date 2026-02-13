// Package model 盈亏分析模块 - 数据传输对象
package model

// DailyPnLPoint 每日盈亏数据点
type DailyPnLPoint struct {
	Date       string  `json:"date"`        // 日期 YYYY-MM-DD
	StartValue float64 `json:"start_value"` // 当日开始总值
	EndValue   float64 `json:"end_value"`   // 当日结束总值
	PnL        float64 `json:"pnl"`         // 盈亏金额
	PnLPercent float64 `json:"pnl_percent"` // 盈亏百分比
}

// PnLSummary 盈亏汇总
type PnLSummary struct {
	TotalPnL        float64 `json:"total_pnl"`         // 总盈亏金额
	TotalPnLPercent float64 `json:"total_pnl_percent"` // 总盈亏比例
	PnL7d           float64 `json:"pnl_7d"`            // 7 日盈亏
	PnL30d          float64 `json:"pnl_30d"`           // 30 日盈亏
	PnL90d          float64 `json:"pnl_90d"`           // 90 日盈亏
	BestDay         string  `json:"best_day"`          // 最佳单日日期
	WorstDay        string  `json:"worst_day"`         // 最差单日日期
	BestDayPnL      float64 `json:"best_day_pnl"`      // 最佳单日盈亏
	WorstDayPnL     float64 `json:"worst_day_pnl"`     // 最差单日盈亏
}

// PnLPeriod 某时段盈亏
type PnLPeriod struct {
	PnL        float64 `json:"pnl"`         // 盈亏金额
	PnLPercent float64 `json:"pnl_percent"` // 盈亏百分比
	StartValue float64 `json:"start_value"` // 起始价值
	EndValue   float64 `json:"end_value"`   // 结束价值
}
