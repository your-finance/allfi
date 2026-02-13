// Package logic 报告模块 - Logic 层导入文件
// 在包加载时自动注册报告服务实现
package logic

import (
	"your-finance/allfi/internal/app/report/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册报告服务
	service.RegisterReport(New())
}
