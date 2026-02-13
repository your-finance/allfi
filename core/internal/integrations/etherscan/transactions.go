// Package etherscan 链上交易记录获取
// 使用 Etherscan txlist 和 tokentx 模块查询地址的链上交易历史
package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"your-finance/allfi/internal/utils"
)

// OnChainTransaction 链上交易记录
type OnChainTransaction struct {
	Hash            string    `json:"hash"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Value           float64   `json:"value"`            // 原生币数量
	Gas             float64   `json:"gas"`              // Gas 用量
	GasPrice        float64   `json:"gas_price"`        // Gas 价格（Gwei）
	GasUsed         float64   `json:"gas_used"`         // 实际 Gas 用量
	GasFee          float64   `json:"gas_fee"`          // Gas 费（原生币计）
	IsError         bool      `json:"is_error"`         // 是否失败
	TokenSymbol     string    `json:"token_symbol"`     // 代币符号（ERC20 转账时有值）
	TokenName       string    `json:"token_name"`       // 代币名称
	TokenAmount     float64   `json:"token_amount"`     // 代币数量
	ContractAddress string    `json:"contract_address"` // 合约地址
	Chain           string    `json:"chain"`            // 链名称
	Timestamp       time.Time `json:"timestamp"`
}

// txListResponse Etherscan txlist API 响应
type txListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Hash             string `json:"hash"`
		From             string `json:"from"`
		To               string `json:"to"`
		Value            string `json:"value"`       // wei
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`    // wei
		GasUsed          string `json:"gasUsed"`
		IsError          string `json:"isError"`     // "0" or "1"
		TimeStamp        string `json:"timeStamp"`   // unix timestamp
		FunctionName     string `json:"functionName"`
		ContractAddress  string `json:"contractAddress"`
	} `json:"result"`
}

// tokenTxResponse Etherscan tokentx API 响应
type tokenTxResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		Hash            string `json:"hash"`
		From            string `json:"from"`
		To              string `json:"to"`
		Value           string `json:"value"`           // 代币最小单位
		TokenName       string `json:"tokenName"`
		TokenSymbol     string `json:"tokenSymbol"`
		TokenDecimal    string `json:"tokenDecimal"`
		ContractAddress string `json:"contractAddress"`
		GasPrice        string `json:"gasPrice"`
		GasUsed         string `json:"gasUsed"`
		TimeStamp       string `json:"timeStamp"`
	} `json:"result"`
}

// GetTransactions 获取普通交易记录
// 使用 Etherscan txlist 模块
func (c *Client) GetTransactions(ctx context.Context, address string, startBlock, endBlock int64) ([]OnChainTransaction, error) {
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return nil, err
	}

	start := "0"
	if startBlock > 0 {
		start = strconv.FormatInt(startBlock, 10)
	}
	end := "99999999"
	if endBlock > 0 {
		end = strconv.FormatInt(endBlock, 10)
	}

	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=%s&endblock=%s&page=1&offset=200&sort=desc&apikey=%s",
		c.baseURL, address, start, end, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Etherscan txlist 失败: %v", err)
	}
	defer resp.Body.Close()

	var result txListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 txlist 响应失败: %v", err)
	}

	return parseTxList(result, c.chainName), nil
}

// parseTxList 解析普通交易列表
func parseTxList(result txListResponse, chain string) []OnChainTransaction {
	txs := make([]OnChainTransaction, 0, len(result.Result))
	for _, tx := range result.Result {
		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		value := weiToEth(tx.Value)
		gasPrice := weiToGwei(tx.GasPrice)
		gasUsed, _ := strconv.ParseFloat(tx.GasUsed, 64)

		// Gas 费 = gasUsed * gasPrice（Gwei） / 1e9
		gasFee := gasUsed * gasPrice / 1e9

		txs = append(txs, OnChainTransaction{
			Hash:      tx.Hash,
			From:      tx.From,
			To:        tx.To,
			Value:     value,
			GasPrice:  gasPrice,
			GasUsed:   gasUsed,
			GasFee:    gasFee,
			IsError:   tx.IsError == "1",
			Chain:     chain,
			Timestamp: time.Unix(ts, 0),
		})
	}
	return txs
}

// GetTokenTransfers 获取 ERC20 代币转账记录
// 使用 Etherscan tokentx 模块
func (c *Client) GetTokenTransfers(ctx context.Context, address string, startBlock, endBlock int64) ([]OnChainTransaction, error) {
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return nil, err
	}

	start := "0"
	if startBlock > 0 {
		start = strconv.FormatInt(startBlock, 10)
	}
	end := "99999999"
	if endBlock > 0 {
		end = strconv.FormatInt(endBlock, 10)
	}

	url := fmt.Sprintf("%s?module=account&action=tokentx&address=%s&startblock=%s&endblock=%s&page=1&offset=200&sort=desc&apikey=%s",
		c.baseURL, address, start, end, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Etherscan tokentx 失败: %v", err)
	}
	defer resp.Body.Close()

	var result tokenTxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 tokentx 响应失败: %v", err)
	}

	return parseTokenTxList(result, c.chainName), nil
}

// parseTokenTxList 解析代币转账列表
func parseTokenTxList(result tokenTxResponse, chain string) []OnChainTransaction {
	txs := make([]OnChainTransaction, 0, len(result.Result))
	for _, tx := range result.Result {
		ts, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		decimals, _ := strconv.Atoi(tx.TokenDecimal)
		if decimals == 0 {
			decimals = 18
		}
		tokenAmount := tokenValueToFloat(tx.Value, decimals)
		gasPrice := weiToGwei(tx.GasPrice)
		gasUsed, _ := strconv.ParseFloat(tx.GasUsed, 64)
		gasFee := gasUsed * gasPrice / 1e9

		txs = append(txs, OnChainTransaction{
			Hash:            tx.Hash,
			From:            tx.From,
			To:              tx.To,
			TokenSymbol:     tx.TokenSymbol,
			TokenName:       tx.TokenName,
			TokenAmount:     tokenAmount,
			ContractAddress: tx.ContractAddress,
			GasPrice:        gasPrice,
			GasUsed:         gasUsed,
			GasFee:          gasFee,
			Chain:           chain,
			Timestamp:       time.Unix(ts, 0),
		})
	}
	return txs
}

// weiToEth 将 wei 字符串转换为 ETH 浮点数
func weiToEth(weiStr string) float64 {
	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		return 0
	}
	eth := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	f, _ := eth.Float64()
	return f
}

// weiToGwei 将 wei 字符串转换为 Gwei 浮点数
func weiToGwei(weiStr string) float64 {
	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		return 0
	}
	gwei := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetInt(big.NewInt(1e9)),
	)
	f, _ := gwei.Float64()
	return f
}

// tokenValueToFloat 将代币最小单位转换为浮点数
func tokenValueToFloat(valueStr string, decimals int) float64 {
	value, ok := new(big.Int).SetString(valueStr, 10)
	if !ok {
		return 0
	}
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	result := new(big.Float).Quo(
		new(big.Float).SetInt(value),
		new(big.Float).SetInt(divisor),
	)
	f, _ := result.Float64()
	return f
}
