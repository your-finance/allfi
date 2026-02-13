// Package provider Local Provider 单元测试
package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalProvider_Name(t *testing.T) {
	p := NewLocalProvider()
	assert.Equal(t, "Local", p.Name())
}

func TestLocalProvider_Priority(t *testing.T) {
	p := NewLocalProvider()
	assert.Equal(t, 999, p.Priority(), "Local Provider应该是最低优先级")
}

func TestLocalProvider_SupportedSymbols(t *testing.T) {
	p := NewLocalProvider()
	symbols := p.SupportedSymbols()

	// 检查支持的币种数量
	assert.Greater(t, len(symbols), 20, "应该支持至少20种币种")

	// 检查是否包含主要币种
	expectedSymbols := []string{"BTC", "ETH", "USDT", "USDC", "BNB", "SOL", "CNY"}
	for _, expected := range expectedSymbols {
		assert.Contains(t, symbols, expected, "应该支持 %s", expected)
	}
}

func TestLocalProvider_SupportsSymbol(t *testing.T) {
	p := NewLocalProvider()

	tests := []struct {
		name     string
		symbol   string
		expected bool
	}{
		{"支持BTC", "BTC", true},
		{"支持ETH", "ETH", true},
		{"支持CNY", "CNY", true},
		{"支持USDT", "USDT", true},
		{"不支持未知币种", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.SupportsSymbol(tt.symbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalProvider_FetchRate(t *testing.T) {
	p := NewLocalProvider()
	ctx := context.Background()

	tests := []struct {
		name   string
		symbol string
	}{
		{"获取BTC汇率", "BTC"},
		{"获取ETH汇率", "ETH"},
		{"获取USDT汇率", "USDT"},
		{"获取CNY汇率", "CNY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, err := p.FetchRate(ctx, tt.symbol)
			require.NoError(t, err)
			require.NotNil(t, rate)

			// 验证返回值
			assert.Equal(t, tt.symbol, rate.Symbol)
			assert.Greater(t, rate.PriceUSD, 0.0, "%s价格应该大于0", tt.symbol)
			assert.Equal(t, "Local", rate.Source)
		})
	}
}

func TestLocalProvider_FetchRate_UnsupportedSymbol(t *testing.T) {
	p := NewLocalProvider()
	ctx := context.Background()

	// 测试不支持的币种
	_, err := p.FetchRate(ctx, "UNKNOWN")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不支持")
}

func TestLocalProvider_FetchRates(t *testing.T) {
	p := NewLocalProvider()
	ctx := context.Background()

	// 测试批量获取汇率
	rates, err := p.FetchRates(ctx)
	require.NoError(t, err)
	require.NotNil(t, rates)

	// 验证返回的汇率数量
	assert.Greater(t, len(rates), 20, "应该返回多个汇率")

	// 验证每个汇率的数据完整性
	for symbol, rate := range rates {
		assert.Equal(t, symbol, rate.Symbol, "symbol应该匹配")
		assert.Greater(t, rate.PriceUSD, 0.0, "%s价格应该大于0", symbol)
		assert.Equal(t, "Local", rate.Source)
	}
}

func TestLocalProvider_IsHealthy(t *testing.T) {
	p := NewLocalProvider()
	ctx := context.Background()

	// Local Provider 总是健康的
	healthy := p.IsHealthy(ctx)
	assert.True(t, healthy, "Local Provider应该总是健康的")
}

func TestLocalProvider_Fallback(t *testing.T) {
	// 测试 Local Provider 作为兜底方案
	p := NewLocalProvider()
	ctx := context.Background()

	// 即使在没有网络的情况下也应该能返回数据
	rate, err := p.FetchRate(ctx, "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	// 验证返回的是参考价格
	assert.Equal(t, "BTC", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0)
}
