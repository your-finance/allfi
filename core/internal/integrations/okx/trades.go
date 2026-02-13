// Package okx OKX 交易历史
// 获取成交记录、充值和提现历史
package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
)

// fillsHistoryResponse OKX 成交历史响应
type fillsHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstID  string `json:"instId"`  // 产品ID（如 BTC-USDT）
		TradeID string `json:"tradeId"` // 成交ID
		Side    string `json:"side"`    // buy/sell
		Px      string `json:"px"`      // 成交价格
		Sz      string `json:"sz"`      // 成交数量
		Fee     string `json:"fee"`     // 手续费（负数）
		FeeCcy  string `json:"feeCcy"`  // 手续费币种
		Ts      string `json:"ts"`      // 成交时间戳（毫秒）
	} `json:"data"`
}

// depositHistoryResponse OKX 充值历史响应
type depositHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		DepId string `json:"depId"` // 充值ID
		Ccy   string `json:"ccy"`   // 币种
		Amt   string `json:"amt"`   // 数量
		From  string `json:"from"`  // 来源地址
		TxId  string `json:"txId"`  // 交易哈希
		State string `json:"state"` // 状态：0-等待, 1-到账, 2-已完成
		Ts    string `json:"ts"`    // 时间戳
	} `json:"data"`
}

// withdrawHistoryResponse OKX 提现历史响应
type withdrawHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		WdId  string `json:"wdId"`  // 提现ID
		Ccy   string `json:"ccy"`   // 币种
		Amt   string `json:"amt"`   // 数量
		Fee   string `json:"fee"`   // 手续费
		To    string `json:"to"`    // 目标地址
		TxId  string `json:"txId"`  // 交易哈希
		State string `json:"state"` // 状态
		Ts    string `json:"ts"`    // 时间戳
	} `json:"data"`
}

// GetTradeHistory 获取成交历史
// 调用 OKX GET /api/v5/trade/fills-history
func (c *Client) GetTradeHistory(ctx context.Context, params integrations.TradeHistoryParams) ([]integrations.Trade, error) {
	path := "/api/v5/trade/fills-history"

	// 构建查询参数
	if params.Symbol != "" {
		path += "?instId=" + params.Symbol
	}
	if !params.StartTime.IsZero() {
		sep := "?"
		if params.Symbol != "" {
			sep = "&"
		}
		path += sep + "begin=" + strconv.FormatInt(params.StartTime.UnixMilli(), 10)
	}
	if !params.EndTime.IsZero() {
		path += "&end=" + strconv.FormatInt(params.EndTime.UnixMilli(), 10)
	}
	if params.Limit > 0 {
		path += "&limit=" + strconv.Itoa(params.Limit)
	}

	body, err := c.doRequest(ctx, "GET", path, "")
	if err != nil {
		return nil, fmt.Errorf("获取 OKX 成交历史失败: %v", err)
	}

	return parseTradeHistory(body, ctx)
}

// parseTradeHistory 解析成交历史响应
func parseTradeHistory(body []byte, _ context.Context) ([]integrations.Trade, error) {
	var result fillsHistoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 OKX 成交历史响应失败: %v", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API 错误: %s", result.Msg)
	}

	trades := make([]integrations.Trade, 0, len(result.Data))
	for _, t := range result.Data {
		price, _ := strconv.ParseFloat(t.Px, 64)
		qty, _ := strconv.ParseFloat(t.Sz, 64)
		fee, _ := strconv.ParseFloat(t.Fee, 64)
		ts, _ := strconv.ParseInt(t.Ts, 10, 64)

		// OKX 手续费为负数，取绝对值
		if fee < 0 {
			fee = -fee
		}

		trades = append(trades, integrations.Trade{
			ID:        t.TradeID,
			Symbol:    t.InstID,
			Side:      t.Side,
			Price:     price,
			Quantity:  qty,
			Fee:       fee,
			FeeCoin:   t.FeeCcy,
			Timestamp: time.UnixMilli(ts),
			Source:    "okx",
		})
	}

	return trades, nil
}

// GetDepositHistory 获取充值历史
// 调用 OKX GET /api/v5/asset/deposit-history
func (c *Client) GetDepositHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	path := "/api/v5/asset/deposit-history"

	if params.Coin != "" {
		path += "?ccy=" + params.Coin
	}

	body, err := c.doRequest(ctx, "GET", path, "")
	if err != nil {
		return nil, fmt.Errorf("获取 OKX 充值历史失败: %v", err)
	}

	return parseDepositHistory(body)
}

// parseDepositHistory 解析充值历史响应
func parseDepositHistory(body []byte) ([]integrations.Transfer, error) {
	var result depositHistoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 OKX 充值历史响应失败: %v", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API 错误: %s", result.Msg)
	}

	transfers := make([]integrations.Transfer, 0, len(result.Data))
	for _, d := range result.Data {
		amt, _ := strconv.ParseFloat(d.Amt, 64)
		ts, _ := strconv.ParseInt(d.Ts, 10, 64)

		// 状态映射
		status := "pending"
		switch d.State {
		case "1", "2":
			status = "completed"
		}

		transfers = append(transfers, integrations.Transfer{
			ID:        d.DepId,
			Type:      "deposit",
			Coin:      d.Ccy,
			Amount:    amt,
			Fee:       0,
			Address:   d.From,
			TxHash:    d.TxId,
			Status:    status,
			Timestamp: time.UnixMilli(ts),
			Source:    "okx",
		})
	}

	return transfers, nil
}

// GetWithdrawHistory 获取提现历史
// 调用 OKX GET /api/v5/asset/withdrawal-history
func (c *Client) GetWithdrawHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	path := "/api/v5/asset/withdrawal-history"

	if params.Coin != "" {
		path += "?ccy=" + params.Coin
	}

	body, err := c.doRequest(ctx, "GET", path, "")
	if err != nil {
		return nil, fmt.Errorf("获取 OKX 提现历史失败: %v", err)
	}

	return parseWithdrawHistory(body)
}

// parseWithdrawHistory 解析提现历史响应
func parseWithdrawHistory(body []byte) ([]integrations.Transfer, error) {
	var result withdrawHistoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 OKX 提现历史响应失败: %v", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API 错误: %s", result.Msg)
	}

	transfers := make([]integrations.Transfer, 0, len(result.Data))
	for _, w := range result.Data {
		amt, _ := strconv.ParseFloat(w.Amt, 64)
		fee, _ := strconv.ParseFloat(w.Fee, 64)
		ts, _ := strconv.ParseInt(w.Ts, 10, 64)

		// 状态映射：-3:撤销中, -2:已取消, -1:提现失败,
		// 0-等待提现, 1-提现中, 2-已完成, 7-审核通过, 10-等待划转
		status := "pending"
		switch w.State {
		case "2":
			status = "completed"
		case "-2":
			status = "cancelled"
		case "-1", "-3":
			status = "failed"
		}

		transfers = append(transfers, integrations.Transfer{
			ID:        w.WdId,
			Type:      "withdraw",
			Coin:      w.Ccy,
			Amount:    amt,
			Fee:       fee,
			Address:   w.To,
			TxHash:    w.TxId,
			Status:    status,
			Timestamp: time.UnixMilli(ts),
			Source:    "okx",
		})
	}

	return transfers, nil
}
