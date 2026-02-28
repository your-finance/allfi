// Package service Gas 优化服务接口
package service

import (
	"context"

	v1 "your-finance/allfi/api/v1/gas"
)

// IGas Gas 优化服务接口
type IGas interface {
	// GetCurrent 获取当前 Gas 价格
	GetCurrent(ctx context.Context) (*v1.GetCurrentRes, error)

	// GetHistory 获取 Gas 价格历史
	GetHistory(ctx context.Context, req *v1.GetHistoryReq) (*v1.GetHistoryRes, error)

	// GetRecommendation 获取最佳交易时间推荐
	GetRecommendation(ctx context.Context, req *v1.GetRecommendationReq) (*v1.GetRecommendationRes, error)

	// GetForecast 获取 Gas 价格预测
	GetForecast(ctx context.Context, req *v1.GetForecastReq) (*v1.GetForecastRes, error)

	// RecordGasPrice 记录 Gas 价格（定时任务调用）
	RecordGasPrice(ctx context.Context) error
}

var localGas IGas

// RegisterGas 注册 Gas 服务
func RegisterGas(i IGas) {
	localGas = i
}

// Gas 获取 Gas 服务实例
func Gas() IGas {
	if localGas == nil {
		panic("implement not found for interface IGas, forgot register?")
	}
	return localGas
}
