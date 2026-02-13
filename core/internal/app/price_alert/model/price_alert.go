// Package model 价格预警模块 - 业务数据传输对象
// 定义价格预警模块内部使用的 DTO 和常量
package model

// 预警条件常量
const (
	ConditionAbove = "above" // 价格高于目标价
	ConditionBelow = "below" // 价格低于目标价
)

// IsValidCondition 验证预警条件是否有效
func IsValidCondition(condition string) bool {
	return condition == ConditionAbove || condition == ConditionBelow
}
