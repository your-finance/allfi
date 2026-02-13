// Package defi Rocket Pool 协议实现
// Rocket Pool 是去中心化的 ETH 质押协议，用户质押 ETH 获得 rETH
package defi

import (
	"context"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// rETH 合约地址（Ethereum 主网，小写）
const rETHContract = "0xae78736cd615f374d3085123a210448e74fc6393"

// RocketPoolProtocol Rocket Pool 质押协议
type RocketPoolProtocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
	// apy Rocket Pool 当前年化收益率（百分比）
	apy float64
}

// NewRocketPoolProtocol 创建 Rocket Pool 协议实例
func NewRocketPoolProtocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *RocketPoolProtocol {
	return &RocketPoolProtocol{
		etherscanClients: clients,
		priceClient:      priceClient,
		apy:              3.2, // Rocket Pool 近似 APY
	}
}

func (r *RocketPoolProtocol) GetName() string        { return "rocketpool" }
func (r *RocketPoolProtocol) GetDisplayName() string  { return "Rocket Pool" }
func (r *RocketPoolProtocol) GetType() string         { return "staking" }
func (r *RocketPoolProtocol) SupportedChains() []string { return []string{"ethereum"} }

// GetPositions 获取用户在 Rocket Pool 中的仓位
// 通过 Etherscan 查询用户持有的 rETH 余额，按兑换率计算 ETH 仓位价值
func (r *RocketPoolProtocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	if chain != "ethereum" {
		return nil, nil
	}

	client, ok := r.etherscanClients["ethereum"]
	if !ok {
		return nil, nil
	}

	// 获取所有 ERC20 代币余额
	balances, err := client.GetTokenBalances(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取 ETH 价格
	ethPrice, err := r.priceClient.GetPrice(ctx, "ETH")
	if err != nil {
		ethPrice = 0
	}

	var positions []Position

	for _, bal := range balances {
		if bal.Protocol != "rocketpool" {
			continue
		}

		// rETH 与 ETH 的兑换率约 1.10:1（随时间增长，因为 rETH 累积质押奖励）
		exchangeRate := 1.10
		ethAmount := bal.Total * exchangeRate
		valueUSD := ethAmount * ethPrice

		positions = append(positions, Position{
			Protocol:     "rocketpool",
			ProtocolName: "Rocket Pool",
			Type:         "staking",
			Chain:        "ethereum",
			DepositTokens: []Token{
				{Symbol: "ETH", Amount: ethAmount, ValueUSD: valueUSD},
			},
			ReceiveTokens: []Token{
				{Symbol: "rETH", Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      r.apy,
		})
	}

	return positions, nil
}

// 确保实现接口
var _ DeFiProtocol = (*RocketPoolProtocol)(nil)
