// Package utils API Key 解析工具
// 提供统一的 API Key 获取逻辑：数据库（加密存储）> 环境变量 > 空字符串
package utils

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/dao"
	"your-finance/allfi/internal/model/entity"
	"your-finance/allfi/utility/crypto"
)

// apiKeyPrefix 数据库中 API Key 的存储前缀
const apiKeyPrefix = "apikey."

// providerEnvMap 服务商到环境变量的映射
var providerEnvMap = map[string]string{
	"etherscan": "ETHERSCAN_API_KEY",
	"bscscan":   "BSCSCAN_API_KEY",
	"coingecko": "COINGECKO_API_KEY",
}

// ResolveAPIKey 获取指定服务商的 API Key
// 优先级：数据库（加密存储）> 环境变量 > 空字符串
// provider: 服务商标识（etherscan / bscscan / coingecko）
func ResolveAPIKey(ctx context.Context, provider string) string {
	// 1. 先查数据库
	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()
	if len(masterKey) == 32 {
		configKey := apiKeyPrefix + provider
		var config entity.SystemConfig
		err := dao.SystemConfig.Ctx(ctx).
			Where(dao.SystemConfig.Columns().ConfigKey, configKey).
			WhereNull(dao.SystemConfig.Columns().DeletedAt).
			Scan(&config)
		if err == nil && config.ConfigValue != "" {
			plainKey, decErr := crypto.DecryptAES(config.ConfigValue, masterKey)
			if decErr == nil && plainKey != "" {
				return plainKey
			}
		}
	}

	// 2. 查环境变量
	if envKey, ok := providerEnvMap[provider]; ok {
		if val := os.Getenv(envKey); val != "" {
			return val
		}
	}

	// 3. 都没有返回空字符串
	return ""
}
