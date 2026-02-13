// Package service 策略引擎模块 - 服务接口定义
package service

import (
	"context"

	strategyApi "your-finance/allfi/api/v1/strategy"
)

// IStrategy 策略服务接口
type IStrategy interface {
	// List 获取策略列表
	List(ctx context.Context) (*strategyApi.ListRes, error)

	// Create 创建策略
	Create(ctx context.Context, name, sType string, config any) (*strategyApi.CreateRes, error)

	// Update 更新策略
	Update(ctx context.Context, id uint, name, sType string, config any, isActive *bool) (*strategyApi.UpdateRes, error)

	// Delete 删除策略
	Delete(ctx context.Context, id uint) error

	// GetAnalysis 获取策略分析（偏离度 + 调仓建议）
	GetAnalysis(ctx context.Context, id uint) (*strategyApi.GetAnalysisRes, error)
}

var localStrategy IStrategy

// Strategy 获取策略服务实例
func Strategy() IStrategy {
	if localStrategy == nil {
		panic("IStrategy 服务未注册，请检查 logic/strategy 包的 init 函数")
	}
	return localStrategy
}

// RegisterStrategy 注册策略服务实现
func RegisterStrategy(i IStrategy) {
	localStrategy = i
}
