// Package defi Uniswap V3 集中流动性协议实现
// Uniswap V3 使用 NFT 表示流动性仓位（ERC-721），每个仓位有独立的价格区间
// 通过 NonfungiblePositionManager 合约管理仓位
// 支持 Ethereum、Polygon、Arbitrum、Optimism、Base 等链
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// Uniswap V3 NonfungiblePositionManager 合约地址（各链相同）
const (
	uniV3PositionManager = "0xc36442b4a4522e871399cd717abdd847ab11fe88"
)

// Uniswap V3 热门池子的代币对映射
var uniV3PoolPairs = map[string][]string{
	"uni-v3-pos":        {"ETH", "USDC"},  // 默认
	"uni-v3-weth-usdc":  {"ETH", "USDC"},
	"uni-v3-weth-usdt":  {"ETH", "USDT"},
	"uni-v3-wbtc-weth":  {"BTC", "ETH"},
	"uni-v3-usdc-usdt":  {"USDC", "USDT"},
	"uni-v3-weth-dai":   {"ETH", "DAI"},
	"uni-v3-wbtc-usdc":  {"BTC", "USDC"},
}

// UniswapV3Protocol Uniswap V3 集中流动性协议
type UniswapV3Protocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
}

// NewUniswapV3Protocol 创建 Uniswap V3 协议实例
func NewUniswapV3Protocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *UniswapV3Protocol {
	return &UniswapV3Protocol{
		etherscanClients: clients,
		priceClient:      priceClient,
	}
}

func (u *UniswapV3Protocol) GetName() string        { return "uniswap_v3" }
func (u *UniswapV3Protocol) GetDisplayName() string  { return "Uniswap V3" }
func (u *UniswapV3Protocol) GetType() string         { return "lp" }
func (u *UniswapV3Protocol) SupportedChains() []string {
	return []string{"ethereum", "polygon", "arbitrum", "optimism", "base"}
}

// GetPositions 获取用户在 Uniswap V3 中的集中流动性仓位
// V3 仓位是 NFT（ERC-721），每个仓位有独立的价格区间 [tickLower, tickUpper]
// 当前价格在区间内时仓位活跃，区间外时仓位不产生手续费
//
// 简化实现：通过 known_tokens 识别 V3 仓位代币，估算仓位价值
// 完整实现需要直接调用 PositionManager 合约读取 tick 和 liquidity
func (u *UniswapV3Protocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	client, ok := u.etherscanClients[chain]
	if !ok {
		return nil, nil
	}

	// 获取所有 ERC20 代币余额（包含已识别的 V3 相关代币）
	balances, err := client.GetTokenBalances(ctx, address)
	if err != nil {
		return nil, err
	}

	var positions []Position

	for _, bal := range balances {
		if bal.Protocol != "uniswap_v3" {
			continue
		}

		// 解析仓位对应的代币对
		token0, token1 := resolveV3PoolPair(bal.Symbol)

		// 估算仓位价值
		// V3 仓位价值取决于当前价格和流动性区间
		// 简化：假设仓位在活跃区间内，50/50 分配
		var valueUSD float64
		if u.priceClient != nil {
			price0 := getTokenPrice(ctx, u.priceClient, token0)
			price1 := getTokenPrice(ctx, u.priceClient, token1)
			if price0 > 0 && price1 > 0 {
				valueUSD = bal.Total * (price0 + price1) / 2
			} else if price0 > 0 {
				valueUSD = bal.Total * price0
			} else if price1 > 0 {
				valueUSD = bal.Total * price1
			}
		}

		// V3 仓位还可能有未领取的手续费奖励
		// 完整实现需要调用 PositionManager.collect() 估算
		var rewards []Token
		if valueUSD > 0 {
			// 粗略估算：假设累积手续费约为仓位价值的 0.5%
			feeEstimate := valueUSD * 0.005
			rewards = []Token{
				{Symbol: token0, Amount: 0, ValueUSD: feeEstimate / 2},
				{Symbol: token1, Amount: 0, ValueUSD: feeEstimate / 2},
			}
		}

		positions = append(positions, Position{
			Protocol:     "uniswap_v3",
			ProtocolName: "Uniswap V3",
			Type:         "lp",
			Chain:        chain,
			DepositTokens: []Token{
				{Symbol: token0, Amount: bal.Total / 2, ValueUSD: valueUSD / 2},
				{Symbol: token1, Amount: bal.Total / 2, ValueUSD: valueUSD / 2},
			},
			ReceiveTokens: []Token{
				{Symbol: "UNI-V3-POS", Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      0, // V3 APY 取决于交易量和价格区间
			Rewards:  rewards,
		})
	}

	return positions, nil
}

// resolveV3PoolPair 从 V3 仓位符号解析代币对
func resolveV3PoolPair(symbol string) (string, string) {
	lower := strings.ToLower(symbol)
	if pair, ok := uniV3PoolPairs[lower]; ok && len(pair) == 2 {
		return pair[0], pair[1]
	}
	return "ETH", "USDC"
}

// 确保实现接口
var _ DeFiProtocol = (*UniswapV3Protocol)(nil)
