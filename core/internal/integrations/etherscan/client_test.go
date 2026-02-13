package etherscan

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewChainClient_SupportedChains 测试所有支持的链都能正确创建客户端
func TestNewChainClient_SupportedChains(t *testing.T) {
	chains := []string{"ethereum", "bsc", "arbitrum", "optimism", "polygon", "base"}
	for _, chain := range chains {
		c, err := NewChainClient(chain, "test-key")
		if err != nil {
			t.Errorf("创建 %s 客户端失败: %v", chain, err)
			continue
		}
		if c.GetChainName() != chain {
			t.Errorf("链名称不匹配: 期望 %s, 实际 %s", chain, c.GetChainName())
		}
		// 验证 baseURL 与配置一致
		expected := SupportedChains[chain].BaseURL
		if c.baseURL != expected {
			t.Errorf("%s baseURL 不匹配: 期望 %s, 实际 %s", chain, expected, c.baseURL)
		}
	}
}

// TestNewChainClient_UnsupportedChain 测试不支持的链返回错误
func TestNewChainClient_UnsupportedChain(t *testing.T) {
	_, err := NewChainClient("solana", "test-key")
	if err == nil {
		t.Error("期望返回错误，但未返回")
	}
}

// TestConvenienceFactories 测试便捷工厂函数
func TestConvenienceFactories(t *testing.T) {
	tests := []struct {
		name      string
		factory   func(string) *Client
		chainName string
	}{
		{"Arbiscan", NewArbiscanClient, "arbitrum"},
		{"Optimism", NewOptimismClient, "optimism"},
		{"Polygonscan", NewPolygonscanClient, "polygon"},
		{"Basescan", NewBasescanClient, "base"},
		{"Etherscan", NewClient, "ethereum"},
		{"BscScan", NewBscScanClient, "bsc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.factory("test-key")
			if c == nil {
				t.Fatal("客户端为 nil")
			}
			if c.GetChainName() != tt.chainName {
				t.Errorf("链名称不匹配: 期望 %s, 实际 %s", tt.chainName, c.GetChainName())
			}
		})
	}
}

// TestGetRateLimitKey 测试限流键名获取
func TestGetRateLimitKey(t *testing.T) {
	tests := []struct {
		chain    string
		expected string
	}{
		{"ethereum", "etherscan"},
		{"bsc", "bscscan"},
		{"arbitrum", "arbiscan"},
		{"optimism", "optimism"},
		{"polygon", "polygonscan"},
		{"base", "basescan"},
	}

	for _, tt := range tests {
		c, _ := NewChainClient(tt.chain, "test-key")
		if got := c.getRateLimitKey(); got != tt.expected {
			t.Errorf("%s 限流键名不匹配: 期望 %s, 实际 %s", tt.chain, tt.expected, got)
		}
	}
}

// TestSupportedChains_Completeness 验证所有链配置完整性
func TestSupportedChains_Completeness(t *testing.T) {
	for name, config := range SupportedChains {
		if config.Name == "" {
			t.Errorf("链 %s 的 Name 为空", name)
		}
		if config.BaseURL == "" {
			t.Errorf("链 %s 的 BaseURL 为空", name)
		}
		if config.NativeSymbol == "" {
			t.Errorf("链 %s 的 NativeSymbol 为空", name)
		}
		if config.RateLimitKey == "" {
			t.Errorf("链 %s 的 RateLimitKey 为空", name)
		}
		if config.Name != name {
			t.Errorf("链 %s 的 Name 字段与 map key 不一致: %s", name, config.Name)
		}
	}

	// 确保有 6 条链
	if len(SupportedChains) != 6 {
		t.Errorf("期望支持 6 条链, 实际 %d 条", len(SupportedChains))
	}
}

// TestGetNativeBalance_MockServer 使用 mock server 测试原生余额获取
func TestGetNativeBalance_MockServer(t *testing.T) {
	// 创建 mock server 返回 1.5 ETH（1500000000000000000 wei）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"status":  "1",
			"message": "OK",
			"result":  "1500000000000000000",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端，替换 baseURL 为 mock server
	client := &Client{
		apiKey:     "test-key",
		baseURL:    server.URL,
		chainName:  "ethereum",
		httpClient: server.Client(),
	}

	balance, err := client.GetNativeBalance(context.Background(), "0x742d35Cc6634C0532925a3b844Bc9e7595f2bD50")
	if err != nil {
		t.Fatalf("获取余额失败: %v", err)
	}

	// 1500000000000000000 wei = 1.5 ETH
	if balance < 1.49 || balance > 1.51 {
		t.Errorf("余额不正确: 期望约 1.5, 实际 %f", balance)
	}
}

// TestLookupDeFiToken_StETH 验证 stETH 被识别为 Lido staking
func TestLookupDeFiToken_StETH(t *testing.T) {
	info, found := LookupDeFiToken("0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84")
	if !found {
		t.Fatal("stETH 未被识别")
	}
	if info.Protocol != "lido" {
		t.Errorf("协议不正确: 期望 lido, 实际 %s", info.Protocol)
	}
	if info.AssetType != "staking" {
		t.Errorf("资产类型不正确: 期望 staking, 实际 %s", info.AssetType)
	}
}

// TestLookupDeFiToken_RETH 验证 rETH 被识别为 Rocket Pool staking
func TestLookupDeFiToken_RETH(t *testing.T) {
	info, found := LookupDeFiToken("0xae78736Cd615f374D3085123A210448E74Fc6393")
	if !found {
		t.Fatal("rETH 未被识别")
	}
	if info.Protocol != "rocketpool" {
		t.Errorf("协议不正确: 期望 rocketpool, 实际 %s", info.Protocol)
	}
	if info.AssetType != "staking" {
		t.Errorf("资产类型不正确: 期望 staking, 实际 %s", info.AssetType)
	}
}

// TestLookupDeFiToken_Unknown 验证普通 Token 不被标记为 DeFi
func TestLookupDeFiToken_Unknown(t *testing.T) {
	// USDT 不在 KnownDeFiTokens 中
	_, found := LookupDeFiToken("0xdac17f958d2ee523a2206206994597c13d831ec7")
	if found {
		t.Error("USDT 不应该被识别为 DeFi Token")
	}
}

// TestLookupDeFiToken_CaseInsensitive 验证地址匹配大小写不敏感
func TestLookupDeFiToken_CaseInsensitive(t *testing.T) {
	// 使用全大写地址（去掉 0x 前缀后）
	info, found := LookupDeFiToken("0xAE7AB96520DE3A18E5E111B5EAAB095312D7FE84")
	if !found {
		t.Fatal("大写地址未被识别")
	}
	if info.Protocol != "lido" {
		t.Errorf("协议不正确: 期望 lido, 实际 %s", info.Protocol)
	}
}

// TestKnownDeFiTokens_AllHaveRequiredFields 验证所有映射条目字段完整
func TestKnownDeFiTokens_AllHaveRequiredFields(t *testing.T) {
	for addr, info := range KnownDeFiTokens {
		if info.Protocol == "" {
			t.Errorf("地址 %s 的 Protocol 为空", addr)
		}
		if info.AssetType == "" {
			t.Errorf("地址 %s 的 AssetType 为空", addr)
		}
		if info.DisplayName == "" {
			t.Errorf("地址 %s 的 DisplayName 为空", addr)
		}
	}
}
