// Package controller 资产模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	assetApi "your-finance/allfi/api/v1/asset"
	"your-finance/allfi/internal/app/asset/service"
)

// Controller 资产控制器
type Controller struct{}

// GetSummary 获取资产概览
func (c *Controller) GetSummary(ctx context.Context, req *assetApi.GetSummaryReq) (res *assetApi.GetSummaryRes, err error) {
	return service.Asset().GetSummary(ctx, req.Currency)
}

// GetDetails 获取资产明细列表
func (c *Controller) GetDetails(ctx context.Context, req *assetApi.GetDetailsReq) (res *assetApi.GetDetailsRes, err error) {
	return service.Asset().GetDetails(ctx, req.SourceType, req.Currency)
}

// GetHistory 获取资产历史趋势
func (c *Controller) GetHistory(ctx context.Context, req *assetApi.GetHistoryReq) (res *assetApi.GetHistoryRes, err error) {
	return service.Asset().GetHistory(ctx, req.Days, req.Currency)
}

// Refresh 强制刷新所有资产
func (c *Controller) Refresh(ctx context.Context, req *assetApi.RefreshReq) (res *assetApi.RefreshRes, err error) {
	return service.Asset().RefreshAll(ctx)
}

// Register 注册资产模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
