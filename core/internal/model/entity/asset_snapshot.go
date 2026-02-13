package entity

import (
	"time"
)

// AssetSnapshot 资产快照实体
// 定期记录资产总值，用于历史趋势分析
type AssetSnapshot struct {
	Id              uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId          uint      `gorm:"not null;index:idx_user_time;comment:用户ID" json:"user_id"`
	TotalValueUSD   float64   `gorm:"type:decimal(20,2);not null;comment:总价值(USD)" json:"total_value_usd"`
	TotalValueBTC   float64   `gorm:"type:decimal(20,8);comment:总价值(BTC)" json:"total_value_btc"`
	TotalValueETH   float64   `gorm:"type:decimal(20,8);comment:总价值(ETH)" json:"total_value_eth"`
	TotalValueCNY   float64   `gorm:"type:decimal(20,2);comment:总价值(CNY)" json:"total_value_cny"`
	CexValueUSD     float64   `gorm:"type:decimal(20,2);comment:CEX资产价值(USD)" json:"cex_value_usd"`
	BlockchainValue float64   `gorm:"type:decimal(20,2);comment:链上资产价值(USD)" json:"blockchain_value_usd"`
	ManualValueUSD  float64   `gorm:"type:decimal(20,2);comment:传统资产价值(USD)" json:"manual_value_usd"`
	SnapshotTime    time.Time `gorm:"not null;index:idx_user_time;comment:快照时间" json:"snapshot_time"`
	CreatedAt       time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

// TableName 指定表名
func (AssetSnapshot) TableName() string {
	return "asset_snapshots"
}
