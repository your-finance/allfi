package v1

import "github.com/gogf/gf/v2/frame/g"

// GetExchangeRatesReq 获取汇率请求
type GetExchangeRatesReq struct {
	g.Meta  `path:"/exchange-rates" method:"get" summary:"获取实时汇率" tags:"汇率"`
	Symbols string `json:"symbols" v:"max-length:500" dc:"币种列表（逗号分隔），例如: BTC,ETH,USDT，为空则返回全部"`
}

// GetExchangeRatesRes 获取汇率响应
type GetExchangeRatesRes struct {
	Rates       map[string]float64 `json:"rates" dc:"汇率映射 {\"BTC\": 45000.5}"`
	Base        string             `json:"base" dc:"基准货币（USD）"`
	LastUpdated int64              `json:"last_updated" dc:"最后更新时间（毫秒时间戳）"`
	Source      string             `json:"source" dc:"数据来源"`
	IsCached    bool               `json:"is_cached" dc:"是否缓存数据"`
	Warning     string             `json:"warning,omitempty" dc:"警告信息（如部分数据源失败）"`
}

// ConvertCurrencyReq 货币转换请求
type ConvertCurrencyReq struct {
	g.Meta `path:"/exchange-rates/convert" method:"get" summary:"货币转换" tags:"汇率"`
	Amount float64 `json:"amount" v:"required|min:0" dc:"金额"`
	From   string  `json:"from" v:"required|length:2,10" dc:"源货币，例如: BTC"`
	To     string  `json:"to" v:"required|length:2,10" dc:"目标货币，例如: USDT"`
}

// ConvertCurrencyRes 货币转换响应
type ConvertCurrencyRes struct {
	Amount      float64 `json:"amount" dc:"原始金额"`
	From        string  `json:"from" dc:"源货币"`
	To          string  `json:"to" dc:"目标货币"`
	Result      float64 `json:"result" dc:"转换结果"`
	Rate        float64 `json:"rate" dc:"使用的汇率"`
	LastUpdated int64   `json:"last_updated" dc:"汇率更新时间（毫秒时间戳）"`
	Calculation string  `json:"calculation" dc:"计算公式"`
	Warning     string  `json:"warning,omitempty" dc:"警告信息"`
}
