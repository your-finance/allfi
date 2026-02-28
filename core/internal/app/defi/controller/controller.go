// Package controller DeFi 仓位模块 - 路由和控制器
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	defiApi "your-finance/allfi/api/v1/defi"
	"your-finance/allfi/internal/app/defi/service"
)

// DefiController DeFi 仓位控制器
type DefiController struct{}

// GetPositions 获取 DeFi 仓位列表
func (c *DefiController) GetPositions(ctx context.Context, req *defiApi.GetPositionsReq) (res *defiApi.GetPositionsRes, err error) {
	// 调用服务层获取 DeFi 仓位
	positions, totalValue, err := service.Defi().GetPositions(ctx, "", "")
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetPositionsRes{
		TotalValue: totalValue,
		Currency:   req.Currency,
		Positions:  make([]defiApi.PositionItem, 0, len(positions)),
	}

	for _, p := range positions {
		res.Positions = append(res.Positions, defiApi.PositionItem{
			Protocol:   p.Protocol,
			Type:       p.Type,
			Token:      p.Token,
			Amount:     p.Amount,
			Value:      p.Value,
			APY:        p.APY,
			Chain:      p.Chain,
			WalletAddr: p.WalletAddr,
		})
	}

	return res, nil
}

// GetProtocols 获取支持的 DeFi 协议列表
func (c *DefiController) GetProtocols(ctx context.Context, req *defiApi.GetProtocolsReq) (res *defiApi.GetProtocolsRes, err error) {
	// 调用服务层获取协议列表
	protocols, err := service.Defi().GetProtocols(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetProtocolsRes{
		Protocols: make([]defiApi.ProtocolItem, 0, len(protocols)),
	}

	for _, p := range protocols {
		res.Protocols = append(res.Protocols, defiApi.ProtocolItem{
			Name:     p.Name,
			Chains:   p.Chains,
			Types:    p.Types,
			IsActive: p.IsActive,
		})
	}

	return res, nil
}

// GetStats 获取 DeFi 统计
func (c *DefiController) GetStats(ctx context.Context, req *defiApi.GetStatsReq) (res *defiApi.GetStatsRes, err error) {
	return service.Defi().GetStats(ctx)
}

// GetLendingPositions 获取借贷仓位列表
func (c *DefiController) GetLendingPositions(ctx context.Context, req *defiApi.GetLendingPositionsReq) (res *defiApi.GetLendingPositionsRes, err error) {
	// 调用服务层获取借贷仓位
	positions, err := service.Defi().GetLendingPositions(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetLendingPositionsRes{
		Positions: make([]defiApi.LendingPositionItem, 0, len(positions)),
	}

	for _, p := range positions {
		res.Positions = append(res.Positions, defiApi.LendingPositionItem{
			Protocol:             p.Protocol,
			Chain:                p.Chain,
			WalletAddr:           p.WalletAddr,
			SupplyToken:          p.SupplyToken,
			SupplyAmount:         p.SupplyAmount,
			SupplyValueUSD:       p.SupplyValueUSD,
			SupplyAPY:            p.SupplyAPY,
			BorrowToken:          p.BorrowToken,
			BorrowAmount:         p.BorrowAmount,
			BorrowValueUSD:       p.BorrowValueUSD,
			BorrowAPY:            p.BorrowAPY,
			HealthFactor:         p.HealthFactor,
			LiquidationThreshold: p.LiquidationThreshold,
			LTV:                  p.LTV,
			NetAPY:               p.NetAPY,
		})
	}

	return res, nil
}

// GetLendingRates 获取借贷利率
func (c *DefiController) GetLendingRates(ctx context.Context, req *defiApi.GetLendingRatesReq) (res *defiApi.GetLendingRatesRes, err error) {
	// 调用服务层获取借贷利率
	rates, err := service.Defi().GetLendingRates(ctx, req.Protocol, req.Chain)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetLendingRatesRes{
		Rates: make([]defiApi.LendingRateItem, 0, len(rates)),
	}

	for _, r := range rates {
		res.Rates = append(res.Rates, defiApi.LendingRateItem{
			Protocol:          r.Protocol,
			Chain:             r.Chain,
			Token:             r.Token,
			SupplyAPY:         r.SupplyAPY,
			BorrowAPYStable:   r.BorrowAPYStable,
			BorrowAPYVariable: r.BorrowAPYVariable,
			UtilizationRate:   r.UtilizationRate,
		})
	}

	return res, nil
}

// GetLendingRateHistory 获取借贷利率历史
func (c *DefiController) GetLendingRateHistory(ctx context.Context, req *defiApi.GetLendingRateHistoryReq) (res *defiApi.GetLendingRateHistoryRes, err error) {
	// 调用服务层获取利率历史
	history, err := service.Defi().GetLendingRateHistory(ctx, req.Protocol, req.Token, req.Days)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetLendingRateHistoryRes{
		History: make([]defiApi.LendingRateHistoryItem, 0, len(history)),
	}

	for _, h := range history {
		res.History = append(res.History, defiApi.LendingRateHistoryItem{
			Date:              h.Date,
			SupplyAPY:         h.SupplyAPY,
			BorrowAPYStable:   h.BorrowAPYStable,
			BorrowAPYVariable: h.BorrowAPYVariable,
			UtilizationRate:   h.UtilizationRate,
		})
	}

	return res, nil
}

// GetLendingOptimization 获取借贷策略优化建议
func (c *DefiController) GetLendingOptimization(ctx context.Context, req *defiApi.GetLendingOptimizationReq) (res *defiApi.GetLendingOptimizationRes, err error) {
	// 调用服务层获取优化建议
	result, err := service.Defi().GetLendingOptimization(ctx)
	if err != nil {
		return nil, err
	}

	// 转换为 API 响应格式
	res = &defiApi.GetLendingOptimizationRes{
		CurrentPositions: make([]defiApi.LendingPositionItem, 0, len(result.CurrentPositions)),
		Recommendations:  make([]defiApi.LendingRecommendation, 0, len(result.Recommendations)),
		PotentialGain:    result.PotentialGain,
		RiskLevel:        result.RiskLevel,
		Summary:          result.Summary,
	}

	// 转换当前仓位
	for _, p := range result.CurrentPositions {
		res.CurrentPositions = append(res.CurrentPositions, defiApi.LendingPositionItem{
			Protocol:             p.Protocol,
			Chain:                p.Chain,
			WalletAddr:           p.WalletAddr,
			SupplyToken:          p.SupplyToken,
			SupplyAmount:         p.SupplyAmount,
			SupplyValueUSD:       p.SupplyValueUSD,
			SupplyAPY:            p.SupplyAPY,
			BorrowToken:          p.BorrowToken,
			BorrowAmount:         p.BorrowAmount,
			BorrowValueUSD:       p.BorrowValueUSD,
			BorrowAPY:            p.BorrowAPY,
			HealthFactor:         p.HealthFactor,
			LiquidationThreshold: p.LiquidationThreshold,
			LTV:                  p.LTV,
			NetAPY:               p.NetAPY,
		})
	}

	// 转换优化建议
	for _, r := range result.Recommendations {
		res.Recommendations = append(res.Recommendations, defiApi.LendingRecommendation{
			Action:       r.Action,
			FromProtocol: r.FromProtocol,
			ToProtocol:   r.ToProtocol,
			Token:        r.Token,
			Amount:       r.Amount,
			Reason:       r.Reason,
			ExpectedGain: r.ExpectedGain,
		})
	}

	return res, nil
}


// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&DefiController{})
}
