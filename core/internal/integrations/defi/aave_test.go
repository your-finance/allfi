// Package defi Aave 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAaveProtocol_BasicInfo 测试 Aave 协议基本信息
func TestAaveProtocol_BasicInfo(t *testing.T) {
	aave := NewAaveProtocol(nil, nil)

	assert.Equal(t, "aave", aave.GetName())
	assert.Equal(t, "Aave", aave.GetDisplayName())
	assert.Equal(t, "lending", aave.GetType())
	assert.Contains(t, aave.SupportedChains(), "ethereum")
	assert.Contains(t, aave.SupportedChains(), "polygon")
	assert.Contains(t, aave.SupportedChains(), "arbitrum")
	assert.Contains(t, aave.SupportedChains(), "optimism")
}

// TestAaveProtocol_NoClient 测试没有 Etherscan 客户端时返回空
func TestAaveProtocol_NoClient(t *testing.T) {
	aave := NewAaveProtocol(nil, nil)

	positions, err := aave.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestResolveAaveUnderlying 测试 aToken 底层资产解析
func TestResolveAaveUnderlying(t *testing.T) {
	tests := []struct {
		symbol   string
		expected string
	}{
		{"aUSDC", "USDC"},
		{"aDAI", "DAI"},
		{"aWETH", "WETH"},
		{"aWBTC", "WBTC"},
		{"aUSDT", "USDT"},
		{"aMATIC", "MATIC"},       // 通用规则
		{"unknownToken", "unknownToken"}, // 无前缀时原样返回
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := resolveAaveUnderlying(tt.symbol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAaveProtocol_MultiChainSupport 测试 Aave 多链支持
func TestAaveProtocol_MultiChainSupport(t *testing.T) {
	aave := NewAaveProtocol(nil, nil)

	// 不支持的链应返回 nil（没有对应客户端）
	positions, err := aave.GetPositions(context.Background(), "0x1234", "bsc")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}
