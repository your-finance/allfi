// Package controller 钱包模块 - 路由和控制器
// 绑定钱包 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	walletApi "your-finance/allfi/api/v1/wallet"
	"your-finance/allfi/internal/app/wallet/service"
	"your-finance/allfi/internal/consts"
)

// Controller 钱包控制器
type Controller struct{}

// Register 注册钱包模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// ListAddresses 获取钱包地址列表
func (c *Controller) ListAddresses(ctx context.Context, req *walletApi.ListAddressesReq) (res *walletApi.ListAddressesRes, err error) {
	userID := consts.GetUserID(ctx)

	addresses, err := service.Wallet().ListWallets(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &walletApi.ListAddressesRes{
		Addresses: addresses,
	}, nil
}

// CreateAddress 添加钱包地址
func (c *Controller) CreateAddress(ctx context.Context, req *walletApi.CreateAddressReq) (res *walletApi.CreateAddressRes, err error) {
	userID := consts.GetUserID(ctx)

	address, err := service.Wallet().CreateWallet(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &walletApi.CreateAddressRes{
		Address: address,
	}, nil
}

// GetAddress 获取单个钱包地址
func (c *Controller) GetAddress(ctx context.Context, req *walletApi.GetAddressReq) (res *walletApi.GetAddressRes, err error) {
	address, err := service.Wallet().GetWallet(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &walletApi.GetAddressRes{
		Address: address,
	}, nil
}

// UpdateAddress 更新钱包地址
func (c *Controller) UpdateAddress(ctx context.Context, req *walletApi.UpdateAddressReq) (res *walletApi.UpdateAddressRes, err error) {
	address, err := service.Wallet().UpdateWallet(ctx, req)
	if err != nil {
		return nil, err
	}

	return &walletApi.UpdateAddressRes{
		Address: address,
	}, nil
}

// DeleteAddress 删除钱包地址
func (c *Controller) DeleteAddress(ctx context.Context, req *walletApi.DeleteAddressReq) (res *walletApi.DeleteAddressRes, err error) {
	err = service.Wallet().DeleteWallet(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &walletApi.DeleteAddressRes{}, nil
}

// GetAddressBalances 获取钱包余额
func (c *Controller) GetAddressBalances(ctx context.Context, req *walletApi.GetAddressBalancesReq) (res *walletApi.GetAddressBalancesRes, err error) {
	nativeBalance, tokenBalances, totalValueUSD, err := service.Wallet().GetBalances(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &walletApi.GetAddressBalancesRes{
		NativeBalance: nativeBalance,
		TokenBalances: tokenBalances,
		TotalValueUSD: totalValueUSD,
	}, nil
}

// BatchImport 批量导入钱包地址
func (c *Controller) BatchImport(ctx context.Context, req *walletApi.BatchImportReq) (res *walletApi.BatchImportRes, err error) {
	userID := consts.GetUserID(ctx)

	imported, failed, err := service.Wallet().BatchImport(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &walletApi.BatchImportRes{
		Imported: imported,
		Failed:   failed,
	}, nil
}

// SyncAddress 同步钱包地址余额到 asset_details
func (c *Controller) SyncAddress(ctx context.Context, req *walletApi.SyncAddressReq) (res *walletApi.SyncAddressRes, err error) {
	err = service.Wallet().SyncAddress(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &walletApi.SyncAddressRes{Message: "同步成功"}, nil
}
