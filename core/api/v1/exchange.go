// Package v1 交易所 API 定义
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/model/entity"
)

// AddAccountReq 添加交易所账户请求
type AddAccountReq struct {
	g.Meta        `path:"/exchanges/accounts" method:"post" summary:"添加交易所账户" tags:"交易所"`
	ExchangeName  string `json:"exchange_name" v:"required|in:binance,okx,coinbase" dc:"交易所名称"`
	ApiKey        string `json:"api_key" v:"required|length:16,128" dc:"API Key"`
	ApiSecret     string `json:"api_secret" v:"required|length:16,256" dc:"API Secret"`
	ApiPassphrase string `json:"api_passphrase" v:"max-length:128" dc:"API Passphrase(可选)"`
	Label         string `json:"label" v:"max-length:100" dc:"账户标签"`
}

// AddAccountRes 添加交易所账户响应
type AddAccountRes struct {
	Account *entity.ExchangeAccount `json:"account" dc:"账户信息"`
}

// ListAccountsReq 获取账户列表请求
type ListAccountsReq struct {
	g.Meta `path:"/exchanges/accounts" method:"get" summary:"获取交易所账户列表" tags:"交易所"`
}

// ListAccountsRes 获取账户列表响应
type ListAccountsRes struct {
	Accounts []*entity.ExchangeAccount `json:"accounts" dc:"账户列表"`
}

// DeleteAccountReq 删除账户请求
type DeleteAccountReq struct {
	g.Meta    `path:"/exchanges/accounts/:id" method:"delete" summary:"删除交易所账户" tags:"交易所"`
	AccountId uint `json:"id" v:"required|min:1" dc:"账户ID"`
}

// DeleteAccountRes 删除账户响应
type DeleteAccountRes struct{}
