package logic

import (
	"context"
	"fmt"
	"time"

	"your-finance/allfi/internal/dao"
	"your-finance/allfi/internal/integrations/bridge"
	"your-finance/allfi/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ICrossChain 跨链交易业务逻辑接口
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
	TotalFee       float64            `json:"total_fee"`        // 总手续费
	AvgFee         float64            `json:"avg_fee"`          // 平均手续费
	FeeByBridge    map[string]float64 `json:"fee_by_bridge"`    // 按跨链桥统计
	FeeByChain     map[string]float64 `json:"fee_by_chain"`     // 按链统计
	TransactionCnt int                `json:"transaction_count"` // 交易次数
}

type sCrossChain struct {
	bridgeManager *bridge.BridgeManager
}

// GetTransactions 获取跨链交易列表
func (s *sCrossChain) GetTransactions(ctx context.Context, userId int64, page, pageSize int) ([]*entity.CrossChainTransaction, int, error) {
	// 查询总数
	total, err := dao.CrossChainTransaction.Ctx(ctx).
		Where("user_id", userId).
		Count()
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var transactions []*entity.CrossChainTransaction
	err = dao.CrossChainTransaction.Ctx(ctx).
		Where("user_id", userId).
		Order("initiated_at DESC").
		Page(page, pageSize).
		Scan(&transactions)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetAssetFlow 获取资产流向数据
func (s *sCrossChain) GetAssetFlow(ctx context.Context, userId int64, startTime, endTime *time.Time) (*AssetFlowData, error) {
	// 构建查询条件
	m := dao.CrossChainTransaction.Ctx(ctx).Where("user_id", userId)
	if startTime != nil {
		m = m.WhereGTE("initiated_at", startTime)
	}
	if endTime != nil {
		m = m.WhereLTE("initiated_at", endTime)
	}

	// 查询交易记录
	var transactions []*entity.CrossChainTransaction
	err := m.Scan(&transactions)
	if err != nil {
		return nil, err
	}

	// 统计节点和连接
	nodeMap := make(map[string]int)
	linkMap := make(map[string]float64)

	for _, tx := range transactions {
		// 统计节点（链）
		nodeMap[tx.SourceChain]++
		nodeMap[tx.DestChain]++

		// 统计连接（流向）
		linkKey := fmt.Sprintf("%s->%s", tx.SourceChain, tx.DestChain)
		linkMap[linkKey] += tx.SourceAmount
	}

	// 构建节点列表
	nodes := make([]FlowNode, 0, len(nodeMap))
	for name, value := range nodeMap {
		nodes = append(nodes, FlowNode{
			Name:  name,
			Value: value,
		})
	}

	// 构建连接列表
	links := make([]FlowLink, 0, len(linkMap))
	for key, value := range linkMap {
		// 解析 source 和 target
		var source, target string
		fmt.Sscanf(key, "%s->%s", &source, &target)
		links = append(links, FlowLink{
			Source: source,
			Target: target,
			Value:  value,
		})
	}

	return &AssetFlowData{
		Nodes: nodes,
		Links: links,
	}, nil
}

// GetFeeStats 获取跨链手续费统计
func (s *sCrossChain) GetFeeStats(ctx context.Context, userId int64, startTime, endTime *time.Time) (*FeeStatsData, error) {
	// 构建查询条件
	m := dao.CrossChainTransaction.Ctx(ctx).Where("user_id", userId)
	if startTime != nil {
		m = m.WhereGTE("initiated_at", startTime)
	}
	if endTime != nil {
		m = m.WhereLTE("initiated_at", endTime)
	}

	// 查询交易记录
	var transactions []*entity.CrossChainTransaction
	err := m.Scan(&transactions)
	if err != nil {
		return nil, err
	}

	// 统计数据
	stats := &FeeStatsData{
		FeeByBridge: make(map[string]float64),
		FeeByChain:  make(map[string]float64),
	}

	for _, tx := range transactions {
		stats.TotalFee += tx.TotalFeeUsd
		stats.TransactionCnt++

		// 按跨链桥统计
		stats.FeeByBridge[tx.BridgeProtocol] += tx.TotalFeeUsd

		// 按链统计
		stats.FeeByChain[tx.SourceChain] += tx.TotalFeeUsd
	}

	// 计算平均值
	if stats.TransactionCnt > 0 {
		stats.AvgFee = stats.TotalFee / float64(stats.TransactionCnt)
	}

	return stats, nil
}

// GetBridges 获取支持的跨链桥列表
func (s *sCrossChain) GetBridges(ctx context.Context) ([]*bridge.BridgeInfo, error) {
	return s.bridgeManager.GetAllBridges(), nil
}

// AddTransaction 添加跨链交易记录
func (s *sCrossChain) AddTransaction(ctx context.Context, userId int64, tx *bridge.BridgeTransaction) error {
	// 检查是否已存在
	count, err := dao.CrossChainTransaction.Ctx(ctx).
		Where("tx_hash", tx.TxHash).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("transaction already exists: %s", tx.TxHash)
	}

	// 插入记录
	entity := &entity.CrossChainTransaction{
		UserId:         userId,
		TxHash:         tx.TxHash,
		BridgeProtocol: string(tx.BridgeProtocol),
		SourceChain:    tx.SourceChain,
		SourceToken:    tx.SourceToken,
		SourceAmount:   tx.SourceAmount,
		DestChain:      tx.DestChain,
		DestToken:      tx.DestToken,
		DestAmount:     tx.DestAmount,
		BridgeFee:      tx.BridgeFee,
		GasFee:         tx.GasFee,
		TotalFeeUsd:    tx.TotalFeeUsd,
		Status:         string(tx.Status),
		InitiatedAt:    tx.InitiatedAt,
		CompletedAt:    tx.CompletedAt,
	}

	_, err = dao.CrossChainTransaction.Ctx(ctx).Data(entity).Insert()
	return err
}

// UpdateTransactionStatus 更新交易状态
func (s *sCrossChain) UpdateTransactionStatus(ctx context.Context, txHash string, status string) error {
	data := g.Map{
		"status":     status,
		"updated_at": gtime.Now(),
	}

	// 如果状态是完成，更新完成时间
	if status == string(bridge.StatusCompleted) {
		data["completed_at"] = gtime.Now()
	}

	_, err := dao.CrossChainTransaction.Ctx(ctx).
		Where("tx_hash", txHash).
		Data(data).
		Update()
	return err
}

// InitMockData 初始化模拟数据（用于测试）
func (s *sCrossChain) InitMockData(ctx context.Context, userId int64) error {
	// 检查是否已有数据
	count, err := dao.CrossChainTransaction.Ctx(ctx).
		Where("user_id", userId).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 已有数据，不再初始化
	}

	// 创建模拟交易
	mockTxs := []*bridge.BridgeTransaction{
		bridge.MockTransaction(bridge.ProtocolStargate, "ethereum", "bsc"),
		bridge.MockTransaction(bridge.ProtocolSynapse, "bsc", "polygon"),
		bridge.MockTransaction(bridge.ProtocolHop, "polygon", "arbitrum"),
		bridge.MockTransaction(bridge.ProtocolAcross, "arbitrum", "optimism"),
		bridge.MockTransaction(bridge.ProtocolStargate, "optimism", "ethereum"),
	}

	// 批量插入
	return dao.CrossChainTransaction.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, mockTx := range mockTxs {
			if err := s.AddTransaction(ctx, userId, mockTx); err != nil {
				return err
			}
		}
		return nil
	})
}
