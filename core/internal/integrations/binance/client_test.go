// Package binance Binance 客户端测试
package binance

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"your-finance/allfi/internal/integrations"
)

// TestMergeBalancesBySymbol_合并同一币种余额 验证同一币种在不同账户类型的余额正确合并
func TestMergeBalancesBySymbol_合并同一币种余额(t *testing.T) {
	balances := []integrations.Balance{
		{Symbol: "USDT", Name: "USDT", Free: 1000, Total: 1000, AssetType: "spot"},
		{Symbol: "USDT", Name: "USDT (合约)", Free: 500, Locked: 100, Total: 500, AssetType: "futures"},
		{Symbol: "BTC", Name: "BTC", Free: 1.0, Total: 1.0, AssetType: "spot"},
		{Symbol: "USDT", Name: "USDT (杠杆)", Free: 200, Total: 200, AssetType: "margin"},
	}

	merged := MergeBalancesBySymbol(balances)

	// 应该只有 2 种币
	assert.Len(t, merged, 2, "合并后应该只有 2 种币种")

	// 找到 USDT 和 BTC
	usdtFound := false
	btcFound := false
	for _, b := range merged {
		switch b.Symbol {
		case "USDT":
			usdtFound = true
			assert.Equal(t, 1700.0, b.Free, "USDT Free 应为 1000+500+200=1700")
			assert.Equal(t, 100.0, b.Locked, "USDT Locked 应为 100")
			assert.Equal(t, 1700.0, b.Total, "USDT Total 应为 1000+500+200=1700")
		case "BTC":
			btcFound = true
			assert.Equal(t, 1.0, b.Free, "BTC Free 应为 1.0")
			assert.Equal(t, 1.0, b.Total, "BTC Total 应为 1.0")
		}
	}
	assert.True(t, usdtFound, "应包含 USDT")
	assert.True(t, btcFound, "应包含 BTC")
}

// TestMergeBalancesBySymbol_空列表 验证空列表不会出错
func TestMergeBalancesBySymbol_空列表(t *testing.T) {
	merged := MergeBalancesBySymbol(nil)
	assert.Empty(t, merged, "空输入应返回空结果")

	merged = MergeBalancesBySymbol([]integrations.Balance{})
	assert.Empty(t, merged, "空切片应返回空结果")
}

// TestMergeBalancesBySymbol_单一币种 验证单一币种不合并
func TestMergeBalancesBySymbol_单一币种(t *testing.T) {
	balances := []integrations.Balance{
		{Symbol: "ETH", Name: "ETH", Free: 10, Total: 10, AssetType: "spot"},
	}

	merged := MergeBalancesBySymbol(balances)
	assert.Len(t, merged, 1, "单一币种不需要合并")
	assert.Equal(t, "ETH", merged[0].Symbol)
	assert.Equal(t, 10.0, merged[0].Total)
}

// TestNewClient_存储ApiSecret 验证 NewClient 正确存储 apiSecret
func TestNewClient_存储ApiSecret(t *testing.T) {
	client := NewClient("test-key", "test-secret")
	assert.Equal(t, "test-key", client.apiKey, "apiKey 应正确存储")
	assert.Equal(t, "test-secret", client.apiSecret, "apiSecret 应正确存储")
	assert.NotNil(t, client.client, "Binance SDK 客户端应已创建")
}
