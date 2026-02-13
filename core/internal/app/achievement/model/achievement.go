// Package model 成就系统模块业务 DTO
// 定义成就相关的数据传输对象
package model

// AchievementStatus 成就状态（含解锁信息）
type AchievementStatus struct {
	ID          string `json:"id"`          // 成就 ID
	Name        string `json:"name"`        // 成就名称
	Description string `json:"description"` // 成就描述
	Icon        string `json:"icon"`        // 图标标识
	Category    string `json:"category"`    // 分类（milestone/persistence/investment）
	IsUnlocked  bool   `json:"is_unlocked"` // 是否已解锁
	UnlockedAt  string `json:"unlocked_at"` // 解锁时间（ISO 格式）
	Progress    int    `json:"progress"`    // 当前进度（百分比，0-100）
}

// UnlockedAchievement 新解锁的成就
type UnlockedAchievement struct {
	ID          string `json:"id"`          // 成就 ID
	Name        string `json:"name"`        // 成就名称
	Description string `json:"description"` // 成就描述
	Icon        string `json:"icon"`        // 图标标识
	Category    string `json:"category"`    // 分类
}

// 成就类别常量
const (
	CategoryMilestone   = "milestone"   // 里程碑
	CategoryPersistence = "persistence" // 坚持
	CategoryInvestment  = "investment"  // 投资
)
