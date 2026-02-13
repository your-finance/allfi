// Package market 市场数据 API 定义
// 提供 Gas 价格等市场数据查询接口
package market

import "github.com/gogf/gf/v2/frame/g"

// GetGasReq 获取 Gas 价格请求
type GetGasReq struct {
	g.Meta `path:"/market/gas" method:"get" summary:"获取多链 Gas 价格" tags:"市场数据"`
}

// GasPrice Gas 价格条目
type GasPrice struct {
	Chain    string  `json:"chain" dc:"链名称（Ethereum/BSC/Polygon）"`
	Low      float64 `json:"low" dc:"低速价格（Gwei）"`
	Standard float64 `json:"standard" dc:"标准价格（Gwei）"`
	Fast     float64 `json:"fast" dc:"快速价格（Gwei）"`
	Instant  float64 `json:"instant" dc:"极速价格（Gwei）"`
	BaseFee  float64 `json:"base_fee" dc:"基础费用（Gwei）"`
	Level    string  `json:"level" dc:"拥堵等级（low/medium/high）"`
}

// GetGasRes 获取 Gas 价格响应
type GetGasRes struct {
	Prices    []GasPrice `json:"prices" dc:"各链 Gas 价格"`
	UpdatedAt int64      `json:"updated_at" dc:"更新时间（Unix 秒）"`
}
