// Package defi Curve Finance LP 协议实现
// Curve 专注于稳定币和同类资产互换，提供低滑点的多资产流动性池
// 支持 Ethereum、Polygon、Arbitrum 等链
package defi

import (
	"context"
	"strings"

	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/etherscan"
)

// Curve 池子的底层资产映射
// key: Curve LP token 符号（小写），value: 池中资产列表
var curveLPAssets = map[string][]string{
	"3crv":       {"DAI", "USDC", "USDT"},       // 3Pool
	"crv3crypto": {"USDT", "BTC", "ETH"},         // TriCrypto
	"crvfrax":    {"FRAX", "USDC"},               // FRAX/USDC
	"steth-lp":   {"ETH", "stETH"},               // stETH/ETH
	"crvusd":     {"crvUSD", "USDC"},             // crvUSD
	"crveth":     {"ETH", "CRV"},                 // CRV/ETH
}

// Curve 池子近似 APY（百分比）
var curveAPYs = map[string]float64{
	"3crv":       2.5,
	"crv3crypto": 5.0,
	"crvfrax":    3.0,
	"steth-lp":   3.5,
	"crvusd":     4.0,
	"crveth":     8.0,
}

// CurveProtocol Curve Finance 协议
type CurveProtocol struct {
	// etherscanClients 按链名称索引的 Etherscan 客户端
	etherscanClients map[string]*etherscan.Client
	// priceClient 价格查询客户端
	priceClient integrations.PriceClient
}

// NewCurveProtocol 创建 Curve 协议实例
func NewCurveProtocol(clients map[string]*etherscan.Client, priceClient integrations.PriceClient) *CurveProtocol {
	return &CurveProtocol{
		etherscanClients: clients,
		priceClient:      priceClient,
	}
}

func (c *CurveProtocol) GetName() string        { return "curve" }
func (c *CurveProtocol) GetDisplayName() string  { return "Curve Finance" }
func (c *CurveProtocol) GetType() string         { return "lp" }
func (c *CurveProtocol) SupportedChains() []string {
	return []string{"ethereum", "polygon", "arbitrum"}
}

// GetPositions 获取用户在 Curve 中的 LP 仓位
// Curve LP Token 代表用户在多资产池中的份额
// 对于稳定币池（3Pool），LP 价值近似 = balance * 1 USD（因为底层都是稳定币）
// 对于混合池（TriCrypto），需要根据池中各资产价格加权计算
func (c *CurveProtocol) GetPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	client, ok := c.etherscanClients[chain]
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
		if bal.Protocol != "curve" {
			continue
		}

		// 解析池中底层资产
		assets := resolveCurveAssets(bal.Symbol)

		// 估算 LP 价值
		var valueUSD float64
		var depositTokens []Token

		if isStablecoinPool(assets) {
			// 稳定币池：每个 LP Token 约值 1 USD
			valueUSD = bal.Total
			sharePerAsset := bal.Total / float64(len(assets))
			for _, asset := range assets {
				depositTokens = append(depositTokens, Token{
					Symbol:   asset,
					Amount:   sharePerAsset,
					ValueUSD: sharePerAsset,
				})
			}
		} else {
			// 混合池：根据各资产价格估算
			if c.priceClient != nil {
				var totalPrice float64
				for _, asset := range assets {
					price := getTokenPrice(ctx, c.priceClient, asset)
					totalPrice += price
				}
				if totalPrice > 0 {
					valueUSD = bal.Total * totalPrice / float64(len(assets))
				}
			}
			sharePerAsset := bal.Total / float64(len(assets))
			for _, asset := range assets {
				assetPrice := getTokenPrice(ctx, c.priceClient, asset)
				depositTokens = append(depositTokens, Token{
					Symbol:   asset,
					Amount:   sharePerAsset,
					ValueUSD: sharePerAsset * assetPrice,
				})
			}
		}

		// 获取池 APY
		apy := curveAPYs[strings.ToLower(bal.Symbol)]

		positions = append(positions, Position{
			Protocol:      "curve",
			ProtocolName:  "Curve Finance",
			Type:          "lp",
			Chain:         chain,
			DepositTokens: depositTokens,
			ReceiveTokens: []Token{
				{Symbol: bal.Symbol, Amount: bal.Total, ValueUSD: valueUSD},
			},
			ValueUSD: valueUSD,
			APY:      apy,
		})
	}

	return positions, nil
}

// resolveCurveAssets 从 Curve LP Token 符号解析池中底层资产
func resolveCurveAssets(symbol string) []string {
	lower := strings.ToLower(symbol)
	if assets, ok := curveLPAssets[lower]; ok {
		return assets
	}
	// 默认返回 3Pool 资产
	return []string{"DAI", "USDC", "USDT"}
}

// isStablecoinPool 判断是否为稳定币池
func isStablecoinPool(assets []string) bool {
	stablecoins := map[string]bool{
		"DAI": true, "USDC": true, "USDT": true, "FRAX": true,
		"BUSD": true, "TUSD": true, "GUSD": true, "crvUSD": true,
	}
	for _, asset := range assets {
		if !stablecoins[asset] {
			return false
		}
	}
	return true
}

// 确保实现接口
var _ DeFiProtocol = (*CurveProtocol)(nil)
