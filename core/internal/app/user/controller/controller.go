// Package controller 用户模块控制器
// 绑定用户设置和缓存管理 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	userApi "your-finance/allfi/api/v1/user"
	"your-finance/allfi/internal/app/user/service"
)

// Controller 用户控制器
type Controller struct{}

// Register 注册用户模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// GetSettings 获取用户设置
func (c *Controller) GetSettings(ctx context.Context, req *userApi.GetSettingsReq) (res *userApi.GetSettingsRes, err error) {
	return service.User().GetSettings(ctx)
}

// UpdateSettings 更新用户设置
func (c *Controller) UpdateSettings(ctx context.Context, req *userApi.UpdateSettingsReq) (res *userApi.UpdateSettingsRes, err error) {
	err = service.User().UpdateSettings(ctx, req.Settings)
	if err != nil {
		return nil, err
	}
	return &userApi.UpdateSettingsRes{Message: "设置已更新"}, nil
}

// ResetSettings 重置用户设置
func (c *Controller) ResetSettings(ctx context.Context, req *userApi.ResetSettingsReq) (res *userApi.ResetSettingsRes, err error) {
	err = service.User().ResetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return &userApi.ResetSettingsRes{Message: "设置已重置"}, nil
}

// ClearCache 清除缓存
func (c *Controller) ClearCache(ctx context.Context, req *userApi.ClearCacheReq) (res *userApi.ClearCacheRes, err error) {
	err = service.User().ClearCache(ctx)
	if err != nil {
		return nil, err
	}
	return &userApi.ClearCacheRes{Message: "缓存已清除"}, nil
}

// ExportData 导出用户数据
func (c *Controller) ExportData(ctx context.Context, req *userApi.ExportDataReq) (res *userApi.ExportDataRes, err error) {
	return service.User().ExportData(ctx)
}
