// Package service 认证模块 - 服务接口定义
// PIN 码认证：bcrypt 哈希存储 + JWT Token + 锁定保护
package service

import (
	"context"

	authApi "your-finance/allfi/api/v1/auth"
)

// IAuth 认证服务接口
type IAuth interface {
	// GetStatus 获取认证状态（PIN 是否已设置）
	GetStatus(ctx context.Context) (*authApi.GetStatusRes, error)

	// Setup 首次设置 PIN（4-8位数字）
	Setup(ctx context.Context, pin string) (*authApi.SetupRes, error)

	// Login 验证 PIN 返回 JWT Token
	Login(ctx context.Context, pin string) (*authApi.LoginRes, error)

	// ChangePin 修改 PIN（需验证旧 PIN）
	ChangePin(ctx context.Context, currentPin string, newPin string) (*authApi.ChangePinRes, error)
}

var localAuth IAuth

// Auth 获取认证服务实例
func Auth() IAuth {
	if localAuth == nil {
		panic("IAuth 服务未注册，请检查 logic/auth 包的 init 函数")
	}
	return localAuth
}

// RegisterAuth 注册认证服务实现
// 由 logic 层在 init 函数中调用
func RegisterAuth(i IAuth) {
	localAuth = i
}
