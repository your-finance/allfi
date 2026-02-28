// Package logic Gas 优化业务逻辑实现
package logic

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "your-finance/allfi/api/v1/gas"
	"your-finance/allfi/internal/app/gas/dao"
	"your-finance/allfi/internal/app/gas/model/entity"
	"your-finance/allfi/internal/app/gas/service"
	marketService "your-finance/allfi/internal/app/market/service"
)

// sGas Gas 优化服务实现
type sGas struct{}

// New 创建 Gas 优化服务实例
func New() service.IGas {
	return &sGas{}
}

// GetCurrent 获取当前 Gas 价格
func (s *sGas) GetCurrent(ctx context.Context) (*v1.GetCurrentRes, error) {
	// 调用 market 服务获取当前 Gas 价格
	gasData, err := marketService.Market().GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	chains := make([]v1.ChainGasPrice, len(gasData.Chains))
	for i, chain := range gasData.Chains {
		chains[i] = v1.ChainGasPrice{
			Chain:    chain.Chain,
			Name:     chain.Name,
			Low:      chain.Low,
			Standard: chain.Standard,
			Fast:     chain.Fast,
			Instant:  chain.Instant,
			BaseFee:  chain.BaseFee,
			Unit:     chain.Unit,
			Level:    chain.Level,
		}
	}

	return &v1.GetCurrentRes{
		Chains: chains,
	}, nil
}

// GetHistory 获取 Gas 价格历史
func (s *sGas) GetHistory(ctx context.Context, req *v1.GetHistoryReq) (*v1.GetHistoryRes, error) {
	// 计算查询时间范围
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(req.Hours) * time.Hour)

	// 查询历史记录
	var records []entity.GasPriceHistory
	err := dao.GasPriceHistory.Ctx(ctx).
		Where(dao.GasPriceHistory.Columns().Chain, req.Chain).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" >= ?", startTime).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" <= ?", endTime).
		Order(dao.GasPriceHistory.Columns().RecordedAt + " ASC").
		Scan(&records)

	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	history := make([]v1.GasPriceHistoryVO, len(records))
	for i, record := range records {
		history[i] = v1.GasPriceHistoryVO{
			Timestamp:  record.RecordedAt.Unix(),
			Low:        float64(record.Low),
			Standard:   float64(record.Standard),
			Fast:       float64(record.Fast),
			Instant:    float64(record.Instant),
			BaseFee:    float64(record.BaseFee),
			RecordedAt: record.RecordedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &v1.GetHistoryRes{
		Chain:   req.Chain,
		Hours:   req.Hours,
		History: history,
	}, nil
}

// GetRecommendation 获取最佳交易时间推荐
func (s *sGas) GetRecommendation(ctx context.Context, req *v1.GetRecommendationReq) (*v1.GetRecommendationRes, error) {
	// 查询最近 24 小时的历史数据
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	var records []entity.GasPriceHistory
	err := dao.GasPriceHistory.Ctx(ctx).
		Where(dao.GasPriceHistory.Columns().Chain, req.Chain).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" >= ?", startTime).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" <= ?", endTime).
		Order(dao.GasPriceHistory.Columns().RecordedAt + " ASC").
		Scan(&records)

	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("暂无历史数据，无法生成推荐")
	}

	// 分析每小时的平均 Gas 价格
	hourlyAvg := make(map[int][]float64) // hour -> prices
	for _, record := range records {
		hour := record.RecordedAt.Hour()
		hourlyAvg[hour] = append(hourlyAvg[hour], float64(record.Standard))
	}

	// 计算每小时的平均值
	hourlyMean := make(map[int]float64)
	for hour, prices := range hourlyAvg {
		sum := 0.0
		for _, price := range prices {
			sum += price
		}
		hourlyMean[hour] = sum / float64(len(prices))
	}

	// 找出 Gas 价格最低的时段
	minHour := -1
	minPrice := math.MaxFloat64
	for hour, price := range hourlyMean {
		if price < minPrice {
			minPrice = price
			minHour = hour
		}
	}

	// 获取当前 Gas 价格
	gasData, err := marketService.Market().GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	var currentPrice float64
	for _, chain := range gasData.Chains {
		if chain.Chain == req.Chain {
			currentPrice = chain.Standard
			break
		}
	}

	// 计算预计节省百分比
	estimatedSavings := 0.0
	if currentPrice > 0 {
		estimatedSavings = ((currentPrice - minPrice) / currentPrice) * 100
	}

	// 计算置信度（基于数据点数量）
	confidence := math.Min(float64(len(records))/288.0, 1.0) // 288 = 24小时 * 12次/小时

	// 生成推荐时间段描述
	recommendedTime := fmt.Sprintf("%02d:00 - %02d:00 UTC", minHour, (minHour+1)%24)

	// 有效期设置为 1 小时后
	validUntil := time.Now().Add(1 * time.Hour)

	return &v1.GetRecommendationRes{
		Chain:            req.Chain,
		RecommendedTime:  recommendedTime,
		EstimatedSavings: estimatedSavings,
		Confidence:       confidence,
		ValidUntil:       validUntil.Format("2006-01-02 15:04:05"),
		CurrentGasPrice:  currentPrice,
		RecommendedPrice: minPrice,
	}, nil
}

// GetForecast 获取 Gas 价格预测
func (s *sGas) GetForecast(ctx context.Context, req *v1.GetForecastReq) (*v1.GetForecastRes, error) {
	// 查询最近 7 天的历史数据用于预测
	endTime := time.Now()
	startTime := endTime.Add(-7 * 24 * time.Hour)

	var records []entity.GasPriceHistory
	err := dao.GasPriceHistory.Ctx(ctx).
		Where(dao.GasPriceHistory.Columns().Chain, req.Chain).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" >= ?", startTime).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" <= ?", endTime).
		Order(dao.GasPriceHistory.Columns().RecordedAt + " ASC").
		Scan(&records)

	if err != nil {
		return nil, err
	}

	if len(records) < 10 {
		return nil, fmt.Errorf("历史数据不足，无法生成预测（至少需要 10 条记录）")
	}

	// 使用简单移动平均法进行预测
	windowSize := 12 // 使用最近 12 个数据点（约 1 小时）
	forecast := make([]v1.GasForecastVO, req.Hours)

	// 计算最近的平均值和标准差
	recentPrices := make([]float64, 0)
	for i := len(records) - windowSize; i < len(records); i++ {
		if i >= 0 {
			recentPrices = append(recentPrices, float64(records[i].Standard))
		}
	}

	mean := calculateMean(recentPrices)
	stdDev := calculateStdDev(recentPrices, mean)

	// 计算趋势（线性回归斜率）
	slope := calculateSlope(recentPrices)
	var trend string
	if slope > 0.5 {
		trend = "上升"
	} else if slope < -0.5 {
		trend = "下降"
	} else {
		trend = "平稳"
	}

	// 生成预测数据
	now := time.Now()
	for i := 0; i < req.Hours; i++ {
		futureTime := now.Add(time.Duration(i+1) * time.Hour)

		// 简单预测：基于趋势和历史平均
		predicted := mean + slope*float64(i+1)

		// 预测区间（±1 标准差）
		lower := predicted - stdDev
		upper := predicted + stdDev

		// 确保价格不为负
		if lower < 0 {
			lower = 0
		}
		if predicted < 0 {
			predicted = 0
		}

		forecast[i] = v1.GasForecastVO{
			Timestamp: futureTime.Unix(),
			Time:      futureTime.Format("2006-01-02 15:04:05"),
			Predicted: predicted,
			Lower:     lower,
			Upper:     upper,
		}
	}

	// 计算预测置信度（基于数据质量）
	confidence := math.Min(float64(len(records))/1000.0, 0.85) // 最高 85%

	return &v1.GetForecastRes{
		Chain:      req.Chain,
		Hours:      req.Hours,
		Forecast:   forecast,
		Trend:      trend,
		Confidence: confidence,
	}, nil
}

// RecordGasPrice 记录 Gas 价格（定时任务调用）
func (s *sGas) RecordGasPrice(ctx context.Context) error {
	// 获取当前 Gas 价格
	gasData, err := marketService.Market().GetGasPrice(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取 Gas 价格失败", err)
		return err
	}

	// 记录每条链的 Gas 价格
	now := gtime.Now()
	for _, chain := range gasData.Chains {
		record := g.Map{
			dao.GasPriceHistory.Columns().Chain:      chain.Chain,
			dao.GasPriceHistory.Columns().Low:        chain.Low,
			dao.GasPriceHistory.Columns().Standard:   chain.Standard,
			dao.GasPriceHistory.Columns().Fast:       chain.Fast,
			dao.GasPriceHistory.Columns().Instant:    chain.Instant,
			dao.GasPriceHistory.Columns().BaseFee:    chain.BaseFee,
			dao.GasPriceHistory.Columns().RecordedAt: now,
		}

		_, err := dao.GasPriceHistory.Ctx(ctx).Data(record).Insert()
		if err != nil {
			g.Log().Error(ctx, "记录 Gas 价格失败", "chain", chain.Chain, "error", err)
			continue
		}
	}

	// 清理 7 天前的历史数据
	cleanupTime := now.Add(-7 * 24 * time.Hour)
	_, err = dao.GasPriceHistory.Ctx(ctx).
		Where(dao.GasPriceHistory.Columns().RecordedAt+" < ?", cleanupTime).
		Delete()
	if err != nil {
		g.Log().Warning(ctx, "清理历史数据失败", err)
	}

	g.Log().Info(ctx, "Gas 价格记录成功", "chains", len(gasData.Chains))
	return nil
}

// calculateMean 计算平均值
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateStdDev 计算标准差
func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	variance := 0.0
	for _, v := range values {
		variance += math.Pow(v-mean, 2)
	}
	return math.Sqrt(variance / float64(len(values)))
}

// calculateSlope 计算线性回归斜率
func calculateSlope(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}

	// 计算 x 和 y 的平均值
	sumX := 0.0
	sumY := 0.0
	for i, y := range values {
		sumX += float64(i)
		sumY += y
	}
	meanX := sumX / float64(n)
	meanY := sumY / float64(n)

	// 计算斜率
	numerator := 0.0
	denominator := 0.0
	for i, y := range values {
		x := float64(i)
		numerator += (x - meanX) * (y - meanY)
		denominator += math.Pow(x-meanX, 2)
	}

	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
