// Package model 汇率模块 - 数据传输对象 (DTO)
package model

// GetRatesInput 获取汇率输入
type GetRatesInput struct {
	Currencies []string // 要查询的货币列表
}

// GetPricesInput 获取价格输入
type GetPricesInput struct {
	Symbols []string // 要查询的币种列表
}

// RefreshRatesInput 刷新汇率输入
type RefreshRatesInput struct {
	Force bool // 是否强制刷新（忽略缓存）
}

// 默认查询的加密货币列表
var DefaultCryptoSymbols = []string{"BTC", "ETH", "USDT", "USDC", "BNB", "SOL"}

// 缓存过期时间（秒）
const CacheTTLSeconds = 300 // 5 分钟
