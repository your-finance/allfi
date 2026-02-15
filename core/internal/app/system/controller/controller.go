// Package controller 系统管理模块路由注册
// 提供版本信息查询、在线更新、回滚等系统管理端点
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	systemApi "your-finance/allfi/api/v1/system"
	"your-finance/allfi/internal/app/system/service"
)

// SystemController 系统管理控制器
type SystemController struct{}

// GetVersion 获取当前版本信息
//
// 对应路由: GET /system/version
func (c *SystemController) GetVersion(ctx context.Context, req *systemApi.GetVersionReq) (res *systemApi.GetVersionRes, err error) {
	return service.System().GetVersion(ctx)
}

// CheckUpdate 检查 GitHub Releases 是否有新版本
//
// 对应路由: GET /system/update/check
func (c *SystemController) CheckUpdate(ctx context.Context, req *systemApi.CheckUpdateReq) (res *systemApi.CheckUpdateRes, err error) {
	return service.System().CheckUpdate(ctx)
}

// ApplyUpdate 执行版本更新
//
// 对应路由: POST /system/update/apply
func (c *SystemController) ApplyUpdate(ctx context.Context, req *systemApi.ApplyUpdateReq) (res *systemApi.ApplyUpdateRes, err error) {
	return service.System().ApplyUpdate(ctx, req.TargetVersion)
}

// Rollback 版本回滚
//
// 对应路由: POST /system/update/rollback
func (c *SystemController) Rollback(ctx context.Context, req *systemApi.RollbackReq) (res *systemApi.RollbackRes, err error) {
	return service.System().Rollback(ctx, req.TargetVersion)
}

// GetUpdateStatus 获取更新/回滚操作进度
//
// 对应路由: GET /system/update/status
func (c *SystemController) GetUpdateStatus(ctx context.Context, req *systemApi.GetUpdateStatusReq) (res *systemApi.GetUpdateStatusRes, err error) {
	return service.System().GetUpdateStatus(ctx)
}

// GetUpdateHistory 获取历史更新记录
//
// 对应路由: GET /system/update/history
func (c *SystemController) GetUpdateHistory(ctx context.Context, req *systemApi.GetUpdateHistoryReq) (res *systemApi.GetUpdateHistoryRes, err error) {
	return service.System().GetUpdateHistory(ctx)
}

// Register 注册系统管理路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&SystemController{})
}
