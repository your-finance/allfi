// Package logic 手动资产业务逻辑单元测试
package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManualAssetNew(t *testing.T) {
	logic := New()
	assert.NotNil(t, logic)
	assert.IsType(t, &sManualAsset{}, logic)
}

func TestManualAssetStruct(t *testing.T) {
	assert.NotNil(t, &sManualAsset{})
}

func TestManualAssetSupportedAssetTypes(t *testing.T) {
	supportedTypes := []string{"bank", "cash", "stock", "fund"}

	for _, assetType := range supportedTypes {
		assert.Contains(t, supportedTypes, assetType)
	}
}

func TestConvertToUSDReturnsOriginalAmountForDollarCurrencies(t *testing.T) {
	logic := &sManualAsset{}

	testCases := []string{"", "USD", "USDC", "USDT"}
	for _, currency := range testCases {
		assert.Equalf(t, 123.45, logic.convertToUSD(t.Context(), 123.45, currency), "currency=%s", currency)
	}
}
