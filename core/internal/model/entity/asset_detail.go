package entity

import (
	"time"
)

// AssetDetail 资产详情实体
// 存储从各来源聚合的实时资产数据
type AssetDetail struct {
	Id           uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId       uint      `gorm:"not null;index:idx_user_source;comment:用户ID" json:"user_id"`
	SourceType   string    `gorm:"size:50;not null;index:idx_user_source;comment:来源类型(cex/blockchain/manual)" json:"source_type"`
	SourceId     uint      `gorm:"not null;index:idx_user_source;comment:来源ID(关联account_id/wallet_id/manual_id)" json:"source_id"`
	SourceName   string    `gorm:"size:100;not null;comment:来源名称" json:"source_name"`
	AssetSymbol  string    `gorm:"size:50;not null;index:idx_symbol;comment:资产符号" json:"asset_symbol"`
	AssetName    string    `gorm:"size:100;comment:资产名称" json:"asset_name"`
	Amount       float64   `gorm:"type:decimal(20,8);not null;comment:持有数量" json:"amount"`
	PriceUSD     float64   `gorm:"type:decimal(20,2);comment:单价(USD)" json:"price_usd"`
	ValueUSD     float64   `gorm:"type:decimal(20,2);comment:总价值(USD)" json:"value_usd"`
	LastSyncTime time.Time `gorm:"not null;index:idx_sync_time;comment:最后同步时间" json:"last_sync_time"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (AssetDetail) TableName() string {
	return "asset_details"
}
