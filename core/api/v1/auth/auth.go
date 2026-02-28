// Package auth 认证 API 定义
// 提供 PIN 设置、登录、修改、状态查询接口
package auth

import "github.com/gogf/gf/v2/frame/g"

// GetStatusReq 获取认证状态请求
type GetStatusReq struct {
	g.Meta `path:"/auth/status" method:"get" summary:"获取认证状态" tags:"认证"`
}

// GetStatusRes 获取认证状态响应
type GetStatusRes struct {
	PinSet       bool `json:"pin_set" dc:"是否已设置 PIN"`
	TwoFAEnabled bool `json:"two_fa_enabled" dc:"是否已启用 2FA"`
}

// SetupReq 首次设置 PIN 请求
type SetupReq struct {
	g.Meta `path:"/auth/setup" method:"post" summary:"首次设置 PIN" tags:"认证"`
	Pin    string `json:"pin" v:"required|length:4,20" dc:"PIN 码"`
}

// SetupRes 首次设置 PIN 响应
type SetupRes struct {
	Token string `json:"token" dc:"JWT Token（设置成功后自动登录）"`
}

// LoginReq 登录请求
type LoginReq struct {
	g.Meta `path:"/auth/login" method:"post" summary:"PIN 登录" tags:"认证"`
	Pin    string `json:"pin" v:"required" dc:"PIN 码"`
}

// LoginRes 登录响应
type LoginRes struct {
	Token       string `json:"token" dc:"JWT Token"`
	Requires2FA bool   `json:"requires_2fa" dc:"是否需要 2FA 验证"`
}

// ChangePinReq 修改 PIN 请求
type ChangePinReq struct {
	g.Meta     `path:"/auth/change" method:"post" summary:"修改 PIN" tags:"认证"`
	CurrentPin string `json:"current_pin" v:"required" dc:"当前 PIN 码"`
	NewPin     string `json:"new_pin" v:"required|length:4,20" dc:"新 PIN 码"`
}

// ChangePinRes 修改 PIN 响应
type ChangePinRes struct {
	Success bool `json:"success" dc:"是否修改成功"`
}

// Setup2FAReq 首次设置 2FA 请求
type Setup2FAReq struct {
	g.Meta `path:"/auth/2fa/setup" method:"post" summary:"获取 2FA 密钥配置" tags:"认证"`
}

// Setup2FARes 首次设置 2FA 响应
type Setup2FARes struct {
	Secret string `json:"secret" dc:"TOTP 密钥（Base32）"`
	QrUrl  string `json:"qr_url" dc:"QR 码扫描链接"`
}

// Enable2FAReq 启用 2FA 请求
type Enable2FAReq struct {
	g.Meta `path:"/auth/2fa/enable" method:"post" summary:"验证并启用 2FA" tags:"认证"`
	Code   string `json:"code" v:"required|length:6,6" dc:"6位 TOTP 验证码"`
}

// Enable2FARes 启用 2FA 响应
type Enable2FARes struct {
	Success bool `json:"success" dc:"是否启用成功"`
}

// Disable2FAReq 禁用 2FA 请求
type Disable2FAReq struct {
	g.Meta `path:"/auth/2fa/disable" method:"post" summary:"验证并禁用 2FA" tags:"认证"`
	Code   string `json:"code" v:"required|length:6,6" dc:"6位 TOTP 验证码"`
}

// Disable2FARes 禁用 2FA 响应
type Disable2FARes struct {
	Success bool `json:"success" dc:"是否禁用成功"`
}

// Verify2FAReq 登录后验证 2FA 请求
type Verify2FAReq struct {
	g.Meta `path:"/auth/2fa/verify" method:"post" summary:"验证 2FA 码" tags:"认证"`
	Code   string `json:"code" v:"required|length:6,6" dc:"6位 TOTP 验证码"`
}

// Verify2FARes 登录后验证 2FA 响应
type Verify2FARes struct {
	Success bool   `json:"success" dc:"是否验证成功"`
	Token   string `json:"token" dc:"重新颁发的完全授权 JWT Token"`
}
