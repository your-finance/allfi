// =================================================================================
// 市场数据服务接口定义
// 提供多链 Gas 价格等市场数据查询能力
// =================================================================================

package service

import (
	"context"

	"your-finance/allfi/internal/app/market/model"
)

// IMarket 市场数据服务接口
type IMarket interface {
	// GetGasPrice 获取多链 Gas 价格
	// 支持 Ethereum（实时 Etherscan API）、BSC 和 Polygon（硬编码估算值）
	GetGasPrice(ctx context.Context) (*model.MultiChainGasResponse, error)
}

// localMarket 市场数据服务实例（延迟注入）
var localMarket IMarket

// Market 获取市场数据服务实例
// 如果服务未注册，会触发 panic
func Market() IMarket {
	if localMarket == nil {
		panic("IMarket 服务未注册，请检查 logic/market 包的 init 函数")
	}
	return localMarket
}

// RegisterMarket 注册市场数据服务实现
// 由 logic 层在 init 函数中调用
func RegisterMarket(i IMarket) {
	localMarket = i
}
