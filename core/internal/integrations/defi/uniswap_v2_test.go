// Package defi Uniswap V2 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUniswapV2Protocol_BasicInfo 测试基本信息
func TestUniswapV2Protocol_BasicInfo(t *testing.T) {
	p := NewUniswapV2Protocol(nil, nil)

	assert.Equal(t, "uniswap_v2", p.GetName())
	assert.Equal(t, "Uniswap V2", p.GetDisplayName())
	assert.Equal(t, "lp", p.GetType())
	assert.Contains(t, p.SupportedChains(), "ethereum")
	assert.Contains(t, p.SupportedChains(), "polygon")
	assert.Contains(t, p.SupportedChains(), "arbitrum")
}

// TestUniswapV2Protocol_NoClient 测试无客户端时返回空
func TestUniswapV2Protocol_NoClient(t *testing.T) {
	p := NewUniswapV2Protocol(nil, nil)

	positions, err := p.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestResolveV2LPPair_已知对 测试已知 LP 对解析
func TestResolveV2LPPair_已知对(t *testing.T) {
	tests := []struct {
		symbol string
		token0 string
		token1 string
	}{
		{"UNI-V2-WETH-USDC", "ETH", "USDC"},
		{"UNI-V2-WBTC-WETH", "BTC", "ETH"},
		{"UNI-V2-USDC-USDT", "USDC", "USDT"},
		{"SLP", "ETH", "USDT"},
	}

	for _, tt := range tests {
		t0, t1 := resolveV2LPPair(tt.symbol)
		assert.Equal(t, tt.token0, t0, "symbol: %s token0", tt.symbol)
		assert.Equal(t, tt.token1, t1, "symbol: %s token1", tt.symbol)
	}
}

// TestResolveV2LPPair_未知对 测试未知符号返回默认值
func TestResolveV2LPPair_未知对(t *testing.T) {
	t0, t1 := resolveV2LPPair("UNKNOWN-LP")
	assert.Equal(t, "ETH", t0)
	assert.Equal(t, "USDC", t1)
}

// TestGetTokenPrice_WrappedMapping 测试 wrapped 代币价格映射
func TestGetTokenPrice_WrappedMapping(t *testing.T) {
	mock := &mockPriceClient{
		prices: map[string]float64{
			"ETH":   3000,
			"BTC":   60000,
			"MATIC": 0.8,
		},
	}

	ctx := context.Background()

	assert.Equal(t, 3000.0, getTokenPrice(ctx, mock, "WETH"))
	assert.Equal(t, 60000.0, getTokenPrice(ctx, mock, "WBTC"))
	assert.Equal(t, 0.8, getTokenPrice(ctx, mock, "WMATIC"))
	assert.Equal(t, 3000.0, getTokenPrice(ctx, mock, "ETH"))
}
