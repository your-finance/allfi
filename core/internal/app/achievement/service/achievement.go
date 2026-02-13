// =================================================================================
// 成就系统服务接口定义
// 提供成就列表查询和解锁检查能力
// =================================================================================

package service

import (
	"context"

	"your-finance/allfi/internal/app/achievement/model"
)

// IAchievement 成就系统服务接口
type IAchievement interface {
	// GetAchievements 获取成就列表（含解锁状态）
	// userID: 用户 ID
	// 返回所有成就定义及其解锁状态
	GetAchievements(ctx context.Context, userID uint) ([]model.AchievementStatus, error)

	// CheckAndUnlock 检查并解锁新成就
	// userID: 用户 ID
	// 返回本次新解锁的成就列表
	CheckAndUnlock(ctx context.Context, userID uint) ([]model.UnlockedAchievement, error)
}

// localAchievement 成就系统服务实例（延迟注入）
var localAchievement IAchievement

// Achievement 获取成就系统服务实例
// 如果服务未注册，会触发 panic
func Achievement() IAchievement {
	if localAchievement == nil {
		panic("IAchievement 服务未注册，请检查 logic/achievement 包的 init 函数")
	}
	return localAchievement
}

// RegisterAchievement 注册成就系统服务实现
// 由 logic 层在 init 函数中调用
func RegisterAchievement(i IAchievement) {
	localAchievement = i
}
