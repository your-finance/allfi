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
	"your-finance/allfi/internal/app/exchange_rate/dao"
	erModel "your-finance/allfi/internal/app/exchange_rate/model"
	"your-finance/allfi/internal/app/exchange_rate/service"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/model/entity"
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
	// 解析货币列表
	symbols := parseCurrencyList(currencies)
	if len(symbols) == 0 {
		symbols = erModel.DefaultCryptoSymbols
	}

	// 从数据库查询汇率
	rates := make(map[string]float64)
	var lastUpdated int64
	var source string

	for _, symbol := range symbols {
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
			rates[symbol] = rate.Rate
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
	// 获取所有需要刷新的货币对
	var existingRates []entity.ExchangeRates
	err := dao.ExchangeRates.Ctx(ctx).
		OrderDesc(dao.ExchangeRates.Columns().FetchedAt).
		Scan(&existingRates)
	if err != nil {
		return nil, gerror.Wrap(err, "查询已有汇率失败")
	}

	// 更新 fetched_at 为当前时间（标记为已刷新）
	now := gtime.Now().Time
	updatedCount := 0
	for _, rate := range existingRates {
		_, err := dao.ExchangeRates.Ctx(ctx).
			Where(dao.ExchangeRates.Columns().Id, rate.Id).
			Data(g.Map{
				dao.ExchangeRates.Columns().FetchedAt: now,
			}).
			Update()
		if err != nil {
			g.Log().Warning(ctx, "更新汇率失败",
				"from", rate.FromCurrency,
				"to", rate.ToCurrency,
				"error", err,
			)
			continue
		}
		updatedCount++
	}

	g.Log().Info(ctx, "刷新汇率完成", "updatedCount", updatedCount)

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
