// Package service 风险管理模块服务接口定义
package service

import (
	"context"

	"your-finance/allfi/internal/app/risk/model"
)

// IRisk 风险管理服务接口
type IRisk interface {
	// GetLatestMetrics 获取最新风险指标
	GetLatestMetrics(ctx context.Context) (*model.RiskMetrics, error)

	// GetHistoryMetrics 获取历史风险指标
	GetHistoryMetrics(ctx context.Context, days int) ([]*model.RiskMetrics, error)

	// CalculateMetrics 计算风险指标
	// period: 计算周期（天数，默认30）
	CalculateMetrics(ctx context.Context, period int) (*model.RiskMetrics, error)
}

var localRisk IRisk

// RegisterRisk 注册风险管理服务
func RegisterRisk(i IRisk) {
	localRisk = i
}

// Risk 获取风险管理服务实例
func Risk() IRisk {
	if localRisk == nil {
		panic("风险管理服务未注册")
	}
	return localRisk
}
