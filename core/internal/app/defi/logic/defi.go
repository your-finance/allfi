// Package logic DeFi 仓位模块 - 业务逻辑实现
package logic

import (
	"context"
	"math"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	defiApi "your-finance/allfi/api/v1/defi"
	"your-finance/allfi/internal/app/defi/model"
	"your-finance/allfi/internal/app/defi/service"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	"your-finance/allfi/internal/consts"
	defiPkg "your-finance/allfi/internal/integrations/defi"
	"your-finance/allfi/internal/model/entity"
)

// sDefi DeFi 仓位服务实现
type sDefi struct{}

// New 创建 DeFi 仓位服务实例
func New() service.IDefi {
	return &sDefi{}
}

// defiRegistry 全局 DeFi 协议注册中心
// 在 main.go 中初始化并设置
var defiRegistry *defiPkg.Registry

// SetRegistry 设置 DeFi 协议注册中心
// 由 main.go 在启动时调用
func SetRegistry(r *defiPkg.Registry) {
	defiRegistry = r
}

// GetPositions 获取用户 DeFi 仓位列表
//
// 业务逻辑:
// 1. 获取用户所有钱包地址
// 2. 并发查询每个钱包在各 DeFi 协议中的仓位
// 3. 聚合所有仓位并计算总价值
func (s *sDefi) GetPositions(ctx context.Context, chain string, protocol string) (positions []*model.PositionItem, totalValue float64, err error) {
	userID := consts.GetUserID(ctx)

	// 获取用户所有钱包地址
	var wallets []*entity.WalletAddresses
	err = walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, userID).
		Scan(&wallets)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询钱包地址失败")
	}

	if len(wallets) == 0 {
		g.Log().Info(ctx, "用户无钱包地址，跳过 DeFi 仓位查询")
		return []*model.PositionItem{}, 0, nil
	}

	// 检查注册中心是否可用
	if defiRegistry == nil {
		g.Log().Warning(ctx, "DeFi 协议注册中心未初始化")
		return []*model.PositionItem{}, 0, nil
	}

	// 并发查询每个钱包的 DeFi 仓位
	type result struct {
		positions []defiPkg.Position
		wallet    string
		err       error
	}

	ch := make(chan result, len(wallets))
	var wg sync.WaitGroup

	for _, w := range wallets {
		wg.Add(1)
		go func(wallet *entity.WalletAddresses) {
			defer wg.Done()

			var pos []defiPkg.Position
			var queryErr error

			if protocol != "" {
				// 指定协议查询
				pos, queryErr = defiRegistry.GetPositionsByProtocol(ctx, wallet.Address, chain, protocol)
			} else {
				// 查询所有协议
				pos, queryErr = defiRegistry.GetAllPositions(ctx, wallet.Address, chain)
			}

			ch <- result{positions: pos, wallet: wallet.Address, err: queryErr}
		}(w)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// 收集所有仓位
	positions = make([]*model.PositionItem, 0)
	for r := range ch {
		if r.err != nil {
			g.Log().Warning(ctx, "查询钱包 DeFi 仓位失败",
				"wallet", r.wallet,
				"error", r.err,
			)
			continue // 单个钱包失败不影响其他
		}

		for _, p := range r.positions {
			// 取第一个存入代币的符号和数量
			var tokenSymbol string
			var amount float64
			if len(p.DepositTokens) > 0 {
				tokenSymbol = p.DepositTokens[0].Symbol
				amount = p.DepositTokens[0].Amount
			}

			item := &model.PositionItem{
				Protocol:   p.ProtocolName,
				Type:       p.Type,
				Token:      tokenSymbol,
				Amount:     math.Round(amount*10000) / 10000,
				Value:      math.Round(p.ValueUSD*100) / 100,
				APY:        math.Round(p.APY*100) / 100,
				Chain:      p.Chain,
				WalletAddr: r.wallet,
			}
			positions = append(positions, item)
			totalValue += p.ValueUSD
		}
	}

	totalValue = math.Round(totalValue*100) / 100

	g.Log().Info(ctx, "获取 DeFi 仓位成功",
		"walletCount", len(wallets),
		"positionCount", len(positions),
		"totalValue", totalValue,
	)

	return positions, totalValue, nil
}

// GetStats 获取 DeFi 统计（按协议/链/类型分组聚合）
//
// 业务逻辑:
// 1. 调用 GetPositions 获取所有仓位
// 2. 遍历仓位列表，按 Protocol/Chain/Type 分组聚合价值
// 3. 返回统计结果
func (s *sDefi) GetStats(ctx context.Context) (*defiApi.GetStatsRes, error) {
	// 获取所有仓位
	positions, totalValue, err := s.GetPositions(ctx, "", "")
	if err != nil {
		return nil, err
	}

	// 按协议/链/类型分组聚合价值
	byProtocol := make(map[string]float64)
	byChain := make(map[string]float64)
	byType := make(map[string]float64)

	for _, p := range positions {
		byProtocol[p.Protocol] += p.Value
		byChain[p.Chain] += p.Value
		byType[p.Type] += p.Value
	}

	// 保留两位小数
	for k, v := range byProtocol {
		byProtocol[k] = math.Round(v*100) / 100
	}
	for k, v := range byChain {
		byChain[k] = math.Round(v*100) / 100
	}
	for k, v := range byType {
		byType[k] = math.Round(v*100) / 100
	}

	g.Log().Info(ctx, "获取 DeFi 统计成功",
		"positionCount", len(positions),
		"totalValueLocked", totalValue,
	)

	return &defiApi.GetStatsRes{
		TotalValueLocked: totalValue,
		PositionCount:    len(positions),
		ByProtocol:       byProtocol,
		ByChain:          byChain,
		ByType:           byType,
	}, nil
}

// GetProtocols 获取支持的 DeFi 协议列表
func (s *sDefi) GetProtocols(ctx context.Context) ([]*model.ProtocolItem, error) {
	if defiRegistry == nil {
		g.Log().Warning(ctx, "DeFi 协议注册中心未初始化")
		return []*model.ProtocolItem{}, nil
	}

	infos := defiRegistry.ListProtocols()
	protocols := make([]*model.ProtocolItem, 0, len(infos))

	for _, info := range infos {
		protocols = append(protocols, &model.ProtocolItem{
			Name:     info.DisplayName,
			Chains:   info.Chains,
			Types:    []string{info.Type},
			IsActive: true,
		})
	}

	g.Log().Info(ctx, "获取 DeFi 协议列表成功", "count", len(protocols))
	return protocols, nil
}
