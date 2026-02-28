package bridge

import "time"

// BridgeProtocol 跨链桥协议类型
type BridgeProtocol string

const (
	ProtocolStargate BridgeProtocol = "Stargate"
	ProtocolSynapse  BridgeProtocol = "Synapse"
	ProtocolHop      BridgeProtocol = "Hop"
	ProtocolAcross   BridgeProtocol = "Across"
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusConfirmed TransactionStatus = "confirmed"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
)

// BridgeTransaction 跨链桥交易信息
type BridgeTransaction struct {
	TxHash         string            `json:"tx_hash"`
	BridgeProtocol BridgeProtocol    `json:"bridge_protocol"`
	SourceChain    string            `json:"source_chain"`
	SourceToken    string            `json:"source_token"`
	SourceAmount   float64           `json:"source_amount"`
	DestChain      string            `json:"dest_chain"`
	DestToken      string            `json:"dest_token"`
	DestAmount     float64           `json:"dest_amount"`
	BridgeFee      float64           `json:"bridge_fee"`
	GasFee         float64           `json:"gas_fee"`
	TotalFeeUsd    float64           `json:"total_fee_usd"`
	Status         TransactionStatus `json:"status"`
	InitiatedAt    time.Time         `json:"initiated_at"`
	CompletedAt    *time.Time        `json:"completed_at"`
}

// BridgeInfo 跨链桥信息
type BridgeInfo struct {
	Protocol      BridgeProtocol `json:"protocol"`
	Name          string         `json:"name"`
	SupportChains []string       `json:"support_chains"`
	FeeRate       float64        `json:"fee_rate"` // 费率百分比
	AvgTime       int            `json:"avg_time"` // 平均完成时间（秒）
}

// BridgeClient 跨链桥客户端接口
type BridgeClient interface {
	// GetTransaction 获取交易信息
	GetTransaction(txHash string) (*BridgeTransaction, error)

	// GetSupportedChains 获取支持的链列表
	GetSupportedChains() []string

	// GetBridgeInfo 获取跨链桥信息
	GetBridgeInfo() *BridgeInfo
}
