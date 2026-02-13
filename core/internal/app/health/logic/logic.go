// Package logic 健康检查模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/health/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册健康检查服务
	service.RegisterHealth(New())
}
