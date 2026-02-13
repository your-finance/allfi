// Package logic 费用分析业务逻辑
// 从交易记录中汇总分析费用，提供按来源/类型分类、趋势分析和优化建议
package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/fee/model"
	"your-finance/allfi/internal/app/fee/service"
	transactionDao "your-finance/allfi/internal/app/transaction/dao"
	"your-finance/allfi/internal/model/entity"
)

// sFee 费用分析服务实现
type sFee struct{}

// New 创建费用分析服务实例
func New() service.IFee {
	return &sFee{}
}

// 交易类型常量
const (
	txTypeBuy      = "buy"
	txTypeSell     = "sell"
	txTypeSwap     = "swap"
	txTypeTransfer = "transfer"
	txTypeDeposit  = "deposit"
	txTypeWithdraw = "withdraw"
)

// GetFeeAnalytics 获取费用分析
func (s *sFee) GetFeeAnalytics(ctx context.Context, userID uint, period string, currency string) (*model.FeeAnalytics, error) {
	// 解析时间范围
	endTime := time.Now()
	startTime, prevStart := parsePeriod(period, endTime)

	// 查询当期交易记录
	var txs []entity.UnifiedTransactions
	err := transactionDao.UnifiedTransactions.Ctx(ctx).
		Where(transactionDao.UnifiedTransactions.Columns().UserId, userID).
		WhereGTE(transactionDao.UnifiedTransactions.Columns().Timestamp, startTime).
		WhereLTE(transactionDao.UnifiedTransactions.Columns().Timestamp, endTime).
		Scan(&txs)
	if err != nil {
		return nil, gerror.Wrap(err, "查询交易记录失败")
	}

	// 汇总费用
	analytics := s.aggregateFees(txs, currency)

	// 计算每日趋势
	analytics.DailyTrend = s.calculateDailyTrend(txs)

	// 计算与上期对比
	if !prevStart.IsZero() {
		var prevTxs []entity.UnifiedTransactions
		err := transactionDao.UnifiedTransactions.Ctx(ctx).
			Where(transactionDao.UnifiedTransactions.Columns().UserId, userID).
			WhereGTE(transactionDao.UnifiedTransactions.Columns().Timestamp, prevStart).
			WhereLT(transactionDao.UnifiedTransactions.Columns().Timestamp, startTime).
			Scan(&prevTxs)
		if err == nil {
			prevAnalytics := s.aggregateFees(prevTxs, currency)
			if prevAnalytics.TotalFees > 0 {
				analytics.ComparePercent = ((analytics.TotalFees - prevAnalytics.TotalFees) / prevAnalytics.TotalFees) * 100
			}
		}
	}

	// 生成优化建议
	analytics.Suggestions = s.generateSuggestions(analytics)

	g.Log().Debug(ctx, "费用分析完成",
		"userID", userID,
		"period", period,
		"totalFees", analytics.TotalFees,
		"breakdownCount", len(analytics.Breakdown),
	)

	return analytics, nil
}

// aggregateFees 汇总交易费用
func (s *sFee) aggregateFees(txs []entity.UnifiedTransactions, currency string) *model.FeeAnalytics {
	analytics := &model.FeeAnalytics{
		Currency: currency,
	}

	type breakdownKey struct {
		source  string
		feeType string
	}
	breakdownMap := make(map[breakdownKey]*model.FeeBreakdown)

	for _, tx := range txs {
		if tx.Fee <= 0 {
			continue
		}

		feeUSD := float64(tx.Fee)
		analytics.TotalFees += feeUSD

		// 判断费用类型
		feeType := classifyFeeType(tx.TxType, tx.Chain)

		switch feeType {
		case model.FeeTypeTrade:
			analytics.TradingFees += feeUSD
		case model.FeeTypeGas:
			analytics.GasFees += feeUSD
		}

		// 按来源+类型分组
		key := breakdownKey{source: tx.Source, feeType: feeType}
		bd, ok := breakdownMap[key]
		if !ok {
			bd = &model.FeeBreakdown{
				Source:   tx.Source,
				Type:     feeType,
				Currency: currency,
			}
			breakdownMap[key] = bd
		}
		bd.Amount += feeUSD
		bd.Count++
	}

	for _, bd := range breakdownMap {
		analytics.Breakdown = append(analytics.Breakdown, *bd)
	}

	return analytics
}

// classifyFeeType 根据交易类型判断费用类别
func classifyFeeType(txType, chain string) string {
	switch txType {
	case txTypeBuy, txTypeSell, txTypeSwap:
		return model.FeeTypeTrade
	case txTypeWithdraw:
		return model.FeeTypeWithdraw
	case txTypeDeposit, txTypeTransfer:
		if chain != "" {
			return model.FeeTypeGas
		}
		return model.FeeTypeTrade
	default:
		return model.FeeTypeTrade
	}
}

// calculateDailyTrend 计算每日费用趋势
func (s *sFee) calculateDailyTrend(txs []entity.UnifiedTransactions) []model.DailyFee {
	dailyMap := make(map[string]float64)

	for _, tx := range txs {
		if tx.Fee <= 0 {
			continue
		}
		date := tx.Timestamp.Format("2006-01-02")
		dailyMap[date] += float64(tx.Fee)
	}

	result := make([]model.DailyFee, 0, len(dailyMap))
	for date, amount := range dailyMap {
		result = append(result, model.DailyFee{
			Date:   date,
			Amount: amount,
		})
	}

	// 按日期升序排序
	sortDailyFees(result)
	return result
}

// sortDailyFees 按日期升序排序
func sortDailyFees(fees []model.DailyFee) {
	for i := 1; i < len(fees); i++ {
		for j := i; j > 0 && fees[j].Date < fees[j-1].Date; j-- {
			fees[j], fees[j-1] = fees[j-1], fees[j]
		}
	}
}

// generateSuggestions 生成智能费用优化建议
func (s *sFee) generateSuggestions(analytics *model.FeeAnalytics) []string {
	var suggestions []string

	if analytics.TotalFees == 0 {
		suggestions = append(suggestions, "暂无费用记录")
		return suggestions
	}

	if analytics.GasFees > 0 && analytics.GasFees/analytics.TotalFees > 0.4 {
		suggestions = append(suggestions,
			"Gas 费占比超过 40%，建议使用 L2 链（如 Arbitrum、Optimism、Base）降低链上交易成本")
	}

	if analytics.TradingFees > 100 {
		suggestions = append(suggestions,
			"交易手续费较高，建议使用交易所 VIP 等级或平台币抵扣手续费")
	}

	if analytics.ComparePercent > 50 {
		suggestions = append(suggestions,
			fmt.Sprintf("费用较上期上升 %.0f%%，请关注是否有异常交易", analytics.ComparePercent))
	}

	sourceCount := 0
	sourceSet := make(map[string]bool)
	for _, bd := range analytics.Breakdown {
		if !sourceSet[bd.Source] {
			sourceSet[bd.Source] = true
			sourceCount++
		}
	}
	if sourceCount > 3 {
		suggestions = append(suggestions,
			"资金分散在多个平台，可考虑集中到手续费更优惠的平台")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "费用支出合理，暂无优化建议")
	}

	return suggestions
}

// parsePeriod 解析时间范围，返回当期开始和上期开始时间
func parsePeriod(period string, endTime time.Time) (startTime, prevStart time.Time) {
	switch period {
	case "7d":
		startTime = endTime.AddDate(0, 0, -7)
		prevStart = startTime.AddDate(0, 0, -7)
	case "30d":
		startTime = endTime.AddDate(0, 0, -30)
		prevStart = startTime.AddDate(0, 0, -30)
	case "90d":
		startTime = endTime.AddDate(0, 0, -90)
		prevStart = startTime.AddDate(0, 0, -90)
	case "1y":
		startTime = endTime.AddDate(-1, 0, 0)
		prevStart = startTime.AddDate(-1, 0, 0)
	default:
		startTime = endTime.AddDate(0, 0, -30)
		prevStart = startTime.AddDate(0, 0, -30)
	}
	return
}
