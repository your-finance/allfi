// Package dao 成就系统模块数据访问层
// 封装 user_achievements 表的 DAO 操作
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// UserAchievements 用户成就表全局访问对象（引用全局 DAO）
var UserAchievements = &globalDao.UserAchievements
