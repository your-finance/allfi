// Package service 资产归因分析模块 - Service 层接口定义
package service

import (
	"context"

	"your-finance/allfi/internal/app/attribution/model"
)

// IAttribution 资产归因分析服务接口
type IAttribution interface {
	// GetAttribution 获取资产归因分析结果
	// days: 分析天数（1/7/30）
	// currency: 计价货币
	// 返回归因分析结果，包含总收益、各资产价格/数量/交叉效应
	GetAttribution(ctx context.Context, days int, currency string) (*model.AttributionResult, error)
}

var localAttribution IAttribution

// Attribution 获取资产归因分析服务实例
func Attribution() IAttribution {
	if localAttribution == nil {
		panic("IAttribution 服务未注册，请检查 logic/attribution 包的 init 函数")
	}
	return localAttribution
}

// RegisterAttribution 注册资产归因分析服务实现
// 由 logic 层在 init 函数中调用
func RegisterAttribution(i IAttribution) {
	localAttribution = i
}
