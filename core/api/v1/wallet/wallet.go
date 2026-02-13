// Package wallet 钱包地址 API 定义
// 提供钱包地址的增删改查、余额查询、批量导入接口
package wallet

import "github.com/gogf/gf/v2/frame/g"

// ListAddressesReq 获取钱包地址列表请求
type ListAddressesReq struct {
	g.Meta `path:"/wallets/addresses" method:"get" summary:"获取钱包地址列表" tags:"钱包"`
}

// ListAddressesRes 获取钱包地址列表响应
type ListAddressesRes struct {
	Addresses []AddressItem `json:"addresses" dc:"钱包地址列表"`
}

// AddressItem 钱包地址条目
type AddressItem struct {
	ID         uint   `json:"id" dc:"地址 ID"`
	Blockchain string `json:"blockchain" dc:"区块链网络"`
	Address    string `json:"address" dc:"钱包地址"`
	Label      string `json:"label" dc:"地址标签"`
	CreatedAt  string `json:"created_at" dc:"创建时间"`
	UpdatedAt  string `json:"updated_at" dc:"更新时间"`
}

// CreateAddressReq 添加钱包地址请求
type CreateAddressReq struct {
	g.Meta     `path:"/wallets/addresses" method:"post" summary:"添加钱包地址" tags:"钱包"`
	Blockchain string `json:"blockchain" v:"required|in:ethereum,bsc,polygon,arbitrum,optimism,base" dc:"区块链网络"`
	Address    string `json:"address" v:"required" dc:"钱包地址（0x 开头）"`
	Label      string `json:"label" dc:"地址标签"`
}

// CreateAddressRes 添加钱包地址响应
type CreateAddressRes struct {
	Address *AddressItem `json:"address" dc:"新建的钱包地址信息"`
}

// GetAddressReq 获取单个钱包地址请求
type GetAddressReq struct {
	g.Meta `path:"/wallets/addresses/{id}" method:"get" summary:"获取单个钱包地址" tags:"钱包"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"地址 ID"`
}

// GetAddressRes 获取单个钱包地址响应
type GetAddressRes struct {
	Address *AddressItem `json:"address" dc:"钱包地址信息"`
}

// UpdateAddressReq 更新钱包地址请求
type UpdateAddressReq struct {
	g.Meta     `path:"/wallets/addresses/{id}" method:"put" summary:"更新钱包地址" tags:"钱包"`
	Id         uint   `json:"id" in:"path" v:"required|min:1" dc:"地址 ID"`
	Blockchain string `json:"blockchain" dc:"区块链网络（可选更新）"`
	Address    string `json:"address" dc:"钱包地址（可选更新）"`
	Label      string `json:"label" dc:"地址标签"`
}

// UpdateAddressRes 更新钱包地址响应
type UpdateAddressRes struct {
	Address *AddressItem `json:"address" dc:"更新后的钱包地址信息"`
}

// DeleteAddressReq 删除钱包地址请求
type DeleteAddressReq struct {
	g.Meta `path:"/wallets/addresses/{id}" method:"delete" summary:"删除钱包地址" tags:"钱包"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"地址 ID"`
}

// DeleteAddressRes 删除钱包地址响应
type DeleteAddressRes struct{}

// GetAddressBalancesReq 获取钱包地址余额请求
type GetAddressBalancesReq struct {
	g.Meta `path:"/wallets/addresses/{id}/balances" method:"get" summary:"获取钱包地址余额" tags:"钱包"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"地址 ID"`
}

// GetAddressBalancesRes 获取钱包地址余额响应
type GetAddressBalancesRes struct {
	NativeBalance float64            `json:"native_balance" dc:"原生代币余额"`
	TokenBalances map[string]float64 `json:"token_balances" dc:"Token 余额映射"`
	TotalValueUSD float64            `json:"total_value_usd" dc:"总价值（USD）"`
}

// BatchAddress 批量导入的单个地址
type BatchAddress struct {
	Blockchain string `json:"blockchain" v:"required" dc:"区块链网络"`
	Address    string `json:"address" v:"required" dc:"钱包地址"`
	Label      string `json:"label" dc:"地址标签"`
}

// BatchImportReq 批量导入钱包地址请求
type BatchImportReq struct {
	g.Meta    `path:"/wallets/batch" method:"post" summary:"批量导入钱包地址" tags:"钱包"`
	Addresses []BatchAddress `json:"addresses" v:"required" dc:"地址列表"`
}

// BatchImportRes 批量导入钱包地址响应
type BatchImportRes struct {
	Imported int `json:"imported" dc:"成功导入数量"`
	Failed   int `json:"failed" dc:"导入失败数量"`
}

// SyncAddressReq 同步钱包地址请求
type SyncAddressReq struct {
	g.Meta `path:"/wallets/addresses/{id}/sync" method:"post" summary:"同步钱包地址" tags:"钱包"`
	Id     int `json:"id" v:"required" in:"path" dc:"钱包 ID"`
}

// SyncAddressRes 同步钱包地址响应
type SyncAddressRes struct {
	Message string `json:"message"`
}
