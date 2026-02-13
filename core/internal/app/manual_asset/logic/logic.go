// Package logic 手动资产模块 - Logic 层导入文件
// 在包加载时自动注册手动资产服务实现
package logic

import (
	"your-finance/allfi/internal/app/manual_asset/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册手动资产服务
	service.RegisterManualAsset(New())
}
