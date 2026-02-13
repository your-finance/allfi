// Package defi Compound 协议实现
// Compound 是老牌去中心化借贷协议，用户存入资产获得 cToken 凭证
// cToken 与底层资产的兑换率随时间增长（累积利息）
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// cToken 底层资产映射表及兑换率
// key: cToken 符号（小写）
type cTokenInfo struct {
	Underlying   string  // 底层资产符号
	ExchangeRate float64 // cToken → 底层资产的兑换率（近似值，实际应从合约读取）
}

var cTokenMapping = map[string]cTokenInfo{
	"cdai":  {Underlying: "DAI", ExchangeRate: 0.0225},
	"cusdc": {Underlying: "USDC", ExchangeRate: 0.0228},
	"cusdt": {Underlying: "USDT", ExchangeRate: 0.0223},
	"ceth":  {Underlying: "ETH", ExchangeRate: 0.0204},
	"cwbtc": {Underlying: "WBTC", ExchangeRate: 0.0205},
}

// Compound 各资产近似年化收益率（百分比）
var compoundAPYs = map[string]float64{
	"DAI":  2.5,
	"USDC": 2.8,
	"USDT": 3.0,
	"ETH":  1.0,
	"WBTC": 0.3,
}

// CompoundProtocol Compound 借贷协议
type CompoundProtocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
}

// NewCompoundProtocol 创建 Compound 协议实例
func NewCompoundProtocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *CompoundProtocol {
	return &CompoundProtocol{
		etherscanClients: clients,
		priceClient:      priceClient,
	}
}

func (c *CompoundProtocol) GetName() string        { return "compound" }
func (c *CompoundProtocol) GetDisplayName() string  { return "Compound" }
func (c *CompoundProtocol) GetType() string         { return "lending" }
func (c *CompoundProtocol) SupportedChains() []string { return []string{"ethereum"} }

// GetPositions 获取用户在 Compound 中的仓位
// cToken 余额需要乘以 exchangeRate 才能得到底层资产数量
func (c *CompoundProtocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
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

		// 获取 APY
		apy := compoundAPYs[info.Underlying]

		positions = append(positions, Position{
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
			APY:      apy,
		})
	}

	return positions, nil
}

// resolveCompoundUnderlying 从 cToken 符号解析底层资产信息
func resolveCompoundUnderlying(symbol string) cTokenInfo {
	lower := strings.ToLower(symbol)
	if info, ok := cTokenMapping[lower]; ok {
		return info
	}
	// 通用规则：去掉 "c" 前缀，使用默认兑换率
	underlying := symbol
	if len(symbol) > 1 && (symbol[0] == 'c' || symbol[0] == 'C') {
		underlying = strings.ToUpper(symbol[1:])
	}
	return cTokenInfo{Underlying: underlying, ExchangeRate: 0.02}
}

// 确保实现接口
var _ DeFiProtocol = (*CompoundProtocol)(nil)
