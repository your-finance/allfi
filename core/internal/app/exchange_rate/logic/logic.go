// Package logic 汇率模块 - Logic 层导入文件
// 在包加载时自动注册汇率服务实现
package logic

import (
	"your-finance/allfi/internal/app/exchange_rate/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册汇率服务
	service.RegisterExchangeRate(New())
}
