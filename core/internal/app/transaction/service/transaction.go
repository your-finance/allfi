// Package service 交易记录模块 - 服务接口定义
package service

import (
	"context"

	transactionApi "your-finance/allfi/api/v1/transaction"
)

// ITransaction 交易记录服务接口
type ITransaction interface {
	// List 获取交易记录列表（分页 + 筛选）
	List(ctx context.Context, page, pageSize int, source, txType, start, end, cursor string) (*transactionApi.ListRes, error)

	// Sync 触发交易记录同步
	Sync(ctx context.Context) (*transactionApi.SyncRes, error)

	// GetStats 获取交易统计
	GetStats(ctx context.Context) (*transactionApi.GetStatsRes, error)

	// GetSyncSettings 获取同步设置
	GetSyncSettings(ctx context.Context) (*transactionApi.GetSyncSettingsRes, error)

	// UpdateSyncSettings 更新同步设置
	UpdateSyncSettings(ctx context.Context, autoSync *bool, syncInterval *int) (*transactionApi.UpdateSyncSettingsRes, error)
}

var localTransaction ITransaction

// Transaction 获取交易记录服务实例
func Transaction() ITransaction {
	if localTransaction == nil {
		panic("ITransaction 服务未注册，请检查 logic/transaction 包的 init 函数")
	}
	return localTransaction
}

// RegisterTransaction 注册交易记录服务实现
func RegisterTransaction(i ITransaction) {
	localTransaction = i
}
