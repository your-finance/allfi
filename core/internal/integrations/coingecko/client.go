// Package coingecko CoinGecko API 客户端
// 获取加密货币价格（免费 API，无需 API Key）
// 内置全局价格缓存（60s TTL）和 HTTP 429 退避重试机制
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

const (
	baseURL = "https://api.coingecko.com/api/v3"

	// priceCacheTTL 价格缓存有效期（60 秒）
	priceCacheTTL = 60 * time.Second

	// max429Retries 429 错误最大重试次数
	max429Retries = 2

	// cooldownAfter429 遇到 429 后全局冷却时间
	cooldownAfter429 = 60 * time.Second
)

// priceCacheEntry 价格缓存条目
type priceCacheEntry struct {
	price     float64
	fetchedAt time.Time
}

// 全局价格缓存（所有 Client 实例共享）
var (
	globalPriceCache = make(map[string]priceCacheEntry)
	priceCacheMu     sync.RWMutex

	// 全局 429 冷却截止时间
	cooldownUntil   time.Time
	cooldownUntilMu sync.RWMutex
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

// getCachedPrice 从缓存获取价格，未命中返回 (0, false)
func getCachedPrice(symbol string) (float64, bool) {
	priceCacheMu.RLock()
	defer priceCacheMu.RUnlock()
	entry, ok := globalPriceCache[symbol]
	if !ok || time.Since(entry.fetchedAt) > priceCacheTTL {
		return 0, false
	}
	return entry.price, true
}

// setCachedPrices 批量写入缓存
func setCachedPrices(prices map[string]float64) {
	priceCacheMu.Lock()
	defer priceCacheMu.Unlock()
	now := time.Now()
	for symbol, price := range prices {
		globalPriceCache[symbol] = priceCacheEntry{price: price, fetchedAt: now}
	}
}

// isInCooldown 检查是否处于 429 冷却期
func isInCooldown() bool {
	cooldownUntilMu.RLock()
	defer cooldownUntilMu.RUnlock()
	return time.Now().Before(cooldownUntil)
}

// setCooldown 设置 429 冷却期
func setCooldown() {
	cooldownUntilMu.Lock()
	defer cooldownUntilMu.Unlock()
	cooldownUntil = time.Now().Add(cooldownAfter429)
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
// 内置缓存：60 秒内不会重复请求同一币种
// 内置 429 退避：遇到限流后自动冷却 60 秒
func (c *Client) GetPrices(ctx context.Context, symbols []string) (map[string]float64, error) {
	result := make(map[string]float64)
	var needFetch []string

	// 1. 先从缓存取，过滤出需要请求的币种
	for _, symbol := range symbols {
		upper := strings.ToUpper(symbol)
		if price, ok := getCachedPrice(upper); ok {
			result[upper] = price
		} else {
			needFetch = append(needFetch, symbol)
		}
	}

	// 全部命中缓存，直接返回
	if len(needFetch) == 0 {
		return result, nil
	}

	// 2. 检查是否处于 429 冷却期
	if isInCooldown() {
		g.Log().Debug(ctx, "CoinGecko 处于 429 冷却期，使用缓存中的已有数据")
		// 冷却期内返回已有缓存数据（可能不完整但不会报错）
		return result, nil
	}

	// 3. 限流等待
	if err := utils.WaitForAPI(ctx, "coingecko"); err != nil {
		return result, err
	}

	// 转换 symbol 为 CoinGecko ID
	ids := make([]string, 0, len(needFetch))
	symbolMap := make(map[string]string) // id -> symbol
	for _, symbol := range needFetch {
		id := GetSymbolID(symbol)
		ids = append(ids, id)
		symbolMap[id] = strings.ToUpper(symbol)
	}

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd",
		baseURL, strings.Join(ids, ","))

	// 4. 请求（带 429 重试）
	fetchedPrices, err := c.doRequestWithRetry(ctx, url)
	if err != nil {
		return result, err
	}

	// 5. 解析并写入缓存
	newPrices := make(map[string]float64)
	for id, priceData := range fetchedPrices {
		if symbol, ok := symbolMap[id]; ok {
			if usdPrice, ok := priceData["usd"]; ok {
				result[symbol] = usdPrice
				newPrices[symbol] = usdPrice
			}
		}
	}

	if len(newPrices) > 0 {
		setCachedPrices(newPrices)
	}

	return result, nil
}

// doRequestWithRetry 执行 HTTP 请求，遇到 429 自动退避重试
func (c *Client) doRequestWithRetry(ctx context.Context, url string) (map[string]map[string]float64, error) {
	var lastErr error

	for attempt := 0; attempt <= max429Retries; attempt++ {
		if attempt > 0 {
			// 指数退避：第 1 次等 10 秒，第 2 次等 30 秒
			backoff := time.Duration(10*(1<<(attempt-1))) * time.Second
			g.Log().Infof(ctx, "CoinGecko 429 退避重试 %d/%d，等待 %v", attempt, max429Retries, backoff)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		if c.apiKey != "" {
			req.Header.Set("x-cg-pro-api-key", c.apiKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("请求 CoinGecko API 失败: %v", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = fmt.Errorf("CoinGecko API 限流: HTTP 429")
			if attempt == max429Retries {
				// 最后一次重试仍然 429，设置全局冷却
				setCooldown()
				g.Log().Warning(ctx, "CoinGecko 连续 429，进入全局冷却期 60 秒")
			}
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("CoinGecko API 错误: HTTP %d", resp.StatusCode)
		}

		var result map[string]map[string]float64
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("解析 CoinGecko 响应失败: %v", err)
		}

		return result, nil
	}

	return nil, lastErr
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
	// 检查 429 冷却期
	if isInCooldown() {
		return 0, fmt.Errorf("CoinGecko 处于限流冷却期")
	}

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

	if resp.StatusCode == http.StatusTooManyRequests {
		setCooldown()
		return 0, fmt.Errorf("CoinGecko API 限流: HTTP 429")
	}

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
	// 检查 429 冷却期
	if isInCooldown() {
		return nil, fmt.Errorf("CoinGecko 处于限流冷却期")
	}

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

	if resp.StatusCode == http.StatusTooManyRequests {
		setCooldown()
		return nil, fmt.Errorf("CoinGecko API 限流: HTTP 429")
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("CoinGecko API 错误: HTTP %d", resp.StatusCode)
	}

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
