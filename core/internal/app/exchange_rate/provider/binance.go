// Package provider Binance 汇率提供者
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
	BinanceCacheTTL = 60 * time.Second
)

// BinanceProvider Binance 数据提供者
type BinanceProvider struct {
	exchange *ccxt.Binance
	cache    *gcache.Cache
}

// NewBinanceProvider 创建 Binance 提供者
func NewBinanceProvider() *BinanceProvider {
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
				g.Log().Info(ctx, "Binance Provider 已启用 SOCKS5 代理", g.Map{
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
					g.Log().Info(ctx, "Binance Provider 已启用 HTTP/HTTPS 代理", g.Map{
						"httpProxy":  config["httpProxy"],
						"httpsProxy": config["httpsProxy"],
					})
				}
			}
		}
	}

	// 创建 Binance 交易所实例
	exchange := ccxt.NewBinance(config)

	return &BinanceProvider{
		exchange: exchange,
		cache:    gcache.New(),
	}
}

func (p *BinanceProvider) Name() string {
	return "Binance"
}

func (p *BinanceProvider) Priority() int {
	return 1 // 最高优先级
}

// FetchRate 获取单个币种汇率（优先级：USDC > USDT）
func (p *BinanceProvider) FetchRate(ctx context.Context, symbol string) (*RateInfo, error) {
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

	// 尝试获取汇率：先 USDC，再 USDT
	quoteSymbols := []string{"USDC", "USDT"}
	var lastErr error

	for _, quote := range quoteSymbols {
		pair := symbol + "/" + quote

		// 尝试从缓存获取
		cacheKey := fmt.Sprintf("binance:ticker:%s", pair)
		if val, err := p.cache.Get(ctx, cacheKey); err == nil && val != nil {
			if rate, ok := val.Val().(*RateInfo); ok {
				return rate, nil
			}
		}

		// 使用 ccxt 获取 ticker
		ticker, err := p.exchange.FetchTicker(pair)
		if err != nil {
			lastErr = err
			g.Log().Debug(ctx, "Binance 交易对不存在，尝试下一个", g.Map{
				"pair":  pair,
				"error": err.Error(),
			})
			continue
		}

		// 成功获取
		price := 0.0
		change := 0.0
		volume := 0.0

		// 解析 ticker 数据
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
		_ = p.cache.Set(ctx, cacheKey, rate, BinanceCacheTTL)

		return rate, nil
	}

	// 所有交易对都失败
	if lastErr != nil {
		return nil, NewProviderError(p.Name(), fmt.Errorf("获取 %s 汇率失败: %w", symbol, lastErr))
	}
	return nil, NewProviderError(p.Name(), fmt.Errorf("不支持的币种: %s", symbol))
}

// FetchRates 批量获取汇率
func (p *BinanceProvider) FetchRates(ctx context.Context) (map[string]*RateInfo, error) {
	rates := make(map[string]*RateInfo)
	symbols := p.SupportedSymbols()

	for _, symbol := range symbols {
		rate, err := p.FetchRate(ctx, symbol)
		if err != nil {
			g.Log().Warning(ctx, "Binance 获取汇率失败", g.Map{
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

func (p *BinanceProvider) IsHealthy(ctx context.Context) bool {
	// 简单健康检查：尝试获取 BTC 价格
	_, err := p.FetchRate(ctx, "BTC")
	return err == nil
}

func (p *BinanceProvider) SupportedSymbols() []string {
	return []string{
		"BTC", "ETH", "BNB", "SOL", "XRP", "ADA",
		"DOGE", "MATIC", "DOT", "SHIB", "TRX", "AVAX",
		"LINK", "ATOM", "UNI", "LTC", "PEPE",
		"USDC", "USDT", "DAI", // 稳定币
	}
}

func (p *BinanceProvider) SupportsSymbol(symbol string) bool {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	for _, s := range p.SupportedSymbols() {
		if s == symbol {
			return true
		}
	}
	return false
}
