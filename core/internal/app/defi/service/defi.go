// Package service DeFi 仓位模块 - Service 层接口定义
package service

import (
	"context"

	defiApi "your-finance/allfi/api/v1/defi"
	"your-finance/allfi/internal/app/defi/model"
)

// IDefi DeFi 仓位服务接口
type IDefi interface {
	// GetPositions 获取用户 DeFi 仓位列表
	// chain: 指定链名（ethereum/polygon 等），空字符串查询所有链
	// protocol: 指定协议名（lido/aave 等），空字符串查询所有协议
	// 返回仓位列表和总价值
	GetPositions(ctx context.Context, chain string, protocol string) (positions []*model.PositionItem, totalValue float64, err error)

	// GetProtocols 获取支持的 DeFi 协议列表
	GetProtocols(ctx context.Context) ([]*model.ProtocolItem, error)

	// GetStats 获取 DeFi 统计（按协议/链/类型分组聚合）
	GetStats(ctx context.Context) (*defiApi.GetStatsRes, error)
}

var localDefi IDefi

// Defi 获取 DeFi 仓位服务实例
func Defi() IDefi {
	if localDefi == nil {
		panic("IDefi 服务未注册，请检查 logic/defi 包的 init 函数")
	}
	return localDefi
}

// RegisterDefi 注册 DeFi 仓位服务实现
// 由 logic 层在 init 函数中调用
func RegisterDefi(i IDefi) {
	localDefi = i
}
