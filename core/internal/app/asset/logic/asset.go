// Package logic 资产业务逻辑
// 聚合 CEX、区块链、手动资产，计算总值和分类统计
package logic

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	assetApi "your-finance/allfi/api/v1/asset"
	"your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/asset/service"
	"your-finance/allfi/internal/consts"
	exchangeDao "your-finance/allfi/internal/app/exchange/dao"
	exchangeService "your-finance/allfi/internal/app/exchange/service"
	manualAssetDao "your-finance/allfi/internal/app/manual_asset/dao"
	walletDao "your-finance/allfi/internal/app/wallet/dao"
	walletService "your-finance/allfi/internal/app/wallet/service"
	exchangeRateDao "your-finance/allfi/internal/app/exchange_rate/dao"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/model/entity"
)

// sAsset 资产服务实现
type sAsset struct{}

// New 创建资产服务实例
func New() service.IAsset {
	return &sAsset{}
}

// GetSummary 获取资产概览
//
// 功能说明:
// 1. 从 asset_details 表聚合所有来源的资产
// 2. 按来源分类统计（cex/blockchain/manual）
// 3. 返回总价值和分类统计
//
// 参数:
//   - ctx: 上下文
//   - currency: 计价货币（USD/BTC/ETH/CNY）
func (s *sAsset) GetSummary(ctx context.Context, currency string) (*assetApi.GetSummaryRes, error) {
	// 设置默认货币
	if currency == "" {
		currency = "USD"
	}

	// 查询所有资产明细
	var assets []entity.AssetDetails
	err := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
		Scan(&assets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产明细失败")
	}

	// 聚合统计
	res := &assetApi.GetSummaryRes{
		Currency:  currency,
		BySource:  make(map[string]float64),
		UpdatedAt: gtime.Now().Format("Y-m-d H:i:s"),
	}

	// 获取货币转换汇率（USD -> 目标货币）
	convRate := s.getConversionRate(ctx, currency)

	var totalValue float64
	for _, asset := range assets {
		value := asset.ValueUsd * convRate

		totalValue += value
		res.BySource[asset.SourceType] += value
	}

	res.TotalValue = totalValue

	g.Log().Info(ctx, "获取资产概览成功",
		"currency", currency,
		"totalValue", totalValue,
		"assetCount", len(assets),
	)

	return res, nil
}

// GetDetails 获取资产明细列表
//
// 功能说明:
// 1. 从 asset_details 表查询资产明细
// 2. 支持按来源类型筛选
// 3. 按价值降序排列
func (s *sAsset) GetDetails(ctx context.Context, sourceType string, currency string) (*assetApi.GetDetailsRes, error) {
	// 设置默认货币
	if currency == "" {
		currency = "USD"
	}

	// 构建查询
	query := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, consts.GetUserID(ctx))

	// 按来源类型筛选
	if sourceType != "" {
		query = query.Where(dao.AssetDetails.Columns().SourceType, sourceType)
	}

	// 执行查询
	var assets []entity.AssetDetails
	err := query.OrderDesc(dao.AssetDetails.Columns().ValueUsd).Scan(&assets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产明细失败")
	}

	// 转换为 API 响应格式
	items := make([]assetApi.AssetDetailItem, 0, len(assets))
	for _, asset := range assets {
		items = append(items, assetApi.AssetDetailItem{
			ID:         uint(asset.Id),
			Symbol:     asset.AssetSymbol,
			Amount:     asset.Balance,
			Value:      asset.ValueUsd,
			Price:      asset.PriceUsd,
			Source:     asset.AssetName,
			SourceType: asset.SourceType,
			UpdatedAt:  asset.LastUpdated.Format("2006-01-02 15:04:05"),
		})
	}

	g.Log().Info(ctx, "获取资产明细成功",
		"sourceType", sourceType,
		"currency", currency,
		"count", len(items),
	)

	return &assetApi.GetDetailsRes{Assets: items}, nil
}

// GetHistory 获取资产历史趋势
//
// 功能说明:
// 1. 从 asset_snapshots 表查询指定天数内的快照
// 2. 按时间升序排列
func (s *sAsset) GetHistory(ctx context.Context, days int, currency string) (*assetApi.GetHistoryRes, error) {
	// 设置默认值
	if days <= 0 {
		days = 30
	}
	if currency == "" {
		currency = "USD"
	}

	// 计算起始时间
	startTime := time.Now().AddDate(0, 0, -days)

	// 查询历史快照
	var snapshots []entity.AssetSnapshots
	err := dao.AssetSnapshots.Ctx(ctx).
		Where(dao.AssetSnapshots.Columns().UserId, consts.GetUserID(ctx)).
		WhereGTE(dao.AssetSnapshots.Columns().SnapshotTime, startTime).
		OrderAsc(dao.AssetSnapshots.Columns().SnapshotTime).
		Scan(&snapshots)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产历史失败")
	}

	// 转换为 API 响应格式
	items := make([]assetApi.SnapshotItem, 0, len(snapshots))
	for _, snap := range snapshots {
		items = append(items, assetApi.SnapshotItem{
			Date:       snap.SnapshotTime.Format("2006-01-02"),
			TotalValue: snap.TotalValueUsd,
			Currency:   currency,
		})
	}

	g.Log().Info(ctx, "获取资产历史成功",
		"days", days,
		"currency", currency,
		"count", len(items),
	)

	return &assetApi.GetHistoryRes{Snapshots: items}, nil
}

// RefreshAll 强制刷新所有资产
//
// 功能说明:
// 1. 清空当前用户资产明细
// 2. 获取所有 CEX 账户余额
// 3. 获取所有钱包余额（待实现）
// 4. 获取手动资产
// 5. 创建新的资产快照
func (s *sAsset) RefreshAll(ctx context.Context) (*assetApi.RefreshRes, error) {
	g.Log().Info(ctx, "开始刷新资产")

	userID := consts.GetUserID(ctx)
	var totalAssetCount int

	// 1. 清空当前用户的资产明细
	_, err := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, userID).
		Delete()
	if err != nil {
		return nil, gerror.Wrap(err, "清空旧资产数据失败")
	}

	// 2. 获取所有 CEX 账户
	var exchangeAccounts []entity.ExchangeAccounts
	err = exchangeDao.ExchangeAccounts.Ctx(ctx).
		Where(exchangeDao.ExchangeAccounts.Columns().UserId, userID).
		Where(exchangeDao.ExchangeAccounts.Columns().IsActive, 1).
		Scan(&exchangeAccounts)
	if err != nil {
		return nil, gerror.Wrap(err, "查询交易所账户失败")
	}

	g.Log().Info(ctx, "查询到交易所账户", "count", len(exchangeAccounts))

	// 3. 调用各 CEX API 获取余额并保存
	for _, account := range exchangeAccounts {
		balances, _, err := exchangeService.Exchange().GetBalances(ctx, account.Id)
		if err != nil {
			g.Log().Warning(ctx, "获取交易所资产失败",
				"accountId", account.Id,
				"exchange", account.ExchangeName,
				"error", err,
			)
			continue
		}

		// 保存资产明细
		for _, bal := range balances {
			// 计算单价：如果总余额大于 0，单价 = 价值 / 余额
			priceUsd := 0.0
			if bal.Total > 0 {
				priceUsd = bal.ValueUSD / bal.Total
			}
			detail := entity.AssetDetails{
				UserId:      userID,
				SourceType:  "cex",
				SourceId:    account.Id,
				AssetSymbol: bal.Symbol,
				AssetName:   account.ExchangeName,
				Balance:     bal.Total,
				PriceUsd:    priceUsd,
				ValueUsd:    bal.ValueUSD,
				LastUpdated: gtime.Now().Time,
			}
			_, insertErr := dao.AssetDetails.Ctx(ctx).Insert(detail)
			if insertErr != nil {
				g.Log().Error(ctx, "保存交易所资产失败", "error", insertErr)
				continue
			}
			totalAssetCount++
		}
	}

	// 4. 获取钱包地址资产
	var walletAddresses []entity.WalletAddresses
	err = walletDao.WalletAddresses.Ctx(ctx).
		Where(walletDao.WalletAddresses.Columns().UserId, userID).
		Where(walletDao.WalletAddresses.Columns().IsActive, 1).
		Scan(&walletAddresses)
	if err == nil && len(walletAddresses) > 0 {
		g.Log().Info(ctx, "查询到钱包地址", "count", len(walletAddresses))

		// 逐一调用钱包服务同步区块链余额到 asset_details 表
		for _, addr := range walletAddresses {
			syncErr := walletService.Wallet().SyncAddress(ctx, addr.Id)
			if syncErr != nil {
				g.Log().Warning(ctx, "同步钱包地址余额失败",
					"walletId", addr.Id,
					"blockchain", addr.Blockchain,
					"address", addr.Address,
					"error", syncErr,
				)
				continue
			}

			// 统计该钱包同步后写入的资产条数
			walletLabel := addr.Blockchain + ":" + addr.Address
			count, countErr := dao.AssetDetails.Ctx(ctx).
				Where(dao.AssetDetails.Columns().UserId, userID).
				Where(dao.AssetDetails.Columns().SourceType, "wallet").
				Where(dao.AssetDetails.Columns().AssetName, walletLabel).
				Count()
			if countErr == nil {
				totalAssetCount += count
			}
		}
	}

	// 5. 获取手动资产
	var manualAssets []entity.ManualAssets
	err = manualAssetDao.ManualAssets.Ctx(ctx).
		Where(manualAssetDao.ManualAssets.Columns().UserId, userID).
		Where(manualAssetDao.ManualAssets.Columns().IsActive, 1).
		Scan(&manualAssets)
	if err == nil {
		for _, ma := range manualAssets {
			detail := entity.AssetDetails{
				UserId:      userID,
				SourceType:  "manual",
				SourceId:    ma.Id,
				AssetSymbol: ma.Currency,
				AssetName:   ma.AssetName,
				Balance:     ma.Amount,
				PriceUsd:    1.0,
				ValueUsd:    ma.AmountUsd,
				LastUpdated: gtime.Now().Time,
			}
			_, insertErr := dao.AssetDetails.Ctx(ctx).Insert(detail)
			if insertErr != nil {
				g.Log().Error(ctx, "保存手动资产失败", "error", insertErr)
				continue
			}
			totalAssetCount++
		}
	}

	// 6. 创建资产快照
	err = s.createAssetSnapshot(ctx, userID)
	if err != nil {
		g.Log().Warning(ctx, "创建资产快照失败", "error", err)
	}

	g.Log().Info(ctx, "刷新资产完成", "totalAssetCount", totalAssetCount)

	return &assetApi.RefreshRes{
		Message:        "资产刷新完成",
		RefreshedCount: totalAssetCount,
	}, nil
}

// createAssetSnapshot 创建资产快照
func (s *sAsset) createAssetSnapshot(ctx context.Context, userID int) error {
	// 计算当前总资产价值
	totalValueUSD, err := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, userID).
		Sum(dao.AssetDetails.Columns().ValueUsd)
	if err != nil {
		return gerror.Wrap(err, "计算总资产失败")
	}

	// 计算各来源资产价值
	cexValue, _ := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, userID).
		Where(dao.AssetDetails.Columns().SourceType, "cex").
		Sum(dao.AssetDetails.Columns().ValueUsd)

	blockchainValue, _ := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, userID).
		Where(dao.AssetDetails.Columns().SourceType, "blockchain").
		Sum(dao.AssetDetails.Columns().ValueUsd)

	manualValue, _ := dao.AssetDetails.Ctx(ctx).
		Where(dao.AssetDetails.Columns().UserId, userID).
		Where(dao.AssetDetails.Columns().SourceType, "manual").
		Sum(dao.AssetDetails.Columns().ValueUsd)

	// 创建快照记录
	snapshot := entity.AssetSnapshots{
		UserId:             userID,
		TotalValueUsd:      totalValueUSD,
		CexValueUsd:        cexValue,
		BlockchainValueUsd: blockchainValue,
		ManualValueUsd:     manualValue,
		SnapshotTime:       gtime.Now().Time,
	}

	_, err = dao.AssetSnapshots.Ctx(ctx).OmitEmpty().Insert(snapshot)
	if err != nil {
		return gerror.Wrap(err, "保存资产快照失败")
	}

	g.Log().Info(ctx, "创建资产快照成功",
		"userId", userID,
		"totalValueUSD", totalValueUSD,
	)

	return nil
}

// getConversionRate 获取 USD 到目标货币的转换汇率
// 支持加密货币（BTC/ETH）和法币（CNY/EUR/JPY 等）
// 如果找不到汇率数据，返回 1.0（降级为 USD）
func (s *sAsset) getConversionRate(ctx context.Context, currency string) float64 {
	if currency == "" || currency == "USD" || currency == "USDC" || currency == "USDT" {
		return 1.0
	}

	// 1. 尝试从 exchange_rates 表查找 USD -> target 的汇率（法币场景：USD->CNY）
	var rate entity.ExchangeRates
	err := exchangeRateDao.ExchangeRates.Ctx(ctx).
		Where(exchangeRateDao.ExchangeRates.Columns().FromCurrency, "USD").
		Where(exchangeRateDao.ExchangeRates.Columns().ToCurrency, currency).
		OrderDesc(exchangeRateDao.ExchangeRates.Columns().FetchedAt).
		Scan(&rate)
	if err == nil && rate.Rate > 0 {
		return rate.Rate
	}

	// 2. 尝试反向查找（target -> USD），汇率取倒数（加密货币场景：BTC->USD = 65000，则 USD->BTC = 1/65000）
	err = exchangeRateDao.ExchangeRates.Ctx(ctx).
		Where(exchangeRateDao.ExchangeRates.Columns().FromCurrency, currency).
		Where(exchangeRateDao.ExchangeRates.Columns().ToCurrency, "USD").
		OrderDesc(exchangeRateDao.ExchangeRates.Columns().FetchedAt).
		Scan(&rate)
	if err == nil && rate.Rate > 0 {
		return 1.0 / rate.Rate
	}

	// 3. 对于常见加密货币，直接调用 CoinGecko 获取实时汇率
	cgClient := coingecko.NewClient("")
	usdPrice, err := cgClient.GetPrice(ctx, currency)
	if err == nil && usdPrice > 0 {
		return 1.0 / usdPrice // USD -> crypto = 1/price
	}

	g.Log().Warning(ctx, "未找到货币转换汇率，降级使用 USD", "currency", currency)
	return 1.0
}
