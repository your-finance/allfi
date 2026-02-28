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
//
// 从资产快照数据计算每日收益率
//
// 算法步骤：
// 1. 遍历快照序列（按时间升序）
// 2. 对于第 i 天，计算收益率 = (第 i 天价值 - 第 i-1 天价值) / 第 i-1 天价值
// 3. 跳过起始价值为 0 的情况（避免除零错误）
//
// 参数：
//   - snapshots: 资产快照序列（按时间升序排列）
//
// 返回：日收益率序列（如 [0.02, -0.01, 0.03] 表示涨 2%、跌 1%、涨 3%）
//
// 示例：
//   - 快照价值：[1000, 1020, 1010, 1040]
//   - 收益率：[(1020-1000)/1000, (1010-1020)/1020, (1040-1010)/1010]
//   - 结果：[0.02, -0.0098, 0.0297]
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
//
// 通过 CoinGecko API 获取 BTC 历史价格，并计算日收益率序列
//
// 算法步骤：
// 1. 调用 CoinGecko API 获取过去 N 天的 BTC/USD 价格数据
//    - API 端点：/coins/bitcoin/market_chart?vs_currency=usd&days=N
//    - 返回格式：{"prices": [[timestamp, price], ...]}
// 2. 遍历价格序列，计算每日收益率
//    - 收益率 = (第 i 天价格 - 第 i-1 天价格) / 第 i-1 天价格
// 3. 返回收益率序列，用于计算 Beta 系数
//
// 参数：
//   - ctx: 上下文（用于请求取消和超时控制）
//   - days: 获取的天数（如 30 表示获取过去 30 天的数据）
//
// 返回：
//   - []float64: BTC 日收益率序列
//   - error: 错误信息（API 请求失败、数据不足等）
//
// 数据来源：CoinGecko 免费 API（无需 API Key）
//
// 示例：
//   - 价格序列：[50000, 51000, 49500]
//   - 收益率：[(51000-50000)/50000, (49500-51000)/51000]
//   - 结果：[0.02, -0.0294]
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
//
// 这是风险管理模块的核心方法，整合了所有风险指标的计算
//
// 计算流程：
// 1. VaR (Value at Risk) - 历史模拟法
//    - 95% 置信度 VaR：有 95% 的把握，损失不会超过该值
//    - 99% 置信度 VaR：有 99% 的把握，损失不会超过该值
//
// 2. 波动率（年化）
//    - 衡量资产价格的波动程度
//    - 波动率越高，风险越大
//
// 3. 夏普比率
//    - 衡量每承担一单位风险所获得的超额回报
//    - 夏普比率越高，投资表现越好
//
// 4. 下行偏差
//    - 只考虑负收益的波动性
//    - 比标准差更能反映投资者关心的"下行风险"
//
// 5. 索提诺比率
//    - 夏普比率的改进版本，只考虑下行风险
//    - 更符合投资者"只关心损失风险"的心理
//
// 6. 最大回撤
//    - 从历史最高点到最低点的最大跌幅
//    - 衡量投资组合可能遭受的最大损失
//
// 7. Beta 系数（相对 BTC）
//    - 衡量投资组合相对于 BTC 的系统性风险
//    - Beta > 1：波动大于 BTC（高风险高收益）
//    - Beta < 1：波动小于 BTC（低风险低收益）
//
// 参数：
//   - returns: 资产日收益率序列
//   - btcReturns: BTC 日收益率序列（用于计算 Beta）
//
// 返回：CalculationResult 结构体，包含所有风险指标
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
//
// VaR (Value at Risk) 表示在给定置信度下，投资组合在未来一段时间内可能遭受的最大损失
//
// 算法步骤：
// 1. 获取历史收益率序列（从 asset_snapshots 表计算得出）
// 2. 将收益率从小到大排序（负收益在前，正收益在后）
// 3. 取第 (1-α) 分位数作为 VaR 值（α 为置信度）
//    - 95% 置信度：取第 5% 分位数（最差的 5% 情况）
//    - 99% 置信度：取第 1% 分位数（最差的 1% 情况）
// 4. VaR = 当前资产价值 × 分位数收益率的绝对值
//
// 参数：
//   - returns: 日收益率序列（如 [-0.05, 0.02, -0.01, 0.03]）
//   - confidence: 置信度（0.95 表示 95%，0.99 表示 99%）
//
// 返回：VaR 百分比（如 5.2 表示在该置信度下最大损失 5.2%）
//
// 示例：
//   - 收益率序列：[-0.08, -0.05, -0.02, 0.01, 0.03, 0.05]
//   - 95% 置信度：取第 5% 分位数 = -0.08
//   - VaR = 8.0%（表示有 95% 的把握，损失不会超过 8%）
func (s *sRisk) calculateVaR(returns []float64, confidence float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 复制并排序收益率（从小到大：负收益 → 正收益）
	sorted := make([]float64, len(returns))
	copy(sorted, returns)
	sort.Float64s(sorted)

	// 计算分位数索引（1-置信度）
	// 例如：95% 置信度 → 取第 5% 分位数
	index := int((1 - confidence) * float64(len(sorted)))
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	// VaR 为负收益率的绝对值（转为百分比）
	// 例如：sorted[index] = -0.08 → VaR = 8.0%
	varValue := -sorted[index] * 100
	return varValue
}

// calculateVolatility 计算波动率（年化）
//
// 波动率衡量资产价格的波动程度，是风险的重要指标
//
// 算法步骤：
// 1. 计算日收益率的标准差（使用 gonum/stat.StdDev）
// 2. 年化波动率 = 日波动率 × sqrt(365)
//    - 根据金融学理论，波动率的年化需要乘以 sqrt(交易日数)
//    - 加密货币市场 7×24 小时交易，使用 365 天
// 3. 转为百分比（× 100）
//
// 参数：
//   - returns: 日收益率序列
//
// 返回：年化波动率百分比（如 25.5 表示年化波动率 25.5%）
//
// 示例：
//   - 日收益率标准差 = 0.02
//   - 年化波动率 = 0.02 × sqrt(365) × 100 ≈ 38.2%
func (s *sRisk) calculateVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}

	// 计算日收益率的标准差
	stdDev := stat.StdDev(returns, nil)

	// 年化波动率 = 日波动率 × sqrt(365) × 100
	annualizedVol := stdDev * math.Sqrt(365) * 100
	return annualizedVol
}

// calculateSharpeRatio 计算夏普比率
//
// 夏普比率衡量每承担一单位风险所获得的超额回报，是评估投资组合表现的重要指标
//
// 算法步骤：
// 1. 计算日收益率的平均值（使用 gonum/stat.Mean）
// 2. 年化平均收益率 = 日均收益率 × 365 × 100
// 3. 夏普比率 = (年化收益率 - 无风险利率) / 年化波动率
//    - 本项目假设无风险利率为 0（加密货币市场无传统意义的无风险资产）
//
// 参数：
//   - returns: 日收益率序列
//   - volatility: 年化波动率（百分比）
//
// 返回：夏普比率（无单位，如 1.5 表示每承担 1% 风险获得 1.5% 超额收益）
//
// 解读：
//   - Sharpe > 1.0：表现良好
//   - Sharpe > 2.0：表现优秀
//   - Sharpe < 0：收益为负，投资失败
//
// 示例：
//   - 年化收益率 = 30%，年化波动率 = 20%
//   - 夏普比率 = 30 / 20 = 1.5
func (s *sRisk) calculateSharpeRatio(returns []float64, volatility float64) float64 {
	if len(returns) == 0 || volatility == 0 {
		return 0
	}

	// 计算日收益率的平均值
	meanReturn := stat.Mean(returns, nil)

	// 年化平均收益率（日均收益率 × 365 × 100）
	annualizedReturn := meanReturn * 365 * 100

	// 夏普比率 = (年化收益率 - 无风险利率) / 年化波动率
	// 无风险利率假设为 0
	sharpe := annualizedReturn / volatility
	return sharpe
}

// calculateDownsideDeviation 计算下行偏差
//
// 下行偏差只考虑负收益的波动性，比标准差更能反映投资者关心的"下行风险"
//
// 算法步骤：
// 1. 从收益率序列中筛选出所有负收益率（r < 0）
// 2. 计算负收益率的标准差（使用 gonum/stat.StdDev）
// 3. 年化下行偏差 = 日下行偏差 × sqrt(365) × 100
//
// 参数：
//   - returns: 日收益率序列
//
// 返回：年化下行偏差百分比（如 15.2 表示下行偏差 15.2%）
//
// 与波动率的区别：
//   - 波动率：考虑所有收益率的波动（包括正收益）
//   - 下行偏差：只考虑负收益的波动（更关注损失风险）
//
// 示例：
//   - 收益率序列：[-0.05, -0.02, 0.03, 0.05]
//   - 负收益率：[-0.05, -0.02]
//   - 下行偏差 = StdDev([-0.05, -0.02]) × sqrt(365) × 100
func (s *sRisk) calculateDownsideDeviation(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 筛选负收益率（只考虑损失的情况）
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

	// 年化下行偏差（日下行偏差 × sqrt(365) × 100）
	annualizedDD := stdDev * math.Sqrt(365) * 100
	return annualizedDD
}

// calculateSortinoRatio 计算索提诺比率
//
// 索提诺比率是夏普比率的改进版本，只考虑下行风险（负收益的波动性）
//
// 算法步骤：
// 1. 计算日收益率的平均值
// 2. 年化平均收益率 = 日均收益率 × 365 × 100
// 3. 索提诺比率 = (年化收益率 - 目标收益率) / 下行偏差
//    - 目标收益率假设为 0（即不接受任何损失）
//
// 参数：
//   - returns: 日收益率序列
//   - downsideDeviation: 下行偏差（百分比）
//
// 返回：索提诺比率（无单位）
//
// 与夏普比率的区别：
//   - 夏普比率：分母为总波动率（包括正负收益的波动）
//   - 索提诺比率：分母为下行偏差（只考虑负收益的波动）
//   - 索提诺比率更符合投资者"只关心损失风险"的心理
//
// 解读：
//   - Sortino > 1.5：表现良好
//   - Sortino > 2.5：表现优秀
//   - Sortino < 0：收益为负，投资失败
//
// 示例：
//   - 年化收益率 = 30%，下行偏差 = 15%
//   - 索提诺比率 = 30 / 15 = 2.0
func (s *sRisk) calculateSortinoRatio(returns []float64, downsideDeviation float64) float64 {
	if len(returns) == 0 || downsideDeviation == 0 {
		return 0
	}

	// 计算日收益率的平均值
	meanReturn := stat.Mean(returns, nil)

	// 年化平均收益率（日均收益率 × 365 × 100）
	annualizedReturn := meanReturn * 365 * 100

	// 索提诺比率 = (年化收益率 - 目标收益率) / 下行偏差
	// 目标收益率假设为 0
	sortino := annualizedReturn / downsideDeviation
	return sortino
}

// calculateMaxDrawdown 计算最大回撤和持续天数
//
// 最大回撤衡量从历史最高点到最低点的最大跌幅，是评估投资风险的关键指标
//
// 算法步骤：
// 1. 根据日收益率序列计算累计净值曲线
//    - 初始净值 = 1.0
//    - 第 i 天净值 = 前一天净值 × (1 + 第 i 天收益率)
// 2. 遍历净值曲线，维护历史最高点（peak）
// 3. 对于每个点，计算从最高点到当前点的回撤：
//    - 回撤 = (peak - current) / peak
// 4. 记录最大回撤及其持续天数（从最高点到最低点的天数）
//
// 参数：
//   - returns: 日收益率序列
//
// 返回：
//   - maxDrawdown: 最大回撤百分比（如 25.5 表示最大回撤 25.5%）
//   - maxDrawdownDays: 最大回撤持续天数
//
// 示例：
//   - 净值曲线：[1.0, 1.1, 1.2, 0.9, 0.8, 1.0]
//   - 最高点：1.2（第 2 天）
//   - 最低点：0.8（第 4 天）
//   - 最大回撤 = (1.2 - 0.8) / 1.2 = 33.3%
//   - 持续天数 = 4 - 2 = 2 天
//
// 解读：
//   - 回撤 < 10%：风险较低
//   - 回撤 10-20%：中等风险
//   - 回撤 > 20%：高风险
func (s *sRisk) calculateMaxDrawdown(returns []float64) (float64, int) {
	if len(returns) == 0 {
		return 0, 0
	}

	// 计算累计净值曲线（初始净值 = 1.0）
	cumValues := make([]float64, len(returns)+1)
	cumValues[0] = 1.0
	for i, r := range returns {
		cumValues[i+1] = cumValues[i] * (1 + r)
	}

	// 计算最大回撤
	maxDrawdown := 0.0
	maxDrawdownDays := 0
	peak := cumValues[0]       // 历史最高点
	peakIndex := 0             // 最高点索引

	for i := 1; i < len(cumValues); i++ {
		// 更新历史最高点
		if cumValues[i] > peak {
			peak = cumValues[i]
			peakIndex = i
		} else {
			// 计算当前回撤
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
//
// Beta 系数衡量投资组合相对于基准指数的系统性风险（市场风险）
//
// 算法步骤：
// 1. 对齐资产收益率和基准收益率的长度（取较短的一方）
// 2. 计算资产收益率与基准收益率的协方差（使用 gonum/stat.Covariance）
// 3. 计算基准收益率的方差（使用 gonum/stat.Variance）
// 4. Beta = Cov(资产收益率, 基准收益率) / Var(基准收益率)
//
// 参数：
//   - returns: 资产日收益率序列
//   - benchmarkReturns: 基准日收益率序列（本项目使用 BTC 作为基准）
//
// 返回：Beta 系数（无单位）
//
// 解读：
//   - Beta = 1.0：与基准同步波动（系统性风险等于市场平均）
//   - Beta > 1.0：波动大于基准（高风险高收益）
//   - Beta < 1.0：波动小于基准（低风险低收益）
//   - Beta < 0：与基准反向波动（对冲资产）
//
// 示例：
//   - Beta = 1.5：当 BTC 涨 10% 时，资产预期涨 15%
//   - Beta = 0.5：当 BTC 涨 10% 时，资产预期涨 5%
//
// 数据来源：
//   - 资产收益率：从 asset_snapshots 表计算
//   - 基准收益率：从 CoinGecko API 获取 BTC 历史价格
func (s *sRisk) calculateBeta(returns, benchmarkReturns []float64) float64 {
	if len(returns) < 2 || len(benchmarkReturns) < 2 {
		return 0
	}

	// 对齐长度（取较短的一方）
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
