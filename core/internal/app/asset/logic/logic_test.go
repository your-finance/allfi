// Package logic 资产总览业务逻辑单元测试
package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetLogicNew(t *testing.T) {
	logic := New()
	assert.NotNil(t, logic)
	assert.IsType(t, &sAsset{}, logic)
}

func TestSAssetStruct(t *testing.T) {
	assert.NotNil(t, &sAsset{})
}

func TestGetConversionRateReturnsOneForDefaultCurrencies(t *testing.T) {
	asset := &sAsset{}
	ctx := context.Background()

	testCases := []string{"", "USD", "USDC", "USDT"}
	for _, currency := range testCases {
		rate := asset.getConversionRate(ctx, currency)
		assert.Equalf(t, 1.0, rate, "currency=%s", currency)
	}
}
