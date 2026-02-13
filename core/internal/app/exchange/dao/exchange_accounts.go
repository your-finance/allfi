// Package dao 交易所模块数据访问层
// 封装 exchange_accounts 表的 DAO 操作
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// ExchangeAccounts 交易所账户表全局访问对象（引用全局 DAO）
var ExchangeAccounts = &globalDao.ExchangeAccounts
