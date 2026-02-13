// Package service 钱包模块 - 服务接口定义
// 定义钱包地址管理的所有服务方法
package service

import (
	"context"

	walletApi "your-finance/allfi/api/v1/wallet"
)

// IWallet 钱包服务接口
type IWallet interface {
	// ListWallets 获取钱包地址列表
	ListWallets(ctx context.Context, userID int) ([]walletApi.AddressItem, error)

	// CreateWallet 添加钱包地址
	CreateWallet(ctx context.Context, userID int, req *walletApi.CreateAddressReq) (*walletApi.AddressItem, error)

	// GetWallet 获取单个钱包地址
	GetWallet(ctx context.Context, walletID int) (*walletApi.AddressItem, error)

	// UpdateWallet 更新钱包地址信息
	UpdateWallet(ctx context.Context, req *walletApi.UpdateAddressReq) (*walletApi.AddressItem, error)

	// DeleteWallet 软删除钱包地址
	DeleteWallet(ctx context.Context, walletID int) error

	// GetBalances 获取钱包余额
	GetBalances(ctx context.Context, walletID int) (float64, map[string]float64, float64, error)

	// BatchImport 批量导入钱包地址
	BatchImport(ctx context.Context, userID int, req *walletApi.BatchImportReq) (int, int, error)

	// SyncAddress 同步钱包地址余额到 asset_details 表
	SyncAddress(ctx context.Context, walletID int) error
}

var localWallet IWallet

// Wallet 获取钱包服务实例
func Wallet() IWallet {
	if localWallet == nil {
		panic("IWallet 服务未注册，请检查 logic/wallet 包的 init 函数")
	}
	return localWallet
}

// RegisterWallet 注册钱包服务实现
// 由 logic 层在 init 函数中调用
func RegisterWallet(i IWallet) {
	localWallet = i
}
