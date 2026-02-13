// Package logic 交易记录业务逻辑
// 提供交易记录查询、同步、统计、同步设置管理
package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	transactionApi "your-finance/allfi/api/v1/transaction"
	"your-finance/allfi/internal/consts"
	exchangeDao "your-finance/allfi/internal/app/exchange/dao"
	"your-finance/allfi/internal/app/transaction/dao"
	txModel "your-finance/allfi/internal/app/transaction/model"
	"your-finance/allfi/internal/app/transaction/service"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	"your-finance/allfi/internal/integrations"
	"your-finance/allfi/internal/integrations/binance"
	"your-finance/allfi/internal/integrations/coinbase"
	"your-finance/allfi/internal/integrations/okx"
	"your-finance/allfi/internal/model/entity"
	"your-finance/allfi/utility/crypto"
)

// sTransaction 交易记录服务实现
type sTransaction struct{}

// New 创建交易记录服务实例
func New() service.ITransaction {
	return &sTransaction{}
}

// List 获取交易记录列表
func (s *sTransaction) List(ctx context.Context, page, pageSize int, source, txType, start, end, cursor string) (*transactionApi.ListRes, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := dao.UnifiedTransactions.Ctx(ctx).
		Where(dao.UnifiedTransactions.Columns().UserId, consts.GetUserID(ctx))

	// 筛选条件
	if source != "" {
		query = query.Where(dao.UnifiedTransactions.Columns().Source, source)
	}
	if txType != "" {
		query = query.Where(dao.UnifiedTransactions.Columns().TxType, txType)
	}
	if start != "" {
		startTime, err := time.Parse("2006-01-02", start)
		if err == nil {
			query = query.WhereGTE(dao.UnifiedTransactions.Columns().Timestamp, startTime)
		}
	}
	if end != "" {
		endTime, err := time.Parse("2006-01-02", end)
		if err == nil {
			query = query.WhereLTE(dao.UnifiedTransactions.Columns().Timestamp, endTime.Add(24*time.Hour))
		}
	}
	if cursor != "" {
		cursorTime, err := time.Parse(time.RFC3339, cursor)
		if err == nil {
			query = query.WhereLT(dao.UnifiedTransactions.Columns().Timestamp, cursorTime)
		}
	}

	// 查询总数和分页数据
	var transactions []entity.UnifiedTransactions
	var total int
	err := query.
		OrderDesc(dao.UnifiedTransactions.Columns().Timestamp).
		Page(page, pageSize).
		ScanAndCount(&transactions, &total, true)
	if err != nil {
		return nil, gerror.Wrap(err, "查询交易记录失败")
	}

	// 转换为 API 格式
	items := make([]transactionApi.TransactionItem, 0, len(transactions))
	for _, tx := range transactions {
		items = append(items, toTransactionItem(&tx))
	}

	return &transactionApi.ListRes{
		Transactions: items,
		Total:        int64(total),
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

// Sync 触发交易记录同步
// 1. 查询所有活跃的交易所账户，解密凭证后拉取交易历史
// 2. 将交易记录写入 unified_transactions 表（按 source_id 去重）
// 3. 更新 sync_metadata 和系统配置中的同步时间
func (s *sTransaction) Sync(ctx context.Context) (*transactionApi.SyncRes, error) {
	syncedCount := 0

	// ===== 1. 同步交易所交易记录 =====
	var accounts []entity.ExchangeAccounts
	err := exchangeDao.ExchangeAccounts.Ctx(ctx).
		Where(exchangeDao.ExchangeAccounts.Columns().UserId, consts.GetUserID(ctx)).
		Where(exchangeDao.ExchangeAccounts.Columns().IsActive, 1).
		WhereNull(exchangeDao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&accounts)
	if err != nil {
		g.Log().Warning(ctx, "查询交易所账户失败", "error", err)
	}

	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()

	for _, acc := range accounts {
		g.Log().Info(ctx, "开始同步交易所交易记录", "exchange", acc.ExchangeName, "accountId", acc.Id)

		// 解密 API 凭证
		apiKey, err := crypto.DecryptAES(acc.ApiKeyEncrypted, masterKey)
		if err != nil {
			g.Log().Warning(ctx, "解密 API Key 失败，跳过该账户", "accountId", acc.Id, "error", err)
			continue
		}
		apiSecret, err := crypto.DecryptAES(acc.ApiSecretEncrypted, masterKey)
		if err != nil {
			g.Log().Warning(ctx, "解密 API Secret 失败，跳过该账户", "accountId", acc.Id, "error", err)
			continue
		}
		var passphrase string
		if acc.ApiPassphraseEncrypted != "" {
			passphrase, _ = crypto.DecryptAES(acc.ApiPassphraseEncrypted, masterKey)
		}

		// 创建交易所客户端
		var client integrations.ExchangeClient
		switch strings.ToLower(acc.ExchangeName) {
		case "binance":
			client = binance.NewClient(apiKey, apiSecret)
		case "okx":
			client = okx.NewClient(apiKey, apiSecret, passphrase)
		case "coinbase":
			client = coinbase.NewClient(apiKey, apiSecret)
		default:
			g.Log().Warning(ctx, "暂不支持的交易所，跳过同步", "exchange", acc.ExchangeName)
			continue
		}

		// 查询该交易所的同步元数据获取上次同步时间
		var meta entity.SyncMetadata
		_ = dao.SyncMetadata.Ctx(ctx).
			Where(dao.SyncMetadata.Columns().Source, acc.ExchangeName).
			Scan(&meta)

		// 构建查询参数：从上次同步时间到当前
		tradeParams := integrations.TradeHistoryParams{
			StartTime: meta.LastSyncTime,
			EndTime:   time.Now(),
			Limit:     500,
		}
		dwParams := integrations.DepositWithdrawParams{
			StartTime: meta.LastSyncTime,
			EndTime:   time.Now(),
			Limit:     500,
		}

		// 拉取交易历史
		trades, err := client.GetTradeHistory(ctx, tradeParams)
		if err != nil {
			g.Log().Warning(ctx, "获取交易历史失败", "exchange", acc.ExchangeName, "error", err)
		} else {
			count := s.saveTrades(ctx, trades)
			syncedCount += count
			g.Log().Info(ctx, "交易记录同步完成", "exchange", acc.ExchangeName, "count", count)
		}

		// 拉取充值历史
		deposits, err := client.GetDepositHistory(ctx, dwParams)
		if err != nil {
			g.Log().Warning(ctx, "获取充值历史失败", "exchange", acc.ExchangeName, "error", err)
		} else {
			count := s.saveTransfers(ctx, deposits)
			syncedCount += count
			g.Log().Info(ctx, "充值记录同步完成", "exchange", acc.ExchangeName, "count", count)
		}

		// 拉取提现历史
		withdraws, err := client.GetWithdrawHistory(ctx, dwParams)
		if err != nil {
			g.Log().Warning(ctx, "获取提现历史失败", "exchange", acc.ExchangeName, "error", err)
		} else {
			count := s.saveTransfers(ctx, withdraws)
			syncedCount += count
			g.Log().Info(ctx, "提现记录同步完成", "exchange", acc.ExchangeName, "count", count)
		}

		// 更新同步元数据
		s.upsertSyncMeta(ctx, acc.ExchangeName, syncedCount)
	}

	// ===== 2. 记录钱包地址同步状态 =====
	var wallets []entity.WalletAddresses
	err = walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, consts.GetUserID(ctx)).
		Where(walletDao.WalletAddresses.Columns().IsActive, 1).
		WhereNull(walletDao.WalletAddresses.Columns().DeletedAt).
		Scan(&wallets)
	if err != nil {
		g.Log().Warning(ctx, "查询钱包地址失败", "error", err)
	}
	for _, w := range wallets {
		g.Log().Info(ctx, "钱包地址已记录（链上交易同步待扩展）",
			"blockchain", w.Blockchain,
			"address", w.Address,
		)
		// 更新钱包对应链的同步元数据
		s.upsertSyncMeta(ctx, w.Blockchain, 0)
	}

	// ===== 3. 更新最后同步时间到系统配置 =====
	s.upsertConfig(ctx, txModel.ConfigKeyLastSyncAt, time.Now().Format(time.RFC3339))

	return &transactionApi.SyncRes{
		Message:     "同步完成",
		SyncedCount: syncedCount,
	}, nil
}

// saveTrades 保存交易记录到 unified_transactions 表
// 按 source_id 去重，返回新增记录数
func (s *sTransaction) saveTrades(ctx context.Context, trades []integrations.Trade) int {
	saved := 0
	now := time.Now()
	cols := dao.UnifiedTransactions.Columns()

	for _, t := range trades {
		sourceID := t.Source + ":" + t.ID
		// 去重检查
		count, _ := dao.UnifiedTransactions.Ctx(ctx).
			Where(cols.SourceId, sourceID).
			Count()
		if count > 0 {
			continue
		}

		// 计算 USD 价值（价格 * 数量）
		valueUSD := t.Price * t.Quantity

		_, err := dao.UnifiedTransactions.Ctx(ctx).Data(g.Map{
			cols.UserId:     consts.GetUserID(ctx),
			cols.TxType:     t.Side,
			cols.Source:     t.Source,
			cols.SourceId:   sourceID,
			cols.FromAsset:  t.Symbol,
			cols.FromAmount: t.Quantity,
			cols.Fee:        t.Fee,
			cols.FeeCoin:    t.FeeCoin,
			cols.ValueUsd:   valueUSD,
			cols.Timestamp:  t.Timestamp,
			cols.CreatedAt:  now,
			cols.UpdatedAt:  now,
		}).Insert()
		if err != nil {
			g.Log().Warning(ctx, "保存交易记录失败", "sourceId", sourceID, "error", err)
			continue
		}
		saved++
	}
	return saved
}

// saveTransfers 保存充提记录到 unified_transactions 表
// 按 source_id 去重，返回新增记录数
func (s *sTransaction) saveTransfers(ctx context.Context, transfers []integrations.Transfer) int {
	saved := 0
	now := time.Now()
	cols := dao.UnifiedTransactions.Columns()

	for _, t := range transfers {
		sourceID := t.Source + ":" + t.Type + ":" + t.ID
		// 去重检查
		count, _ := dao.UnifiedTransactions.Ctx(ctx).
			Where(cols.SourceId, sourceID).
			Count()
		if count > 0 {
			continue
		}

		_, err := dao.UnifiedTransactions.Ctx(ctx).Data(g.Map{
			cols.UserId:     consts.GetUserID(ctx),
			cols.TxType:     t.Type,
			cols.Source:     t.Source,
			cols.SourceId:   sourceID,
			cols.FromAsset:  t.Coin,
			cols.FromAmount: t.Amount,
			cols.Fee:        t.Fee,
			cols.FeeCoin:    t.Coin,
			cols.TxHash:     t.TxHash,
			cols.Timestamp:  t.Timestamp,
			cols.CreatedAt:  now,
			cols.UpdatedAt:  now,
		}).Insert()
		if err != nil {
			g.Log().Warning(ctx, "保存充提记录失败", "sourceId", sourceID, "error", err)
			continue
		}
		saved++
	}
	return saved
}

// upsertSyncMeta 更新或插入同步元数据
func (s *sTransaction) upsertSyncMeta(ctx context.Context, source string, txCount int) {
	now := time.Now()
	cols := dao.SyncMetadata.Columns()

	// 先尝试更新
	result, err := dao.SyncMetadata.Ctx(ctx).
		Where(cols.Source, source).
		Data(g.Map{
			cols.LastSyncTime: now,
			cols.TxCount:      txCount,
			cols.UpdatedAt:    now,
		}).Update()
	if err != nil {
		g.Log().Warning(ctx, "更新同步元数据失败", "source", source, "error", err)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		// 不存在则插入
		_, err = dao.SyncMetadata.Ctx(ctx).Data(g.Map{
			cols.Source:       source,
			cols.LastSyncTime: now,
			cols.TxCount:      txCount,
			cols.CreatedAt:    now,
			cols.UpdatedAt:    now,
		}).Insert()
		if err != nil {
			g.Log().Warning(ctx, "插入同步元数据失败", "source", source, "error", err)
		}
	}
}

// GetStats 获取交易统计
// 使用 SQL 聚合查询替代全表扫描，避免将大量记录加载到内存
func (s *sTransaction) GetStats(ctx context.Context) (*transactionApi.GetStatsRes, error) {
	cols := dao.UnifiedTransactions.Columns()
	baseQuery := dao.UnifiedTransactions.Ctx(ctx).
		Where(cols.UserId, consts.GetUserID(ctx))

	// 1. 总交易数
	total, err := baseQuery.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计交易总数失败")
	}

	// 2. 使用 SQL SUM 聚合计算总交易量和总手续费
	type SumResult struct {
		TotalVolume float64 `json:"total_volume"`
		TotalFees   float64 `json:"total_fees"`
	}
	var sumResult SumResult
	err = dao.UnifiedTransactions.Ctx(ctx).
		Where(cols.UserId, consts.GetUserID(ctx)).
		Fields("COALESCE(SUM(value_usd), 0) as total_volume, COALESCE(SUM(fee), 0) as total_fees").
		Scan(&sumResult)
	if err != nil {
		return nil, gerror.Wrap(err, "聚合交易量和手续费失败")
	}

	// 3. 按交易类型分组统计
	type CountResult struct {
		Key string `json:"key"`
		Cnt int    `json:"cnt"`
	}
	var typeResults []CountResult
	err = dao.UnifiedTransactions.Ctx(ctx).
		Where(cols.UserId, consts.GetUserID(ctx)).
		Fields("tx_type as `key`, COUNT(*) as cnt").
		Group("tx_type").
		Scan(&typeResults)
	if err != nil {
		return nil, gerror.Wrap(err, "按类型统计失败")
	}
	byType := make(map[string]int, len(typeResults))
	for _, r := range typeResults {
		byType[r.Key] = r.Cnt
	}

	// 4. 按来源分组统计
	var sourceResults []CountResult
	err = dao.UnifiedTransactions.Ctx(ctx).
		Where(cols.UserId, consts.GetUserID(ctx)).
		Fields("source as `key`, COUNT(*) as cnt").
		Group("source").
		Scan(&sourceResults)
	if err != nil {
		return nil, gerror.Wrap(err, "按来源统计失败")
	}
	bySource := make(map[string]int, len(sourceResults))
	for _, r := range sourceResults {
		bySource[r.Key] = r.Cnt
	}

	return &transactionApi.GetStatsRes{
		TotalTransactions: total,
		TotalVolume:       sumResult.TotalVolume,
		TotalFees:         sumResult.TotalFees,
		ByType:            byType,
		BySource:          bySource,
	}, nil
}

// GetSyncSettings 获取同步设置
func (s *sTransaction) GetSyncSettings(ctx context.Context) (*transactionApi.GetSyncSettingsRes, error) {
	settings := &transactionApi.SyncSettingsItem{
		AutoSync:     txModel.DefaultAutoSync,
		SyncInterval: txModel.DefaultSyncInterval,
	}

	// 从系统配置读取
	var configs []entity.SystemConfig
	err := dao.SystemConfig.Ctx(ctx).
		WhereIn(dao.SystemConfig.Columns().ConfigKey, g.Slice{
			txModel.ConfigKeyAutoSync,
			txModel.ConfigKeySyncInterval,
			txModel.ConfigKeyLastSyncAt,
		}).Scan(&configs)
	if err != nil {
		g.Log().Warning(ctx, "读取同步设置失败", "error", err)
		return &transactionApi.GetSyncSettingsRes{Settings: settings}, nil
	}

	for _, cfg := range configs {
		switch cfg.ConfigKey {
		case txModel.ConfigKeyAutoSync:
			settings.AutoSync = cfg.ConfigValue == "true" || cfg.ConfigValue == "1"
		case txModel.ConfigKeySyncInterval:
			if v, err := strconv.Atoi(cfg.ConfigValue); err == nil {
				settings.SyncInterval = v
			}
		case txModel.ConfigKeyLastSyncAt:
			settings.LastSyncAt = cfg.ConfigValue
		}
	}

	return &transactionApi.GetSyncSettingsRes{Settings: settings}, nil
}

// UpdateSyncSettings 更新同步设置
func (s *sTransaction) UpdateSyncSettings(ctx context.Context, autoSync *bool, syncInterval *int) (*transactionApi.UpdateSyncSettingsRes, error) {
	if autoSync != nil {
		value := "false"
		if *autoSync {
			value = "true"
		}
		s.upsertConfig(ctx, txModel.ConfigKeyAutoSync, value)
	}

	if syncInterval != nil {
		s.upsertConfig(ctx, txModel.ConfigKeySyncInterval, strconv.Itoa(*syncInterval))
	}

	// 返回最新设置
	return s.getUpdatedSettings(ctx)
}

// upsertConfig 更新或插入系统配置
func (s *sTransaction) upsertConfig(ctx context.Context, key, value string) {
	// 先尝试更新
	result, err := dao.SystemConfig.Ctx(ctx).
		Where(dao.SystemConfig.Columns().ConfigKey, key).
		Data(g.Map{
			dao.SystemConfig.Columns().ConfigValue: value,
			dao.SystemConfig.Columns().UpdatedAt:   time.Now(),
		}).Update()
	if err != nil {
		g.Log().Warning(ctx, "更新系统配置失败", "key", key, "error", err)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		// 不存在则插入
		_, err = dao.SystemConfig.Ctx(ctx).Insert(g.Map{
			dao.SystemConfig.Columns().ConfigKey:   key,
			dao.SystemConfig.Columns().ConfigValue: value,
			dao.SystemConfig.Columns().CreatedAt:   time.Now(),
			dao.SystemConfig.Columns().UpdatedAt:   time.Now(),
		})
		if err != nil {
			g.Log().Warning(ctx, "插入系统配置失败", "key", key, "error", err)
		}
	}
}

// getUpdatedSettings 获取更新后的同步设置
func (s *sTransaction) getUpdatedSettings(ctx context.Context) (*transactionApi.UpdateSyncSettingsRes, error) {
	getRes, err := s.GetSyncSettings(ctx)
	if err != nil {
		return nil, err
	}
	return &transactionApi.UpdateSyncSettingsRes{Settings: getRes.Settings}, nil
}

// toTransactionItem 将实体转换为 API 条目
// Price = ValueUSD / Amount（当 Amount 不为 0 时）
func toTransactionItem(tx *entity.UnifiedTransactions) transactionApi.TransactionItem {
	// 计算单价：ValueUSD / Amount
	var price float64
	amount := float64(tx.FromAmount)
	if amount != 0 {
		price = float64(tx.ValueUsd) / amount
	}

	return transactionApi.TransactionItem{
		ID:        uint(tx.Id),
		Source:    tx.Source,
		TxType:    tx.TxType,
		Symbol:    tx.FromAsset,
		Amount:    amount,
		Price:     price,
		Total:     float64(tx.ValueUsd),
		Fee:       float64(tx.Fee),
		FeeCoin:   tx.FeeCoin,
		TxHash:    tx.TxHash,
		Timestamp: tx.Timestamp.Format(time.RFC3339),
	}
}
