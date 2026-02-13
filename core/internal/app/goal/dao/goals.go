// Package dao 目标追踪模块数据访问层
// 封装 goals 表的 DAO 操作
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// Goals 目标表全局访问对象（引用全局 DAO）
var Goals = &globalDao.Goals
