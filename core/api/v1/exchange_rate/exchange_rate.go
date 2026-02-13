// Package exchange_rate 汇率 API 定义
// 提供实时汇率查询、价格查询、汇率刷新接口
package exchange_rate

import "github.com/gogf/gf/v2/frame/g"

// GetCurrentReq 获取实时汇率请求
type GetCurrentReq struct {
	g.Meta     `path:"/rates/current" method:"get" summary:"获取实时汇率" tags:"汇率"`
	Currencies string `json:"currencies" in:"query" dc:"币种列表（逗号分隔），例如: BTC,ETH,USD"`
}

// GetCurrentRes 获取实时汇率响应
type GetCurrentRes struct {
	Rates       map[string]float64 `json:"rates" dc:"汇率映射 {\"BTC\": 45000.5}"`
	Base        string             `json:"base" dc:"基准货币（USD）"`
	LastUpdated int64              `json:"last_updated" dc:"最后更新时间（毫秒时间戳）"`
	Source      string             `json:"source" dc:"数据来源"`
	IsCached    bool               `json:"is_cached" dc:"是否缓存数据"`
}

// GetPricesReq 获取币种价格请求
type GetPricesReq struct {
	g.Meta  `path:"/rates/prices" method:"get" summary:"获取币种价格" tags:"汇率"`
	Symbols string `json:"symbols" in:"query" dc:"币种列表（逗号分隔），例如: BTC,ETH"`
}

// PriceItem 价格条目
type PriceItem struct {
	Symbol      string  `json:"symbol" dc:"币种符号"`
	PriceUSD    float64 `json:"price_usd" dc:"USD 价格"`
	Change24h   float64 `json:"change_24h" dc:"24 小时涨跌幅"`
	LastUpdated int64   `json:"last_updated" dc:"最后更新时间"`
}

// GetPricesRes 获取币种价格响应
type GetPricesRes struct {
	Prices []PriceItem `json:"prices" dc:"价格列表"`
}

// RefreshReq 刷新汇率请求
type RefreshReq struct {
	g.Meta `path:"/rates/refresh" method:"post" summary:"强制刷新汇率缓存" tags:"汇率"`
}

// RefreshRes 刷新汇率响应
type RefreshRes struct {
	Message string `json:"message" dc:"刷新结果消息"`
}

// GetHistoryReq 历史汇率请求
type GetHistoryReq struct {
	g.Meta `path:"/rates/history" method:"get" summary:"获取历史汇率" tags:"汇率"`
	Base   string `json:"base" in:"query" dc:"基准货币" d:"USD"`
	Quote  string `json:"quote" in:"query" dc:"目标货币"`
	Days   int    `json:"days" in:"query" dc:"天数" d:"30"`
}

// RateHistoryItem 历史汇率项
type RateHistoryItem struct {
	Date  string  `json:"date" dc:"日期"`
	Rate  float64 `json:"rate" dc:"汇率"`
	Base  string  `json:"base" dc:"基准货币"`
	Quote string  `json:"quote" dc:"目标货币"`
}

// GetHistoryRes 历史汇率响应
type GetHistoryRes struct {
	History []RateHistoryItem `json:"history" dc:"历史汇率列表"`
	Base    string            `json:"base" dc:"基准货币"`
	Quote   string            `json:"quote" dc:"目标货币"`
	Days    int               `json:"days" dc:"天数"`
}
