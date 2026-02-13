// Package logic 交易所业务逻辑
// 实现交易所账户的增删改查、连接测试和余额查询
package logic

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	exchangeApi "your-finance/allfi/api/v1/exchange"
	assetDao "your-finance/allfi/internal/app/asset/dao"
	"your-finance/allfi/internal/app/exchange/dao"
	"your-finance/allfi/internal/app/exchange/service"
	"your-finance/allfi/internal/consts"
	"your-finance/allfi/internal/integrations/binance"
	"your-finance/allfi/internal/integrations/coinbase"
	"your-finance/allfi/internal/integrations/okx"
	"your-finance/allfi/internal/model/entity"
	"your-finance/allfi/utility/crypto"
)

// sExchange 交易所服务实现
type sExchange struct{}

// New 创建交易所服务实例
func New() service.IExchange {
	return &sExchange{}
}

// ListAccounts 获取用户的交易所账户列表
// 过滤加密字段，只返回安全信息
func (s *sExchange) ListAccounts(ctx context.Context, userID int) ([]exchangeApi.AccountItem, error) {
	var accounts []entity.ExchangeAccounts
	err := dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().UserId, userID).
		Where(dao.ExchangeAccounts.Columns().IsActive, 1).
		WhereNull(dao.ExchangeAccounts.Columns().DeletedAt).
		OrderDesc(dao.ExchangeAccounts.Columns().CreatedAt).
		Scan(&accounts)
	if err != nil {
		return nil, gerror.Wrap(err, "查询交易所账户列表失败")
	}

	// 转换为 API 响应格式（过滤敏感字段）
	items := make([]exchangeApi.AccountItem, 0, len(accounts))
	for _, acc := range accounts {
		items = append(items, s.toAccountItem(&acc))
	}
	return items, nil
}

// CreateAccount 添加交易所账户
// 加密 API 凭证后存储到数据库
func (s *sExchange) CreateAccount(ctx context.Context, userID int, req *exchangeApi.CreateAccountReq) (*exchangeApi.AccountItem, error) {
	// 获取加密密钥
	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()
	if len(masterKey) != 32 {
		return nil, gerror.New("主加密密钥配置错误，长度必须为32字节")
	}

	// 加密 API Key
	encryptedKey, err := crypto.EncryptAES(req.ApiKey, masterKey)
	if err != nil {
		return nil, gerror.Wrap(err, "加密 API Key 失败")
	}

	// 加密 API Secret
	encryptedSecret, err := crypto.EncryptAES(req.ApiSecret, masterKey)
	if err != nil {
		return nil, gerror.Wrap(err, "加密 API Secret 失败")
	}

	// 加密 Passphrase（如有）
	var encryptedPassphrase string
	if req.Passphrase != "" {
		encryptedPassphrase, err = crypto.EncryptAES(req.Passphrase, masterKey)
		if err != nil {
			return nil, gerror.Wrap(err, "加密 API Passphrase 失败")
		}
	}

	// 插入数据库
	now := gtime.Now()
	result, err := dao.ExchangeAccounts.Ctx(ctx).Data(g.Map{
		dao.ExchangeAccounts.Columns().UserId:                 userID,
		dao.ExchangeAccounts.Columns().ExchangeName:           req.ExchangeName,
		dao.ExchangeAccounts.Columns().ApiKeyEncrypted:        encryptedKey,
		dao.ExchangeAccounts.Columns().ApiSecretEncrypted:     encryptedSecret,
		dao.ExchangeAccounts.Columns().ApiPassphraseEncrypted: encryptedPassphrase,
		dao.ExchangeAccounts.Columns().Label:                  req.Label,
		dao.ExchangeAccounts.Columns().Note:                   req.Note,
		dao.ExchangeAccounts.Columns().IsActive:               1,
		dao.ExchangeAccounts.Columns().CreatedAt:              now,
		dao.ExchangeAccounts.Columns().UpdatedAt:              now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存交易所账户失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取账户ID失败")
	}

	g.Log().Info(ctx, "添加交易所账户成功",
		"accountId", id,
		"exchange", req.ExchangeName,
		"userId", userID,
	)

	// 查询并返回新建的账户
	return s.GetAccount(ctx, int(id))
}

// GetAccount 获取单个交易所账户详情
func (s *sExchange) GetAccount(ctx context.Context, accountID int) (*exchangeApi.AccountItem, error) {
	var account entity.ExchangeAccounts
	err := dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, accountID).
		WhereNull(dao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&account)
	if err != nil {
		return nil, gerror.Wrap(err, "查询账户信息失败")
	}
	if account.Id == 0 {
		return nil, gerror.Newf("账户不存在: %d", accountID)
	}

	item := s.toAccountItem(&account)
	return &item, nil
}

// UpdateAccount 更新交易所账户信息
func (s *sExchange) UpdateAccount(ctx context.Context, req *exchangeApi.UpdateAccountReq) (*exchangeApi.AccountItem, error) {
	// 检查账户是否存在
	var existing entity.ExchangeAccounts
	err := dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, req.Id).
		WhereNull(dao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&existing)
	if err != nil {
		return nil, gerror.Wrap(err, "查询账户信息失败")
	}
	if existing.Id == 0 {
		return nil, gerror.Newf("账户不存在: %d", req.Id)
	}

	// 构建更新数据
	updateData := g.Map{
		dao.ExchangeAccounts.Columns().Label:     req.Label,
		dao.ExchangeAccounts.Columns().Note:      req.Note,
		dao.ExchangeAccounts.Columns().UpdatedAt: gtime.Now(),
	}

	// 如果提供了新的 API 凭证，则加密更新
	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()
	if req.ApiKey != "" {
		encrypted, err := crypto.EncryptAES(req.ApiKey, masterKey)
		if err != nil {
			return nil, gerror.Wrap(err, "加密 API Key 失败")
		}
		updateData[dao.ExchangeAccounts.Columns().ApiKeyEncrypted] = encrypted
	}
	if req.ApiSecret != "" {
		encrypted, err := crypto.EncryptAES(req.ApiSecret, masterKey)
		if err != nil {
			return nil, gerror.Wrap(err, "加密 API Secret 失败")
		}
		updateData[dao.ExchangeAccounts.Columns().ApiSecretEncrypted] = encrypted
	}
	if req.Passphrase != "" {
		encrypted, err := crypto.EncryptAES(req.Passphrase, masterKey)
		if err != nil {
			return nil, gerror.Wrap(err, "加密 API Passphrase 失败")
		}
		updateData[dao.ExchangeAccounts.Columns().ApiPassphraseEncrypted] = encrypted
	}

	// 执行更新
	_, err = dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, req.Id).
		Data(updateData).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新交易所账户失败")
	}

	g.Log().Info(ctx, "更新交易所账户成功", "accountId", req.Id)

	return s.GetAccount(ctx, int(req.Id))
}

// DeleteAccount 软删除交易所账户
func (s *sExchange) DeleteAccount(ctx context.Context, accountID int) error {
	_, err := dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, accountID).
		Data(g.Map{
			dao.ExchangeAccounts.Columns().DeletedAt: gtime.Now(),
			dao.ExchangeAccounts.Columns().IsActive:  0,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "删除交易所账户失败")
	}

	g.Log().Info(ctx, "删除交易所账户成功", "accountId", accountID)
	return nil
}

// TestConnection 测试交易所 API 连接
func (s *sExchange) TestConnection(ctx context.Context, accountID int) (bool, string, error) {
	// 获取账户并解密凭证
	apiKey, apiSecret, passphrase, exchangeName, err := s.decryptCredentials(ctx, accountID)
	if err != nil {
		return false, "", err
	}

	// 根据交易所类型创建客户端并测试连接
	switch exchangeName {
	case "binance":
		client := binance.NewClient(apiKey, apiSecret)
		err = client.TestConnection(ctx)
	case "okx":
		client := okx.NewClient(apiKey, apiSecret, passphrase)
		err = client.TestConnection(ctx)
	case "coinbase":
		client := coinbase.NewClient(apiKey, apiSecret)
		err = client.TestConnection(ctx)
	default:
		return false, "", gerror.Newf("暂不支持的交易所: %s", exchangeName)
	}

	if err != nil {
		g.Log().Error(ctx, "交易所连接测试失败",
			"accountId", accountID,
			"exchange", exchangeName,
			"error", err,
		)
		return false, "API 连接测试失败: " + err.Error(), nil
	}

	g.Log().Info(ctx, "交易所连接测试成功", "accountId", accountID, "exchange", exchangeName)
	return true, "连接成功", nil
}

// GetBalances 获取交易所账户余额
func (s *sExchange) GetBalances(ctx context.Context, accountID int) ([]exchangeApi.BalanceItem, float64, error) {
	// 获取账户并解密凭证
	apiKey, apiSecret, passphrase, exchangeName, err := s.decryptCredentials(ctx, accountID)
	if err != nil {
		return nil, 0, err
	}

	// 根据交易所类型获取余额
	var balanceItems []exchangeApi.BalanceItem
	var totalValue float64

	switch exchangeName {
	case "binance":
		client := binance.NewClient(apiKey, apiSecret)
		balances, bErr := client.GetAllBalances(ctx)
		if bErr != nil {
			return nil, 0, gerror.Wrap(bErr, "获取 Binance 余额失败")
		}
		for _, b := range balances {
			if b.Total == 0 {
				continue
			}
			item := exchangeApi.BalanceItem{
				Symbol:   b.Symbol,
				Free:     b.Free,
				Locked:   b.Locked,
				Total:    b.Total,
				ValueUSD: b.ValueUSD,
			}
			totalValue += b.ValueUSD
			balanceItems = append(balanceItems, item)
		}
	case "okx":
		client := okx.NewClient(apiKey, apiSecret, passphrase)
		balances, bErr := client.GetBalances(ctx)
		if bErr != nil {
			return nil, 0, gerror.Wrap(bErr, "获取 OKX 余额失败")
		}
		for _, b := range balances {
			if b.Total == 0 {
				continue
			}
			item := exchangeApi.BalanceItem{
				Symbol:   b.Symbol,
				Free:     b.Free,
				Locked:   b.Locked,
				Total:    b.Total,
				ValueUSD: b.ValueUSD,
			}
			totalValue += b.ValueUSD
			balanceItems = append(balanceItems, item)
		}
	case "coinbase":
		client := coinbase.NewClient(apiKey, apiSecret)
		balances, bErr := client.GetBalances(ctx)
		if bErr != nil {
			return nil, 0, gerror.Wrap(bErr, "获取 Coinbase 余额失败")
		}
		for _, b := range balances {
			if b.Total == 0 {
				continue
			}
			item := exchangeApi.BalanceItem{
				Symbol:   b.Symbol,
				Free:     b.Free,
				Locked:   b.Locked,
				Total:    b.Total,
				ValueUSD: b.ValueUSD,
			}
			totalValue += b.ValueUSD
			balanceItems = append(balanceItems, item)
		}
	default:
		return nil, 0, gerror.Newf("暂不支持的交易所: %s", exchangeName)
	}

	if balanceItems == nil {
		balanceItems = []exchangeApi.BalanceItem{}
	}

	return balanceItems, totalValue, nil
}

// SyncAccount 同步交易所账户余额到 asset_details 表
// 1. 调用 GetBalances 获取最新余额
// 2. 查询账户信息获取交易所名称
// 3. 删除该账户在 asset_details 中的旧记录
// 4. 遍历余额列表，插入新的 asset_details 记录
func (s *sExchange) SyncAccount(ctx context.Context, accountID int) error {
	// 获取最新余额
	balances, _, err := s.GetBalances(ctx, accountID)
	if err != nil {
		return gerror.Wrap(err, "获取交易所余额失败")
	}

	// 查询账户信息获取交易所名称
	var account entity.ExchangeAccounts
	err = dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, accountID).
		WhereNull(dao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&account)
	if err != nil {
		return gerror.Wrap(err, "查询交易所账户信息失败")
	}
	if account.Id == 0 {
		return gerror.Newf("账户不存在: %d", accountID)
	}

	exchangeName := account.ExchangeName

	// 删除该账户在 asset_details 中的旧记录
	_, err = assetDao.AssetDetails.Ctx(ctx).
		Where(assetDao.AssetDetails.Columns().SourceType, "cex").
		Where(assetDao.AssetDetails.Columns().AssetName, exchangeName).
		Where(assetDao.AssetDetails.Columns().UserId, consts.GetUserID(ctx)).
		Delete()
	if err != nil {
		return gerror.Wrap(err, "删除旧资产明细失败")
	}

	// 遍历余额列表，插入新的 asset_details 记录
	now := time.Now()
	for _, bal := range balances {
		_, err = assetDao.AssetDetails.Ctx(ctx).Data(g.Map{
			assetDao.AssetDetails.Columns().UserId:      consts.GetUserID(ctx),
			assetDao.AssetDetails.Columns().AssetSymbol: bal.Symbol,
			assetDao.AssetDetails.Columns().AssetName:   exchangeName,
			assetDao.AssetDetails.Columns().Balance:     bal.Total,
			assetDao.AssetDetails.Columns().PriceUsd:    bal.ValueUSD / bal.Total,
			assetDao.AssetDetails.Columns().ValueUsd:    bal.ValueUSD,
			assetDao.AssetDetails.Columns().SourceType:  "cex",
			assetDao.AssetDetails.Columns().LastUpdated: now,
			assetDao.AssetDetails.Columns().CreatedAt:   now,
			assetDao.AssetDetails.Columns().UpdatedAt:   now,
		}).Insert()
		if err != nil {
			g.Log().Error(ctx, "插入资产明细失败",
				"symbol", bal.Symbol,
				"exchange", exchangeName,
				"error", err,
			)
			continue
		}
	}

	g.Log().Info(ctx, "同步交易所账户余额完成",
		"accountId", accountID,
		"exchange", exchangeName,
		"balanceCount", len(balances),
	)

	return nil
}

// decryptCredentials 获取并解密交易所账户的 API 凭证
func (s *sExchange) decryptCredentials(ctx context.Context, accountID int) (apiKey, apiSecret, passphrase, exchangeName string, err error) {
	var account entity.ExchangeAccounts
	err = dao.ExchangeAccounts.Ctx(ctx).
		Where(dao.ExchangeAccounts.Columns().Id, accountID).
		WhereNull(dao.ExchangeAccounts.Columns().DeletedAt).
		Scan(&account)
	if err != nil {
		return "", "", "", "", gerror.Wrap(err, "查询账户信息失败")
	}
	if account.Id == 0 {
		return "", "", "", "", gerror.Newf("账户不存在: %d", accountID)
	}

	masterKey := g.Cfg().MustGet(ctx, "security.masterKey").String()

	apiKey, err = crypto.DecryptAES(account.ApiKeyEncrypted, masterKey)
	if err != nil {
		return "", "", "", "", gerror.Wrap(err, "解密 API Key 失败")
	}

	apiSecret, err = crypto.DecryptAES(account.ApiSecretEncrypted, masterKey)
	if err != nil {
		return "", "", "", "", gerror.Wrap(err, "解密 API Secret 失败")
	}

	if account.ApiPassphraseEncrypted != "" {
		passphrase, err = crypto.DecryptAES(account.ApiPassphraseEncrypted, masterKey)
		if err != nil {
			return "", "", "", "", gerror.Wrap(err, "解密 API Passphrase 失败")
		}
	}

	return apiKey, apiSecret, passphrase, account.ExchangeName, nil
}

// toAccountItem 将数据库实体转换为 API 响应格式
// 过滤所有加密字段，确保安全
func (s *sExchange) toAccountItem(acc *entity.ExchangeAccounts) exchangeApi.AccountItem {
	status := "inactive"
	if acc.IsActive == 1 {
		status = "active"
	}
	return exchangeApi.AccountItem{
		ID:           uint(acc.Id),
		ExchangeName: acc.ExchangeName,
		Label:        acc.Label,
		Note:         acc.Note,
		Status:       status,
		CreatedAt:    acc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    acc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
