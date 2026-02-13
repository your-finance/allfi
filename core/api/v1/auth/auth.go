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
	PinSet bool `json:"pin_set" dc:"是否已设置 PIN"`
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
	Token string `json:"token" dc:"JWT Token"`
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
