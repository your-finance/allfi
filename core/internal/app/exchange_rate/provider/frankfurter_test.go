// Package provider Frankfurter Provider 单元测试
package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrankfurterProvider_Name(t *testing.T) {
	p := NewFrankfurterProvider()
	assert.Equal(t, "Frankfurter", p.Name())
}

func TestFrankfurterProvider_Priority(t *testing.T) {
	p := NewFrankfurterProvider()
	assert.Equal(t, 3, p.Priority())
}

func TestFrankfurterProvider_SupportedSymbols(t *testing.T) {
	p := NewFrankfurterProvider()
	symbols := p.SupportedSymbols()

	// Frankfurter 只支持 CNY
	assert.Equal(t, 1, len(symbols), "Frankfurter只应该支持CNY")
	assert.Contains(t, symbols, "CNY", "应该支持CNY")
}

func TestFrankfurterProvider_SupportsSymbol(t *testing.T) {
	p := NewFrankfurterProvider()

	tests := []struct {
		name     string
		symbol   string
		expected bool
	}{
		{"支持CNY", "CNY", true},
		{"不支持USD", "USD", false},
		{"不支持BTC", "BTC", false},
		{"不支持ETH", "ETH", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.SupportsSymbol(tt.symbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFrankfurterProvider_FetchRate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewFrankfurterProvider()
	ctx := context.Background()

	// 测试获取 CNY 汇率
	rate, err := p.FetchRate(ctx, "CNY")
	require.NoError(t, err)
	require.NotNil(t, rate)

	// 验证返回值
	assert.Equal(t, "CNY", rate.Symbol)
	assert.Greater(t, rate.PriceUSD, 0.0, "CNY汇率应该大于0")
	assert.Equal(t, "Frankfurter", rate.Source)
}

func TestFrankfurterProvider_FetchRate_UnsupportedSymbol(t *testing.T) {
	p := NewFrankfurterProvider()
	ctx := context.Background()

	// 测试不支持的币种
	_, err := p.FetchRate(ctx, "BTC")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不支持")
}

func TestFrankfurterProvider_FetchRates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewFrankfurterProvider()
	ctx := context.Background()

	// 测试批量获取汇率
	rates, err := p.FetchRates(ctx)
	require.NoError(t, err)
	require.NotNil(t, rates)

	// Frankfurter 只返回 CNY 汇率
	assert.Equal(t, 1, len(rates), "应该只返回CNY汇率")

	cnyRate, ok := rates["CNY"]
	require.True(t, ok, "应该包含CNY汇率")
	assert.Greater(t, cnyRate.PriceUSD, 0.0, "CNY汇率应该大于0")
}

func TestFrankfurterProvider_IsHealthy(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	p := NewFrankfurterProvider()
	ctx := context.Background()

	// 测试健康检查
	healthy := p.IsHealthy(ctx)
	assert.True(t, healthy, "Frankfurter API应该是健康的")
}
