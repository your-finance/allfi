// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package dao

import (
	"your-finance/allfi/internal/dao/internal"
)

// internalCrossChainTransactionDao 跨链交易数据访问对象内部实例
var internalCrossChainTransactionDao = internal.NewCrossChainTransactionDao()

// CrossChainTransaction 跨链交易数据访问对象
var CrossChainTransaction = internalCrossChainTransactionDao

// CrossChainTransactionColumns 跨链交易字段列表
var CrossChainTransactionColumns = internalCrossChainTransactionDao.Columns()
