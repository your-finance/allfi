// Package logic 通知模块 - Logic 层导入文件
// 在包加载时自动注册通知服务实现
package logic

import (
	"your-finance/allfi/internal/app/notification/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册通知服务
	service.RegisterNotification(New())
}
