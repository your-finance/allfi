// Package logic 用户模块业务逻辑
// 实现用户设置的读写和缓存清除
package logic

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	userApi "your-finance/allfi/api/v1/user"
	"your-finance/allfi/internal/consts"
	exchangeDao "your-finance/allfi/internal/app/exchange/dao"
	goalDao "your-finance/allfi/internal/app/goal/dao"
	manualAssetDao "your-finance/allfi/internal/app/manual_asset/dao"
	priceAlertDao "your-finance/allfi/internal/app/price_alert/dao"
	strategyDao "your-finance/allfi/internal/app/strategy/dao"
	"your-finance/allfi/internal/app/user/dao"
	"your-finance/allfi/internal/app/user/service"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	"your-finance/allfi/internal/model/entity"
	"your-finance/allfi/internal/utils"
)

// 用户设置 key 前缀
const userSettingPrefix = "user."

// sUser 用户服务实现
type sUser struct{}

// New 创建用户服务实例
func New() service.IUser {
	return &sUser{}
}

// GetSettings 获取用户设置
// 从 system_config 表查询所有 config_key 以 "user." 开头的记录
// 返回去掉前缀后的 key-value 映射
func (s *sUser) GetSettings(ctx context.Context) (*userApi.GetSettingsRes, error) {
	var configs []entity.SystemConfig
	err := dao.SystemConfig.Ctx(ctx).
		WhereLike(dao.SystemConfig.Columns().ConfigKey, userSettingPrefix+"%").
		WhereNull(dao.SystemConfig.Columns().DeletedAt).
		Scan(&configs)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户设置失败")
	}

	// 构建设置映射，去掉 "user." 前缀
	settings := make(map[string]string, len(configs))
	for _, cfg := range configs {
		key := strings.TrimPrefix(cfg.ConfigKey, userSettingPrefix)
		settings[key] = cfg.ConfigValue
	}

	return &userApi.GetSettingsRes{
		Settings: settings,
	}, nil
}

// UpdateSettings 更新用户设置
// 遍历传入的 settings map，对每个 key 执行 upsert 操作
// key 会自动添加 "user." 前缀存储到 system_config 表
func (s *sUser) UpdateSettings(ctx context.Context, settings map[string]string) error {
	if len(settings) == 0 {
		return nil
	}

	now := gtime.Now()
	columns := dao.SystemConfig.Columns()

	for key, value := range settings {
		configKey := userSettingPrefix + key

		// 查询是否已存在该配置项
		count, err := dao.SystemConfig.Ctx(ctx).
			Where(columns.ConfigKey, configKey).
			WhereNull(columns.DeletedAt).
			Count()
		if err != nil {
			return gerror.Wrapf(err, "查询设置项 [%s] 失败", key)
		}

		if count > 0 {
			// 已存在则更新
			_, err = dao.SystemConfig.Ctx(ctx).
				Where(columns.ConfigKey, configKey).
				WhereNull(columns.DeletedAt).
				Data(g.Map{
					columns.ConfigValue: value,
					columns.UpdatedAt:   now,
				}).
				Update()
		} else {
			// 不存在则插入
			_, err = dao.SystemConfig.Ctx(ctx).
				Data(g.Map{
					columns.ConfigKey:   configKey,
					columns.ConfigValue: value,
					columns.Description: "用户设置: " + key,
					columns.CreatedAt:   now,
					columns.UpdatedAt:   now,
				}).
				Insert()
		}

		if err != nil {
			return gerror.Wrapf(err, "更新设置项 [%s] 失败", key)
		}
	}

	g.Log().Info(ctx, "用户设置更新成功", "count", len(settings))
	return nil
}

// ResetSettings 重置所有用户设置
// 删除 system_config 表中所有 config_key 以 "user." 开头的记录
func (s *sUser) ResetSettings(ctx context.Context) error {
	_, err := dao.SystemConfig.Ctx(ctx).
		WhereLike(dao.SystemConfig.Columns().ConfigKey, userSettingPrefix+"%").
		Delete()
	if err != nil {
		return gerror.Wrap(err, "重置用户设置失败")
	}

	g.Log().Info(ctx, "用户设置已重置")
	return nil
}

// ClearCache 清除全局内存缓存
// 调用 utils.GetGlobalCache().Clear() 清空所有缓存项
func (s *sUser) ClearCache(ctx context.Context) error {
	cache := utils.GetGlobalCache()
	cache.Clear()

	g.Log().Info(ctx, "全局缓存已清除")
	return nil
}

// ExportData 导出用户数据
// 聚合查询各模块数据，不包含加密的 API Key/Secret/Passphrase 等敏感信息
func (s *sUser) ExportData(ctx context.Context) (*userApi.ExportDataRes, error) {
	res := &userApi.ExportDataRes{
		ExportedAt: time.Now().Format(time.RFC3339),
	}

	// ===== 交易所账户（不导出加密凭证） =====
	var accounts []entity.ExchangeAccounts
	err := exchangeDao.ExchangeAccounts.Ctx(ctx).
		Where(exchangeDao.ExchangeAccounts.Columns().UserId, consts.GetUserID(ctx)).
		Where(exchangeDao.ExchangeAccounts.Columns().IsActive, 1).
		WhereNull(exchangeDao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&accounts)
	if err != nil {
		g.Log().Warning(ctx, "导出交易所账户失败", "error", err)
	}
	res.ExchangeAccounts = make([]userApi.ExportExchangeAccount, 0, len(accounts))
	for _, acc := range accounts {
		res.ExchangeAccounts = append(res.ExchangeAccounts, userApi.ExportExchangeAccount{
			ID:           uint(acc.Id),
			ExchangeName: acc.ExchangeName,
			Label:        acc.Label,
			Note:         acc.Note,
		})
	}

	// ===== 钱包地址 =====
	var wallets []entity.WalletAddresses
	err = walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, consts.GetUserID(ctx)).
		Where(walletDao.WalletAddresses.Columns().IsActive, 1).
		WhereNull(walletDao.WalletAddresses.Columns().DeletedAt).
		Scan(&wallets)
	if err != nil {
		g.Log().Warning(ctx, "导出钱包地址失败", "error", err)
	}
	res.WalletAddresses = make([]userApi.ExportWalletAddress, 0, len(wallets))
	for _, w := range wallets {
		res.WalletAddresses = append(res.WalletAddresses, userApi.ExportWalletAddress{
			ID:         uint(w.Id),
			Blockchain: w.Blockchain,
			Address:    w.Address,
			Label:      w.Label,
		})
	}

	// ===== 手动资产 =====
	var manualAssets []entity.ManualAssets
	err = manualAssetDao.ManualAssets.Ctx(ctx).
		Where(manualAssetDao.ManualAssets.Columns().UserId, consts.GetUserID(ctx)).
		Where(manualAssetDao.ManualAssets.Columns().IsActive, 1).
		WhereNull(manualAssetDao.ManualAssets.Columns().DeletedAt).
		Scan(&manualAssets)
	if err != nil {
		g.Log().Warning(ctx, "导出手动资产失败", "error", err)
	}
	res.ManualAssets = make([]userApi.ExportManualAsset, 0, len(manualAssets))
	for _, ma := range manualAssets {
		res.ManualAssets = append(res.ManualAssets, userApi.ExportManualAsset{
			ID:        uint(ma.Id),
			AssetType: ma.AssetType,
			AssetName: ma.AssetName,
			Amount:    ma.Amount,
			AmountUSD: ma.AmountUsd,
			Currency:  ma.Currency,
			Notes:     ma.Notes,
		})
	}

	// ===== 策略 =====
	var strategies []entity.Strategies
	err = strategyDao.Strategies.Ctx(ctx).
		Where(strategyDao.Strategies.Columns().UserId, consts.GetUserID(ctx)).
		WhereNull(strategyDao.Strategies.Columns().DeletedAt).
		Scan(&strategies)
	if err != nil {
		g.Log().Warning(ctx, "导出策略失败", "error", err)
	}
	res.Strategies = make([]userApi.ExportStrategy, 0, len(strategies))
	for _, st := range strategies {
		res.Strategies = append(res.Strategies, userApi.ExportStrategy{
			ID:       uint(st.Id),
			Name:     st.Name,
			Type:     st.Type,
			Config:   st.Config,
			IsActive: st.IsActive == 1,
		})
	}

	// ===== 目标 =====
	var goals []entity.Goals
	err = goalDao.Goals.Ctx(ctx).
		WhereNull(goalDao.Goals.Columns().DeletedAt).
		Scan(&goals)
	if err != nil {
		g.Log().Warning(ctx, "导出目标失败", "error", err)
	}
	res.Goals = make([]userApi.ExportGoal, 0, len(goals))
	for _, gl := range goals {
		deadline := ""
		if !gl.Deadline.IsZero() {
			deadline = gl.Deadline.Format("2006-01-02")
		}
		res.Goals = append(res.Goals, userApi.ExportGoal{
			ID:          uint(gl.Id),
			Title:       gl.Title,
			Type:        gl.Type,
			TargetValue: float64(gl.TargetValue),
			Currency:    gl.Currency,
			Deadline:    deadline,
		})
	}

	// ===== 价格预警 =====
	var alerts []entity.PriceAlerts
	err = priceAlertDao.PriceAlerts.Ctx(ctx).
		Where(priceAlertDao.PriceAlerts.Columns().UserId, consts.GetUserID(ctx)).
		WhereNull(priceAlertDao.PriceAlerts.Columns().DeletedAt).
		Scan(&alerts)
	if err != nil {
		g.Log().Warning(ctx, "导出价格预警失败", "error", err)
	}
	res.PriceAlerts = make([]userApi.ExportPriceAlert, 0, len(alerts))
	for _, a := range alerts {
		res.PriceAlerts = append(res.PriceAlerts, userApi.ExportPriceAlert{
			ID:          uint(a.Id),
			Symbol:      a.Symbol,
			Condition:   a.Condition,
			TargetPrice: float64(a.TargetPrice),
			IsActive:    a.IsActive == 1,
			Note:        a.Note,
		})
	}

	// ===== 用户设置 =====
	settingsRes, err := s.GetSettings(ctx)
	if err != nil {
		g.Log().Warning(ctx, "导出用户设置失败", "error", err)
		res.Settings = make(map[string]string)
	} else {
		res.Settings = settingsRes.Settings
	}

	g.Log().Info(ctx, "用户数据导出完成",
		"exchangeAccounts", len(res.ExchangeAccounts),
		"walletAddresses", len(res.WalletAddresses),
		"manualAssets", len(res.ManualAssets),
		"strategies", len(res.Strategies),
		"goals", len(res.Goals),
		"priceAlerts", len(res.PriceAlerts),
	)

	return res, nil
}
