// Package dao 交易记录模块 DAO 封装
// 对全局 DAO 的模块级引用，供本模块和跨模块调用
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// UnifiedTransactions 统一交易记录表访问对象（引用全局 DAO）
var UnifiedTransactions = &globalDao.UnifiedTransactions
