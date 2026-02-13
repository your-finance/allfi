// Package forecast 趋势预测 API 定义
// 提供基于线性回归的资产趋势预测接口
package forecast

import "github.com/gogf/gf/v2/frame/g"

// GetReq 获取趋势预测请求
type GetReq struct {
	g.Meta `path:"/analytics/forecast" method:"get" summary:"获取趋势预测" tags:"趋势预测"`
	Days   int `json:"days" in:"query" d:"30" dc:"预测天数"`
}

// ForecastPoint 预测数据点
type ForecastPoint struct {
	Date  string  `json:"date" dc:"日期"`
	Value float64 `json:"value" dc:"预测值"`
	Lower float64 `json:"lower" dc:"预测下界"`
	Upper float64 `json:"upper" dc:"预测上界"`
}

// GetRes 获取趋势预测响应
type GetRes struct {
	ForecastPoints []ForecastPoint `json:"forecast_points" dc:"预测数据点"`
	Trend          string          `json:"trend" dc:"趋势方向（up/down/stable）"`
	Confidence     float64         `json:"confidence" dc:"置信度（R² 值，0-1）"`
	Slope          float64         `json:"slope" dc:"斜率（每日变化量）"`
	Days           int             `json:"days" dc:"预测天数"`
}
