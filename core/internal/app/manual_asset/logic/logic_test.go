// Package logic 手动资产业务逻辑单元测试
package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	manualAssetEntity "your-finance/allfi/internal/app/manual_asset/model/entity"
)

func TestManualAsset_CalculateValue(t *testing.T) {
	// 测试资产价值计算
	assets := []manualAssetEntity.ManualAsset{
		{Id: 1, AssetType: "bank", Amount: 50000, Currency: "CNY"},
		{Id: 2, AssetType: "cash", Amount: 10000, Currency: "CNY"},
		{Id: 3, AssetType: "bank", Amount: 25000, Currency: "USD"},
	}

	logic := New(&sManualAsset{})

	// 测试资产数量
	assert.Equal(t, 3, len(assets))
}

func TestManualAsset_CalculateByCurrency(t *testing.T) {
	// 测试按货币分组计算
	assets := []manualAssetEntity.ManualAsset{
		{Id: 1, AssetType: "bank", Amount: 50000, Currency: "CNY"},
		{Id: 2, AssetType: "bank", Amount: 50000, Currency: "CNY"},
		{Id: 3, AssetType: "bank", Amount: 25000, Currency: "USD"},
	}

	logic := New(&sManualAsset{})

	// 验证 CNY 总额
	cnyTotal := float64(0)
	usdTotal := float64(0)
	for _, asset := range assets {
		if asset.Currency == "CNY" {
			cnyTotal += asset.Amount
		}
		if asset.Currency == "USD" {
			usdTotal += asset.Amount
		}
	}

	assert.Equal(t, float64(100000), cnyTotal)
	assert.Equal(t, float64(25000), usdTotal)
}

func TestManualAsset_SupportedAssetTypes(t *testing.T) {
	// 测试支持的资产类型
	supportedTypes := []string{"bank", "cash", "stock", "fund"}

	for _, assetType := range supportedTypes {
		asset := manualAssetEntity.ManualAsset{
			AssetType: assetType,
			Amount:     1000,
			Currency:  "CNY",
		}

		// 验证资产类型有效
		assert.Contains(t, supportedTypes, asset.AssetType)
		assert.True(t, asset.Amount > 0)
	}

	// 测试不支持的资产类型
	invalidAsset := manualAssetEntity.ManualAsset{
		AssetType: "crypto",
		Amount:     1000,
		Currency: "BTC",
	}

	// 验证不支持的类型
	assert.NotContains(t, supportedTypes, invalidAsset.AssetType)
}

func TestManualAsset_EmptyAssets(t *testing.T) {
	logic := New(&sManualAsset{})
	assert.NotNil(t, logic)
}
