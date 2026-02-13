// Package integrations 第三方 API 集成基础设施
// 提供 HTTP 客户端基类和通用工具
package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"your-finance/allfi/internal/utils"
)

// BaseClient HTTP 基础客户端
type BaseClient struct {
	httpClient *http.Client
	baseURL    string
	apiName    string // 用于限流
}

// NewBaseClient 创建基础客户端
func NewBaseClient(baseURL, apiName string, timeout time.Duration) *BaseClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &BaseClient{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		apiName:    apiName,
	}
}

// DoRequest 发送 HTTP 请求
func (c *BaseClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// 应用限流
	if err := utils.WaitForAPI(ctx, c.apiName); err != nil {
		return nil, fmt.Errorf("限流等待被中断: %v", err)
	}

	// 设置请求上下文
	req = req.WithContext(ctx)

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DoGet 发送 GET 请求
func (c *BaseClient) DoGet(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.DoRequest(ctx, req)
}

// ReadJSON 读取 JSON 响应
func ReadJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API 错误 [%d]: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("解析 JSON 失败: %v", err)
	}

	return nil
}

// Balance 通用余额结构
type Balance struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Free      float64 `json:"free"`                  // 可用余额
	Locked    float64 `json:"locked"`                // 冻结余额
	Total     float64 `json:"total"`                 // 总余额
	ValueUSD  float64 `json:"value_usd"`             // USD 价值
	Protocol  string  `json:"protocol,omitempty"`    // DeFi 协议名称（lido, rocketpool, aave）
	AssetType string  `json:"asset_type,omitempty"`  // 资产类型（spot, futures, margin, staking, lending, lp）
}

// ExchangeClient 交易所客户端接口
type ExchangeClient interface {
	GetBalances(ctx context.Context) ([]Balance, error)
	TestConnection(ctx context.Context) error
	GetName() string
	// 交易历史
	GetTradeHistory(ctx context.Context, params TradeHistoryParams) ([]Trade, error)
	// 充值历史
	GetDepositHistory(ctx context.Context, params DepositWithdrawParams) ([]Transfer, error)
	// 提现历史
	GetWithdrawHistory(ctx context.Context, params DepositWithdrawParams) ([]Transfer, error)
}

// Trade 交易记录
type Trade struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`      // 交易对（如 BTCUSDT）
	Side      string    `json:"side"`        // buy/sell
	Price     float64   `json:"price"`       // 成交价
	Quantity  float64   `json:"quantity"`    // 成交量
	Fee       float64   `json:"fee"`         // 手续费
	FeeCoin   string    `json:"fee_coin"`    // 手续费币种
	Timestamp time.Time `json:"timestamp"`   // 成交时间
	Source    string    `json:"source"`      // 来源交易所
}

// Transfer 充提记录
type Transfer struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`        // deposit/withdraw
	Coin      string    `json:"coin"`        // 币种
	Amount    float64   `json:"amount"`      // 数量
	Fee       float64   `json:"fee"`         // 手续费
	Address   string    `json:"address"`     // 地址
	TxHash    string    `json:"tx_hash"`     // 交易哈希
	Status    string    `json:"status"`      // 状态
	Timestamp time.Time `json:"timestamp"`   // 时间
	Source    string    `json:"source"`      // 来源交易所
}

// TradeHistoryParams 交易历史查询参数
type TradeHistoryParams struct {
	Symbol    string    `json:"symbol"`      // 交易对（可选）
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间
	Limit     int       `json:"limit"`       // 数量限制
}

// DepositWithdrawParams 充提历史查询参数
type DepositWithdrawParams struct {
	Coin      string    `json:"coin"`        // 币种（可选）
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间
	Limit     int       `json:"limit"`       // 数量限制
}

// BlockchainClient 区块链客户端接口
type BlockchainClient interface {
	GetNativeBalance(ctx context.Context, address string) (float64, error)
	GetTokenBalances(ctx context.Context, address string) ([]Balance, error)
	GetChainName() string
}

// PriceClient 价格服务客户端接口
type PriceClient interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
	GetPrices(ctx context.Context, symbols []string) (map[string]float64, error)
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
}
