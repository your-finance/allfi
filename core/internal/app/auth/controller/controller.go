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

// Setup2FA 获取 2FA 密钥与二维码
func (c *Controller) Setup2FA(ctx context.Context, req *authApi.Setup2FAReq) (res *authApi.Setup2FARes, err error) {
	return service.Auth().Setup2FA(ctx)
}

// Enable2FA 启用 2FA
func (c *Controller) Enable2FA(ctx context.Context, req *authApi.Enable2FAReq) (res *authApi.Enable2FARes, err error) {
	return service.Auth().Enable2FA(ctx, req.Code)
}

// Disable2FA 禁用 2FA
func (c *Controller) Disable2FA(ctx context.Context, req *authApi.Disable2FAReq) (res *authApi.Disable2FARes, err error) {
	return service.Auth().Disable2FA(ctx, req.Code)
}

// Verify2FA 验证 2FA 发放完整 Token
func (c *Controller) Verify2FA(ctx context.Context, req *authApi.Verify2FAReq) (res *authApi.Verify2FARes, err error) {
	return service.Auth().Verify2FA(ctx, req.Code)
}

// SwitchType 切换密码类型
func (c *Controller) SwitchType(ctx context.Context, req *authApi.SwitchTypeReq) (res *authApi.SwitchTypeRes, err error) {
	return service.Auth().SwitchType(ctx, req.CurrentPassword, req.NewType, req.NewPassword, req.TwoFACode)
}

// Register 注册认证模块路由（公开接口，无需认证）
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
