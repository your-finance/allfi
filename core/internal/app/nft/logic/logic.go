// Package logic NFT 资产模块 - Logic 层导入文件
package logic

import (
	"your-finance/allfi/internal/app/nft/service"
)

// init 在包加载时自动注册所有服务
func init() {
	// 注册 NFT 资产服务
	service.RegisterNft(New())
}
