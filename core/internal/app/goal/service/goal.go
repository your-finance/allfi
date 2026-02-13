// Package service 目标追踪模块 - 服务接口定义
// 定义投资目标的增删改查和进度计算
package service

import (
	"context"

	goalApi "your-finance/allfi/api/v1/goal"
)

// IGoal 目标追踪服务接口
type IGoal interface {
	// GetGoals 获取目标列表（带进度百分比）
	GetGoals(ctx context.Context) ([]goalApi.GoalItem, error)

	// CreateGoal 创建目标
	CreateGoal(ctx context.Context, req *goalApi.CreateReq) (*goalApi.GoalItem, error)

	// UpdateGoal 更新目标
	UpdateGoal(ctx context.Context, req *goalApi.UpdateReq) (*goalApi.GoalItem, error)

	// DeleteGoal 删除目标
	DeleteGoal(ctx context.Context, goalID int) error
}

var localGoal IGoal

// Goal 获取目标追踪服务实例
func Goal() IGoal {
	if localGoal == nil {
		panic("IGoal 服务未注册，请检查 logic/goal 包的 init 函数")
	}
	return localGoal
}

// RegisterGoal 注册目标追踪服务实现
// 由 logic 层在 init 函数中调用
func RegisterGoal(i IGoal) {
	localGoal = i
}
