// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package entity

import (
	"time"
)

// NotificationPreferences is the golang structure for table notification_preferences.
type NotificationPreferences struct {
	Id                  int       `json:"id"                    orm:"id"                    description:""` //
	CreatedAt           time.Time `json:"created_at"            orm:"created_at"            description:""` //
	UpdatedAt           time.Time `json:"updated_at"            orm:"updated_at"            description:""` //
	DeletedAt           time.Time `json:"deleted_at"            orm:"deleted_at"            description:""` //
	UserId              int       `json:"user_id"               orm:"user_id"               description:""` //
	EnableDailyDigest   float64   `json:"enable_daily_digest"   orm:"enable_daily_digest"   description:""` //
	DigestTime          string    `json:"digest_time"           orm:"digest_time"           description:""` //
	EnablePriceAlert    float64   `json:"enable_price_alert"    orm:"enable_price_alert"    description:""` //
	EnableAssetAlert    float64   `json:"enable_asset_alert"    orm:"enable_asset_alert"    description:""` //
	AssetAlertThreshold float32   `json:"asset_alert_threshold" orm:"asset_alert_threshold" description:""` //
	WebhookUrl          string    `json:"webhook_url"           orm:"webhook_url"           description:""` //
}
