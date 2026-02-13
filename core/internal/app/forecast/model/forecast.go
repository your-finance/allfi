// Package model 趋势预测模块 - 数据传输对象
package model

import "time"

// ForecastResult 趋势预测结果
type ForecastResult struct {
	CurrentValue  float64    `json:"current_value"`  // 当前资产总值
	TargetValue   float64    `json:"target_value"`   // 目标值
	Currency      string     `json:"currency"`       // 计价货币
	DailyGrowth   float64    `json:"daily_growth"`   // 日均增长额
	GrowthRate    float64    `json:"growth_rate"`    // 日均增长率（%）
	EstimatedDate *time.Time `json:"estimated_date"` // 预计达成日期（nil 表示无法预测）
	DaysToTarget  int        `json:"days_to_target"` // 距离目标天数
	Confidence    float64    `json:"confidence"`     // 置信度（R²，0-1）
	Trend         string     `json:"trend"`          // 趋势方向（up/down/flat）
	DataPoints    int        `json:"data_points"`    // 用于预测的数据点数量
}
