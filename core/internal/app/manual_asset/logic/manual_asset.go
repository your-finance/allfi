// Package logic 手动资产业务逻辑
// 实现手动资产（银行账户/现金/股票/基金）的增删改查
package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	manualAssetApi "your-finance/allfi/api/v1/manual_asset"
	exchangeRateDao "your-finance/allfi/internal/app/exchange_rate/dao"
	"your-finance/allfi/internal/app/manual_asset/dao"
	"your-finance/allfi/internal/app/manual_asset/service"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/model/entity"
)

// sManualAsset 手动资产服务实现
type sManualAsset struct{}

// New 创建手动资产服务实例
func New() service.IManualAsset {
	return &sManualAsset{}
}

// ListManualAssets 获取手动资产列表
func (s *sManualAsset) ListManualAssets(ctx context.Context, userID int) ([]manualAssetApi.ManualAssetItem, error) {
	var assets []entity.ManualAssets
	err := dao.ManualAssets.Ctx(ctx).
		Where(dao.ManualAssets.Columns().UserId, userID).
		Where(dao.ManualAssets.Columns().IsActive, 1).
		WhereNull(dao.ManualAssets.Columns().DeletedAt).
		OrderDesc(dao.ManualAssets.Columns().CreatedAt).
		Scan(&assets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询手动资产列表失败")
	}

	items := make([]manualAssetApi.ManualAssetItem, 0, len(assets))
	for _, a := range assets {
		items = append(items, s.toAssetItem(&a))
	}
	return items, nil
}

// CreateManualAsset 添加手动资产
func (s *sManualAsset) CreateManualAsset(ctx context.Context, userID int, req *manualAssetApi.CreateReq) (*manualAssetApi.ManualAssetItem, error) {
	// 插入数据库
	now := gtime.Now()
	result, err := dao.ManualAssets.Ctx(ctx).Data(g.Map{
		dao.ManualAssets.Columns().UserId:    userID,
		dao.ManualAssets.Columns().AssetType: req.AssetType,
		dao.ManualAssets.Columns().AssetName: req.AssetName,
		dao.ManualAssets.Columns().Amount:    req.Amount,
		dao.ManualAssets.Columns().AmountUsd: s.convertToUSD(ctx, req.Amount, req.Currency),
		dao.ManualAssets.Columns().Currency:  req.Currency,
		dao.ManualAssets.Columns().Notes:     req.Notes,
		dao.ManualAssets.Columns().IsActive:  1,
		dao.ManualAssets.Columns().CreatedAt: now,
		dao.ManualAssets.Columns().UpdatedAt: now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存手动资产失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取资产ID失败")
	}

	g.Log().Info(ctx, "添加手动资产成功",
		"assetId", id,
		"assetType", req.AssetType,
		"assetName", req.AssetName,
	)

	return s.getAssetByID(ctx, int(id))
}

// UpdateManualAsset 更新手动资产
func (s *sManualAsset) UpdateManualAsset(ctx context.Context, req *manualAssetApi.UpdateReq) (*manualAssetApi.ManualAssetItem, error) {
	// 检查是否存在
	var existing entity.ManualAssets
	err := dao.ManualAssets.Ctx(ctx).
		Where(dao.ManualAssets.Columns().Id, req.Id).
		WhereNull(dao.ManualAssets.Columns().DeletedAt).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产信息失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("资产不存在: %d", req.Id)
	}

	// 构建更新数据
	updateData := g.Map{
		dao.ManualAssets.Columns().UpdatedAt: gtime.Now(),
	}
	if req.AssetType != "" {
		updateData[dao.ManualAssets.Columns().AssetType] = req.AssetType
	}
	if req.AssetName != "" {
		updateData[dao.ManualAssets.Columns().AssetName] = req.AssetName
	}
	if req.Amount > 0 {
		updateData[dao.ManualAssets.Columns().Amount] = req.Amount
		// 使用当前货币或更新后的货币计算 USD 值
		currency := req.Currency
		if currency == "" {
			currency = existing.Currency
		}
		updateData[dao.ManualAssets.Columns().AmountUsd] = s.convertToUSD(ctx, req.Amount, currency)
	}
	if req.Currency != "" {
		updateData[dao.ManualAssets.Columns().Currency] = req.Currency
	}
	if req.Notes != "" {
		updateData[dao.ManualAssets.Columns().Notes] = req.Notes
	}

	_, err = dao.ManualAssets.Ctx(ctx).
		Where(dao.ManualAssets.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新手动资产失败")
	}

	g.Log().Info(ctx, "更新手动资产成功", "assetId", req.Id)

	return s.getAssetByID(ctx, int(req.Id))
}

// DeleteManualAsset 删除手动资产（软删除）
func (s *sManualAsset) DeleteManualAsset(ctx context.Context, assetID int) error {
	_, err := dao.ManualAssets.Ctx(ctx).
		Where(dao.ManualAssets.Columns().Id, assetID).
		Data(g.Map{
			dao.ManualAssets.Columns().DeletedAt: gtime.Now(),
			dao.ManualAssets.Columns().IsActive:  0,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除手动资产失败")
	}

	g.Log().Info(ctx, "删除手动资产成功", "assetId", assetID)
	return nil
}

// getAssetByID 根据 ID 查询手动资产
func (s *sManualAsset) getAssetByID(ctx context.Context, assetID int) (*manualAssetApi.ManualAssetItem, error) {
	var asset entity.ManualAssets
	err := dao.ManualAssets.Ctx(ctx).
		Where(dao.ManualAssets.Columns().Id, assetID).
		WhereNull(dao.ManualAssets.Columns().DeletedAt).
		Scan(&asset)
	if err != nil {
		return nil, gerror.Wrap(err, "查询资产信息失败")
	}
	if asset.Id == 0 {
		return nil, gerror.Newf("资产不存在: %d", assetID)
	}

	item := s.toAssetItem(&asset)
	return &item, nil
}

// toAssetItem 将数据库实体转换为 API 响应格式
func (s *sManualAsset) toAssetItem(a *entity.ManualAssets) manualAssetApi.ManualAssetItem {
	return manualAssetApi.ManualAssetItem{
		ID:        uint(a.Id),
		AssetType: a.AssetType,
		AssetName: a.AssetName,
		Amount:    a.Amount,
		Currency:  a.Currency,
		Notes:     a.Notes,
		CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// convertToUSD 将指定货币金额转换为 USD
// 对于 USD/USDC/USDT 直接返回原值
// 对于其他货币，查询 exchange_rates 表或调用 CoinGecko
func (s *sManualAsset) convertToUSD(ctx context.Context, amount float64, currency string) float64 {
	if currency == "" || currency == "USD" || currency == "USDC" || currency == "USDT" {
		return amount
	}

	// 1. 尝试从 exchange_rates 表查找 currency -> USD 的汇率
	var rate entity.ExchangeRates
	err := exchangeRateDao.ExchangeRates.Ctx(ctx).
		Where(exchangeRateDao.ExchangeRates.Columns().FromCurrency, currency).
		Where(exchangeRateDao.ExchangeRates.Columns().ToCurrency, "USD").
		OrderDesc(exchangeRateDao.ExchangeRates.Columns().FetchedAt).
		Scan(&rate)
	if err == nil && rate.Rate > 0 {
		return amount * rate.Rate
	}

	// 2. 尝试反向查找（USD -> currency），汇率取倒数
	err = exchangeRateDao.ExchangeRates.Ctx(ctx).
		Where(exchangeRateDao.ExchangeRates.Columns().FromCurrency, "USD").
		Where(exchangeRateDao.ExchangeRates.Columns().ToCurrency, currency).
		OrderDesc(exchangeRateDao.ExchangeRates.Columns().FetchedAt).
		Scan(&rate)
	if err == nil && rate.Rate > 0 {
		return amount / rate.Rate
	}

	// 3. 对于加密货币，调用 CoinGecko 获取 USD 价格
	cgClient := coingecko.NewClient("")
	usdPrice, err := cgClient.GetPrice(ctx, currency)
	if err == nil && usdPrice > 0 {
		return amount * usdPrice
	}

	g.Log().Warning(ctx, "无法获取汇率，AmountUSD 使用原值", "currency", currency)
	return amount
}
