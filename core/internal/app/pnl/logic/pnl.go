// Package logic 盈亏分析模块 - 业务逻辑实现
package logic

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/pnl/model"
	"your-finance/allfi/internal/app/pnl/service"
	"your-finance/allfi/internal/model/entity"

	assetDao "your-finance/allfi/internal/app/asset/dao"
)

// sPnl 盈亏分析服务实现
type sPnl struct{}

// New 创建盈亏分析服务实例
func New() service.IPnl {
	return &sPnl{}
}

// GetDailyPnL 获取每日盈亏列表
//
// 业务逻辑:
// 1. 从 asset_snapshots 表获取指定天数内的快照数据
// 2. 按日期分组，取每天最早和最晚的快照
// 3. 计算每日盈亏: pnl = endVal - startVal, percent = pnl/startVal*100
// 4. 使用前一天结束值作为当天的开始值
func (s *sPnl) GetDailyPnL(ctx context.Context, days int) (daily []*model.DailyPnLPoint, totalPnL float64, err error) {
	if days <= 0 {
		days = 30
	}

	// 计算起始日期
	startDate := time.Now().AddDate(0, 0, -days)

	// 获取快照数据
	var snapshots []*entity.AssetSnapshots
	err = assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().SnapshotTime+" >= ?", startDate).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询快照数据失败")
	}

	if len(snapshots) < 2 {
		g.Log().Info(ctx, "快照数据不足，无法计算每日盈亏", "snapshotCount", len(snapshots))
		return []*model.DailyPnLPoint{}, 0, nil
	}

	// 按日期分组快照，取每天最早和最晚的值
	type dayEntry struct {
		first float64
		last  float64
	}
	dayMap := make(map[string]*dayEntry)
	dayOrder := make([]string, 0)

	for _, snap := range snapshots {
		dateKey := snap.SnapshotTime.Format("2006-01-02")
		val := snap.TotalValueUsd

		entry, exists := dayMap[dateKey]
		if !exists {
			dayOrder = append(dayOrder, dateKey)
			dayMap[dateKey] = &dayEntry{first: val, last: val}
		} else {
			entry.last = val
		}
	}

	// 计算每日盈亏（使用前一天的最后值作为当天起始值）
	daily = make([]*model.DailyPnLPoint, 0, len(dayOrder))
	var prevEndValue float64

	for i, dateKey := range dayOrder {
		entry := dayMap[dateKey]
		startVal := entry.first
		if i > 0 {
			startVal = prevEndValue
		}
		endVal := entry.last
		pnl := endVal - startVal

		var pnlPercent float64
		if startVal > 0 {
			pnlPercent = math.Round(pnl/startVal*10000) / 100 // 保留两位小数
		}
		totalPnL += pnl

		daily = append(daily, &model.DailyPnLPoint{
			Date:       dateKey,
			StartValue: startVal,
			EndValue:   endVal,
			PnL:        math.Round(pnl*100) / 100,
			PnLPercent: pnlPercent,
		})
		prevEndValue = endVal
	}

	g.Log().Info(ctx, "获取每日盈亏成功",
		"days", days,
		"dataPoints", len(daily),
		"totalPnL", totalPnL,
	)

	return daily, math.Round(totalPnL*100) / 100, nil
}

// GetPnLSummary 获取盈亏汇总
//
// 业务逻辑:
// 1. 获取最新快照作为当前值
// 2. 分别获取 7d/30d/90d 时间窗口的起始快照
// 3. 计算各时段盈亏金额和百分比
// 4. 遍历 90 天数据找出最佳/最差单日
func (s *sPnl) GetPnLSummary(ctx context.Context) (*model.PnLSummary, error) {
	// 获取 90 天快照数据（用于计算所有指标）
	startDate := time.Now().AddDate(0, 0, -90)

	var snapshots []*entity.AssetSnapshots
	err := assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().SnapshotTime+" >= ?", startDate).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, gerror.Wrap(err, "查询快照数据失败")
	}

	if len(snapshots) == 0 {
		return &model.PnLSummary{}, nil
	}

	// 获取当前值（最新快照）
	currentValue := snapshots[len(snapshots)-1].TotalValueUsd

	// 计算各时段盈亏
	now := time.Now()
	day7Start := now.AddDate(0, 0, -7)
	day30Start := now.AddDate(0, 0, -30)
	day90Start := now.AddDate(0, 0, -90)

	pnl7d := s.calculatePeriodPnL(snapshots, day7Start, currentValue)
	pnl30d := s.calculatePeriodPnL(snapshots, day30Start, currentValue)
	pnl90d := s.calculatePeriodPnL(snapshots, day90Start, currentValue)

	// 计算每日盈亏，找出最佳/最差单日
	daily, _, _ := s.GetDailyPnL(ctx, 90)

	summary := &model.PnLSummary{
		TotalPnL:        pnl90d.PnL,
		TotalPnLPercent: pnl90d.PnLPercent,
		PnL7d:           pnl7d.PnL,
		PnL30d:          pnl30d.PnL,
		PnL90d:          pnl90d.PnL,
	}

	// 查找最佳和最差单日
	if len(daily) > 0 {
		// 按盈亏排序找出最佳和最差
		sorted := make([]*model.DailyPnLPoint, len(daily))
		copy(sorted, daily)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].PnL > sorted[j].PnL
		})

		summary.BestDay = sorted[0].Date
		summary.BestDayPnL = sorted[0].PnL
		summary.WorstDay = sorted[len(sorted)-1].Date
		summary.WorstDayPnL = sorted[len(sorted)-1].PnL
	}

	g.Log().Info(ctx, "获取盈亏汇总成功",
		"totalPnL", summary.TotalPnL,
		"pnl7d", summary.PnL7d,
		"pnl30d", summary.PnL30d,
	)

	return summary, nil
}

// calculatePeriodPnL 计算指定时段的盈亏
// 找到最接近 periodStart 的快照作为起始值
func (s *sPnl) calculatePeriodPnL(snapshots []*entity.AssetSnapshots, periodStart time.Time, currentValue float64) model.PnLPeriod {
	if len(snapshots) == 0 {
		return model.PnLPeriod{EndValue: currentValue}
	}

	// 找到最接近 periodStart 的快照
	var startValue float64
	for _, snap := range snapshots {
		if !snap.SnapshotTime.Before(periodStart) {
			startValue = snap.TotalValueUsd
			break
		}
		startValue = snap.TotalValueUsd
	}

	if startValue == 0 {
		startValue = snapshots[0].TotalValueUsd
	}

	pnl := currentValue - startValue
	var pnlPercent float64
	if startValue > 0 {
		pnlPercent = math.Round(pnl/startValue*10000) / 100
	}

	return model.PnLPeriod{
		PnL:        math.Round(pnl*100) / 100,
		PnLPercent: pnlPercent,
		StartValue: startValue,
		EndValue:   currentValue,
	}
}
