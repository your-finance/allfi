// Package defi Aave 借贷协议扩展实现
// 实现 LendingProtocol 接口,支持获取借贷仓位、利率、健康因子
package defi

import (
	"context"
	"fmt"
	"math"
)

// GetLendingPositions 获取用户的借贷仓位（包含存款、借款、健康因子）
// 实现 LendingProtocol 接口
func (a *AaveProtocol) GetLendingPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	client, ok := a.etherscanClients[chain]
	if !ok {
		return nil, nil
	}

	// 获取所有 ERC20 代币余额
	balances, err := client.GetTokenBalances(ctx, address)
	if err != nil {
		return nil, err
	}

	var positions []Position

	// 存款仓位（aToken）
	supplyPositions := make(map[string]*Position) // key: underlying token

	for _, bal := range balances {
		if bal.Protocol != "aave" {
			continue
		}

		// 解析底层资产符号
		underlying := resolveAaveUnderlying(bal.Symbol)

		// 获取底层资产价格
		var valueUSD float64
		if a.priceClient != nil {
			priceSymbol := underlying
			if priceSymbol == "WETH" {
				priceSymbol = "ETH"
			}
			if priceSymbol == "WBTC" {
				priceSymbol = "BTC"
			}
			price, err := a.priceClient.GetPrice(ctx, priceSymbol)
			if err == nil {
				valueUSD = bal.Total * price
			}
		}

		// 获取存款 APY
		supplyAPY := aaveAPYs[underlying]

		pos := &Position{
			Protocol:     "aave",
			ProtocolName: "Aave",
			Type:         "lending",
			Chain:        chain,
			DepositTokens: []Token{
				{Symbol: underlying, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ReceiveTokens: []Token{
				{Symbol: bal.Symbol, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      supplyAPY,
		}

		supplyPositions[underlying] = pos
	}

	// TODO: 从链上获取借款数据
	// 这里需要调用 Aave V3 Pool 合约的 getUserAccountData 方法
	// 暂时使用模拟数据演示结构

	// 计算健康因子（需要从合约获取真实数据）
	// healthFactor := calculateHealthFactor(totalCollateralUSD, totalDebtUSD, liquidationThreshold)

	// 将存款仓位添加到结果
	for _, pos := range supplyPositions {
		// 如果有借款,设置借贷相关字段
		// pos.BorrowTokens = []Token{...}
		// pos.HealthFactor = healthFactor
		// pos.LiquidationThreshold = 0.85 // 85%
		// pos.LTV = totalDebtUSD / totalCollateralUSD * 100
		// pos.NetAPY = supplyAPY - borrowAPY

		positions = append(positions, *pos)
	}

	return positions, nil
}

// GetSupplyAPY 获取指定代币的存款年化收益率
func (a *AaveProtocol) GetSupplyAPY(ctx context.Context, token string, chain string) (float64, error) {
	// 从静态映射表获取（生产环境应从 Aave API 或合约获取）
	if apy, ok := aaveAPYs[token]; ok {
		return apy, nil
	}
	return 0, fmt.Errorf("token %s not supported", token)
}

// GetBorrowAPY 获取指定代币的借款年化利率
// 返回 (稳定利率, 浮动利率, error)
func (a *AaveProtocol) GetBorrowAPY(ctx context.Context, token string, chain string) (stable float64, variable float64, err error) {
	// 借款利率通常高于存款利率
	// 这里使用简化计算：借款利率 = 存款利率 * 1.5
	supplyAPY, err := a.GetSupplyAPY(ctx, token, chain)
	if err != nil {
		return 0, 0, err
	}

	// 稳定利率通常高于浮动利率
	variable = supplyAPY * 1.5
	stable = variable * 1.2

	return stable, variable, nil
}

// GetHealthFactor 获取用户的健康因子
// 健康因子 = (抵押品价值 * 清算阈值) / 借款价值
// > 1: 安全；< 1: 可能被清算
func (a *AaveProtocol) GetHealthFactor(ctx context.Context, address string, chain string) (float64, error) {
	// TODO: 从 Aave V3 Pool 合约获取真实数据
	// 调用 getUserAccountData(address) 方法
	// 返回: totalCollateralBase, totalDebtBase, availableBorrowsBase, currentLiquidationThreshold, ltv, healthFactor

	// 暂时返回模拟数据
	// 生产环境需要通过 RPC 调用合约
	return 2.5, nil // 健康因子 2.5 表示安全
}

// calculateHealthFactor 计算健康因子
// healthFactor = (totalCollateral * liquidationThreshold) / totalDebt
func calculateHealthFactor(totalCollateral, totalDebt, liquidationThreshold float64) float64 {
	if totalDebt == 0 {
		return math.MaxFloat64 // 无借款,健康因子无限大
	}
	return (totalCollateral * liquidationThreshold) / totalDebt
}

// 确保 AaveProtocol 实现 LendingProtocol 接口
var _ LendingProtocol = (*AaveProtocol)(nil)
