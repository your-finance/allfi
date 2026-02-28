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

// GetAPIKeys 获取 API Key 配置列表（脱敏显示）
//
// 对应路由: GET /system/apikeys
func (c *SystemController) GetAPIKeys(ctx context.Context, req *systemApi.GetAPIKeysReq) (res *systemApi.GetAPIKeysRes, err error) {
	return service.System().GetAPIKeys(ctx)
}

// UpdateAPIKey 更新 API Key
//
// 对应路由: PUT /system/apikeys
func (c *SystemController) UpdateAPIKey(ctx context.Context, req *systemApi.UpdateAPIKeyReq) (res *systemApi.UpdateAPIKeyRes, err error) {
	err = service.System().UpdateAPIKey(ctx, req.Provider, req.APIKey)
	if err != nil {
		return nil, err
	}
	return &systemApi.UpdateAPIKeyRes{Success: true}, nil
}

// DeleteAPIKey 删除 API Key
//
// 对应路由: DELETE /system/apikeys
func (c *SystemController) DeleteAPIKey(ctx context.Context, req *systemApi.DeleteAPIKeyReq) (res *systemApi.DeleteAPIKeyRes, err error) {
	err = service.System().DeleteAPIKey(ctx, req.Provider)
	if err != nil {
		return nil, err
	}
	return &systemApi.DeleteAPIKeyRes{Success: true}, nil
}

// RegisterPublic 注册免认证的系统管理路由
// 仅包含版本信息、更新检查等只读、非敏感接口
func RegisterPublic(group *ghttp.RouterGroup) {
	publicCtrl := &SystemController{}
	group.GET("/system/version", publicCtrl.GetVersion)
	group.GET("/system/update/check", publicCtrl.CheckUpdate)
	group.GET("/system/update/status", publicCtrl.GetUpdateStatus)
	group.GET("/system/update/history", publicCtrl.GetUpdateHistory)
}

// RegisterProtected 注册需认证的系统管理路由
// 包含 API Key 管理、系统更新/回滚等敏感操作
func RegisterProtected(group *ghttp.RouterGroup) {
	protectedCtrl := &SystemController{}
	group.GET("/system/apikeys", protectedCtrl.GetAPIKeys)
	group.PUT("/system/apikeys", protectedCtrl.UpdateAPIKey)
	group.DELETE("/system/apikeys", protectedCtrl.DeleteAPIKey)
	group.POST("/system/update/apply", protectedCtrl.ApplyUpdate)
	group.POST("/system/update/rollback", protectedCtrl.Rollback)
}

// Register 注册所有系统管理路由（兼容旧代码）
// 注意：建议使用 RegisterPublic + RegisterProtected 分开注册
func Register(group *ghttp.RouterGroup) {
	RegisterPublic(group)
}
