// Package model 通知模块 - 业务数据传输对象
// 定义通知模块内部使用的 DTO 和常量
package model

// 通知类型常量
const (
	NotifyPriceAlert  = "price_alert"  // 价格预警通知
	NotifyAssetChange = "asset_change" // 资产变动通知
	NotifyDailyDigest = "daily_digest" // 每日摘要通知
	NotifySystem      = "system"       // 系统通知
)

// DefaultPreference 默认通知偏好设置
type DefaultPreference struct {
	EmailEnabled   bool // 邮件通知
	PushEnabled    bool // 推送通知
	PriceAlert     bool // 价格预警
	PortfolioAlert bool // 资产变动
	SystemNotice   bool // 系统通知
}

// GetDefault 获取默认偏好设置
func GetDefault() DefaultPreference {
	return DefaultPreference{
		EmailEnabled:   false,
		PushEnabled:    true,
		PriceAlert:     true,
		PortfolioAlert: true,
		SystemNotice:   true,
	}
}
