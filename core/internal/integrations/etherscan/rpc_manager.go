package etherscan

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"your-finance/allfi/internal/utils"
)

var (
	dynamicRPCMap = make(map[int]string)
	rpcMapMutex   sync.RWMutex
	lastFetch     time.Time
	fetchMutex    sync.Mutex
)

// UpdateDynamicRPCs 从 chainlist.org 获取免费 RPC 列表
func UpdateDynamicRPCs(ctx context.Context) {
	fetchMutex.Lock()
	defer fetchMutex.Unlock()

	// 限制频率, 1小时更新一次即可
	rpcMapMutex.RLock()
	if time.Since(lastFetch) < time.Hour {
		rpcMapMutex.RUnlock()
		return
	}
	rpcMapMutex.RUnlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://chainlist.org/rpcs.json", nil)
	if err != nil {
		g.Log().Warning(ctx, "构造 chainlist rpc 请求失败: ", err)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Warning(ctx, "获取 chainlist rpc 失败: ", err)
		return
	}
	defer resp.Body.Close()

	type ChainlistRPC struct {
		URL      string `json:"url"`
		Tracking string `json:"tracking"`
	}
	type ChainlistEntry struct {
		ChainId int            `json:"chainId"`
		RPCs    []ChainlistRPC `json:"rpc"`
	}

	var entries []ChainlistEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		g.Log().Warning(ctx, "解析 chainlist rpc 失败: ", err)
		return
	}

	newMap := make(map[int]string)
	for _, entry := range entries {
		for _, rpc := range entry.RPCs {
			// 简单筛选没有特殊字符且可访问的纯公共 URL
			if strings.HasPrefix(rpc.URL, "https://") &&
				!strings.Contains(rpc.URL, "API_KEY") &&
				!strings.Contains(rpc.URL, "${") {
				newMap[entry.ChainId] = rpc.URL
				break // 选第一个
			}
		}
	}

	rpcMapMutex.Lock()
	dynamicRPCMap = newMap
	lastFetch = time.Now()
	rpcMapMutex.Unlock()

	g.Log().Info(ctx, "免费 RPC 列表动态更新完成")
}

// GetRPCURL 获取指定链的有效 RPC URL
func GetRPCURL(ctx context.Context, chainName string) string {
	// 1. 获取用户自定义的 RPC
	provider := chainName + "_rpc"
	customRPC := utils.ResolveAPIKey(ctx, provider)
	if customRPC != "" {
		return customRPC
	}

	// 2. 动态获取免费 RPC
	// 异步触发更新以不阻塞主流程
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

	// 3. Fallback 到旧的硬编码配置
	if ok {
		return config.PublicRPC
	}
	return ""
}
