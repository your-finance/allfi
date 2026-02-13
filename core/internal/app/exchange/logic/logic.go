// Package logic 交易所模块 - Logic 层导入文件
// 在包加载时自动注册交易所服务实现
package logic

import (
	"your-finance/allfi/internal/app/exchange/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册交易所服务
	service.RegisterExchange(New())
}
