// Package model 策略引擎模块 - 数据传输对象
package model

// TargetAllocation 目标配比
type TargetAllocation struct {
	Symbol     string  `json:"symbol"`     // 币种
	Percentage float64 `json:"percentage"` // 目标比例（0-100）
}

// RebalanceConfig 再平衡策略配置
type RebalanceConfig struct {
	Allocations []TargetAllocation `json:"allocations"` // 目标配比列表
	Threshold   float64            `json:"threshold"`   // 触发阈值（偏离百分比）
}

// 策略类型常量
const (
	StrategyTypeRebalance = "rebalance"
	StrategyTypeDCA       = "dca"
	StrategyTypeStopLimit = "stop_limit"
)
