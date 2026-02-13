// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2026-02-12 11:01:49
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationPreferences is the golang structure of table notification_preferences for DAO operations like Where/Data.
type NotificationPreferences struct {
	g.Meta              `orm:"table:notification_preferences, do:true"`
	Id                  any //
	CreatedAt           any //
	UpdatedAt           any //
	DeletedAt           any //
	UserId              any //
	EnableDailyDigest   any //
	DigestTime          any //
	EnablePriceAlert    any //
	EnableAssetAlert    any //
	AssetAlertThreshold any //
	WebhookUrl          any //
}
