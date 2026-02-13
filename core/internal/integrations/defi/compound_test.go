// Package defi Compound 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompoundProtocol_BasicInfo 测试 Compound 协议基本信息
func TestCompoundProtocol_BasicInfo(t *testing.T) {
	compound := NewCompoundProtocol(nil, nil)

	assert.Equal(t, "compound", compound.GetName())
	assert.Equal(t, "Compound", compound.GetDisplayName())
	assert.Equal(t, "lending", compound.GetType())
	assert.Equal(t, []string{"ethereum"}, compound.SupportedChains())
}

// TestCompoundProtocol_UnsupportedChain 测试 Compound 在不支持的链上返回空
func TestCompoundProtocol_UnsupportedChain(t *testing.T) {
	compound := NewCompoundProtocol(nil, nil)

	positions, err := compound.GetPositions(context.Background(), "0x1234", "polygon")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestCompoundProtocol_NoClient 测试没有 Etherscan 客户端时返回空
func TestCompoundProtocol_NoClient(t *testing.T) {
	compound := NewCompoundProtocol(nil, nil)

	positions, err := compound.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestResolveCompoundUnderlying 测试 cToken 底层资产解析
func TestResolveCompoundUnderlying(t *testing.T) {
	tests := []struct {
		symbol       string
		underlying   string
		hasRate      bool // 是否有已知兑换率
	}{
		{"cDAI", "DAI", true},
		{"cUSDC", "USDC", true},
		{"cUSDT", "USDT", true},
		{"cETH", "ETH", true},
		{"cWBTC", "WBTC", true},
		{"cLINK", "LINK", false}, // 未知 cToken，使用通用规则
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			info := resolveCompoundUnderlying(tt.symbol)
			assert.Equal(t, tt.underlying, info.Underlying)
			if tt.hasRate {
				assert.Greater(t, info.ExchangeRate, 0.0)
			}
		})
	}
}

// TestCompoundExchangeRate 测试 Compound 兑换率计算逻辑
func TestCompoundExchangeRate(t *testing.T) {
	// cDAI 兑换率约 0.0225，即 100 cDAI ≈ 2.25 DAI
	info := resolveCompoundUnderlying("cDAI")
	cTokenBalance := 1000.0
	underlyingAmount := cTokenBalance * info.ExchangeRate

	assert.Equal(t, "DAI", info.Underlying)
	assert.InDelta(t, 22.5, underlyingAmount, 0.1)
}
