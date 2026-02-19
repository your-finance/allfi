// Package logic API Key 管理业务逻辑
// 实现 API Key 的加密存储、脱敏查询、删除功能
// 使用 system_config 表存储，key 前缀为 "apikey."
package logic

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	systemApi "your-finance/allfi/api/v1/system"
	"your-finance/allfi/internal/app/system/dao"
	"your-finance/allfi/internal/model/entity"
	"your-finance/allfi/utility/crypto"
)

// API Key 存储前缀
const apiKeyPrefix = "apikey."

// 支持的 API Key 服务商配置
var apiKeyProviders = []struct {
	Provider    string
	DisplayName string
	Description string
	EnvKey      string // 对应的环境变量名
}{
	{"etherscan", "Etherscan", "以太坊区块链浏览器 API（余额查询、Gas 价格、交易记录）", "ETHERSCAN_API_KEY"},
	{"bscscan", "BscScan", "BSC 区块链浏览器 API（余额查询、Gas 价格、交易记录）", "BSCSCAN_API_KEY"},
	{"coingecko", "CoinGecko", "加密货币价格 API（可选，免费 API 无需 Key）", "COINGECKO_API_KEY"},
}

// getMasterKey 获取加密主密钥
func getMasterKey(ctx context.Context) (string, error) {
	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()
	if len(masterKey) != 32 {
		return "", gerror.New("主加密密钥配置错误，长度必须为32字节")
	}
	return masterKey, nil
}

// maskAPIKey 脱敏 API Key（显示前4位和后4位）
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// GetAPIKeys 获取所有 API Key 配置（脱敏显示）
func (s *sSystem) GetAPIKeys(ctx context.Context) (*systemApi.GetAPIKeysRes, error) {
	masterKey, err := getMasterKey(ctx)
	if err != nil {
		return nil, err
	}

	keys := make([]systemApi.APIKeyItem, 0, len(apiKeyProviders))

	for _, p := range apiKeyProviders {
		item := systemApi.APIKeyItem{
			Provider:    p.Provider,
			DisplayName: p.DisplayName,
			Description: p.Description,
			Configured:  false,
			MaskedKey:   "",
		}

		// 先查数据库
		configKey := apiKeyPrefix + p.Provider
		var config entity.SystemConfig
		err := dao.SystemConfig.Ctx(ctx).
			Where(dao.SystemConfig.Columns().ConfigKey, configKey).
			WhereNull(dao.SystemConfig.Columns().DeletedAt).
			Scan(&config)
		if err == nil && config.ConfigValue != "" {
			// 解密并脱敏显示
			plainKey, decErr := crypto.DecryptAES(config.ConfigValue, masterKey)
			if decErr == nil && plainKey != "" {
				item.Configured = true
				item.MaskedKey = maskAPIKey(plainKey)
			}
		}

		// 数据库没有则检查环境变量
		if !item.Configured {
			if envVal := os.Getenv(p.EnvKey); envVal != "" {
				item.Configured = true
				item.MaskedKey = maskAPIKey(envVal) + " (env)"
			}
		}

		keys = append(keys, item)
	}

	return &systemApi.GetAPIKeysRes{Keys: keys}, nil
}

// UpdateAPIKey 更新指定服务商的 API Key（加密存储到数据库）
func (s *sSystem) UpdateAPIKey(ctx context.Context, provider string, apiKey string) error {
	// 验证 provider 合法性
	valid := false
	for _, p := range apiKeyProviders {
		if p.Provider == provider {
			valid = true
			break
		}
	}
	if !valid {
		return gerror.Newf("不支持的服务商: %s", provider)
	}

	masterKey, err := getMasterKey(ctx)
	if err != nil {
		return err
	}

	// 加密 API Key
	encrypted, err := crypto.EncryptAES(apiKey, masterKey)
	if err != nil {
		return gerror.Wrap(err, "加密 API Key 失败")
	}

	configKey := apiKeyPrefix + provider
	columns := dao.SystemConfig.Columns()
	now := gtime.Now()

	// Upsert：查询是否已存在
	count, err := dao.SystemConfig.Ctx(ctx).
		Where(columns.ConfigKey, configKey).
		WhereNull(columns.DeletedAt).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询 API Key 配置失败")
	}

	if count > 0 {
		_, err = dao.SystemConfig.Ctx(ctx).
			Where(columns.ConfigKey, configKey).
			WhereNull(columns.DeletedAt).
			Data(g.Map{
				columns.ConfigValue: encrypted,
				columns.UpdatedAt:   now,
			}).Update()
	} else {
		_, err = dao.SystemConfig.Ctx(ctx).Data(g.Map{
			columns.ConfigKey:   configKey,
			columns.ConfigValue: encrypted,
			columns.CreatedAt:   now,
			columns.UpdatedAt:   now,
		}).Insert()
	}

	if err != nil {
		return gerror.Wrap(err, "保存 API Key 失败")
	}
	return nil
}

// DeleteAPIKey 删除指定服务商的 API Key（软删除）
func (s *sSystem) DeleteAPIKey(ctx context.Context, provider string) error {
	configKey := apiKeyPrefix + provider
	columns := dao.SystemConfig.Columns()

	_, err := dao.SystemConfig.Ctx(ctx).
		Where(columns.ConfigKey, configKey).
		WhereNull(columns.DeletedAt).
		Data(g.Map{
			columns.DeletedAt: gtime.Now(),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "删除 API Key 失败")
	}
	return nil
}

// GetAPIKeyPlain 获取指定服务商的 API Key 明文
// 优先级：数据库 > 环境变量 > 空字符串
// 内部使用，不对外暴露
func (s *sSystem) GetAPIKeyPlain(ctx context.Context, provider string) string {
	// 先查数据库
	masterKey, err := getMasterKey(ctx)
	if err == nil {
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

	// 数据库没有则查环境变量
	for _, p := range apiKeyProviders {
		if p.Provider == provider {
			return os.Getenv(p.EnvKey)
		}
	}

	return ""
}
