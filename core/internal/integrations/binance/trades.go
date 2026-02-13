// Package binance Binance 交易历史
// 获取现货交易记录、充值和提现历史
package binance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

// GetTradeHistory 获取现货交易历史
// 调用 Binance GET /api/v3/myTrades
func (c *Client) GetTradeHistory(ctx context.Context, params integrations.TradeHistoryParams) ([]integrations.Trade, error) {
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	svc := c.client.NewListTradesService()

	// 交易对必填，默认 BTCUSDT
	symbol := params.Symbol
	if symbol == "" {
		symbol = "BTCUSDT"
	}
	svc.Symbol(symbol)

	if !params.StartTime.IsZero() {
		svc.StartTime(params.StartTime.UnixMilli())
	}
	if !params.EndTime.IsZero() {
		svc.EndTime(params.EndTime.UnixMilli())
	}

	limit := params.Limit
	if limit <= 0 || limit > 1000 {
		limit = 500
	}
	svc.Limit(limit)

	result, err := svc.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Binance 交易历史失败: %v", err)
	}

	trades := make([]integrations.Trade, 0, len(result))
	for _, t := range result {
		price, _ := strconv.ParseFloat(t.Price, 64)
		qty, _ := strconv.ParseFloat(t.Quantity, 64)
		fee, _ := strconv.ParseFloat(t.Commission, 64)

		side := "buy"
		if !t.IsBuyer {
			side = "sell"
		}

		trades = append(trades, integrations.Trade{
			ID:        strconv.FormatInt(t.ID, 10),
			Symbol:    t.Symbol,
			Side:      side,
			Price:     price,
			Quantity:  qty,
			Fee:       fee,
			FeeCoin:   t.CommissionAsset,
			Timestamp: time.UnixMilli(t.Time),
			Source:    "binance",
		})
	}

	return trades, nil
}

// GetDepositHistory 获取充值历史
// 调用 Binance GET /sapi/v1/capital/deposit/hisrec
func (c *Client) GetDepositHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	svc := c.client.NewListDepositsService()

	if params.Coin != "" {
		svc.Coin(params.Coin)
	}
	if !params.StartTime.IsZero() {
		svc.StartTime(params.StartTime.UnixMilli())
	}
	if !params.EndTime.IsZero() {
		svc.EndTime(params.EndTime.UnixMilli())
	}

	result, err := svc.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Binance 充值历史失败: %v", err)
	}

	transfers := make([]integrations.Transfer, 0, len(result))
	for _, d := range result {
		// 状态映射：0-待确认, 1-成功, 6-已入账
		status := "pending"
		if d.Status == 1 || d.Status == 6 {
			status = "completed"
		}

		transfers = append(transfers, integrations.Transfer{
			ID:        d.TxID,
			Type:      "deposit",
			Coin:      d.Coin,
			Amount:    parseFloat(d.Amount),
			Fee:       0, // 充值一般无手续费
			Address:   d.Address,
			TxHash:    d.TxID,
			Status:    status,
			Timestamp: time.UnixMilli(d.InsertTime),
			Source:    "binance",
		})
	}

	return transfers, nil
}

// GetWithdrawHistory 获取提现历史
// 调用 Binance GET /sapi/v1/capital/withdraw/history
func (c *Client) GetWithdrawHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	svc := c.client.NewListWithdrawsService()

	if params.Coin != "" {
		svc.Coin(params.Coin)
	}
	if !params.StartTime.IsZero() {
		svc.StartTime(params.StartTime.UnixMilli())
	}
	if !params.EndTime.IsZero() {
		svc.EndTime(params.EndTime.UnixMilli())
	}

	result, err := svc.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Binance 提现历史失败: %v", err)
	}

	transfers := make([]integrations.Transfer, 0, len(result))
	for _, w := range result {
		// 状态映射：0-确认邮件已发, 1-已取消, 2-等待确认,
		// 3-已拒绝, 4-处理中, 5-提现失败, 6-已完成
		status := "pending"
		switch w.Status {
		case 6:
			status = "completed"
		case 1:
			status = "cancelled"
		case 3, 5:
			status = "failed"
		}

		transfers = append(transfers, integrations.Transfer{
			ID:        w.ID,
			Type:      "withdraw",
			Coin:      w.Coin,
			Amount:    parseFloat(w.Amount),
			Fee:       parseFloat(w.TransactionFee),
			Address:   w.Address,
			TxHash:    w.TxID,
			Status:    status,
			Timestamp: time.Now(), // Binance 提现历史无直接时间字段
			Source:    "binance",
		})
	}

	return transfers, nil
}

// parseFloat 解析浮点数，忽略错误
func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
