// Package etherscan DeFi 协议 Token 映射表
// 将已知的 DeFi 协议代币合约地址映射到协议信息和资产类型
package etherscan

import "strings"

// DeFiTokenInfo DeFi 协议 Token 信息
type DeFiTokenInfo struct {
	Protocol    string // 协议名称（lido, rocketpool, aave, compound）
	AssetType   string // 资产类型（staking, lending, lp）
	DisplayName string // 显示名称
}

// KnownDeFiTokens 已知的 DeFi 协议 Token（键为小写合约地址）
// 仅包含 Ethereum 主网的合约地址
var KnownDeFiTokens = map[string]DeFiTokenInfo{
	// === Liquid Staking ===

	// Lido — stETH
	"0xae7ab96520de3a18e5e111b5eaab095312d7fe84": {
		Protocol: "lido", AssetType: "staking", DisplayName: "Lido Staked ETH",
	},
	// Lido — wstETH
	"0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0": {
		Protocol: "lido", AssetType: "staking", DisplayName: "Lido Wrapped stETH",
	},
	// Rocket Pool — rETH
	"0xae78736cd615f374d3085123a210448e74fc6393": {
		Protocol: "rocketpool", AssetType: "staking", DisplayName: "Rocket Pool ETH",
	},
	// Coinbase — cbETH
	"0xbe9895146f7af43049ca1c1ae358b0541ea49704": {
		Protocol: "coinbase", AssetType: "staking", DisplayName: "Coinbase Staked ETH",
	},
	// Frax — sfrxETH
	"0xac3e018457b222d93114458476f3e3416abbe38f": {
		Protocol: "frax", AssetType: "staking", DisplayName: "Frax Staked ETH",
	},

	// === Lending / Borrowing ===

	// Aave v2 aTokens（Ethereum 主网）
	"0x028171bca77440897b824ca71d1c56cac55b68a3": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave DAI",
	},
	"0xbcca60bb61934080951369a648fb03df4f96263c": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave USDC",
	},
	"0x3ed3b47dd13ec9a98b44e6204a523e766b225811": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave USDT",
	},
	// Aave v2 aWETH
	"0x030ba81f1c18d280636f32af80b9aad02cf0854e": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave WETH",
	},
	// Aave v2 aWBTC
	"0x9ff58f4ffb29fa2266ab25e75e2a8b3503311656": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave WBTC",
	},
	// Aave v3 aTokens（Ethereum 主网）
	"0x4d5f47fa6a74757f35c14fd3a6ef8e3c9bc514e8": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave v3 WETH",
	},
	"0x98c23e9d8f34fefb1b7bd6a91b7ff122f4e16f5c": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave v3 USDC",
	},
	"0x23878914efe38d27c4d67ab83ed1b93a74d4086a": {
		Protocol: "aave", AssetType: "lending", DisplayName: "Aave v3 USDT",
	},

	// Compound cTokens（Ethereum 主网）
	"0x5d3a536e4d6dbd6114cc1ead35777bab948e3643": {
		Protocol: "compound", AssetType: "lending", DisplayName: "Compound DAI",
	},
	"0x39aa39c021dfbae8fac545936693ac917d5e7563": {
		Protocol: "compound", AssetType: "lending", DisplayName: "Compound USDC",
	},
	// Compound cETH
	"0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5": {
		Protocol: "compound", AssetType: "lending", DisplayName: "Compound ETH",
	},
	// Compound cUSDT
	"0xf650c3d88d12db855b8bf7d11be6c55a4e07dcc9": {
		Protocol: "compound", AssetType: "lending", DisplayName: "Compound USDT",
	},
	// Compound cWBTC
	"0xccf4429db6322d5c611ee964527d42e5d685dd6a": {
		Protocol: "compound", AssetType: "lending", DisplayName: "Compound WBTC",
	},

	// === Uniswap V2 LP Tokens ===

	// Uniswap V2 WETH/USDT
	"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "Uniswap V2 WETH/USDT",
	},
	// Uniswap V2 WETH/USDC
	"0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "Uniswap V2 WETH/USDC",
	},
	// Uniswap V2 WETH/DAI
	"0xa478c2975ab1ea89e8196811f51a7b7ade33eb11": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "Uniswap V2 WETH/DAI",
	},
	// Uniswap V2 WBTC/WETH
	"0xbb2b8038a1640196fbe3e38816f3e67cba72d940": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "Uniswap V2 WBTC/WETH",
	},
	// Uniswap V2 USDC/USDT
	"0x3041cbd36888becc7bbcbc0045e3b1f144466f5f": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "Uniswap V2 USDC/USDT",
	},

	// === SushiSwap LP Tokens（兼容 Uniswap V2） ===

	// SushiSwap WETH/USDT
	"0x06da0fd433c1a5d7a4faa01111c044910a184553": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "SushiSwap WETH/USDT",
	},
	// SushiSwap WETH/USDC
	"0x397ff1542f962076d0bfe58ea045ffa2d347aca0": {
		Protocol: "uniswap_v2", AssetType: "lp", DisplayName: "SushiSwap WETH/USDC",
	},

	// === Uniswap V3 仓位代币 ===

	// Uniswap V3 NonfungiblePositionManager（仓位 NFT 管理合约）
	"0xc36442b4a4522e871399cd717abdd847ab11fe88": {
		Protocol: "uniswap_v3", AssetType: "lp", DisplayName: "Uniswap V3 Position",
	},

	// === Curve Finance LP Tokens ===

	// Curve 3Pool（DAI/USDC/USDT）
	"0x6c3f90f043a72fa612cbac8115ee7e52bde6e490": {
		Protocol: "curve", AssetType: "lp", DisplayName: "Curve 3Pool",
	},
	// Curve stETH/ETH
	"0x06325440d014e39736583c165c2963ba99faf14e": {
		Protocol: "curve", AssetType: "lp", DisplayName: "Curve stETH/ETH",
	},
	// Curve TriCrypto2（USDT/WBTC/WETH）
	"0xc4ad29ba4b3c580e6d59105fff484999997675ff": {
		Protocol: "curve", AssetType: "lp", DisplayName: "Curve TriCrypto2",
	},
	// Curve FRAX/USDC
	"0x3175df0976dfa876431c2e9ee6bc45b65d3473cc": {
		Protocol: "curve", AssetType: "lp", DisplayName: "Curve FRAX/USDC",
	},
	// Curve crvUSD/USDC
	"0x4dece678ceceb27446b35c672dc7d61f30bad69e": {
		Protocol: "curve", AssetType: "lp", DisplayName: "Curve crvUSD/USDC",
	},
}

// LookupDeFiToken 查找合约地址对应的 DeFi Token 信息
// 返回 info 和 found 标志
func LookupDeFiToken(contractAddress string) (DeFiTokenInfo, bool) {
	info, ok := KnownDeFiTokens[strings.ToLower(contractAddress)]
	return info, ok
}
