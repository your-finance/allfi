// Package logic 成就系统业务逻辑
// 管理成就定义、查询解锁状态、检查解锁条件并持久化记录
package logic

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/achievement/dao"
	"your-finance/allfi/internal/app/achievement/model"
	"your-finance/allfi/internal/app/achievement/service"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	exchangeDao "your-finance/allfi/internal/app/exchange/dao"
	priceAlertDao "your-finance/allfi/internal/app/price_alert/dao"
	reportDao "your-finance/allfi/internal/app/report/dao"
	strategyDao "your-finance/allfi/internal/app/strategy/dao"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	"your-finance/allfi/internal/model/entity"
)

// sAchievement 成就系统服务实现
type sAchievement struct{}

// New 创建成就系统服务实例
func New() service.IAchievement {
	return &sAchievement{}
}

// achievementDef 成就定义（内部使用）
type achievementDef struct {
	ID          string
	Name        string
	Description string
	Category    string
	Icon        string
}

// allAchievementDefs 所有成就定义（硬编码）
var allAchievementDefs = []achievementDef{
	// === 里程碑类 ===
	{ID: "first_account", Name: "初学者", Description: "添加第一个交易所账户或钱包地址", Category: model.CategoryMilestone, Icon: "rocket"},
	{ID: "multi_source", Name: "多元投资者", Description: "同时拥有交易所账户和钱包地址", Category: model.CategoryMilestone, Icon: "layers"},
	{ID: "collector_10", Name: "收藏家", Description: "持有 10 种以上不同资产", Category: model.CategoryMilestone, Icon: "grid"},
	{ID: "hodl_1btc", Name: "比特信仰者", Description: "持有至少 1 BTC", Category: model.CategoryMilestone, Icon: "bitcoin"},
	{ID: "whale", Name: "巨鲸", Description: "总资产超过 $100,000", Category: model.CategoryMilestone, Icon: "whale"},
	// === 坚持类 ===
	{ID: "first_strategy", Name: "策略师", Description: "创建第一个自动化策略", Category: model.CategoryPersistence, Icon: "target"},
	{ID: "alert_setter", Name: "警觉守卫", Description: "设置 3 个以上价格预警", Category: model.CategoryPersistence, Icon: "bell"},
	{ID: "data_master", Name: "数据大师", Description: "生成 5 份以上报告", Category: model.CategoryPersistence, Icon: "chart"},
	// === 投资类 ===
	{ID: "defi_explorer", Name: "DeFi 探索者", Description: "拥有至少 1 个 DeFi 仓位", Category: model.CategoryInvestment, Icon: "compass"},
	{ID: "nft_collector", Name: "NFT 收藏家", Description: "持有至少 1 个 NFT", Category: model.CategoryInvestment, Icon: "image"},
	{ID: "diversified", Name: "分散投资", Description: "资产覆盖 CEX、链上钱包和手动资产三种来源", Category: model.CategoryInvestment, Icon: "pie-chart"},
}

// GetAchievements 获取成就列表（含解锁状态）
func (s *sAchievement) GetAchievements(ctx context.Context, userID uint) ([]model.AchievementStatus, error) {
	// 直接从 DAO 查询用户已解锁的成就记录
	var unlocked []entity.UserAchievements
	err := dao.UserAchievements.Ctx(ctx).
		Where(dao.UserAchievements.Columns().UserId, userID).
		Scan(&unlocked)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户成就记录失败")
	}

	// 构建已解锁成就 map
	unlockedMap := make(map[string]time.Time, len(unlocked))
	for _, u := range unlocked {
		unlockedMap[u.AchievementId] = u.UnlockedAt
	}

	// 合并成就定义和解锁状态
	result := make([]model.AchievementStatus, 0, len(allAchievementDefs))
	for _, def := range allAchievementDefs {
		as := model.AchievementStatus{
			ID:          def.ID,
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Category:    def.Category,
			IsUnlocked:  false,
			Progress:    0,
		}
		if unlockedAt, ok := unlockedMap[def.ID]; ok {
			as.IsUnlocked = true
			as.Progress = 100
			as.UnlockedAt = unlockedAt.Format(time.RFC3339)
		}
		result = append(result, as)
	}

	return result, nil
}

// CheckAndUnlock 检查并解锁新成就
func (s *sAchievement) CheckAndUnlock(ctx context.Context, userID uint) ([]model.UnlockedAchievement, error) {
	// 查询用户已解锁的成就
	var unlocked []entity.UserAchievements
	err := dao.UserAchievements.Ctx(ctx).
		Where(dao.UserAchievements.Columns().UserId, userID).
		Scan(&unlocked)
	if err != nil {
		return nil, gerror.Wrap(err, "查询用户成就失败")
	}

	unlockedSet := make(map[string]bool, len(unlocked))
	for _, u := range unlocked {
		unlockedSet[u.AchievementId] = true
	}

	// 检查各成就条件
	var newlyUnlocked []model.UnlockedAchievement

	// 检查 first_account：有交易所账户或钱包
	if !unlockedSet["first_account"] {
		exchangeCount, _ := exchangeDao.ExchangeAccounts.Ctx(ctx).Where(exchangeDao.ExchangeAccounts.Columns().UserId, userID).Count()
		walletCount, _ := walletDao.WalletAddresses.Ctx(ctx).Where(walletDao.WalletAddresses.Columns().UserId, userID).Count()
		if exchangeCount > 0 || walletCount > 0 {
			s.unlock(ctx, userID, "first_account", &newlyUnlocked)
		}
	}

	// 检查 multi_source：同时有交易所和钱包
	if !unlockedSet["multi_source"] {
		exchangeCount, _ := exchangeDao.ExchangeAccounts.Ctx(ctx).Where(exchangeDao.ExchangeAccounts.Columns().UserId, userID).Count()
		walletCount, _ := walletDao.WalletAddresses.Ctx(ctx).Where(walletDao.WalletAddresses.Columns().UserId, userID).Count()
		if exchangeCount > 0 && walletCount > 0 {
			s.unlock(ctx, userID, "multi_source", &newlyUnlocked)
		}
	}

	// 检查 collector_10：持有 10+ 种不同资产
	if !unlockedSet["collector_10"] {
		type symbolResult struct {
			AssetSymbol string
		}
		var symbols []symbolResult
		_ = assetDao.AssetDetails.Ctx(ctx).
			Where(assetDao.AssetDetails.Columns().UserId, userID).
			Fields(assetDao.AssetDetails.Columns().AssetSymbol).
			Distinct().
			Scan(&symbols)
		if len(symbols) >= 10 {
			s.unlock(ctx, userID, "collector_10", &newlyUnlocked)
		}
	}

	// 检查 whale：总资产 > $100,000
	if !unlockedSet["whale"] {
		var totalValue float64
		v, _ := assetDao.AssetDetails.Ctx(ctx).
			Where(assetDao.AssetDetails.Columns().UserId, userID).
			Sum(assetDao.AssetDetails.Columns().ValueUsd)
		totalValue = v
		if totalValue > 100000 {
			s.unlock(ctx, userID, "whale", &newlyUnlocked)
		}
	}

	// 检查 first_strategy：有策略
	if !unlockedSet["first_strategy"] {
		count, _ := strategyDao.Strategies.Ctx(ctx).Where(strategyDao.Strategies.Columns().UserId, userID).Count()
		if count > 0 {
			s.unlock(ctx, userID, "first_strategy", &newlyUnlocked)
		}
	}

	// 检查 alert_setter：3+ 个价格预警
	if !unlockedSet["alert_setter"] {
		count, _ := priceAlertDao.PriceAlerts.Ctx(ctx).Where(priceAlertDao.PriceAlerts.Columns().UserId, userID).Count()
		if count >= 3 {
			s.unlock(ctx, userID, "alert_setter", &newlyUnlocked)
		}
	}

	// 检查 data_master：5+ 份报告
	if !unlockedSet["data_master"] {
		count, _ := reportDao.Reports.Ctx(ctx).Where(reportDao.Reports.Columns().UserId, userID).Count()
		if count >= 5 {
			s.unlock(ctx, userID, "data_master", &newlyUnlocked)
		}
	}

	if len(newlyUnlocked) > 0 {
		g.Log().Info(ctx, "成就检查完成，新解锁成就数",
			"userID", userID,
			"count", len(newlyUnlocked),
		)
	}

	return newlyUnlocked, nil
}

// unlock 解锁指定成就并添加到结果列表
func (s *sAchievement) unlock(ctx context.Context, userID uint, achievementID string, result *[]model.UnlockedAchievement) {
	// 插入解锁记录
	_, err := dao.UserAchievements.Ctx(ctx).Insert(g.Map{
		dao.UserAchievements.Columns().UserId:        userID,
		dao.UserAchievements.Columns().AchievementId: achievementID,
		dao.UserAchievements.Columns().UnlockedAt:    time.Now(),
	})
	if err != nil {
		g.Log().Warning(ctx, "插入成就解锁记录失败", "achievementID", achievementID, "error", err)
		return
	}

	// 查找成就定义
	for _, def := range allAchievementDefs {
		if def.ID == achievementID {
			*result = append(*result, model.UnlockedAchievement{
				ID:          def.ID,
				Name:        def.Name,
				Description: def.Description,
				Icon:        def.Icon,
				Category:    def.Category,
			})
			break
		}
	}
}
