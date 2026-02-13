// Package logic 交易记录模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/transaction/service"
)

// init 在包加载时自动注册交易记录服务
func init() {
	service.RegisterTransaction(New())
}
