// Package dao 系统管理模块数据访问层
// 封装 system_config 表的 DAO 操作（引用全局 DAO）
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// SystemConfig 系统配置表访问对象（引用全局 DAO）
var SystemConfig = &globalDao.SystemConfig
