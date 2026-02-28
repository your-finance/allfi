// Package logic 手动资产业务逻辑单元测试
package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManualAsset_New(t *testing.T) {
	logic := New()
	assert.NotNil(t, logic)
}

func TestManualAsset_Struct(t *testing.T) {
	asset := &sManualAsset{}
	assert.NotNil(t, asset)
}
