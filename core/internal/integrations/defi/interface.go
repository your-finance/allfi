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

// LendingProtocol 借贷协议扩展接口
// 支持借贷功能的协议（Aave、Compound）需要实现此接口
type LendingProtocol interface {
	DeFiProtocol

	// GetLendingPositions 获取用户的借贷仓位（包含存款、借款、健康因子）
	GetLendingPositions(ctx context.Context, address string, chain string) ([]Position, error)

	// GetSupplyAPY 获取指定代币的存款年化收益率
	GetSupplyAPY(ctx context.Context, token string, chain string) (float64, error)

	// GetBorrowAPY 获取指定代币的借款年化利率
	// 返回 (稳定利率, 浮动利率, error)
	GetBorrowAPY(ctx context.Context, token string, chain string) (stable float64, variable float64, err error)

	// GetHealthFactor 获取用户的健康因子
	// 健康因子 = (抵押品价值 * 清算阈值) / 借款价值
	// > 1: 安全；< 1: 可能被清算
	GetHealthFactor(ctx context.Context, address string, chain string) (float64, error)
}
