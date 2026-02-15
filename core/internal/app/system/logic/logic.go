// Package logic 系统管理模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/system/service"
)

// init 在包加载时自动注册系统管理服务
func init() {
	service.RegisterSystem(New())
}
