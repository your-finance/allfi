// Package logic 资产总览业务逻辑单元测试
package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 New 函数返回有效的实例
func TestAssetLogic_New(t *testing.T) {
	logic := New()
	assert.NotNil(t, logic)
	assert.IsType(t, &sAsset{}, logic)
}

// 测试 sAsset 结构体创建
func TestSAsset_Struct(t *testing.T) {
	asset := &sAsset{}
	assert.NotNil(t, asset)
}
