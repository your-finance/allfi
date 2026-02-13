// Package defi 协议注册中心
// 管理所有已注册的 DeFi 协议，提供聚合查询
package defi

import (
	"context"
	"fmt"
	"sync"
)

// Registry DeFi 协议注册中心
// 集中管理所有 DeFi 协议，支持按地址聚合查询所有仓位
type Registry struct {
	mu        sync.RWMutex
	protocols map[string]DeFiProtocol
}

// NewRegistry 创建协议注册中心
func NewRegistry() *Registry {
	return &Registry{
		protocols: make(map[string]DeFiProtocol),
	}
}

// Register 注册 DeFi 协议
func (r *Registry) Register(protocol DeFiProtocol) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.protocols[protocol.GetName()] = protocol
}

// GetProtocol 获取指定协议
func (r *Registry) GetProtocol(name string) (DeFiProtocol, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.protocols[name]
	if !ok {
		return nil, fmt.Errorf("未注册的 DeFi 协议: %s", name)
	}
	return p, nil
}

// ListProtocols 列出所有已注册的协议信息
func (r *Registry) ListProtocols() []ProtocolInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]ProtocolInfo, 0, len(r.protocols))
	for _, p := range r.protocols {
		infos = append(infos, ProtocolInfo{
			Name:        p.GetName(),
			DisplayName: p.GetDisplayName(),
			Type:        p.GetType(),
			Chains:      p.SupportedChains(),
		})
	}
	return infos
}

// ProtocolInfo 协议基本信息（用于 API 展示）
type ProtocolInfo struct {
	Name        string   `json:"name"`         // 协议标识
	DisplayName string   `json:"display_name"` // 显示名
	Type        string   `json:"type"`         // 协议类型
	Chains      []string `json:"chains"`       // 支持的链
}

// GetAllPositions 聚合查询用户在所有协议中的仓位
// 并发查询所有协议，收集全部仓位
func (r *Registry) GetAllPositions(ctx context.Context, address string, chain string) ([]Position, error) {
	r.mu.RLock()
	protocols := make([]DeFiProtocol, 0, len(r.protocols))
	for _, p := range r.protocols {
		protocols = append(protocols, p)
	}
	r.mu.RUnlock()

	type result struct {
		positions []Position
		err       error
	}

	ch := make(chan result, len(protocols))
	var wg sync.WaitGroup

	for _, p := range protocols {
		// 检查协议是否支持指定的链
		if chain != "" && !supportsChain(p, chain) {
			continue
		}
		wg.Add(1)
		go func(protocol DeFiProtocol) {
			defer wg.Done()
			targetChain := chain
			if targetChain == "" {
				targetChain = "ethereum" // 默认查询 Ethereum
			}
			positions, err := protocol.GetPositions(ctx, address, targetChain)
			ch <- result{positions: positions, err: err}
		}(p)
	}

	// 等待所有协程完成后关闭通道
	go func() {
		wg.Wait()
		close(ch)
	}()

	var allPositions []Position
	for r := range ch {
		if r.err != nil {
			// 单个协议失败不影响其他协议，跳过
			continue
		}
		allPositions = append(allPositions, r.positions...)
	}

	return allPositions, nil
}

// GetPositionsByProtocol 查询用户在指定协议中的仓位
func (r *Registry) GetPositionsByProtocol(ctx context.Context, address string, chain string, protocolName string) ([]Position, error) {
	protocol, err := r.GetProtocol(protocolName)
	if err != nil {
		return nil, err
	}
	return protocol.GetPositions(ctx, address, chain)
}

// supportsChain 检查协议是否支持指定的链
func supportsChain(p DeFiProtocol, chain string) bool {
	for _, c := range p.SupportedChains() {
		if c == chain {
			return true
		}
	}
	return false
}
