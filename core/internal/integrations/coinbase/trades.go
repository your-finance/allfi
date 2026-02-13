// Package coinbase Coinbase 交易历史
// 获取交易记录（通过 Coinbase transactions API）
package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

// transactionsResponse Coinbase 交易列表响应
type transactionsResponse struct {
	Data []struct {
		ID     string `json:"id"`
		Type   string `json:"type"`   // send, receive, buy, sell, trade, transfer
		Status string `json:"status"` // pending, completed, failed, cancelled
		Amount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"amount"`
		NativeAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"native_amount"`
		Network struct {
			Hash string `json:"hash"`
		} `json:"network"`
		To struct {
			Address string `json:"address"`
		} `json:"to"`
		From struct {
			Address string `json:"address"`
		} `json:"from"`
		CreatedAt string `json:"created_at"` // ISO 8601
	} `json:"data"`
	Pagination struct {
		NextURI string `json:"next_uri"`
	} `json:"pagination"`
}

// GetTradeHistory 获取交易历史
// Coinbase 使用统一的 transactions API，从中筛选 buy/sell/trade 类型
func (c *Client) GetTradeHistory(ctx context.Context, params integrations.TradeHistoryParams) ([]integrations.Trade, error) {
	if err := utils.WaitForAPI(ctx, "coinbase"); err != nil {
		return nil, err
	}

	// 获取所有账户
	accounts, err := c.getAccountIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Coinbase 账户列表失败: %v", err)
	}

	var allTrades []integrations.Trade

	for _, accountID := range accounts {
		trades, err := c.getAccountTransactions(ctx, accountID, params)
		if err != nil {
			continue // 单个账户失败不影响其他
		}
		allTrades = append(allTrades, trades...)
	}

	return allTrades, nil
}

// GetDepositHistory 获取充值历史
// 从 transactions API 筛选 receive 类型
func (c *Client) GetDepositHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	if err := utils.WaitForAPI(ctx, "coinbase"); err != nil {
		return nil, err
	}

	accounts, err := c.getAccountIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Coinbase 账户列表失败: %v", err)
	}

	var allTransfers []integrations.Transfer

	for _, accountID := range accounts {
		transfers, err := c.getAccountTransfers(ctx, accountID, "deposit")
		if err != nil {
			continue
		}
		allTransfers = append(allTransfers, transfers...)
	}

	return allTransfers, nil
}

// GetWithdrawHistory 获取提现历史
// 从 transactions API 筛选 send 类型
func (c *Client) GetWithdrawHistory(ctx context.Context, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	if err := utils.WaitForAPI(ctx, "coinbase"); err != nil {
		return nil, err
	}

	accounts, err := c.getAccountIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Coinbase 账户列表失败: %v", err)
	}

	var allTransfers []integrations.Transfer

	for _, accountID := range accounts {
		transfers, err := c.getAccountTransfers(ctx, accountID, "withdraw")
		if err != nil {
			continue
		}
		allTransfers = append(allTransfers, transfers...)
	}

	return allTransfers, nil
}

// getAccountIDs 获取所有账户 ID
func (c *Client) getAccountIDs(ctx context.Context) ([]string, error) {
	path := "/v2/accounts"
	url := baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-VERSION", "2021-08-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var ids []string
	for _, acc := range result.Data {
		ids = append(ids, acc.ID)
	}

	return ids, nil
}

// getAccountTransactions 获取账户交易记录并转换为 Trade 类型
func (c *Client) getAccountTransactions(ctx context.Context, accountID string, params integrations.TradeHistoryParams) ([]integrations.Trade, error) {
	path := fmt.Sprintf("/v2/accounts/%s/transactions", accountID)
	url := baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-VERSION", "2021-08-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result transactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var trades []integrations.Trade
	for _, tx := range result.Data {
		// 筛选 buy/sell/trade 类型
		if tx.Type != "buy" && tx.Type != "sell" && tx.Type != "trade" {
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, tx.CreatedAt)

		// 时间范围过滤
		if !params.StartTime.IsZero() && createdAt.Before(params.StartTime) {
			continue
		}
		if !params.EndTime.IsZero() && createdAt.After(params.EndTime) {
			continue
		}

		amount, _ := strconv.ParseFloat(tx.Amount.Amount, 64)
		if amount < 0 {
			amount = -amount
		}

		side := "buy"
		if tx.Type == "sell" {
			side = "sell"
		}

		trades = append(trades, integrations.Trade{
			ID:        tx.ID,
			Symbol:    tx.Amount.Currency,
			Side:      side,
			Price:     0, // Coinbase transactions API 不直接提供单价
			Quantity:  amount,
			Fee:       0,
			FeeCoin:   tx.NativeAmount.Currency,
			Timestamp: createdAt,
			Source:    "coinbase",
		})
	}

	return trades, nil
}

// getAccountTransfers 获取账户充提记录
func (c *Client) getAccountTransfers(ctx context.Context, accountID string, transferType string) ([]integrations.Transfer, error) {
	path := fmt.Sprintf("/v2/accounts/%s/transactions", accountID)
	url := baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-VERSION", "2021-08-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result transactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var transfers []integrations.Transfer
	for _, tx := range result.Data {
		// 充值筛选 receive 类型，提现筛选 send 类型
		targetType := "receive"
		if transferType == "withdraw" {
			targetType = "send"
		}
		if tx.Type != targetType {
			continue
		}

		createdAt, _ := time.Parse(time.RFC3339, tx.CreatedAt)
		amount, _ := strconv.ParseFloat(tx.Amount.Amount, 64)
		if amount < 0 {
			amount = -amount
		}

		address := tx.To.Address
		if transferType == "deposit" {
			address = tx.From.Address
		}

		transfers = append(transfers, integrations.Transfer{
			ID:        tx.ID,
			Type:      transferType,
			Coin:      tx.Amount.Currency,
			Amount:    amount,
			Fee:       0,
			Address:   address,
			TxHash:    tx.Network.Hash,
			Status:    tx.Status,
			Timestamp: createdAt,
			Source:    "coinbase",
		})
	}

	return transfers, nil
}
