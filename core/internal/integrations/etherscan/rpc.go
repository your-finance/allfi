// Package etherscan 公共 RPC Gas 价格查询
// 当 Etherscan API Key 未配置时，通过公共 RPC 节点获取 Gas 价格
package etherscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"

	"your-finance/allfi/internal/utils"
)

// rpcRequest JSON-RPC 请求结构
type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// rpcResponse JSON-RPC 响应结构
type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error,omitempty"`
}

// rpcError JSON-RPC 错误
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UnmarshalJSON 支持字符串和标准的对象结构，兼容不规范的 RPC 节点返回
func (e *rpcError) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		e.Code = -1
		e.Message = s
		return nil
	}

	type Alias rpcError
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*e = rpcError(aux)
	return nil
}

// GetGasPriceViaRPC 通过公共 RPC 节点获取 Gas 价格
// 调用 eth_gasPrice 方法，返回与 Etherscan Gas Oracle 兼容的 GasPrice 结构
// 由于 RPC 只返回单一建议价格，Safe/Normal/Fast 按比例估算
func GetGasPriceViaRPC(ctx context.Context, chainName string) (*GasPrice, error) {
	config, ok := SupportedChains[chainName]
	if !ok {
		return nil, fmt.Errorf("不支持的链: %s", chainName)
	}
	if config.PublicRPC == "" {
		return nil, fmt.Errorf("链 %s 未配置公共 RPC 端点", chainName)
	}

	// 限流
	if err := utils.WaitForAPI(ctx, config.RateLimitKey+"_rpc"); err != nil {
		return nil, err
	}

	// 构造 JSON-RPC 请求
	reqBody, _ := json.Marshal(rpcRequest{
		JSONRPC: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      1,
	})

	httpClient := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.PublicRPC, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("构造 RPC 请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RPC 请求失败: %v", err)
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("解析 RPC 响应失败: %v", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC 错误: %s", rpcResp.Error.Message)
	}

	// 解析十六进制 Gas 价格（单位: Wei）
	var hexPrice string
	if err := json.Unmarshal(rpcResp.Result, &hexPrice); err != nil {
		return nil, fmt.Errorf("解析 Gas 价格失败: %v", err)
	}

	gasPriceWei, ok := new(big.Int).SetString(hexPrice[2:], 16) // 去掉 0x 前缀
	if !ok {
		return nil, fmt.Errorf("解析十六进制 Gas 价格失败: %s", hexPrice)
	}

	// Wei → Gwei（除以 1e9）
	gasPriceGwei := new(big.Float).Quo(
		new(big.Float).SetInt(gasPriceWei),
		new(big.Float).SetFloat64(1e9),
	)
	normalGwei, _ := gasPriceGwei.Float64()

	// RPC 只返回单一建议价格，按比例估算 Safe/Fast
	return &GasPrice{
		Safe:    math.Max(normalGwei*0.8, 0.1),
		Normal:  normalGwei,
		Fast:    normalGwei * 1.3,
		BaseFee: 0, // RPC 不直接返回 BaseFee
	}, nil
}
