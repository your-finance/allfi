// Package controller 费用分析模块路由注册
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	feeApi "your-finance/allfi/api/v1/fee"
	"your-finance/allfi/internal/app/fee/service"
	"your-finance/allfi/internal/consts"
)

// FeeController 费用分析控制器
type FeeController struct{}

// GetAnalytics 获取费用分析数据
//
// 对应路由: GET /analytics/fees
// 查询参数:
//   - range: 时间范围（7d/30d/90d/1y），默认 30d
//   - currency: 计价货币，默认 USD
func (c *FeeController) GetAnalytics(ctx context.Context, req *feeApi.GetAnalyticsReq) (res *feeApi.GetAnalyticsRes, err error) {
	// 从上下文中获取当前用户 ID
	userID := uint(consts.GetUserID(ctx))

	// 调用 Service 层
	result, err := service.Fee().GetFeeAnalytics(ctx, userID, req.Range, req.Currency)
	if err != nil {
		return nil, gerror.Wrap(err, "获取费用分析失败")
	}

	// 将业务 DTO 转换为 API 响应
	var breakdown []feeApi.FeeBreakdown
	for _, bd := range result.Breakdown {
		breakdown = append(breakdown, feeApi.FeeBreakdown{
			Source:   bd.Source,
			Type:     bd.Type,
			Amount:   bd.Amount,
			Currency: bd.Currency,
			Count:    bd.Count,
		})
	}

	var dailyTrend []feeApi.DailyFee
	for _, df := range result.DailyTrend {
		dailyTrend = append(dailyTrend, feeApi.DailyFee{
			Date:   df.Date,
			Amount: df.Amount,
		})
	}

	res = &feeApi.GetAnalyticsRes{
		TotalFees:   result.TotalFees,
		TradingFees: result.TradingFees,
		GasFees:     result.GasFees,
		Currency:    result.Currency,
		Breakdown:   breakdown,
		DailyTrend:  dailyTrend,
	}

	return res, nil
}

// Register 注册费用分析路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&FeeController{})
}
