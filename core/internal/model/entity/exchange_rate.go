package entity

import (
	"time"
)

// ExchangeRate 汇率实体
// 缓存各币种之间的汇率数据
type ExchangeRate struct {
	Id         uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	FromSymbol string    `gorm:"size:50;not null;index:idx_from_to;comment:源币种" json:"from_symbol"`
	ToSymbol   string    `gorm:"size:50;not null;index:idx_from_to;comment:目标币种" json:"to_symbol"`
	Rate       float64   `gorm:"type:decimal(20,8);not null;comment:汇率" json:"rate"`
	Source     string    `gorm:"size:50;comment:数据源(coingecko/yahoo)" json:"source"`
	ExpireAt   time.Time `gorm:"not null;index:idx_expire;comment:过期时间" json:"expire_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (ExchangeRate) TableName() string {
	return "exchange_rates"
}
