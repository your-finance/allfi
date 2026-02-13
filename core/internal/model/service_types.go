package model

// AssetBalance 资产余额（来自交易所或区块链）
type AssetBalance struct {
	Source     string  `json:"source"`      // 来源：交易所名称或区块链名称
	SourceID   uint    `json:"source_id"`   // 来源 ID
	SourceType string  `json:"source_type"` // 来源类型：cex/blockchain
	Symbol     string  `json:"symbol"`      // 币种符号
	Name       string  `json:"name"`        // 币种名称
	Balance    float64 `json:"balance"`     // 余额
	PriceUSD   float64 `json:"price_usd"`   // USD 价格
	ValueUSD   float64 `json:"value_usd"`   // USD 价值
}

// AssetSummary 资产概览
type AssetSummary struct {
	TotalValue      float64            `json:"total_value"`      // 总价值
	Currency        string             `json:"currency"`         // 计价货币
	CEXValue        float64            `json:"cex_value"`        // CEX 资产价值
	BlockchainValue float64            `json:"blockchain_value"` // 区块链资产价值
	ManualValue     float64            `json:"manual_value"`     // 手动资产价值
	Change24h       float64            `json:"change_24h"`       // 24小时变化
	ChangePercent   float64            `json:"change_percent"`   // 变化百分比
	TopAssets       []AssetDetail      `json:"top_assets"`       // 前5大资产
	ByCategory      map[string]float64 `json:"by_category"`      // 按类别分组
	LastUpdated     int64              `json:"last_updated"`     // 最后更新时间戳
}

// AssetDetail 资产详情
type AssetDetail struct {
	ID         uint    `json:"id"`
	Symbol     string  `json:"symbol"`      // 币种符号
	Name       string  `json:"name"`        // 币种名称
	Balance    float64 `json:"balance"`     // 余额
	Price      float64 `json:"price"`       // 当前价格
	Value      float64 `json:"value"`       // 当前价值
	Currency   string  `json:"currency"`    // 计价货币
	Source     string  `json:"source"`      // 来源
	SourceType string  `json:"source_type"` // 来源类型
	Change24h  float64 `json:"change_24h"`  // 24小时变化
	Percentage float64 `json:"percentage"`  // 占比
}

// ImportedWallet 成功导入的钱包
type ImportedWallet struct {
	ID         uint   `json:"id"`
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
	Label      string `json:"label"`
}

// FailedImport 导入失败的记录
type FailedImport struct {
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

// BatchImportResult 批量导入结果
type BatchImportResult struct {
	Success []ImportedWallet `json:"success"`
	Failed  []FailedImport   `json:"failed"`
	Total   int              `json:"total"`
}

// HistoryPoint 历史数据点
type HistoryPoint struct {
	Timestamp int64   `json:"timestamp"` // 时间戳
	Value     float64 `json:"value"`     // 价值
	Currency  string  `json:"currency"`  // 计价货币
}
