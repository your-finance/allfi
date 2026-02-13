// Package provider 本地配置提供者
package provider

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// LocalProvider 本地配置提供者（兜底策略）
type LocalProvider struct {
	rates map[string]*RateInfo
}

// NewLocalProvider 创建本地配置提供者
func NewLocalProvider() *LocalProvider {
	now := time.Now()
	return &LocalProvider{
		rates: map[string]*RateInfo{
			// 稳定币：固定为 1.0
			"USDT": {Symbol: "USDT", PriceUSD: 1.0, Source: "Local", LastUpdated: now},
			"USDC": {Symbol: "USDC", PriceUSD: 1.0, Source: "Local", LastUpdated: now},
			"DAI":  {Symbol: "DAI", PriceUSD: 1.0, Source: "Local", LastUpdated: now},

			// 主流币：参考价格（需要定期手动更新）
			"BTC": {Symbol: "BTC", PriceUSD: 100000.0, Source: "Local", LastUpdated: now},
			"ETH": {Symbol: "ETH", PriceUSD: 3500.0, Source: "Local", LastUpdated: now},
			"BNB": {Symbol: "BNB", PriceUSD: 650.0, Source: "Local", LastUpdated: now},
			"SOL": {Symbol: "SOL", PriceUSD: 200.0, Source: "Local", LastUpdated: now},
			"XRP": {Symbol: "XRP", PriceUSD: 2.5, Source: "Local", LastUpdated: now},

			// 其他主流币
			"ADA":   {Symbol: "ADA", PriceUSD: 1.0, Source: "Local", LastUpdated: now},
			"DOGE":  {Symbol: "DOGE", PriceUSD: 0.35, Source: "Local", LastUpdated: now},
			"MATIC": {Symbol: "MATIC", PriceUSD: 0.5, Source: "Local", LastUpdated: now},
			"DOT":   {Symbol: "DOT", PriceUSD: 7.0, Source: "Local", LastUpdated: now},
			"SHIB":  {Symbol: "SHIB", PriceUSD: 0.000024, Source: "Local", LastUpdated: now},
			"TRX":   {Symbol: "TRX", PriceUSD: 0.25, Source: "Local", LastUpdated: now},
			"AVAX":  {Symbol: "AVAX", PriceUSD: 40.0, Source: "Local", LastUpdated: now},
			"LINK":  {Symbol: "LINK", PriceUSD: 20.0, Source: "Local", LastUpdated: now},
			"ATOM":  {Symbol: "ATOM", PriceUSD: 8.0, Source: "Local", LastUpdated: now},
			"UNI":   {Symbol: "UNI", PriceUSD: 10.0, Source: "Local", LastUpdated: now},
			"LTC":   {Symbol: "LTC", PriceUSD: 100.0, Source: "Local", LastUpdated: now},
			"PEPE":  {Symbol: "PEPE", PriceUSD: 0.000018, Source: "Local", LastUpdated: now},

			// 法币：人民币
			"CNY": {Symbol: "CNY", PriceUSD: 0.14, Source: "Local", LastUpdated: now}, // 1 CNY ≈ 0.14 USD
		},
	}
}

func (p *LocalProvider) Name() string {
	return "Local"
}

func (p *LocalProvider) Priority() int {
	return 999 // 最低优先级，仅作为兜底
}

func (p *LocalProvider) FetchRate(ctx context.Context, symbol string) (*RateInfo, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	rate, ok := p.rates[symbol]
	if !ok {
		return nil, NewProviderError(p.Name(), fmt.Errorf("本地配置不支持币种: %s", symbol))
	}

	// 返回拷贝，避免被修改
	rateCopy := *rate
	return &rateCopy, nil
}

func (p *LocalProvider) FetchRates(ctx context.Context) (map[string]*RateInfo, error) {
	rates := make(map[string]*RateInfo)

	for symbol, rate := range p.rates {
		rateCopy := *rate
		rates[symbol] = &rateCopy
	}

	return rates, nil
}

func (p *LocalProvider) IsHealthy(ctx context.Context) bool {
	return true // 本地配置始终健康
}

func (p *LocalProvider) SupportedSymbols() []string {
	symbols := make([]string, 0, len(p.rates))
	for symbol := range p.rates {
		symbols = append(symbols, symbol)
	}
	return symbols
}

func (p *LocalProvider) SupportsSymbol(symbol string) bool {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	_, ok := p.rates[symbol]
	return ok
}
