// Package logic DeFi 借贷管理 - 业务逻辑实现
package logic

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"your-finance/allfi/internal/app/defi/dao"
	"your-finance/allfi/internal/app/defi/model"
	defiEntity "your-finance/allfi/internal/app/defi/model/entity"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	walletEntity "your-finance/allfi/internal/app/wallet/model/entity"
	"your-finance/allfi/internal/consts"
	defiPkg "your-finance/allfi/internal/integrations/defi"
)

// GetLendingPositions 获取用户的借贷仓位（包含健康因子）
func (s *sDefi) GetLendingPositions(ctx context.Context) ([]*model.LendingPositionItem, error) {
	userID := consts.GetUserID(ctx)

	// 获取用户所有钱包地址
	var wallets []*walletEntity.WalletAddresses
	err := walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, userID).
		Scan(&wallets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询钱包地址失败")
	}

	if len(wallets) == 0 {
		return []*model.LendingPositionItem{}, nil
	}

	// 检查注册中心是否可用
	if defiRegistry == nil {
		g.Log().Warning(ctx, "DeFi 协议注册中心未初始化")
		return []*model.LendingPositionItem{}, nil
	}

	// 并发查询每个钱包的借贷仓位
	type result struct {
		positions []*model.LendingPositionItem
		err       error
	}

	ch := make(chan result, len(wallets))
	var wg sync.WaitGroup

	for _, w := range wallets {
		wg.Add(1)
		go func(wallet *walletEntity.WalletAddresses) {
			defer wg.Done()

			positions, err := s.fetchLendingPositionsForWallet(ctx, wallet.Address, wallet.Blockchain)
			ch <- result{positions: positions, err: err}
		}(w)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// 收集所有仓位
	var allPositions []*model.LendingPositionItem
	for r := range ch {
		if r.err != nil {
			g.Log().Warning(ctx, "查询钱包借贷仓位失败", "error", r.err)
			continue
		}
		allPositions = append(allPositions, r.positions...)
	}

	// 保存到数据库
	if err := s.saveLendingPositions(ctx, userID, allPositions); err != nil {
		g.Log().Warning(ctx, "保存借贷仓位失败", "error", err)
	}

	return allPositions, nil
}

// fetchLendingPositionsForWallet 获取单个钱包的借贷仓位
func (s *sDefi) fetchLendingPositionsForWallet(ctx context.Context, address string, chain string) ([]*model.LendingPositionItem, error) {
	var positions []*model.LendingPositionItem

	// 获取所有借贷协议
	protocols := defiRegistry.ListProtocols()
	for _, info := range protocols {
		if info.Type != "lending" {
			continue
		}

		protocol, err := defiRegistry.GetProtocol(info.Name)
		if err != nil {
			continue
		}

		// 检查是否实现 LendingProtocol 接口
		lendingProtocol, ok := protocol.(defiPkg.LendingProtocol)
		if !ok {
			continue
		}

		// 获取借贷仓位
		defiPositions, err := lendingProtocol.GetLendingPositions(ctx, address, chain)
		if err != nil {
			g.Log().Warning(ctx, "获取借贷仓位失败",
				"protocol", info.Name,
				"address", address,
				"error", err,
			)
			continue
		}

		// 转换为业务模型
		for _, pos := range defiPositions {
			item := &model.LendingPositionItem{
				Protocol:             pos.ProtocolName,
				Chain:                pos.Chain,
				WalletAddr:           address,
				SupplyValueUSD:       pos.ValueUSD,
				SupplyAPY:            pos.APY,
				HealthFactor:         pos.HealthFactor,
				LiquidationThreshold: pos.LiquidationThreshold,
				LTV:                  pos.LTV,
				NetAPY:               pos.NetAPY,
			}

			// 存款代币
			if len(pos.DepositTokens) > 0 {
				item.SupplyToken = pos.DepositTokens[0].Symbol
				item.SupplyAmount = pos.DepositTokens[0].Amount
			}

			// 借款代币
			if len(pos.BorrowTokens) > 0 {
				item.BorrowToken = pos.BorrowTokens[0].Symbol
				item.BorrowAmount = pos.BorrowTokens[0].Amount
				item.BorrowValueUSD = pos.BorrowTokens[0].ValueUSD
			}

			positions = append(positions, item)
		}
	}

	return positions, nil
}

// saveLendingPositions 保存借贷仓位到数据库
func (s *sDefi) saveLendingPositions(ctx context.Context, userID int, positions []*model.LendingPositionItem) error {
	// 删除旧数据
	_, err := dao.LendingPositions.Ctx(ctx).
		Where(dao.LendingPositions.Columns().UserId, userID).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "删除旧借贷仓位失败")
	}

	// 插入新数据
	for _, pos := range positions {
		_, err := dao.LendingPositions.Ctx(ctx).Data(g.Map{
			dao.LendingPositions.Columns().UserId:               userID,
			dao.LendingPositions.Columns().Protocol:             pos.Protocol,
			dao.LendingPositions.Columns().Chain:                pos.Chain,
			dao.LendingPositions.Columns().WalletAddress:        pos.WalletAddr,
			dao.LendingPositions.Columns().SupplyToken:          pos.SupplyToken,
			dao.LendingPositions.Columns().SupplyAmount:         pos.SupplyAmount,
			dao.LendingPositions.Columns().SupplyValueUsd:       pos.SupplyValueUSD,
			dao.LendingPositions.Columns().SupplyApy:            pos.SupplyAPY,
			dao.LendingPositions.Columns().BorrowToken:          pos.BorrowToken,
			dao.LendingPositions.Columns().BorrowAmount:         pos.BorrowAmount,
			dao.LendingPositions.Columns().BorrowValueUsd:       pos.BorrowValueUSD,
			dao.LendingPositions.Columns().BorrowApy:            pos.BorrowAPY,
			dao.LendingPositions.Columns().HealthFactor:         pos.HealthFactor,
			dao.LendingPositions.Columns().LiquidationThreshold: pos.LiquidationThreshold,
			dao.LendingPositions.Columns().Ltv:                  pos.LTV,
			dao.LendingPositions.Columns().NetApy:               pos.NetAPY,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "插入借贷仓位失败")
		}
	}

	return nil
}

// GetLendingRates 获取借贷利率（存款APY、借款APY）
func (s *sDefi) GetLendingRates(ctx context.Context, protocol string, chain string) ([]*model.LendingRateItem, error) {
	if defiRegistry == nil {
		return []*model.LendingRateItem{}, nil
	}

	// 获取指定协议
	var protocols []defiPkg.DeFiProtocol
	if protocol != "" {
		p, err := defiRegistry.GetProtocol(protocol)
		if err == nil {
			protocols = append(protocols, p)
		}
	} else {
		// 获取所有借贷协议
		infos := defiRegistry.ListProtocols()
		for _, info := range infos {
			if info.Type == "lending" {
				p, err := defiRegistry.GetProtocol(info.Name)
				if err == nil {
					protocols = append(protocols, p)
				}
			}
		}
	}

	var rates []*model.LendingRateItem

	// 常见代币列表
	tokens := []string{"USDC", "USDT", "DAI", "ETH", "WBTC"}

	for _, p := range protocols {
		lendingProtocol, ok := p.(defiPkg.LendingProtocol)
		if !ok {
			continue
		}

		chains := p.SupportedChains()
		if chain != "" {
			chains = []string{chain}
		}

		for _, c := range chains {
			for _, token := range tokens {
				// 获取存款利率
				supplyAPY, err := lendingProtocol.GetSupplyAPY(ctx, token, c)
				if err != nil {
					continue
				}

				// 获取借款利率
				stableAPY, variableAPY, err := lendingProtocol.GetBorrowAPY(ctx, token, c)
				if err != nil {
					continue
				}

				rates = append(rates, &model.LendingRateItem{
					Protocol:          p.GetDisplayName(),
					Chain:             c,
					Token:             token,
					SupplyAPY:         math.Round(supplyAPY*100) / 100,
					BorrowAPYStable:   math.Round(stableAPY*100) / 100,
					BorrowAPYVariable: math.Round(variableAPY*100) / 100,
					UtilizationRate:   0, // TODO: 从合约获取
				})
			}
		}
	}

	// 保存利率历史
	if err := s.saveLendingRateHistory(ctx, rates); err != nil {
		g.Log().Warning(ctx, "保存利率历史失败", "error", err)
	}

	return rates, nil
}

// saveLendingRateHistory 保存借贷利率历史
func (s *sDefi) saveLendingRateHistory(ctx context.Context, rates []*model.LendingRateItem) error {
	for _, rate := range rates {
		_, err := dao.LendingRateHistory.Ctx(ctx).Data(g.Map{
			dao.LendingRateHistory.Columns().Protocol:          rate.Protocol,
			dao.LendingRateHistory.Columns().Chain:             rate.Chain,
			dao.LendingRateHistory.Columns().Token:             rate.Token,
			dao.LendingRateHistory.Columns().SupplyApy:         rate.SupplyAPY,
			dao.LendingRateHistory.Columns().BorrowApyStable:   rate.BorrowAPYStable,
			dao.LendingRateHistory.Columns().BorrowApyVariable: rate.BorrowAPYVariable,
			dao.LendingRateHistory.Columns().UtilizationRate:   rate.UtilizationRate,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "插入利率历史失败")
		}
	}
	return nil
}

// GetLendingRateHistory 获取借贷利率历史
func (s *sDefi) GetLendingRateHistory(ctx context.Context, protocol string, token string, days int) ([]*model.LendingRateHistoryItem, error) {
	if days <= 0 {
		days = 30
	}

	startDate := gtime.Now().AddDate(0, 0, -days)

	var records []*defiEntity.LendingRateHistory
	err := dao.LendingRateHistory.Ctx(ctx).
		Where(dao.LendingRateHistory.Columns().Protocol, protocol).
		Where(dao.LendingRateHistory.Columns().Token, token).
		Where(dao.LendingRateHistory.Columns().RecordedAt+" >= ?", startDate).
		OrderAsc(dao.LendingRateHistory.Columns().RecordedAt).
		Scan(&records)
	if err != nil {
		return nil, gerror.Wrap(err, "查询利率历史失败")
	}

	var history []*model.LendingRateHistoryItem
	for _, r := range records {
		history = append(history, &model.LendingRateHistoryItem{
			Date:              r.RecordedAt.Format("2006-01-02"),
			SupplyAPY:         r.SupplyApy,
			BorrowAPYStable:   r.BorrowApyStable,
			BorrowAPYVariable: r.BorrowApyVariable,
			UtilizationRate:   r.UtilizationRate,
		})
	}

	return history, nil
}

// CheckHealthFactors 检查所有借贷仓位的健康因子,返回低于阈值的仓位
func (s *sDefi) CheckHealthFactors(ctx context.Context, threshold float64) ([]*model.LendingPositionItem, error) {
	if threshold <= 0 {
		threshold = 1.8 // 默认阈值 1.8
	}

	userID := consts.GetUserID(ctx)

	var positions []*defiEntity.LendingPositions
	err := dao.LendingPositions.Ctx(ctx).
		Where(dao.LendingPositions.Columns().UserId, userID).
		Where(dao.LendingPositions.Columns().HealthFactor+" < ?", threshold).
		Where(dao.LendingPositions.Columns().HealthFactor+" > 0"). // 排除无借款的仓位
		Scan(&positions)
	if err != nil {
		return nil, gerror.Wrap(err, "查询低健康因子仓位失败")
	}

	var result []*model.LendingPositionItem
	for _, pos := range positions {
		result = append(result, &model.LendingPositionItem{
			Protocol:             pos.Protocol,
			Chain:                pos.Chain,
			WalletAddr:           pos.WalletAddress,
			SupplyToken:          pos.SupplyToken,
			SupplyAmount:         pos.SupplyAmount,
			SupplyValueUSD:       pos.SupplyValueUsd,
			SupplyAPY:            pos.SupplyApy,
			BorrowToken:          pos.BorrowToken,
			BorrowAmount:         pos.BorrowAmount,
			BorrowValueUSD:       pos.BorrowValueUsd,
			BorrowAPY:            pos.BorrowApy,
			HealthFactor:         pos.HealthFactor,
			LiquidationThreshold: pos.LiquidationThreshold,
			LTV:                  pos.Ltv,
			NetAPY:               pos.NetApy,
		})
	}

	return result, nil
}

// GetLendingOptimization 获取最优借贷策略推荐
func (s *sDefi) GetLendingOptimization(ctx context.Context) (*model.LendingOptimizationResult, error) {
	// 获取当前仓位
	currentPositions, err := s.GetLendingPositions(ctx)
	if err != nil {
		return nil, err
	}

	if len(currentPositions) == 0 {
		return &model.LendingOptimizationResult{
			CurrentPositions: []*model.LendingPositionItem{},
			Recommendations:  []*model.LendingRecommendation{},
			PotentialGain:    0,
			RiskLevel:        "low",
			Summary:          "您当前没有借贷仓位",
		}, nil
	}

	// 获取所有协议的利率
	allRates, err := s.GetLendingRates(ctx, "", "")
	if err != nil {
		return nil, err
	}

	// 构建利率查找表
	rateMap := make(map[string]*model.LendingRateItem) // key: protocol_chain_token
	for _, rate := range allRates {
		key := rate.Protocol + "_" + rate.Chain + "_" + rate.Token
		rateMap[key] = rate
	}

	var recommendations []*model.LendingRecommendation
	var totalPotentialGain float64

	// 分析每个仓位
	for _, pos := range currentPositions {
		// 1. 检查健康因子
		if pos.HealthFactor > 0 && pos.HealthFactor < 2.0 {
			recommendations = append(recommendations, &model.LendingRecommendation{
				Action:       "reduce_borrow",
				FromProtocol: pos.Protocol,
				ToProtocol:   pos.Protocol,
				Token:        pos.BorrowToken,
				Amount:       pos.BorrowAmount * 0.3, // 建议减少 30% 借款
				Reason:       "健康因子偏低，建议减少借款以降低清算风险",
				ExpectedGain: 0,
			})
		}

		// 2. 寻找更高收益的存款协议
		currentKey := pos.Protocol + "_" + pos.Chain + "_" + pos.SupplyToken
		_ = currentKey // 用于后续查找

		for _, rate := range rateMap {
			if rate.Chain == pos.Chain && rate.Token == pos.SupplyToken && rate.Protocol != pos.Protocol {
				if rate.SupplyAPY > pos.SupplyAPY+0.5 { // 收益率高出 0.5% 以上
					gain := pos.SupplyValueUSD * (rate.SupplyAPY - pos.SupplyAPY) / 100
					recommendations = append(recommendations, &model.LendingRecommendation{
						Action:       "migrate",
						FromProtocol: pos.Protocol,
						ToProtocol:   rate.Protocol,
						Token:        pos.SupplyToken,
						Amount:       pos.SupplyAmount,
						Reason:       fmt.Sprintf("迁移到 %s 可获得更高收益率 (%.2f%% vs %.2f%%)", rate.Protocol, rate.SupplyAPY, pos.SupplyAPY),
						ExpectedGain: gain,
					})
					totalPotentialGain += gain
				}
			}
		}

		// 3. 寻找更低利率的借款协议
		if pos.BorrowAmount > 0 {
			for _, rate := range rateMap {
				if rate.Chain == pos.Chain && rate.Token == pos.BorrowToken && rate.Protocol != pos.Protocol {
					if rate.BorrowAPYVariable < pos.BorrowAPY-0.5 { // 利率低 0.5% 以上
						gain := pos.BorrowValueUSD * (pos.BorrowAPY - rate.BorrowAPYVariable) / 100
						recommendations = append(recommendations, &model.LendingRecommendation{
							Action:       "migrate",
							FromProtocol: pos.Protocol,
							ToProtocol:   rate.Protocol,
							Token:        pos.BorrowToken,
							Amount:       pos.BorrowAmount,
							Reason:       fmt.Sprintf("迁移到 %s 可降低借款成本 (%.2f%% vs %.2f%%)", rate.Protocol, rate.BorrowAPYVariable, pos.BorrowAPY),
							ExpectedGain: gain,
						})
						totalPotentialGain += gain
					}
				}
			}
		}
	}

	// 评估风险等级
	riskLevel := "low"
	for _, pos := range currentPositions {
		if pos.HealthFactor > 0 && pos.HealthFactor < 1.5 {
			riskLevel = "high"
			break
		} else if pos.HealthFactor > 0 && pos.HealthFactor < 2.0 {
			riskLevel = "medium"
		}
	}

	// 生成总结
	summary := fmt.Sprintf("分析了 %d 个借贷仓位，发现 %d 条优化建议，预计年化收益提升 $%.2f",
		len(currentPositions), len(recommendations), totalPotentialGain)

	return &model.LendingOptimizationResult{
		CurrentPositions: currentPositions,
		Recommendations:  recommendations,
		PotentialGain:    math.Round(totalPotentialGain*100) / 100,
		RiskLevel:        riskLevel,
		Summary:          summary,
	}, nil
}
