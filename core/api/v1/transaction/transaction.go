// Package transaction 交易记录 API 定义
// 提供交易记录查询、同步、统计，以及同步设置接口
package transaction

import "github.com/gogf/gf/v2/frame/g"

// ListReq 获取交易记录列表请求
type ListReq struct {
	g.Meta   `path:"/transactions" method:"get" summary:"获取交易记录列表" tags:"交易记录"`
	Page     int    `json:"page" in:"query" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" in:"query" d:"20" dc:"每页数量"`
	Source   string `json:"source" in:"query" dc:"来源筛选（binance/okx/coinbase/ethereum 等）"`
	Type     string `json:"type" in:"query" dc:"交易类型筛选（buy/sell/transfer/deposit/withdraw）"`
	Start    string `json:"start" in:"query" dc:"开始日期（格式: 2024-01-01）"`
	End      string `json:"end" in:"query" dc:"结束日期（格式: 2024-12-31）"`
	Cursor   string `json:"cursor" in:"query" dc:"游标分页（上一页最后一条记录的时间戳，RFC3339 格式）"`
}

// TransactionItem 交易记录条目
type TransactionItem struct {
	ID        uint    `json:"id" dc:"交易 ID"`
	Source    string  `json:"source" dc:"来源"`
	TxType    string  `json:"tx_type" dc:"交易类型"`
	Symbol    string  `json:"symbol" dc:"交易币种"`
	Amount    float64 `json:"amount" dc:"交易数量"`
	Price     float64 `json:"price" dc:"成交价格"`
	Total     float64 `json:"total" dc:"总金额"`
	Fee       float64 `json:"fee" dc:"手续费"`
	FeeCoin   string  `json:"fee_coin" dc:"手续费币种"`
	TxHash    string  `json:"tx_hash" dc:"交易哈希（链上交易）"`
	Timestamp string  `json:"timestamp" dc:"交易时间"`
}

// ListRes 获取交易记录列表响应
type ListRes struct {
	Transactions []TransactionItem `json:"transactions" dc:"交易记录列表"`
	Total        int64             `json:"total" dc:"总记录数"`
	Page         int               `json:"page" dc:"当前页码"`
	PageSize     int               `json:"page_size" dc:"每页数量"`
}

// SyncReq 同步交易记录请求
type SyncReq struct {
	g.Meta `path:"/transactions/sync" method:"post" summary:"触发交易记录同步" tags:"交易记录"`
}

// SyncRes 同步交易记录响应
type SyncRes struct {
	Message    string `json:"message" dc:"同步结果消息"`
	SyncedCount int   `json:"synced_count" dc:"同步的记录数"`
}

// GetStatsReq 获取交易统计请求
type GetStatsReq struct {
	g.Meta `path:"/transactions/stats" method:"get" summary:"获取交易统计" tags:"交易记录"`
}

// GetStatsRes 获取交易统计响应
type GetStatsRes struct {
	TotalTransactions int                `json:"total_transactions" dc:"总交易次数"`
	TotalVolume       float64            `json:"total_volume" dc:"总交易量（USD）"`
	TotalFees         float64            `json:"total_fees" dc:"总手续费（USD）"`
	ByType            map[string]int     `json:"by_type" dc:"按类型统计"`
	BySource          map[string]int     `json:"by_source" dc:"按来源统计"`
}

// GetSyncSettingsReq 获取交易同步设置请求
type GetSyncSettingsReq struct {
	g.Meta `path:"/settings/tx-sync" method:"get" summary:"获取交易同步设置" tags:"交易记录"`
}

// SyncSettingsItem 同步设置条目
type SyncSettingsItem struct {
	AutoSync     bool   `json:"auto_sync" dc:"是否自动同步"`
	SyncInterval int    `json:"sync_interval" dc:"同步间隔（分钟）"`
	LastSyncAt   string `json:"last_sync_at" dc:"上次同步时间"`
}

// GetSyncSettingsRes 获取交易同步设置响应
type GetSyncSettingsRes struct {
	Settings *SyncSettingsItem `json:"settings" dc:"同步设置"`
}

// UpdateSyncSettingsReq 更新交易同步设置请求
type UpdateSyncSettingsReq struct {
	g.Meta       `path:"/settings/tx-sync" method:"put" summary:"更新交易同步设置" tags:"交易记录"`
	AutoSync     *bool `json:"auto_sync" dc:"是否自动同步"`
	SyncInterval *int  `json:"sync_interval" dc:"同步间隔（分钟）"`
}

// UpdateSyncSettingsRes 更新交易同步设置响应
type UpdateSyncSettingsRes struct {
	Settings *SyncSettingsItem `json:"settings" dc:"更新后的同步设置"`
}
