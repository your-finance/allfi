// Package provider 汇率数据提供者
package provider

import (
	"context"
	"time"
)

// Provider 汇率数据提供者接口
type Provider interface {
	// Name 提供者名称
	Name() string

	// Priority 优先级（数字越小优先级越高）
	Priority() int

	// FetchRates 获取所有支持的汇率
	FetchRates(ctx context.Context) (map[string]*RateInfo, error)

	// FetchRate 获取单个币种汇率
	FetchRate(ctx context.Context, symbol string) (*RateInfo, error)

	// IsHealthy 健康检查
	IsHealthy(ctx context.Context) bool

	// SupportedSymbols 支持的币种列表
	SupportedSymbols() []string

	// SupportsSymbol 是否支持指定币种
	SupportsSymbol(symbol string) bool
}

// RateInfo 汇率信息
type RateInfo struct {
	Symbol      string    `json:"symbol"`      // 币种符号
	PriceUSD    float64   `json:"priceUsd"`    // 对USD价格
	Change24H   float64   `json:"change24h"`   // 24小时涨跌幅
	Volume24H   float64   `json:"volume24h"`   // 24小时交易量
	MarketCap   float64   `json:"marketCap"`   // 市值
	LastUpdated time.Time `json:"lastUpdated"` // 更新时间
	Source      string    `json:"source"`      // 数据来源
}

// ProviderError 数据源错误
type ProviderError struct {
	Provider string
	Err      error
}

func (e *ProviderError) Error() string {
	return e.Provider + ": " + e.Err.Error()
}

// NewProviderError 创建数据源错误
func NewProviderError(provider string, err error) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Err:      err,
	}
}
