// Package logic 盈亏分析模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/pnl/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册盈亏分析服务
	service.RegisterPnl(New())
}
