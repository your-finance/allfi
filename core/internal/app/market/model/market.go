// Package model 市场数据模块业务 DTO
// 定义 Gas 价格等市场数据相关的数据传输对象
package model

// MultiChainGasResponse 多链 Gas 价格响应
type MultiChainGasResponse struct {
	Chains    []ChainGasPrice `json:"chains"`     // 各链 Gas 价格列表
	UpdatedAt int64           `json:"updated_at"` // 更新时间（Unix 秒）
}

// ChainGasPrice 单链 Gas 价格
type ChainGasPrice struct {
	Chain    string  `json:"chain"`    // 链标识（ethereum/bsc/polygon）
	Name     string  `json:"name"`     // 显示名称
	Low      float64 `json:"low"`      // 低速价格（Gwei）
	Standard float64 `json:"standard"` // 标准价格（Gwei）
	Fast     float64 `json:"fast"`     // 快速价格（Gwei）
	Instant  float64 `json:"instant"`  // 极速价格（Gwei）
	BaseFee  float64 `json:"base_fee"` // 基础费用（Gwei）
	Unit     string  `json:"unit"`     // 单位（Gwei）
	Level    string  `json:"level"`    // 拥堵等级（low/medium/high）
}

// determineLevel 根据标准 Gas 价格判断拥堵等级
func DetermineLevel(standardGwei float64) string {
	switch {
	case standardGwei <= 20:
		return "low"
	case standardGwei <= 60:
		return "medium"
	default:
		return "high"
	}
}
