// Package etherscan 公共 RPC Gas 价格查询
// 当 Etherscan API Key 未配置时，通过公共 RPC 节点获取 Gas 价格
package etherscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"

	"your-finance/allfi/internal/utils"
)

// rpcRequest JSON-RPC 请求结构
type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// rpcResponse JSON-RPC 响应结构
type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error,omitempty"`
}

// rpcError JSON-RPC 错误
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UnmarshalJSON 支持字符串和标准的对象结构，兼容不规范的 RPC 节点返回
func (e *rpcError) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		e.Code = -1
		e.Message = s
		return nil
	}

	type Alias rpcError
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*e = rpcError(aux)
	return nil
}

// GetGasPriceViaRPC 通过公共 RPC 节点获取 Gas 价格
// 调用 eth_gasPrice 方法，返回与 Etherscan Gas Oracle 兼容的 GasPrice 结构
// 由于 RPC 只返回单一建议价格，Safe/Normal/Fast 按比例估算
func GetGasPriceViaRPC(ctx context.Context, chainName string) (*GasPrice, error) {
	config, ok := SupportedChains[chainName]
	if !ok {
		return nil, fmt.Errorf("不支持的链: %s", chainName)
	}
	rpcUrl := GetRPCURL(ctx, chainName)
	if rpcUrl == "" {
		return nil, fmt.Errorf("链 %s 未配置公共 RPC 端点", chainName)
	}

	// 限流
	if err := utils.WaitForAPI(ctx, config.RateLimitKey+"_rpc"); err != nil {
		return nil, err
	}

	// 构造 JSON-RPC 请求
	reqBody, _ := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      1,
	})

	httpClient := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcUrl, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("构造 RPC 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RPC 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RPC HTTP 错误: 状态码 %d", resp.StatusCode)
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("解析 RPC 响应失败: %v", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC 错误: %s", rpcResp.Error.Message)
	}

	// 解析十六进制 Gas 价格（单位: Wei）
	var hexPrice string
	if err := json.Unmarshal(rpcResp.Result, &hexPrice); err != nil {
		return nil, fmt.Errorf("解析 Gas 价格失败: %v", err)
	}

	gasPriceWei, ok := new(big.Int).SetString(hexPrice[2:], 16) // 去掉 0x 前缀
	if !ok {
		return nil, fmt.Errorf("解析十六进制 Gas 价格失败: %s", hexPrice)
	}

	// Wei → Gwei（除以 1e9）
	gasPriceGwei := new(big.Float).Quo(
		new(big.Float).SetInt(gasPriceWei),
		new(big.Float).SetFloat64(1e9),
	)
	normalGwei, _ := gasPriceGwei.Float64()

	// RPC 只返回单一建议价格，按比例估算 Safe/Fast
	return &GasPrice{
		Safe:    math.Max(normalGwei*0.8, 0.1),
		Normal:  normalGwei,
		Fast:    normalGwei * 1.3,
		BaseFee: 0, // RPC 不直接返回 BaseFee
	}, nil
}

// GetNativeBalanceViaRPC 通过公共 RPC 节点获取原生代币（ETH/BNB等）余额
func GetNativeBalanceViaRPC(ctx context.Context, chainName string, address string) (float64, error) {
	config, ok := SupportedChains[chainName]
	if !ok {
		return 0, fmt.Errorf("不支持的链: %s", chainName)
	}
	rpcUrl := GetRPCURL(ctx, chainName)
	if rpcUrl == "" {
		return 0, fmt.Errorf("链 %s 未配置公共 RPC 端点", chainName)
	}

	// 限流
	if err := utils.WaitForAPI(ctx, config.RateLimitKey+"_rpc"); err != nil {
		return 0, err
	}

	// 构造 JSON-RPC 请求
	reqBody, _ := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
		ID:      1,
	})

	httpClient := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcUrl, bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("构造 RPC 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("RPC 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("RPC HTTP 错误: 状态码 %d", resp.StatusCode)
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("解析 RPC 响应失败: %v", err)
	}
	if rpcResp.Error != nil {
		return 0, fmt.Errorf("RPC 错误: %s", rpcResp.Error.Message)
	}

	// 解析十六进制余额
	var hexBalance string
	if err := json.Unmarshal(rpcResp.Result, &hexBalance); err != nil {
		return 0, fmt.Errorf("解析余额结果失败: %v", err)
	}

	if hexBalance == "0x" {
		return 0, nil
	}

	balanceWei, ok := new(big.Int).SetString(strings.TrimPrefix(hexBalance, "0x"), 16)
	if !ok {
		return 0, fmt.Errorf("解析十六进制余额失败: %s", hexBalance)
	}

	// Wei → ETH (1e18)
	balanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(balanceWei),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)
	balance, _ := balanceEth.Float64()

	return balance, nil
}

// ERC20TokenInfo 主流 ERC20 代币信息
type ERC20TokenInfo struct {
	Symbol   string
	Decimals int
	Contract string
}

// 各链主流 ERC20 代币合约地址
// 只包含交易量最大的稳定币和主流代币，避免过多 RPC 请求
var WellKnownTokens = map[string][]ERC20TokenInfo{
	"ethereum": {
		{Symbol: "USDT", Decimals: 6, Contract: "0xdac17f958d2ee523a2206206994597c13d831ec7"},
		{Symbol: "USDC", Decimals: 6, Contract: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},
		{Symbol: "DAI", Decimals: 18, Contract: "0x6b175474e89094c44da98b954eedeac495271d0f"},
		{Symbol: "WETH", Decimals: 18, Contract: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},
		{Symbol: "WBTC", Decimals: 8, Contract: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"},
		{Symbol: "UNI", Decimals: 18, Contract: "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
		{Symbol: "LINK", Decimals: 18, Contract: "0x514910771af9ca656af840dff83e8264ecf986ca"},
	},
	"bsc": {
		{Symbol: "USDT", Decimals: 18, Contract: "0x55d398326f99059ff775485246999027b3197955"},
		{Symbol: "USDC", Decimals: 18, Contract: "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"},
		{Symbol: "BUSD", Decimals: 18, Contract: "0xe9e7cea3dedca5984780bafc599bd69add087d56"},
		{Symbol: "WBNB", Decimals: 18, Contract: "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"},
		{Symbol: "ETH", Decimals: 18, Contract: "0x2170ed0880ac9a755fd29b2688956bd959f933f8"},
		{Symbol: "BTCB", Decimals: 18, Contract: "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c"},
		{Symbol: "CAKE", Decimals: 18, Contract: "0x0e09fabb73bd3ade0a17ecc321fd13a19e81ce82"},
	},
	"polygon": {
		{Symbol: "USDT", Decimals: 6, Contract: "0xc2132d05d31c914a87c6611c10748aeb04b58e8f"},
		{Symbol: "USDC", Decimals: 6, Contract: "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359"},
		{Symbol: "WETH", Decimals: 18, Contract: "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619"},
		{Symbol: "WMATIC", Decimals: 18, Contract: "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"},
		{Symbol: "WBTC", Decimals: 8, Contract: "0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6"},
	},
	"arbitrum": {
		{Symbol: "USDT", Decimals: 6, Contract: "0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9"},
		{Symbol: "USDC", Decimals: 6, Contract: "0xaf88d065e77c8cc2239327c5edb3a432268e5831"},
		{Symbol: "WETH", Decimals: 18, Contract: "0x82af49447d8a07e3bd95bd0d56f35241523fbab1"},
		{Symbol: "WBTC", Decimals: 8, Contract: "0x2f2a2543b76a4166549f7aab2e75bef0aefc5b0f"},
		{Symbol: "ARB", Decimals: 18, Contract: "0x912ce59144191c1204e64559fe8253a0e49e6548"},
	},
	"optimism": {
		{Symbol: "USDT", Decimals: 6, Contract: "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58"},
		{Symbol: "USDC", Decimals: 6, Contract: "0x0b2c639c533813f4aa9d7837caf62653d097ff85"},
		{Symbol: "WETH", Decimals: 18, Contract: "0x4200000000000000000000000000000000000006"},
		{Symbol: "WBTC", Decimals: 8, Contract: "0x68f180fcce6836688e9084f035309e29bf0a2095"},
		{Symbol: "OP", Decimals: 18, Contract: "0x4200000000000000000000000000000000000042"},
	},
	"base": {
		{Symbol: "USDC", Decimals: 6, Contract: "0x833589fcd6edb6e08f4c7c32d4f71b54bda02913"},
		{Symbol: "WETH", Decimals: 18, Contract: "0x4200000000000000000000000000000000000006"},
		{Symbol: "DAI", Decimals: 18, Contract: "0x50c5725949a6f0c72e6c4a641f24049a917db0cb"},
	},
}

// GetERC20BalancesViaRPC 通过 RPC eth_call 批量查询主流 ERC20 代币余额
// 无需 Etherscan API Key，直接调用合约的 balanceOf(address) 方法
// 返回: map[symbol]balance（已按 decimals 转换）
func GetERC20BalancesViaRPC(ctx context.Context, chainName string, walletAddress string) (map[string]float64, error) {
	tokens, ok := WellKnownTokens[chainName]
	if !ok || len(tokens) == 0 {
		return nil, nil
	}

	rpcUrl := GetRPCURL(ctx, chainName)
	if rpcUrl == "" {
		return nil, fmt.Errorf("链 %s 未配置 RPC 端点", chainName)
	}

	config := SupportedChains[chainName]

	// balanceOf(address) 的函数选择器: 0x70a08231
	// 参数: 地址左填充到 32 字节
	paddedAddr := fmt.Sprintf("%064s", strings.TrimPrefix(strings.ToLower(walletAddress), "0x"))
	callData := "0x70a08231" + paddedAddr

	balances := make(map[string]float64)
	httpClient := &http.Client{Timeout: 15 * time.Second}

	for _, token := range tokens {
		// 限流
		if err := utils.WaitForAPI(ctx, config.RateLimitKey+"_rpc"); err != nil {
			continue
		}

		reqBody, _ := json.Marshal(rpcRequest{
			JSONRPC: "2.0",
			Method:  "eth_call",
			Params: []interface{}{
				map[string]string{
					"to":   token.Contract,
					"data": callData,
				},
				"latest",
			},
			ID: 1,
		})

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, rpcUrl, bytes.NewReader(reqBody))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			continue
		}

		var rpcResp rpcResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if rpcResp.Error != nil {
			continue
		}

		var hexResult string
		if err := json.Unmarshal(rpcResp.Result, &hexResult); err != nil {
			continue
		}

		// 解析余额
		hexStr := strings.TrimPrefix(hexResult, "0x")
		if hexStr == "" || hexStr == "0" {
			continue
		}

		balanceWei, ok := new(big.Int).SetString(hexStr, 16)
		if !ok || balanceWei.Sign() == 0 {
			continue
		}

		// 按 decimals 转换
		divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token.Decimals)), nil))
		balanceFloat := new(big.Float).Quo(new(big.Float).SetInt(balanceWei), divisor)
		bal, _ := balanceFloat.Float64()

		if bal > 0 {
			balances[token.Symbol] = bal
		}
	}

	return balances, nil
}
