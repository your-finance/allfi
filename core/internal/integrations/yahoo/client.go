// Package yahoo Yahoo Finance API 客户端
// 获取法币汇率和股票/加密货币价格
package yahoo

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
	// Yahoo Finance API（非官方，免费）
	baseURL = "https://query1.finance.yahoo.com/v8/finance/chart"
)

// Client Yahoo Finance 客户端
type Client struct {
	httpClient *http.Client
}

// NewClient 创建 Yahoo Finance 客户端
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ChartResponse Yahoo Finance 图表响应
type ChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				Currency           string  `json:"currency"`
				Symbol             string  `json:"symbol"`
			} `json:"meta"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// GetPrice 获取价格
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

// GetPrices 批量获取价格
func (c *Client) GetPrices(ctx context.Context, symbols []string) (map[string]float64, error) {
	prices := make(map[string]float64)

	for _, symbol := range symbols {
		price, err := c.fetchPrice(ctx, symbol)
		if err != nil {
			continue // 忽略单个失败
		}
		prices[symbol] = price
	}

	return prices, nil
}

// fetchPrice 获取单个价格
func (c *Client) fetchPrice(ctx context.Context, symbol string) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "yahoo"); err != nil {
		return 0, err
	}

	// 转换符号格式
	yahooSymbol := convertToYahooSymbol(symbol)

	url := fmt.Sprintf("%s/%s?interval=1d&range=1d", baseURL, yahooSymbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	// 设置 User-Agent（Yahoo 可能会拒绝没有 UA 的请求）
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求 Yahoo Finance API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("Yahoo Finance API 错误: HTTP %d", resp.StatusCode)
	}

	var result ChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("解析 Yahoo Finance 响应失败: %v", err)
	}

	if result.Chart.Error != nil {
		return 0, fmt.Errorf("Yahoo Finance API 错误: %s", result.Chart.Error.Description)
	}

	if len(result.Chart.Result) == 0 {
		return 0, fmt.Errorf("未找到 %s 的数据", symbol)
	}

	return result.Chart.Result[0].Meta.RegularMarketPrice, nil
}

// GetExchangeRate 获取汇率
func (c *Client) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "yahoo"); err != nil {
		return 0, err
	}

	// 构建货币对符号
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	// 特殊处理稳定币
	if isStableCoin(from) {
		from = "USD"
	}
	if isStableCoin(to) {
		to = "USD"
	}

	// 如果相同，返回 1
	if from == to {
		return 1.0, nil
	}

	// Yahoo 货币对格式：EURUSD=X
	symbol := from + to + "=X"

	price, err := c.fetchPrice(ctx, symbol)
	if err != nil {
		return 0, err
	}

	return price, nil
}

// GetForexRates 获取多个货币对的汇率
func (c *Client) GetForexRates(ctx context.Context, baseCurrency string, targetCurrencies []string) (map[string]float64, error) {
	rates := make(map[string]float64)

	for _, target := range targetCurrencies {
		rate, err := c.GetExchangeRate(ctx, baseCurrency, target)
		if err != nil {
			continue
		}
		rates[target] = rate
	}

	return rates, nil
}

// convertToYahooSymbol 转换为 Yahoo Finance 符号格式
func convertToYahooSymbol(symbol string) string {
	symbol = strings.ToUpper(symbol)

	// 加密货币
	cryptoSymbols := map[string]string{
		"BTC":   "BTC-USD",
		"ETH":   "ETH-USD",
		"USDT":  "USDT-USD",
		"USDC":  "USDC-USD",
		"BNB":   "BNB-USD",
		"XRP":   "XRP-USD",
		"ADA":   "ADA-USD",
		"DOGE":  "DOGE-USD",
		"SOL":   "SOL-USD",
		"DOT":   "DOT-USD",
		"MATIC": "MATIC-USD",
		"SHIB":  "SHIB-USD",
		"TRX":   "TRX-USD",
		"AVAX":  "AVAX-USD",
		"LINK":  "LINK-USD",
		"UNI":   "UNI-USD",
		"LTC":   "LTC-USD",
	}

	if yahooSymbol, ok := cryptoSymbols[symbol]; ok {
		return yahooSymbol
	}

	// 法币对 USD
	if isFiat(symbol) && symbol != "USD" {
		return symbol + "=X"
	}

	return symbol
}

// isStableCoin 判断是否为稳定币
func isStableCoin(symbol string) bool {
	stableCoins := []string{"USDT", "USDC", "BUSD", "DAI", "TUSD", "USDP", "GUSD"}
	symbol = strings.ToUpper(symbol)
	for _, s := range stableCoins {
		if s == symbol {
			return true
		}
	}
	return false
}

// isFiat 判断是否为法币
func isFiat(symbol string) bool {
	fiats := []string{"USD", "CNY", "EUR", "JPY", "GBP", "KRW", "HKD", "TWD", "SGD", "AUD", "CAD", "CHF"}
	symbol = strings.ToUpper(symbol)
	for _, f := range fiats {
		if f == symbol {
			return true
		}
	}
	return false
}

// HistoricalChartResponse Yahoo Finance 历史图表响应（含时间序列数据）
type HistoricalChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				Currency           string  `json:"currency"`
				Symbol             string  `json:"symbol"`
			} `json:"meta"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// GetHistoricalReturn 获取指定标的的历史区间收益率
// 调用 Yahoo Finance chart API 获取历史收盘价，计算收益率百分比
// symbol: Yahoo Finance 标的代码（如 ^GSPC 代表 S&P 500）
// days: 回溯天数
// 返回: 收益率百分比 = (最新价 - 起始价) / 起始价 * 100
func (c *Client) GetHistoricalReturn(ctx context.Context, symbol string, days int) (float64, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "yahoo"); err != nil {
		return 0, err
	}

	// 将天数转为 Yahoo range 参数
	rangeStr := fmt.Sprintf("%dd", days)
	if days >= 365 {
		rangeStr = "1y"
	} else if days >= 180 {
		rangeStr = "6mo"
	} else if days >= 90 {
		rangeStr = "3mo"
	}

	url := fmt.Sprintf("%s/%s?range=%s&interval=1d", baseURL, symbol, rangeStr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求 Yahoo Finance 历史数据失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("Yahoo Finance API 错误: HTTP %d", resp.StatusCode)
	}

	var result HistoricalChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("解析 Yahoo Finance 历史数据响应失败: %v", err)
	}

	if result.Chart.Error != nil {
		return 0, fmt.Errorf("Yahoo Finance API 错误: %s", result.Chart.Error.Description)
	}

	if len(result.Chart.Result) == 0 {
		return 0, fmt.Errorf("未找到 %s 的历史数据", symbol)
	}

	chartResult := result.Chart.Result[0]
	if len(chartResult.Indicators.Quote) == 0 || len(chartResult.Indicators.Quote[0].Close) < 2 {
		return 0, fmt.Errorf("历史价格数据不足（%s，%d天）", symbol, days)
	}

	closePrices := chartResult.Indicators.Quote[0].Close

	// 找到第一个非零收盘价（跳过可能的空数据点）
	var startPrice float64
	for _, p := range closePrices {
		if p > 0 {
			startPrice = p
			break
		}
	}

	// 找到最后一个非零收盘价
	var endPrice float64
	for i := len(closePrices) - 1; i >= 0; i-- {
		if closePrices[i] > 0 {
			endPrice = closePrices[i]
			break
		}
	}

	if startPrice <= 0 {
		return 0, fmt.Errorf("起始价格为零（%s）", symbol)
	}

	// 收益率 = (最新价 - 起始价) / 起始价 * 100
	returnPct := (endPrice - startPrice) / startPrice * 100
	return returnPct, nil
}

// 确保实现接口
var _ integrations.PriceClient = (*Client)(nil)
