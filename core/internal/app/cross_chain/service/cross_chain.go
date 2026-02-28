package service

import (
	"context"
	"time"

	"your-finance/allfi/internal/integrations/bridge"
	"your-finance/allfi/internal/model/entity"
)

// AssetFlowData 资产流向数据
type AssetFlowData struct {
	Nodes []FlowNode `json:"nodes"` // 节点列表（链）
	Links []FlowLink `json:"links"` // 连接列表（流向）
}

// FlowNode 流向节点
type FlowNode struct {
	Name  string `json:"name"`  // 节点名称（链名）
	Value int    `json:"value"` // 节点值（交易次数）
}

// FlowLink 流向连接
type FlowLink struct {
	Source string  `json:"source"` // 源节点
	Target string  `json:"target"` // 目标节点
	Value  float64 `json:"value"`  // 流向值（金额）
}

// FeeStatsData 手续费统计数据
type FeeStatsData struct {
	TotalFee       float64            `json:"total_fee"`         // 总手续费
	AvgFee         float64            `json:"avg_fee"`           // 平均手续费
	FeeByBridge    map[string]float64 `json:"fee_by_bridge"`     // 按跨链桥统计
	FeeByChain     map[string]float64 `json:"fee_by_chain"`      // 按链统计
	TransactionCnt int                `json:"transaction_count"` // 交易次数
}

// ICrossChain 跨链交易服务接口
type ICrossChain interface {
	// GetTransactions 获取跨链交易列表
	GetTransactions(ctx context.Context, userId int64, page, pageSize int) ([]*entity.CrossChainTransaction, int, error)

	// GetAssetFlow 获取资产流向数据
	GetAssetFlow(ctx context.Context, userId int64, startTime, endTime *time.Time) (*AssetFlowData, error)

	// GetFeeStats 获取跨链手续费统计
	GetFeeStats(ctx context.Context, userId int64, startTime, endTime *time.Time) (*FeeStatsData, error)

	// GetBridges 获取支持的跨链桥列表
	GetBridges(ctx context.Context) ([]*bridge.BridgeInfo, error)

	// AddTransaction 添加跨链交易记录
	AddTransaction(ctx context.Context, userId int64, tx *bridge.BridgeTransaction) error

	// UpdateTransactionStatus 更新交易状态
	UpdateTransactionStatus(ctx context.Context, txHash string, status string) error
}

var localCrossChain ICrossChain

// RegisterCrossChain 注册跨链交易服务
func RegisterCrossChain(i ICrossChain) {
	if localCrossChain == nil {
		localCrossChain = i
	}
}

// CrossChain 获取跨链交易服务实例
func CrossChain() ICrossChain {
	if localCrossChain == nil {
		panic("implement not found for interface ICrossChain, forgot register?")
	}
	return localCrossChain
}
