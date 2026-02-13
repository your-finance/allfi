// Package defi Uniswap V2 LP 协议实现
// Uniswap V2 是最经典的 AMM DEX，LP Token 代表用户在流动性池中的份额
// 支持 Ethereum、Polygon、Arbitrum 等 EVM 链
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// Uniswap V2 热门 LP 对的底层资产映射
// key: LP token 符号（小写），value: [token0符号, token1符号]
var uniV2LPPairs = map[string][]string{
	"uni-v2":      {"ETH", "USDT"},     // 默认 fallback
	"uni-v2-weth-usdc": {"ETH", "USDC"},
	"uni-v2-weth-usdt": {"ETH", "USDT"},
	"uni-v2-weth-dai":  {"ETH", "DAI"},
	"uni-v2-wbtc-weth": {"BTC", "ETH"},
	"uni-v2-usdc-usdt": {"USDC", "USDT"},
	"slp":              {"ETH", "USDT"}, // SushiSwap（兼容 V2）
}

// UniswapV2Protocol Uniswap V2 LP 协议
type UniswapV2Protocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
}

// NewUniswapV2Protocol 创建 Uniswap V2 协议实例
func NewUniswapV2Protocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *UniswapV2Protocol {
	return &UniswapV2Protocol{
		etherscanClients: clients,
		priceClient:      priceClient,
	}
}

func (u *UniswapV2Protocol) GetName() string        { return "uniswap_v2" }
func (u *UniswapV2Protocol) GetDisplayName() string  { return "Uniswap V2" }
func (u *UniswapV2Protocol) GetType() string         { return "lp" }
func (u *UniswapV2Protocol) SupportedChains() []string {
	return []string{"ethereum", "polygon", "arbitrum"}
}

// GetPositions 获取用户在 Uniswap V2 中的 LP 仓位
// 通过 Etherscan 查询已知 LP Token 余额，估算底层资产价值
func (u *UniswapV2Protocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	client, ok := u.etherscanClients[chain]
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
		if bal.Protocol != "uniswap_v2" {
			continue
		}

		// 解析 LP 对应的底层资产对
		token0, token1 := resolveV2LPPair(bal.Symbol)

		// 获取底层资产价格
		// LP Token 价值 ≈ 2 * sqrt(reserve0_value * reserve1_value) / totalSupply * userBalance
		// 简化估算：假设池子 50/50 对称，LP 价值 = balance * 2 * token0Price（以一半为基准）
		var valueUSD float64
		if u.priceClient != nil {
			price0 := getTokenPrice(ctx, u.priceClient, token0)
			price1 := getTokenPrice(ctx, u.priceClient, token1)
			// LP 价值估算：假设用户余额代表的是总 LP 供应量的一定比例
			// 由于无法直接获取 reserve，使用价格估算
			if price0 > 0 && price1 > 0 {
				// 对于每个 LP Token，底层包含等值的 token0 和 token1
				// 假设 LP Token 的价格近似 = (token0Price + token1Price) * amount / 2
				valueUSD = bal.Total * (price0 + price1) / 2
			} else if price0 > 0 {
				valueUSD = bal.Total * price0
			} else if price1 > 0 {
				valueUSD = bal.Total * price1
			}
		}

		positions = append(positions, Position{
			Protocol:     "uniswap_v2",
			ProtocolName: "Uniswap V2",
			Type:         "lp",
			Chain:        chain,
			DepositTokens: []Token{
				{Symbol: token0, Amount: bal.Total / 2, ValueUSD: valueUSD / 2},
				{Symbol: token1, Amount: bal.Total / 2, ValueUSD: valueUSD / 2},
			},
			ReceiveTokens: []Token{
				{Symbol: bal.Symbol, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      0, // V2 LP 无固定 APY，取决于交易量
		})
	}

	return positions, nil
}

// resolveV2LPPair 从 LP Token 符号解析底层资产对
func resolveV2LPPair(symbol string) (string, string) {
	lower := strings.ToLower(symbol)
	if pair, ok := uniV2LPPairs[lower]; ok && len(pair) == 2 {
		return pair[0], pair[1]
	}
	// 默认返回 ETH/USDC
	return "ETH", "USDC"
}

// getTokenPrice 获取代币价格的辅助函数
func getTokenPrice(ctx context.Context, priceClient integrations.PriceClient, symbol string) float64 {
	// 将 wrapped 代币映射到原始代币
	priceSymbol := symbol
	switch strings.ToUpper(symbol) {
	case "WETH":
		priceSymbol = "ETH"
	case "WBTC":
		priceSymbol = "BTC"
	case "WMATIC":
		priceSymbol = "MATIC"
	}
	price, err := priceClient.GetPrice(ctx, priceSymbol)
	if err != nil {
		return 0
	}
	return price
}

// 确保实现接口
var _ DeFiProtocol = (*UniswapV2Protocol)(nil)
