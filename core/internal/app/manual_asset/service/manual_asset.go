// Package service 手动资产模块 - 服务接口定义
// 定义手动资产（银行账户/现金/股票/基金）管理的所有服务方法
package service

import (
	"context"

	manualAssetApi "your-finance/allfi/api/v1/manual_asset"
)

// IManualAsset 手动资产服务接口
type IManualAsset interface {
	// ListManualAssets 获取手动资产列表
	ListManualAssets(ctx context.Context, userID int) ([]manualAssetApi.ManualAssetItem, error)

	// CreateManualAsset 添加手动资产
	CreateManualAsset(ctx context.Context, userID int, req *manualAssetApi.CreateReq) (*manualAssetApi.ManualAssetItem, error)

	// UpdateManualAsset 更新手动资产
	UpdateManualAsset(ctx context.Context, req *manualAssetApi.UpdateReq) (*manualAssetApi.ManualAssetItem, error)

	// DeleteManualAsset 删除手动资产
	DeleteManualAsset(ctx context.Context, assetID int) error
}

var localManualAsset IManualAsset

// ManualAsset 获取手动资产服务实例
func ManualAsset() IManualAsset {
	if localManualAsset == nil {
		panic("IManualAsset 服务未注册，请检查 logic/manual_asset 包的 init 函数")
	}
	return localManualAsset
}

// RegisterManualAsset 注册手动资产服务实现
// 由 logic 层在 init 函数中调用
func RegisterManualAsset(i IManualAsset) {
	localManualAsset = i
}
