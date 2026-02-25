package etherscan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"

	"your-finance/allfi/internal/utils"
)

const (
	// chainlistCacheFile Chainlist 原始 JSON 缓存文件
	chainlistCacheFile = "data/chainlist_rpcs.json"
	// verifiedRPCCacheFile 已验证的 RPC 映射缓存文件
	verifiedRPCCacheFile = "data/chainlist_verified_rpcs.json"
	// chainlistCacheTTL 缓存有效期：24 小时
	chainlistCacheTTL = 24 * time.Hour
)

// ChainlistRPC Chainlist 中单个 RPC 节点信息
type ChainlistRPC struct {
	URL      string `json:"url"`
	Tracking string `json:"tracking"`
}

// ChainlistEntry Chainlist 中单条链的数据
type ChainlistEntry struct {
	ChainId int            `json:"chainId"`
	RPCs    []ChainlistRPC `json:"rpc"`
}

// chainlistCacheData 缓存文件结构
type chainlistCacheData struct {
	FetchedAt time.Time `json:"fetched_at"`
	Data      json.RawMessage `json:"data"`
}

// verifiedRPCCache 已验证 RPC 缓存结构
type verifiedRPCCache struct {
	VerifiedAt time.Time      `json:"verified_at"`
	RPCs       map[int]string `json:"rpcs"`
}

var (
	dynamicRPCMap = make(map[int]string)
	rpcMapMutex   sync.RWMutex
	lastFetch     time.Time
	fetchMutex    sync.Mutex
	// initOnce 确保启动时只从文件加载一次
	initOnce sync.Once
)

// loadCachedRPCs 从本地文件加载已验证的 RPC 缓存
// 如果缓存存在且未过期，直接使用；否则返回 false
func loadCachedRPCs(ctx context.Context) bool {
	if !gfile.Exists(verifiedRPCCacheFile) {
		return false
	}

	data, err := os.ReadFile(verifiedRPCCacheFile)
	if err != nil {
		g.Log().Warning(ctx, "读取 RPC 缓存文件失败: ", err)
		return false
	}

	var cache verifiedRPCCache
	if err := json.Unmarshal(data, &cache); err != nil {
		g.Log().Warning(ctx, "解析 RPC 缓存文件失败: ", err)
		return false
	}

	// 检查缓存是否过期
	if time.Since(cache.VerifiedAt) > chainlistCacheTTL {
		return false
	}

	rpcMapMutex.Lock()
	dynamicRPCMap = cache.RPCs
	lastFetch = cache.VerifiedAt
	rpcMapMutex.Unlock()

	g.Log().Info(ctx, "从本地缓存加载 RPC 列表成功，缓存时间: ", cache.VerifiedAt.Format(time.RFC3339))
	return true
}

// saveVerifiedRPCs 将已验证的 RPC 映射保存到本地文件
func saveVerifiedRPCs(ctx context.Context, rpcs map[int]string) {
	dir := filepath.Dir(verifiedRPCCacheFile)
	if !gfile.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			g.Log().Warning(ctx, "创建缓存目录失败: ", err)
			return
		}
	}

	cache := verifiedRPCCache{
		VerifiedAt: time.Now(),
		RPCs:       rpcs,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		g.Log().Warning(ctx, "序列化 RPC 缓存失败: ", err)
		return
	}

	if err := os.WriteFile(verifiedRPCCacheFile, data, 0644); err != nil {
		g.Log().Warning(ctx, "写入 RPC 缓存文件失败: ", err)
	}
}

// loadChainlistFromFile 从本地缓存文件加载 Chainlist 原始数据
// 如果缓存存在且未过期，返回数据；否则返回 nil
func loadChainlistFromFile(ctx context.Context) []ChainlistEntry {
	if !gfile.Exists(chainlistCacheFile) {
		return nil
	}

	data, err := os.ReadFile(chainlistCacheFile)
	if err != nil {
		g.Log().Warning(ctx, "读取 Chainlist 缓存文件失败: ", err)
		return nil
	}

	var cache chainlistCacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		g.Log().Warning(ctx, "解析 Chainlist 缓存文件失败: ", err)
		return nil
	}

	// 检查缓存是否过期
	if time.Since(cache.FetchedAt) > chainlistCacheTTL {
		g.Log().Info(ctx, "Chainlist 缓存已过期，需要重新获取")
		return nil
	}

	var entries []ChainlistEntry
	if err := json.Unmarshal(cache.Data, &entries); err != nil {
		g.Log().Warning(ctx, "解析 Chainlist 缓存数据失败: ", err)
		return nil
	}

	g.Log().Info(ctx, "从本地缓存加载 Chainlist 数据成功")
	return entries
}

// saveChainlistToFile 将 Chainlist 原始数据保存到本地文件
func saveChainlistToFile(ctx context.Context, rawData []byte) {
	dir := filepath.Dir(chainlistCacheFile)
	if !gfile.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			g.Log().Warning(ctx, "创建缓存目录失败: ", err)
			return
		}
	}

	cache := chainlistCacheData{
		FetchedAt: time.Now(),
		Data:      rawData,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		g.Log().Warning(ctx, "序列化 Chainlist 缓存失败: ", err)
		return
	}

	if err := os.WriteFile(chainlistCacheFile, data, 0644); err != nil {
		g.Log().Warning(ctx, "写入 Chainlist 缓存文件失败: ", err)
	}
}

// fetchChainlistFromAPI 从 Chainlist API 获取数据
func fetchChainlistFromAPI(ctx context.Context) ([]ChainlistEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://chainlist.org/rpcs.json", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取原始数据用于缓存
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 保存原始数据到本地文件
	saveChainlistToFile(ctx, rawData)

	var entries []ChainlistEntry
	if err := json.Unmarshal(rawData, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// verifyRPC 测试 RPC 节点是否可用
func verifyRPC(ctx context.Context, rpcUrl string) bool {
	reqBody := []byte(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}`)

	testCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	testReq, err := http.NewRequestWithContext(testCtx, http.MethodPost, rpcUrl, bytes.NewReader(reqBody))
	if err != nil {
		return false
	}
	testReq.Header.Set("Content-Type", "application/json")

	testClient := &http.Client{Timeout: 3 * time.Second}
	testResp, err := testClient.Do(testReq)
	if err != nil {
		return false
	}
	defer testResp.Body.Close()

	if testResp.StatusCode != http.StatusOK {
		return false
	}

	var rpcResp struct {
		Result interface{} `json:"result"`
		Error  interface{} `json:"error"`
	}
	if err := json.NewDecoder(testResp.Body).Decode(&rpcResp); err != nil {
		return false
	}

	return rpcResp.Error == nil && rpcResp.Result != nil
}

// filterAndVerifyRPCs 从 Chainlist 数据中筛选并验证 RPC 节点
func filterAndVerifyRPCs(ctx context.Context, entries []ChainlistEntry) map[int]string {
	newMap := make(map[int]string)
	for _, entry := range entries {
		// 只检查我们支持的链
		isSupported := false
		for _, cfg := range SupportedChains {
			if cfg.ChainID == entry.ChainId {
				isSupported = true
				break
			}
		}
		if !isSupported {
			continue
		}

		var preferredRPCs []string
		var otherRPCs []string

		for _, rpc := range entry.RPCs {
			urlLower := strings.ToLower(rpc.URL)
			if strings.HasPrefix(urlLower, "https://") &&
				!strings.Contains(urlLower, "api_key") &&
				!strings.Contains(urlLower, "${") &&
				!strings.Contains(urlLower, "ankr.com") &&
				!strings.Contains(urlLower, "alchemy.com") &&
				!strings.Contains(urlLower, "infura.io") &&
				!strings.Contains(urlLower, "blastapi.io") &&
				!strings.Contains(urlLower, "nodereal.io") &&
				!strings.Contains(urlLower, "getblock.io") &&
				!strings.Contains(urlLower, "tenderly.co") &&
				!strings.Contains(urlLower, "quiknode.pro") &&
				!strings.Contains(urlLower, "rpcfast.com") &&
				!strings.Contains(urlLower, "gateway.fm") &&
				!strings.Contains(urlLower, "1rpc.io") &&
				!strings.Contains(urlLower, "lava.build") {

				if rpc.Tracking == "none" {
					preferredRPCs = append(preferredRPCs, rpc.URL)
				} else {
					otherRPCs = append(otherRPCs, rpc.URL)
				}
			}
		}

		// 验证节点并选择第一个可用的
		for _, url := range append(preferredRPCs, otherRPCs...) {
			if verifyRPC(ctx, url) {
				newMap[entry.ChainId] = url
				break
			}
		}
	}
	return newMap
}

// UpdateDynamicRPCs 从本地缓存或 Chainlist API 获取免费 RPC 列表
// 优先使用本地文件缓存，24 小时刷新一次，避免频繁请求导致 IP 被封
func UpdateDynamicRPCs(ctx context.Context) {
	fetchMutex.Lock()
	defer fetchMutex.Unlock()

	// 限制频率：24 小时更新一次
	rpcMapMutex.RLock()
	if time.Since(lastFetch) < chainlistCacheTTL {
		rpcMapMutex.RUnlock()
		return
	}
	rpcMapMutex.RUnlock()

	// 优先尝试从已验证的 RPC 缓存文件加载
	if loadCachedRPCs(ctx) {
		return
	}

	// 尝试从 Chainlist 原始数据缓存加载（避免重复请求 API）
	entries := loadChainlistFromFile(ctx)
	if entries == nil {
		// 缓存不存在或已过期，从 API 获取
		var err error
		entries, err = fetchChainlistFromAPI(ctx)
		if err != nil {
			g.Log().Warning(ctx, "获取 Chainlist RPC 失败: ", err)
			return
		}
	}

	// 筛选并验证 RPC 节点
	newMap := filterAndVerifyRPCs(ctx, entries)

	// 更新内存中的 RPC 映射
	rpcMapMutex.Lock()
	dynamicRPCMap = newMap
	lastFetch = time.Now()
	rpcMapMutex.Unlock()

	// 保存已验证的 RPC 到本地文件
	saveVerifiedRPCs(ctx, newMap)

	g.Log().Info(ctx, "免费 RPC 列表更新完成，共 ", len(newMap), " 条")
}

// GetRPCURL 获取指定链的有效 RPC URL
func GetRPCURL(ctx context.Context, chainName string) string {
	// 1. 获取用户自定义的 RPC
	provider := chainName + "_rpc"
	customRPC := utils.ResolveAPIKey(ctx, provider)
	if customRPC != "" {
		return customRPC
	}

	// 2. 启动时从本地缓存初始化（只执行一次）
	initOnce.Do(func() {
		loadCachedRPCs(ctx)
	})

	// 3. 异步触发更新（不阻塞主流程）
	go UpdateDynamicRPCs(context.Background())

	config, ok := SupportedChains[chainName]
	if ok && config.ChainID > 0 {
		rpcMapMutex.RLock()
		dynRPC, dynOk := dynamicRPCMap[config.ChainID]
		rpcMapMutex.RUnlock()
		if dynOk && dynRPC != "" {
			return dynRPC
		}
	}

	// 4. Fallback 到硬编码配置
	if ok {
		return config.PublicRPC
	}
	return ""
}
