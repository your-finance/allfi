// Package entity 定义数据库实体模型
package entity

import (
	"time"
)

// ExchangeAccount 交易所账户实体
// 存储用户在各交易所的 API 凭证（加密存储）
type ExchangeAccount struct {
	Id                     uint           `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId                 uint           `gorm:"not null;index:idx_user_exchange;comment:用户ID" json:"user_id"`
	ExchangeName           string         `gorm:"size:50;not null;index:idx_user_exchange;comment:交易所名称(binance/okx/coinbase)" json:"exchange_name"`
	ApiKeyEncrypted        string         `gorm:"type:text;not null;comment:加密后的API Key" json:"-"`
	ApiSecretEncrypted     string         `gorm:"type:text;not null;comment:加密后的API Secret" json:"-"`
	ApiPassphraseEncrypted string         `gorm:"type:text;comment:加密后的API Passphrase(部分交易所需要)" json:"-"`
	Label                  string         `gorm:"size:100;comment:账户标签" json:"label"`
	IsActive               bool           `gorm:"not null;default:true;index:idx_active;comment:是否启用" json:"is_active"`
	CreatedAt              time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt              *time.Time     `gorm:"index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (ExchangeAccount) TableName() string {
	return "exchange_accounts"
}
