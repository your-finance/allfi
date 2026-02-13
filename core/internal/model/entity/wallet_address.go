package entity

import (
	"time"
)

// WalletAddress 区块链钱包地址实体
// 存储用户的链上钱包地址
type WalletAddress struct {
	Id         uint           `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId     uint           `gorm:"not null;index:idx_user_blockchain;comment:用户ID" json:"user_id"`
	Blockchain string         `gorm:"size:50;not null;index:idx_user_blockchain;comment:区块链名称(ethereum/bsc/polygon)" json:"blockchain"`
	Address    string         `gorm:"size:100;not null;index:idx_address;comment:钱包地址" json:"address"`
	Label      string         `gorm:"size:100;comment:地址标签" json:"label"`
	IsActive   bool           `gorm:"not null;default:true;index:idx_active;comment:是否启用" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt  *time.Time     `gorm:"index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (WalletAddress) TableName() string {
	return "wallet_addresses"
}
