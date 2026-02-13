// Package provider Frankfurter 汇率提供者
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

const (
	// Frankfurter API 基础 URL
	FrankfurterAPIBase = "https://api.frankfurter.dev"

	// 缓存过期时间
	FrankfurterCacheTTL = 5 * time.Minute
)

// FrankfurterProvider Frankfurter API 提供者（仅法币汇率）
type FrankfurterProvider struct {
	baseURL string
	cache   *gcache.Cache
}

// FrankfurterResponse Frankfurter API 响应
type FrankfurterResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// NewFrankfurterProvider 创建 Frankfurter 提供者
func NewFrankfurterProvider() *FrankfurterProvider {
	return &FrankfurterProvider{
		baseURL: FrankfurterAPIBase,
		cache:   gcache.New(),
	}
}

func (p *FrankfurterProvider) Name() string {
	return "Frankfurter"
}

func (p *FrankfurterProvider) Priority() int {
	return 3 // 较低优先级（仅用于法币）
}

// FetchRate 获取单个币种汇率（仅支持 CNY）
func (p *FrankfurterProvider) FetchRate(ctx context.Context, symbol string) (*RateInfo, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// 仅支持 CNY
	if symbol != "CNY" {
		return nil, NewProviderError(p.Name(), fmt.Errorf("仅支持 CNY，不支持: %s", symbol))
	}

	// 尝试从缓存获取
	cacheKey := "frankfurter:rate:CNY"
	if val, err := p.cache.Get(ctx, cacheKey); err == nil && val != nil {
		if rate, ok := val.Val().(*RateInfo); ok {
			return rate, nil
		}
	}

	// 调用 Frankfurter API：获取 USD 对 CNY 的汇率
	url := fmt.Sprintf("%s/v1/latest?base=USD&symbols=CNY", p.baseURL)
	resp, err := g.Client().Get(ctx, url)
	if err != nil {
		return nil, NewProviderError(p.Name(), fmt.Errorf("API 调用失败: %w", err))
	}
	defer resp.Close()

	if resp.StatusCode != 200 {
		return nil, NewProviderError(p.Name(), fmt.Errorf("API 返回错误: %d", resp.StatusCode))
	}

	// 解析响应
	var result FrankfurterResponse
	if err := json.Unmarshal(resp.ReadAll(), &result); err != nil {
		return nil, NewProviderError(p.Name(), fmt.Errorf("解析响应失败: %w", err))
	}

	// 获取 CNY 汇率
	cnyRate, ok := result.Rates["CNY"]
	if !ok {
		return nil, NewProviderError(p.Name(), fmt.Errorf("未找到 CNY 汇率"))
	}

	// Frankfurter 返回的是 1 USD = X CNY
	// 我们需要转换为 1 CNY = Y USD
	priceUSD := 1.0 / cnyRate

	rate := &RateInfo{
		Symbol:      "CNY",
		PriceUSD:    priceUSD,
		Change24H:   0, // Frankfurter 不提供24小时变化
		Volume24H:   0,
		LastUpdated: time.Now(),
		Source:      p.Name(),
	}

	// 缓存结果
	_ = p.cache.Set(ctx, cacheKey, rate, FrankfurterCacheTTL)

	return rate, nil
}

// FetchRates 批量获取汇率（仅返回 CNY）
func (p *FrankfurterProvider) FetchRates(ctx context.Context) (map[string]*RateInfo, error) {
	rate, err := p.FetchRate(ctx, "CNY")
	if err != nil {
		return nil, err
	}

	return map[string]*RateInfo{
		"CNY": rate,
	}, nil
}

func (p *FrankfurterProvider) IsHealthy(ctx context.Context) bool {
	// 健康检查：尝试获取 CNY 汇率
	_, err := p.FetchRate(ctx, "CNY")
	return err == nil
}

func (p *FrankfurterProvider) SupportedSymbols() []string {
	// 仅支持 CNY
	return []string{"CNY"}
}

func (p *FrankfurterProvider) SupportsSymbol(symbol string) bool {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	return symbol == "CNY"
}
