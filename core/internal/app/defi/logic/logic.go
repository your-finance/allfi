// Package logic DeFi 仓位模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/defi/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册 DeFi 仓位服务
	service.RegisterDefi(New())
}
