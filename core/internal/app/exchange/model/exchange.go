// Package model 交易所模块 - 业务数据传输对象
// 定义交易所模块内部使用的 DTO 和转换方法
package model

// CreateAccountInput 创建账户的内部输入参数
type CreateAccountInput struct {
	UserID       int    // 用户 ID
	ExchangeName string // 交易所名称
	ApiKey       string // API Key（明文，待加密）
	ApiSecret    string // API Secret（明文，待加密）
	Passphrase   string // API Passphrase（明文，待加密，OKX 必填）
	Label        string // 账户标签
	Note         string // 备注
}

// UpdateAccountInput 更新账户的内部输入参数
type UpdateAccountInput struct {
	AccountID  int    // 账户 ID
	ApiKey     string // API Key（可选更新）
	ApiSecret  string // API Secret（可选更新）
	Passphrase string // API Passphrase（可选更新）
	Label      string // 账户标签
	Note       string // 备注
}
