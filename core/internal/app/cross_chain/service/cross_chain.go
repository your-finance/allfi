package service

import (
	"context"
	"time"

	"your-finance/allfi/internal/app/cross_chain/logic"
	"your-finance/allfi/internal/integrations/bridge"
	"your-finance/allfi/internal/model/entity"
)

// ICrossChain 跨链交易服务接口
type ICrossChain interface {
	// GetTransactions 获取跨链交易列表
	GetTransactions(ctx context.Context, userId int64, page, pageSize int) ([]*entity.CrossChainTransaction, int, error)

	// GetAssetFlow 获取资产流向数据
	GetAssetFlow(ctx context.Context, userId int64, startTime, endTime *time.Time) (*logic.AssetFlowData, error)

	// GetFeeStats 获取跨链手续费统计
	GetFeeStats(ctx context.Context, userId int64, startTime, endTime *time.Time) (*logic.FeeStatsData, error)

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
