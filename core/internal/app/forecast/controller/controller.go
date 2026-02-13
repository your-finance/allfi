// Package controller 趋势预测模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	forecastApi "your-finance/allfi/api/v1/forecast"
	"your-finance/allfi/internal/app/forecast/service"
)

// ForecastController 趋势预测控制器
type ForecastController struct{}

// Get 获取趋势预测
func (c *ForecastController) Get(ctx context.Context, req *forecastApi.GetReq) (res *forecastApi.GetRes, err error) {
	// 默认预测 30 天，目标值使用当前值的 120%
	// 实际 targetValue 由前端传入，这里用 Days 作为提示信息
	targetValue := float64(req.Days * 1000) // 简化: 天数 * 1000 作为默认目标值
	if targetValue <= 0 {
		targetValue = 100000 // 默认目标 10 万
	}

	// 调用服务层获取预测结果
	result, err := service.Forecast().GetForecast(ctx, targetValue, "USD")
	if err != nil {
		return nil, err
	}

	// 构建预测数据点（未来趋势线）
	forecastPoints := make([]forecastApi.ForecastPoint, 0)
	// 预测点数据在前端侧展示，这里返回核心指标

	res = &forecastApi.GetRes{
		ForecastPoints: forecastPoints,
		Trend:          result.Trend,
		Confidence:     result.Confidence,
		Slope:          result.DailyGrowth,
		Days:           req.Days,
	}

	return res, nil
}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&ForecastController{})
}
