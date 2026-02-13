// Package logic 市场数据业务逻辑
// 获取多链 Gas 价格：
// - Ethereum: 通过 Etherscan API 实时查询（带 15 秒缓存）
// - BSC: 通过 BscScan API 实时查询（带 15 秒缓存）
// - Polygon: 通过 PolygonScan API 实时查询（带 15 秒缓存）
package logic

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/app/market/model"
	"your-finance/allfi/internal/app/market/service"
	"your-finance/allfi/internal/integrations/etherscan"
)

// chainGasCache 单条链的 Gas 价格缓存
type chainGasCache struct {
	data *etherscan.GasPrice
	time time.Time
}

// sMarket 市场数据服务实现
type sMarket struct {
	// Gas 价格缓存（避免频繁调用 API），按链名索引
	gasCacheMu sync.RWMutex
	gasCache   map[string]*chainGasCache
	gasCacheTTL time.Duration
}

// New 创建市场数据服务实例
func New() service.IMarket {
	return &sMarket{
		gasCache:    make(map[string]*chainGasCache),
		gasCacheTTL: 15 * time.Second, // 15 秒缓存 TTL
	}
}

// GetGasPrice 获取多链 Gas 价格
//
// 功能说明:
// 1. 查询 ETH Gas 价格（Etherscan API + 15 秒缓存）
// 2. 查询 BSC Gas 价格（BscScan API + 15 秒缓存）
// 3. 查询 Polygon Gas 价格（PolygonScan API + 15 秒缓存）
// 4. 判断各链拥堵等级
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - *model.MultiChainGasResponse: 多链 Gas 价格数据
//   - error: 错误信息
func (s *sMarket) GetGasPrice(ctx context.Context) (*model.MultiChainGasResponse, error) {
	// 获取 ETH Gas 价格（带缓存）
	ethGas := s.getChainGasPrice(ctx, "ethereum", "ETHERSCAN_API_KEY")

	// 获取 BSC Gas 价格（带缓存），失败时使用硬编码兜底
	bscGas := s.getChainGasPrice(ctx, "bsc", "BSCSCAN_API_KEY")
	if bscGas.Normal == 0 {
		// API 未返回有效数据，使用硬编码兜底值
		bscGas = &etherscan.GasPrice{Safe: 1.0, Normal: 3.0, Fast: 5.0, BaseFee: 0}
	}

	// 获取 Polygon Gas 价格（带缓存），失败时使用硬编码兜底
	polygonGas := s.getChainGasPrice(ctx, "polygon", "POLYGONSCAN_API_KEY")
	if polygonGas.Normal == 0 {
		// API 未返回有效数据，使用硬编码兜底值
		polygonGas = &etherscan.GasPrice{Safe: 25.0, Normal: 30.0, Fast: 50.0, BaseFee: 0}
	}

	// 构建各链 Gas 价格列表
	chains := []model.ChainGasPrice{
		{
			Chain:    "ethereum",
			Name:     "Ethereum",
			Low:      ethGas.Safe,
			Standard: ethGas.Normal,
			Fast:     ethGas.Fast,
			Instant:  ethGas.Fast * 1.2, // 极速约为快速的 1.2 倍
			BaseFee:  ethGas.BaseFee,
			Unit:     "Gwei",
			Level:    model.DetermineLevel(ethGas.Normal),
		},
		{
			Chain:    "bsc",
			Name:     "BSC",
			Low:      bscGas.Safe,
			Standard: bscGas.Normal,
			Fast:     bscGas.Fast,
			Instant:  bscGas.Fast * 1.2,
			BaseFee:  bscGas.BaseFee,
			Unit:     "Gwei",
			Level:    model.DetermineLevel(bscGas.Normal),
		},
		{
			Chain:    "polygon",
			Name:     "Polygon",
			Low:      polygonGas.Safe,
			Standard: polygonGas.Normal,
			Fast:     polygonGas.Fast,
			Instant:  polygonGas.Fast * 1.2,
			BaseFee:  polygonGas.BaseFee,
			Unit:     "Gwei",
			Level:    model.DetermineLevel(polygonGas.Normal),
		},
	}

	return &model.MultiChainGasResponse{
		Chains:    chains,
		UpdatedAt: time.Now().Unix(),
	}, nil
}

// getChainGasPrice 获取指定链的 Gas 价格（带 15 秒缓存）
//
// 通用方法，支持 Ethereum / BSC / Polygon 等所有 EVM 兼容链。
// 通过 etherscan.NewChainClient 创建对应链的客户端，复用相同的
// Gas Oracle API 格式（各 *Scan 平台接口一致）。
//
// 缓存策略:
// 1. 先检查该链的缓存是否有效
// 2. 缓存过期则调用对应链的 *Scan API
// 3. API 失败时返回过期缓存
// 4. 都没有时返回零值
//
// 参数:
//   - chainName: 链标识（ethereum / bsc / polygon）
//   - envKey: API Key 对应的环境变量名
func (s *sMarket) getChainGasPrice(ctx context.Context, chainName string, envKey string) *etherscan.GasPrice {
	// 先检查缓存（读锁）
	s.gasCacheMu.RLock()
	if cached, ok := s.gasCache[chainName]; ok && cached.data != nil && time.Since(cached.time) < s.gasCacheTTL {
		data := cached.data
		s.gasCacheMu.RUnlock()
		return data
	}
	s.gasCacheMu.RUnlock()

	// 缓存过期，创建对应链的客户端并调用 API
	apiKey := os.Getenv(envKey)
	if apiKey == "" {
		g.Log().Warning(ctx, "API Key 未配置，Gas 价格将返回零值",
			"chain", chainName,
			"envKey", envKey,
		)
		return &etherscan.GasPrice{}
	}

	client, err := etherscan.NewChainClient(chainName, apiKey)
	if err != nil {
		g.Log().Warning(ctx, "创建链客户端失败",
			"chain", chainName,
			"error", err,
		)
		return &etherscan.GasPrice{}
	}

	gasPrice, err := client.GetGasPrice(ctx)
	if err != nil {
		g.Log().Warning(ctx, "获取 Gas 价格失败",
			"chain", chainName,
			"error", err,
		)

		// API 失败时返回过期缓存
		s.gasCacheMu.RLock()
		if cached, ok := s.gasCache[chainName]; ok && cached.data != nil {
			data := cached.data
			s.gasCacheMu.RUnlock()
			return data
		}
		s.gasCacheMu.RUnlock()

		// 都没有就返回零值
		return &etherscan.GasPrice{}
	}

	// 更新缓存（写锁）
	s.gasCacheMu.Lock()
	s.gasCache[chainName] = &chainGasCache{
		data: gasPrice,
		time: time.Now(),
	}
	s.gasCacheMu.Unlock()

	return gasPrice
}
