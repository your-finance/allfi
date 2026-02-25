package provider

import (
	"context"

	"your-finance/allfi/internal/integrations/yahoo"
)

// YahooProvider Yahoo Finance API 汇率数据源
type YahooProvider struct{}

// NewYahooProvider 创建 YahooProvider
func NewYahooProvider() *YahooProvider {
	return &YahooProvider{}
}

// Name 名称
func (p *YahooProvider) Name() string {
	return "Yahoo Finance"
}

// Priority 优先级
func (p *YahooProvider) Priority() int {
	return 3 // 法币获取兜底，CNY 首选
}

// IsHealthy 健康状态
func (p *YahooProvider) IsHealthy(ctx context.Context) bool {
	return true
}

// SupportsSymbol 是否支持
func (p *YahooProvider) SupportsSymbol(symbol string) bool {
	fiatList := map[string]bool{
		"CNY": true, "EUR": true, "GBP": true, "JPY": true,
		"SGD": true, "HKD": true, "AUD": true, "CAD": true,
	}
	return fiatList[symbol]
}

// SupportedSymbols 支持的法币列表
func (p *YahooProvider) SupportedSymbols() []string {
	return []string{"CNY", "EUR", "GBP", "JPY", "SGD", "HKD", "AUD", "CAD"}
}

// FetchRates 批量获取汇率
func (p *YahooProvider) FetchRates(ctx context.Context) (map[string]*RateInfo, error) {
	symbols := p.SupportedSymbols()
	rates := make(map[string]*RateInfo)
	for _, sym := range symbols {
		rate, err := p.FetchRate(ctx, sym)
		if err == nil {
			rates[sym] = rate
		}
	}
	return rates, nil
}

// FetchRate 获取汇率
func (p *YahooProvider) FetchRate(ctx context.Context, symbol string) (*RateInfo, error) {
	client := yahoo.NewClient()

	// Yahoo 获取的时候
	rate, err := client.GetExchangeRate(ctx, symbol, "USD")
	if err != nil {
		return nil, err
	}

	return &RateInfo{
		Symbol:   symbol,
		PriceUSD: rate,
		Source:   p.Name(),
	}, nil
}
