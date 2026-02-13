// Package achievement 成就系统 API 定义
// 提供成就列表查询和成就检查接口
package achievement

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取成就列表请求
type ListReq struct {
	g.Meta `path:"/achievements" method:"get" summary:"获取成就列表" tags:"成就系统"`
}

// AchievementItem 成就条目
type AchievementItem struct {
	ID          string `json:"id" dc:"成就 ID"`
	Name        string `json:"name" dc:"成就名称"`
	Description string `json:"description" dc:"成就描述"`
	Icon        string `json:"icon" dc:"图标"`
	Category    string `json:"category" dc:"分类"`
	IsUnlocked  bool   `json:"is_unlocked" dc:"是否已解锁"`
	UnlockedAt  string `json:"unlocked_at" dc:"解锁时间"`
	Progress    int    `json:"progress" dc:"当前进度（百分比）"`
}

// ListRes 获取成就列表响应
type ListRes struct {
	Achievements []AchievementItem `json:"achievements" dc:"成就列表"`
	TotalCount   int               `json:"total_count" dc:"总成就数"`
	UnlockedCount int              `json:"unlocked_count" dc:"已解锁数量"`
}

// CheckReq 检查成就解锁请求
type CheckReq struct {
	g.Meta `path:"/achievements/check" method:"post" summary:"检查并解锁成就" tags:"成就系统"`
}

// CheckRes 检查成就解锁响应
type CheckRes struct {
	NewlyUnlocked []AchievementItem `json:"newly_unlocked" dc:"本次新解锁的成就"`
}
