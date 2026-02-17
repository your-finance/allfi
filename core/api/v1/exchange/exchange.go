// Package exchange 交易所账户 API 定义
// 提供交易所账户的增删改查、连接测试、余额查询接口
package exchange

import "github.com/gogf/gf/v2/frame/g"

// ListSupportedExchangesReq 获取支持的交易所列表请求
type ListSupportedExchangesReq struct {
	g.Meta `path:"/exchanges/supported" method:"get" summary:"获取支持的交易所列表" tags:"交易所"`
}

// ListSupportedExchangesRes 获取支持的交易所列表响应
type ListSupportedExchangesRes struct {
	Exchanges []ExchangeInfo `json:"exchanges" dc:"交易所列表"`
}

// ExchangeInfo 交易所信息
type ExchangeInfo struct {
	ID       string `json:"id" dc:"交易所 ID（用于 API 调用）"`
	Name     string `json:"name" dc:"交易所显示名称"`
	Category string `json:"category" dc:"交易所分类（spot/futures/derivatives 等）"`
}

// ListAccountsReq 获取交易所账户列表请求
type ListAccountsReq struct {
	g.Meta `path:"/exchanges/accounts" method:"get" summary:"获取交易所账户列表" tags:"交易所"`
}

// ListAccountsRes 获取交易所账户列表响应
type ListAccountsRes struct {
	Accounts []AccountItem `json:"accounts" dc:"账户列表"`
}

// AccountItem 交易所账户条目
type AccountItem struct {
	ID           uint   `json:"id" dc:"账户 ID"`
	ExchangeName string `json:"exchange_name" dc:"交易所名称"`
	Label        string `json:"label" dc:"账户标签"`
	Note         string `json:"note" dc:"备注"`
	Status       string `json:"status" dc:"账户状态"`
	CreatedAt    string `json:"created_at" dc:"创建时间"`
	UpdatedAt    string `json:"updated_at" dc:"更新时间"`
}

// CreateAccountReq 添加交易所账户请求
type CreateAccountReq struct {
	g.Meta       `path:"/exchanges/accounts" method:"post" summary:"添加交易所账户" tags:"交易所"`
	ExchangeName string `json:"exchange_name" v:"required|in:binance,okx,coinbase" dc:"交易所名称"`
	ApiKey       string `json:"api_key" v:"required" dc:"API Key"`
	ApiSecret    string `json:"api_secret" v:"required" dc:"API Secret"`
	Passphrase   string `json:"passphrase" dc:"API Passphrase（OKX 必填）"`
	Label        string `json:"label" dc:"账户标签"`
	Note         string `json:"note" dc:"备注"`
}

// CreateAccountRes 添加交易所账户响应
type CreateAccountRes struct {
	Account *AccountItem `json:"account" dc:"新建的账户信息"`
}

// GetAccountReq 获取单个交易所账户请求
type GetAccountReq struct {
	g.Meta `path:"/exchanges/accounts/{id}" method:"get" summary:"获取单个交易所账户" tags:"交易所"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"账户 ID"`
}

// GetAccountRes 获取单个交易所账户响应
type GetAccountRes struct {
	Account *AccountItem `json:"account" dc:"账户信息"`
}

// UpdateAccountReq 更新交易所账户请求
type UpdateAccountReq struct {
	g.Meta     `path:"/exchanges/accounts/{id}" method:"put" summary:"更新交易所账户" tags:"交易所"`
	Id         uint   `json:"id" in:"path" v:"required|min:1" dc:"账户 ID"`
	ApiKey     string `json:"api_key" dc:"API Key（可选更新）"`
	ApiSecret  string `json:"api_secret" dc:"API Secret（可选更新）"`
	Passphrase string `json:"passphrase" dc:"API Passphrase（可选更新）"`
	Label      string `json:"label" dc:"账户标签"`
	Note       string `json:"note" dc:"备注"`
}

// UpdateAccountRes 更新交易所账户响应
type UpdateAccountRes struct {
	Account *AccountItem `json:"account" dc:"更新后的账户信息"`
}

// DeleteAccountReq 删除交易所账户请求
type DeleteAccountReq struct {
	g.Meta `path:"/exchanges/accounts/{id}" method:"delete" summary:"删除交易所账户" tags:"交易所"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"账户 ID"`
}

// DeleteAccountRes 删除交易所账户响应
type DeleteAccountRes struct{}

// TestConnectionReq 测试交易所连接请求
type TestConnectionReq struct {
	g.Meta `path:"/exchanges/accounts/{id}/test" method:"post" summary:"测试交易所连接" tags:"交易所"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"账户 ID"`
}

// TestConnectionRes 测试交易所连接响应
type TestConnectionRes struct {
	Success bool   `json:"success" dc:"连接是否成功"`
	Message string `json:"message" dc:"测试结果消息"`
}

// GetBalancesReq 获取交易所账户余额请求
type GetBalancesReq struct {
	g.Meta `path:"/exchanges/accounts/{id}/balances" method:"get" summary:"获取交易所账户余额" tags:"交易所"`
	Id     uint `json:"id" in:"path" v:"required|min:1" dc:"账户 ID"`
}

// BalanceItem 余额条目
type BalanceItem struct {
	Symbol    string  `json:"symbol" dc:"币种符号"`
	Free      float64 `json:"free" dc:"可用余额"`
	Locked    float64 `json:"locked" dc:"冻结余额"`
	Total     float64 `json:"total" dc:"总余额"`
	ValueUSD  float64 `json:"value_usd" dc:"USD 估值"`
}

// GetBalancesRes 获取交易所账户余额响应
type GetBalancesRes struct {
	Balances    []BalanceItem `json:"balances" dc:"余额列表"`
	TotalValue  float64       `json:"total_value" dc:"总价值（USD）"`
}

// SyncAccountReq 同步交易所账户请求
type SyncAccountReq struct {
	g.Meta `path:"/exchanges/accounts/{id}/sync" method:"post" summary:"同步交易所账户" tags:"交易所"`
	Id     int `json:"id" v:"required" in:"path" dc:"账户 ID"`
}

// SyncAccountRes 同步交易所账户响应
type SyncAccountRes struct {
	Message string `json:"message"`
}
