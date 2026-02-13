// Package logic 认证模块 - Logic 层导入文件
// 在包加载时自动注册认证服务实现
package logic

import (
	"your-finance/allfi/internal/app/auth/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册认证服务
	service.RegisterAuth(New())
}
