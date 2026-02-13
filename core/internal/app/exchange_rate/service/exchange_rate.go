// Package service 汇率模块 - 服务接口定义
// 提供汇率查询、价格查询、汇率刷新功能
package service

import (
	"context"

	exchangeRateApi "your-finance/allfi/api/v1/exchange_rate"
)

// IExchangeRate 汇率服务接口
type IExchangeRate interface {
	// GetRates 获取汇率（基于 exchange_rates 表的 from_currency/to_currency）
	GetRates(ctx context.Context, currencies string) (*exchangeRateApi.GetCurrentRes, error)

	// GetPrices 获取加密货币价格
	GetPrices(ctx context.Context, symbols string) (*exchangeRateApi.GetPricesRes, error)

	// RefreshRates 刷新汇率缓存
	RefreshRates(ctx context.Context) (*exchangeRateApi.RefreshRes, error)

	// GetHistory 获取历史汇率
	GetHistory(ctx context.Context, base, quote string, days int) (*exchangeRateApi.GetHistoryRes, error)
}

var localExchangeRate IExchangeRate

// ExchangeRate 获取汇率服务实例
func ExchangeRate() IExchangeRate {
	if localExchangeRate == nil {
		panic("IExchangeRate 服务未注册，请检查 logic/exchange_rate 包的 init 函数")
	}
	return localExchangeRate
}

// RegisterExchangeRate 注册汇率服务实现
// 由 logic 层在 init 函数中调用
func RegisterExchangeRate(i IExchangeRate) {
	localExchangeRate = i
}
