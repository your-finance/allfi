// Package logic 目标追踪模块 - Logic 层导入文件
// 在包加载时自动注册目标追踪服务实现
package logic

import (
	"your-finance/allfi/internal/app/goal/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册目标追踪服务
	service.RegisterGoal(New())
}
