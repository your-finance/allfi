// Package logic 钱包模块 - Logic 层导入文件
// 在包加载时自动注册钱包服务实现
package logic

import (
	"your-finance/allfi/internal/app/wallet/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册钱包服务
	service.RegisterWallet(New())
}
