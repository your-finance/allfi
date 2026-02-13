// Package logic 资产模块 - Logic 层导入文件
// 在包加载时自动注册资产服务实现
package logic

import (
	"your-finance/allfi/internal/app/asset/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册资产服务
	service.RegisterAsset(New())
}
