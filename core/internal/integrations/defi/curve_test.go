// Package defi Curve Finance 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCurveProtocol_BasicInfo 测试基本信息
func TestCurveProtocol_BasicInfo(t *testing.T) {
	p := NewCurveProtocol(nil, nil)

	assert.Equal(t, "curve", p.GetName())
	assert.Equal(t, "Curve Finance", p.GetDisplayName())
	assert.Equal(t, "lp", p.GetType())
	assert.Contains(t, p.SupportedChains(), "ethereum")
	assert.Contains(t, p.SupportedChains(), "polygon")
}

// TestCurveProtocol_NoClient 测试无客户端时返回空
func TestCurveProtocol_NoClient(t *testing.T) {
	p := NewCurveProtocol(nil, nil)

	positions, err := p.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestResolveCurveAssets_已知池 测试已知 Curve 池解析
func TestResolveCurveAssets_已知池(t *testing.T) {
	assets := resolveCurveAssets("3CRV")
	assert.Equal(t, []string{"DAI", "USDC", "USDT"}, assets)

	assets = resolveCurveAssets("crv3crypto")
	assert.Equal(t, []string{"USDT", "BTC", "ETH"}, assets)

	assets = resolveCurveAssets("stETH-LP")
	assert.Equal(t, []string{"ETH", "stETH"}, assets)
}

// TestResolveCurveAssets_未知池 测试未知池返回默认 3Pool 资产
func TestResolveCurveAssets_未知池(t *testing.T) {
	assets := resolveCurveAssets("UNKNOWN")
	assert.Equal(t, []string{"DAI", "USDC", "USDT"}, assets)
}

// TestIsStablecoinPool_纯稳定币池 测试纯稳定币池判断
func TestIsStablecoinPool_纯稳定币池(t *testing.T) {
	assert.True(t, isStablecoinPool([]string{"DAI", "USDC", "USDT"}))
	assert.True(t, isStablecoinPool([]string{"FRAX", "USDC"}))
	assert.True(t, isStablecoinPool([]string{"crvUSD", "USDC"}))
}

// TestIsStablecoinPool_混合池 测试混合池判断
func TestIsStablecoinPool_混合池(t *testing.T) {
	assert.False(t, isStablecoinPool([]string{"USDT", "BTC", "ETH"}))
	assert.False(t, isStablecoinPool([]string{"ETH", "stETH"}))
	assert.False(t, isStablecoinPool([]string{"ETH", "CRV"}))
}

// TestCurveAPYs_覆盖主要池 测试 APY 映射存在
func TestCurveAPYs_覆盖主要池(t *testing.T) {
	assert.Greater(t, curveAPYs["3crv"], 0.0)
	assert.Greater(t, curveAPYs["crv3crypto"], 0.0)
	assert.Greater(t, curveAPYs["steth-lp"], 0.0)
}
