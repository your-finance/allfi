// Package controller 汇率模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	exchangeRateApi "your-finance/allfi/api/v1/exchange_rate"
	"your-finance/allfi/internal/app/exchange_rate/service"
)

// Controller 汇率控制器
type Controller struct{}

// GetCurrent 获取实时汇率
func (c *Controller) GetCurrent(ctx context.Context, req *exchangeRateApi.GetCurrentReq) (res *exchangeRateApi.GetCurrentRes, err error) {
	return service.ExchangeRate().GetRates(ctx, req.Currencies)
}

// GetPrices 获取币种价格
func (c *Controller) GetPrices(ctx context.Context, req *exchangeRateApi.GetPricesReq) (res *exchangeRateApi.GetPricesRes, err error) {
	return service.ExchangeRate().GetPrices(ctx, req.Symbols)
}

// Refresh 强制刷新汇率缓存
func (c *Controller) Refresh(ctx context.Context, req *exchangeRateApi.RefreshReq) (res *exchangeRateApi.RefreshRes, err error) {
	return service.ExchangeRate().RefreshRates(ctx)
}

// GetHistory 获取历史汇率
func (c *Controller) GetHistory(ctx context.Context, req *exchangeRateApi.GetHistoryReq) (res *exchangeRateApi.GetHistoryRes, err error) {
	return service.ExchangeRate().GetHistory(ctx, req.Base, req.Quote, req.Days)
}

// Register 注册汇率模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}
