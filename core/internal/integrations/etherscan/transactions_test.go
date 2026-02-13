// Package etherscan 链上交易记录测试
package etherscan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestParseTxList_解析普通交易 验证普通交易列表解析逻辑
func TestParseTxList_解析普通交易(t *testing.T) {
	result := txListResponse{
		Status:  "1",
		Message: "OK",
		Result: []struct {
			Hash             string `json:"hash"`
			From             string `json:"from"`
			To               string `json:"to"`
			Value            string `json:"value"`
			Gas              string `json:"gas"`
			GasPrice         string `json:"gasPrice"`
			GasUsed          string `json:"gasUsed"`
			IsError          string `json:"isError"`
			TimeStamp        string `json:"timeStamp"`
			FunctionName     string `json:"functionName"`
			ContractAddress  string `json:"contractAddress"`
		}{
			{
				Hash:      "0xabc123",
				From:      "0xsender",
				To:        "0xreceiver",
				Value:     "1000000000000000000", // 1 ETH
				Gas:       "21000",
				GasPrice:  "20000000000", // 20 Gwei
				GasUsed:   "21000",
				IsError:   "0",
				TimeStamp: "1707580800",
			},
			{
				Hash:      "0xdef456",
				From:      "0xsender",
				To:        "0xreceiver",
				Value:     "500000000000000000", // 0.5 ETH
				Gas:       "21000",
				GasPrice:  "30000000000", // 30 Gwei
				GasUsed:   "21000",
				IsError:   "1", // 失败交易
				TimeStamp: "1707584400",
			},
		},
	}

	txs := parseTxList(result, "ethereum")

	assert.Len(t, txs, 2, "应返回 2 条交易")

	// 验证第一条交易
	assert.Equal(t, "0xabc123", txs[0].Hash)
	assert.Equal(t, "0xsender", txs[0].From)
	assert.Equal(t, "0xreceiver", txs[0].To)
	assert.InDelta(t, 1.0, txs[0].Value, 0.001, "1 ETH")
	assert.InDelta(t, 20.0, txs[0].GasPrice, 0.001, "20 Gwei")
	assert.InDelta(t, 21000.0, txs[0].GasUsed, 0.1)
	assert.False(t, txs[0].IsError, "不是失败交易")
	assert.Equal(t, "ethereum", txs[0].Chain)

	// Gas 费 = 21000 * 20 Gwei / 1e9 = 0.00042 ETH
	assert.InDelta(t, 0.00042, txs[0].GasFee, 0.00001)

	// 验证失败交易
	assert.True(t, txs[1].IsError, "应标记为失败")
}

// TestParseTokenTxList_解析代币转账 验证 ERC20 转账解析逻辑
func TestParseTokenTxList_解析代币转账(t *testing.T) {
	result := tokenTxResponse{
		Status:  "1",
		Message: "OK",
		Result: []struct {
			Hash            string `json:"hash"`
			From            string `json:"from"`
			To              string `json:"to"`
			Value           string `json:"value"`
			TokenName       string `json:"tokenName"`
			TokenSymbol     string `json:"tokenSymbol"`
			TokenDecimal    string `json:"tokenDecimal"`
			ContractAddress string `json:"contractAddress"`
			GasPrice        string `json:"gasPrice"`
			GasUsed         string `json:"gasUsed"`
			TimeStamp       string `json:"timeStamp"`
		}{
			{
				Hash:            "0xtokentx1",
				From:            "0xsender",
				To:              "0xreceiver",
				Value:           "1000000000", // 1000 USDT (6 decimals)
				TokenName:       "Tether USD",
				TokenSymbol:     "USDT",
				TokenDecimal:    "6",
				ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
				GasPrice:        "25000000000",
				GasUsed:         "65000",
				TimeStamp:       "1707580800",
			},
		},
	}

	txs := parseTokenTxList(result, "ethereum")

	assert.Len(t, txs, 1)
	assert.Equal(t, "USDT", txs[0].TokenSymbol)
	assert.Equal(t, "Tether USD", txs[0].TokenName)
	assert.InDelta(t, 1000.0, txs[0].TokenAmount, 0.01, "1000 USDT")
	assert.Equal(t, "0xdac17f958d2ee523a2206206994597c13d831ec7", txs[0].ContractAddress)
	assert.Equal(t, "ethereum", txs[0].Chain)
}

// TestWeiToEth_转换 验证 wei 到 ETH 的转换
func TestWeiToEth_转换(t *testing.T) {
	tests := []struct {
		wei      string
		expected float64
	}{
		{"1000000000000000000", 1.0},      // 1 ETH
		{"500000000000000000", 0.5},        // 0.5 ETH
		{"100000000000000", 0.0001},        // 0.0001 ETH
		{"0", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		result := weiToEth(tt.wei)
		assert.InDelta(t, tt.expected, result, 0.00001, "weiToEth(%s) 应为 %f", tt.wei, tt.expected)
	}
}

// TestWeiToGwei_转换 验证 wei 到 Gwei 的转换
func TestWeiToGwei_转换(t *testing.T) {
	tests := []struct {
		wei      string
		expected float64
	}{
		{"20000000000", 20.0},   // 20 Gwei
		{"1000000000", 1.0},     // 1 Gwei
		{"0", 0},
	}

	for _, tt := range tests {
		result := weiToGwei(tt.wei)
		assert.InDelta(t, tt.expected, result, 0.001)
	}
}

// TestTokenValueToFloat_转换 验证代币值转换
func TestTokenValueToFloat_转换(t *testing.T) {
	tests := []struct {
		value    string
		decimals int
		expected float64
	}{
		{"1000000000", 6, 1000.0},         // 1000 USDT
		{"1000000000000000000", 18, 1.0},  // 1 ETH
		{"0", 18, 0},
	}

	for _, tt := range tests {
		result := tokenValueToFloat(tt.value, tt.decimals)
		assert.InDelta(t, tt.expected, result, 0.001)
	}
}

// TestOnChainTransaction_字段 验证交易结构体
func TestOnChainTransaction_字段(t *testing.T) {
	tx := OnChainTransaction{
		Hash:      "0xhash",
		From:      "0xfrom",
		To:        "0xto",
		Value:     1.5,
		GasFee:    0.0042,
		Chain:     "ethereum",
		Timestamp: time.Now(),
	}

	assert.Equal(t, "0xhash", tx.Hash)
	assert.Equal(t, 1.5, tx.Value)
	assert.Equal(t, "ethereum", tx.Chain)
}
