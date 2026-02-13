// Package utils 配置加载工具
// 封装 GoFrame 配置管理，提供便捷的配置获取方法
package utils

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// Config 配置结构体
type Config struct {
	ctx context.Context
}

// NewConfig 创建配置实例
func NewConfig(ctx context.Context) *Config {
	return &Config{ctx: ctx}
}

// GetString 获取字符串配置
func (c *Config) GetString(key string, defaultValue ...string) string {
	value := g.Cfg().MustGet(c.ctx, key).String()
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string, defaultValue ...int) int {
	value := g.Cfg().MustGet(c.ctx, key).Int()
	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value := g.Cfg().MustGet(c.ctx, key)
	if value.IsNil() || value.IsEmpty() {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return value.Bool()
}

// GetFloat64 获取浮点数配置
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	value := g.Cfg().MustGet(c.ctx, key).Float64()
	if value == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// 环境变量支持

// GetEnvString 优先从环境变量获取，否则从配置文件获取
func (c *Config) GetEnvString(envKey, configKey string, defaultValue ...string) string {
	// 优先环境变量
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	// 然后配置文件
	return c.GetString(configKey, defaultValue...)
}

// GetEnvInt 优先从环境变量获取整数
func (c *Config) GetEnvInt(envKey, configKey string, defaultValue ...int) int {
	if value := os.Getenv(envKey); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return c.GetInt(configKey, defaultValue...)
}

// GetEnvBool 优先从环境变量获取布尔值
func (c *Config) GetEnvBool(envKey, configKey string, defaultValue ...bool) bool {
	if value := os.Getenv(envKey); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return c.GetBool(configKey, defaultValue...)
}

// 便捷方法 - 服务器配置

// GetServerPort 获取服务器端口
func (c *Config) GetServerPort() int {
	return c.GetEnvInt("ALLFI_PORT", "server.port", 8080)
}

// GetServerMode 获取服务器模式（development/production）
func (c *Config) GetServerMode() string {
	return c.GetEnvString("ALLFI_MODE", "server.mode", "development")
}

// IsProduction 判断是否生产环境
func (c *Config) IsProduction() bool {
	return c.GetServerMode() == "production"
}

// 便捷方法 - 数据库配置

// GetDatabaseType 获取数据库类型（sqlite/mysql）
func (c *Config) GetDatabaseType() string {
	return c.GetEnvString("ALLFI_DB_TYPE", "database.type", "sqlite")
}

// GetSQLitePath 获取 SQLite 数据库路径
func (c *Config) GetSQLitePath() string {
	return c.GetEnvString("ALLFI_DB_PATH", "database.sqlite.path", "../data/allfi.db")
}

// GetMySQLDSN 获取 MySQL 连接字符串
func (c *Config) GetMySQLDSN() string {
	if dsn := os.Getenv("ALLFI_MYSQL_DSN"); dsn != "" {
		return dsn
	}
	host := c.GetString("database.mysql.host", "localhost")
	port := c.GetInt("database.mysql.port", 3306)
	user := c.GetString("database.mysql.user", "root")
	password := c.GetString("database.mysql.password", "")
	database := c.GetString("database.mysql.database", "allfi")
	charset := c.GetString("database.mysql.charset", "utf8mb4")

	return user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database + "?charset=" + charset + "&parseTime=True&loc=Local"
}

// 便捷方法 - 安全配置

// GetMasterKey 获取加密主密钥
func (c *Config) GetMasterKey() string {
	return c.GetEnvString("ALLFI_MASTER_KEY", "security.master_key", "")
}

// IsHTTPSEnabled 是否启用 HTTPS
func (c *Config) IsHTTPSEnabled() bool {
	return c.GetBool("security.enable_https", false)
}

// 便捷方法 - 外部 API 配置

// GetEtherscanAPIKey 获取 Etherscan API Key
func (c *Config) GetEtherscanAPIKey() string {
	return c.GetEnvString("ETHERSCAN_API_KEY", "external_apis.etherscan.api_key", "")
}

// GetBscScanAPIKey 获取 BscScan API Key
func (c *Config) GetBscScanAPIKey() string {
	return c.GetEnvString("BSCSCAN_API_KEY", "external_apis.bscscan.api_key", "")
}

// GetCoinGeckoAPIKey 获取 CoinGecko API Key（可选）
func (c *Config) GetCoinGeckoAPIKey() string {
	return c.GetEnvString("COINGECKO_API_KEY", "external_apis.coingecko.api_key", "")
}

// 便捷方法 - 定时任务配置

// GetSnapshotInterval 获取快照间隔（秒）
func (c *Config) GetSnapshotInterval() int {
	return c.GetInt("cron.snapshot_interval", 3600)
}

// GetPriceCacheTTL 获取价格缓存时间（秒）
func (c *Config) GetPriceCacheTTL() int {
	return c.GetInt("cron.price_cache_ttl", 300)
}

// 便捷方法 - 默认设置

// GetDefaultCurrency 获取默认计价货币
func (c *Config) GetDefaultCurrency() string {
	return c.GetString("defaults.currency", "USDC")
}

// GetHistoryRetentionDays 获取历史数据保留天数
func (c *Config) GetHistoryRetentionDays() int {
	return c.GetInt("defaults.history_retention_days", 90)
}

// 全局配置实例
var globalConfig *Config

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig(ctx context.Context) {
	globalConfig = NewConfig(ctx)
}

// GetGlobalConfig 获取全局配置实例
func GetGlobalConfig() *Config {
	if globalConfig == nil {
		globalConfig = NewConfig(context.Background())
	}
	return globalConfig
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	ctx := context.Background()
	// 设置配置文件路径
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath(path)
	return NewConfig(ctx), nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return NewConfig(context.Background())
}

// GetDuration 获取时间间隔配置
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	seconds := c.GetInt(key)
	if seconds == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return time.Duration(seconds) * time.Second
}
