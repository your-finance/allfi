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

// Register 注册路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&DefiController{})
}
