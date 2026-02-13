// Package logic 策略引擎业务逻辑
// 提供策略 CRUD + 配比分析 + 调仓建议
package logic

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	strategyApi "your-finance/allfi/api/v1/strategy"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/strategy/dao"
	strategyModel "your-finance/allfi/internal/app/strategy/model"
	"your-finance/allfi/internal/app/strategy/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/model/entity"
)

// sStrategy 策略服务实现
type sStrategy struct{}

// New 创建策略服务实例
func New() service.IStrategy {
	return &sStrategy{}
}

// List 获取策略列表
func (s *sStrategy) List(ctx context.Context) (*strategyApi.ListRes, error) {
	var strategies []entity.Strategies
	err := dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().UserId, consts.GetUserID(ctx)).
		OrderDesc(dao.Strategies.Columns().CreatedAt).
		Scan(&strategies)
	if err != nil {
		return nil, gerror.Wrap(err, "查询策略列表失败")
	}

	items := make([]strategyApi.StrategyItem, 0, len(strategies))
	for _, st := range strategies {
		items = append(items, toStrategyItem(&st))
	}
	return &strategyApi.ListRes{Strategies: items}, nil
}

// Create 创建策略
func (s *sStrategy) Create(ctx context.Context, name, sType string, config any) (*strategyApi.CreateRes, error) {
	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, gerror.Wrap(err, "序列化策略配置失败")
	}

	now := gtime.Now().Time
	st := &entity.Strategies{
		UserId:    consts.GetUserID(ctx),
		Name:      name,
		Type:      sType,
		Config:    string(configJSON),
		IsActive:  1, // SQLite: 1 = true
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := dao.Strategies.Ctx(ctx).Insert(st)
	if err != nil {
		return nil, gerror.Wrap(err, "创建策略失败")
	}
	id, _ := result.LastInsertId()
	st.Id = int(id)

	item := toStrategyItem(st)
	return &strategyApi.CreateRes{Strategy: &item}, nil
}

// Update 更新策略
func (s *sStrategy) Update(ctx context.Context, id uint, name, sType string, config any, isActive *bool) (*strategyApi.UpdateRes, error) {
	// 检查策略是否存在
	var existing entity.Strategies
	err := dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().Id, id).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询策略失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("策略不存在: %d", id)
	}

	// 构建更新数据
	data := g.Map{
		dao.Strategies.Columns().UpdatedAt: gtime.Now().Time,
	}
	if name != "" {
		data[dao.Strategies.Columns().Name] = name
	}
	if sType != "" {
		data[dao.Strategies.Columns().Type] = sType
	}
	if config != nil {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, gerror.Wrap(err, "序列化策略配置失败")
		}
		data[dao.Strategies.Columns().Config] = string(configJSON)
	}
	if isActive != nil {
		if *isActive {
			data[dao.Strategies.Columns().IsActive] = 1
		} else {
			data[dao.Strategies.Columns().IsActive] = 0
		}
	}

	_, err = dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().Id, id).
		Data(data).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新策略失败")
	}

	// 重新查询返回
	var updated entity.Strategies
	_ = dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().Id, id).
		Scan(&updated)
	item := toStrategyItem(&updated)
	return &strategyApi.UpdateRes{Strategy: &item}, nil
}

// Delete 删除策略（软删除）
func (s *sStrategy) Delete(ctx context.Context, id uint) error {
	_, err := dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().Id, id).
		Data(g.Map{
			dao.Strategies.Columns().DeletedAt: gtime.Now().Time,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "删除策略失败")
	}
	return nil
}

// GetAnalysis 获取策略分析（偏离度 + 调仓建议）
func (s *sStrategy) GetAnalysis(ctx context.Context, id uint) (*strategyApi.GetAnalysisRes, error) {
	// 获取策略
	var st entity.Strategies
	err := dao.Strategies.Ctx(ctx).
		Where(dao.Strategies.Columns().Id, id).
		Scan(&st)
	if err != nil || st.Id == 0 {
		return nil, gerror.Newf("策略不存在: %d", id)
	}

	// 解析策略配置
	var config strategyModel.RebalanceConfig
	if err := json.Unmarshal([]byte(st.Config), &config); err != nil {
		return nil, gerror.Wrap(err, "解析策略配置失败")
	}

	// 获取当前资产明细
	var details []entity.AssetDetails
	err = assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
		Scan(&details)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产明细失败")
	}

	// 计算当前各币种占比
	totalValue := 0.0
	assetValues := make(map[string]float64)
	for _, d := range details {
		assetValues[d.AssetSymbol] += d.ValueUsd
		totalValue += d.ValueUsd
	}

	currentAlloc := make(map[string]float64)
	if totalValue > 0 {
		for symbol, value := range assetValues {
			currentAlloc[symbol] = (value / totalValue) * 100
		}
	}

	// 目标配比
	targetAlloc := make(map[string]float64)
	for _, alloc := range config.Allocations {
		targetAlloc[alloc.Symbol] = alloc.Percentage
	}

	// 计算偏离度和调仓建议
	deviation := make(map[string]float64)
	recommendations := make([]strategyApi.RecommendationItem, 0)

	// 收集需要查询价格的币种
	needPriceSymbols := make([]string, 0)
	for symbol, target := range targetAlloc {
		current := currentAlloc[symbol]
		dev := current - target
		threshold := config.Threshold
		if threshold <= 0 {
			threshold = 5
		}
		if math.Abs(dev) > threshold {
			needPriceSymbols = append(needPriceSymbols, symbol)
		}
	}

	// 通过 CoinGecko 批量获取当前价格
	prices := make(map[string]float64)
	if len(needPriceSymbols) > 0 {
		cgClient := coingecko.NewClient("")
		fetched, err := cgClient.GetPrices(ctx, needPriceSymbols)
		if err != nil {
			g.Log().Warning(ctx, "获取 CoinGecko 价格失败，调仓建议的 Amount 将为 0", "error", err)
		} else {
			prices = fetched
		}
	}

	for symbol, target := range targetAlloc {
		current := currentAlloc[symbol]
		dev := current - target
		deviation[symbol] = dev

		// 生成建议（偏离超过阈值）
		threshold := config.Threshold
		if threshold <= 0 {
			threshold = 5 // 默认 5%
		}

		if math.Abs(dev) > threshold {
			action := "buy"
			reason := "低于目标配比"
			if dev > 0 {
				action = "sell"
				reason = "高于目标配比"
			}
			adjustValue := math.Abs(dev) / 100 * totalValue

			// 通过价格计算调仓数量：Amount = ValueUSD / Price
			var amount float64
			upperSymbol := strings.ToUpper(symbol)
			if price, ok := prices[upperSymbol]; ok && price > 0 {
				amount = adjustValue / price
			}

			recommendations = append(recommendations, strategyApi.RecommendationItem{
				Symbol:   symbol,
				Action:   action,
				Amount:   amount,
				ValueUSD: adjustValue,
				Reason:   reason,
			})
		}
	}

	return &strategyApi.GetAnalysisRes{
		Analysis: &strategyApi.AnalysisItem{
			StrategyID:      uint(st.Id),
			CurrentAlloc:    currentAlloc,
			TargetAlloc:     targetAlloc,
			Deviation:       deviation,
			Recommendations: recommendations,
		},
	}, nil
}

// toStrategyItem 将实体转换为 API 条目
func toStrategyItem(st *entity.Strategies) strategyApi.StrategyItem {
	// 解析配置 JSON
	var config any
	if st.Config != "" {
		var parsed map[string]any
		if json.Unmarshal([]byte(st.Config), &parsed) == nil {
			config = parsed
		} else {
			config = st.Config
		}
	}

	return strategyApi.StrategyItem{
		ID:        uint(st.Id),
		Name:      st.Name,
		Type:      st.Type,
		Config:    config,
		IsActive:  st.IsActive == 1, // SQLite: 1 = true
		CreatedAt: st.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: st.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
