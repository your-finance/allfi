// Package controller 市场数据模块控制器
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	marketApi "your-finance/allfi/api/v1/market"
	"your-finance/allfi/internal/app/market/service"
)

// MarketController 市场数据控制器
type MarketController struct{}

// GetGas 获取多链 Gas 价格
//
// 对应路由: GET /market/gas
func (c *MarketController) GetGas(ctx context.Context, req *marketApi.GetGasReq) (res *marketApi.GetGasRes, err error) {
	// 调用 Service 层
	result, err := service.Market().GetGasPrice(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "获取 Gas 价格失败")
	}

	// 将业务 DTO 转换为 API 响应
	var prices []marketApi.GasPrice
	for _, chain := range result.Chains {
		prices = append(prices, marketApi.GasPrice{
			Chain:    chain.Chain,
			Low:      chain.Low,
			Standard: chain.Standard,
			Fast:     chain.Fast,
			Instant:  chain.Instant,
			BaseFee:  chain.BaseFee,
			Level:    chain.Level,
		})
	}

	res = &marketApi.GetGasRes{
		Prices:    prices,
		UpdatedAt: result.UpdatedAt,
	}

	return res, nil
}

// Register 注册市场数据路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&MarketController{})
}
