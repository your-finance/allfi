// Package service 盈亏分析模块 - Service 层接口定义
package service

import (
	"context"

	"your-finance/allfi/internal/app/pnl/model"
)

// IPnl 盈亏分析服务接口
type IPnl interface {
	// GetDailyPnL 获取每日盈亏列表
	// days: 查询天数（7/30/90/365）
	// 返回每日盈亏数据点列表和总盈亏金额
	GetDailyPnL(ctx context.Context, days int) (daily []*model.DailyPnLPoint, totalPnL float64, err error)

	// GetPnLSummary 获取盈亏汇总
	// 返回今日/7日/30日/90日盈亏统计，以及最佳/最差单日
	GetPnLSummary(ctx context.Context) (*model.PnLSummary, error)
}

var localPnl IPnl

// Pnl 获取盈亏分析服务实例
func Pnl() IPnl {
	if localPnl == nil {
		panic("IPnl 服务未注册，请检查 logic/pnl 包的 init 函数")
	}
	return localPnl
}

// RegisterPnl 注册盈亏分析服务实现
// 由 logic 层在 init 函数中调用
func RegisterPnl(i IPnl) {
	localPnl = i
}
