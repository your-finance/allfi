// Package defi Aave 协议实现
// Aave 是最大的去中心化借贷协议，用户存入资产获得 aToken 凭证
// 支持 Ethereum、Polygon、Arbitrum、Optimism 多链
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// aToken 底层资产映射表
// key: aToken 符号（小写），value: 底层资产符号
var aTokenUnderlying = map[string]string{
	"adai":   "DAI",
	"ausdc":  "USDC",
	"ausdt":  "USDT",
	"aweth":  "WETH",
	"awbtc":  "WBTC",
	"aeth":   "ETH",
	"amatic": "MATIC",
}

// Aave 各资产近似年化收益率（百分比）
// 实际生产环境应从 Aave 利率合约或 API 动态获取
var aaveAPYs = map[string]float64{
	"DAI":  3.5,
	"USDC": 3.8,
	"USDT": 4.0,
	"WETH": 1.5,
	"WBTC": 0.5,
	"ETH":  1.5,
}

// AaveProtocol Aave 借贷协议
type AaveProtocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
}

// NewAaveProtocol 创建 Aave 协议实例
func NewAaveProtocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *AaveProtocol {
	return &AaveProtocol{
		etherscanClients: clients,
		priceClient:      priceClient,
	}
}

func (a *AaveProtocol) GetName() string       { return "aave" }
func (a *AaveProtocol) GetDisplayName() string { return "Aave" }
func (a *AaveProtocol) GetType() string        { return "lending" }
func (a *AaveProtocol) SupportedChains() []string {
	return []string{"ethereum", "polygon", "arbitrum", "optimism"}
}

// GetPositions 获取用户在 Aave 中的仓位
// aToken 是 1:1 映射底层资产的凭证代币，余额即为存款金额
func (a *AaveProtocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
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

		// 获取 APY
		apy := aaveAPYs[underlying]

		positions = append(positions, Position{
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
			APY:      apy,
		})
	}

	return positions, nil
}

// resolveAaveUnderlying 从 aToken 符号解析底层资产
// 如 "aUSDC" → "USDC"，"aWETH" → "WETH"
func resolveAaveUnderlying(symbol string) string {
	lower := strings.ToLower(symbol)
	if underlying, ok := aTokenUnderlying[lower]; ok {
		return underlying
	}
	// 通用规则：去掉 "a" 前缀
	if len(symbol) > 1 && (symbol[0] == 'a' || symbol[0] == 'A') {
		return strings.ToUpper(symbol[1:])
	}
	return symbol
}

// 确保实现接口
var _ DeFiProtocol = (*AaveProtocol)(nil)
