// Package controller Gas 优化控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	v1 "your-finance/allfi/api/v1/gas"
	"your-finance/allfi/internal/app/gas/service"
)

// Controller Gas 优化控制器
type Controller struct{}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	c := &Controller{}
	group.Group("/gas", func(group *ghttp.RouterGroup) {
		group.Bind(
			c.GetCurrent,
			c.GetHistory,
			c.GetRecommendation,
			c.GetForecast,
		)
	})
}

// GetCurrent 获取当前 Gas 价格
func (c *Controller) GetCurrent(ctx context.Context, req *v1.GetCurrentReq) (res *v1.GetCurrentRes, err error) {
	return service.Gas().GetCurrent(ctx)
}

// GetHistory 获取 Gas 价格历史
func (c *Controller) GetHistory(ctx context.Context, req *v1.GetHistoryReq) (res *v1.GetHistoryRes, err error) {
	return service.Gas().GetHistory(ctx, req)
}

// GetRecommendation 获取最佳交易时间推荐
func (c *Controller) GetRecommendation(ctx context.Context, req *v1.GetRecommendationReq) (res *v1.GetRecommendationRes, err error) {
	return service.Gas().GetRecommendation(ctx, req)
}

// GetForecast 获取 Gas 价格预测
func (c *Controller) GetForecast(ctx context.Context, req *v1.GetForecastReq) (res *v1.GetForecastRes, err error) {
	return service.Gas().GetForecast(ctx, req)
}
