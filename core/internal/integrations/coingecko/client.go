// Package coingecko CoinGecko API 客户端
// 获取加密货币价格（免费 API，无需 API Key）
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

const (
	baseURL = "https://api.coingecko.com/api/v3"
)

// Client CoinGecko 客户端
type Client struct {
	apiKey     string // 可选，Pro 版本需要
	httpClient *http.Client
}

// NewClient 创建 CoinGecko 客户端
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// 币种 ID 映射（CoinGecko 使用 ID 而不是 Symbol）
var symbolToID = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"USDT":  "tether",
	"USDC":  "usd-coin",
	"BNB":   "binancecoin",
	"XRP":   "ripple",
	"ADA":   "cardano",
	"DOGE":  "dogecoin",
	"SOL":   "solana",
	"DOT":   "polkadot",
	"MATIC": "matic-network",
	"SHIB":  "shiba-inu",
	"TRX":   "tron",
	"AVAX":  "avalanche-2",
	"LINK":  "chainlink",
	"UNI":   "uniswap",
	"ATOM":  "cosmos",
	"LTC":   "litecoin",
	"ETC":   "ethereum-classic",
	"XMR":   "monero",
	"DAI":   "dai",
	"AAVE":  "aave",
	"MKR":   "maker",
	"COMP":  "compound-governance-token",
	"SUSHI": "sushi",
}

// GetSymbolID 获取币种的 CoinGecko ID
func GetSymbolID(symbol string) string {
	symbol = strings.ToUpper(symbol)
	if id, ok := symbolToID[symbol]; ok {
		return id
	}
	return strings.ToLower(symbol)
}

// GetPrice 获取单个币种价格
func (c *Client) GetPrice(ctx context.Context, symbol string) (float64, error) {
	prices, err := c.GetPrices(ctx, []string{symbol})
	if err != nil {
		return 0, err
	}
	if price, ok := prices[symbol]; ok {
		return price, nil
	}
	return 0, fmt.Errorf("未找到 %s 的价格", symbol)
}

// GetPrices 批量获取币种价格
func (c *Client) GetPrices(ctx context.Context, symbols []string) (map[string]float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "coingecko"); err != nil {
		return nil, err
	}

	// 转换 symbol 为 CoinGecko ID
	ids := make([]string, 0, len(symbols))
	symbolMap := make(map[string]string) // id -> symbol
	for _, symbol := range symbols {
		id := GetSymbolID(symbol)
		ids = append(ids, id)
		symbolMap[id] = strings.ToUpper(symbol)
	}

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd",
		baseURL, strings.Join(ids, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 如果有 API Key，添加到请求头
	if c.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 CoinGecko API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("CoinGecko API 错误: HTTP %d", resp.StatusCode)
	}

	// 解析响应
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 CoinGecko 响应失败: %v", err)
	}

	// 转换回 symbol
	prices := make(map[string]float64)
	for id, priceData := range result {
		if symbol, ok := symbolMap[id]; ok {
			if usdPrice, ok := priceData["usd"]; ok {
				prices[symbol] = usdPrice
			}
		}
	}

	return prices, nil
}

// GetExchangeRate 获取汇率
func (c *Client) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// 如果是稳定币到 USD，返回 1
	stableCoins := []string{"USDT", "USDC", "BUSD", "DAI", "TUSD"}
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	for _, stable := range stableCoins {
		if from == stable && to == "USD" {
			return 1.0, nil
		}
	}

	// 获取 from 币种的 USD 价格
	fromPrice, err := c.GetPrice(ctx, from)
	if err != nil {
		return 0, err
	}

	// 如果目标是 USD，直接返回
	if to == "USD" || to == "USDC" || to == "USDT" {
		return fromPrice, nil
	}

	// 获取 to 币种的 USD 价格
	toPrice, err := c.GetPrice(ctx, to)
	if err != nil {
		return 0, err
	}

	if toPrice == 0 {
		return 0, fmt.Errorf("目标货币 %s 价格为 0", to)
	}

	// 计算汇率
	return fromPrice / toPrice, nil
}

// GetHistoricalPrices 获取指定币种的历史价格并计算区间收益率
// 调用 CoinGecko /coins/{id}/market_chart API 获取历史价格数据
// 返回收益率百分比：(最新价 - 起始价) / 起始价 * 100
func (c *Client) GetHistoricalPrices(ctx context.Context, symbol string, days int) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "coingecko"); err != nil {
		return 0, err
	}

	// 将 symbol 转为 CoinGecko ID
	coinID := GetSymbolID(symbol)

	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=%d",
		baseURL, coinID, days)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	// 如果有 API Key，添加到请求头
	if c.apiKey != "" {
		req.Header.Set("x-cg-pro-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求 CoinGecko 历史价格失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("CoinGecko API 错误: HTTP %d", resp.StatusCode)
	}

	// 解析响应：{"prices": [[timestamp, price], ...]}
	var result struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("解析 CoinGecko 历史价格响应失败: %v", err)
	}

	if len(result.Prices) < 2 {
		return 0, fmt.Errorf("历史价格数据不足（%s，%d天）", symbol, days)
	}

	// 取第一个和最后一个数据点计算收益率
	startPrice := result.Prices[0][1]
	endPrice := result.Prices[len(result.Prices)-1][1]

	if startPrice <= 0 {
		return 0, fmt.Errorf("起始价格为零（%s）", symbol)
	}

	// 收益率 = (最新价 - 起始价) / 起始价 * 100
	returnPct := (endPrice - startPrice) / startPrice * 100
	return returnPct, nil
}

// GetMarketData 获取市场数据（包含 24h 变化等）
func (c *Client) GetMarketData(ctx context.Context, symbols []string) ([]MarketData, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "coingecko"); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		ids = append(ids, GetSymbolID(symbol))
	}

	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&sparkline=false",
		baseURL, strings.Join(ids, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 CoinGecko API 失败: %v", err)
	}
	defer resp.Body.Close()

	var result []MarketData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 CoinGecko 响应失败: %v", err)
	}

	return result, nil
}

// MarketData 市场数据
type MarketData struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCap                float64 `json:"market_cap"`
	TotalVolume              float64 `json:"total_volume"`
	PriceChange24h           float64 `json:"price_change_24h"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

// 确保实现接口
var _ integrations.PriceClient = (*Client)(nil)
