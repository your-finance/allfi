// Package defi Lido 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"your-finance/allfi/internal/integrations"
)

// mockPriceClient 模拟价格客户端
type mockPriceClient struct {
	prices map[string]float64
}

func (m *mockPriceClient) GetPrice(_ context.Context, symbol string) (float64, error) {
	if p, ok := m.prices[symbol]; ok {
		return p, nil
	}
	return 0, nil
}

func (m *mockPriceClient) GetPrices(_ context.Context, symbols []string) (map[string]float64, error) {
	result := make(map[string]float64)
	for _, s := range symbols {
		if p, ok := m.prices[s]; ok {
			result[s] = p
		}
	}
	return result, nil
}

func (m *mockPriceClient) GetExchangeRate(_ context.Context, _, _ string) (float64, error) {
	return 1.0, nil
}

// 确保 mock 实现接口
var _ integrations.PriceClient = (*mockPriceClient)(nil)

// TestLidoProtocol_BasicInfo 测试 Lido 协议基本信息
func TestLidoProtocol_BasicInfo(t *testing.T) {
	lido := NewLidoProtocol(nil, nil)

	assert.Equal(t, "lido", lido.GetName())
	assert.Equal(t, "Lido Finance", lido.GetDisplayName())
	assert.Equal(t, "staking", lido.GetType())
	assert.Equal(t, []string{"ethereum"}, lido.SupportedChains())
}

// TestLidoProtocol_UnsupportedChain 测试 Lido 在不支持的链上返回空
func TestLidoProtocol_UnsupportedChain(t *testing.T) {
	lido := NewLidoProtocol(nil, nil)

	positions, err := lido.GetPositions(context.Background(), "0x1234", "polygon")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestLidoProtocol_NoClient 测试没有 Etherscan 客户端时返回空
func TestLidoProtocol_NoClient(t *testing.T) {
	lido := NewLidoProtocol(nil, nil)

	positions, err := lido.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}
