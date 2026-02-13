// Package defi Uniswap V3 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUniswapV3Protocol_BasicInfo 测试基本信息
func TestUniswapV3Protocol_BasicInfo(t *testing.T) {
	p := NewUniswapV3Protocol(nil, nil)

	assert.Equal(t, "uniswap_v3", p.GetName())
	assert.Equal(t, "Uniswap V3", p.GetDisplayName())
	assert.Equal(t, "lp", p.GetType())
	assert.Contains(t, p.SupportedChains(), "ethereum")
	assert.Contains(t, p.SupportedChains(), "base")
	assert.Contains(t, p.SupportedChains(), "optimism")
}

// TestUniswapV3Protocol_NoClient 测试无客户端时返回空
func TestUniswapV3Protocol_NoClient(t *testing.T) {
	p := NewUniswapV3Protocol(nil, nil)

	positions, err := p.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestResolveV3PoolPair_已知池 测试已知池子解析
func TestResolveV3PoolPair_已知池(t *testing.T) {
	tests := []struct {
		symbol string
		token0 string
		token1 string
	}{
		{"UNI-V3-WETH-USDC", "ETH", "USDC"},
		{"UNI-V3-WBTC-WETH", "BTC", "ETH"},
		{"UNI-V3-WBTC-USDC", "BTC", "USDC"},
	}

	for _, tt := range tests {
		t0, t1 := resolveV3PoolPair(tt.symbol)
		assert.Equal(t, tt.token0, t0, "symbol: %s token0", tt.symbol)
		assert.Equal(t, tt.token1, t1, "symbol: %s token1", tt.symbol)
	}
}

// TestResolveV3PoolPair_未知池 测试未知符号返回默认值
func TestResolveV3PoolPair_未知池(t *testing.T) {
	t0, t1 := resolveV3PoolPair("UNKNOWN-V3")
	assert.Equal(t, "ETH", t0)
	assert.Equal(t, "USDC", t1)
}

// TestUniswapV3_PositionManagerAddress 测试仓位管理合约地址常量
func TestUniswapV3_PositionManagerAddress(t *testing.T) {
	assert.Equal(t, "0xc36442b4a4522e871399cd717abdd847ab11fe88", uniV3PositionManager)
}
