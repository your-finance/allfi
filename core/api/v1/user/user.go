// Package user 用户模块 API 定义
// 提供用户设置的读写和缓存清除接口
package user

import "github.com/gogf/gf/v2/frame/g"

// GetSettingsReq 获取用户设置请求
type GetSettingsReq struct {
	g.Meta `path:"/users/settings" method:"get" summary:"获取用户设置" tags:"用户"`
}

// GetSettingsRes 获取用户设置响应
type GetSettingsRes struct {
	Settings map[string]string `json:"settings" dc:"用户设置键值对"`
}

// UpdateSettingsReq 更新用户设置请求
type UpdateSettingsReq struct {
	g.Meta   `path:"/users/settings" method:"put" summary:"更新用户设置" tags:"用户"`
	Settings map[string]string `json:"settings" v:"required" dc:"设置键值对"`
}

// UpdateSettingsRes 更新用户设置响应
type UpdateSettingsRes struct {
	Message string `json:"message" dc:"操作结果消息"`
}

// ResetSettingsReq 重置用户设置请求
type ResetSettingsReq struct {
	g.Meta `path:"/users/reset-settings" method:"post" summary:"重置用户设置" tags:"用户"`
}

// ResetSettingsRes 重置用户设置响应
type ResetSettingsRes struct {
	Message string `json:"message" dc:"操作结果消息"`
}

// ClearCacheReq 清除缓存请求
type ClearCacheReq struct {
	g.Meta `path:"/users/clear-cache" method:"post" summary:"清除缓存" tags:"用户"`
}

// ClearCacheRes 清除缓存响应
type ClearCacheRes struct {
	Message string `json:"message" dc:"操作结果消息"`
}

// ExportDataReq 导出用户数据请求
type ExportDataReq struct {
	g.Meta `path:"/users/export" method:"get" summary:"导出用户数据" tags:"用户"`
}

// ExportExchangeAccount 导出用的交易所账户（不含加密凭证）
type ExportExchangeAccount struct {
	ID           uint   `json:"id" dc:"账户 ID"`
	ExchangeName string `json:"exchange_name" dc:"交易所名称"`
	Label        string `json:"label" dc:"标签"`
	Note         string `json:"note" dc:"备注"`
}

// ExportWalletAddress 导出用的钱包地址
type ExportWalletAddress struct {
	ID         uint   `json:"id" dc:"钱包 ID"`
	Blockchain string `json:"blockchain" dc:"区块链网络"`
	Address    string `json:"address" dc:"钱包地址"`
	Label      string `json:"label" dc:"标签"`
}

// ExportManualAsset 导出用的手动资产
type ExportManualAsset struct {
	ID        uint    `json:"id" dc:"资产 ID"`
	AssetType string  `json:"asset_type" dc:"资产类型"`
	AssetName string  `json:"asset_name" dc:"资产名称"`
	Amount    float64 `json:"amount" dc:"数量"`
	AmountUSD float64 `json:"amount_usd" dc:"USD 价值"`
	Currency  string  `json:"currency" dc:"币种"`
	Notes     string  `json:"notes" dc:"备注"`
}

// ExportStrategy 导出用的策略
type ExportStrategy struct {
	ID       uint   `json:"id" dc:"策略 ID"`
	Name     string `json:"name" dc:"策略名称"`
	Type     string `json:"type" dc:"策略类型"`
	Config   string `json:"config" dc:"策略配置（JSON）"`
	IsActive bool   `json:"is_active" dc:"是否启用"`
}

// ExportGoal 导出用的目标
type ExportGoal struct {
	ID          uint    `json:"id" dc:"目标 ID"`
	Title       string  `json:"title" dc:"目标标题"`
	Type        string  `json:"type" dc:"目标类型"`
	TargetValue float64 `json:"target_value" dc:"目标值"`
	Currency    string  `json:"currency" dc:"币种"`
	Deadline    string  `json:"deadline" dc:"截止日期"`
}

// ExportPriceAlert 导出用的价格预警
type ExportPriceAlert struct {
	ID          uint    `json:"id" dc:"预警 ID"`
	Symbol      string  `json:"symbol" dc:"币种"`
	Condition   string  `json:"condition" dc:"触发条件"`
	TargetPrice float64 `json:"target_price" dc:"目标价格"`
	IsActive    bool    `json:"is_active" dc:"是否启用"`
	Note        string  `json:"note" dc:"备注"`
}

// ExportDataRes 导出用户数据响应
type ExportDataRes struct {
	ExchangeAccounts []ExportExchangeAccount `json:"exchange_accounts" dc:"交易所账户"`
	WalletAddresses  []ExportWalletAddress   `json:"wallet_addresses" dc:"钱包地址"`
	ManualAssets     []ExportManualAsset     `json:"manual_assets" dc:"手动资产"`
	Strategies       []ExportStrategy        `json:"strategies" dc:"策略"`
	Goals            []ExportGoal            `json:"goals" dc:"目标"`
	PriceAlerts      []ExportPriceAlert      `json:"price_alerts" dc:"价格预警"`
	Settings         map[string]string       `json:"settings" dc:"用户设置"`
	ExportedAt       string                  `json:"exported_at" dc:"导出时间"`
}
