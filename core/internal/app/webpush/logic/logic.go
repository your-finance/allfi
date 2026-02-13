// Package logic WebPush 推送模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/webpush/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册 WebPush 推送服务
	service.RegisterWebpush(New())
}
