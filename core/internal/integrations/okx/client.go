// Package okx OKX 交易所 API 客户端
// 使用 REST API 获取账户余额
package okx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

const (
	baseURL = "https://www.okx.com"
)

// Client OKX 客户端
type Client struct {
	apiKey     string
	apiSecret  string
	passphrase string
	httpClient *http.Client
}

// NewClient 创建 OKX 客户端
func NewClient(apiKey, apiSecret, passphrase string) *Client {
	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		passphrase: passphrase,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetName 获取交易所名称
func (c *Client) GetName() string {
	return "okx"
}

// sign 生成签名
func (c *Client) sign(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// doRequest 发送请求
func (c *Client) doRequest(ctx context.Context, method, path string, body string) ([]byte, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "okx"); err != nil {
		return nil, err
	}

	url := baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// 生成时间戳（ISO 8601 格式）
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	// 设置请求头
	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", c.sign(timestamp, method, path, body))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 OKX 响应失败: %v", err)
	}

	return result, nil
}

// TestConnection 测试 API 连接
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.GetBalances(ctx)
	if err != nil {
		return fmt.Errorf("OKX API 连接失败: %v", err)
	}
	return nil
}

// BalanceResponse OKX 余额响应
type BalanceResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Details []struct {
			Ccy       string `json:"ccy"`       // 币种
			AvailBal  string `json:"availBal"`  // 可用余额
			FrozenBal string `json:"frozenBal"` // 冻结余额
			Bal       string `json:"bal"`       // 总余额
		} `json:"details"`
	} `json:"data"`
}

// GetBalances 获取账户余额
func (c *Client) GetBalances(ctx context.Context) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "okx"); err != nil {
		return nil, err
	}

	url := baseURL + "/api/v5/account/balance"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	path := "/api/v5/account/balance"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", c.sign(timestamp, "GET", path, ""))
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 OKX API 失败: %v", err)
	}
	defer resp.Body.Close()

	var result BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 OKX 响应失败: %v", err)
	}

	if result.Code != "0" {
		return nil, fmt.Errorf("OKX API 错误: %s", result.Msg)
	}

	var balances []integrations.Balance

	for _, data := range result.Data {
		for _, detail := range data.Details {
			availBal, _ := strconv.ParseFloat(detail.AvailBal, 64)
			frozenBal, _ := strconv.ParseFloat(detail.FrozenBal, 64)
			bal, _ := strconv.ParseFloat(detail.Bal, 64)

			if bal <= 0 {
				continue
			}

			balances = append(balances, integrations.Balance{
				Symbol: detail.Ccy,
				Name:   detail.Ccy,
				Free:   availBal,
				Locked: frozenBal,
				Total:  bal,
			})
		}
	}

	return balances, nil
}

// 确保实现接口
var _ integrations.ExchangeClient = (*Client)(nil)
