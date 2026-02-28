// Package gas Gas 优化 API 定义
package gas

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetCurrentReq 获取当前 Gas 价格请求
type GetCurrentReq struct {
	g.Meta `path:"/gas/current" method:"get" tags:"Gas优化" summary:"获取当前 Gas 价格"`
}

// GetCurrentRes 获取当前 Gas 价格响应
type GetCurrentRes struct {
	Chains []ChainGasPrice `json:"chains" dc:"各链 Gas 价格列表"`
}

// ChainGasPrice 单链 Gas 价格
type ChainGasPrice struct {
	Chain    string  `json:"chain"    dc:"链名称"`
	Name     string  `json:"name"     dc:"链显示名称"`
	Low      float64 `json:"low"      dc:"低速 Gas 价格（Gwei）"`
	Standard float64 `json:"standard" dc:"标准 Gas 价格（Gwei）"`
	Fast     float64 `json:"fast"     dc:"快速 Gas 价格（Gwei）"`
	Instant  float64 `json:"instant"  dc:"极速 Gas 价格（Gwei）"`
	BaseFee  float64 `json:"base_fee" dc:"基础费用（Gwei）"`
	Unit     string  `json:"unit"     dc:"单位"`
	Level    string  `json:"level"    dc:"拥堵等级"`
}

// GetHistoryReq 获取 Gas 价格历史请求
type GetHistoryReq struct {
	g.Meta `path:"/gas/history" method:"get" tags:"Gas优化" summary:"获取 Gas 价格历史"`
	Chain  string `json:"chain" v:"required|in:ethereum,bsc,polygon,arbitrum,optimism,base" dc:"链名称"`
	Hours  int    `json:"hours" v:"required|min:1|max:168" dc:"查询小时数（1-168）"`
}

// GetHistoryRes 获取 Gas 价格历史响应
type GetHistoryRes struct {
	Chain   string              `json:"chain"   dc:"链名称"`
	Hours   int                 `json:"hours"   dc:"查询小时数"`
	History []GasPriceHistoryVO `json:"history" dc:"历史记录"`
}

// GasPriceHistoryVO Gas 价格历史记录视图对象
type GasPriceHistoryVO struct {
	Timestamp  int64   `json:"timestamp"  dc:"时间戳（秒）"`
	Low        float64 `json:"low"        dc:"低速 Gas 价格"`
	Standard   float64 `json:"standard"   dc:"标准 Gas 价格"`
	Fast       float64 `json:"fast"       dc:"快速 Gas 价格"`
	Instant    float64 `json:"instant"    dc:"极速 Gas 价格"`
	BaseFee    float64 `json:"base_fee"   dc:"基础费用"`
	RecordedAt string  `json:"recorded_at" dc:"记录时间"`
}

// GetRecommendationReq 获取最佳交易时间推荐请求
type GetRecommendationReq struct {
	g.Meta `path:"/gas/recommendation" method:"get" tags:"Gas优化" summary:"获取最佳交易时间推荐"`
	Chain  string `json:"chain" v:"required|in:ethereum,bsc,polygon,arbitrum,optimism,base" dc:"链名称"`
}

// GetRecommendationRes 获取最佳交易时间推荐响应
type GetRecommendationRes struct {
	Chain            string  `json:"chain"             dc:"链名称"`
	RecommendedTime  string  `json:"recommended_time"  dc:"推荐交易时间段"`
	EstimatedSavings float64 `json:"estimated_savings" dc:"预计节省百分比"`
	Confidence       float64 `json:"confidence"        dc:"置信度（0-1）"`
	ValidUntil       string  `json:"valid_until"       dc:"有效期至"`
	CurrentGasPrice  float64 `json:"current_gas_price" dc:"当前 Gas 价格"`
	RecommendedPrice float64 `json:"recommended_price" dc:"推荐时段预计价格"`
}

// GetForecastReq 获取 Gas 价格预测请求
type GetForecastReq struct {
	g.Meta `path:"/gas/forecast" method:"get" tags:"Gas优化" summary:"获取 Gas 价格预测"`
	Chain  string `json:"chain" v:"required|in:ethereum,bsc,polygon,arbitrum,optimism,base" dc:"链名称"`
	Hours  int    `json:"hours" v:"required|min:1|max:24" dc:"预测小时数（1-24）"`
}

// GetForecastRes 获取 Gas 价格预测响应
type GetForecastRes struct {
	Chain    string            `json:"chain"    dc:"链名称"`
	Hours    int               `json:"hours"    dc:"预测小时数"`
	Forecast []GasForecastVO   `json:"forecast" dc:"预测数据"`
	Trend    string            `json:"trend"    dc:"趋势（上升/下降/平稳）"`
	Confidence float64         `json:"confidence" dc:"预测置信度"`
}

// GasForecastVO Gas 价格预测视图对象
type GasForecastVO struct {
	Timestamp int64   `json:"timestamp" dc:"时间戳（秒）"`
	Time      string  `json:"time"      dc:"时间"`
	Predicted float64 `json:"predicted" dc:"预测价格"`
	Lower     float64 `json:"lower"     dc:"预测下限"`
	Upper     float64 `json:"upper"     dc:"预测上限"`
}
