// Package defi 协议注册中心测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockProtocol 模拟 DeFi 协议（用于注册中心测试）
type mockProtocol struct {
	name        string
	displayName string
	protocolType string
	chains      []string
	positions   []Position
}

func (m *mockProtocol) GetName() string        { return m.name }
func (m *mockProtocol) GetDisplayName() string  { return m.displayName }
func (m *mockProtocol) GetType() string         { return m.protocolType }
func (m *mockProtocol) SupportedChains() []string { return m.chains }
func (m *mockProtocol) GetPositions(_ context.Context, _ string, _ string) ([]Position, error) {
	return m.positions, nil
}

// TestRegistry_RegisterAndList 测试注册协议和列出协议
func TestRegistry_RegisterAndList(t *testing.T) {
	registry := NewRegistry()

	p1 := &mockProtocol{
		name: "lido", displayName: "Lido Finance",
		protocolType: "staking", chains: []string{"ethereum"},
	}
	p2 := &mockProtocol{
		name: "aave", displayName: "Aave",
		protocolType: "lending", chains: []string{"ethereum", "polygon"},
	}

	registry.Register(p1)
	registry.Register(p2)

	infos := registry.ListProtocols()
	assert.Len(t, infos, 2)

	// 验证能获取到指定协议
	protocol, err := registry.GetProtocol("lido")
	assert.NoError(t, err)
	assert.Equal(t, "lido", protocol.GetName())

	// 获取不存在的协议应返回错误
	_, err = registry.GetProtocol("unknown")
	assert.Error(t, err)
}

// TestRegistry_GetAllPositions 测试聚合所有协议仓位
func TestRegistry_GetAllPositions(t *testing.T) {
	registry := NewRegistry()

	p1 := &mockProtocol{
		name: "proto_a", displayName: "Protocol A",
		protocolType: "staking", chains: []string{"ethereum"},
		positions: []Position{
			{Protocol: "proto_a", Type: "staking", ValueUSD: 1000},
		},
	}
	p2 := &mockProtocol{
		name: "proto_b", displayName: "Protocol B",
		protocolType: "lending", chains: []string{"ethereum", "polygon"},
		positions: []Position{
			{Protocol: "proto_b", Type: "lending", ValueUSD: 2000},
			{Protocol: "proto_b", Type: "lending", ValueUSD: 500},
		},
	}

	registry.Register(p1)
	registry.Register(p2)

	// 查询 Ethereum 链上所有仓位
	positions, err := registry.GetAllPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Len(t, positions, 3) // 1 + 2

	// 查询 Polygon 链上仓位（只有 proto_b 支持 Polygon）
	positions, err = registry.GetAllPositions(context.Background(), "0x1234", "polygon")
	assert.NoError(t, err)
	assert.Len(t, positions, 2) // 只有 proto_b 的仓位
}

// TestRegistry_GetPositionsByProtocol 测试按协议查询仓位
func TestRegistry_GetPositionsByProtocol(t *testing.T) {
	registry := NewRegistry()

	p := &mockProtocol{
		name: "lido", displayName: "Lido",
		protocolType: "staking", chains: []string{"ethereum"},
		positions: []Position{
			{Protocol: "lido", Type: "staking", ValueUSD: 5000},
		},
	}

	registry.Register(p)

	positions, err := registry.GetPositionsByProtocol(context.Background(), "0x1234", "ethereum", "lido")
	assert.NoError(t, err)
	assert.Len(t, positions, 1)
	assert.Equal(t, float64(5000), positions[0].ValueUSD)

	// 查询未注册的协议
	_, err = registry.GetPositionsByProtocol(context.Background(), "0x1234", "ethereum", "unknown")
	assert.Error(t, err)
}

// TestSupportsChain 测试链支持检查
func TestSupportsChain(t *testing.T) {
	p := &mockProtocol{chains: []string{"ethereum", "polygon"}}

	assert.True(t, supportsChain(p, "ethereum"))
	assert.True(t, supportsChain(p, "polygon"))
	assert.False(t, supportsChain(p, "bsc"))
}
