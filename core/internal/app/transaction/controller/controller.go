// Package controller 交易记录模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	transactionApi "your-finance/allfi/api/v1/transaction"
	"your-finance/allfi/internal/app/transaction/service"
)

// Controller 交易记录控制器
type Controller struct{}

// List 获取交易记录列表
func (c *Controller) List(ctx context.Context, req *transactionApi.ListReq) (res *transactionApi.ListRes, err error) {
	return service.Transaction().List(ctx, req.Page, req.PageSize, req.Source, req.Type, req.Start, req.End, req.Cursor)
}

// Sync 触发交易记录同步
func (c *Controller) Sync(ctx context.Context, req *transactionApi.SyncReq) (res *transactionApi.SyncRes, err error) {
	return service.Transaction().Sync(ctx)
}

// GetStats 获取交易统计
func (c *Controller) GetStats(ctx context.Context, req *transactionApi.GetStatsReq) (res *transactionApi.GetStatsRes, err error) {
	return service.Transaction().GetStats(ctx)
}

// GetSyncSettings 获取同步设置
func (c *Controller) GetSyncSettings(ctx context.Context, req *transactionApi.GetSyncSettingsReq) (res *transactionApi.GetSyncSettingsRes, err error) {
	return service.Transaction().GetSyncSettings(ctx)
}

// UpdateSyncSettings 更新同步设置
func (c *Controller) UpdateSyncSettings(ctx context.Context, req *transactionApi.UpdateSyncSettingsReq) (res *transactionApi.UpdateSyncSettingsRes, err error) {
	return service.Transaction().UpdateSyncSettings(ctx, req.AutoSync, req.SyncInterval)
}

// Register 注册交易记录模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
