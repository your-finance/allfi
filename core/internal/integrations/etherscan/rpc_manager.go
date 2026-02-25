package etherscan

import (
	"bytes"
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

	// 用于测试 RPC 是否可用的内部函数
	verifyRPC := func(rpcUrl string) bool {
		// 构造极简的 JSON-RPC 请求
		reqBody := []byte(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}`)

		// 测试请求通常需要较短超时，避免拖慢整体速度
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

		// 如果包含 result 且不包含 error，则认为有效
		return rpcResp.Error == nil && rpcResp.Result != nil
	}

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
			// 简单筛选没有特殊字符、不含已知需鉴权的商业节点名的 HTTPs URL
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
			if verifyRPC(url) {
				newMap[entry.ChainId] = url
				break
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
