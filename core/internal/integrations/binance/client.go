// Package binance Binance 交易所 API 客户端
// 使用官方 SDK 获取账户余额
package binance

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/adshao/go-binance/v2"
	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/utils"
)

// Client Binance 客户端
type Client struct {
	client    *binance.Client
	apiKey    string
	apiSecret string
}

// NewClient 创建 Binance 客户端
func NewClient(apiKey, apiSecret string) *Client {
	client := binance.NewClient(apiKey, apiSecret)
	return &Client{
		client:    client,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// GetName 获取交易所名称
func (c *Client) GetName() string {
	return "binance"
}

// TestConnection 测试 API 连接
func (c *Client) TestConnection(ctx context.Context) error {
	// 限流
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return err
	}

	// 尝试获取账户信息（需要有效的 API Key）
	_, err := c.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return fmt.Errorf("Binance API 连接失败: %v", err)
	}
	return nil
}

// GetBalances 获取账户余额
func (c *Client) GetBalances(ctx context.Context) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	// 获取账户信息
	account, err := c.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Binance 账户信息失败: %v", err)
	}

	var balances []integrations.Balance

	for _, balance := range account.Balances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)
		total := free + locked

		// 过滤掉余额为 0 的资产
		if total <= 0 {
			continue
		}

		balances = append(balances, integrations.Balance{
			Symbol: balance.Asset,
			Name:   balance.Asset,
			Free:   free,
			Locked: locked,
			Total:  total,
		})
	}

	return balances, nil
}

// GetSpotBalance 获取现货账户余额
func (c *Client) GetSpotBalance(ctx context.Context) ([]integrations.Balance, error) {
	return c.GetBalances(ctx)
}

// GetFuturesBalance 获取 U 本位合约账户余额
func (c *Client) GetFuturesBalance(ctx context.Context) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	// 使用期货客户端（需要传入正确的 apiSecret）
	futuresClient := binance.NewFuturesClient(c.apiKey, c.apiSecret)
	account, err := futuresClient.NewGetAccountService().Do(ctx)
	if err != nil {
		// 如果没有期货账户或 API 权限不足，返回空列表（优雅降级）
		return []integrations.Balance{}, nil
	}

	var balances []integrations.Balance

	for _, asset := range account.Assets {
		walletBalance, _ := strconv.ParseFloat(asset.WalletBalance, 64)
		unrealizedProfit, _ := strconv.ParseFloat(asset.UnrealizedProfit, 64)

		if walletBalance <= 0 {
			continue
		}

		balances = append(balances, integrations.Balance{
			Symbol:    asset.Asset,
			Name:      asset.Asset + " (合约)",
			Free:      walletBalance,
			Locked:    unrealizedProfit,
			Total:     walletBalance,
			AssetType: "futures",
		})
	}

	return balances, nil
}

// GetMarginBalance 获取杠杆账户余额
func (c *Client) GetMarginBalance(ctx context.Context) ([]integrations.Balance, error) {
	// 限流
	if err := utils.WaitForAPI(ctx, "binance"); err != nil {
		return nil, err
	}

	// 获取杠杆账户信息
	account, err := c.client.NewGetMarginAccountService().Do(ctx)
	if err != nil {
		// 如果没有杠杆账户或权限不足，返回空列表（优雅降级）
		return []integrations.Balance{}, nil
	}

	var balances []integrations.Balance

	for _, asset := range account.UserAssets {
		free, _ := strconv.ParseFloat(asset.Free, 64)
		locked, _ := strconv.ParseFloat(asset.Locked, 64)
		total := free + locked

		if total <= 0 {
			continue
		}

		balances = append(balances, integrations.Balance{
			Symbol:    asset.Asset,
			Name:      asset.Asset + " (杠杆)",
			Free:      free,
			Locked:    locked,
			Total:     total,
			AssetType: "margin",
		})
	}

	return balances, nil
}

// GetAllBalances 获取所有账户类型的综合余额（现货+合约+杠杆）
// 并发获取三种账户余额，合并同一币种的持仓
// 任何子账户获取失败都不会影响其他账户的数据返回
func (c *Client) GetAllBalances(ctx context.Context) ([]integrations.Balance, error) {
	var (
		allBalances []integrations.Balance
		mu          sync.Mutex
		wg          sync.WaitGroup
	)

	// 收集非致命错误用于日志记录
	var errs []error

	wg.Add(3)

	// 现货
	go func() {
		defer wg.Done()
		balances, err := c.GetBalances(ctx)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("现货余额获取失败: %v", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		allBalances = append(allBalances, balances...)
		mu.Unlock()
	}()

	// U 本位合约
	go func() {
		defer wg.Done()
		balances, err := c.GetFuturesBalance(ctx)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("合约余额获取失败: %v", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		allBalances = append(allBalances, balances...)
		mu.Unlock()
	}()

	// 杠杆
	go func() {
		defer wg.Done()
		balances, err := c.GetMarginBalance(ctx)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("杠杆余额获取失败: %v", err))
			mu.Unlock()
			return
		}
		mu.Lock()
		allBalances = append(allBalances, balances...)
		mu.Unlock()
	}()

	wg.Wait()

	// 如果所有账户都获取失败，返回错误
	if len(errs) == 3 {
		return nil, fmt.Errorf("获取 Binance 所有账户余额均失败: %v", errs[0])
	}

	return allBalances, nil
}

// MergeBalancesBySymbol 合并同一币种在不同账户类型中的余额
// 将现货、合约、杠杆中同一币种的余额合并为一条记录
func MergeBalancesBySymbol(balances []integrations.Balance) []integrations.Balance {
	merged := make(map[string]*integrations.Balance)

	for _, b := range balances {
		existing, ok := merged[b.Symbol]
		if ok {
			existing.Free += b.Free
			existing.Locked += b.Locked
			existing.Total += b.Total
		} else {
			copy := b
			copy.Name = b.Symbol // 合并后使用纯符号名称
			copy.AssetType = ""  // 合并后不标注单一类型
			merged[b.Symbol] = &copy
		}
	}

	result := make([]integrations.Balance, 0, len(merged))
	for _, b := range merged {
		result = append(result, *b)
	}

	return result
}

// 确保实现接口
var _ integrations.ExchangeClient = (*Client)(nil)
