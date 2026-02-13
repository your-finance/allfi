// Package provider Gate.io Provider 单元测试
package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGateioProvider_Name(t *testing.T) {
	p := NewGateioProvider()
	assert.Equal(t, "Gate.io", p.Name())
}

func TestGateioProvider_Priority(t *testing.T) {
	p := NewGateioProvider()
	assert.Equal(t, 2, p.Priority())
}

func TestGateioProvider_SupportedSymbols(t *testing.T) {
	p := NewGateioProvider()
	symbols := p.SupportedSymbols()

	// 检查支持的币种数量
	assert.Greater(t, len(symbols), 10, "应该支持至少10种加密货币")

	// 检查是否包含主要币种
	expectedSymbols := []string{"BTC", "ETH", "USDT"}
	for _, expected := range expectedSymbols {
		assert.Contains(t, symbols, expected, "应该支持 %s", expected)
	}
}

func TestGateioProvider_SupportsSymbol(t *testing.T) {
	p := NewGateioProvider()

	tests := []struct {
		name     string
		symbol   string
		expected bool
	}{
		{"支持BTC", "BTC", true},
		{"支持ETH", "ETH", true},
		{"不支持CNY", "CNY", false},
		{"不支持未知币种", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.SupportsSymbol(tt.symbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGateioProvider_FetchRate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewGateioProvider()
	ctx := context.Background()

	// 测试获取 BTC 汇率
	rate, err := p.FetchRate(ctx, "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	// 验证返回值
	assert.Equal(t, "BTC", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0, "BTC价格应该大于0")
	assert.Equal(t, "Gate.io", rate.Source)
}

func TestGateioProvider_FetchRates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewGateioProvider()
	ctx := context.Background()

	// 测试批量获取汇率
	rates, err := p.FetchRates(ctx)
	require.NoError(t, err)
	require.NotNil(t, rates)

	// 验证返回的汇率数量
	assert.Greater(t, len(rates), 5, "应该返回多个汇率")

	// 验证每个汇率的数据完整性
	for symbol, rate := range rates {
		assert.Equal(t, symbol, rate.Symbol, "symbol应该匹配")
		assert.Greater(t, rate.PriceUSD, 0.0, "%s价格应该大于0", symbol)
		assert.Equal(t, "Gate.io", rate.Source)
	}
}

func TestGateioProvider_IsHealthy(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewGateioProvider()
	ctx := context.Background()

	// 测试健康检查
	healthy := p.IsHealthy(ctx)
	assert.True(t, healthy, "Gate.io API应该是健康的")
}
