// Package logic 钱包业务逻辑
// 实现钱包地址的增删改查、余额查询和批量导入
package logic

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	walletApi "your-finance/allfi/api/v1/wallet"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/wallet/dao"
	"your-finance/allfi/internal/app/wallet/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/coingecko"
	"your-finance/allfi/internal/integrations/etherscan"
	"your-finance/allfi/internal/model/entity"
)

// sWallet 钱包服务实现
type sWallet struct{}

// New 创建钱包服务实例
func New() service.IWallet {
	return &sWallet{}
}

// ListWallets 获取钱包地址列表
func (s *sWallet) ListWallets(ctx context.Context, userID int) ([]walletApi.AddressItem, error) {
	var wallets []entity.WalletAddresses
	err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().UserId, userID).
		Where(dao.WalletAddresses.Columns().IsActive, 1).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		OrderDesc(dao.WalletAddresses.Columns().CreatedAt).
		Scan(&wallets)
	if err != nil {
		return nil, gerror.Wrap(err, "查询钱包地址列表失败")
	}

	items := make([]walletApi.AddressItem, 0, len(wallets))
	for _, w := range wallets {
		items = append(items, s.toAddressItem(&w))
	}
	return items, nil
}

// CreateWallet 添加钱包地址
func (s *sWallet) CreateWallet(ctx context.Context, userID int, req *walletApi.CreateAddressReq) (*walletApi.AddressItem, error) {
	// 验证地址格式（必须以 0x 开头）
	if !strings.HasPrefix(req.Address, "0x") {
		return nil, gerror.New("钱包地址必须以 0x 开头")
	}

	// 地址转小写统一格式
	address := strings.ToLower(req.Address)

	// 检查地址是否已存在
	count, err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().UserId, userID).
		Where(dao.WalletAddresses.Columns().Blockchain, req.Blockchain).
		Where(dao.WalletAddresses.Columns().Address, address).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "检查钱包地址失败")
	}
	if count > 0 {
		return nil, gerror.New("该钱包地址已存在")
	}

	// 插入数据库
	now := gtime.Now()
	result, err := dao.WalletAddresses.Ctx(ctx).Data(g.Map{
		dao.WalletAddresses.Columns().UserId:     userID,
		dao.WalletAddresses.Columns().Blockchain: req.Blockchain,
		dao.WalletAddresses.Columns().Address:    address,
		dao.WalletAddresses.Columns().Label:      req.Label,
		dao.WalletAddresses.Columns().IsActive:   1,
		dao.WalletAddresses.Columns().CreatedAt:  now,
		dao.WalletAddresses.Columns().UpdatedAt:  now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存钱包地址失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取钱包ID失败")
	}

	g.Log().Info(ctx, "添加钱包地址成功",
		"walletId", id,
		"blockchain", req.Blockchain,
		"address", address,
	)

	return s.GetWallet(ctx, int(id))
}

// GetWallet 获取单个钱包地址
func (s *sWallet) GetWallet(ctx context.Context, walletID int) (*walletApi.AddressItem, error) {
	var wallet entity.WalletAddresses
	err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, walletID).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		Scan(&wallet)
	if err != nil {
		return nil, gerror.Wrap(err, "查询钱包信息失败")
	}
	if wallet.Id == 0 {
		return nil, gerror.Newf("钱包不存在: %d", walletID)
	}

	item := s.toAddressItem(&wallet)
	return &item, nil
}

// UpdateWallet 更新钱包地址信息
func (s *sWallet) UpdateWallet(ctx context.Context, req *walletApi.UpdateAddressReq) (*walletApi.AddressItem, error) {
	// 检查是否存在
	var existing entity.WalletAddresses
	err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, req.Id).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询钱包信息失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("钱包不存在: %d", req.Id)
	}

	// 构建更新数据
	updateData := g.Map{
		dao.WalletAddresses.Columns().Label:     req.Label,
		dao.WalletAddresses.Columns().UpdatedAt: gtime.Now(),
	}
	if req.Blockchain != "" {
		updateData[dao.WalletAddresses.Columns().Blockchain] = req.Blockchain
	}
	if req.Address != "" {
		updateData[dao.WalletAddresses.Columns().Address] = strings.ToLower(req.Address)
	}

	_, err = dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新钱包地址失败")
	}

	g.Log().Info(ctx, "更新钱包地址成功", "walletId", req.Id)
	return s.GetWallet(ctx, int(req.Id))
}

// DeleteWallet 软删除钱包地址
func (s *sWallet) DeleteWallet(ctx context.Context, walletID int) error {
	_, err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, walletID).
		Data(g.Map{
			dao.WalletAddresses.Columns().DeletedAt: gtime.Now(),
			dao.WalletAddresses.Columns().IsActive:  0,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除钱包地址失败")
	}

	g.Log().Info(ctx, "删除钱包地址成功", "walletId", walletID)
	return nil
}

// GetBalances 获取钱包余额
// 返回: 原生代币余额, Token 余额映射, 总价值（USD）, 错误
func (s *sWallet) GetBalances(ctx context.Context, walletID int) (float64, map[string]float64, float64, error) {
	// 查询钱包信息
	var wallet entity.WalletAddresses
	err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, walletID).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		Scan(&wallet)
	if err != nil {
		return 0, nil, 0, gerror.Wrap(err, "查询钱包信息失败")
	}
	if wallet.Id == 0 {
		return 0, nil, 0, gerror.Newf("钱包不存在: %d", walletID)
	}

	// 根据区块链选择对应客户端
	chainConfig, ok := etherscan.SupportedChains[wallet.Blockchain]
	if !ok {
		return 0, nil, 0, gerror.Newf("暂不支持的区块链: %s", wallet.Blockchain)
	}

	configKey := "external." + chainConfig.RateLimitKey + "ApiKey"
	apiKey := g.Cfg().MustGet(ctx, configKey).String()
	if apiKey == "" {
		return 0, nil, 0, gerror.Newf("%s API Key 未配置（配置项: %s）", wallet.Blockchain, configKey)
	}

	client, err := etherscan.NewChainClient(wallet.Blockchain, apiKey)
	if err != nil {
		return 0, nil, 0, gerror.Wrap(err, "创建区块链客户端失败")
	}

	// 查询原生代币余额
	nativeBalance, err := client.GetNativeBalance(ctx, wallet.Address)
	if err != nil {
		return 0, nil, 0, gerror.Wrap(err, "获取原生代币余额失败")
	}

	// 查询 Token 余额
	tokenBalances := make(map[string]float64)
	tokenBalances[chainConfig.NativeSymbol] = nativeBalance

	tokens, err := client.GetTokenBalances(ctx, wallet.Address)
	if err != nil {
		g.Log().Warning(ctx, "获取代币余额失败", "walletId", walletID, "error", err)
	} else {
		for _, token := range tokens {
			// Token 余额已由 etherscan 客户端根据 decimal 转换，直接使用
			tokenBalances[token.Symbol] = token.Total
		}
	}

	// 使用 CoinGecko 批量获取价格，计算 USD 总价值
	var totalUSD float64
	symbols := make([]string, 0, len(tokenBalances))
	for sym := range tokenBalances {
		symbols = append(symbols, sym)
	}

	if len(symbols) > 0 {
		cgClient := coingecko.NewClient("")
		prices, priceErr := cgClient.GetPrices(ctx, symbols)
		if priceErr != nil {
			g.Log().Warning(ctx, "获取价格失败，USD 价值可能不准确", "error", priceErr)
		} else {
			for sym, balance := range tokenBalances {
				if price, ok := prices[strings.ToUpper(sym)]; ok {
					totalUSD += balance * price
				}
			}
		}
	}

	g.Log().Info(ctx, "获取钱包余额成功",
		"walletId", walletID,
		"blockchain", wallet.Blockchain,
		"nativeBalance", nativeBalance,
		"totalUSD", totalUSD,
	)

	return nativeBalance, tokenBalances, totalUSD, nil
}

// SyncAddress 同步钱包地址余额到 asset_details 表
// 1. 调用 GetBalances 获取最新余额和价格
// 2. 删除该钱包在 asset_details 中的旧记录
// 3. 遍历余额列表，插入新的 asset_details 记录
func (s *sWallet) SyncAddress(ctx context.Context, walletID int) error {
	// 查询钱包信息
	var wallet entity.WalletAddresses
	err := dao.WalletAddresses.Ctx(ctx).
		Where(dao.WalletAddresses.Columns().Id, walletID).
		WhereNull(dao.WalletAddresses.Columns().DeletedAt).
		Scan(&wallet)
	if err != nil {
		return gerror.Wrap(err, "查询钱包信息失败")
	}
	if wallet.Id == 0 {
		return gerror.Newf("钱包不存在: %d", walletID)
	}

	// 获取最新余额
	_, tokenBalances, _, err := s.GetBalances(ctx, walletID)
	if err != nil {
		return gerror.Wrap(err, "获取钱包余额失败")
	}

	// 钱包地址作为资产名称标识
	walletLabel := wallet.Blockchain + ":" + wallet.Address

	// 删除该钱包在 asset_details 中的旧记录
	_, err = assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().SourceType, "wallet").
		Where(assetDao.AssetDetails.Columns().AssetName, walletLabel).
		Where(assetDao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "删除旧资产明细失败")
	}

	// 使用 CoinGecko 获取价格
	symbols := make([]string, 0, len(tokenBalances))
	for sym := range tokenBalances {
		symbols = append(symbols, sym)
	}

	var prices map[string]float64
	if len(symbols) > 0 {
		cgClient := coingecko.NewClient("")
		prices, err = cgClient.GetPrices(ctx, symbols)
		if err != nil {
			g.Log().Warning(ctx, "获取价格失败，价格信息可能不完整", "error", err)
			prices = make(map[string]float64)
		}
	} else {
		prices = make(map[string]float64)
	}

	// 遍历余额列表，插入新的 asset_details 记录
	now := time.Now()
	insertCount := 0
	for sym, balance := range tokenBalances {
		if balance <= 0 {
			continue
		}

		priceUSD := prices[strings.ToUpper(sym)]
		valueUSD := balance * priceUSD

		_, insertErr := assetDao.AssetDetails.Ctx(ctx).Data(g.Map{
			assetDao.AssetDetails.Columns().UserId:      consts.GetUserID(ctx),
			assetDao.AssetDetails.Columns().AssetSymbol: sym,
			assetDao.AssetDetails.Columns().AssetName:   walletLabel,
			assetDao.AssetDetails.Columns().Balance:     balance,
			assetDao.AssetDetails.Columns().PriceUsd:    priceUSD,
			assetDao.AssetDetails.Columns().ValueUsd:    valueUSD,
			assetDao.AssetDetails.Columns().SourceType:  "wallet",
			assetDao.AssetDetails.Columns().LastUpdated: now,
			assetDao.AssetDetails.Columns().CreatedAt:   now,
			assetDao.AssetDetails.Columns().UpdatedAt:   now,
		}).Insert()
		if insertErr != nil {
			g.Log().Error(ctx, "插入资产明细失败",
				"symbol", sym,
				"wallet", walletLabel,
				"error", insertErr,
			)
			continue
		}
		insertCount++
	}

	g.Log().Info(ctx, "同步钱包地址余额完成",
		"walletId", walletID,
		"blockchain", wallet.Blockchain,
		"address", wallet.Address,
		"insertCount", insertCount,
	)

	return nil
}

// BatchImport 批量导入钱包地址
// 返回: 成功数, 失败数, 错误
func (s *sWallet) BatchImport(ctx context.Context, userID int, req *walletApi.BatchImportReq) (int, int, error) {
	imported := 0
	failed := 0

	for _, addr := range req.Addresses {
		// 验证地址格式
		if !strings.HasPrefix(addr.Address, "0x") {
			failed++
			continue
		}

		address := strings.ToLower(addr.Address)

		// 检查是否已存在
		count, err := dao.WalletAddresses.Ctx(ctx).
			Where(dao.WalletAddresses.Columns().UserId, userID).
			Where(dao.WalletAddresses.Columns().Blockchain, addr.Blockchain).
			Where(dao.WalletAddresses.Columns().Address, address).
			WhereNull(dao.WalletAddresses.Columns().DeletedAt).
			Count()
		if err != nil || count > 0 {
			failed++
			continue
		}

		// 插入数据库
		now := gtime.Now()
		_, err = dao.WalletAddresses.Ctx(ctx).Data(g.Map{
			dao.WalletAddresses.Columns().UserId:     userID,
			dao.WalletAddresses.Columns().Blockchain: addr.Blockchain,
			dao.WalletAddresses.Columns().Address:    address,
			dao.WalletAddresses.Columns().Label:      addr.Label,
			dao.WalletAddresses.Columns().IsActive:   1,
			dao.WalletAddresses.Columns().CreatedAt:  now,
			dao.WalletAddresses.Columns().UpdatedAt:  now,
		}).Insert()
		if err != nil {
			failed++
			g.Log().Error(ctx, "批量导入钱包地址失败", "address", address, "error", err)
			continue
		}
		imported++
	}

	g.Log().Info(ctx, "批量导入钱包地址完成",
		"userId", userID,
		"imported", imported,
		"failed", failed,
	)

	return imported, failed, nil
}

// toAddressItem 将数据库实体转换为 API 响应格式
func (s *sWallet) toAddressItem(w *entity.WalletAddresses) walletApi.AddressItem {
	return walletApi.AddressItem{
		ID:         uint(w.Id),
		Blockchain: w.Blockchain,
		Address:    w.Address,
		Label:      w.Label,
		CreatedAt:  w.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:  w.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
