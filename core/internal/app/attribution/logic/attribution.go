// Package logic 资产归因分析模块 - 业务逻辑实现
package logic

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/attribution/model"
	"your-finance/allfi/internal/app/attribution/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/model/entity"
)

// sAttribution 资产归因分析服务实现
type sAttribution struct{}

// New 创建资产归因分析服务实例
func New() service.IAttribution {
	return &sAttribution{}
}

// GetAttribution 获取资产归因分析结果
//
// 业务逻辑:
// 1. 获取时间范围内首尾快照（确定总值变化）
// 2. 获取当前资产明细作为结束状态
// 3. 计算归因三效应:
//   - 价格效应 = startBalance * (endPrice - startPrice)
//   - 数量效应 = startPrice * (endBalance - startBalance)
//   - 交叉效应 = (endBalance - startBalance) * (endPrice - startPrice)
func (s *sAttribution) GetAttribution(ctx context.Context, days int, currency string) (*model.AttributionResult, error) {
	if days <= 0 {
		days = 7
	}
	if currency == "" {
		currency = "USD"
	}

	// 计算时间范围
	startDate := time.Now().AddDate(0, 0, -(days + 1))

	// 获取快照列表
	var snapshots []*entity.AssetSnapshots
	err := assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().SnapshotTime+" >= ?", startDate).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, gerror.Wrap(err, "查询快照数据失败")
	}

	// 快照数据不足时返回空结果
	if len(snapshots) < 2 {
		g.Log().Info(ctx, "快照数据不足，无法进行归因分析", "snapshotCount", len(snapshots))
		return &model.AttributionResult{
			Range:    formatRange(days),
			Currency: currency,
			Assets:   []model.AssetAttribution{},
		}, nil
	}

	// 取起始和结束快照
	startSnap := snapshots[0]
	endSnap := snapshots[len(snapshots)-1]

	// 获取当前资产详情作为结束状态
	userID := consts.GetUserID(ctx)
	var endDetails []*entity.AssetDetails
	err = assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().UserId, userID).
		Scan(&endDetails)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产详情失败")
	}

	// 计算总值变化
	totalStartValue := startSnap.TotalValueUsd
	totalEndValue := endSnap.TotalValueUsd
	totalChange := totalEndValue - totalStartValue

	// 构建各资产归因结果
	var assets []model.AssetAttribution
	var totalPriceEffect, totalQuantityEffect float64

	for _, endDetail := range endDetails {
		endBalance := endDetail.Balance
		endPrice := endDetail.PriceUsd
		endValue := endDetail.ValueUsd

		// 估算起始状态（基于快照比例推算）
		// 如果没有历史明细，假设起始时数量相同，价格按总值比例变化
		var startBalance, startPrice, startValue float64
		if totalEndValue > 0 {
			startBalance = endBalance
			ratio := totalStartValue / totalEndValue
			startPrice = endPrice * ratio
			startValue = startBalance * startPrice
		}

		// 价格效应 = 起始数量 * (结束价格 - 起始价格)
		priceEffect := startBalance * (endPrice - startPrice)
		// 数量效应 = 起始价格 * (结束数量 - 起始数量)
		quantityEffect := startPrice * (endBalance - startBalance)
		// 交叉效应 = (结束数量 - 起始数量) * (结束价格 - 起始价格)
		interactionEffect := (endBalance - startBalance) * (endPrice - startPrice)

		totalPriceEffect += priceEffect
		totalQuantityEffect += quantityEffect

		assets = append(assets, model.AssetAttribution{
			Symbol:            endDetail.AssetSymbol,
			Name:              endDetail.AssetName,
			StartBalance:      startBalance,
			EndBalance:        endBalance,
			StartPrice:        math.Round(startPrice*100) / 100,
			EndPrice:          endPrice,
			StartValue:        math.Round(startValue*100) / 100,
			EndValue:          endValue,
			TotalChange:       math.Round((endValue-startValue)*100) / 100,
			PriceEffect:       math.Round(priceEffect*100) / 100,
			QuantityEffect:    math.Round(quantityEffect*100) / 100,
			InteractionEffect: math.Round(interactionEffect*100) / 100,
		})
	}

	result := &model.AttributionResult{
		Range:          formatRange(days),
		StartTime:      startSnap.SnapshotTime,
		EndTime:        endSnap.SnapshotTime,
		TotalChange:    math.Round(totalChange*100) / 100,
		PriceEffect:    math.Round(totalPriceEffect*100) / 100,
		QuantityEffect: math.Round(totalQuantityEffect*100) / 100,
		Assets:         assets,
		Currency:       currency,
	}

	g.Log().Info(ctx, "获取归因分析成功",
		"days", days,
		"totalChange", totalChange,
		"priceEffect", totalPriceEffect,
		"quantityEffect", totalQuantityEffect,
		"assetCount", len(assets),
	)

	return result, nil
}

// formatRange 将天数转换为范围字符串
func formatRange(days int) string {
	switch days {
	case 1:
		return "1d"
	case 7:
		return "7d"
	case 30:
		return "30d"
	default:
		return "7d"
	}
}
