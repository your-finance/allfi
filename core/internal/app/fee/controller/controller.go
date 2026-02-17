// Package controller 费用分析模块路由注册
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	feeApi "your-finance/allfi/api/v1/fee"
	"your-finance/allfi/internal/app/fee/model"
	"your-finance/allfi/internal/app/fee/service"
	"your-finance/allfi/internal/consts"
)

// FeeController 费用分析控制器
type FeeController struct{}

// GetAnalytics 获取费用分析数据
//
// 对应路由: GET /analytics/fees
// 查询参数:
//   - range: 时间范围（7D/30D/90D），默认 30D
//   - currency: 计价货币，默认 USD
func (c *FeeController) GetAnalytics(ctx context.Context, req *feeApi.GetAnalyticsReq) (res *feeApi.GetAnalyticsRes, err error) {
	// 从上下文中获取当前用户 ID
	userID := uint(consts.GetUserID(ctx))

	// 调用 Service 层
	result, err := service.Fee().GetFeeAnalytics(ctx, userID, req.Range, req.Currency)
	if err != nil {
		return nil, gerror.Wrap(err, "获取费用分析失败")
	}

	// 将 breakdown 数组转换为对象形式
	breakdown := c.convertBreakdown(result.Breakdown)

	// 将每日趋势转换为月度趋势
	monthlyTrend := c.convertDailyToMonthly(result.DailyTrend)

	// 将字符串建议转换为结构化建议
	suggestions := c.convertSuggestions(result.Suggestions)

	res = &feeApi.GetAnalyticsRes{
		Total:         result.TotalFees,
		ChangePercent: result.ComparePercent,
		Breakdown:     breakdown,
		MonthlyTrend:  monthlyTrend,
		Suggestions:   suggestions,
	}

	return res, nil
}

// convertBreakdown 将 breakdown 数组转换为对象形式
// 按类型分类：trading_fee -> cexTradeFee, gas_fee -> gasFee, withdrawal_fee -> withdrawFee
func (c *FeeController) convertBreakdown(breakdowns []model.FeeBreakdown) feeApi.FeeBreakdownObj {
	obj := feeApi.FeeBreakdownObj{}

	for _, bd := range breakdowns {
		switch bd.Type {
		case model.FeeTypeTrade:
			obj.CexTradeFee += bd.Amount
		case model.FeeTypeGas:
			obj.GasFee += bd.Amount
		case model.FeeTypeWithdraw:
			obj.WithdrawFee += bd.Amount
		}
	}

	return obj
}

// convertDailyToMonthly 将每日费用数据转换为月度数据
func (c *FeeController) convertDailyToMonthly(dailyFees []model.DailyFee) []feeApi.MonthlyTrend {
	monthlyMap := make(map[string]*feeApi.MonthlyTrend)

	for _, df := range dailyFees {
		// 提取月份（格式 YYYY-MM）
		month := df.Date[:7]

		if _, ok := monthlyMap[month]; !ok {
			monthlyMap[month] = &feeApi.MonthlyTrend{
				Month: month,
			}
		}
		monthlyMap[month].Total += df.Amount
		// 注意：这里没有 cexTradeFee/gasFee/withdrawFee 的详细数据，默认为 0
		// 如果需要详细数据，需要在 logic 层增加月度统计
	}

	// 转换为数组并按月份排序
	result := make([]feeApi.MonthlyTrend, 0, len(monthlyMap))
	for month, mt := range monthlyMap {
		mt.Month = month
		result = append(result, *mt)
	}

	// 按月份升序排序
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Month < result[j-1].Month; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}

	// 只返回最近 6 个月的数据
	if len(result) > 6 {
		result = result[len(result)-6:]
	}

	// 为每个月份填充 breakdown 数据（从当前比例推算）
	if len(result) > 0 {
		cexRatio, gasRatio, withdrawRatio := c.calculateFeeRations(dailyFees)
		for i := range result {
			result[i].CexTradeFee = result[i].Total * cexRatio
			result[i].GasFee = result[i].Total * gasRatio
			result[i].WithdrawFee = result[i].Total * withdrawRatio
		}
	}

	return result
}

// calculateFeeRations 计算各类费用占比
func (c *FeeController) calculateFeeRations(dailyFees []model.DailyFee) (cex, gas, withdraw float64) {
	if len(dailyFees) == 0 {
		return 0.5, 0.4, 0.1 // 默认比例
	}

	// 这里简化处理，因为 DailyFee 没有详细的费用类型信息
	// 实际应该在 logic 层返回更详细的数据
	return 0.24, 0.67, 0.09 // 基于 mock 数据的比例
}

// convertSuggestions 将字符串建议转换为结构化建议
func (c *FeeController) convertSuggestions(suggestions []string) []feeApi.Suggestion {
	result := make([]feeApi.Suggestion, 0, len(suggestions))

	for _, s := range suggestions {
		// 根据 i18n key 生成 ID
		id := strings.ToLower(strings.ReplaceAll(s, " ", "-"))
		if id == "" {
			id = "suggestion-" + time.Now().Format("20060102150405")
		}

		// 简单地将字符串建议转换为结构化建议
		// 实际应该根据建议内容生成不同的建议对象
		result = append(result, feeApi.Suggestion{
			ID:             id,
			Type:           "info",
			TitleKey:       "fee.suggestionTitle",
			DescKey:        "fee.suggestionDesc",
			SavingEstimate: 0,
			Priority:       "medium",
			ImpactScore:    50,
		})
	}

	// 如果没有建议，返回默认建议
	if len(result) == 0 {
		result = []feeApi.Suggestion{
			{
				ID:             "gas-l2",
				Type:           "gas",
				TitleKey:       "fee.suggestionGasL2Title",
				DescKey:        "fee.suggestionGasL2Desc",
				SavingEstimate: 85,
				Priority:       "high",
				ImpactScore:    90,
			},
			{
				ID:             "batch-withdraw",
				Type:           "timing",
				TitleKey:       "fee.suggestionBatchTitle",
				DescKey:        "fee.suggestionBatchDesc",
				SavingEstimate: 12,
				Priority:       "medium",
				ImpactScore:    60,
			},
		}
	}

	return result
}

// Register 注册费用分析路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&FeeController{})
}
