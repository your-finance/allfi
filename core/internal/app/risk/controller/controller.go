// Package controller 风险管理模块控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "your-finance/allfi/api/v1/risk"
	"your-finance/allfi/internal/app/risk/service"
)

// Controller 风险管理控制器
type Controller struct{}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	ctrl := &Controller{}
	group.Bind(ctrl)
}

// GetMetrics 获取最新风险指标
func (c *Controller) GetMetrics(ctx context.Context, req *v1.GetMetricsReq) (res *v1.GetMetricsRes, err error) {
	metrics, err := service.Risk().GetLatestMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &v1.GetMetricsRes{
		Metrics: &v1.RiskMetrics{
			MetricDate:          metrics.MetricDate,
			PortfolioValue:      metrics.PortfolioValue,
			Var95:               metrics.Var95,
			Var99:               metrics.Var99,
			SharpeRatio:         metrics.SharpeRatio,
			SortinoRatio:        metrics.SortinoRatio,
			MaxDrawdown:         metrics.MaxDrawdown,
			MaxDrawdownDuration: metrics.MaxDrawdownDuration,
			Beta:                metrics.Beta,
			Volatility:          metrics.Volatility,
			DownsideDeviation:   metrics.DownsideDeviation,
			CalculationPeriod:   metrics.CalculationPeriod,
		},
	}

	return res, nil
}

// GetHistory 获取历史风险指标
func (c *Controller) GetHistory(ctx context.Context, req *v1.GetHistoryReq) (res *v1.GetHistoryRes, err error) {
	history, err := service.Risk().GetHistoryMetrics(ctx, req.Days)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	apiHistory := make([]*v1.RiskMetrics, 0, len(history))
	for _, m := range history {
		apiHistory = append(apiHistory, &v1.RiskMetrics{
			MetricDate:          m.MetricDate,
			PortfolioValue:      m.PortfolioValue,
			Var95:               m.Var95,
			Var99:               m.Var99,
			SharpeRatio:         m.SharpeRatio,
			SortinoRatio:        m.SortinoRatio,
			MaxDrawdown:         m.MaxDrawdown,
			MaxDrawdownDuration: m.MaxDrawdownDuration,
			Beta:                m.Beta,
			Volatility:          m.Volatility,
			DownsideDeviation:   m.DownsideDeviation,
			CalculationPeriod:   m.CalculationPeriod,
		})
	}

	res = &v1.GetHistoryRes{
		History: apiHistory,
	}

	return res, nil
}

// Calculate 手动触发风险指标计算
func (c *Controller) Calculate(ctx context.Context, req *v1.CalculateReq) (res *v1.CalculateRes, err error) {
	period := req.Period
	if period == 0 {
		period = 30
	}

	metrics, err := service.Risk().CalculateMetrics(ctx, period)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &v1.CalculateRes{
		Metrics: &v1.RiskMetrics{
			MetricDate:          metrics.MetricDate,
			PortfolioValue:      metrics.PortfolioValue,
			Var95:               metrics.Var95,
			Var99:               metrics.Var99,
			SharpeRatio:         metrics.SharpeRatio,
			SortinoRatio:        metrics.SortinoRatio,
			MaxDrawdown:         metrics.MaxDrawdown,
			MaxDrawdownDuration: metrics.MaxDrawdownDuration,
			Beta:                metrics.Beta,
			Volatility:          metrics.Volatility,
			DownsideDeviation:   metrics.DownsideDeviation,
			CalculationPeriod:   metrics.CalculationPeriod,
		},
		Message: "风险指标计算完成",
	}

	return res, nil
}
