// Package service 趋势预测模块 - Service 层接口定义
package service

import (
	"context"

	"your-finance/allfi/internal/app/forecast/model"
)

// IForecast 趋势预测服务接口
type IForecast interface {
	// GetForecast 获取趋势预测结果
	// targetValue: 目标资产总值
	// currency: 计价货币
	// 使用 90 天历史快照做线性回归，预测到达目标值的时间
	GetForecast(ctx context.Context, targetValue float64, currency string) (*model.ForecastResult, error)
}

var localForecast IForecast

// Forecast 获取趋势预测服务实例
func Forecast() IForecast {
	if localForecast == nil {
		panic("IForecast 服务未注册，请检查 logic/forecast 包的 init 函数")
	}
	return localForecast
}

// RegisterForecast 注册趋势预测服务实现
// 由 logic 层在 init 函数中调用
func RegisterForecast(i IForecast) {
	localForecast = i
}
