// Package binance Binance 交易历史测试
package binance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"your-finance/allfi/internal/integrations"
)

// TestTradeHistoryParams_默认值 验证查询参数默认值处理
func TestTradeHistoryParams_默认值(t *testing.T) {
	params := integrations.TradeHistoryParams{}

	assert.Empty(t, params.Symbol, "默认 Symbol 应为空")
	assert.True(t, params.StartTime.IsZero(), "默认 StartTime 应为零值")
	assert.True(t, params.EndTime.IsZero(), "默认 EndTime 应为零值")
	assert.Equal(t, 0, params.Limit, "默认 Limit 应为 0")
}

// TestTradeHistoryParams_自定义值 验证查询参数设置
func TestTradeHistoryParams_自定义值(t *testing.T) {
	now := time.Now()
	params := integrations.TradeHistoryParams{
		Symbol:    "ETHUSDT",
		StartTime: now.Add(-24 * time.Hour),
		EndTime:   now,
		Limit:     100,
	}

	assert.Equal(t, "ETHUSDT", params.Symbol)
	assert.False(t, params.StartTime.IsZero())
	assert.False(t, params.EndTime.IsZero())
	assert.Equal(t, 100, params.Limit)
}

// TestTrade_结构体字段 验证 Trade 结构体字段完整性
func TestTrade_结构体字段(t *testing.T) {
	trade := integrations.Trade{
		ID:        "12345",
		Symbol:    "BTCUSDT",
		Side:      "buy",
		Price:     50000.0,
		Quantity:  0.1,
		Fee:       0.001,
		FeeCoin:   "BNB",
		Timestamp: time.Now(),
		Source:    "binance",
	}

	assert.Equal(t, "12345", trade.ID)
	assert.Equal(t, "BTCUSDT", trade.Symbol)
	assert.Equal(t, "buy", trade.Side)
	assert.Equal(t, 50000.0, trade.Price)
	assert.Equal(t, 0.1, trade.Quantity)
	assert.Equal(t, 0.001, trade.Fee)
	assert.Equal(t, "BNB", trade.FeeCoin)
	assert.Equal(t, "binance", trade.Source)
}

// TestTransfer_充值记录 验证 Transfer 充值类型
func TestTransfer_充值记录(t *testing.T) {
	transfer := integrations.Transfer{
		ID:        "tx-001",
		Type:      "deposit",
		Coin:      "USDT",
		Amount:    1000.0,
		Fee:       0,
		Address:   "0xabc123",
		TxHash:    "0xhash123",
		Status:    "completed",
		Timestamp: time.Now(),
		Source:    "binance",
	}

	assert.Equal(t, "deposit", transfer.Type)
	assert.Equal(t, "USDT", transfer.Coin)
	assert.Equal(t, 1000.0, transfer.Amount)
	assert.Equal(t, "completed", transfer.Status)
	assert.Equal(t, "binance", transfer.Source)
}

// TestTransfer_提现记录 验证 Transfer 提现类型
func TestTransfer_提现记录(t *testing.T) {
	transfer := integrations.Transfer{
		ID:        "wd-001",
		Type:      "withdraw",
		Coin:      "ETH",
		Amount:    5.0,
		Fee:       0.005,
		Address:   "0xdef456",
		TxHash:    "0xhash456",
		Status:    "pending",
		Timestamp: time.Now(),
		Source:    "binance",
	}

	assert.Equal(t, "withdraw", transfer.Type)
	assert.Equal(t, "ETH", transfer.Coin)
	assert.Equal(t, 5.0, transfer.Amount)
	assert.Equal(t, 0.005, transfer.Fee)
	assert.Equal(t, "pending", transfer.Status)
}

// TestParseFloat_正确解析 验证 parseFloat 工具函数
func TestParseFloat_正确解析(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"100.5", 100.5},
		{"0", 0},
		{"", 0},
		{"abc", 0},
		{"1000000.123456", 1000000.123456},
	}

	for _, tt := range tests {
		result := parseFloat(tt.input)
		assert.Equal(t, tt.expected, result, "parseFloat(%q) 应返回 %f", tt.input, tt.expected)
	}
}

// TestDepositWithdrawParams_参数设置 验证充提查询参数
func TestDepositWithdrawParams_参数设置(t *testing.T) {
	now := time.Now()
	params := integrations.DepositWithdrawParams{
		Coin:      "BTC",
		StartTime: now.Add(-7 * 24 * time.Hour),
		EndTime:   now,
		Limit:     50,
	}

	assert.Equal(t, "BTC", params.Coin)
	assert.False(t, params.StartTime.IsZero())
	assert.Equal(t, 50, params.Limit)
}
