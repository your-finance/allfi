// Package logic 资产健康评分业务逻辑
// 对用户资产组合进行多维度健康评估：
// - 现金缓冲（25 分）：稳定币占比
// - 集中度（30 分）：最大单一资产占比
// - 平台多样性（20 分）：使用的平台数量
// - 波动性（25 分）：蓝筹资产占比
package logic

import (
	"context"
	"math"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/health_score/model"
	"your-finance/allfi/internal/app/health_score/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/model/entity"

	assetDao "your-finance/allfi/internal/app/asset/dao"
)

// sHealthScore 健康评分服务实现
type sHealthScore struct{}

// New 创建健康评分服务实例
func New() service.IHealthScore {
	return &sHealthScore{}
}

// stableCoins 稳定币列表
var stableCoins = map[string]bool{
	"USDC": true,
	"USDT": true,
	"DAI":  true,
	"BUSD": true,
}

// blueChips 蓝筹币种列表
var blueChips = map[string]bool{
	"BTC": true,
	"ETH": true,
}

// GetHealthScore 计算资产健康评分
func (s *sHealthScore) GetHealthScore(ctx context.Context, currency string) (*model.HealthScoreResult, error) {
	// 直接从 DAO 查询资产详情
	var details []entity.AssetDetails
	err := assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
		Scan(&details)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产详情失败")
	}

	// 统计各类指标
	var totalValue float64
	var stableValue float64
	var blueChipValue float64
	var maxSingleValue float64
	platformSet := make(map[string]bool)

	for _, d := range details {
		totalValue += d.ValueUsd
		symbol := strings.ToUpper(d.AssetSymbol)

		if stableCoins[symbol] {
			stableValue += d.ValueUsd
		}
		if blueChips[symbol] {
			blueChipValue += d.ValueUsd
		}
		if d.ValueUsd > maxSingleValue {
			maxSingleValue = d.ValueUsd
		}
		if d.SourceType != "" {
			platformSet[d.SourceType] = true
		}
	}

	// 避免除零错误
	if totalValue == 0 {
		return &model.HealthScoreResult{
			OverallScore: 0,
			Level:        "poor",
			Details:      []model.HealthScoreDimension{},
			Currency:     currency,
			Weakest:      "cash_buffer",
			Advice:       []string{"暂无资产数据，请添加资产后再进行健康评估"},
		}, nil
	}

	// 计算各维度得分
	dimensions := make([]model.HealthScoreDimension, 4)

	// 1. 现金缓冲（25 分）
	stableRatio := stableValue / totalValue
	cashScore := math.Min(stableRatio/0.25*25, 25)
	cashSuggestion := ""
	if stableRatio < 0.1 {
		cashSuggestion = "建议增加稳定币持仓至总资产的 10%-25%，以应对市场波动"
	} else if stableRatio < 0.25 {
		cashSuggestion = "当前稳定币占比偏低，建议适当增持 USDC/USDT 等稳定币"
	}
	dimensions[0] = model.HealthScoreDimension{
		Category:    "cash_buffer",
		Score:       math.Round(cashScore*100) / 100,
		Weight:      0.25,
		MaxScore:    25,
		Description: "稳定币(USDC/USDT/DAI/BUSD)占比，25% 以上满分",
		Suggestion:  cashSuggestion,
		Value:       math.Round(stableRatio*10000) / 100,
	}

	// 2. 集中度（30 分）
	concentrationRatio := maxSingleValue / totalValue
	var concentrationScore float64
	if concentrationRatio <= 0.3 {
		concentrationScore = 30
	} else if concentrationRatio >= 0.8 {
		concentrationScore = 0
	} else {
		concentrationScore = 30 * (0.8 - concentrationRatio) / 0.5
	}
	concSuggestion := ""
	if concentrationRatio > 0.5 {
		concSuggestion = "资产集中度过高，建议分散投资以降低单一资产风险"
	} else if concentrationRatio > 0.3 {
		concSuggestion = "最大单一资产占比略高，可考虑进一步分散配置"
	}
	dimensions[1] = model.HealthScoreDimension{
		Category:    "concentration",
		Score:       math.Round(concentrationScore*100) / 100,
		Weight:      0.30,
		MaxScore:    30,
		Description: "最大单一资产占比，<30% 满分，>80% 为 0 分",
		Suggestion:  concSuggestion,
		Value:       math.Round(concentrationRatio*10000) / 100,
	}

	// 3. 平台多样性（20 分）
	platformCount := len(platformSet)
	var platformScore float64
	switch {
	case platformCount >= 4:
		platformScore = 20
	case platformCount == 3:
		platformScore = 15
	case platformCount == 2:
		platformScore = 10
	case platformCount == 1:
		platformScore = 5
	default:
		platformScore = 0
	}
	platSuggestion := ""
	if platformCount < 2 {
		platSuggestion = "仅使用单一平台，建议将资产分散到多个平台以降低平台风险"
	} else if platformCount < 4 {
		platSuggestion = "可考虑增加更多平台分散存储，降低单平台风险"
	}
	dimensions[2] = model.HealthScoreDimension{
		Category:    "platform_diversity",
		Score:       platformScore,
		Weight:      0.20,
		MaxScore:    20,
		Description: "使用的平台数量，>=4 个满分，1 个为 5 分",
		Suggestion:  platSuggestion,
		Value:       float64(platformCount),
	}

	// 4. 波动性（25 分）
	blueChipRatio := blueChipValue / totalValue
	volatilityScore := math.Min(blueChipRatio/0.6*25, 25)
	volSuggestion := ""
	if blueChipRatio < 0.3 {
		volSuggestion = "蓝筹资产(BTC/ETH)占比较低，建议增持以提升组合稳定性"
	} else if blueChipRatio < 0.6 {
		volSuggestion = "可适当增加 BTC/ETH 等蓝筹资产的配置比例"
	}
	dimensions[3] = model.HealthScoreDimension{
		Category:    "volatility",
		Score:       math.Round(volatilityScore*100) / 100,
		Weight:      0.25,
		MaxScore:    25,
		Description: "蓝筹(BTC/ETH)占比，>60% 满分",
		Suggestion:  volSuggestion,
		Value:       math.Round(blueChipRatio*10000) / 100,
	}

	// 计算总分
	overallScore := cashScore + concentrationScore + platformScore + volatilityScore

	// 确定等级
	var level string
	switch {
	case overallScore >= 80:
		level = "excellent"
	case overallScore >= 60:
		level = "good"
	case overallScore >= 40:
		level = "fair"
	default:
		level = "poor"
	}

	// 找最弱维度
	weakest := dimensions[0]
	for _, d := range dimensions[1:] {
		if d.MaxScore > 0 && d.Score/d.MaxScore < weakest.Score/weakest.MaxScore {
			weakest = d
		}
	}

	// 汇总建议
	var advice []string
	for _, d := range dimensions {
		if d.Suggestion != "" {
			advice = append(advice, d.Suggestion)
		}
	}
	if len(advice) == 0 {
		advice = append(advice, "资产配置健康，继续保持当前策略")
	}

	g.Log().Debug(ctx, "健康评分计算完成",
		"overallScore", overallScore,
		"level", level,
		"weakest", weakest.Category,
	)

	return &model.HealthScoreResult{
		OverallScore: math.Round(overallScore*100) / 100,
		Level:        level,
		Details:      dimensions,
		Currency:     currency,
		Weakest:      weakest.Category,
		Advice:       advice,
	}, nil
}
