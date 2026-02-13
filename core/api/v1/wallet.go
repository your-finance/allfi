// Package v1 钱包 API 定义
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/model/entity"
)

// AddWalletReq 添加钱包地址请求
type AddWalletReq struct {
	g.Meta     `path:"/wallets/addresses" method:"post" summary:"添加钱包地址" tags:"钱包"`
	Blockchain string `json:"blockchain" v:"required|in:ethereum,bsc,polygon,arbitrum,optimism" dc:"区块链网络"`
	Address    string `json:"address" v:"required|length:42,42" dc:"钱包地址（0x开头）"`
	Label      string `json:"label" v:"max-length:100" dc:"地址标签"`
}

// AddWalletRes 添加钱包地址响应
type AddWalletRes struct {
	Wallet *entity.WalletAddress `json:"wallet" dc:"钱包地址信息"`
}

// ListWalletsReq 获取钱包列表请求
type ListWalletsReq struct {
	g.Meta     `path:"/wallets/addresses" method:"get" summary:"获取钱包地址列表" tags:"钱包"`
	Blockchain string `json:"blockchain" v:"in:ethereum,bsc,polygon,arbitrum,optimism" dc:"区块链网络（可选）"`
}

// ListWalletsRes 获取钱包列表响应
type ListWalletsRes struct {
	Wallets []*entity.WalletAddress `json:"wallets" dc:"钱包地址列表"`
}

// DeleteWalletReq 删除钱包地址请求
type DeleteWalletReq struct {
	g.Meta   `path:"/wallets/addresses/:id" method:"delete" summary:"删除钱包地址" tags:"钱包"`
	WalletId uint `json:"id" v:"required|min:1" dc:"钱包ID"`
}

// DeleteWalletRes 删除钱包地址响应
type DeleteWalletRes struct{}

// GetWalletBalanceReq 获取钱包余额请求
type GetWalletBalanceReq struct {
	g.Meta   `path:"/wallets/addresses/:id/balance" method:"get" summary:"获取钱包余额" tags:"钱包"`
	WalletId uint `json:"id" v:"required|min:1" dc:"钱包ID"`
}

// GetWalletBalanceRes 获取钱包余额响应
type GetWalletBalanceRes struct {
	NativeBalance float64           `json:"native_balance" dc:"原生代币余额"`
	TokenBalances map[string]float64 `json:"token_balances" dc:"Token 余额"`
	TotalValueUSD float64           `json:"total_value_usd" dc:"总价值（USD）"`
}
