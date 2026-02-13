// Package logic 基准对比业务逻辑
// 从快照数据计算用户收益率，并与 BTC/ETH 基准指数对比
package logic

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/benchmark/model"
	"your-finance/allfi/internal/app/benchmark/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/integrations/yahoo"
	"your-finance/allfi/internal/model/entity"
)

// sBenchmark 基准对比服务实现
type sBenchmark struct{}

// New 创建基准对比服务实例
func New() service.IBenchmark {
	return &sBenchmark{}
}

// GetBenchmark 获取用户收益率与基准指数的对比
func (s *sBenchmark) GetBenchmark(ctx context.Context, period string) (*model.BenchmarkResult, error) {
	days := parsePeriodDays(period)

	// 获取用户收益率（从快照数据计算）
	userReturn, startDate, endDate, err := s.calculateUserReturn(ctx, days)
	if err != nil {
		g.Log().Warning(ctx, "计算用户收益率失败，使用默认值0", "error", err)
		userReturn = 0
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -days)
	}

	// 获取基准指数收益率
	benchmarks := s.calculateBenchmarks(ctx, days, userReturn)

	return &model.BenchmarkResult{
		Period:     period,
		UserReturn: userReturn,
		Benchmarks: benchmarks,
		StartDate:  startDate.Format("2006-01-02"),
		EndDate:    endDate.Format("2006-01-02"),
	}, nil
}

// calculateUserReturn 从快照数据计算用户收益率
func (s *sBenchmark) calculateUserReturn(ctx context.Context, days int) (float64, time.Time, time.Time, error) {
	startTime := time.Now().AddDate(0, 0, -days)

	// 直接从 DAO 查询快照
	var snapshots []entity.AssetSnapshots
	err := assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, consts.GetUserID(ctx)).
		WhereGTE(assetDao.AssetSnapshots.Columns().SnapshotTime, startTime).
		OrderDesc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return 0, time.Time{}, time.Time{}, gerror.Wrap(err, "查询快照数据失败")
	}

	if len(snapshots) < 2 {
		return 0, time.Time{}, time.Time{}, gerror.New("快照数据不足（至少需要 2 个快照）")
	}

	// 快照按时间降序：[0] = 最新，[n-1] = 最旧
	latest := snapshots[0]
	oldest := snapshots[len(snapshots)-1]

	if oldest.TotalValueUsd <= 0 {
		return 0, oldest.SnapshotTime, latest.SnapshotTime, gerror.New("起始资产为零")
	}

	// 收益率 = (最新 - 最旧) / 最旧 * 100
	returnPct := (latest.TotalValueUsd - oldest.TotalValueUsd) / oldest.TotalValueUsd * 100

	return returnPct, oldest.SnapshotTime, latest.SnapshotTime, nil
}

// calculateBenchmarks 计算各基准指数的收益率
// 通过 CoinGecko 历史价格 API 获取 BTC/ETH 的区间收益率
// 通过 Yahoo Finance API 获取 S&P 500 的区间收益率
func (s *sBenchmark) calculateBenchmarks(ctx context.Context, days int, userReturn float64) []model.BenchmarkIndex {
	cryptoSymbols := []struct {
		symbol string
		name   string
	}{
		{"BTC", "Bitcoin"},
		{"ETH", "Ethereum"},
	}

	// 创建 CoinGecko 客户端（免费 API，无需 Key）
	cgClient := coingecko.NewClient("")

	var benchmarks []model.BenchmarkIndex

	// 获取加密货币基准收益率（BTC/ETH）
	for _, sym := range cryptoSymbols {
		// 调用 CoinGecko 获取历史价格并计算收益率
		returnPct, err := cgClient.GetHistoricalPrices(ctx, sym.symbol, days)
		if err != nil {
			g.Log().Warning(ctx, "获取基准历史价格失败，使用默认值0",
				"symbol", sym.symbol, "days", days, "error", err)
			returnPct = 0.0
		}

		benchmarks = append(benchmarks, model.BenchmarkIndex{
			Name:   sym.name,
			Symbol: sym.symbol,
			Return: returnPct,
			Diff:   userReturn - returnPct,
		})
	}

	// 获取 S&P 500 基准收益率（通过 Yahoo Finance）
	yahooClient := yahoo.NewClient()
	sp500Return, err := yahooClient.GetHistoricalReturn(ctx, "^GSPC", days)
	if err != nil {
		g.Log().Warning(ctx, "获取 S&P 500 历史收益率失败，使用默认值0",
			"days", days, "error", err)
		sp500Return = 0.0
	}

	benchmarks = append(benchmarks, model.BenchmarkIndex{
		Name:   "S&P 500",
		Symbol: "SPX",
		Return: sp500Return,
		Diff:   userReturn - sp500Return,
	})

	return benchmarks
}

// parsePeriodDays 将周期字符串转为天数
func parsePeriodDays(period string) int {
	switch period {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	case "1y":
		return 365
	default:
		return 30
	}
}
