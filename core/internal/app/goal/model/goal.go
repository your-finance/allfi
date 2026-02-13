// Package model 目标追踪模块 - 业务数据传输对象
// 定义目标追踪模块内部使用的 DTO 和常量
package model

// 目标类型常量
const (
	GoalTypeAssetValue    = "asset_value"    // 资产总值
	GoalTypeHoldingAmount = "holding_amount" // 持仓数量
	GoalTypeReturnRate    = "return_rate"    // 收益率
)

// IsValidGoalType 验证目标类型是否有效
func IsValidGoalType(goalType string) bool {
	switch goalType {
	case GoalTypeAssetValue, GoalTypeHoldingAmount, GoalTypeReturnRate:
		return true
	default:
		return false
	}
}
