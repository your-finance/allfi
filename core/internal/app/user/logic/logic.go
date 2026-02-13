// Package logic 用户模块逻辑层入口
// 在包加载时自动注册用户服务实现
package logic

import (
	"your-finance/allfi/internal/app/user/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册用户服务
	service.RegisterUser(New())
}
