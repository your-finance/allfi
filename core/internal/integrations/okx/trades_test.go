// Package okx OKX 交易历史测试
package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"your-finance/allfi/internal/integrations"
)

// TestGetTradeHistory_解析成交记录 验证成交历史 API 响应解析
func TestGetTradeHistory_解析成交记录(t *testing.T) {
	// 模拟 OKX API 响应
	mockResponse := `{
		"code": "0",
		"msg": "",
		"data": [
			{
				"instId": "BTC-USDT",
				"tradeId": "123456",
				"side": "buy",
				"px": "50000.5",
				"sz": "0.1",
				"fee": "-0.05",
				"feeCcy": "USDT",
				"ts": "1707580800000"
			},
			{
				"instId": "ETH-USDT",
				"tradeId": "123457",
				"side": "sell",
				"px": "3000.0",
				"sz": "1.5",
				"fee": "-0.03",
				"feeCcy": "USDT",
				"ts": "1707584400000"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// 创建使用模拟服务器的客户端
	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// 直接调用内部方法测试解析逻辑
	trades, err := testGetTradeHistory(client, server.URL, integrations.TradeHistoryParams{})

	require.NoError(t, err, "解析成交记录不应出错")
	assert.Len(t, trades, 2, "应返回 2 条成交记录")

	// 验证第一条记录
	assert.Equal(t, "123456", trades[0].ID)
	assert.Equal(t, "BTC-USDT", trades[0].Symbol)
	assert.Equal(t, "buy", trades[0].Side)
	assert.Equal(t, 50000.5, trades[0].Price)
	assert.Equal(t, 0.1, trades[0].Quantity)
	assert.Equal(t, 0.05, trades[0].Fee, "手续费应取绝对值")
	assert.Equal(t, "USDT", trades[0].FeeCoin)
	assert.Equal(t, "okx", trades[0].Source)

	// 验证第二条记录
	assert.Equal(t, "sell", trades[1].Side)
	assert.Equal(t, 3000.0, trades[1].Price)
}

// TestGetTradeHistory_API错误 验证 API 错误处理
func TestGetTradeHistory_API错误(t *testing.T) {
	mockResponse := `{
		"code": "50011",
		"msg": "Invalid parameter",
		"data": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	trades, err := testGetTradeHistory(client, server.URL, integrations.TradeHistoryParams{})

	assert.Error(t, err, "API 错误应返回 error")
	assert.Nil(t, trades, "API 错误时应返回 nil")
	assert.Contains(t, err.Error(), "OKX API 错误")
}

// TestGetDepositHistory_解析充值记录 验证充值历史解析
func TestGetDepositHistory_解析充值记录(t *testing.T) {
	mockResponse := `{
		"code": "0",
		"msg": "",
		"data": [
			{
				"depId": "dep-001",
				"ccy": "USDT",
				"amt": "1000.0",
				"from": "0xabc123",
				"txId": "0xhash123",
				"state": "2",
				"ts": "1707580800000"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	transfers, err := testGetDepositHistory(client, server.URL, integrations.DepositWithdrawParams{})

	require.NoError(t, err)
	assert.Len(t, transfers, 1)
	assert.Equal(t, "deposit", transfers[0].Type)
	assert.Equal(t, "USDT", transfers[0].Coin)
	assert.Equal(t, 1000.0, transfers[0].Amount)
	assert.Equal(t, "completed", transfers[0].Status)
	assert.Equal(t, "okx", transfers[0].Source)
}

// TestGetWithdrawHistory_解析提现记录 验证提现历史解析
func TestGetWithdrawHistory_解析提现记录(t *testing.T) {
	mockResponse := `{
		"code": "0",
		"msg": "",
		"data": [
			{
				"wdId": "wd-001",
				"ccy": "ETH",
				"amt": "5.0",
				"fee": "0.005",
				"to": "0xdef456",
				"txId": "0xhash456",
				"state": "2",
				"ts": "1707580800000"
			},
			{
				"wdId": "wd-002",
				"ccy": "BTC",
				"amt": "0.5",
				"fee": "0.0001",
				"to": "bc1q...",
				"txId": "txhash789",
				"state": "-2",
				"ts": "1707584400000"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "test-key",
		apiSecret:  "test-secret",
		passphrase: "test-pass",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	transfers, err := testGetWithdrawHistory(client, server.URL, integrations.DepositWithdrawParams{})

	require.NoError(t, err)
	assert.Len(t, transfers, 2)

	// 第一条：已完成
	assert.Equal(t, "withdraw", transfers[0].Type)
	assert.Equal(t, "ETH", transfers[0].Coin)
	assert.Equal(t, 5.0, transfers[0].Amount)
	assert.Equal(t, 0.005, transfers[0].Fee)
	assert.Equal(t, "completed", transfers[0].Status)

	// 第二条：已取消
	assert.Equal(t, "BTC", transfers[1].Coin)
	assert.Equal(t, "cancelled", transfers[1].Status)
}

// testGetTradeHistory 辅助函数：使用自定义 URL 测试成交历史解析
func testGetTradeHistory(c *Client, serverURL string, params integrations.TradeHistoryParams) ([]integrations.Trade, error) {
	ctx := context.Background()
	path := "/api/v5/trade/fills-history"

	body, err := testDoRequestWithURL(c, serverURL, path)
	if err != nil {
		return nil, err
	}

	// 复用解析逻辑
	return parseTradeHistory(body, ctx)
}

// testGetDepositHistory 辅助函数
func testGetDepositHistory(c *Client, serverURL string, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	path := "/api/v5/asset/deposit-history"
	body, err := testDoRequestWithURL(c, serverURL, path)
	if err != nil {
		return nil, err
	}
	return parseDepositHistory(body)
}

// testGetWithdrawHistory 辅助函数
func testGetWithdrawHistory(c *Client, serverURL string, params integrations.DepositWithdrawParams) ([]integrations.Transfer, error) {
	path := "/api/v5/asset/withdrawal-history"
	body, err := testDoRequestWithURL(c, serverURL, path)
	if err != nil {
		return nil, err
	}
	return parseWithdrawHistory(body)
}
