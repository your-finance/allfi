// Package defi Compound 借贷协议扩展实现
// 实现 LendingProtocol 接口,支持获取借贷仓位、利率、健康因子
package defi

import (
	"context"
	"fmt"
	"math"
)

// GetLendingPositions 获取用户的借贷仓位（包含存款、借款、健康因子）
// 实现 LendingProtocol 接口
func (c *CompoundProtocol) GetLendingPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	if chain != "ethereum" {
		return nil, nil // Compound v2 仅支持 Ethereum 主网
	}

	client, ok := c.etherscanClients["ethereum"]
	if !ok {
		return nil, nil
	}

	// 获取所有 ERC20 代币余额
	balances, err := client.GetTokenBalances(ctx, address)
	if err != nil {
		return nil, err
	}

	var positions []Position

	// 存款仓位（cToken）
	supplyPositions := make(map[string]*Position) // key: underlying token

	for _, bal := range balances {
		if bal.Protocol != "compound" {
			continue
		}

		// 解析 cToken 信息
		info := resolveCompoundUnderlying(bal.Symbol)

		// 计算底层资产数量 = cToken 余额 × 兑换率
		underlyingAmount := bal.Total * info.ExchangeRate

		// 获取底层资产价格
		var valueUSD float64
		if c.priceClient != nil {
			priceSymbol := info.Underlying
			if priceSymbol == "WBTC" {
				priceSymbol = "BTC"
			}
			price, err := c.priceClient.GetPrice(ctx, priceSymbol)
			if err == nil {
				valueUSD = underlyingAmount * price
			}
		}

		// 获取存款 APY
		supplyAPY := compoundAPYs[info.Underlying]

		pos := &Position{
			Protocol:     "compound",
			ProtocolName: "Compound",
			Type:         "lending",
			Chain:        "ethereum",
			DepositTokens: []Token{
				{Symbol: info.Underlying, Amount: underlyingAmount, ValueUSD: valueUSD},
			},
			ReceiveTokens: []Token{
				{Symbol: bal.Symbol, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      supplyAPY,
		}

		supplyPositions[info.Underlying] = pos
	}

	// TODO: 从链上获取借款数据
	// 这里需要调用 Compound Comptroller 合约的 getAccountLiquidity 方法
	// 暂时使用模拟数据演示结构

	// 将存款仓位添加到结果
	for _, pos := range supplyPositions {
		// 如果有借款,设置借贷相关字段
		// pos.BorrowTokens = []Token{...}
		// pos.HealthFactor = healthFactor
		// pos.LiquidationThreshold = 0.75 // 75%
		// pos.LTV = totalDebtUSD / totalCollateralUSD * 100
		// pos.NetAPY = supplyAPY - borrowAPY

		positions = append(positions, *pos)
	}

	return positions, nil
}

// GetSupplyAPY 获取指定代币的存款年化收益率
func (c *CompoundProtocol) GetSupplyAPY(ctx context.Context, token string, chain string) (float64, error) {
	if chain != "ethereum" {
		return 0, fmt.Errorf("chain %s not supported", chain)
	}

	// 从静态映射表获取（生产环境应从 Compound API 或合约获取）
	if apy, ok := compoundAPYs[token]; ok {
		return apy, nil
	}
	return 0, fmt.Errorf("token %s not supported", token)
}

// GetBorrowAPY 获取指定代币的借款年化利率
// Compound V2 只有浮动利率,稳定利率返回 0
func (c *CompoundProtocol) GetBorrowAPY(ctx context.Context, token string, chain string) (stable float64, variable float64, err error) {
	if chain != "ethereum" {
		return 0, 0, fmt.Errorf("chain %s not supported", chain)
	}

	// 借款利率通常高于存款利率
	supplyAPY, err := c.GetSupplyAPY(ctx, token, chain)
	if err != nil {
		return 0, 0, err
	}

	// Compound V2 只有浮动利率
	variable = supplyAPY * 1.8
	stable = 0 // Compound V2 不支持稳定利率

	return stable, variable, nil
}

// GetHealthFactor 获取用户的健康因子
// Compound 使用 "账户流动性" 概念,这里转换为健康因子
// healthFactor = (totalCollateral * collateralFactor) / totalBorrow
func (c *CompoundProtocol) GetHealthFactor(ctx context.Context, address string, chain string) (float64, error) {
	if chain != "ethereum" {
		return 0, fmt.Errorf("chain %s not supported", chain)
	}

	// TODO: 从 Compound Comptroller 合约获取真实数据
	// 调用 getAccountLiquidity(address) 方法
	// 返回: error, liquidity, shortfall
	// healthFactor = liquidity / (liquidity + shortfall)

	// 暂时返回模拟数据
	return 3.0, nil // 健康因子 3.0 表示安全
}

// calculateCompoundHealthFactor 计算 Compound 健康因子
// healthFactor = (totalCollateral * collateralFactor) / totalBorrow
func calculateCompoundHealthFactor(totalCollateral, totalBorrow, collateralFactor float64) float64 {
	if totalBorrow == 0 {
		return math.MaxFloat64 // 无借款,健康因子无限大
	}
	return (totalCollateral * collateralFactor) / totalBorrow
}

// 确保 CompoundProtocol 实现 LendingProtocol 接口
var _ LendingProtocol = (*CompoundProtocol)(nil)
