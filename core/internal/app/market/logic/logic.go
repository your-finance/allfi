// Package logic 市场数据模块 - Logic 层导入文件
// init 函数自动注册服务实现
package logic

import (
	"your-finance/allfi/internal/app/market/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册市场数据服务
	service.RegisterMarket(New())
}
