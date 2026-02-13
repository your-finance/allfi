// Package logic 趋势预测模块 - 业务逻辑实现
package logic

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/forecast/model"
	"your-finance/allfi/internal/app/forecast/service"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/model/entity"
)

// sForecast 趋势预测服务实现
type sForecast struct{}

// New 创建趋势预测服务实例
func New() service.IForecast {
	return &sForecast{}
}

// GetForecast 获取趋势预测
//
// 业务逻辑:
// 1. 获取 90 天的快照数据做回归分析
// 2. 按天去重（取每天最后的快照值）
// 3. 线性回归: y = slope * x + intercept
// 4. 计算 R²（决定系数/拟合优度）
// 5. 预测达成目标的天数，限制最多 10 年
// 6. 趋势判断: slope > 1 → "up", < -1 → "down", 其他 → "flat"
func (s *sForecast) GetForecast(ctx context.Context, targetValue float64, currency string) (*model.ForecastResult, error) {
	if targetValue <= 0 {
		return nil, gerror.New("目标值必须为正数")
	}
	if currency == "" {
		currency = "USD"
	}

	// 获取 90 天的快照数据
	startDate := time.Now().AddDate(0, 0, -90)

	var snapshots []*entity.AssetSnapshots
	err := assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().SnapshotTime+" >= ?", startDate).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, gerror.Wrap(err, "查询快照数据失败")
	}

	// 数据不足时返回默认结果
	if len(snapshots) < 3 {
		g.Log().Info(ctx, "数据不足，需要至少 3 个快照才能预测", "snapshotCount", len(snapshots))
		return &model.ForecastResult{
			TargetValue: targetValue,
			Currency:    currency,
			Trend:       "flat",
			DataPoints:  len(snapshots),
		}, nil
	}

	// 按天去重（取每天最后的快照值）
	type dayPoint struct {
		dayIndex int     // 从第一天开始的天数
		value    float64 // 当天总值
	}

	firstTime := snapshots[0].SnapshotTime
	dayMap := make(map[int]float64)
	for _, snap := range snapshots {
		dayIdx := int(snap.SnapshotTime.Sub(firstTime).Hours() / 24)
		dayMap[dayIdx] = snap.TotalValueUsd
	}

	points := make([]dayPoint, 0, len(dayMap))
	for idx, val := range dayMap {
		points = append(points, dayPoint{dayIndex: idx, value: val})
	}

	if len(points) < 3 {
		return &model.ForecastResult{
			TargetValue: targetValue,
			Currency:    currency,
			Trend:       "flat",
			DataPoints:  len(points),
		}, nil
	}

	// 线性回归: y = slope * x + intercept
	// x = 天数索引, y = 总资产值
	n := float64(len(points))
	var sumX, sumY, sumXY, sumX2 float64
	for _, p := range points {
		x := float64(p.dayIndex)
		y := p.value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return &model.ForecastResult{
			TargetValue: targetValue,
			Currency:    currency,
			Trend:       "flat",
			DataPoints:  len(points),
		}, nil
	}

	slope := (n*sumXY - sumX*sumY) / denominator
	intercept := (sumY - slope*sumX) / n

	// 计算 R²（决定系数/拟合优度）
	meanY := sumY / n
	var ssTotal, ssResidual float64
	for _, p := range points {
		x := float64(p.dayIndex)
		predicted := slope*x + intercept
		ssTotal += (p.value - meanY) * (p.value - meanY)
		ssResidual += (p.value - predicted) * (p.value - predicted)
	}
	var rSquared float64
	if ssTotal > 0 {
		rSquared = 1.0 - ssResidual/ssTotal
	}

	// 当前值（最新快照）
	latestSnap := snapshots[len(snapshots)-1]
	currentValue := latestSnap.TotalValueUsd
	currentDayIdx := int(latestSnap.SnapshotTime.Sub(firstTime).Hours() / 24)

	// 趋势方向
	trend := "flat"
	if slope > 1 {
		trend = "up"
	} else if slope < -1 {
		trend = "down"
	}

	// 日均增长率
	var growthRate float64
	if currentValue > 0 {
		growthRate = slope / currentValue * 100
	}

	result := &model.ForecastResult{
		CurrentValue: currentValue,
		TargetValue:  targetValue,
		Currency:     currency,
		DailyGrowth:  math.Round(slope*100) / 100,
		GrowthRate:   math.Round(growthRate*1000) / 1000,
		Confidence:   math.Round(rSquared*1000) / 1000,
		Trend:        trend,
		DataPoints:   len(points),
	}

	// 预测达成目标的天数
	if slope > 0 && targetValue > currentValue {
		// target = slope * futureDay + intercept
		futureDay := (targetValue - intercept) / slope
		daysFromNow := int(math.Ceil(futureDay - float64(currentDayIdx)))
		if daysFromNow > 0 && daysFromNow < 3650 { // 最多预测 10 年
			estimatedDate := time.Now().AddDate(0, 0, daysFromNow)
			result.EstimatedDate = &estimatedDate
			result.DaysToTarget = daysFromNow
		}
	} else if currentValue >= targetValue {
		// 已经达到目标
		now := time.Now()
		result.EstimatedDate = &now
		result.DaysToTarget = 0
	}

	g.Log().Info(ctx, "获取趋势预测成功",
		"currentValue", currentValue,
		"targetValue", targetValue,
		"slope", slope,
		"rSquared", rSquared,
		"trend", trend,
		"daysToTarget", result.DaysToTarget,
	)

	return result, nil
}
