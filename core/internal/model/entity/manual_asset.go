package entity

import (
	"time"
)

// ManualAsset 手动资产实体
// 存储传统金融资产（银行存款、现金、股票、基金）
type ManualAsset struct {
	Id          uint           `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserId      uint           `gorm:"not null;index:idx_user_type;comment:用户ID" json:"user_id"`
	AssetType   string         `gorm:"size:50;not null;index:idx_user_type;comment:资产类型(bank/cash/stock/fund)" json:"asset_type"`
	AssetName   string         `gorm:"size:100;not null;comment:资产名称" json:"asset_name"`
	Institution string         `gorm:"size:100;comment:机构名称（银行/券商/基金平台）" json:"institution"`
	AssetSymbol string         `gorm:"size:50;not null;comment:资产符号/货币" json:"asset_symbol"`
	Amount      float64        `gorm:"type:decimal(20,8);not null;comment:数量" json:"amount"`
	ValueUSD    float64        `gorm:"type:decimal(20,2);comment:美元估值" json:"value_usd"`
	Label       string         `gorm:"size:100;comment:备注标签" json:"label"`
	IsActive    bool           `gorm:"not null;default:true;index:idx_active;comment:是否启用" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   *time.Time     `gorm:"index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (ManualAsset) TableName() string {
	return "manual_assets"
}
