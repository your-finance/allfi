// Package provider Gate.io 汇率提供者
package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

const (
	// 缓存过期时间
	GateioCacheTTL = 60 * time.Second
)

// GateioProvider Gate.io 数据提供者
type GateioProvider struct {
	exchange *ccxt.Gateio
	cache    *gcache.Cache
}

// NewGateioProvider 创建 Gate.io 提供者
func NewGateioProvider() *GateioProvider {
	config := map[string]interface{}{
		"enableRateLimit": true, // 启用请求频率限制
	}

	// 从配置读取代理设置（如果配置文件存在）
	ctx := context.Background()
	if cfg := g.Cfg(); cfg != nil {
		if enabled, err := cfg.Get(ctx, "proxy.enabled"); err == nil && !enabled.IsNil() && enabled.Bool() {
			// 优先使用 SOCKS5 代理
			if socksProxy, err := cfg.Get(ctx, "proxy.socksProxy"); err == nil && !socksProxy.IsEmpty() {
				socksProxyStr := socksProxy.String()
				config["httpProxy"] = socksProxyStr
				config["httpsProxy"] = socksProxyStr
				g.Log().Info(ctx, "Gate.io Provider 已启用 SOCKS5 代理", g.Map{
					"socksProxy": socksProxyStr,
				})
			} else {
				// 否则使用 HTTP/HTTPS 代理
				if httpProxy, err := cfg.Get(ctx, "proxy.httpProxy"); err == nil && !httpProxy.IsEmpty() {
					config["httpProxy"] = httpProxy.String()
				}
				if httpsProxy, err := cfg.Get(ctx, "proxy.httpsProxy"); err == nil && !httpsProxy.IsEmpty() {
					config["httpsProxy"] = httpsProxy.String()
				}
				if config["httpProxy"] != nil || config["httpsProxy"] != nil {
					g.Log().Info(ctx, "Gate.io Provider 已启用 HTTP/HTTPS 代理", g.Map{
						"httpProxy":  config["httpProxy"],
						"httpsProxy": config["httpsProxy"],
					})
				}
			}
		}
	}

	// 创建 Gate.io 交易所实例
	exchange := ccxt.NewGateio(config)

	return &GateioProvider{
		exchange: exchange,
		cache:    gcache.New(),
	}
}

func (p *GateioProvider) Name() string {
	return "Gate.io"
}

func (p *GateioProvider) Priority() int {
	return 2 // 次优先级（Binance 降级方案）
}

// FetchRate 获取单个币种汇率（仅支持 USDT 交易对）
func (p *GateioProvider) FetchRate(ctx context.Context, symbol string) (*RateInfo, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// 稳定币直接返回 1.0
	if symbol == "USDC" || symbol == "USDT" || symbol == "DAI" {
		return &RateInfo{
			Symbol:      symbol,
			PriceUSD:    1.0,
			Change24H:   0,
			Volume24H:   0,
			LastUpdated: time.Now(),
			Source:      p.Name(),
		}, nil
	}

	// Gate.io 使用斜杠分隔交易对，如 BTC/USDT
	pair := symbol + "/USDT"

	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("gateio:ticker:%s", pair)
	if val, err := p.cache.Get(ctx, cacheKey); err == nil && val != nil {
		if rate, ok := val.Val().(*RateInfo); ok {
			return rate, nil
		}
	}

	// 使用 ccxt 获取 ticker
	ticker, err := p.exchange.FetchTicker(pair)
	if err != nil {
		return nil, NewProviderError(p.Name(), fmt.Errorf("获取 %s 汇率失败: %w", symbol, err))
	}

	// 解析 ticker 数据
	price := 0.0
	change := 0.0
	volume := 0.0

	if ticker.Last != nil {
		price = *ticker.Last
	}
	if ticker.Percentage != nil {
		change = *ticker.Percentage
	}
	if ticker.BaseVolume != nil {
		volume = *ticker.BaseVolume
	}

	rate := &RateInfo{
		Symbol:      symbol,
		PriceUSD:    price,
		Change24H:   change,
		Volume24H:   volume,
		LastUpdated: time.Now(),
		Source:      fmt.Sprintf("%s(%s)", p.Name(), pair),
	}

	// 缓存结果
	_ = p.cache.Set(ctx, cacheKey, rate, GateioCacheTTL)

	return rate, nil
}

// FetchRates 批量获取汇率
func (p *GateioProvider) FetchRates(ctx context.Context) (map[string]*RateInfo, error) {
	rates := make(map[string]*RateInfo)
	symbols := p.SupportedSymbols()

	for _, symbol := range symbols {
		rate, err := p.FetchRate(ctx, symbol)
		if err != nil {
			g.Log().Warning(ctx, "Gate.io 获取汇率失败", g.Map{
				"symbol": symbol,
				"error":  err.Error(),
			})
			continue
		}
		rates[symbol] = rate
	}

	if len(rates) == 0 {
		return nil, NewProviderError(p.Name(), fmt.Errorf("所有币种汇率获取失败"))
	}

	return rates, nil
}

func (p *GateioProvider) IsHealthy(ctx context.Context) bool {
	// 简单健康检查：尝试获取 BTC 价格
	_, err := p.FetchRate(ctx, "BTC")
	return err == nil
}

func (p *GateioProvider) SupportedSymbols() []string {
	// Gate.io 支持的主流币种
	return []string{
		"BTC", "ETH", "BNB", "SOL", "XRP", "ADA",
		"DOGE", "MATIC", "DOT", "SHIB", "TRX", "AVAX",
		"LINK", "ATOM", "UNI", "LTC", "PEPE",
		"USDC", "USDT", "DAI",
	}
}

func (p *GateioProvider) SupportsSymbol(symbol string) bool {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	for _, s := range p.SupportedSymbols() {
		if s == symbol {
			return true
		}
	}
	return false
}
