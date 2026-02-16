// Package database 提供数据库初始化和迁移功能
// 在应用启动时自动创建所有必要的表
package database

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	// 确保 SQLite 驱动已加载
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

// init 在包加载时自动初始化数据库
// 必须在其他业务包的 init() 之前执行（通过 import 顺序保证）
func init() {
	ctx := context.Background()
	if err := Initialize(ctx); err != nil {
		g.Log().Fatalf(ctx, "数据库初始化失败: %v", err)
	}
}

// tables 定义所有需要创建的表及其 DDL
var tables = []string{
	// 系统配置表
	`CREATE TABLE IF NOT EXISTS system_config (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at  DATETIME,
		updated_at  DATETIME,
		deleted_at  DATETIME,
		config_key  TEXT NOT NULL,
		config_value TEXT NOT NULL DEFAULT '',
		description TEXT DEFAULT ''
	)`,

	// 用户表
	`CREATE TABLE IF NOT EXISTS users (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at    DATETIME,
		updated_at    DATETIME,
		deleted_at    DATETIME,
		username      TEXT NOT NULL,
		email         TEXT NOT NULL DEFAULT '',
		password_hash TEXT NOT NULL DEFAULT '',
		nickname      TEXT DEFAULT '',
		avatar        TEXT DEFAULT '',
		status        TEXT NOT NULL DEFAULT 'active',
		last_login_at DATETIME,
		last_login_ip TEXT DEFAULT ''
	)`,

	// 交易所账户表
	`CREATE TABLE IF NOT EXISTS exchange_accounts (
		id                       INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at               DATETIME,
		updated_at               DATETIME,
		deleted_at               DATETIME,
		user_id                  INTEGER NOT NULL,
		exchange_name            TEXT NOT NULL,
		api_key_encrypted        TEXT NOT NULL DEFAULT '',
		api_secret_encrypted     TEXT NOT NULL DEFAULT '',
		api_passphrase_encrypted TEXT DEFAULT '',
		label                    TEXT DEFAULT '',
		note                     TEXT DEFAULT '',
		is_active                INTEGER NOT NULL DEFAULT 1
	)`,

	// 钱包地址表
	`CREATE TABLE IF NOT EXISTS wallet_addresses (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id    INTEGER NOT NULL,
		blockchain TEXT NOT NULL,
		address    TEXT NOT NULL,
		label      TEXT DEFAULT '',
		is_active  REAL NOT NULL DEFAULT 1
	)`,

	// 手动资产表
	`CREATE TABLE IF NOT EXISTS manual_assets (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id    INTEGER NOT NULL,
		asset_type TEXT NOT NULL DEFAULT '',
		asset_name TEXT NOT NULL DEFAULT '',
		amount     REAL NOT NULL DEFAULT 0,
		amount_usd REAL NOT NULL DEFAULT 0,
		currency   TEXT DEFAULT '',
		notes      TEXT DEFAULT '',
		is_active  INTEGER NOT NULL DEFAULT 1
	)`,

	// 资产详情表
	`CREATE TABLE IF NOT EXISTS asset_details (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at   DATETIME,
		updated_at   DATETIME,
		deleted_at   DATETIME,
		user_id      INTEGER NOT NULL,
		source_type  TEXT NOT NULL DEFAULT '',
		source_id    INTEGER NOT NULL DEFAULT 0,
		asset_symbol TEXT NOT NULL DEFAULT '',
		asset_name   TEXT DEFAULT '',
		balance      REAL NOT NULL DEFAULT 0,
		price_usd    REAL DEFAULT 0,
		value_usd    REAL DEFAULT 0,
		last_updated DATETIME
	)`,

	// 资产快照表
	`CREATE TABLE IF NOT EXISTS asset_snapshots (
		id                   INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at           DATETIME,
		updated_at           DATETIME,
		deleted_at           DATETIME,
		user_id              INTEGER NOT NULL,
		snapshot_time        DATETIME NOT NULL,
		total_value_usd      REAL DEFAULT 0,
		total_value_cny      REAL DEFAULT 0,
		total_value_btc      REAL DEFAULT 0,
		cex_value_usd        REAL DEFAULT 0,
		blockchain_value_usd REAL DEFAULT 0,
		manual_value_usd     REAL DEFAULT 0,
		exchange_rates_json  TEXT DEFAULT ''
	)`,

	// 汇率表
	`CREATE TABLE IF NOT EXISTS exchange_rates (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at    DATETIME,
		updated_at    DATETIME,
		deleted_at    DATETIME,
		from_currency TEXT NOT NULL,
		to_currency   TEXT NOT NULL,
		rate          REAL NOT NULL DEFAULT 0,
		source        TEXT DEFAULT '',
		fetched_at    DATETIME
	)`,

	// 统一交易记录表
	`CREATE TABLE IF NOT EXISTS unified_transactions (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at  DATETIME,
		updated_at  DATETIME,
		deleted_at  DATETIME,
		user_id     INTEGER NOT NULL,
		tx_type     TEXT NOT NULL DEFAULT '',
		source      TEXT DEFAULT '',
		source_id   TEXT DEFAULT '',
		from_asset  TEXT DEFAULT '',
		from_amount REAL DEFAULT 0,
		to_asset    TEXT DEFAULT '',
		to_amount   REAL DEFAULT 0,
		fee         REAL DEFAULT 0,
		fee_coin    TEXT DEFAULT '',
		value_usd   REAL DEFAULT 0,
		tx_hash     TEXT DEFAULT '',
		chain       TEXT DEFAULT '',
		timestamp   DATETIME
	)`,

	// 交易日汇总表
	`CREATE TABLE IF NOT EXISTS transaction_daily_summaries (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at    DATETIME,
		updated_at    DATETIME,
		deleted_at    DATETIME,
		user_id       INTEGER NOT NULL,
		date          TEXT NOT NULL,
		buy_count     INTEGER DEFAULT 0,
		sell_count    INTEGER DEFAULT 0,
		total_count   INTEGER DEFAULT 0,
		total_fee_usd REAL DEFAULT 0,
		net_flow_usd  REAL DEFAULT 0
	)`,

	// 通知表
	`CREATE TABLE IF NOT EXISTS notifications (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id    INTEGER NOT NULL,
		type       TEXT NOT NULL DEFAULT '',
		title      TEXT DEFAULT '',
		content    TEXT DEFAULT '',
		is_read    INTEGER NOT NULL DEFAULT 0
	)`,

	// 通知偏好表
	`CREATE TABLE IF NOT EXISTS notification_preferences (
		id                    INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at            DATETIME,
		updated_at            DATETIME,
		deleted_at            DATETIME,
		user_id               INTEGER NOT NULL,
		enable_daily_digest   INTEGER NOT NULL DEFAULT 0,
		digest_time           TEXT DEFAULT '08:00',
		enable_price_alert    INTEGER NOT NULL DEFAULT 1,
		enable_asset_alert    INTEGER NOT NULL DEFAULT 1,
		asset_alert_threshold REAL DEFAULT 5.0,
		webhook_url           TEXT DEFAULT ''
	)`,

	// 价格预警表
	`CREATE TABLE IF NOT EXISTS price_alerts (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at   DATETIME,
		updated_at   DATETIME,
		deleted_at   DATETIME,
		user_id      INTEGER NOT NULL,
		symbol       TEXT NOT NULL DEFAULT '',
		condition    TEXT NOT NULL DEFAULT '',
		target_price REAL NOT NULL DEFAULT 0,
		is_active    INTEGER NOT NULL DEFAULT 1,
		triggered    INTEGER NOT NULL DEFAULT 0,
		triggered_at DATETIME,
		note         TEXT DEFAULT ''
	)`,

	// 报告表
	`CREATE TABLE IF NOT EXISTS reports (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at     DATETIME,
		updated_at     DATETIME,
		deleted_at     DATETIME,
		user_id        INTEGER NOT NULL,
		type           TEXT NOT NULL DEFAULT '',
		period         TEXT DEFAULT '',
		total_value    REAL DEFAULT 0,
		change         REAL DEFAULT 0,
		change_percent REAL DEFAULT 0,
		top_gainers    TEXT DEFAULT '',
		top_losers     TEXT DEFAULT '',
		btc_benchmark  REAL DEFAULT 0,
		eth_benchmark  REAL DEFAULT 0,
		content        TEXT DEFAULT '',
		generated_at   DATETIME
	)`,

	// 策略表
	`CREATE TABLE IF NOT EXISTS strategies (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at     DATETIME,
		updated_at     DATETIME,
		deleted_at     DATETIME,
		user_id        INTEGER NOT NULL,
		name           TEXT NOT NULL DEFAULT '',
		type           TEXT NOT NULL DEFAULT '',
		config         TEXT DEFAULT '',
		is_active      INTEGER NOT NULL DEFAULT 1,
		last_checked   DATETIME,
		last_triggered DATETIME,
		trigger_count  INTEGER DEFAULT 0
	)`,

	// 目标表
	`CREATE TABLE IF NOT EXISTS goals (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at   DATETIME,
		updated_at   DATETIME,
		deleted_at   DATETIME,
		title        TEXT NOT NULL DEFAULT '',
		type         TEXT NOT NULL DEFAULT '',
		target_value REAL NOT NULL DEFAULT 0,
		currency     TEXT DEFAULT 'USD',
		deadline     DATETIME
	)`,

	// 同步元数据表
	`CREATE TABLE IF NOT EXISTS sync_metadata (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at     DATETIME,
		updated_at     DATETIME,
		deleted_at     DATETIME,
		source         TEXT NOT NULL DEFAULT '',
		last_sync_time DATETIME,
		last_sync_id   TEXT DEFAULT '',
		last_block     INTEGER DEFAULT 0,
		tx_count       INTEGER DEFAULT 0
	)`,

	// NFT 缓存表
	`CREATE TABLE IF NOT EXISTS nft_caches (
		id                INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at        DATETIME,
		updated_at        DATETIME,
		deleted_at        DATETIME,
		user_id           INTEGER NOT NULL,
		wallet_address    TEXT NOT NULL DEFAULT '',
		contract_address  TEXT DEFAULT '',
		token_id          TEXT DEFAULT '',
		name              TEXT DEFAULT '',
		description       TEXT DEFAULT '',
		image_url         TEXT DEFAULT '',
		collection        TEXT DEFAULT '',
		collection_slug   TEXT DEFAULT '',
		chain             TEXT DEFAULT '',
		floor_price       REAL DEFAULT 0,
		floor_currency    TEXT DEFAULT '',
		floor_price_usd   REAL DEFAULT 0,
		cached_at         DATETIME
	)`,

	// 用户成就表
	`CREATE TABLE IF NOT EXISTS user_achievements (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at     DATETIME,
		updated_at     DATETIME,
		deleted_at     DATETIME,
		user_id        INTEGER NOT NULL,
		achievement_id TEXT NOT NULL DEFAULT '',
		unlocked_at    DATETIME
	)`,

	// Web Push 订阅表
	`CREATE TABLE IF NOT EXISTS web_push_subscriptions (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		user_id    INTEGER NOT NULL,
		endpoint   TEXT NOT NULL DEFAULT '',
		p256dh     TEXT NOT NULL DEFAULT '',
		auth       TEXT NOT NULL DEFAULT ''
	)`,
}

// indexes 定义常用索引
var indexes = []string{
	`CREATE INDEX IF NOT EXISTS idx_system_config_key ON system_config(config_key)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
	`CREATE INDEX IF NOT EXISTS idx_exchange_accounts_user ON exchange_accounts(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_wallet_addresses_user ON wallet_addresses(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_asset_details_user ON asset_details(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_asset_snapshots_user_time ON asset_snapshots(user_id, snapshot_time)`,
	`CREATE INDEX IF NOT EXISTS idx_exchange_rates_pair ON exchange_rates(from_currency, to_currency)`,
	`CREATE INDEX IF NOT EXISTS idx_unified_transactions_user ON unified_transactions(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_price_alerts_user ON price_alerts(user_id)`,
	`CREATE INDEX IF NOT EXISTS idx_manual_assets_user ON manual_assets(user_id)`,
}

// Initialize 初始化数据库：创建所有表和索引
// 使用 CREATE TABLE IF NOT EXISTS，幂等安全
func Initialize(ctx context.Context) error {
	db := g.DB()

	return db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建表
		for _, ddl := range tables {
			if _, err := tx.Exec(ddl); err != nil {
				g.Log().Errorf(ctx, "创建表失败: %v\nSQL: %s", err, ddl)
				return err
			}
		}

		// 创建索引
		for _, idx := range indexes {
			if _, err := tx.Exec(idx); err != nil {
				g.Log().Warningf(ctx, "创建索引警告: %v\nSQL: %s", err, idx)
				// 索引创建失败不阻止启动
			}
		}

		g.Log().Info(ctx, "数据库表初始化完成")
		return nil
	})
}
