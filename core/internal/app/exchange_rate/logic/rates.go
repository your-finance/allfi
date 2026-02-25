// Package logic 汇率模块业务逻辑
// 基于 exchange_rates 表（from_currency/to_currency/rate/source/fetched_at）
package logic

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	exchangeRateApi "your-finance/allfi/api/v1/exchange_rate"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/exchange_rate/dao"
	erModel "your-finance/allfi/internal/app/exchange_rate/model"
	"your-finance/allfi/internal/app/exchange_rate/model/entity"
	"your-finance/allfi/internal/app/exchange_rate/provider"
	"your-finance/allfi/internal/app/exchange_rate/service"
	"your-finance/allfi/internal/integrations/coingecko"
)

// sExchangeRate 汇率服务实现
type sExchangeRate struct{}

// New 创建汇率服务实例
func New() service.IExchangeRate {
	return &sExchangeRate{}
}

// GetRates 获取汇率
//
// 功能说明:
// 1. 解析要查询的货币列表
// 2. 从 exchange_rates 表查询以 USD 为基准的汇率
// 3. 如果数据过期（>5分钟），尝试刷新
func (s *sExchangeRate) GetRates(ctx context.Context, currencies string) (*exchangeRateApi.GetCurrentRes, error) {
	// 解析货币列表；未指定时默认查询加密货币 + 法币
	symbols := parseCurrencyList(currencies)
	if len(symbols) == 0 {
		symbols = append(append([]string{}, erModel.DefaultCryptoSymbols...), erModel.DefaultFiatSymbols...)
	}

	// 从数据库查询汇率
	rates := make(map[string]float64)
	var lastUpdated int64
	var source string

	// 法币列表：这些货币在 DB 中存储为 from=CNY/to=USD（即 1 CNY = 0.14 USD），
	// 但前端期望的是 USD→法币 的汇率（即 1 USD = 7.2 CNY），需要取倒数
	fiatSet := map[string]bool{
		"CNY": true, "EUR": true, "GBP": true, "JPY": true,
		"SGD": true, "HKD": true, "AUD": true, "CAD": true,
	}

	for _, symbol := range symbols {
		if symbol == "USDC" || symbol == "USDT" || symbol == "DAI" || symbol == "USD" {
			rates[symbol] = 1.0
			continue
		}

		var rate entity.ExchangeRates
		err := dao.ExchangeRates.Ctx(ctx).
			Where(dao.ExchangeRates.Columns().FromCurrency, symbol).
			Where(dao.ExchangeRates.Columns().ToCurrency, "USD").
			OrderDesc(dao.ExchangeRates.Columns().FetchedAt).
			Scan(&rate)
		if err != nil {
			// 尝试反向查询
			err = dao.ExchangeRates.Ctx(ctx).
				Where(dao.ExchangeRates.Columns().FromCurrency, "USD").
				Where(dao.ExchangeRates.Columns().ToCurrency, symbol).
				OrderDesc(dao.ExchangeRates.Columns().FetchedAt).
				Scan(&rate)
			if err != nil {
				g.Log().Debug(ctx, "未找到汇率", "symbol", symbol)
				continue
			}
			// 反向汇率
			if rate.Rate > 0 {
				rates[symbol] = 1.0 / rate.Rate
			}
		} else {
			// 法币需要取倒数：DB 存的是 CNY→USD（0.14），前端需要 USD→CNY（7.2）
			if fiatSet[symbol] && rate.Rate > 0 {
				rates[symbol] = 1.0 / rate.Rate
			} else {
				rates[symbol] = rate.Rate
			}
		}

		// 记录最后更新时间
		fetchedAt := rate.FetchedAt.UnixMilli()
		if fetchedAt > lastUpdated {
			lastUpdated = fetchedAt
			source = rate.Source
		}
	}

	// 判断是否缓存数据
	isCached := false
	if lastUpdated > 0 {
		elapsed := time.Since(time.UnixMilli(lastUpdated))
		isCached = elapsed < time.Duration(erModel.CacheTTLSeconds)*time.Second
	}

	if source == "" {
		source = "database"
	}

	return &exchangeRateApi.GetCurrentRes{
		Rates:       rates,
		Base:        "USD",
		LastUpdated: lastUpdated,
		Source:      source,
		IsCached:    isCached,
	}, nil
}

// GetPrices 获取加密货币价格
//
// 功能说明:
// 1. 从 exchange_rates 表查询 {symbol} -> USD 的汇率作为价格
// 2. 通过 CoinGecko 获取 24 小时涨跌幅
// 3. 返回价格列表
func (s *sExchangeRate) GetPrices(ctx context.Context, symbols string) (*exchangeRateApi.GetPricesRes, error) {
	// 解析币种列表
	symbolList := parseCurrencyList(symbols)
	if len(symbolList) == 0 {
		symbolList = erModel.DefaultCryptoSymbols
	}

	// 通过 CoinGecko 批量获取 24 小时涨跌幅
	// 构建 symbol -> Change24h 映射
	change24hMap := make(map[string]float64)
	cgClient := coingecko.NewClient("")
	marketData, err := cgClient.GetMarketData(ctx, symbolList)
	if err != nil {
		// 获取失败不阻断主流程，仅记录日志，Change24h 保持为 0
		g.Log().Warning(ctx, "获取 CoinGecko 市场数据失败", "error", err)
	} else {
		for _, md := range marketData {
			// GetMarketData 返回的 Symbol 是小写，需转为大写匹配
			change24hMap[strings.ToUpper(md.Symbol)] = md.PriceChangePercentage24h
		}
	}

	prices := make([]exchangeRateApi.PriceItem, 0, len(symbolList))

	for _, symbol := range symbolList {
		if symbol == "USDC" || symbol == "USDT" || symbol == "DAI" || symbol == "USD" {
			prices = append(prices, exchangeRateApi.PriceItem{
				Symbol:      symbol,
				PriceUSD:    1.0,
				Change24h:   0,
				LastUpdated: gtime.Now().UnixMilli(),
			})
			continue
		}

		var rate entity.ExchangeRates
		err := dao.ExchangeRates.Ctx(ctx).
			Where(dao.ExchangeRates.Columns().FromCurrency, symbol).
			Where(dao.ExchangeRates.Columns().ToCurrency, "USD").
			OrderDesc(dao.ExchangeRates.Columns().FetchedAt).
			Scan(&rate)
		if err != nil {
			continue
		}

		prices = append(prices, exchangeRateApi.PriceItem{
			Symbol:      symbol,
			PriceUSD:    rate.Rate,
			Change24h:   change24hMap[symbol],
			LastUpdated: rate.FetchedAt.UnixMilli(),
		})
	}

	return &exchangeRateApi.GetPricesRes{Prices: prices}, nil
}

// RefreshRates 刷新汇率缓存
//
// 功能说明:
// 1. 从 Provider（Binance/Gate.io/Frankfurter）获取最新汇率
// 2. 更新 exchange_rates 表
func (s *sExchangeRate) RefreshRates(ctx context.Context) (*exchangeRateApi.RefreshRes, error) {
	// 获取我们需要刷新的货币对
	// 1. 默认的加密货币列表
	symbols := append([]string{}, erModel.DefaultCryptoSymbols...)
	// 2. 默认的法币列表 (CNY)
	symbols = append(symbols, "CNY")

	// 3. 从数据库查询曾经获取过的符号
	var existingCurrencies []string
	dao.ExchangeRates.Ctx(ctx).Distinct().Fields(dao.ExchangeRates.Columns().FromCurrency).Scan(&existingCurrencies)

	// 4. 从 asset_details 表提取用户实际持有的资产符号，动态添加到刷新列表
	var assetSymbols []string
	assetDao.AssetDetails.Ctx(ctx).Distinct().Fields(assetDao.AssetDetails.Columns().AssetSymbol).Scan(&assetSymbols)

	// 合并并去重
	symbolSet := make(map[string]bool)
	for _, sym := range symbols {
		symbolSet[sym] = true
	}
	for _, sym := range existingCurrencies {
		// 避免将 USD 当作从货币
		if sym != "USD" {
			symbolSet[sym] = true
		}
	}
	for _, sym := range assetSymbols {
		upper := strings.ToUpper(sym)
		if upper != "" && upper != "USD" {
			symbolSet[upper] = true
		}
	}

	uniqueSymbols := make([]string, 0, len(symbolSet))
	for sym := range symbolSet {
		uniqueSymbols = append(uniqueSymbols, sym)
	}

	// 从 ProviderManager 并发获取最新汇率
	pm := provider.GetProviderManager()
	rates, _, err := pm.FetchRates(ctx, uniqueSymbols)
	if err != nil && len(rates) == 0 {
		g.Log().Warning(ctx, "获取部分或全部最新汇率失败，但将继续处理稳定币", "error", err)
		if rates == nil {
			rates = make(map[string]*provider.RateInfo)
		}
	}

	// usdc=usd 确保稳定币永远有数据，无论 Provider 是否正常
	for _, sym := range uniqueSymbols {
		if sym == "USDC" || sym == "USDT" || sym == "DAI" || sym == "USD" {
			rates[sym] = &provider.RateInfo{
				Symbol:   sym,
				PriceUSD: 1.0,
				Source:   "System",
			}
		}
	}

	// 插入 exchange_rates 表（保留历史记录，只新增）
	now := gtime.Now().Time
	insertedCount := 0
	for sym, rate := range rates {
		if rate == nil || rate.PriceUSD <= 0 {
			continue
		}

		_, dbErr := dao.ExchangeRates.Ctx(ctx).Data(g.Map{
			dao.ExchangeRates.Columns().FromCurrency: sym,
			dao.ExchangeRates.Columns().ToCurrency:   "USD",
			dao.ExchangeRates.Columns().Rate:         rate.PriceUSD,
			dao.ExchangeRates.Columns().Source:       rate.Source,
			dao.ExchangeRates.Columns().FetchedAt:    now,
			// createdAt 和 updatedAt 由 ORM 处理
		}).Insert()

		if dbErr != nil {
			g.Log().Warning(ctx, "插入汇率历史记录失败",
				"from", sym,
				"error", dbErr,
			)
			continue
		}
		insertedCount++
	}

	g.Log().Info(ctx, "刷新汇率完成", "insertedCount", insertedCount, "totalSymbols", len(uniqueSymbols))

	return &exchangeRateApi.RefreshRes{
		Message: "汇率刷新完成",
	}, nil
}

// GetHistory 获取历史汇率
//
// 功能说明:
// 1. 计算时间范围（当前时间 - days 天）
// 2. 查询 exchange_rates 表中匹配的历史记录
// 3. 转换为历史汇率列表返回
func (s *sExchangeRate) GetHistory(ctx context.Context, base, quote string, days int) (*exchangeRateApi.GetHistoryRes, error) {
	if days <= 0 {
		days = 30
	}
	if base == "" {
		base = "USD"
	}
	if quote == "" {
		return nil, gerror.New("目标货币（quote）不能为空")
	}

	// 计算起始时间
	startTime := time.Now().AddDate(0, 0, -days)

	// 查询历史汇率记录
	var rates []entity.ExchangeRates
	err := dao.ExchangeRates.Ctx(ctx).
		Where(dao.ExchangeRates.Columns().FromCurrency, base).
		Where(dao.ExchangeRates.Columns().ToCurrency, quote).
		WhereGTE(dao.ExchangeRates.Columns().FetchedAt, startTime).
		OrderAsc(dao.ExchangeRates.Columns().FetchedAt).
		Scan(&rates)
	if err != nil {
		return nil, gerror.Wrap(err, "查询历史汇率失败")
	}

	// 如果正向查询无结果，尝试反向查询（交换 from/to，汇率取倒数）
	reverse := false
	if len(rates) == 0 {
		err = dao.ExchangeRates.Ctx(ctx).
			Where(dao.ExchangeRates.Columns().FromCurrency, quote).
			Where(dao.ExchangeRates.Columns().ToCurrency, base).
			WhereGTE(dao.ExchangeRates.Columns().FetchedAt, startTime).
			OrderAsc(dao.ExchangeRates.Columns().FetchedAt).
			Scan(&rates)
		if err != nil {
			return nil, gerror.Wrap(err, "查询反向历史汇率失败")
		}
		reverse = true
	}

	// 转换为历史汇率列表
	history := make([]exchangeRateApi.RateHistoryItem, 0, len(rates))
	for _, r := range rates {
		rate := r.Rate
		if reverse && rate > 0 {
			rate = 1.0 / rate
		}
		history = append(history, exchangeRateApi.RateHistoryItem{
			Date:  r.FetchedAt.Format("2006-01-02"),
			Rate:  rate,
			Base:  base,
			Quote: quote,
		})
	}

	return &exchangeRateApi.GetHistoryRes{
		History: history,
		Base:    base,
		Quote:   quote,
		Days:    days,
	}, nil
}

// parseCurrencyList 解析逗号分隔的货币列表
func parseCurrencyList(currencies string) []string {
	if currencies == "" {
		return nil
	}
	parts := strings.Split(currencies, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		symbol := strings.TrimSpace(strings.ToUpper(part))
		if symbol != "" {
			result = append(result, symbol)
		}
	}
	return result
}
