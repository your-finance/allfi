// Package defi Rocket Pool 协议测试
package defi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRocketPoolProtocol_BasicInfo 测试 Rocket Pool 协议基本信息
func TestRocketPoolProtocol_BasicInfo(t *testing.T) {
	rp := NewRocketPoolProtocol(nil, nil)

	assert.Equal(t, "rocketpool", rp.GetName())
	assert.Equal(t, "Rocket Pool", rp.GetDisplayName())
	assert.Equal(t, "staking", rp.GetType())
	assert.Equal(t, []string{"ethereum"}, rp.SupportedChains())
}

// TestRocketPoolProtocol_UnsupportedChain 测试 Rocket Pool 在不支持的链上返回空
func TestRocketPoolProtocol_UnsupportedChain(t *testing.T) {
	rp := NewRocketPoolProtocol(nil, nil)

	positions, err := rp.GetPositions(context.Background(), "0x1234", "polygon")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}

// TestRocketPoolProtocol_NoClient 测试没有 Etherscan 客户端时返回空
func TestRocketPoolProtocol_NoClient(t *testing.T) {
	rp := NewRocketPoolProtocol(nil, nil)

	positions, err := rp.GetPositions(context.Background(), "0x1234", "ethereum")
	assert.NoError(t, err)
	assert.Nil(t, positions)
}
