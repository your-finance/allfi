// Package etherscan Etherscan API 客户端
// 查询 ETH 余额和 ERC20 代币余额
package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

const (
	etherscanBaseURL = "https://api.etherscan.io/api"
)

// Client Etherscan 客户端
type Client struct {
	apiKey     string
	baseURL    string
	chainName  string
	httpClient *http.Client
}

// NewClient 创建 Etherscan 客户端（Ethereum 主网）
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    etherscanBaseURL,
		chainName:  "ethereum",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewBscScanClient 创建 BscScan 客户端（BSC 链）
func NewBscScanClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://api.bscscan.com/api",
		chainName:  "bsc",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewChainClient 通用 EVM 链客户端工厂函数
// 根据链名称从 SupportedChains 中查找配置并创建客户端
func NewChainClient(chainName string, apiKey string) (*Client, error) {
	config, ok := SupportedChains[chainName]
	if !ok {
		return nil, fmt.Errorf("不支持的链: %s", chainName)
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    config.BaseURL,
		chainName:  config.Name,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// NewArbiscanClient 创建 Arbiscan 客户端（Arbitrum 链）
func NewArbiscanClient(apiKey string) *Client {
	c, _ := NewChainClient("arbitrum", apiKey)
	return c
}

// NewOptimismClient 创建 Optimism Explorer 客户端（Optimism 链）
func NewOptimismClient(apiKey string) *Client {
	c, _ := NewChainClient("optimism", apiKey)
	return c
}

// NewPolygonscanClient 创建 Polygonscan 客户端（Polygon 链）
func NewPolygonscanClient(apiKey string) *Client {
	c, _ := NewChainClient("polygon", apiKey)
	return c
}

// NewBasescanClient 创建 Basescan 客户端（Base 链）
func NewBasescanClient(apiKey string) *Client {
	c, _ := NewChainClient("base", apiKey)
	return c
}

// GetChainName 获取链名称
func (c *Client) GetChainName() string {
	return c.chainName
}

// getRateLimitKey 获取当前链对应的限流键名
func (c *Client) getRateLimitKey() string {
	if config, ok := SupportedChains[c.chainName]; ok {
		return config.RateLimitKey
	}
	return "etherscan"
}

// EtherscanResponse Etherscan API 响应
type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

// TokenBalanceResponse 代币余额响应
type TokenBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		ContractAddress string `json:"contractAddress"`
		TokenName       string `json:"tokenName"`
		TokenSymbol     string `json:"tokenSymbol"`
		TokenDecimal    string `json:"tokenDecimal"`
		Balance         string `json:"balance"`
	} `json:"result"`
}

// GetNativeBalance 获取原生代币余额（ETH/BNB 等）
func (c *Client) GetNativeBalance(ctx context.Context, address string) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s?module=account&action=balance&address=%s&tag=latest&apikey=%s",
		c.baseURL, address, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求 Etherscan API 失败: %v", err)
	}
	defer resp.Body.Close()

	var result EtherscanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("解析 Etherscan 响应失败: %v", err)
	}

	if result.Status != "1" {
		return 0, fmt.Errorf("Etherscan API 错误: %s", result.Message)
	}

	// 将 wei 转换为 ETH（18 位小数）
	balanceWei, ok := new(big.Int).SetString(result.Result, 10)
	if !ok {
		return 0, fmt.Errorf("解析余额失败")
	}

	// 转换为 ETH
	balanceEth := new(big.Float).Quo(
		new(big.Float).SetInt(balanceWei),
		new(big.Float).SetInt(big.NewInt(1e18)),
	)

	balance, _ := balanceEth.Float64()
	return balance, nil
}

// GetTokenBalances 获取 ERC20 代币余额
func (c *Client) GetTokenBalances(ctx context.Context, address string) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return nil, err
	}

	// 使用 tokentx 接口获取代币交易历史，从中提取代币列表
	url := fmt.Sprintf("%s?module=account&action=tokentx&address=%s&page=1&offset=100&sort=desc&apikey=%s",
		c.baseURL, address, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Etherscan API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析代币交易记录
	var txResult struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			ContractAddress string `json:"contractAddress"`
			TokenName       string `json:"tokenName"`
			TokenSymbol     string `json:"tokenSymbol"`
			TokenDecimal    string `json:"tokenDecimal"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&txResult); err != nil {
		return nil, fmt.Errorf("解析 Etherscan 响应失败: %v", err)
	}

	// 收集唯一的代币合约
	tokenContracts := make(map[string]struct {
		Name     string
		Symbol   string
		Decimals int
	})

	for _, tx := range txResult.Result {
		if _, exists := tokenContracts[tx.ContractAddress]; !exists {
			decimals, _ := strconv.Atoi(tx.TokenDecimal)
			tokenContracts[tx.ContractAddress] = struct {
				Name     string
				Symbol   string
				Decimals int
			}{
				Name:     tx.TokenName,
				Symbol:   tx.TokenSymbol,
				Decimals: decimals,
			}
		}
	}

	// 查询每个代币的余额
	var balances []integrations.Balance
	for contractAddr, token := range tokenContracts {
		balance, err := c.getTokenBalance(ctx, address, contractAddr, token.Decimals)
		if err != nil {
			continue // 忽略单个代币查询失败
		}

		if balance <= 0 {
			continue
		}

		bal := integrations.Balance{
			Symbol: token.Symbol,
			Name:   token.Name,
			Free:   balance,
			Total:  balance,
		}

		// 检查是否为已知 DeFi 协议 Token，标注协议和资产类型
		if info, ok := LookupDeFiToken(contractAddr); ok {
			bal.Protocol = info.Protocol
			bal.AssetType = info.AssetType
		}

		balances = append(balances, bal)
	}

	return balances, nil
}

// getTokenBalance 获取单个代币余额
func (c *Client) getTokenBalance(ctx context.Context, address, contractAddress string, decimals int) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest&apikey=%s",
		c.baseURL, contractAddress, address, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result EtherscanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Status != "1" {
		return 0, fmt.Errorf("API 错误: %s", result.Message)
	}

	// 转换余额
	balanceRaw, ok := new(big.Int).SetString(result.Result, 10)
	if !ok {
		return 0, fmt.Errorf("解析余额失败")
	}

	if decimals == 0 {
		decimals = 18 // 默认 18 位小数
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	balanceFloat := new(big.Float).Quo(
		new(big.Float).SetInt(balanceRaw),
		new(big.Float).SetInt(divisor),
	)

	balance, _ := balanceFloat.Float64()
	return balance, nil
}

// 常见 ERC20 代币合约地址
var CommonTokens = map[string]string{
	"USDT":  "0xdac17f958d2ee523a2206206994597c13d831ec7",
	"USDC":  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
	"WBTC":  "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
	"DAI":   "0x6b175474e89094c44da98b954eedeac495271d0f",
	"LINK":  "0x514910771af9ca656af840dff83e8264ecf986ca",
	"UNI":   "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
	"AAVE":  "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9",
	"SHIB":  "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
}

// GasOracleResult Gas Oracle 返回结构
type GasOracleResult struct {
	SafeGasPrice    string `json:"SafeGasPrice"`
	ProposeGasPrice string `json:"ProposeGasPrice"`
	FastGasPrice    string `json:"FastGasPrice"`
	SuggestBaseFee  string `json:"suggestBaseFee"`
}

// GasOracleResponse Gas Oracle API 响应
type GasOracleResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  GasOracleResult `json:"result"`
}

// GasPrice 解析后的 Gas 价格（单位: Gwei）
type GasPrice struct {
	Safe    float64 `json:"safe"`
	Normal  float64 `json:"normal"`
	Fast    float64 `json:"fast"`
	BaseFee float64 `json:"baseFee"`
}

// GetGasPrice 获取当前 Gas 价格（通过 Gas Oracle API）
func (c *Client) GetGasPrice(ctx context.Context) (*GasPrice, error) {
	if err := utils.WaitForAPI(ctx, c.getRateLimitKey()); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s?module=gastracker&action=gasoracle&apikey=%s", c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GasOracleResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("Gas Oracle API 错误: %s", result.Message)
	}

	safe, _ := strconv.ParseFloat(result.Result.SafeGasPrice, 64)
	normal, _ := strconv.ParseFloat(result.Result.ProposeGasPrice, 64)
	fast, _ := strconv.ParseFloat(result.Result.FastGasPrice, 64)
	baseFee, _ := strconv.ParseFloat(result.Result.SuggestBaseFee, 64)

	return &GasPrice{
		Safe:    safe,
		Normal:  normal,
		Fast:    fast,
		BaseFee: baseFee,
	}, nil
}

// 确保实现接口
var _ integrations.BlockchainClient = (*Client)(nil)
