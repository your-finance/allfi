// Package controller 成就系统模块路由注册
// 使用子目录 API 包定义的请求/响应类型
package controller

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	achievementApi "your-finance/allfi/api/v1/achievement"
	"your-finance/allfi/internal/app/achievement/service"
	"your-finance/allfi/internal/consts"
)

// AchievementController 成就系统控制器
type AchievementController struct{}

// List 获取成就列表
//
// 对应路由: GET /achievements
// 返回所有成就及其解锁状态
func (c *AchievementController) List(ctx context.Context, req *achievementApi.ListReq) (res *achievementApi.ListRes, err error) {
	// 从上下文中获取当前用户 ID
	userID := uint(consts.GetUserID(ctx))

	// 调用 Service 层
	statuses, err := service.Achievement().GetAchievements(ctx, userID)
	if err != nil {
		return nil, gerror.Wrap(err, "获取成就列表失败")
	}

	// 将业务 DTO 转换为 API 响应
	var achievements []achievementApi.AchievementItem
	unlockedCount := 0

	for _, status := range statuses {
		achievements = append(achievements, achievementApi.AchievementItem{
			ID:          status.ID,
			Name:        status.Name,
			Description: status.Description,
			Icon:        status.Icon,
			Category:    status.Category,
			IsUnlocked:  status.IsUnlocked,
			UnlockedAt:  status.UnlockedAt,
			Progress:    status.Progress,
		})
		if status.IsUnlocked {
			unlockedCount++
		}
	}

	res = &achievementApi.ListRes{
		Achievements:  achievements,
		TotalCount:    len(achievements),
		UnlockedCount: unlockedCount,
	}

	return res, nil
}

// Check 检查并解锁成就
//
// 对应路由: POST /achievements/check
// 检查当前用户是否满足新成就解锁条件
func (c *AchievementController) Check(ctx context.Context, req *achievementApi.CheckReq) (res *achievementApi.CheckRes, err error) {
	// 从上下文中获取当前用户 ID
	userID := uint(consts.GetUserID(ctx))

	// 调用 Service 层
	newlyUnlocked, err := service.Achievement().CheckAndUnlock(ctx, userID)
	if err != nil {
		return nil, gerror.Wrap(err, "检查成就解锁失败")
	}

	// 转换为 API 响应
	var items []achievementApi.AchievementItem
	for _, unlocked := range newlyUnlocked {
		items = append(items, achievementApi.AchievementItem{
			ID:          unlocked.ID,
			Name:        unlocked.Name,
			Description: unlocked.Description,
			Icon:        unlocked.Icon,
			Category:    unlocked.Category,
			IsUnlocked:  true,
			Progress:    100,
		})
	}

	res = &achievementApi.CheckRes{
		NewlyUnlocked: items,
	}

	return res, nil
}

// Register 注册成就系统路由
// 使用 group.Bind 自动绑定控制器方法到路由
func Register(group *ghttp.RouterGroup) {
	group.Bind(&AchievementController{})
}
