// Package etherscan EVM 链配置
// 集中管理所有支持的 EVM 兼容链的 Explorer API 地址
package etherscan

// ChainConfig EVM 链配置
type ChainConfig struct {
	Name         string // 链名称
	ChainID      int    // 链ID
	BaseURL      string // Explorer API 地址
	NativeSymbol string // 原生代币符号
	RateLimitKey string // 限流键名
	PublicRPC    string // 公共 RPC 端点（免费，无需 API Key）
}

// SupportedChains 支持的链配置
var SupportedChains = map[string]ChainConfig{
	"ethereum": {Name: "ethereum", ChainID: 1, BaseURL: "https://api.etherscan.io/api", NativeSymbol: "ETH", RateLimitKey: "etherscan", PublicRPC: "https://eth.llamarpc.com"},
	"bsc":      {Name: "bsc", ChainID: 56, BaseURL: "https://api.bscscan.com/api", NativeSymbol: "BNB", RateLimitKey: "bscscan", PublicRPC: "https://bsc-dataseed.binance.org"},
	"arbitrum": {Name: "arbitrum", ChainID: 42161, BaseURL: "https://api.arbiscan.io/api", NativeSymbol: "ETH", RateLimitKey: "arbiscan", PublicRPC: "https://arb1.arbitrum.io/rpc"},
	"optimism": {Name: "optimism", ChainID: 10, BaseURL: "https://api-optimistic.etherscan.io/api", NativeSymbol: "ETH", RateLimitKey: "optimism", PublicRPC: "https://mainnet.optimism.io"},
	"polygon":  {Name: "polygon", ChainID: 137, BaseURL: "https://api.polygonscan.com/api", NativeSymbol: "MATIC", RateLimitKey: "polygonscan", PublicRPC: ""},
	"base":     {Name: "base", ChainID: 8453, BaseURL: "https://api.basescan.org/api", NativeSymbol: "ETH", RateLimitKey: "basescan", PublicRPC: "https://mainnet.base.org"},
}
