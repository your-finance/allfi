// Package service 资产模块 - 服务接口定义
// 定义资产概览、明细、历史趋势、刷新等功能
package service

import (
	"context"

	assetApi "your-finance/allfi/api/v1/asset"
)

// IAsset 资产服务接口
type IAsset interface {
	// GetSummary 获取资产概览
	// 聚合 CEX/区块链/手动三个来源的资产，按指定货币计价
	GetSummary(ctx context.Context, currency string) (*assetApi.GetSummaryRes, error)

	// GetDetails 获取资产明细列表
	// 可按来源类型筛选（cex/blockchain/manual）
	GetDetails(ctx context.Context, sourceType string, currency string) (*assetApi.GetDetailsRes, error)

	// GetHistory 获取资产历史趋势
	// 返回指定天数内的资产快照
	GetHistory(ctx context.Context, days int, currency string) (*assetApi.GetHistoryRes, error)

	// RefreshAll 强制刷新所有资产
	// 从 CEX API 和区块链查询最新余额
	RefreshAll(ctx context.Context) (*assetApi.RefreshRes, error)
}

var localAsset IAsset

// Asset 获取资产服务实例
func Asset() IAsset {
	if localAsset == nil {
		panic("IAsset 服务未注册，请检查 logic/asset 包的 init 函数")
	}
	return localAsset
}

// RegisterAsset 注册资产服务实现
// 由 logic 层在 init 函数中调用
func RegisterAsset(i IAsset) {
	localAsset = i
}
