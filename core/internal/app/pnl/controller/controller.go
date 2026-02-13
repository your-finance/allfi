// Package controller 盈亏分析模块 - 控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	pnlApi "your-finance/allfi/api/v1/pnl"
	"your-finance/allfi/internal/app/pnl/service"
)

// PnlController 盈亏分析控制器
type PnlController struct{}

// GetDaily 获取每日盈亏
func (c *PnlController) GetDaily(ctx context.Context, req *pnlApi.GetDailyReq) (res *pnlApi.GetDailyRes, err error) {
	// 调用服务层获取每日盈亏
	daily, _, err := service.Pnl().GetDailyPnL(ctx, req.Days)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &pnlApi.GetDailyRes{
		Daily: make([]pnlApi.DailyPnLItem, 0, len(daily)),
	}
	for _, d := range daily {
		res.Daily = append(res.Daily, pnlApi.DailyPnLItem{
			Date:       d.Date,
			PnL:        d.PnL,
			PnLPercent: d.PnLPercent,
			TotalValue: d.EndValue,
		})
	}

	return res, nil
}

// GetSummary 获取盈亏汇总
func (c *PnlController) GetSummary(ctx context.Context, req *pnlApi.GetSummaryReq) (res *pnlApi.GetSummaryRes, err error) {
	// 调用服务层获取盈亏汇总
	summary, err := service.Pnl().GetPnLSummary(ctx)
	if err != nil {
		return nil, err
	}

	res = &pnlApi.GetSummaryRes{
		TotalPnL:        summary.TotalPnL,
		TotalPnLPercent: summary.TotalPnLPercent,
		PnL7d:           summary.PnL7d,
		PnL30d:          summary.PnL30d,
		PnL90d:          summary.PnL90d,
		BestDay:         summary.BestDay,
		WorstDay:        summary.WorstDay,
		BestDayPnL:      summary.BestDayPnL,
		WorstDayPnL:     summary.WorstDayPnL,
	}

	return res, nil
}

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&PnlController{})
}
