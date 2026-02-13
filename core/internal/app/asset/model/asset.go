// Package model 资产模块 - 数据传输对象 (DTO)
package model

// GetSummaryInput 获取资产概览输入
type GetSummaryInput struct {
	Currency string // 计价货币（USD/BTC/ETH/CNY）
}

// GetDetailsInput 获取资产明细输入
type GetDetailsInput struct {
	SourceType string // 来源类型筛选（cex/blockchain/manual）
	Currency   string // 计价货币
}

// GetHistoryInput 获取资产历史输入
type GetHistoryInput struct {
	Days     int    // 查询天数
	Currency string // 计价货币
}
