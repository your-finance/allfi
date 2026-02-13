// Package controller 策略引擎模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	strategyApi "your-finance/allfi/api/v1/strategy"
	"your-finance/allfi/internal/app/strategy/service"
)

// Controller 策略控制器
type Controller struct{}

// List 获取策略列表
func (c *Controller) List(ctx context.Context, req *strategyApi.ListReq) (res *strategyApi.ListRes, err error) {
	return service.Strategy().List(ctx)
}

// Create 创建策略
func (c *Controller) Create(ctx context.Context, req *strategyApi.CreateReq) (res *strategyApi.CreateRes, err error) {
	return service.Strategy().Create(ctx, req.Name, req.Type, req.Config)
}

// Update 更新策略
func (c *Controller) Update(ctx context.Context, req *strategyApi.UpdateReq) (res *strategyApi.UpdateRes, err error) {
	return service.Strategy().Update(ctx, req.Id, req.Name, req.Type, req.Config, req.IsActive)
}

// Delete 删除策略
func (c *Controller) Delete(ctx context.Context, req *strategyApi.DeleteReq) (res *strategyApi.DeleteRes, err error) {
	err = service.Strategy().Delete(ctx, req.Id)
	return
}

// GetAnalysis 获取策略分析
func (c *Controller) GetAnalysis(ctx context.Context, req *strategyApi.GetAnalysisReq) (res *strategyApi.GetAnalysisRes, err error) {
	return service.Strategy().GetAnalysis(ctx, req.Id)
}

// GetRebalance 获取再平衡建议（复用 GetAnalysis 逻辑）
func (c *Controller) GetRebalance(ctx context.Context, req *strategyApi.GetRebalanceReq) (res *strategyApi.GetRebalanceRes, err error) {
	return service.Strategy().GetAnalysis(ctx, req.Id)
}

// Register 注册策略模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
