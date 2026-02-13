// Package service 交易所模块 - 服务接口定义
// 定义交易所账户管理的所有服务方法
package service

import (
	"context"

	exchangeApi "your-finance/allfi/api/v1/exchange"
)

// IExchange 交易所服务接口
type IExchange interface {
	// ListAccounts 获取用户的交易所账户列表
	// 返回账户列表（已过滤敏感字段）
	ListAccounts(ctx context.Context, userID int) ([]exchangeApi.AccountItem, error)

	// CreateAccount 添加交易所账户
	// API 凭证使用 AES-256-GCM 加密存储
	CreateAccount(ctx context.Context, userID int, req *exchangeApi.CreateAccountReq) (*exchangeApi.AccountItem, error)

	// GetAccount 获取单个交易所账户详情
	GetAccount(ctx context.Context, accountID int) (*exchangeApi.AccountItem, error)

	// UpdateAccount 更新交易所账户信息
	// 支持更新 API 凭证、标签、备注
	UpdateAccount(ctx context.Context, req *exchangeApi.UpdateAccountReq) (*exchangeApi.AccountItem, error)

	// DeleteAccount 软删除交易所账户
	DeleteAccount(ctx context.Context, accountID int) error

	// TestConnection 测试交易所 API 连接
	TestConnection(ctx context.Context, accountID int) (bool, string, error)

	// GetBalances 获取交易所账户余额
	GetBalances(ctx context.Context, accountID int) ([]exchangeApi.BalanceItem, float64, error)

	// SyncAccount 同步交易所账户余额到 asset_details 表
	SyncAccount(ctx context.Context, accountID int) error
}

var localExchange IExchange

// Exchange 获取交易所服务实例
func Exchange() IExchange {
	if localExchange == nil {
		panic("IExchange 服务未注册，请检查 logic/exchange 包的 init 函数")
	}
	return localExchange
}

// RegisterExchange 注册交易所服务实现
// 由 logic 层在 init 函数中调用
func RegisterExchange(i IExchange) {
	localExchange = i
}
