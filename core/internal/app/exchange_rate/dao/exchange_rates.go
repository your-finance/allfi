// Package dao 汇率模块数据访问层
// 封装 exchange_rates 表的 DAO 操作
package dao

import (
	globalDao "your-finance/allfi/internal/dao"
)

// ExchangeRates 汇率表全局访问对象（引用全局 DAO）
var ExchangeRates = &globalDao.ExchangeRates
