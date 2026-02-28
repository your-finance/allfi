// Package logic 风险管理模块 - 业务逻辑实现
package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"gonum.org/v1/gonum/stat"

	assetDao "your-finance/allfi/internal/app/asset/dao"
	assetEntity "your-finance/allfi/internal/app/asset/model/entity"
	"your-finance/allfi/internal/app/risk/dao"
	riskEntity "your-finance/allfi/internal/app/risk/model/entity"
	"your-finance/allfi/internal/app/risk/model"
	"your-finance/allfi/internal/app/risk/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/coingecko"
)

// sRisk 风险管理服务实现
type sRisk struct{}

// New 创建风险管理服务实例
func New() service.IRisk {
	return &sRisk{}
}

// GetLatestMetrics 获取最新风险指标
func (s *sRisk) GetLatestMetrics(ctx context.Context) (*model.RiskMetrics, error) {
	userID := consts.GetUserID(ctx)

	// 从数据库获取最新风险指标
	var entity riskEntity.RiskMetrics
	err := dao.RiskMetrics.Ctx(ctx).
		Where(dao.RiskMetrics.Columns().UserId, userID).
		OrderDesc(dao.RiskMetrics.Columns().MetricDate).
		Limit(1).
		Scan(&entity)
	if err != nil {
		return nil, gerror.Wrap(err, "查询最新风险指标失败")
	}

	// 如果没有数据，返回计算结果
	if entity.Id == 0 {
		return s.CalculateMetrics(ctx, 30)
	}

	// 转换为业务模型
	riskLevel := model.GetRiskLevel(float64(entity.Volatility), float64(entity.MaxDrawdown))
	return &model.RiskMetrics{
		MetricDate:          entity.MetricDate.Format("2006-01-02"),
		PortfolioValue:      float64(entity.PortfolioValue),
		Var95:               float64(entity.Var95),
		Var99:               float64(entity.Var99),
		SharpeRatio:         float64(entity.SharpeRatio),
		SortinoRatio:        float64(entity.SortinoRatio),
		MaxDrawdown:         float64(entity.MaxDrawdown),
		MaxDrawdownDuration: entity.MaxDrawdownDuration,
		Beta:                float64(entity.Beta),
		Volatility:          float64(entity.Volatility),
		DownsideDeviation:   float64(entity.DownsideDeviation),
		CalculationPeriod:   entity.CalculationPeriod,
		RiskLevel:           string(riskLevel),
	}, nil
}

// GetHistoryMetrics 获取历史风险指标
func (s *sRisk) GetHistoryMetrics(ctx context.Context, days int) ([]*model.RiskMetrics, error) {
	if days <= 0 {
		days = 30
	}
	if days > 365 {
		days = 365
	}

	userID := consts.GetUserID(ctx)
	startDate := time.Now().AddDate(0, 0, -days)

	// 从数据库查询历史风险指标
	var entities []*riskEntity.RiskMetrics
	err := dao.RiskMetrics.Ctx(ctx).
		Where(dao.RiskMetrics.Columns().UserId, userID).
		Where(dao.RiskMetrics.Columns().MetricDate+" >= ?", startDate).
		OrderDesc(dao.RiskMetrics.Columns().MetricDate).
		Scan(&entities)
	if err != nil {
		return nil, gerror.Wrap(err, "查询历史风险指标失败")
	}

	// 转换为业务模型
	result := make([]*model.RiskMetrics, 0, len(entities))
	for _, entity := range entities {
		riskLevel := model.GetRiskLevel(float64(entity.Volatility), float64(entity.MaxDrawdown))
		result = append(result, &model.RiskMetrics{
			MetricDate:          entity.MetricDate.Format("2006-01-02"),
			PortfolioValue:      float64(entity.PortfolioValue),
			Var95:               float64(entity.Var95),
			Var99:               float64(entity.Var99),
			SharpeRatio:         float64(entity.SharpeRatio),
			SortinoRatio:        float64(entity.SortinoRatio),
			MaxDrawdown:         float64(entity.MaxDrawdown),
			MaxDrawdownDuration: entity.MaxDrawdownDuration,
			Beta:                float64(entity.Beta),
			Volatility:          float64(entity.Volatility),
			DownsideDeviation:   float64(entity.DownsideDeviation),
			CalculationPeriod:   entity.CalculationPeriod,
			RiskLevel:           string(riskLevel),
		})
	}

	g.Log().Info(ctx, "查询历史风险指标成功", "userID", userID, "days", days, "count", len(result))

	return result, nil
}

// CalculateMetrics 计算风险指标
func (s *sRisk) CalculateMetrics(ctx context.Context, period int) (*model.RiskMetrics, error) {
	if period <= 0 {
		period = 30
	}
	if period < 7 {
		return nil, gerror.New("计算周期至少需要 7 天")
	}
	if period > 365 {
		period = 365
	}

	userID := consts.GetUserID(ctx)

	// 1. 获取资产快照数据
	startDate := time.Now().AddDate(0, 0, -period)
	var snapshots []*assetEntity.AssetSnapshots
	err := assetDao.AssetSnapshots.Ctx(ctx).
		Where(assetDao.AssetSnapshots.Columns().UserId, userID).
		Where(assetDao.AssetSnapshots.Columns().SnapshotTime+" >= ?", startDate).
		OrderAsc(assetDao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产快照失败")
	}

	if len(snapshots) < 7 {
		return nil, gerror.New("快照数据不足（至少需要 7 个快照）")
	}

	// 2. 计算日收益率序列
	returns := s.calculateReturns(snapshots)
	if len(returns) < 2 {
		return nil, gerror.New("收益率数据不足")
	}

	// 3. 获取 BTC 基准收益率（用于计算 Beta）
	btcReturns, err := s.getBTCReturns(ctx, period)
	if err != nil {
		g.Log().Warning(ctx, "获取 BTC 收益率失败，Beta 将设为 0", "error", err)
		btcReturns = make([]float64, len(returns))
	}

	// 4. 计算各项风险指标
	result := s.calculateRiskMetrics(returns, btcReturns)

	// 5. 构建返回结果
	currentValue := float64(snapshots[len(snapshots)-1].TotalValueUsd)
	riskLevel := model.GetRiskLevel(result.Volatility, result.MaxDrawdown)

	metrics := &model.RiskMetrics{
		MetricDate:          time.Now().Format("2006-01-02"),
		PortfolioValue:      math.Round(currentValue*100) / 100,
		Var95:               math.Round(result.Var95*100) / 100,
		Var99:               math.Round(result.Var99*100) / 100,
		SharpeRatio:         math.Round(result.SharpeRatio*100) / 100,
		SortinoRatio:        math.Round(result.SortinoRatio*100) / 100,
		MaxDrawdown:         math.Round(result.MaxDrawdown*100) / 100,
		MaxDrawdownDuration: result.MaxDrawdownDays,
		Beta:                math.Round(result.Beta*100) / 100,
		Volatility:          math.Round(result.Volatility*100) / 100,
		DownsideDeviation:   math.Round(result.DownsideDeviation*100) / 100,
		CalculationPeriod:   period,
		RiskLevel:           string(riskLevel),
	}

	g.Log().Info(ctx, "风险指标计算完成",
		"period", period,
		"snapshots", len(snapshots),
		"sharpeRatio", metrics.SharpeRatio,
		"maxDrawdown", metrics.MaxDrawdown,
		"volatility", metrics.Volatility,
	)

	// 6. 保存到数据库
	metricDate, _ := time.Parse("2006-01-02", metrics.MetricDate)
	_, err = dao.RiskMetrics.Ctx(ctx).Data(g.Map{
		dao.RiskMetrics.Columns().UserId:              userID,
		dao.RiskMetrics.Columns().MetricDate:          metricDate,
		dao.RiskMetrics.Columns().PortfolioValue:      metrics.PortfolioValue,
		dao.RiskMetrics.Columns().Var95:               metrics.Var95,
		dao.RiskMetrics.Columns().Var99:               metrics.Var99,
		dao.RiskMetrics.Columns().SharpeRatio:         metrics.SharpeRatio,
		dao.RiskMetrics.Columns().SortinoRatio:        metrics.SortinoRatio,
		dao.RiskMetrics.Columns().MaxDrawdown:         metrics.MaxDrawdown,
		dao.RiskMetrics.Columns().MaxDrawdownDuration: metrics.MaxDrawdownDuration,
		dao.RiskMetrics.Columns().Beta:                metrics.Beta,
		dao.RiskMetrics.Columns().Volatility:          metrics.Volatility,
		dao.RiskMetrics.Columns().DownsideDeviation:   metrics.DownsideDeviation,
		dao.RiskMetrics.Columns().CalculationPeriod:   metrics.CalculationPeriod,
	}).OnConflict(dao.RiskMetrics.Columns().UserId, dao.RiskMetrics.Columns().MetricDate).Save()
	if err != nil {
		g.Log().Warning(ctx, "保存风险指标失败", "error", err)
	}

	return metrics, nil
}

// calculateReturns 计算日收益率序列
func (s *sRisk) calculateReturns(snapshots []*assetEntity.AssetSnapshots) []float64 {
	if len(snapshots) < 2 {
		return []float64{}
	}

	returns := make([]float64, 0, len(snapshots)-1)
	for i := 1; i < len(snapshots); i++ {
		prev := float64(snapshots[i-1].TotalValueUsd)
		curr := float64(snapshots[i].TotalValueUsd)
		if prev > 0 {
			ret := (curr - prev) / prev
			returns = append(returns, ret)
		}
	}
	return returns
}

// getBTCReturns 获取 BTC 基准收益率
func (s *sRisk) getBTCReturns(ctx context.Context, days int) ([]float64, error) {
	// 获取 BTC 历史价格数据
	coinID := coingecko.GetSymbolID("BTC")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=usd&days=%d", coinID, days)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, gerror.Wrap(err, "请求 BTC 历史价格失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, gerror.Newf("CoinGecko API 错误: HTTP %d", resp.StatusCode)
	}

	// 解析响应：{"prices": [[timestamp, price], ...]}
	var result struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, gerror.Wrap(err, "解析 BTC 历史价格响应失败")
	}

	if len(result.Prices) < 2 {
		return nil, gerror.New("BTC 价格数据不足")
	}

	// 计算日收益率
	returns := make([]float64, 0, len(result.Prices)-1)
	for i := 1; i < len(result.Prices); i++ {
		if result.Prices[i-1][1] > 0 {
			ret := (result.Prices[i][1] - result.Prices[i-1][1]) / result.Prices[i-1][1]
			returns = append(returns, ret)
		}
	}

	return returns, nil
}

// calculateRiskMetrics 计算所有风险指标
func (s *sRisk) calculateRiskMetrics(returns, btcReturns []float64) model.CalculationResult {
	result := model.CalculationResult{}

	// 1. VaR (Value at Risk) - 历史模拟法
	result.Var95 = s.calculateVaR(returns, 0.95)
	result.Var99 = s.calculateVaR(returns, 0.99)

	// 2. 波动率（年化）
	result.Volatility = s.calculateVolatility(returns)

	// 3. 夏普比率（假设无风险利率为 0）
	result.SharpeRatio = s.calculateSharpeRatio(returns, result.Volatility)

	// 4. 下行偏差
	result.DownsideDeviation = s.calculateDownsideDeviation(returns)

	// 5. 索提诺比率
	result.SortinoRatio = s.calculateSortinoRatio(returns, result.DownsideDeviation)

	// 6. 最大回撤
	result.MaxDrawdown, result.MaxDrawdownDays = s.calculateMaxDrawdown(returns)

	// 7. Beta 系数（相对 BTC）
	result.Beta = s.calculateBeta(returns, btcReturns)

	return result
}

// calculateVaR 计算 VaR（历史模拟法）
// confidence: 置信度（如 0.95 表示 95%）
func (s *sRisk) calculateVaR(returns []float64, confidence float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 复制并排序收益率
	sorted := make([]float64, len(returns))
	copy(sorted, returns)
	sort.Float64s(sorted)

	// 计算分位数索引
	index := int((1 - confidence) * float64(len(sorted)))
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	// VaR 为负收益率的绝对值（百分比）
	var_ := -sorted[index] * 100
	return var_
}

// calculateVolatility 计算波动率（年化）
func (s *sRisk) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// 计算标准差
	stdDev := stat.StdDev(returns, nil)

	// 年化波动率 = 日波动率 * sqrt(365)
	annualizedVol := stdDev * math.Sqrt(365) * 100
	return annualizedVol
}

// calculateSharpeRatio 计算夏普比率
// 假设无风险利率为 0
func (s *sRisk) calculateSharpeRatio(returns []float64, volatility float64) float64 {
	if len(returns) == 0 || volatility == 0 {
		return 0
	}

	// 计算平均收益率
	meanReturn := stat.Mean(returns, nil)

	// 年化平均收益率
	annualizedReturn := meanReturn * 365 * 100

	// 夏普比率 = (年化收益率 - 无风险利率) / 年化波动率
	sharpe := annualizedReturn / volatility
	return sharpe
}

// calculateDownsideDeviation 计算下行偏差
// 只考虑负收益率的标准差
func (s *sRisk) calculateDownsideDeviation(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 筛选负收益率
	negativeReturns := make([]float64, 0)
	for _, r := range returns {
		if r < 0 {
			negativeReturns = append(negativeReturns, r)
		}
	}

	if len(negativeReturns) < 2 {
		return 0
	}

	// 计算负收益率的标准差
	stdDev := stat.StdDev(negativeReturns, nil)

	// 年化下行偏差
	annualizedDD := stdDev * math.Sqrt(365) * 100
	return annualizedDD
}

// calculateSortinoRatio 计算索提诺比率
func (s *sRisk) calculateSortinoRatio(returns []float64, downsideDeviation float64) float64 {
	if len(returns) == 0 || downsideDeviation == 0 {
		return 0
	}

	// 计算平均收益率
	meanReturn := stat.Mean(returns, nil)

	// 年化平均收益率
	annualizedReturn := meanReturn * 365 * 100

	// 索提诺比率 = (年化收益率 - 无风险利率) / 下行偏差
	sortino := annualizedReturn / downsideDeviation
	return sortino
}

// calculateMaxDrawdown 计算最大回撤和持续天数
func (s *sRisk) calculateMaxDrawdown(returns []float64) (float64, int) {
	if len(returns) == 0 {
		return 0, 0
	}

	// 计算累计净值曲线
	cumValues := make([]float64, len(returns)+1)
	cumValues[0] = 1.0
	for i, r := range returns {
		cumValues[i+1] = cumValues[i] * (1 + r)
	}

	// 计算最大回撤
	maxDrawdown := 0.0
	maxDrawdownDays := 0
	peak := cumValues[0]
	peakIndex := 0

	for i := 1; i < len(cumValues); i++ {
		if cumValues[i] > peak {
			peak = cumValues[i]
			peakIndex = i
		} else {
			drawdown := (peak - cumValues[i]) / peak
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
				maxDrawdownDays = i - peakIndex
			}
		}
	}

	// 转为百分比
	maxDrawdownPct := maxDrawdown * 100
	return maxDrawdownPct, maxDrawdownDays
}

// calculateBeta 计算 Beta 系数（相对基准）
func (s *sRisk) calculateBeta(returns, benchmarkReturns []float64) float64 {
	if len(returns) < 2 || len(benchmarkReturns) < 2 {
		return 0
	}

	// 对齐长度
	minLen := len(returns)
	if len(benchmarkReturns) < minLen {
		minLen = len(benchmarkReturns)
	}
	returns = returns[:minLen]
	benchmarkReturns = benchmarkReturns[:minLen]

	// 计算协方差和基准方差
	covariance := stat.Covariance(returns, benchmarkReturns, nil)
	benchmarkVariance := stat.Variance(benchmarkReturns, nil)

	if benchmarkVariance == 0 {
		return 0
	}

	// Beta = Cov(R, Rm) / Var(Rm)
	beta := covariance / benchmarkVariance
	return beta
}
