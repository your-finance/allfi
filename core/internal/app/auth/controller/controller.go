// Package controller 认证模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	authApi "your-finance/allfi/api/v1/auth"
	"your-finance/allfi/internal/app/auth/service"
)

// Controller 认证控制器
type Controller struct{}

// GetStatus 获取认证状态
func (c *Controller) GetStatus(ctx context.Context, req *authApi.GetStatusReq) (res *authApi.GetStatusRes, err error) {
	return service.Auth().GetStatus(ctx)
}

// Setup 首次设置 PIN
func (c *Controller) Setup(ctx context.Context, req *authApi.SetupReq) (res *authApi.SetupRes, err error) {
	return service.Auth().Setup(ctx, req.Pin)
}

// Login PIN 登录
func (c *Controller) Login(ctx context.Context, req *authApi.LoginReq) (res *authApi.LoginRes, err error) {
	return service.Auth().Login(ctx, req.Pin)
}

// ChangePin 修改 PIN
func (c *Controller) ChangePin(ctx context.Context, req *authApi.ChangePinReq) (res *authApi.ChangePinRes, err error) {
	return service.Auth().ChangePin(ctx, req.CurrentPin, req.NewPin)
}

// Register 注册认证模块路由（公开接口，无需认证）
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
