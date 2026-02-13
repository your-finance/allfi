// Package defi Lido 协议实现
// Lido 是最大的 ETH 流动性质押协议，用户质押 ETH 获得 stETH
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// stETH 和 wstETH 的合约地址（Ethereum 主网，小写）
const (
	stETHContract  = "0xae7ab96520de3a18e5e111b5eaab095312d7fe84"
	wstETHContract = "0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0"
)

// LidoProtocol Lido 质押协议
type LidoProtocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
	// apy Lido 当前年化收益率（百分比）
	// 实际生产环境应从 Lido API 动态获取，此处使用近似值
	apy float64
}

// NewLidoProtocol 创建 Lido 协议实例
// clients: Etherscan 客户端映射（key 为链名称）
// priceClient: 价格查询客户端（用于获取 ETH 价格）
func NewLidoProtocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *LidoProtocol {
	return &LidoProtocol{
		etherscanClients: clients,
		priceClient:      priceClient,
		apy:              3.5, // Lido 近似 APY
	}
}

func (l *LidoProtocol) GetName() string        { return "lido" }
func (l *LidoProtocol) GetDisplayName() string  { return "Lido Finance" }
func (l *LidoProtocol) GetType() string         { return "staking" }
func (l *LidoProtocol) SupportedChains() []string { return []string{"ethereum"} }

// GetPositions 获取用户在 Lido 中的仓位
// 通过 Etherscan 查询用户持有的 stETH/wstETH 余额，计算仓位价值
func (l *LidoProtocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	if chain != "ethereum" {
		return nil, nil // Lido 仅支持 Ethereum 主网
	}

	client, ok := l.etherscanClients["ethereum"]
	if !ok {
		return nil, nil
	}

	// 获取所有 ERC20 代币余额
	balances, err := client.GetTokenBalances(ctx, address)
	if err != nil {
		return nil, err
	}

	// 获取 ETH 价格
	ethPrice, err := l.priceClient.GetPrice(ctx, "ETH")
	if err != nil {
		ethPrice = 0 // 价格获取失败时仍返回仓位，价值为 0
	}

	var positions []Position

	for _, bal := range balances {
		if bal.Protocol != "lido" {
			continue
		}

		// stETH 与 ETH 近似 1:1 兑换（实际存在微小偏差）
		// wstETH 与 ETH 的兑换率约 1.15:1（随时间增长）
		exchangeRate := 1.0
		if strings.ToLower(bal.Symbol) == "wsteth" {
			exchangeRate = 1.15 // wstETH → ETH 近似兑换率
		}

		ethAmount := bal.Total * exchangeRate
		valueUSD := ethAmount * ethPrice

		positions = append(positions, Position{
			Protocol:     "lido",
			ProtocolName: "Lido Finance",
			Type:         "staking",
			Chain:        "ethereum",
			DepositTokens: []Token{
				{Symbol: "ETH", Amount: ethAmount, ValueUSD: valueUSD},
			},
			ReceiveTokens: []Token{
				{Symbol: bal.Symbol, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      l.apy,
		})
	}

	return positions, nil
}

// 确保实现接口
var _ DeFiProtocol = (*LidoProtocol)(nil)
