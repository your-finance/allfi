// Package controller 交易所模块 - 路由和控制器
// 绑定交易所 API 请求到对应的服务方法
package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	exchangeApi "your-finance/allfi/api/v1/exchange"
	"your-finance/allfi/internal/app/exchange/service"
	"your-finance/allfi/internal/consts"
)

// Controller 交易所控制器
type Controller struct{}

// Register 注册交易所模块路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&Controller{})
}

// ListSupportedExchanges 获取支持的交易所列表
func (c *Controller) ListSupportedExchanges(ctx context.Context, req *exchangeApi.ListSupportedExchangesReq) (res *exchangeApi.ListSupportedExchangesRes, err error) {
	exchanges, err := service.Exchange().ListSupportedExchanges(ctx)
	if err != nil {
		return nil, err
	}

	return &exchangeApi.ListSupportedExchangesRes{
		Exchanges: exchanges,
	}, nil
}

// ListAccounts 获取交易所账户列表
func (c *Controller) ListAccounts(ctx context.Context, req *exchangeApi.ListAccountsReq) (res *exchangeApi.ListAccountsRes, err error) {
	userID := consts.GetUserID(ctx)

	accounts, err := service.Exchange().ListAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &exchangeApi.ListAccountsRes{
		Accounts: accounts,
	}, nil
}

// CreateAccount 添加交易所账户
func (c *Controller) CreateAccount(ctx context.Context, req *exchangeApi.CreateAccountReq) (res *exchangeApi.CreateAccountRes, err error) {
	userID := consts.GetUserID(ctx)

	account, err := service.Exchange().CreateAccount(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &exchangeApi.CreateAccountRes{
		Account: account,
	}, nil
}

// GetAccount 获取单个交易所账户
func (c *Controller) GetAccount(ctx context.Context, req *exchangeApi.GetAccountReq) (res *exchangeApi.GetAccountRes, err error) {
	account, err := service.Exchange().GetAccount(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &exchangeApi.GetAccountRes{
		Account: account,
	}, nil
}

// UpdateAccount 更新交易所账户
func (c *Controller) UpdateAccount(ctx context.Context, req *exchangeApi.UpdateAccountReq) (res *exchangeApi.UpdateAccountRes, err error) {
	account, err := service.Exchange().UpdateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return &exchangeApi.UpdateAccountRes{
		Account: account,
	}, nil
}

// DeleteAccount 删除交易所账户
func (c *Controller) DeleteAccount(ctx context.Context, req *exchangeApi.DeleteAccountReq) (res *exchangeApi.DeleteAccountRes, err error) {
	err = service.Exchange().DeleteAccount(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &exchangeApi.DeleteAccountRes{}, nil
}

// TestConnection 测试交易所连接
func (c *Controller) TestConnection(ctx context.Context, req *exchangeApi.TestConnectionReq) (res *exchangeApi.TestConnectionRes, err error) {
	success, message, err := service.Exchange().TestConnection(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &exchangeApi.TestConnectionRes{
		Success: success,
		Message: message,
	}, nil
}

// GetBalances 获取交易所账户余额
func (c *Controller) GetBalances(ctx context.Context, req *exchangeApi.GetBalancesReq) (res *exchangeApi.GetBalancesRes, err error) {
	balances, totalValue, err := service.Exchange().GetBalances(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &exchangeApi.GetBalancesRes{
		Balances:   balances,
		TotalValue: totalValue,
	}, nil
}

// SyncAccount 同步交易所账户余额到 asset_details
func (c *Controller) SyncAccount(ctx context.Context, req *exchangeApi.SyncAccountReq) (res *exchangeApi.SyncAccountRes, err error) {
	err = service.Exchange().SyncAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &exchangeApi.SyncAccountRes{Message: "同步成功"}, nil
}
