// Package model 交易记录模块 - 数据传输对象
package model

// 默认同步设置
const (
	DefaultSyncInterval = 60    // 默认同步间隔（分钟）
	DefaultAutoSync     = false // 默认不启用自动同步
	DefaultLookbackDays = 90    // 默认回溯天数
)

// 系统配置键名
const (
	ConfigKeyAutoSync     = "tx_auto_sync"
	ConfigKeySyncInterval = "tx_sync_interval"
	ConfigKeyLastSyncAt   = "tx_last_sync_at"
	ConfigKeyLookbackDays = "tx_lookback_days"
)
