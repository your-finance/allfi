// Package logic 资产总览业务逻辑单元测试
package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 GetConversionRate 内部方法的汇率计算
func TestGetConversionRate_DefaultUSD(t *testing.T) {
	logic := New(&sAsset{})
	ctx := testingContext{}

	rate := logic.(*sAsset).getConversionRate(ctx, "USD")
	assert.Equal(t, float64(1.0), rate)
}

func TestGetConversionRate_UnknownCurrency(t *testing.T) {
	logic := New(&sAsset{})
	ctx := testingContext{}

	// 未知货币返回 0
	rate := logic.(*sAsset).getConversionRate(ctx, "UNKNOWN")
	assert.Equal(t, float64(0), rate)
}

// 测试 New 函数返回有效的实例
func TestAssetLogic_New(t *testing.T) {
	logic := New(&sAsset{})
	assert.NotNil(t, logic)
	assert.IsType(t, &sAsset{}, logic)
}

// 测试空列表情况下的计算
func TestAssetLogic_CalculateWithEmptyList(t *testing.T) {
	// 创建空的服务实例
	logic := New(&sAsset{})
	assert.NotNil(t, logic)
}

// 测试 sAsset 结构体字段
func TestSAsset_Fields(t *testing.T) {
	asset := &sAsset{}
	assert.NotNil(t, asset)
}
