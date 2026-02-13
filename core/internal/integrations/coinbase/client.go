// Package coinbase Coinbase 交易所 API 客户端
// 使用 REST API 获取账户余额
package coinbase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/utils"
)

const (
	baseURL = "https://api.coinbase.com"
)

// Client Coinbase 客户端
type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

// NewClient 创建 Coinbase 客户端
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetName 获取交易所名称
func (c *Client) GetName() string {
	return "coinbase"
}

// sign 生成签名
func (c *Client) sign(timestamp, method, path, body string) string {
	message := timestamp + method + path + body
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// TestConnection 测试 API 连接
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.GetBalances(ctx)
	if err != nil {
		return fmt.Errorf("Coinbase API 连接失败: %v", err)
	}
	return nil
}

// AccountsResponse Coinbase 账户列表响应
type AccountsResponse struct {
	Data []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Currency struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"currency"`
		Balance struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"balance"`
	} `json:"data"`
	Pagination struct {
		NextURI string `json:"next_uri"`
	} `json:"pagination"`
}

// GetBalances 获取账户余额
func (c *Client) GetBalances(ctx context.Context) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "coinbase"); err != nil {
		return nil, err
	}

	path := "/v2/accounts"
	url := baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-VERSION", "2021-08-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Coinbase API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Coinbase API 错误: HTTP %d", resp.StatusCode)
	}

	var result AccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 Coinbase 响应失败: %v", err)
	}

	// 第一步：收集所有非零余额
	type rawBalance struct {
		symbol string
		name   string
		amount float64
	}
	var rawBalances []rawBalance
	var symbols []string

	for _, account := range result.Data {
		amount, _ := strconv.ParseFloat(account.Balance.Amount, 64)

		if amount <= 0 {
			continue
		}

		rawBalances = append(rawBalances, rawBalance{
			symbol: account.Currency.Code,
			name:   account.Currency.Name,
			amount: amount,
		})
		symbols = append(symbols, account.Currency.Code)
	}

	// 第二步：通过 CoinGecko 批量获取 USD 价格
	var prices map[string]float64
	if len(symbols) > 0 {
		cgClient := coingecko.NewClient("")
		var priceErr error
		prices, priceErr = cgClient.GetPrices(ctx, symbols)
		if priceErr != nil {
			// 获取价格失败不影响余额返回，仅 ValueUSD 为 0
			prices = make(map[string]float64)
		}
	} else {
		prices = make(map[string]float64)
	}

	// 第三步：组装余额列表，计算 USD 价值
	var balances []integrations.Balance
	for _, rb := range rawBalances {
		priceUSD := prices[rb.symbol]
		valueUSD := rb.amount * priceUSD

		balances = append(balances, integrations.Balance{
			Symbol:   rb.symbol,
			Name:     rb.name,
			Free:     rb.amount,
			Locked:   0,
			Total:    rb.amount,
			ValueUSD: valueUSD,
		})
	}

	return balances, nil
}

// 确保实现接口
var _ integrations.ExchangeClient = (*Client)(nil)
