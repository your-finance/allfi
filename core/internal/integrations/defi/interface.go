// Package defi DeFi 协议集成
// 定义 DeFi 协议通用接口，支持 Staking、Lending、LP 等仓位识别
package defi

import "context"

// DeFiProtocol DeFi 协议通用接口
// 每个 DeFi 协议（Lido、Aave、Uniswap 等）需要实现此接口
type DeFiProtocol interface {
	// GetName 获取协议标识名称（如 "lido", "aave"）
	GetName() string
	// GetDisplayName 获取显示名称（如 "Lido Finance"）
	GetDisplayName() string
	// GetPositions 获取用户在该协议中的所有仓位
	// address: 用户钱包地址
	// chain: 链名称（ethereum, polygon 等）
	GetPositions(ctx context.Context, address string, chain string) ([]Position, error)
	// SupportedChains 返回该协议支持的链列表
	SupportedChains() []string
	// GetType 获取协议类型（staking/lending/lp/vault）
	GetType() string
}
